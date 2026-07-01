// Package physics runs the billiards simulation via the event-based engine in sim/.
package physics

import (
	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/physics/sim"
	"github.com/user/pooltest/pool-go/internal/vec"
)

// Frame duration at 60 TPS (seconds).
const frameDT = 1.0 / 60.0

// Step advances the simulation by one frame and returns pocketed ball numbers.
func Step(balls []*ball.Ball) []int {
	return sim.Advance(balls, frameDT)
}

// SetHooks assigns impact callbacks (called from game.NewGame).
func SetHooks(onBall, onRail func(vec.Vec2, float64), onPocket func(vec.Vec2)) {
	sim.OnBallImpact = onBall
	sim.OnRailImpact = onRail
	sim.OnPocketDrop = onPocket
}

// AllStopped reports whether every ball has come to rest.
func AllStopped(balls []*ball.Ball) bool {
	return sim.AllStopped(balls)
}

// RayToRail returns where a ray meets the cushion rectangle.
func RayToRail(start, dir vec.Vec2) vec.Vec2 {
	return sim.RayToRail(start, dir)
}

// EstimateThrowDeg returns approximate object-ball throw for aim preview (degrees).
func EstimateThrowDeg(cutAngle, speedFrac, sideSpin, followDraw float64) float64 {
	return sim.EstimateThrowDeg(cutAngle, speedFrac, sideSpin, followDraw)
}
