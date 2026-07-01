package ballcushion

import (
	"math"

	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/physics/ptmath"
	"github.com/user/pooltest/pool-go/internal/vec"
)

// CushionHeight is nominal cushion height (m) for cushion models.
const CushionHeight = 0.03654

func han2005(rvw ball.RVW, xyNormal vec.Vec3, radius, m, h, ec, fc float64) ball.RVW {
	psi := ptmath.Angle(xyNormal, vec.Vec3{X: 1})
	rot := func(v vec.Vec3) vec.Vec3 { return ptmath.CoordinateRotation(v, -psi) }
	unrot := func(v vec.Vec3) vec.Vec3 { return ptmath.CoordinateRotation(v, psi) }

	rvwR := ball.RVW{
		V: rot(rvw.V),
		W: rot(rvw.W),
	}

	e := ec
	mu := fc
	thetaA := math.Asin(h/radius - 1)

	sx := rvwR.V.X*math.Sin(thetaA) - rvwR.V.Z*math.Cos(thetaA) + radius*rvwR.W.Y
	sy := -rvwR.V.Y - radius*rvwR.W.Z*math.Cos(thetaA) + radius*rvwR.W.X*math.Sin(thetaA)
	c := -rvwR.V.X * math.Cos(thetaA)

	II := 2.0 / 5.0 * m * radius * radius
	PzE := -(1 + e) * c / (1.0 / m)
	absS0 := math.Hypot(sx, sy)
	PzS := absS0 / (7.0 / 2.0 / m)

	var PxE, PyE float64
	if PzS <= mu*PzE {
		PxE = sx / (7.0 / 2.0 / m)
		PyE = sy / (7.0 / 2.0 / m)
	} else {
		PxE = mu * PzE * sx / absS0
		PyE = mu * PzE * sy / absS0
	}

	PX := -PxE*math.Sin(thetaA) - PzE*math.Cos(thetaA)
	PY := PyE
	PZ := PxE*math.Cos(thetaA) - PzE*math.Sin(thetaA)
	_ = PZ

	rvwR.V.X += PX / m
	rvwR.V.Y += PY / m
	rvwR.W.X += -radius / II * PY * math.Sin(thetaA)
	rvwR.W.Y += radius / II * (PX*math.Sin(thetaA) - PZ*math.Cos(thetaA))
	rvwR.W.Z += radius / II * PY * math.Cos(thetaA)

	return ball.RVW{
		R: rvw.R,
		V: unrot(rvwR.V),
		W: unrot(rvwR.W),
	}
}
