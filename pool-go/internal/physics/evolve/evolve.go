// Package evolve integrates ball motion through sliding, rolling, and spinning phases.
package evolve

import (
	"math"

	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/physics/ptmath"
	"github.com/user/pooltest/pool-go/internal/vec"
)

// EvolveBallMotion advances rvw by t seconds and returns updated motion state.
func EvolveBallMotion(state ball.MotionState, rvw ball.RVW, p ball.Params, t float64) (ball.RVW, ball.MotionState) {
	if t <= 0 || state == ball.MotionStationary || state == ball.MotionPocketed {
		return rvw, state
	}
	usp := p.USp()
	for t > 0 {
		switch state {
		case ball.MotionSliding:
			dtau := ptmath.GetSlideTime(ptmath.RVWFrom(rvw.R, rvw.V, rvw.W), p.R, p.Us, p.G)
			if t >= dtau {
				rvw = evolveSlideState(rvw, p, usp, dtau)
				state = ball.MotionRolling
				t -= dtau
			} else {
				return evolveSlideState(rvw, p, usp, t), ball.MotionSliding
			}
		case ball.MotionRolling:
			dtau := ptmath.GetRollTime(ptmath.RVWFrom(rvw.R, rvw.V, rvw.W), p.Ur, p.G)
			if t >= dtau {
				rvw = evolveRollState(rvw, p, usp, dtau)
				state = ball.MotionSpinning
				t -= dtau
			} else {
				return evolveRollState(rvw, p, usp, t), ball.MotionRolling
			}
		case ball.MotionSpinning:
			dtau := ptmath.GetSpinTime(ptmath.RVWFrom(rvw.R, rvw.V, rvw.W), p.R, usp, p.G)
			if t >= dtau {
				rvw = evolvePerpendicularSpinState(rvw, p, usp, dtau)
				return rvw, ball.MotionStationary
			}
			return evolvePerpendicularSpinState(rvw, p, usp, t), ball.MotionSpinning
		default:
			return rvw, state
		}
	}
	return rvw, state
}

func evolveSlideState(rvw ball.RVW, p ball.Params, usp, t float64) ball.RVW {
	if t == 0 {
		return rvw
	}
	phi := ptmath.Angle(rvw.V, vec.Vec3{X: 1})
	rot := func(v vec.Vec3) vec.Vec3 { return ptmath.CoordinateRotation(v, -phi) }
	unrot := func(v vec.Vec3) vec.Vec3 { return ptmath.CoordinateRotation(v, phi) }

	rvwB0R := rot(rvw.R)
	rvwB0V := rot(rvw.V)
	rvwB0W := rot(rvw.W)

	rel := ptmath.RelVelocity(ptmath.RVWFrom(rvw.R, rvw.V, rvw.W), p.R)
	u0 := rot(ptmath.Unit3(rel))

	g := p.G
	us := p.Us
	drag := us * g * t

	rvwB := ball.RVW{}
	rvwB.R = vec.Vec3{
		X: rvwB0R.X + rvwB0V.X*t - 0.5*drag*t*u0.X,
		Y: rvwB0R.Y + rvwB0V.Y*t - 0.5*drag*t*u0.Y,
		Z: p.R,
	}
	rvwB.V = rvwB0V.Sub(u0.Scale(drag))
	spinCross := ptmath.Cross(u0, vec.Vec3{Z: 1}).Scale(5.0 / 2.0 / p.R * drag)
	rvwB.W = rvwB0W.Sub(spinCross)
	rvwB.W.Z = rvwB0W.Z
	rvwB = evolvePerpendicularSpinState(rvwB, p, usp, t)

	out := ball.RVW{
		R: unrot(rvwB.R.Sub(rvwB0R)).Add(rvw.R),
		V: unrot(rvwB.V),
		W: unrot(rvwB.W),
	}
	out.R.Z = p.R
	return out
}

func evolveRollState(rvw ball.RVW, p ball.Params, usp, t float64) ball.RVW {
	if t == 0 {
		return rvw
	}
	v0 := rvw.V
	v0hat := ptmath.Unit3(v0)
	drag := p.Ur * p.G * t

	r := rvw.R.Add(v0.Scale(t)).Sub(v0hat.Scale(0.5 * drag * t))
	v := v0.Sub(v0hat.Scale(drag))
	// Natural roll: ω_xy from v/R.
	w := ptmath.CoordinateRotation(v.Scale(1/p.R), math.Pi/2)
	temp := evolvePerpendicularSpinState(ball.RVW{W: rvw.W}, p, usp, t)
	w.Z = temp.W.Z

	return ball.RVW{R: r, V: v, W: w}
}

func evolvePerpendicularSpinState(rvw ball.RVW, p ball.Params, usp, t float64) ball.RVW {
	if t == 0 {
		return rvw
	}
	wz := rvw.W.Z
	if math.Abs(wz) < ptmath.Eps {
		return rvw
	}
	alpha := 5 * usp * p.G / (2 * p.R)
	maxT := math.Abs(wz) / alpha
	if t > maxT {
		t = maxT
	}
	sign := 1.0
	if wz < 0 {
		sign = -1
	}
	rvw.W.Z = wz - sign*alpha*t
	return rvw
}

// FinalMotionState after a strike returns sliding unless airborne.
func FinalMotionState(rvw ball.RVW) ball.MotionState {
	if rvw.V.Z != 0 {
		return ball.MotionSliding // no airborne in v1
	}
	return ball.MotionSliding
}

// MotionAfterEvolve picks state when velocity/spin are near zero.
func MotionAfterEvolve(rvw ball.RVW, prev ball.MotionState) ball.MotionState {
	if ptmath.Norm3(rvw.V) < ptmath.Eps && math.Abs(rvw.W.Z) < ptmath.Eps {
		return ball.MotionStationary
	}
	return prev
}
