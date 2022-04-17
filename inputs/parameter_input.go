package inputs

import (
	"encoding/json"
	"fmt"
	"math"
	"sync"
)

type inputLine map[string]interface{}
type ParameterInput struct {
	rwlock      sync.RWMutex
	input_block map[string]inputLine
}

//----------------------------------------------------------------------------------------
//! \fn error ParameterInput.LoadFromByte(str []byte)
//! \brief Load input parameters from a string in bytes format
//!
//! If a parameters already exist, the value are replaced. If not, the parameter and the value
//! will be inserted into the original input block. Comments aren't allowed.

func (this *ParameterInput) LoadFromByte(str []byte) error {
	temp := make(map[string]inputLine)
	if err := json.Unmarshal(str, &temp); err != nil {
		return err
	}
	// Check whether the input value format is legal.
	for block_name, block := range temp {
		for para_name, para := range block {
			switch para.(type) {
			case bool, string, float64, int:
			default:
				return fmt.Errorf("Load Parameter Error: Parameter %s in block %s is illegal.", para_name, block_name)
			}
		}
	}
	// If no error exists, add the input parameter to this. Seperating the check and
	// write is to make sure changes will not applied until no error is found.
	this.rwlock.Lock()
	defer this.rwlock.Unlock()
	// If the input_block hasn't been initialized, just use the temp.
	if len(this.input_block) == 0 {
		this.input_block = temp
		return nil
	}
	for block_name, block := range temp {
		_, ok := this.input_block[block_name]
		if ok {
			for para_name, para := range block {
				this.input_block[block_name][para_name] = para
			}
		} else {
			this.input_block[block_name] = block
		}
	}
	return nil
}

//----------------------------------------------------------------------------------------
//! \fn (interface{}, error) ParameterInput.GetParameter(block string, name string)
//! \brief returns parameter value stored in block/name; return error if it does not exist

func (this *ParameterInput) GetParameter(block_name string, para_name string) (interface{}, error) {
	this.rwlock.RLock()
	defer this.rwlock.RUnlock()
	if input_line, ok := this.input_block[block_name]; ok {
		if para, ok := input_line[para_name]; ok {
			// Deepcopy for secure
			var result interface{}
			temp, _ := json.Marshal(para)
			json.Unmarshal(temp, result)
			return result, nil
		} else {
			return nil, fmt.Errorf("Get Parameter Error: %s isn't in block %s.", para_name, block_name)
		}
	}
	return nil, fmt.Errorf("Get Parameter Error: Block %s doesn't exist.", block_name)
}

//----------------------------------------------------------------------------------------
//! \fn ParameterInput.SetParameter(block string, name string, value interface{})
//! \brief updates a parameter; creates it if it does not exist

func (this *ParameterInput) SetParameter(block_name string, para_name string, value interface{}) {
	this.rwlock.Lock()
	defer this.rwlock.Unlock()
	if _, ok := this.input_block[block_name]; !ok {
		this.input_block[block_name] = make(inputLine)
	}
	this.input_block[block_name][para_name] = value
}

//----------------------------------------------------------------------------------------
//! \fn string ParameterInput.ParameterDump()
//! \brief output entire InputBlock/InputLine hierarchy to specified stream

func (this *ParameterInput) ParameterDump() string {
	this.rwlock.RLock()
	defer this.rwlock.RUnlock()
	if len(this.input_block) == 0 {return ""}
	result, _ := json.MarshalIndent(this.input_block, "", "    ")
	return string(result)
}

//----------------------------------------------------------------------------------------
//! \fn ParameterInput.ChangeNextTime(mesh_time float64)
//! \brief rollback next_time by dt for each output block if mesh_time = -1, or add dt to
//! next_time until next_time > mesh_time for each output block else
//!
//! If the user has added a new/fresh output round to multiple of dt, make sure that
//! mesh_time - dt0 < next_time <= mesh_time, to ensure immediate writing.
//! 
//! It's assumed that you have checked whether some parameters exsits.

func (this *ParameterInput) ChangeNextTime(mesh_time float64) {
	this.rwlock.Lock()
	defer this.rwlock.Unlock()
	for block_name, block := range this.input_block {
		// Slice of string means that everything should be in English.
		if len(block_name) > 5 && block_name[:6] == "output" {
			// dt must exist while next_time just need to exist if mesh_time = -1
			dt, ok := block["dt"]
			if !ok { panic("Get Parameter Error: \"dt\" isn't in block " + block_name) }
			dt_value, ok := dt.(float64)
			if !ok { panic("Parameter Type Error: \"dt\" in block " + block_name + " isn't an integer") }

			next_time, ok := block["next_time"]
			var next_time_value float64
			if !ok {
				if mesh_time == -1 {
					panic("Get Parameter Error: \"next_time\" isn't in block " + block_name)
				} else {
					this.input_block[block_name]["next_time"] = math.Floor(mesh_time/dt_value) * dt_value
					return
				}
			} else {
				next_time_value, ok = next_time.(float64)
				if !ok { panic("Parameter Type Error: \"next_time\" in block " + block_name + " isn't an integer") }
			}

			if mesh_time == -1 {
				this.input_block[block_name]["next_time"] = next_time_value - dt_value
			} else {
				dt0 := dt_value * (math.Floor((mesh_time-next_time_value)/dt_value) + 1)
				if dt0 > 0 {
					this.input_block[block_name]["next_time"] = next_time_value + dt0
				}
			}
		}
	}
	return
}
