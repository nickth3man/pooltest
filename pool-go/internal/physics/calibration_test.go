package physics

import (
	"math"
	"testing"

	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/physics/resolve/ballball"
	"github.com/user/pooltest/pool-go/internal/physics/resolve/stickball"
	"github.com/user/pooltest/pool-go/internal/vec"
)

func TestAlciatoreMuCurve(t *testing.T) {
	slow := ballball.AlciatoreMu(0.05)
	fast := ballball.AlciatoreMu(2.0)
	if slow <= fast {
		t.Fatalf("slow mu %.4f should exceed fast mu %.4f", slow, fast)
	}
	if slow < 0.03 || slow > 0.12 {
		t.Fatalf("slow mu %.4f outside Dr. Dave 0.03-0.08 band (extended at low v)", slow)
	}
}

func TestSlowHalfBallThrowInRange(t *testing.T) {
	const cut = math.Pi / 6
	striker := vec.Vec2{X: 200, Y: 200}
	target := vec.Vec2{
		X: 200 + (2*ball.CollisionRadius()-0.5)*math.Cos(cut),
		Y: 200 + (2*ball.CollisionRadius()-0.5)*math.Sin(cut),
	}
	throw := PredictThrowDeg(striker, target, vec.Vec2{X: 1, Y: 0}, 3, 0, 0)
	if throw < 2 || throw > 8 {
		t.Fatalf("slow half-ball throw %.2f° outside 2-8° band", throw)
	}
}

func TestSquirtAngleIncreasesWithSideSpin(t *testing.T) {
	b := &ball.Ball{Number: 0, Active: true, Params: ball.DefaultParams}
	b.InitFromXY(vec.Vec2{X: 300, Y: 300}, ball.DefaultParams)
	stickball.StrikeFromGame(b, vec.Vec2{X: 1, Y: 0}, 8, 0, 0)
	v0 := b.RVW.V
	b.InitFromXY(vec.Vec2{X: 300, Y: 300}, ball.DefaultParams)
	stickball.StrikeFromGame(b, vec.Vec2{X: 1, Y: 0}, 8, 0.8, 0)
	v1 := b.RVW.V
	angle0 := math.Atan2(v0.Y, v0.X)
	angle1 := math.Atan2(v1.Y, v1.X)
	if math.Abs(angle1-angle0) < 0.005 {
		t.Fatalf("squirt deflection too small: d=%v", angle1-angle0)
	}
}

func TestDrawFollowThrowReduction(t *testing.T) {
	const cut = math.Pi / 6
	striker := vec.Vec2{X: 200, Y: 200}
	target := vec.Vec2{
		X: 200 + (2*ball.CollisionRadius()-0.5)*math.Cos(cut),
		Y: 200 + (2*ball.CollisionRadius()-0.5)*math.Sin(cut),
	}
	stun := PredictThrowDeg(striker, target, vec.Vec2{X: 1, Y: 0}, 5, 0, 0)
	draw := PredictThrowDeg(striker, target, vec.Vec2{X: 1, Y: 0}, 5, 0, -0.6)
	if draw >= stun*0.75 {
		t.Fatalf("draw throw %.2f should be <75%% of stun %.2f", draw, stun)
	}
}

func TestSlideRollTransitionDistance(t *testing.T) {
	b := &ball.Ball{Number: 0, Active: true, Params: ball.PoolGeneric()}
	b.InitFromXY(vec.Vec2{X: 300, Y: 300}, b.Params)
	stickball.StrikeFromGame(b, vec.Vec2{X: 1, Y: 0}, 10, 0, 0)
	start := b.Pos.X
	balls := []*ball.Ball{b}
	rolled := false
	for range 500 {
		Step(balls)
		if b.Motion == ball.MotionRolling || b.Motion == ball.MotionSpinning || b.Motion == ball.MotionStationary {
			rolled = true
			break
		}
	}
	if !rolled {
		t.Fatal("ball never left sliding phase")
	}
	dist := math.Abs(b.Pos.X - start)
	if dist < 150 || dist > 400 {
		t.Fatalf("slide distance %.1f px outside expected band", dist)
	}
}

func TestArcadePresetSnappierDecay(t *testing.T) {
	g := ball.PoolGeneric()
	a := ball.Arcade()
	if a.Ur <= g.Ur {
		t.Fatal("arcade rolling friction should be higher than realistic")
	}
}

func TestRealisticPresetMatchesGeneric(t *testing.T) {
	if ball.Realistic().Eb != ball.PoolGeneric().Eb {
		t.Fatal("realistic preset should match POOL_GENERIC")
	}
}
