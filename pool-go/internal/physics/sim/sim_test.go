package sim_test

import (
	"math"
	"testing"

	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/physics"
	"github.com/user/pooltest/pool-go/internal/table"
	"github.com/user/pooltest/pool-go/internal/vec"
)

func TestRailBounceAtDistance(t *testing.T) {
	cr := ball.CollisionRadius()
	startX := float64(table.PlayRight) - cr - 50
	b := &ball.Ball{Number: 0, Active: true, Params: ball.DefaultParams}
	b.InitFromXY(vec.Vec2{X: startX, Y: 300}, ball.DefaultParams)
	speedPx := 8.0
	b.RVW.V = vec.Vec3{X: ball.PxToMeters(speedPx * 60)}
	b.Motion = ball.MotionSliding

	start := b.Pos.X
	balls := []*ball.Ball{b}
	for range 120 {
		physics.Step(balls)
		if b.Vel.X < 0 {
			break
		}
	}
	travel := b.Pos.X - start
	if travel < 40 || travel > 55 {
		t.Fatalf("rail bounce after %.1f px travel, want ~50 px", travel)
	}
}

func TestBallPassesSidePocketMouth(t *testing.T) {
	cr := ball.CollisionRadius()
	b := &ball.Ball{Number: 0, Active: true, Params: ball.DefaultParams}
	b.InitFromXY(vec.Vec2{X: table.MidX, Y: float64(table.PlayTop) + cr + 5}, ball.DefaultParams)
	b.RVW.V = vec.Vec3{Y: ball.PxToMeters(6 * 60)}
	b.Motion = ball.MotionSliding

	balls := []*ball.Ball{b}
	for range 200 {
		physics.Step(balls)
		if b.Pos.Y > float64(table.PlayBottom)-cr-10 {
			break
		}
	}
	if b.Vel.Y < 0 && math.Abs(b.Pos.X-table.MidX) < 5 {
		t.Fatal("ball bounced off phantom rail at side pocket mouth")
	}
}
