// Package sim runs event-based billiards simulation (ported from pooltool).
package sim

import (
	"math"

	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/physics/evolve"
	"github.com/user/pooltest/pool-go/internal/physics/ptmath"
	"github.com/user/pooltest/pool-go/internal/physics/resolve/ballball"
	"github.com/user/pooltest/pool-go/internal/physics/resolve/ballcushion"
	"github.com/user/pooltest/pool-go/internal/table"
	"github.com/user/pooltest/pool-go/internal/vec"
)

const maxEventsPerFrame = 200

const spinStopEps = 2e-2

// Hooks mirror legacy physics callbacks.
var (
	OnBallImpact func(pos vec.Vec2, strength float64)
	OnRailImpact func(pos vec.Vec2, strength float64)
	OnPocketDrop func(pos vec.Vec2)
)

var defaultTable = table.DefaultPhysicsTable()

// Advance simulates dt seconds and returns pocketed ball numbers (0 = cue).
type eventAux struct {
	i      int
	j      int
	nx, ny float64
	cx, cy float64
}

func Advance(balls []*ball.Ball, dt float64) []int {
	var pocketed []int
	remaining := dt
	events := 0

	for remaining > 0 && events < maxEventsPerFrame {
		for i, b := range balls {
			if b.Active && pocketOne(b) {
				pocketed = append(pocketed, b.Number)
				if OnPocketDrop != nil {
					OnPocketDrop(b.Pos)
				}
				_ = i
			}
		}
		if allStopped(balls) {
			break
		}
		next, kind, aux := nextEvent(balls, remaining)
		if next <= 0 || math.IsInf(next, 1) {
			evolveAll(balls, remaining)
			break
		}
		if next > remaining {
			next = remaining
			kind = eventEvolve
		}
		evolveAll(balls, next)
		remaining -= next

		switch kind {
		case eventBallBall:
			strength := ptmath.Norm3(balls[aux.j].RVW.V.Sub(balls[aux.i].RVW.V))
			ballball.Resolve(balls[aux.i], balls[aux.j])
			if OnBallImpact != nil {
				mid := currentPosPx(balls[aux.i]).Add(currentPosPx(balls[aux.j])).Scale(0.5)
				OnBallImpact(mid, ball.MetersToPx(strength)*60)
			}
		case eventRail:
			ballcushion.ResolveLinear(balls[aux.i], vec.Vec2{X: aux.nx, Y: aux.ny})
			if OnRailImpact != nil {
				OnRailImpact(currentPosPx(balls[aux.i]), currentSpeedPxFrame(balls[aux.i]))
			}
		case eventJaw:
			ballcushion.ResolveCircular(balls[aux.i], vec.Vec2{X: aux.cx, Y: aux.cy})
			if OnRailImpact != nil {
				OnRailImpact(currentPosPx(balls[aux.i]), currentSpeedPxFrame(balls[aux.i]))
			}
		case eventPocket:
			b := balls[aux.i]
			if pocketOne(b) {
				pocketed = append(pocketed, b.Number)
				if OnPocketDrop != nil {
					OnPocketDrop(b.Pos)
				}
			}
		}
		events++
	}

	for _, b := range balls {
		b.SyncXY(dt)
	}
	return pocketed
}

type eventKind int

const (
	eventEvolve eventKind = iota
	eventBallBall
	eventRail
	eventJaw
	eventPocket
)

func evolveAll(balls []*ball.Ball, t float64) {
	for _, b := range balls {
		if !b.Active || b.Motion == ball.MotionPocketed {
			continue
		}
		rvw, state := evolve.EvolveBallMotion(b.Motion, b.RVW, b.Params, t)
		b.RVW, b.Motion = rvw, state
		if ptmath.Norm3(b.RVW.V) < ptmath.Eps && math.Abs(b.RVW.W.Z) < spinStopEps {
			b.Motion = ball.MotionStationary
			b.RVW.V = vec.Vec3{}
			b.RVW.W = vec.Vec3{}
		}
	}
}

func nextEvent(balls []*ball.Ball, horizon float64) (float64, eventKind, eventAux) {
	best := horizon
	kind := eventEvolve
	var aux eventAux

	for i := range balls {
		b := balls[i]
		if !b.Active {
			continue
		}
		if t := nextTransition(b); t > ptmath.Eps && t < best {
			best, kind, aux = t, eventEvolve, eventAux{i: i}
		}
		if t, pk := nextPocket(b); pk && t < best {
			best, kind, aux = t, eventPocket, eventAux{i: i}
		}
		if t, n := nextRail(b); t > 0 && t < best {
			best, kind, aux = t, eventRail, eventAux{i: i, nx: n.X, ny: n.Y}
		}
		for _, seg := range defaultTable.Circular {
			if t := nextJaw(b, seg); t > 0 && t < best {
				best, kind, aux = t, eventJaw, eventAux{i: i, cx: seg.Center.X, cy: seg.Center.Y}
			}
		}
	}

	for i := range balls {
		for j := i + 1; j < len(balls); j++ {
			a, bb := balls[i], balls[j]
			if !a.Active || !bb.Active {
				continue
			}
			if t, ok := nextBallBall(a, bb); ok && t > 0 && t < best {
				best, kind, aux = t, eventBallBall, eventAux{i: i, j: j}
			}
		}
	}
	return best, kind, aux
}

func nextPocket(b *ball.Ball) (float64, bool) {
	pos := currentPosPx(b)
	for _, p := range defaultTable.PocketPts {
		if pos.Sub(p).LenSq() <= defaultTable.PocketR*defaultTable.PocketR {
			return ptmath.Eps, true
		}
	}
	return nextPocketTime(b)
}

func nextRail(b *ball.Ball) (float64, vec.Vec2) {
	best := math.Inf(1)
	var normal vec.Vec2
	for _, seg := range defaultTable.Linear {
		if t, ok := nextLinearCushion(b, seg); ok && t < best {
			best, normal = t, seg.Normal
		}
	}
	if math.IsInf(best, 1) {
		return 0, vec.Vec2{}
	}
	return best, normal
}

func nextJaw(b *ball.Ball, seg table.CushionSegment) float64 {
	return nextCircularCushion(b, seg)
}

func pocketOne(b *ball.Ball) bool {
	pos := currentPosPx(b)
	for _, p := range defaultTable.PocketPts {
		if pos.Sub(p).LenSq() <= defaultTable.PocketR*defaultTable.PocketR {
			b.Active = false
			b.Motion = ball.MotionPocketed
			b.RVW.V = vec.Vec3{}
			b.RVW.W = vec.Vec3{}
			b.Sinking = true
			b.Pos = pos
			b.SinkFrom = pos
			b.SinkPos = p
			b.SinkT = 0
			return true
		}
	}
	return false
}

// AllStopped reports whether every ball has come to rest.
func AllStopped(balls []*ball.Ball) bool {
	return allStopped(balls)
}

func allStopped(balls []*ball.Ball) bool {
	for _, b := range balls {
		if b.Moving() {
			return false
		}
	}
	return true
}

// RayToRail delegates to the physics table.
func RayToRail(start, dir vec.Vec2) vec.Vec2 {
	return defaultTable.RayToRail(start, dir)
}

// EstimateThrowDeg provides aim-preview throw estimate (degrees).
func EstimateThrowDeg(cutAngle, speedFrac, sideSpin, followDraw float64) float64 {
	// Use Alciatore friction at representative speed.
	speed := 3 + speedFrac*9
	mu := ballball.AlciatoreMu(speed * 0.05)
	throw := mu / 0.06 * 5
	if followDraw != 0 {
		throw *= 0.35
	}
	if sideSpin > 0 {
		throw -= sideSpin * 1.5
	} else if sideSpin < 0 {
		throw += math.Abs(sideSpin) * 0.8
	}
	if throw > 6 {
		throw = 6
	}
	if throw < -6 {
		throw = -6
	}
	return throw
}
