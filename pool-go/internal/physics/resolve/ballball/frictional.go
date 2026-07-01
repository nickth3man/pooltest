package ballball

import (
	"math"

	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/physics/ptmath"
	"github.com/user/pooltest/pool-go/internal/vec"
)

// Alciatore friction fit: u_b = a + b*exp(-c*v_rel).
const (
	fricA = 9.951e-3
	fricB = 0.108
	fricC = 1.088
)

// AlciatoreMu returns ball-ball friction from relative surface speed (m/s).
func AlciatoreMu(vRel float64) float64 {
	return fricA + fricB*math.Exp(-fricC*vRel)
}

// Resolve applies frictional inelastic equal-mass collision (pooltool).
func Resolve(a, b *ball.Ball) {
	p := a.Params
	if b.Params.R > 0 {
		p = b.Params
	}
	r1, r2 := a.RVW, b.RVW
	r1out, r2out := resolvePair(r1, r2, p.R, frictionBetween(a, b), p.Eb)
	a.RVW, b.RVW = r1out, r2out
	a.Motion, b.Motion = ball.MotionSliding, ball.MotionSliding
}

func frictionBetween(a, b *ball.Ball) float64 {
	n := ptmath.Unit3(b.RVW.R.Sub(a.RVW.R))
	v1 := ptmath.TangentSurfaceVelocity(ptmath.RVWFrom(a.RVW.R, a.RVW.V, a.RVW.W), n, a.Params.R)
	v2 := ptmath.TangentSurfaceVelocity(ptmath.RVWFrom(b.RVW.R, b.RVW.V, b.RVW.W), n.Scale(-1), b.Params.R)
	vRel := ptmath.Norm3(v1.Sub(v2))
	return AlciatoreMu(vRel)
}

func resolvePair(rvw1, rvw2 ball.RVW, radius, ub, eb float64) (ball.RVW, ball.RVW) {
	unitX := vec.Vec3{X: 1}
	delta := rvw2.R.Sub(rvw1.R)
	theta := ptmath.Angle(delta, unitX)

	rot := func(v vec.Vec3) vec.Vec3 { return ptmath.CoordinateRotation(v, -theta) }
	unrot := func(v vec.Vec3) vec.Vec3 { return ptmath.CoordinateRotation(v, theta) }

	v1, w1 := rot(rvw1.V), rot(rvw1.W)
	v2, w2 := rot(rvw2.V), rot(rvw2.W)

	v1nF := 0.5 * ((1-eb)*v1.X + (1+eb)*v2.X)
	v2nF := 0.5 * ((1+eb)*v1.X + (1-eb)*v2.X)
	w1nF, w2nF := w1.X, w2.X

	v1 = vec.Vec3{Y: v1.Y, Z: v1.Z}
	v2 = vec.Vec3{Y: v2.Y, Z: v2.Z}
	w1 = vec.Vec3{Y: w1.Y, Z: w1.Z}
	w2 = vec.Vec3{Y: w2.Y, Z: w2.Z}

	rvw1f := ball.RVW{R: rvw1.R, V: v1, W: w1}
	rvw2f := ball.RVW{R: rvw2.R, V: v2, W: w2}

	v1c := ptmath.SurfaceVelocity(ptmath.RVWFrom(rvw1f.R, rvw1f.V, rvw1f.W), unitX, radius)
	v2c := ptmath.SurfaceVelocity(ptmath.RVWFrom(rvw2f.R, rvw2f.V, rvw2f.W), unitX.Scale(-1), radius)
	v12c := v1c.Sub(v2c)
	hasRel := ptmath.Norm3(v12c) > ptmath.Eps

	if hasRel {
		v12hat := ptmath.Unit3(v12c)
		dvn := math.Abs(v2nF - v1nF)
		dv1t := v12hat.Scale(-ub * dvn)
		dw1 := ptmath.Cross(unitX, dv1t).Scale(2.5 / radius)
		rvw1f.V = rvw1f.V.Add(dv1t)
		rvw1f.W = rvw1f.W.Add(dw1)
		rvw2f.V = rvw2f.V.Sub(dv1t)
		rvw2f.W = rvw2f.W.Add(dw1)

		v1slip := ptmath.SurfaceVelocity(ptmath.RVWFrom(rvw1f.R, rvw1f.V, rvw1f.W), unitX, radius)
		v2slip := ptmath.SurfaceVelocity(ptmath.RVWFrom(rvw2f.R, rvw2f.V, rvw2f.W), unitX.Scale(-1), radius)
		v12slip := v1slip.Sub(v2slip)
		if !hasRel || v12c.Dot(v12slip) <= 0 {
			hasRel = false
		}
	}
	if !hasRel {
		wSum := w1.Add(w2)
		crossW := ptmath.Cross(wSum, unitX)
		dv1t := v1.Sub(v2).Add(vec.Vec3{Y: crossW.Y * radius, Z: crossW.Z * radius}).Scale(-1.0 / 7.0)
		dw1 := ptmath.Cross(unitX, v1.Sub(v2).Scale(1/radius)).Scale(-5.0 / 14.0).Sub(wSum.Scale(5.0 / 14.0))
		rvw1f.V = rvw1f.V.Add(dv1t)
		rvw1f.W = rvw1f.W.Add(dw1)
		rvw2f.V = rvw2f.V.Sub(dv1t)
		rvw2f.W = rvw2f.W.Add(dw1)
	}

	rvw1f.V.X = v1nF
	rvw2f.V.X = v2nF
	rvw1f.W.X = w1nF
	rvw2f.W.X = w2nF

	rvw1f.V = unrot(rvw1f.V)
	rvw1f.W = unrot(rvw1f.W)
	rvw2f.V = unrot(rvw2f.V)
	rvw2f.W = unrot(rvw2f.W)
	rvw1f.V.Z = 0
	rvw2f.V.Z = 0

	// Separate overlap along line of centers.
	n := ptmath.Unit3(delta)
	overlap := 2*radius - ptmath.Norm3(delta)
	if overlap > 0 {
		shift := n.Scale(overlap / 2)
		rvw1f.R = rvw1f.R.Sub(shift)
		rvw2f.R = rvw2f.R.Add(shift)
	}

	return rvw1f, rvw2f
}
