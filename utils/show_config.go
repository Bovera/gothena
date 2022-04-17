package utils

import "fmt"

//----------------------------------------------------------------------------------------
//! \fn void ShowConfig()
//! \brief prints diagnostic messages about the configuration of an Athena++ executable

func ShowConfig() {
	// To match configure.py output: use 2 space indent for option, value output starts on
	// column 30
	fmt.Println("This Gothena executable is configured with:")
	fmt.Printf("  Problem generator:          %s\n", PROBLEM_GENERATOR)
	fmt.Printf("  Coordinate system:          %s\n", COORDINATE_SYSTEM)
	fmt.Printf("  Equation of state:          %s\n", EQUATION_OF_STATE)
	fmt.Printf("  Riemann solver:             %s\n", RIEMANN_SOLVER)
	if MAGNETIC_FIELDS_ENABLED {
		fmt.Println("  Magnetic fields:            ON")
	} else {
		fmt.Println("  Magnetic fields:            OFF")
	}
	if RELATIVISTIC_DYNAMICS { // configure.py output: "Special relativity"
		fmt.Println("  Relativistic dynamics:      ON ")
	} else {
		fmt.Println("  Relativistic dynamics:      OFF ")
	}
	if GENERAL_RELATIVITY {
		fmt.Println("  General relativity:         ON ")
	} else {
		fmt.Println("  General relativity:         OFF ")
	}
	// configure.py output: "Frame transformations"
	if SELF_GRAVITY_ENABLED == 1 {
		fmt.Println("  Self-Gravity:               FFT")
	} else if SELF_GRAVITY_ENABLED == 2 {
		fmt.Println("  Self-Gravity:               Multigrid")
	} else {
		fmt.Println("  Self-Gravity:               OFF")
	}
	if STS_ENABLED {
		fmt.Println("  Super-Time-Stepping:        ON")
	} else {
		fmt.Println("  Super-Time-Stepping:        OFF")
	}
	// configure.py output: +"Debug flags"
	// configure.py output: +"Code coverage flags"
	// configure.py output: +"Linker flags"
	if SINGLE_PRECISION_ENABLED {
		fmt.Println("  Floating-point precision:   single")
	} else {
		fmt.Println("  Floating-point precision:   double")
	}
	fmt.Printf("  Number of ghost cells:      %d\n", NGHOST)

	if FFT {
		fmt.Println("  FFT:                        ON")
	} else {
		fmt.Println("  FFT:                        OFF")
	}

	if HDF5OUTPUT {
		fmt.Println("  HDF5 output:                ON")
		if H5_DOUBLE_PRECISION_ENABLED {
			fmt.Println("  HDF5 precision:             double")
		} else {
			fmt.Println("  HDF5 precision:             single")
		}
	} else {
		fmt.Println("  HDF5 output:                OFF")
	}

	fmt.Printf("  Compiler:                   %s\n", COMPILED_WITH)
	// configure.py output: Doesnt append "Linker flags" in prev. output (excessive space!)
	return
}
