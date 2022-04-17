package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"runtime"
	"time"
)

import (
	"gothena/inputs"
	"gothena/utils"
)

var mbcnt uint64

func main() {

	//--- Step 1. --------------------------------------------------------------------
	// Set CPU number for parallel and check for command line options and respond.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Input -h will automatically get the help menu
	input_filename := flag.String("i", "", "specify input file [athinput]")
	restart_filename := flag.String("r", "", "restart with this file")
	// prundir := flag.String("d", "", "specify run dir [current dir]") TODO
	// narg_flag := flag.Bool("n", false, "parse input file and quit") TODO
	config_flag := flag.Bool("c", false, "show configuration and quit")
	// mesh_flag := flag.Int("m", 0, "output mesh structure and quit") TODO
	// set to <nproc> if -m <nproc> argument is on cmdline
	wtlim := flag.Duration("t", 0, "wall time limit for final output")

	flag.Parse()

	if *config_flag { utils.ShowConfig(); return }

	if *input_filename == "" && *restart_filename == "" {
		panic("No input file or restart file is specified.")
	}

	// Set up the timer
	if *wtlim != 0 {
		timeout := make(chan bool)
		go func() { time.Sleep(*wtlim); timeout <- true }()
	}

	// my_rank? TODO

	//--- Step 2. --------------------------------------------------------------------------
	// Construct object to store input parameters, then parse input file and command line.
	// With MPI, the input is read by every process in parallel using MPI-IO.

	var pinput inputs.ParameterInput

	if *restart_filename != "" {
		str, err := ioutil.ReadFile(*restart_filename)
		if err != nil {
			panic(err)
		}
		err = pinput.LoadFromByte(str)
		if err != nil {
			panic(err)
		}
		// If both -r and -i are specified, make sure next_time gets corrected.
		// This needs to be corrected on the restart file because we need the old dt.
		if *input_filename != "" {
			pinput.ChangeNextTime(-1)
		}
	}
	if *input_filename != "" {
		// if both -r and -i are specified, override the parameters using the input file
		str, err := ioutil.ReadFile(*input_filename)
		if err != nil {
			panic(err)
		}
		err = pinput.LoadFromByte(str)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println(pinput.ParameterDump())
}

// Modify from cmd line isn't surported.
/*
	//--- Step 3. --------------------------------------------------------------------------
	// Construct and initialize Mesh

  Mesh *pmesh;
#ifdef ENABLE_EXCEPTIONS
  try {
#endif
    if (res_flag == 0) {
      pmesh = new Mesh(pinput, mesh_flag);
    } else {
      pmesh = new Mesh(pinput, restartfile, mesh_flag);
    }
#ifdef ENABLE_EXCEPTIONS
  }
  catch(std::bad_alloc& ba) {
    std::cout << "### FATAL ERROR in main" << std::endl
              << "memory allocation failed initializing class Mesh: "
              << ba.what() << std::endl;
    if (res_flag == 1) restartfile.Close();
#ifdef MPI_PARALLEL
    MPI_Finalize();
#endif
    return(0);
  }
  catch(std::exception const& ex) {
    std::cout << ex.what() << std::endl;  // prints diagnostic message
    if (res_flag == 1) restartfile.Close();
#ifdef MPI_PARALLEL
    MPI_Finalize();
#endif
    return(0);
  }
#endif // ENABLE_EXCEPTIONS

  // With current mesh time possibly read from restart file, correct next_time for outputs
  if (iarg_flag == 1 && res_flag == 1) {
    // if both -r and -i are specified, ensure that next_time  >= mesh_time - dt
    pinput->ForwardNextTime(pmesh->time);
  }

  // Dump input parameters and quit if code was run with -n option.
  if (narg_flag) {
    if (Globals::my_rank == 0) pinput->ParameterDump(std::cout);
    if (res_flag == 1) restartfile.Close();
#ifdef MPI_PARALLEL
    MPI_Finalize();
#endif
    return(0);
  }

  if (res_flag == 1) restartfile.Close(); // close the restart file here

  // Quit if -m was on cmdline.  This option builds and outputs mesh structure.
  if (mesh_flag > 0) {
#ifdef MPI_PARALLEL
    MPI_Finalize();
#endif
    return(0);
  }

  //--- Step 4. --------------------------------------------------------------------------
  // Construct and initialize TaskList

  TimeIntegratorTaskList *ptlist;
#ifdef ENABLE_EXCEPTIONS
  try {
#endif
    ptlist = new TimeIntegratorTaskList(pinput, pmesh);
#ifdef ENABLE_EXCEPTIONS
  }
  catch(std::bad_alloc& ba) {
    std::cout << "### FATAL ERROR in main" << std::endl << "memory allocation failed "
              << "in creating task list " << ba.what() << std::endl;
#ifdef MPI_PARALLEL
    MPI_Finalize();
#endif
    return(0);
  }
#endif // ENABLE_EXCEPTIONS

  SuperTimeStepTaskList *pststlist = nullptr;
  if (STS_ENABLED) {
#ifdef ENABLE_EXCEPTIONS
    try {
#endif
      pststlist = new SuperTimeStepTaskList(pinput, pmesh, ptlist);
#ifdef ENABLE_EXCEPTIONS
    }
    catch(std::bad_alloc& ba) {
      std::cout << "### FATAL ERROR in main" << std::endl << "memory allocation failed "
                << "in creating task list " << ba.what() << std::endl;
#ifdef MPI_PARALLEL
      MPI_Finalize();
#endif
      return(0);
    }
#endif // ENABLE_EXCEPTIONS
  }

  //--- Step 5. --------------------------------------------------------------------------
  // Set initial conditions by calling problem generator, or reading restart file

#ifdef ENABLE_EXCEPTIONS
  try {
#endif
    pmesh->Initialize(res_flag, pinput);
#ifdef ENABLE_EXCEPTIONS
  }
  catch(std::bad_alloc& ba) {
    std::cout << "### FATAL ERROR in main" << std::endl << "memory allocation failed "
              << "in problem generator " << ba.what() << std::endl;
#ifdef MPI_PARALLEL
    MPI_Finalize();
#endif
    return(0);
  }
  catch(std::exception const& ex) {
    std::cout << ex.what() << std::endl;  // prints diagnostic message
#ifdef MPI_PARALLEL
    MPI_Finalize();
#endif
    return(0);
  }
#endif // ENABLE_EXCEPTIONS

  //--- Step 6. --------------------------------------------------------------------------
  // Change to run directory, initialize outputs object, and make output of ICs

  Outputs *pouts;
#ifdef ENABLE_EXCEPTIONS
  try {
#endif
    ChangeRunDir(prundir);
    pouts = new Outputs(pmesh, pinput);
    if (res_flag == 0) pouts->MakeOutputs(pmesh, pinput);
#ifdef ENABLE_EXCEPTIONS
  }
  catch(std::bad_alloc& ba) {
    std::cout << "### FATAL ERROR in main" << std::endl
              << "memory allocation failed setting initial conditions: "
              << ba.what() << std::endl;
#ifdef MPI_PARALLEL
    MPI_Finalize();
#endif
    return(0);
  }
  catch(std::exception const& ex) {
    std::cout << ex.what() << std::endl;  // prints diagnostic message
#ifdef MPI_PARALLEL
    MPI_Finalize();
#endif
    return(0);
  }
#endif // ENABLE_EXCEPTIONS

  //=== Step 7. === START OF MAIN INTEGRATION LOOP =======================================
  // For performance, there is no error handler protecting this step (except outputs)

  if (Globals::my_rank == 0) {
    std::cout << "\nSetup complete, entering main loop...\n" << std::endl;
  }

  clock_t tstart = clock();
#ifdef OPENMP_PARALLEL
  double omp_start_time = omp_get_wtime();
#endif

  while ((pmesh->time < pmesh->tlim) &&
         (pmesh->nlim < 0 || pmesh->ncycle < pmesh->nlim)) {
    if (Globals::my_rank == 0)
      pmesh->OutputCycleDiagnostics();

    if (STS_ENABLED) {
      pmesh->sts_loc = TaskType::op_split_before;
      // compute nstages for this STS
      if (pmesh->sts_integrator == "rkl2") { // default
        pststlist->nstages =
            static_cast<int>
              (0.5*(-1. + std::sqrt(9. + 16.*(0.5*pmesh->dt)/pmesh->dt_parabolic))) + 1;
      } else { // rkl1
        pststlist->nstages =
            static_cast<int>
              (0.5*(-1. + std::sqrt(1. + 8.*pmesh->dt/pmesh->dt_parabolic))) + 1;
      }
      if (pststlist->nstages % 2 == 0) { // guarantee odd nstages for STS
        pststlist->nstages += 1;
      }
      // take super-timestep
      for (int stage=1; stage<=pststlist->nstages; ++stage)
        pststlist->DoTaskListOneStage(pmesh, stage);

      pmesh->sts_loc = TaskType::main_int;
    }

    if (pmesh->turb_flag > 1) pmesh->ptrbd->Driving(); // driven turbulence

    for (int stage=1; stage<=ptlist->nstages; ++stage) {
      ptlist->DoTaskListOneStage(pmesh, stage);
      if (ptlist->CheckNextMainStage(stage)) {
        if (SELF_GRAVITY_ENABLED == 1) // fft (0: discrete kernel, 1: continuous kernel)
          pmesh->pfgrd->Solve(stage, 0);
        else if (SELF_GRAVITY_ENABLED == 2) // multigrid
          pmesh->pmgrd->Solve(stage);
      }
    }

    if (STS_ENABLED && pmesh->sts_integrator == "rkl2") {
      pmesh->sts_loc = TaskType::op_split_after;
      // take super-timestep
      for (int stage=1; stage<=pststlist->nstages; ++stage)
        pststlist->DoTaskListOneStage(pmesh, stage);
    }

    pmesh->UserWorkInLoop();

    pmesh->ncycle++;
    pmesh->time += pmesh->dt;
    mbcnt += pmesh->nbtotal;
    pmesh->step_since_lb++;

    pmesh->LoadBalancingAndAdaptiveMeshRefinement(pinput);

    pmesh->NewTimeStep();

#ifdef ENABLE_EXCEPTIONS
    try {
#endif
      if (pmesh->time < pmesh->tlim) // skip the final output as it happens later
        pouts->MakeOutputs(pmesh,pinput);
#ifdef ENABLE_EXCEPTIONS
    }
    catch(std::bad_alloc& ba) {
      std::cout << "### FATAL ERROR in main" << std::endl
                << "memory allocation failed during output: " << ba.what() <<std::endl;
#ifdef MPI_PARALLEL
      MPI_Finalize();
#endif
      return(0);
    }
    catch(std::exception const& ex) {
      std::cout << ex.what() << std::endl;  // prints diagnostic message
#ifdef MPI_PARALLEL
      MPI_Finalize();
#endif
      return(0);
    }
#endif // ENABLE_EXCEPTIONS

    // check for signals
    if (SignalHandler::CheckSignalFlags() != 0) {
      break;
    }
  } // END OF MAIN INTEGRATION LOOP ======================================================
  // Make final outputs, print diagnostics, clean up and terminate

  if (Globals::my_rank == 0 && wtlim > 0)
    SignalHandler::CancelWallTimeAlarm();


  //--- Step 8. --------------------------------------------------------------------------
  // Output the final cycle diagnostics and make the final outputs and print diagnostic
  // messages related to the end of the simulation

  if (Globals::my_rank == 0)
    pmesh->OutputCycleDiagnostics();

  pmesh->UserWorkAfterLoop(pinput);

  pouts->MakeOutputs(pmesh,pinput,true);

  if (Globals::my_rank == 0) {
    if (SignalHandler::GetSignalFlag(SIGTERM) != 0) {
      std::cout << std::endl << "Terminating on Terminate signal" << std::endl;
    } else if (SignalHandler::GetSignalFlag(SIGINT) != 0) {
      std::cout << std::endl << "Terminating on Interrupt signal" << std::endl;
    } else if (SignalHandler::GetSignalFlag(SIGALRM) != 0) {
      std::cout << std::endl << "Terminating on wall-time limit" << std::endl;
    } else if (pmesh->ncycle == pmesh->nlim) {
      std::cout << std::endl << "Terminating on cycle limit" << std::endl;
    } else {
      std::cout << std::endl << "Terminating on time limit" << std::endl;
    }

    std::cout << "time=" << pmesh->time << " cycle=" << pmesh->ncycle << std::endl;
    std::cout << "tlim=" << pmesh->tlim << " nlim=" << pmesh->nlim << std::endl;

    if (pmesh->adaptive) {
      std::cout << std::endl << "Number of MeshBlocks = " << pmesh->nbtotal
                << "; " << pmesh->nbnew << "  created, " << pmesh->nbdel
                << " destroyed during this simulation." << std::endl;
    }

    // Calculate and print the zone-cycles/cpu-second and wall-second
#ifdef OPENMP_PARALLEL
    double omp_time = omp_get_wtime() - omp_start_time;
#endif
    clock_t tstop = clock();
    double cpu_time = (tstop>tstart ? static_cast<double> (tstop-tstart) :
                       1.0)/static_cast<double> (CLOCKS_PER_SEC);
    std::uint64_t zonecycles = mbcnt
      *static_cast<std::uint64_t> (pmesh->my_blocks(0)->GetNumberOfMeshBlockCells());
    double zc_cpus = static_cast<double> (zonecycles) / cpu_time;

    std::cout << std::endl << "zone-cycles = " << zonecycles << std::endl;
    std::cout << "cpu time used  = " << cpu_time << std::endl;
    std::cout << "zone-cycles/cpu_second = " << zc_cpus << std::endl;
#ifdef OPENMP_PARALLEL
    double zc_omps = static_cast<double> (zonecycles) / omp_time;
    std::cout << std::endl << "omp wtime used = " << omp_time << std::endl;
    std::cout << "zone-cycles/omp_wsecond = " << zc_omps << std::endl;
#endif
  }
}*/
