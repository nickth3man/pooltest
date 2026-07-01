// Package ball defines the pool-ball model: kinematic state, color, and
// pocket-drop animation.
package ball

import (
	"image/color"
	"math"

	"github.com/user/pooltest/pool-go/internal/vec"
)

// Ball is a single pool ball. Number 0 is the cue ball; 1-7 are solids, 8 is
// the black ball, and 9-15 are stripes.
type Ball struct {
	Number int
	Pos    vec.Vec2 // pixels (synced from RVW after each physics step)
	Vel    vec.Vec2 // px/frame display velocity (synced from RVW)
	Color  color.RGBA
	Stripe bool
	Active bool

	Motion MotionState
	RVW    RVW
	Params Params

	// Angle is visual roll (radians) from integrated ω_z.
	Angle float64

	frameDt float64 // last step dt in seconds (for angle integration)

	// Pocket-drop animation (cosmetic).
	Sinking  bool
	SinkFrom vec.Vec2
	SinkPos  vec.Vec2
	SinkT    float64
}

// Moving reports whether the ball still has motion the simulation must resolve.
func (b *Ball) Moving() bool {
	if !b.Active || b.Motion == MotionPocketed {
		return false
	}
	if b.Motion != MotionStationary {
		return true
	}
	return b.RVW.V.LenSq() > epsV*epsV || math.Abs(b.RVW.W.Z) > epsW
}

// CollisionR returns physics radius in pixels.
func (b *Ball) CollisionR() float64 {
	return CollisionRadiusPx(b.Params)
}

// MetersToPx converts a scalar distance from meters to pixels.
func MetersToPx(m float64) float64 { return m * PxPerM() }

// PxToMeters converts pixels to meters.
func PxToMeters(px float64) float64 { return px / PxPerM() }

// PosM returns position in meters.
func (b *Ball) PosM() vec.Vec2 {
	return vec.Vec2{X: PxToMeters(b.Pos.X), Y: PxToMeters(b.Pos.Y)}
}

// SetPosPx sets pixel position and RVW.R xy; z = ball radius in meters.
func (b *Ball) SetPosPx(p vec.Vec2) {
	b.Pos = p
	b.RVW.R.X = PxToMeters(p.X)
	b.RVW.R.Y = PxToMeters(p.Y)
	b.RVW.R.Z = b.Params.R
}
