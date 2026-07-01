package ballcushion

import (
	"math"
	"testing"

	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/vec"
)

func cushionTestBall(v vec.Vec3, w vec.Vec3) *ball.Ball {
	b := &ball.Ball{Number: 0, Active: true, Params: ball.DefaultParams}
	b.InitFromXY(vec.Vec2{X: 700, Y: 300}, b.Params)
	b.RVW.V = v
	b.RVW.W = w
	b.Motion = ball.MotionSliding
	return b
}

func TestMathavanNoSpinReboundsAndLosesSpeed(t *testing.T) {
	b := cushionTestBall(vec.Vec3{X: 1.2}, vec.Vec3{})
	before := b.RVW.V.Len()
	ResolveLinear(b, vec.Vec2{X: -1})
	after := b.RVW.V.Len()
	if b.RVW.V.X >= 0 {
		t.Fatalf("vx = %.6f, want rebound away from cushion", b.RVW.V.X)
	}
	if after <= 0 || after >= before {
		t.Fatalf("speed after %.6f, want positive and below %.6f", after, before)
	}
}

func TestMathavanSpinChangesReboundTangent(t *testing.T) {
	noSpin := cushionTestBall(vec.Vec3{X: 1.2}, vec.Vec3{})
	sideSpin := cushionTestBall(vec.Vec3{X: 1.2}, vec.Vec3{Z: 25})
	ResolveLinear(noSpin, vec.Vec2{X: -1})
	ResolveLinear(sideSpin, vec.Vec2{X: -1})
	if math.Abs(sideSpin.RVW.V.Y-noSpin.RVW.V.Y) < 0.01 {
		t.Fatalf("sidespin did not change tangent rebound: noSpin=%.6f sideSpin=%.6f", noSpin.RVW.V.Y, sideSpin.RVW.V.Y)
	}
}

func TestMathavanDrawAndFollowRemainFinite(t *testing.T) {
	for name, wy := range map[string]float64{"draw": -30, "follow": 30} {
		b := cushionTestBall(vec.Vec3{X: 1.2}, vec.Vec3{Y: wy})
		ResolveLinear(b, vec.Vec2{X: -1})
		if b.RVW.V.X >= 0 {
			t.Fatalf("%s vx = %.6f, want rebound", name, b.RVW.V.X)
		}
		if math.IsNaN(b.RVW.V.Len()) || math.IsInf(b.RVW.V.Len(), 0) || b.RVW.V.Len() == 0 {
			t.Fatalf("%s produced invalid velocity %+v", name, b.RVW.V)
		}
	}
}
