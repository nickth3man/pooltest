package ball

import (
	"github.com/user/pooltest/pool-go/internal/vec"
)

// MotionState labels ball kinematic phase (pooltool constants).
type MotionState int

const (
	MotionStationary MotionState = iota
	MotionSpinning
	MotionSliding
	MotionRolling
	MotionPocketed
)

// RVW is position, velocity, and angular velocity (SI: m, m/s, rad/s).
type RVW struct {
	R vec.Vec3
	V vec.Vec3
	W vec.Vec3
}

// Copy returns a deep copy.
func (r RVW) Copy() RVW { return r }

// SyncXY copies RVW into pixel Pos/Vel and advances visual roll angle.
func (b *Ball) SyncXY(dtSec float64) {
	b.frameDt = dtSec
	b.Pos = vec.Vec2{
		X: MetersToPx(b.RVW.R.X),
		Y: MetersToPx(b.RVW.R.Y),
	}
	tps := 1.0 / dtSec
	b.Vel = vec.Vec2{
		X: MetersToPx(b.RVW.V.X) / tps,
		Y: MetersToPx(b.RVW.V.Y) / tps,
	}
	b.Angle += b.RVW.W.Z * dtSec
}

// InitFromXY sets RVW from pixel position with zero motion at table height.
func (b *Ball) InitFromXY(pos vec.Vec2, params Params) {
	b.Params = params
	b.SetPosPx(pos)
	b.RVW.V = vec.Vec3{}
	b.RVW.W = vec.Vec3{}
	b.Motion = MotionStationary
	b.Vel = vec.Vec2{}
}

const (
	epsV = 1e-6
	epsW = 2e-2
)
