package physics

import (
	"math"

	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/physics/resolve/stickball"
	"github.com/user/pooltest/pool-go/internal/vec"
)

// PredictThrowDeg simulates a single ball-ball hit and returns object-ball throw (degrees).
func PredictThrowDeg(strikerPos, targetPos, aimDir vec.Vec2, speedPxPerFrame, sideSpin, followDraw float64) float64 {
	dir := aimDir.Normalize()
	a := &ball.Ball{Number: 0, Active: true, Params: ball.DefaultParams}
	b := &ball.Ball{Number: 1, Active: true, Params: ball.DefaultParams}
	a.InitFromXY(strikerPos, ball.DefaultParams)
	b.InitFromXY(targetPos, ball.DefaultParams)

	ghost := a.Pos.Add(dir.Scale(a.Pos.Sub(targetPos).Len() - 2*ball.CollisionRadius()))
	a.SetPosPx(ghost)
	stickball.StrikeFromGame(a, dir, speedPxPerFrame, sideSpin, followDraw)

	ResolveBallCollisions([]*ball.Ball{a, b})
	b.SyncXY(1.0 / 60.0)
	n := targetPos.Sub(ghost).Normalize()
	if n.LenSq() == 0 {
		return 0
	}
	dot := n.Dot(b.Vel.Normalize())
	if dot > 1 {
		dot = 1
	}
	if dot < -1 {
		dot = -1
	}
	return math.Acos(dot) * 180 / math.Pi
}

// PreviewCBPath returns cue-ball positions along a short forward integration for aim preview.
func PreviewCBPath(start vec.Vec2, dir vec.Vec2, speedPxPerFrame, sideSpin, followDraw float64, steps int) []vec.Vec2 {
	b := &ball.Ball{Number: 0, Active: true, Params: ball.DefaultParams}
	b.InitFromXY(start, ball.DefaultParams)
	stickball.StrikeFromGame(b, dir.Normalize(), speedPxPerFrame, sideSpin, followDraw)
	balls := []*ball.Ball{b}
	pts := []vec.Vec2{start}
	for range steps {
		Step(balls)
		pts = append(pts, b.Pos)
		if AllStopped(balls) {
			break
		}
	}
	return pts
}
