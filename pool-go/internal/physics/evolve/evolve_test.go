package evolve

import (
	"math"
	"testing"

	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/physics/ptmath"
	"github.com/user/pooltest/pool-go/internal/vec"
)

func TestSlideToStop(t *testing.T) {
	p := ball.PoolGeneric()
	rvw := ball.RVW{
		R: vec.Vec3{Z: p.R},
		V: vec.Vec3{X: 2},
	}
	state := ball.MotionSliding
	total := ptmath.GetSlideTime(ptmath.RVWFrom(rvw.R, rvw.V, rvw.W), p.R, p.Us, p.G)
	total += ptmath.GetRollTime(ptmath.RVWFrom(rvw.R, rvw.V, rvw.W), p.Ur, p.G)
	rvw, state = EvolveBallMotion(state, rvw, p, total*1.5)
	if state != ball.MotionStationary && ptmath.Norm3(rvw.V) > 1e-4 {
		t.Fatalf("ball still moving: state=%v v=%v", state, rvw.V)
	}
}

func TestDrawCurvesBack(t *testing.T) {
	p := ball.PoolGeneric()
	rvw := ball.RVW{
		R: vec.Vec3{Z: p.R},
		V: vec.Vec3{X: 0.5},
		W: vec.Vec3{Y: -80},
	}
	startX := rvw.R.X
	rvw, state := EvolveBallMotion(ball.MotionSliding, rvw, p, 8.0)
	if state == ball.MotionSliding && rvw.R.X > startX+0.05 {
		t.Fatalf("draw should not advance far forward during slide: x=%v", rvw.R.X)
	}
}

func TestPerpendicularSpinDecays(t *testing.T) {
	p := ball.PoolGeneric()
	rvw := ball.RVW{
		R: vec.Vec3{Z: p.R},
		W: vec.Vec3{Z: 5},
	}
	tSpin := ptmath.GetSpinTime(ptmath.RVWFrom(rvw.R, rvw.V, rvw.W), p.R, p.USp(), p.G)
	rvw, state := EvolveBallMotion(ball.MotionSpinning, rvw, p, tSpin*1.2)
	if math.Abs(rvw.W.Z) > 1e-3 {
		t.Fatalf("spin not decayed: wz=%v", rvw.W.Z)
	}
	if state != ball.MotionStationary {
		t.Fatalf("expected stationary, got %v", state)
	}
}
