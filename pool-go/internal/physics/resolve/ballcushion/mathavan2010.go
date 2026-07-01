package ballcushion

import (
	"math"

	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/physics/ptmath"
	"github.com/user/pooltest/pool-go/internal/vec"
)

// UseMathavan2010 selects the Mathavan cushion model when true (pooltool default).
var UseMathavan2010 = true

// ResolveLinear reflects a ball off a cushion with outward normal n (unit, table xy).
func ResolveLinear(b *ball.Ball, n vec.Vec2) {
	if UseMathavan2010 {
		resolveMathavanLinear(b, n)
		return
	}
	resolveHanLinear(b, n)
}

func resolveHanLinear(b *ball.Ball, n vec.Vec2) {
	p := b.Params
	xyNormal := vec.Vec3{X: n.X, Y: n.Y}
	if b.RVW.V.X*xyNormal.X+b.RVW.V.Y*xyNormal.Y >= 0 {
		vn := b.RVW.V.X*xyNormal.X + b.RVW.V.Y*xyNormal.Y
		if vn >= 0 {
			return
		}
	}
	b.RVW = han2005(b.RVW, xyNormal, p.R, p.M, CushionHeight, p.Ec, p.Fc)
	b.Motion = ball.MotionSliding
}

// ResolveCircular reflects off a jaw peg at center c (pixels).
func ResolveCircular(b *ball.Ball, center vec.Vec2) {
	if UseMathavan2010 {
		resolveMathavanCircular(b, center)
		return
	}
	resolveHanCircular(b, center)
}

func resolveHanCircular(b *ball.Ball, center vec.Vec2) {
	p := b.Params
	cM := vec.Vec3{X: ball.PxToMeters(center.X), Y: ball.PxToMeters(center.Y), Z: p.R}
	d := b.RVW.R.Sub(cM)
	dist := ptmath.Norm3(d)
	sumR := p.R + ball.PxToMeters(5)
	if dist >= sumR || dist == 0 {
		return
	}
	n := d.Scale(1 / dist)
	b.RVW.R = cM.Add(n.Scale(sumR))
	vn := b.RVW.V.Dot(n)
	if vn >= 0 {
		return
	}
	b.RVW.V = b.RVW.V.Sub(n.Scale((1 + p.Ec) * vn))
	b.Motion = ball.MotionSliding
}

func resolveMathavanLinear(b *ball.Ball, n vec.Vec2) {
	p := b.Params
	xyNormal := vec.Vec3{X: n.X, Y: n.Y}
	vn := b.RVW.V.X*xyNormal.X + b.RVW.V.Y*xyNormal.Y
	if vn >= 0 {
		return
	}
	contactNormal := xyNormal
	if b.RVW.V.Dot(contactNormal) <= 0 {
		contactNormal = contactNormal.Scale(-1)
	}
	var ok bool
	if b.RVW, ok = mathavan2010(b.RVW, contactNormal, p.R, p.M, CushionHeight, p.Ec, p.Fc, p.Us); !ok {
		b.RVW = han2005(b.RVW, xyNormal, p.R, p.M, CushionHeight, p.Ec, p.Fc)
	}
	b.Motion = ball.MotionSliding
}

func resolveMathavanCircular(b *ball.Ball, center vec.Vec2) {
	p := b.Params
	cM := vec.Vec3{X: ball.PxToMeters(center.X), Y: ball.PxToMeters(center.Y), Z: p.R}
	d := b.RVW.R.Sub(cM)
	dist := ptmath.Norm3(d)
	sumR := p.R + ball.PxToMeters(5)
	if dist >= sumR || dist == 0 {
		return
	}
	n := d.Scale(1 / dist)
	b.RVW.R = cM.Add(n.Scale(sumR))
	if b.RVW.V.Dot(n) >= 0 {
		return
	}
	contactNormal := n
	if b.RVW.V.Dot(contactNormal) <= 0 {
		contactNormal = contactNormal.Scale(-1)
	}
	var ok bool
	if b.RVW, ok = mathavan2010(b.RVW, contactNormal, p.R, p.M, CushionHeight, p.Ec, p.Fc, p.Us); !ok {
		b.RVW.V = b.RVW.V.Sub(n.Scale((1 + p.Ec) * b.RVW.V.Dot(n)))
	}
	b.Motion = ball.MotionSliding
}

// mathavan2010 implements pooltool's 2D Mathavan 2010 cushion impact model.
func mathavan2010(rvw ball.RVW, xyNormal vec.Vec3, radius, m, h, ec, fc, us float64) (ball.RVW, bool) {
	psi := ptmath.Angle(xyNormal, vec.Vec3{X: 1})
	angleToRotate := math.Pi/2 - psi
	rot := func(v vec.Vec3) vec.Vec3 { return ptmath.CoordinateRotation(v, angleToRotate) }
	unrot := func(v vec.Vec3) vec.Vec3 { return ptmath.CoordinateRotation(v, -angleToRotate) }

	r := ball.RVW{V: rot(rvw.V), W: rot(rvw.W)}
	vx, vy, wx, wy, wz, ok := solveMathavan(m, radius, h, ec, us, fc, r.V.X, r.V.Y, r.W.X, r.W.Y, r.W.Z)
	if !ok {
		return rvw, false
	}
	r.V.X = vx
	r.V.Y = vy
	r.V.Z = 0
	r.W.X = wx
	r.W.Y = wy
	r.W.Z = wz

	return ball.RVW{
		R: rvw.R,
		V: unrot(r.V),
		W: unrot(r.W),
	}, true
}

func solveMathavan(m, radius, h, ec, us, fc, vx, vy, wx, wy, wz float64) (float64, float64, float64, float64, float64, bool) {
	sinTheta, cosTheta := sinCosTheta(h, radius)
	vx, vy, wx, wy, wz, work, ok := compressionPhase(m, radius, us, fc, sinTheta, cosTheta, vx, vy, wx, wy, wz)
	if !ok {
		return 0, 0, 0, 0, 0, false
	}
	vx, vy, wx, wy, wz, ok = restitutionPhase(m, radius, us, fc, sinTheta, cosTheta, vx, vy, wx, wy, wz, ec*ec*work)
	return vx, vy, wx, wy, wz, ok
}

func sinCosTheta(h, radius float64) (float64, float64) {
	sinTheta := (h - radius) / radius
	if sinTheta > 1 {
		sinTheta = 1
	}
	if sinTheta < -1 {
		sinTheta = -1
	}
	return sinTheta, math.Sqrt(math.Max(0, 1-sinTheta*sinTheta))
}

func slipAngles(radius, sinTheta, cosTheta, vx, vy, wx, wy, wz float64) (float64, float64) {
	vxI := vx + wy*radius*sinTheta - wz*radius*cosTheta
	vyI := -vy*sinTheta + wx*radius
	vxC := vx - wy*radius
	vyC := vy + wx*radius
	return positiveAtan2(vyI, vxI), positiveAtan2(vyC, vxC)
}

func positiveAtan2(y, x float64) float64 {
	a := math.Atan2(y, x)
	if a < 0 {
		a += 2 * math.Pi
	}
	return a
}

func updateVelocity(m, us, fc, sinTheta, cosTheta, vx, vy, slipAngle, slipAnglePrime, deltaP float64) (float64, float64) {
	vx -= (fc*math.Cos(slipAngle) + us*math.Cos(slipAnglePrime)*(sinTheta+fc*math.Sin(slipAngle)*cosTheta)) * deltaP / m
	vy -= (cosTheta - fc*sinTheta*math.Sin(slipAngle) + us*math.Sin(slipAnglePrime)*(sinTheta+fc*math.Sin(slipAngle)*cosTheta)) * deltaP / m
	return vx, vy
}

func updateAngularVelocity(m, radius, us, fc, sinTheta, cosTheta, wx, wy, wz, slipAngle, slipAnglePrime, deltaP float64) (float64, float64, float64) {
	factor := 5 / (2 * m * radius)
	wx += -factor * (fc*math.Sin(slipAngle) + us*math.Sin(slipAnglePrime)*(sinTheta+fc*math.Sin(slipAngle)*cosTheta)) * deltaP
	wy += -factor * (fc*math.Cos(slipAngle)*sinTheta - us*math.Cos(slipAnglePrime)*(sinTheta+fc*math.Sin(slipAngle)*cosTheta)) * deltaP
	wz += factor * fc * math.Cos(slipAngle) * cosTheta * deltaP
	return wx, wy, wz
}

func workDone(vy, cosTheta, deltaP float64) float64 {
	return deltaP * math.Abs(vy) * cosTheta
}

func compressionPhase(m, radius, us, fc, sinTheta, cosTheta, vx, vy, wx, wy, wz float64) (float64, float64, float64, float64, float64, float64, bool) {
	const maxSteps = 1000
	deltaP := math.Max(m*vy/maxSteps, 0.001)
	work := 0.0
	for step := 0; vy > 0; step++ {
		if step > 10*maxSteps {
			return 0, 0, 0, 0, 0, 0, false
		}
		slipAngle, slipAnglePrime := slipAngles(radius, sinTheta, cosTheta, vx, vy, wx, wy, wz)
		nextVX, nextVY := updateVelocity(m, us, fc, sinTheta, cosTheta, vx, vy, slipAngle, slipAnglePrime, deltaP)
		if nextVY <= 0 {
			return refineCompression(m, radius, us, fc, sinTheta, cosTheta, vx, vy, wx, wy, wz, work, deltaP)
		}
		vx, vy = nextVX, nextVY
		wx, wy, wz = updateAngularVelocity(m, radius, us, fc, sinTheta, cosTheta, wx, wy, wz, slipAngle, slipAnglePrime, deltaP)
		work += workDone(vy, cosTheta, deltaP)
	}
	return vx, vy, wx, wy, wz, work, true
}

func refineCompression(m, radius, us, fc, sinTheta, cosTheta, vx, vy, wx, wy, wz, work, deltaP float64) (float64, float64, float64, float64, float64, float64, bool) {
	for range 8 {
		deltaP /= 2
		slipAngle, slipAnglePrime := slipAngles(radius, sinTheta, cosTheta, vx, vy, wx, wy, wz)
		testVX, testVY := updateVelocity(m, us, fc, sinTheta, cosTheta, vx, vy, slipAngle, slipAnglePrime, deltaP)
		if testVY <= 0 {
			continue
		}
		vx, vy = testVX, testVY
		wx, wy, wz = updateAngularVelocity(m, radius, us, fc, sinTheta, cosTheta, wx, wy, wz, slipAngle, slipAnglePrime, deltaP)
		work += workDone(vy, cosTheta, deltaP)
	}
	return vx, vy, wx, wy, wz, work, true
}

func restitutionPhase(m, radius, us, fc, sinTheta, cosTheta, vx, vy, wx, wy, wz, targetWork float64) (float64, float64, float64, float64, float64, bool) {
	const maxSteps = 1000
	if targetWork <= 0 {
		return vx, vy, wx, wy, wz, true
	}
	deltaP := math.Max(targetWork/maxSteps, 0.001)
	work := 0.0
	for step := 0; work < targetWork; step++ {
		if step > 10*maxSteps {
			return 0, 0, 0, 0, 0, false
		}
		slipAngle, slipAnglePrime := slipAngles(radius, sinTheta, cosTheta, vx, vy, wx, wy, wz)
		nextWork := workDone(vy, cosTheta, deltaP)
		if work+nextWork > targetWork {
			remaining := targetWork - work
			den := math.Abs(vy) * cosTheta
			if den <= ptmath.Eps {
				return vx, vy, wx, wy, wz, true
			}
			deltaP = remaining / den
			slipAngle, slipAnglePrime = slipAngles(radius, sinTheta, cosTheta, vx, vy, wx, wy, wz)
			vx, vy = updateVelocity(m, us, fc, sinTheta, cosTheta, vx, vy, slipAngle, slipAnglePrime, deltaP)
			wx, wy, wz = updateAngularVelocity(m, radius, us, fc, sinTheta, cosTheta, wx, wy, wz, slipAngle, slipAnglePrime, deltaP)
			return vx, vy, wx, wy, wz, true
		}
		vx, vy = updateVelocity(m, us, fc, sinTheta, cosTheta, vx, vy, slipAngle, slipAnglePrime, deltaP)
		wx, wy, wz = updateAngularVelocity(m, radius, us, fc, sinTheta, cosTheta, wx, wy, wz, slipAngle, slipAnglePrime, deltaP)
		work += workDone(vy, cosTheta, deltaP)
	}
	return vx, vy, wx, wy, wz, true
}
