package sim

import (
	"math"
	"testing"

	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/physics/ptmath"
	"github.com/user/pooltest/pool-go/internal/table"
	"github.com/user/pooltest/pool-go/internal/vec"
)

func testBall(num int, pos vec.Vec2, vel vec.Vec2, state ball.MotionState) *ball.Ball {
	b := &ball.Ball{Number: num, Active: true, Params: ball.DefaultParams}
	b.InitFromXY(pos, b.Params)
	b.RVW.V = vec.Vec3{
		X: ball.PxToMeters(vel.X * 60),
		Y: ball.PxToMeters(vel.Y * 60),
	}
	b.Motion = state
	return b
}

func TestSlidingRailTimeUsesAcceleration(t *testing.T) {
	cr := ball.CollisionRadius()
	gapPx := 20.0
	b := testBall(0, vec.Vec2{X: float64(table.PlayRight) - cr - gapPx, Y: 300}, vec.Vec2{X: 5}, ball.MotionSliding)

	got, normal := nextRail(b)
	if normal.X != -1 || normal.Y != 0 {
		t.Fatalf("normal = %+v, want right cushion", normal)
	}

	dist := ball.PxToMeters(gapPx)
	speed := ball.PxToMeters(5 * 60)
	t1, t2, ok := ptmath.SolveQuadratic(-0.5*b.Params.Us*b.Params.G, speed, -dist)
	if !ok {
		t.Fatal("expected quadratic rail root")
	}
	want := math.Min(t1, t2)
	if want <= 0 {
		want = math.Max(t1, t2)
	}
	if math.Abs(got-want) > 1e-6 {
		t.Fatalf("rail time %.9f, want %.9f", got, want)
	}
	if constant := dist / speed; math.Abs(got-constant) < 1e-4 {
		t.Fatalf("rail time %.9f looks constant-velocity, constant %.9f", got, constant)
	}
}

func TestBallBallTimeUsesAcceleration(t *testing.T) {
	cr := ball.CollisionRadius()
	gapPx := 15.0
	a := testBall(0, vec.Vec2{X: 200, Y: 300}, vec.Vec2{X: 5}, ball.MotionSliding)
	b := testBall(1, vec.Vec2{X: 200 + 2*cr + gapPx, Y: 300}, vec.Vec2{}, ball.MotionStationary)

	got, ok := nextBallBall(a, b)
	if !ok {
		t.Fatal("expected ball-ball event")
	}

	dist := ball.PxToMeters(gapPx)
	speed := ball.PxToMeters(5 * 60)
	t1, t2, ok := ptmath.SolveQuadratic(-0.5*a.Params.Us*a.Params.G, speed, -dist)
	if !ok {
		t.Fatal("expected quadratic ball-ball root")
	}
	want := math.Min(t1, t2)
	if want <= 0 {
		want = math.Max(t1, t2)
	}
	if math.Abs(got-want) > 1e-6 {
		t.Fatalf("ball-ball time %.9f, want %.9f", got, want)
	}
}

func TestPocketTimeUsesQuarticPrediction(t *testing.T) {
	pocket := table.PocketCenters[0]
	start := pocket.Add(vec.Vec2{X: 42, Y: 0})
	b := testBall(3, start, vec.Vec2{X: -4}, ball.MotionRolling)

	got, ok := nextPocketTime(b)
	if !ok {
		t.Fatal("expected pocket event")
	}
	if got <= 0 || got >= 0.2 {
		t.Fatalf("pocket time %.6f outside expected short approach", got)
	}

	evolved := ballPositionPolynomial(b).at(got)
	end := vec.Vec2{X: ball.MetersToPx(evolved.X), Y: ball.MetersToPx(evolved.Y)}
	if d := end.Sub(pocket).Len(); math.Abs(d-defaultTable.PocketR) > 1e-4 {
		t.Fatalf("predicted pocket distance %.6f, want %.6f", d, defaultTable.PocketR)
	}
}
