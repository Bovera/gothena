// configure.py dict(definitions) string values:
package utils

const (
	// problem generator
	PROBLEM_GENERATOR = "shock_tube"

	// coordinate system
	COORDINATE_SYSTEM = "cartesian"

	// Riemann solver
	RIEMANN_SOLVER = "hllc"

	// configure.py dict(definitions) Boolean values:
	// Equation of state
	EQUATION_OF_STATE = "adiabatic"

	// use general EOS framework default=0 (false).
	GENERAL_EOS = false

	// use EOS table default=0 (false).
	EOS_TABLE_ENABLED = false

	// non-barotropic equation of state (i.e. P not simply a func of rho)? default=1 (true)
	NON_BAROTROPIC_EOS = true

	// include magnetic fields? default=0 (false)
	MAGNETIC_FIELDS_ENABLED = false

	// include super-time-stepping? default=0 (false)
	STS_ENABLED = false

	// include self gravity? default=0 (false)
	SELF_GRAVITY_ENABLED = 0

	// include radiative transfer? default=0 (false)
	RADIATION_ENABLED = false

	// enable special or general relativity? default=0 (false)
	RELATIVISTIC_DYNAMICS = false

	// enable general relativity? default=0 (false)
	GENERAL_RELATIVITY = false

	// enable GR frame transformations? default=0 (false)
	FRAME_TRANSFORMATIONS = false

	// use single precision floating-point values (binary32)? default=0 (false; use binary64)
	SINGLE_PRECISION_ENABLED = false

	// use double precision for HDF5 output? default=0 (false; write out binary32)
	H5_DOUBLE_PRECISION_ENABLED = false

	// configure.py dict(definitions) Boolean string macros:
	// (these options have the latter (false) option as defaults, unless noted otherwise)
	// make use of FFT? (FFT or NO_FFT)
	FFT = false

	// HDF5 output (HDF5OUTPUT or NO_HDF5OUTPUT)
	HDF5OUTPUT = false

	// debug build macros (DEBUG or NOT_DEBUG)
	DEBUG = false

	// compiler options
	COMPILED_WITH = "Go1.18"

	//----------------------------------------------------------------------------------------
	// macros associated with numerical algorithm (rarely modified)

	NHYDRO        = 5
	NFIELD        = 0
	NWAVE         = 5
	NSCALARS      = 0
	NGHOST        = 2
	MAX_NSTAGE    = 6 // maximum number of stages per cycle for time-integrator
	MAX_NREGISTER = 3 // maximum number of (u, b) register pairs for time-integrator

	//----------------------------------------------------------------------------------------
	// general purpose macros (never modified)

	// all constants specified to 17 total digits of precision = max_digits10 for "double"
	PI             = 3.1415926535897932
	TWO_PI         = 6.2831853071795862
	SQRT2          = 1.4142135623730951
	ONE_OVER_SQRT2 = 0.70710678118654752
	ONE_3RD        = 0.33333333333333333
	TWO_3RD        = 0.66666666666666667
	TINY_NUMBER    = 1.0e-20
	HUGE_NUMBER    = 1.0e+36
)
