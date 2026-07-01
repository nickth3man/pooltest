// Package physics runs the billiards simulation: ball motion, spin transfer,
// rail and jaw collisions, and pocket capture. It is independent of rendering
// and audio; it reports impact / drop events through callback hooks so the
// orchestrator can wire its own feedback.
package physics

import (
	"math"

	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/table"
	"github.com/user/pooltest/pool-go/internal/vec"
)

const (
	railRestitution = 0.9  // energy kept when bouncing off a cushion
	ballRestitution = 0.96 // energy kept in a ball-to-ball collision
	minSpeed        = 0.05 // speeds below this are snapped to rest
	substeps        = 4    // collision sub-steps per frame to limit tunneling

	// rollDecel is a constant deceleration (px/frame²). Unlike an exponential
	// drag it models real rolling friction: a struck ball travels a predictable
	// distance and comes to a definite stop instead of asymptotically crawling.
	rollDecel = 0.10

	// Spin model. Spin is converted to linear motion by the cloth a little each
	// frame (spinGrip), most strongly once the ball has slowed — so draw/follow
	// reveal themselves after contact rather than during the fast initial roll.
	spinGrip         = 0.10
	spinGripSpeed    = 11.0 // above this speed the cloth barely grips the spin
	spinDecay        = 0.985
	spinRest         = 0.05
	sideDecay        = 0.93
	sideRest         = 0.001
	railSpinTransfer = 0.6 // how much english bends the rebound off a cushion

	// Pocket jaws. Each pocket mouth is a gap in the rail flanked by two cushion
	// tips (pegs); balls entering off-angle rattle against the tips.
	cornerMouth = 30.0
	sideMouth   = 26.0
	jawRadius   = 5.0
)

// Frame FX hooks. They are nil under test and set by the running game so the
// physics layer can stay free of rendering/audio concerns while still reporting
// impacts, rail hits, and pocket drops.
var (
	OnBallImpact func(pos vec.Vec2, strength float64)
	OnRailImpact func(pos vec.Vec2, strength float64)
	OnPocketDrop func(pos vec.Vec2)
)

// Step advances the simulation by one frame and returns the numbers of the
// balls pocketed during the frame (the cue ball is reported as 0).
func Step(balls []*ball.Ball) []int {
	var pocketed []int
	inv := 1.0 / float64(substeps)

	for range substeps {
		for _, b := range balls {
			if !b.Active {
				continue
			}
			b.Pos = b.Pos.Add(b.Vel.Scale(inv))
		}
		// Capture into pockets before resolving rails so a ball entering a
		// pocket mouth drops instead of bouncing off the cushion behind it.
		for _, b := range balls {
			if !b.Active {
				continue
			}
			for _, p := range table.PocketCenters {
				if b.Pos.Sub(p).LenSq() <= table.PocketRadius*table.PocketRadius {
					b.Active = false
					b.Vel = vec.Vec2{}
					b.Spin = vec.Vec2{}
					b.Sinking = true
					b.SinkFrom = b.Pos
					b.SinkPos = p
					b.SinkT = 0
					pocketed = append(pocketed, b.Number)
					if OnPocketDrop != nil {
						OnPocketDrop(p)
					}
					break
				}
			}
		}
		ResolveBallCollisions(balls)
		ResolveRailCollisions(balls)
		ResolveJawCollisions(balls)
	}

	for _, b := range balls {
		if !b.Active {
			continue
		}
		applySpin(b)
		applyFriction(b)
	}
	return pocketed
}

// applySpin lets the cloth convert a ball's stored spin into linear motion and
// advance its visible roll, then decays both. The grip is weakest at high speed
// so follow and draw manifest as the ball slows after contact.
func applySpin(b *ball.Ball) {
	if b.Spin.LenSq() > 0 {
		grip := spinGrip * (1 - math.Min(1, b.Vel.Len()/spinGripSpeed))
		d := b.Spin.Scale(grip)
		b.Vel = b.Vel.Add(d)
		b.Spin = b.Spin.Sub(d).Scale(spinDecay)
		if b.Spin.LenSq() < spinRest*spinRest {
			b.Spin = vec.Vec2{}
		}
	}
	b.Angle += b.SideSpin
	b.SideSpin *= sideDecay
	if math.Abs(b.SideSpin) < sideRest {
		b.SideSpin = 0
	}
}

// applyFriction removes a fixed amount of speed each frame (rolling friction)
// and snaps slow, spin-free balls to a full stop.
func applyFriction(b *ball.Ball) {
	sp := b.Vel.Len()
	if sp <= rollDecel {
		b.Vel = vec.Vec2{}
	} else {
		b.Vel = b.Vel.Sub(b.Vel.Scale(rollDecel / sp))
	}
	if b.Vel.LenSq() < minSpeed*minSpeed && b.Spin.LenSq() == 0 {
		b.Vel = vec.Vec2{}
	}
}

// ResolveBallCollisions handles equal-mass elastic collisions between every
// pair of active balls, separating overlaps and exchanging normal momentum.
func ResolveBallCollisions(balls []*ball.Ball) {
	const minDist = 2 * ball.Radius
	for i := range balls {
		a := balls[i]
		if !a.Active {
			continue
		}
		for j := i + 1; j < len(balls); j++ {
			b := balls[j]
			if !b.Active {
				continue
			}
			d := b.Pos.Sub(a.Pos)
			distSq := d.LenSq()
			if distSq >= minDist*minDist || distSq == 0 {
				continue
			}
			dist := math.Sqrt(distSq)
			n := d.Scale(1 / dist)

			// Push the pair apart so they no longer overlap.
			overlap := minDist - dist
			a.Pos = a.Pos.Sub(n.Scale(overlap / 2))
			b.Pos = b.Pos.Add(n.Scale(overlap / 2))

			// Exchange velocity along the contact normal (equal masses).
			vn := b.Vel.Sub(a.Vel).Dot(n)
			if vn < 0 {
				imp := -(1 + ballRestitution) * vn / 2
				a.Vel = a.Vel.Sub(n.Scale(imp))
				b.Vel = b.Vel.Add(n.Scale(imp))
				if OnBallImpact != nil {
					OnBallImpact(a.Pos.Add(n.Scale(ball.Radius)), -vn)
				}
			}
		}
	}
}

// ResolveRailCollisions keeps ball centers inside the cushion rectangle,
// reflecting velocity off whichever rail was crossed — except across the pocket
// mouths, where balls are allowed to pass through toward the pocket.
func ResolveRailCollisions(balls []*ball.Ball) {
	minX := table.PlayLeft + ball.Radius
	maxX := table.PlayRight - ball.Radius
	minY := table.PlayTop + ball.Radius
	maxY := table.PlayBottom - ball.Radius
	for _, b := range balls {
		if !b.Active {
			continue
		}
		if !InLeftRightMouth(b.Pos.Y) {
			switch {
			case b.Pos.X < minX:
				b.Pos.X = minX
				if b.Vel.X < 0 {
					emitRail(b.Pos, -b.Vel.X)
					b.Vel.X = -b.Vel.X * railRestitution
					b.Vel.Y += b.SideSpin * railSpinTransfer
					b.SideSpin *= 0.5
				}
			case b.Pos.X > maxX:
				b.Pos.X = maxX
				if b.Vel.X > 0 {
					emitRail(b.Pos, b.Vel.X)
					b.Vel.X = -b.Vel.X * railRestitution
					b.Vel.Y -= b.SideSpin * railSpinTransfer
					b.SideSpin *= 0.5
				}
			}
		}
		if !InTopBottomMouth(b.Pos.X) {
			switch {
			case b.Pos.Y < minY:
				b.Pos.Y = minY
				if b.Vel.Y < 0 {
					emitRail(b.Pos, -b.Vel.Y)
					b.Vel.Y = -b.Vel.Y * railRestitution
					b.Vel.X -= b.SideSpin * railSpinTransfer
					b.SideSpin *= 0.5
				}
			case b.Pos.Y > maxY:
				b.Pos.Y = maxY
				if b.Vel.Y > 0 {
					emitRail(b.Pos, b.Vel.Y)
					b.Vel.Y = -b.Vel.Y * railRestitution
					b.Vel.X += b.SideSpin * railSpinTransfer
					b.SideSpin *= 0.5
				}
			}
		}
	}
}

// ResolveJawCollisions bounces balls off the cushion tips flanking each pocket
// mouth, producing the characteristic rattle on off-angle shots.
func ResolveJawCollisions(balls []*ball.Ball) {
	for _, b := range balls {
		if !b.Active {
			continue
		}
		for _, c := range JawPegs {
			if s := ReflectOffPeg(b, c); s > 0 {
				emitRail(b.Pos, s)
			}
		}
	}
}

// ReflectOffPeg treats a cushion tip as a fixed circle and reflects the ball
// off it, returning the impact speed (0 if there was no contact).
func ReflectOffPeg(b *ball.Ball, c vec.Vec2) float64 {
	d := b.Pos.Sub(c)
	distSq := d.LenSq()
	rr := ball.Radius + jawRadius
	if distSq >= rr*rr || distSq == 0 {
		return 0
	}
	dist := math.Sqrt(distSq)
	n := d.Scale(1 / dist)
	b.Pos = c.Add(n.Scale(rr))
	vn := b.Vel.Dot(n)
	if vn >= 0 {
		return 0
	}
	b.Vel = b.Vel.Sub(n.Scale((1 + railRestitution) * vn))
	return -vn
}

func emitRail(pos vec.Vec2, strength float64) {
	if OnRailImpact != nil {
		OnRailImpact(pos, strength)
	}
}

// InTopBottomMouth reports whether x falls within the open mouth of a
// top/bottom pocket (the two corners and the side pocket on that rail).
func InTopBottomMouth(x float64) bool {
	return x <= table.PlayLeft+cornerMouth || x >= table.PlayRight-cornerMouth ||
		math.Abs(x-table.MidX) <= sideMouth
}

// InLeftRightMouth reports whether y falls within the open mouth of a corner
// pocket on the left/right rail.
func InLeftRightMouth(y float64) bool {
	return y <= table.PlayTop+cornerMouth || y >= table.PlayBottom-cornerMouth
}

// SideMouth is exported for tests that need to place balls precisely outside
// the open pocket mouth.
const SideMouth = sideMouth

// JawPegs holds the twelve cushion tips: two flanking each of the six pockets.
// The positions are fixed, so they are computed once.
var JawPegs = []vec.Vec2{
	{X: table.PlayLeft + cornerMouth, Y: table.PlayTop},
	{X: table.PlayLeft, Y: table.PlayTop + cornerMouth},
	{X: table.PlayRight - cornerMouth, Y: table.PlayTop},
	{X: table.PlayRight, Y: table.PlayTop + cornerMouth},
	{X: table.PlayLeft + cornerMouth, Y: table.PlayBottom},
	{X: table.PlayLeft, Y: table.PlayBottom - cornerMouth},
	{X: table.PlayRight - cornerMouth, Y: table.PlayBottom},
	{X: table.PlayRight, Y: table.PlayBottom - cornerMouth},
	{X: table.MidX - sideMouth, Y: table.PlayTop},
	{X: table.MidX + sideMouth, Y: table.PlayTop},
	{X: table.MidX - sideMouth, Y: table.PlayBottom},
	{X: table.MidX + sideMouth, Y: table.PlayBottom},
}

// AllStopped reports whether every ball has come to rest.
func AllStopped(balls []*ball.Ball) bool {
	for _, b := range balls {
		if b.Moving() {
			return false
		}
	}
	return true
}

// RayToRail returns where a ray from start in direction dir meets the cushion
// rectangle (inset by the ball radius), used to draw the aim line when the
// shot will not hit a ball.
func RayToRail(start, dir vec.Vec2) vec.Vec2 {
	minX := table.PlayLeft + ball.Radius
	maxX := table.PlayRight - ball.Radius
	minY := table.PlayTop + ball.Radius
	maxY := table.PlayBottom - ball.Radius
	t := math.Inf(1)
	if dir.X > 1e-9 {
		t = math.Min(t, (maxX-start.X)/dir.X)
	} else if dir.X < -1e-9 {
		t = math.Min(t, (minX-start.X)/dir.X)
	}
	if dir.Y > 1e-9 {
		t = math.Min(t, (maxY-start.Y)/dir.Y)
	} else if dir.Y < -1e-9 {
		t = math.Min(t, (minY-start.Y)/dir.Y)
	}
	if math.IsInf(t, 1) {
		return start
	}
	return start.Add(dir.Scale(t))
}
