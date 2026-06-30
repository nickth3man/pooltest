package main

import "math"

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
	onBallImpact func(pos Vec2, strength float64)
	onRailImpact func(pos Vec2, strength float64)
	onPocketDrop func(pos Vec2)
)

// step advances the simulation by one frame and returns the numbers of the
// balls pocketed during the frame (the cue ball is reported as 0).
func step(balls []*Ball) []int {
	var pocketed []int
	inv := 1.0 / float64(substeps)
	ps := pockets()

	for s := 0; s < substeps; s++ {
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
			for _, p := range ps {
				if b.Pos.Sub(p).LenSq() <= pocketRadius*pocketRadius {
					b.Active = false
					b.Vel = Vec2{}
					b.Spin = Vec2{}
					b.Sinking = true
					b.SinkFrom = b.Pos
					b.SinkPos = p
					b.SinkT = 0
					pocketed = append(pocketed, b.Number)
					if onPocketDrop != nil {
						onPocketDrop(p)
					}
					break
				}
			}
		}
		resolveBallCollisions(balls)
		resolveRailCollisions(balls)
		resolveJawCollisions(balls)
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
func applySpin(b *Ball) {
	if b.Spin.LenSq() > 0 {
		grip := spinGrip * (1 - math.Min(1, b.Vel.Len()/spinGripSpeed))
		d := b.Spin.Scale(grip)
		b.Vel = b.Vel.Add(d)
		b.Spin = b.Spin.Sub(d).Scale(spinDecay)
		if b.Spin.LenSq() < spinRest*spinRest {
			b.Spin = Vec2{}
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
func applyFriction(b *Ball) {
	sp := b.Vel.Len()
	if sp <= rollDecel {
		b.Vel = Vec2{}
	} else {
		b.Vel = b.Vel.Sub(b.Vel.Scale(rollDecel / sp))
	}
	if b.Vel.LenSq() < minSpeed*minSpeed && b.Spin.LenSq() == 0 {
		b.Vel = Vec2{}
	}
}

// resolveBallCollisions handles equal-mass elastic collisions between every
// pair of active balls, separating overlaps and exchanging normal momentum.
func resolveBallCollisions(balls []*Ball) {
	const minDist = 2 * ballRadius
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
				if onBallImpact != nil {
					onBallImpact(a.Pos.Add(n.Scale(ballRadius)), -vn)
				}
			}
		}
	}
}

// resolveRailCollisions keeps ball centers inside the cushion rectangle,
// reflecting velocity off whichever rail was crossed — except across the pocket
// mouths, where balls are allowed to pass through toward the pocket.
func resolveRailCollisions(balls []*Ball) {
	minX := playLeft + ballRadius
	maxX := playRight - ballRadius
	minY := playTop + ballRadius
	maxY := playBottom - ballRadius
	for _, b := range balls {
		if !b.Active {
			continue
		}
		if !inLeftRightMouth(b.Pos.Y) {
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
		if !inTopBottomMouth(b.Pos.X) {
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

// resolveJawCollisions bounces balls off the cushion tips flanking each pocket
// mouth, producing the characteristic rattle on off-angle shots.
func resolveJawCollisions(balls []*Ball) {
	pegs := jawPegs()
	for _, b := range balls {
		if !b.Active {
			continue
		}
		for _, c := range pegs {
			if s := reflectOffPeg(b, c); s > 0 {
				emitRail(b.Pos, s)
			}
		}
	}
}

// reflectOffPeg treats a cushion tip as a fixed circle and reflects the ball off
// it, returning the impact speed (0 if there was no contact).
func reflectOffPeg(b *Ball, c Vec2) float64 {
	d := b.Pos.Sub(c)
	distSq := d.LenSq()
	rr := ballRadius + jawRadius
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

func emitRail(pos Vec2, strength float64) {
	if onRailImpact != nil {
		onRailImpact(pos, strength)
	}
}

// inTopBottomMouth reports whether x falls within the open mouth of a top/bottom
// pocket (the two corners and the side pocket on that rail).
func inTopBottomMouth(x float64) bool {
	midX := float64(playLeft+playRight) / 2
	return x <= playLeft+cornerMouth || x >= playRight-cornerMouth ||
		math.Abs(x-midX) <= sideMouth
}

// inLeftRightMouth reports whether y falls within the open mouth of a corner
// pocket on the left/right rail.
func inLeftRightMouth(y float64) bool {
	return y <= playTop+cornerMouth || y >= playBottom-cornerMouth
}

// jawPegs returns the twelve cushion tips: two flanking each of the six pockets.
func jawPegs() []Vec2 {
	midX := float64(playLeft+playRight) / 2
	return []Vec2{
		{playLeft + cornerMouth, playTop},
		{playLeft, playTop + cornerMouth},
		{playRight - cornerMouth, playTop},
		{playRight, playTop + cornerMouth},
		{playLeft + cornerMouth, playBottom},
		{playLeft, playBottom - cornerMouth},
		{playRight - cornerMouth, playBottom},
		{playRight, playBottom - cornerMouth},
		{midX - sideMouth, playTop},
		{midX + sideMouth, playTop},
		{midX - sideMouth, playBottom},
		{midX + sideMouth, playBottom},
	}
}

// allStopped reports whether every ball has come to rest.
func allStopped(balls []*Ball) bool {
	for _, b := range balls {
		if b.moving() {
			return false
		}
	}
	return true
}

// firstContact walks the aim ray from the cue ball and returns the point at
// which it would first touch another ball, that ball, and whether a hit exists.
func (g *Game) firstContact(dir Vec2) (Vec2, *Ball, bool) {
	best := math.Inf(1)
	var target *Ball
	rad := 2 * ballRadius
	for _, b := range g.balls {
		if b == g.cue || !b.Active {
			continue
		}
		f := b.Pos.Sub(g.cue.Pos)
		proj := f.Dot(dir)
		if proj < 0 {
			continue // behind the shot
		}
		perpSq := f.LenSq() - proj*proj
		if perpSq > rad*rad {
			continue // ray misses this ball
		}
		t := proj - math.Sqrt(rad*rad-perpSq)
		if t < best {
			best, target = t, b
		}
	}
	if target == nil {
		return Vec2{}, nil, false
	}
	return g.cue.Pos.Add(dir.Scale(best)), target, true
}

// rayToRail returns where a ray from start in direction dir meets the cushion
// rectangle (inset by the ball radius), used to draw the aim line when the shot
// will not hit a ball.
func rayToRail(start, dir Vec2) Vec2 {
	minX := playLeft + ballRadius
	maxX := playRight - ballRadius
	minY := playTop + ballRadius
	maxY := playBottom - ballRadius
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
