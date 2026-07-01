package stickball

import (
	"math"

	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/physics/evolve"
	"github.com/user/pooltest/pool-go/internal/physics/ptmath"
	"github.com/user/pooltest/pool-go/internal/vec"
)

// TipOffset is normalized cue contact on ball surface (a=side, b=vertical).
type TipOffset struct {
	A, B float64 // side, follow(+)/draw(-)
}

// Strike applies an instantaneous cue strike to the cue ball.
func Strike(b *ball.Ball, phiRad, v0 float64, tip TipOffset) {
	p := b.Params
	a, bb := tip.A, tip.B
	thetaRad := p.CueElevation
	cueC := math.Sqrt(math.Max(0, 1-a*a-bb*bb))
	ballA := a
	ballC := math.Cos(thetaRad)*cueC - math.Sin(thetaRad)*bb
	ballB := math.Sin(thetaRad)*cueC + math.Cos(thetaRad)*bb

	v, w := cueStrike(p.M, p.CueMass, p.R, v0, phiRad, thetaRad, ballA, ballB, ballC)
	alpha := squirtAngle(p.M, p.CueEndMass, ballA, 1.0)
	v = ptmath.CoordinateRotation(v, alpha)

	b.RVW.V = v
	b.RVW.W = w
	b.RVW.R.Z = p.R
	b.Motion = evolve.FinalMotionState(b.RVW)
}

func cueStrike(m, cueMass, radius, v0, phi, theta, a, b, c float64) (vec.Vec3, vec.Vec3) {
	inertia := 2.0 / 5.0 * radius * radius
	num := 2 * v0
	temp := a*a + math.Pow(b*math.Cos(theta), 2) + math.Pow(c*math.Sin(theta), 2) - 2*b*c*math.Cos(theta)*math.Sin(theta)
	den := 1 + m/cueMass + temp/inertia
	vMag := num / den

	vB := vec.Vec3{Y: -vMag * math.Cos(theta), Z: vMag * math.Sin(theta)}
	vecX := -c*math.Sin(theta) + b*math.Cos(theta)
	vecY := a * math.Sin(theta)
	vecZ := -a * math.Cos(theta)
	spinVec := vec.Vec3{X: vecX, Y: vecY, Z: vecZ}
	wB := spinVec.Scale(vMag / inertia)

	rot := phi + math.Pi/2
	vT := ptmath.CoordinateRotation(vB, rot)
	wT := ptmath.CoordinateRotation(wB, rot)
	return vT, wT
}

func squirtAngle(mB, mE, a, throttle float64) float64 {
	mR := mB / mE
	A := 1 - a*a
	num := 5.0 / 2.0 * a * math.Sqrt(A)
	den := 1 + mR + 5.0/2.0*A
	return -throttle * math.Atan2(num, den)
}

// StrikeFromGame converts game dial + pull direction into a strike.
func StrikeFromGame(b *ball.Ball, dir vec.Vec2, speedPxPerFrame, side, followDraw float64) {
	phi := math.Atan2(dir.Y, dir.X)
	v0 := ball.PxToMeters(speedPxPerFrame * 60) // px/frame -> m/s at 60 TPS
	Strike(b, phi, v0, TipOffset{A: side, B: followDraw})
}
