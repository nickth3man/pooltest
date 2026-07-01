package sim

import (
	"math"

	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/physics/ptmath"
	"github.com/user/pooltest/pool-go/internal/table"
	"github.com/user/pooltest/pool-go/internal/vec"
)

type positionPolynomial struct {
	c0 vec.Vec3
	c1 vec.Vec3
	c2 vec.Vec3
}

func ballPositionPolynomial(b *ball.Ball) positionPolynomial {
	p := positionPolynomial{c0: b.RVW.R, c1: b.RVW.V}
	switch b.Motion {
	case ball.MotionRolling:
		p.c2 = ptmath.Unit3(b.RVW.V).Scale(-0.5 * b.Params.Ur * b.Params.G)
	case ball.MotionSliding:
		rel := ptmath.RelVelocity(ptmath.RVWFrom(b.RVW.R, b.RVW.V, b.RVW.W), b.Params.R)
		p.c2 = ptmath.Unit3(rel).Scale(-0.5 * b.Params.Us * b.Params.G)
	}
	return p
}

func (p positionPolynomial) at(t float64) vec.Vec3 {
	return p.c0.Add(p.c1.Scale(t)).Add(p.c2.Scale(t * t))
}

func (p positionPolynomial) sub(o positionPolynomial) positionPolynomial {
	return positionPolynomial{
		c0: p.c0.Sub(o.c0),
		c1: p.c1.Sub(o.c1),
		c2: p.c2.Sub(o.c2),
	}
}

func nextTransition(b *ball.Ball) float64 {
	switch b.Motion {
	case ball.MotionSliding:
		return ptmath.GetSlideTime(ptmath.RVWFrom(b.RVW.R, b.RVW.V, b.RVW.W), b.Params.R, b.Params.Us, b.Params.G)
	case ball.MotionRolling:
		return ptmath.GetRollTime(ptmath.RVWFrom(b.RVW.R, b.RVW.V, b.RVW.W), b.Params.Ur, b.Params.G)
	case ball.MotionSpinning:
		return ptmath.GetSpinTime(ptmath.RVWFrom(b.RVW.R, b.RVW.V, b.RVW.W), b.Params.R, b.Params.USp(), b.Params.G)
	default:
		return math.Inf(1)
	}
}

func translating(b *ball.Ball) bool {
	return b.Active && (b.Motion == ball.MotionSliding || b.Motion == ball.MotionRolling)
}

func currentPosPx(b *ball.Ball) vec.Vec2 {
	return vec.Vec2{X: ball.MetersToPx(b.RVW.R.X), Y: ball.MetersToPx(b.RVW.R.Y)}
}

func currentSpeedPxFrame(b *ball.Ball) float64 {
	return ball.MetersToPx(ptmath.Norm3(b.RVW.V)) / 60
}

func earliestSphereTime(p positionPolynomial, radius float64) (float64, bool) {
	a := p.c2
	b := p.c1
	c := p.c0
	roots := ptmath.SolveQuartic(
		a.Dot(a),
		2*a.Dot(b),
		b.Dot(b)+2*a.Dot(c),
		2*b.Dot(c),
		c.Dot(c)-radius*radius,
	)
	candidate, ok := ptmath.EarliestPositiveRoot(roots)
	if !ok {
		candidate = 30
	}
	if t, refined := firstBracketedSphereRoot(p, radius, candidate); refined {
		return t, true
	}
	return candidate, ok
}

func firstBracketedSphereRoot(p positionPolynomial, radius, maxT float64) (float64, bool) {
	if maxT <= ptmath.Eps || math.IsInf(maxT, 0) || math.IsNaN(maxT) {
		return 0, false
	}
	f := func(t float64) float64 {
		pos := p.at(t)
		return pos.Dot(pos) - radius*radius
	}

	prevT := 0.0
	prev := f(prevT)
	const samples = 256
	for i := 1; i <= samples; i++ {
		t := maxT * float64(i) / samples
		cur := f(t)
		switch {
		case math.Abs(cur) <= 1e-10 && t > ptmath.Eps:
			return t, true
		case prev > 0 && cur < 0:
			return bisectRoot(f, prevT, t), true
		case prev < 0 && cur > 0:
			return prevT, true
		}
		prevT, prev = t, cur
	}
	return 0, false
}

func bisectRoot(f func(float64) float64, lo, hi float64) float64 {
	for range 64 {
		mid := (lo + hi) / 2
		if f(mid) <= 0 {
			hi = mid
		} else {
			lo = mid
		}
	}
	return hi
}

func nextBallBall(a, b *ball.Ball) (float64, bool) {
	if !translating(a) && !translating(b) {
		return math.Inf(1), false
	}
	delta := ballPositionPolynomial(a).sub(ballPositionPolynomial(b))
	return earliestSphereTime(delta, a.Params.R+b.Params.R)
}

func nextCircularCushion(b *ball.Ball, seg table.CushionSegment) float64 {
	if !translating(b) {
		return math.Inf(1)
	}
	center := vec.Vec3{
		X: ball.PxToMeters(seg.Center.X),
		Y: ball.PxToMeters(seg.Center.Y),
		Z: b.Params.R,
	}
	p := ballPositionPolynomial(b)
	p.c0 = p.c0.Sub(center)
	t, ok := earliestSphereTime(p, b.Params.R+ball.PxToMeters(seg.Radius))
	if !ok {
		return math.Inf(1)
	}
	return t
}

func nextPocketTime(b *ball.Ball) (float64, bool) {
	if !translating(b) {
		return math.Inf(1), false
	}
	best := math.Inf(1)
	for _, pocket := range defaultTable.PocketPts {
		center := vec.Vec3{
			X: ball.PxToMeters(pocket.X),
			Y: ball.PxToMeters(pocket.Y),
			Z: b.Params.R,
		}
		p := ballPositionPolynomial(b)
		p.c0 = p.c0.Sub(center)
		if t, ok := earliestSphereTime(p, ball.PxToMeters(defaultTable.PocketR)); ok && t < best {
			best = t
		}
	}
	return best, best < math.Inf(1)
}

func nextLinearCushion(b *ball.Ball, seg table.CushionSegment) (float64, bool) {
	if !translating(b) {
		return math.Inf(1), false
	}

	poly := ballPositionPolynomial(b)
	p1 := vec.Vec3{X: ball.PxToMeters(seg.P1.X), Y: ball.PxToMeters(seg.P1.Y), Z: b.Params.R}
	normal := vec.Vec3{X: seg.Normal.X, Y: seg.Normal.Y}

	a := normal.Dot(poly.c2)
	bb := normal.Dot(poly.c1)
	c := normal.Dot(poly.c0.Sub(p1)) - b.Params.R

	t1, t2, ok := ptmath.SolveQuadratic(a, bb, c)
	if !ok {
		return math.Inf(1), false
	}

	best := math.Inf(1)
	for _, t := range []float64{t1, t2} {
		if t <= ptmath.Eps || t >= best {
			continue
		}
		if linearHitInSegment(poly.at(t), seg) {
			best = t
		}
	}
	return best, best < math.Inf(1)
}

func linearHitInSegment(pos vec.Vec3, seg table.CushionSegment) bool {
	hit := vec.Vec2{X: ball.MetersToPx(pos.X), Y: ball.MetersToPx(pos.Y)}
	axis := seg.P2.Sub(seg.P1)
	den := axis.LenSq()
	if den == 0 {
		return false
	}
	s := hit.Sub(seg.P1).Dot(axis) / den
	const tol = 1e-7
	return s >= -tol && s <= 1+tol
}
