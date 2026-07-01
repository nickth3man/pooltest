package ball

import "math"

// Params holds physical constants for a pool ball (SI internally, scaled via units).
type Params struct {
	M  float64 // mass (kg)
	R  float64 // radius (m)
	Us float64 // sliding friction (ball-cloth)
	Ur float64 // rolling friction
	Ub float64 // ball-ball friction (average model fallback)
	Eb float64 // ball-ball restitution
	Ec float64 // cushion restitution
	Fc float64 // cushion friction
	G  float64 // gravity (m/s²)

	// USpProp is spinning friction proportionality; USp = USpProp * R.
	USpProp float64

	// CueEndMass is effective cue endmass for squirt (kg).
	CueEndMass float64
	// CueMass is cue stick mass (kg).
	CueMass float64
	// CueElevation is implicit cue elevation (radians) for swerve in top-down view.
	CueElevation float64
}

// USp returns the spinning friction coefficient.
func (p Params) USp() float64 { return p.USpProp * p.R }

// IOverM returns moment of inertia divided by mass (2/5 R²).
func (p Params) IOverM() float64 { return 2.0 / 5.0 * p.R * p.R }

// PoolGeneric returns pooltool POOL_GENERIC defaults.
func PoolGeneric() Params {
	return Params{
		M:            0.170097,
		R:            0.028575,
		Us:           0.2,
		Ur:           0.01,
		Ub:           0.05,
		Eb:           0.95,
		Ec:           0.85,
		Fc:           0.2,
		G:            9.81,
		USpProp:      10 * 2 / 5 / 9,
		CueEndMass:   0.05,
		CueMass:      0.567,
		CueElevation: 6 * math.Pi / 180,
	}
}

// Arcade scales friction for snappier play while keeping phase structure.
func Arcade() Params {
	p := PoolGeneric()
	p.Ur *= 1.4
	p.Us *= 0.9
	return p
}

// Realistic returns pooltool POOL_GENERIC (alias for clarity).
func Realistic() Params { return PoolGeneric() }

// DefaultParams is the active simulation preset.
var DefaultParams = Realistic()
