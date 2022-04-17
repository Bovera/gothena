package utils

/* forward declarations needed for function pointer type aliases
class MeshBlock;
class Coordinates;
class ParameterInput;
class HydroDiffusion;
class FieldDiffusion;
class OrbitalAdvection;*/

//--------------------------------------------------------------------------------------
//! \struct LogicalLocation
//! \brief stores logical location and level of MeshBlock
//!
//! These values can exceed the range of std::int32_t even if the root grid has only a
//! single MeshBlock if >30 levels of AMR are used, since the corresponding max index =
//! 1*2^31 > INT_MAX = 2^31 -1 for most 32-bit signed integer type impelementations.
//! Note that it is default to overload the comparison operator.

type LogicalLocation struct { // aggregate and POD type
	lx1, lx2, lx3 int64
	level         int
}

func (left *LogicalLocation) Lesser(right *LogicalLocation) bool  { return left.level < right.level }
func (left *LogicalLocation) Greater(right *LogicalLocation) bool { return left.level > right.level }

//----------------------------------------------------------------------------------------
//! \struct RegionSize
//! \brief physical size and number of cells in a Mesh or a MeshBlock

type RegionSize struct { // aggregate and POD type; do NOT reorder member declarations:
	x1min, x2min, x3min float64
	x1max, x2max, x3max float64
	x1rat, x2rat, x3rat float64 // ratio of dxf(i)/dxf(i-1)
	// the size of the root grid or a MeshBlock should not exceed std::int32_t limits
	nx1, nx2, nx3 int // number of active cells (not including ghost zones)
}

//---------------------------------------------------------------------------------------
//! \struct FaceField
//! \brief container for face-centered fields

type FaceField struct {
	x1f, x2f, x3f Array[float64]
}

func (this *FaceField) Init(ncells3 int, ncells2 int, ncells1 int) {
	this.x1f = AthenaArray[float64](ncells3, ncells2, ncells1+1)
	this.x2f = AthenaArray[float64](ncells3, ncells2+1, ncells1)
	this.x3f = AthenaArray[float64](ncells3+1, ncells2, ncells1)
}

//----------------------------------------------------------------------------------------
//! \struct EdgeField
//! \brief container for edge-centered fields

type EdgeField struct {
	x1e, x2e, x3e Array[float64]
}

func (this *EdgeField) Init(ncells3 int, ncells2 int, ncells1 int) {
	this.x1e = AthenaArray[float64](ncells3+1, ncells2+1, ncells1)
	this.x2e = AthenaArray[float64](ncells3+1, ncells2, ncells1+1)
	this.x3e = AthenaArray[float64](ncells3, ncells2+1, ncells1+1)
}

//----------------------------------------------------------------------------------------
// enums used everywhere
// (not specifying underlying integral type (C++11) for portability & performance)

//! \todo (felker):
//! - C++ Core Guidelines Enum.5: Don’t use ALL_CAPS for enumerators
//!   (avoid clashes with preprocessor macros).
//! - Enumerated type definitions in this file and:
//!   athena_fft.hpp, io_wrapper.hpp, bvals.hpp, hydro_diffusion.hpp, field_diffusion.hpp,
//!   task_list.hpp, ???

//------------------
// named, weakly typed / unscoped enums:
//------------------

// enumerators only used for indexing AthenaArray and regular arrays; typename and
// explicitly specified enumerator values aare unnecessary, but provided for clarity:

/*-- C++ parts
//! array indices for conserved: density, momemtum, total energy
enum ConsIndex {IDN=0, IM1=1, IM2=2, IM3=3, IEN=4};
//! array indices for face-centered field
enum MagneticIndex {IB1=0, IB2=1, IB3=2};

//! array indices for 1D primitives: velocity, transverse components of field
enum PrimIndex {IVX=1, IVY=2, IVZ=3, IPR=4, IBY=(NHYDRO), IBZ=((NHYDRO)+1)};

//! array indices for face-centered electric fields returned by Riemann solver
enum ElectricIndex {X1E2=0, X1E3=1, X2E3=0, X2E1=1, X3E1=0, X3E2=1};

//! array indices for metric matrices in GR
enum MetricIndex {I00=0, I01=1, I02=2, I03=3, I11=4, I12=5, I13=6, I22=7, I23=8, I33=9,
                  NMETRIC=10};
//! array indices for triangular matrices in GR
enum TriangleIndex {T00=0, T10=1, T11=2, T20=3, T21=4, T22=5, T30=6, T31=7, T32=8, T33=9,
                    NTRIANGULAR=10};

// enumerator types that are used for variables and function parameters:

// needed for arrays dimensioned over grid directions
// enumerator type only used in Mesh::EnrollUserMeshGenerator()

//! array indices for grid directions
enum CoordinateDirection {X1DIR=0, X2DIR=1, X3DIR=2};

//------------------
// strongly typed / scoped enums (C++11):
//------------------
// KGF: Except for the 2x MG* enums, these may be unnessary w/ the new class inheritance
// Now, only passed to BoundaryVariable::InitBoundaryData(); could replace w/ bool switch
// TODO(tomo-ono): consider necessity of orbita_cc and orbital_fc
enum class BoundaryQuantity {cc, fc, cc_flcor, fc_flcor, mggrav,
                             mggrav_f, orbital_cc, orbital_fc};
enum class HydroBoundaryQuantity {cons, prim};
enum class BoundaryCommSubset {mesh_init, gr_amr, all, orbital};
// TODO(felker): consider generalizing/renaming to QuantityFormulation
enum class FluidFormulation {evolve, background, disabled}; // rename background -> fixed?
enum class TaskType {op_split_before, main_int, op_split_after};
enum class UserHistoryOperation {sum, max, min};

//----------------------------------------------------------------------------------------
// function pointer prototypes for user-defined modules set at runtime

using BValFunc = void (*)(
    MeshBlock *pmb, Coordinates *pco, AthenaArray<Real> &prim, FaceField &b,
    Real time, Real dt,
    int is, int ie, int js, int je, int ks, int ke, int ngh);
using AMRFlagFunc = int (*)(MeshBlock *pmb);
using MeshGenFunc = Real (*)(Real x, RegionSize rs);
using SrcTermFunc = void (*)(
    MeshBlock *pmb, const Real time, const Real dt, const AthenaArray<Real> &prim,
    const AthenaArray<Real> &prim_scalar, const AthenaArray<Real> &bcc,
    AthenaArray<Real> &cons, AthenaArray<Real> &cons_scalar);
using TimeStepFunc = Real (*)(MeshBlock *pmb);
using HistoryOutputFunc = Real (*)(MeshBlock *pmb, int iout);
using MetricFunc = void (*)(
    Real x1, Real x2, Real x3, ParameterInput *pin,
    AthenaArray<Real> &g, AthenaArray<Real> &g_inv,
    AthenaArray<Real> &dg_dx1, AthenaArray<Real> &dg_dx2, AthenaArray<Real> &dg_dx3);
using MGBoundaryFunc = void (*)(
    AthenaArray<Real> &dst,Real time, int nvar,
    int is, int ie, int js, int je, int ks, int ke, int ngh,
    Real x0, Real y0, Real z0, Real dx, Real dy, Real dz);
using ViscosityCoeffFunc = void (*)(
    HydroDiffusion *phdif, MeshBlock *pmb,
    const  AthenaArray<Real> &w, const AthenaArray<Real> &bc,
    int is, int ie, int js, int je, int ks, int ke);
using ConductionCoeffFunc = void (*)(
    HydroDiffusion *phdif, MeshBlock *pmb,
    const AthenaArray<Real> &w, const AthenaArray<Real> &bc,
    int is, int ie, int js, int je, int ks, int ke);
using FieldDiffusionCoeffFunc = void (*)(
    FieldDiffusion *pfdif, MeshBlock *pmb,
    const AthenaArray<Real> &w,
    const AthenaArray<Real> &bmag,
    int is, int ie, int js, int je, int ks, int ke);
using OrbitalVelocityFunc = Real (*)(
    OrbitalAdvection *porb, Real x1, Real x2, Real x3);
--*/
