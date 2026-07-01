// Package ball defines the pool-ball model: position, velocity, spin, color,
// and pocket-drop animation state.
package ball

import (
	"image/color"

	"github.com/user/pooltest/pool-go/internal/vec"
)

// Radius is the visual and collision radius shared by all balls.
const Radius = 11.0

// Ball is a single pool ball. Number 0 is the cue ball; 1-7 are solids, 8 is
// the black ball, and 9-15 are stripes.
type Ball struct {
	Number int
	Pos    vec.Vec2
	Vel    vec.Vec2
	Color  color.RGBA
	Stripe bool
	Active bool // false once pocketed

	// Spin is a reservoir of "stored" linear motion that the cloth grips and
	// bleeds back into Vel over time — this is what makes draw, follow, and
	// swerve work after the cue ball strikes an object ball.
	Spin vec.Vec2
	// SideSpin is english about the vertical axis: it drives the visible roll
	// and tweaks the rebound angle off a cushion.
	SideSpin float64
	// Angle is the accumulated visual roll (radians) used when blitting the
	// rotated ball sprite.
	Angle float64

	// Pocket-drop animation (cosmetic only — physics ignores inactive balls).
	Sinking  bool
	SinkFrom vec.Vec2
	SinkPos  vec.Vec2
	SinkT    float64
}

// Moving reports whether the ball still has motion the simulation must resolve
// (linear velocity or a spin reservoir that has not yet bled out).
func (b *Ball) Moving() bool {
	return b.Active && (b.Vel.LenSq() > 0 || b.Spin.LenSq() > 0)
}
