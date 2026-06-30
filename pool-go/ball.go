package main

import "image/color"

const ballRadius = 11.0

// Ball is a single pool ball. Number 0 is the cue ball; 1-7 are solids, 8 is
// the black ball, and 9-15 are stripes.
type Ball struct {
	Number int
	Pos    Vec2
	Vel    Vec2
	Color  color.RGBA
	Stripe bool
	Active bool // false once pocketed

	// Spin is a reservoir of "stored" linear motion that the cloth grips and
	// bleeds back into Vel over time — this is what makes draw, follow, and
	// swerve work after the cue ball strikes an object ball.
	Spin Vec2
	// SideSpin is english about the vertical axis: it drives the visible roll
	// and tweaks the rebound angle off a cushion.
	SideSpin float64
	// Angle is the accumulated visual roll (radians) used when blitting the
	// rotated ball sprite.
	Angle float64

	// Pocket-drop animation (cosmetic only — physics ignores inactive balls).
	Sinking  bool
	SinkFrom Vec2
	SinkPos  Vec2
	SinkT    float64
}

// moving reports whether the ball still has motion the simulation must resolve
// (linear velocity or a spin reservoir that has not yet bled out).
func (b *Ball) moving() bool {
	return b.Active && (b.Vel.LenSq() > 0 || b.Spin.LenSq() > 0)
}

// ballColors maps the seven base colors to ball numbers 1-7 (and, reused, 9-15
// for the matching stripes). The 8-ball is black.
var ballColors = map[int]color.RGBA{
	1: {0xE6, 0xC2, 0x29, 0xFF}, // yellow
	2: {0x1E, 0x4F, 0xD8, 0xFF}, // blue
	3: {0xD1, 0x2B, 0x1E, 0xFF}, // red
	4: {0x6A, 0x2C, 0x91, 0xFF}, // purple
	5: {0xE0, 0x6A, 0x1B, 0xFF}, // orange
	6: {0x1F, 0x8A, 0x3E, 0xFF}, // green
	7: {0x8A, 0x2B, 0x2B, 0xFF}, // maroon
	8: {0x1A, 0x1A, 0x1A, 0xFF}, // black
}
