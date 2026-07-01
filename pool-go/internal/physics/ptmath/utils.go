// Package ptmath provides pool physics math utilities (ported from pooltool/ptmath).
package ptmath

import (
	"math"

	"github.com/user/pooltest/pool-go/internal/vec"
)

const Eps = 1e-9

// Norm3 returns |v|.
func Norm3(v vec.Vec3) float64 { return v.Len() }

// Unit3 returns the unit vector; zero maps to zero.
func Unit3(v vec.Vec3) vec.Vec3 { return v.Normalize() }

// Angle returns counter-clockwise angle of v's xy projection from v1 (default x-axis).
func Angle(v, v1 vec.Vec3) float64 {
	if v1.LenSq() == 0 {
		v1 = vec.Vec3{X: 1}
	}
	ang := math.Atan2(v.Y, v.X) - math.Atan2(v1.Y, v1.X)
	if ang < 0 {
		ang += 2 * math.Pi
	}
	return ang
}

// CoordinateRotation rotates v about z by phi radians.
func CoordinateRotation(v vec.Vec3, phi float64) vec.Vec3 {
	c, s := math.Cos(phi), math.Sin(phi)
	return vec.Vec3{
		X: c*v.X - s*v.Y,
		Y: s*v.X + c*v.Y,
		Z: v.Z,
	}
}

// Cross returns u × v.
func Cross(u, v vec.Vec3) vec.Vec3 { return u.Cross(v) }

// GetSlideTime returns seconds until sliding ends (pooltool).
func GetSlideTime(rvw ballRVW, radius, us, g float64) float64 {
	if us == 0 {
		return math.Inf(1)
	}
	rel := RelVelocity(rvw, radius)
	return 2 * Norm3(rel) / (7 * us * g)
}

// GetRollTime returns seconds until rolling ends.
func GetRollTime(rvw ballRVW, ur, g float64) float64 {
	if ur == 0 {
		return math.Inf(1)
	}
	return Norm3(rvw.V) / (ur * g)
}

// GetSpinTime returns seconds until spinning ends.
func GetSpinTime(rvw ballRVW, radius, usp, g float64) float64 {
	if usp == 0 {
		return math.Inf(1)
	}
	return math.Abs(rvw.W.Z) * 2 / 5 * radius / usp / g
}

// ballRVW is a minimal view to avoid import cycles.
type ballRVW struct {
	R, V, W vec.Vec3
}

// RVWFrom copies ball RVW fields.
func RVWFrom(r, v, w vec.Vec3) ballRVW {
	return ballRVW{R: r, V: v, W: w}
}

// SurfaceVelocity returns velocity at surface point in direction d.
func SurfaceVelocity(rvw ballRVW, d vec.Vec3, radius float64) vec.Vec3 {
	return rvw.V.Add(Cross(rvw.W, d.Scale(radius)))
}

// RelVelocity returns contact-point velocity relative to cloth (sliding indicator).
func RelVelocity(rvw ballRVW, radius float64) vec.Vec3 {
	down := vec.Vec3{Z: -1}
	return SurfaceVelocity(rvw, down, radius)
}

// TangentSurfaceVelocity returns tangent velocity at contact normal d.
func TangentSurfaceVelocity(rvw ballRVW, d vec.Vec3, radius float64) vec.Vec3 {
	vn := d.Scale(rvw.V.Dot(d))
	vt := rvw.V.Sub(vn)
	return vt.Add(Cross(rvw.W, d.Scale(radius)))
}

// SolveQuadratic solves a t² + b t + c = 0; returns smaller and larger real roots.
func SolveQuadratic(a, b, c float64) (t1, t2 float64, ok bool) {
	if math.Abs(a) < Eps {
		if math.Abs(b) < Eps {
			return 0, 0, false
		}
		t := -c / b
		return t, t, true
	}
	disc := b*b - 4*a*c
	if disc < 0 {
		return 0, 0, false
	}
	s := math.Sqrt(disc)
	return (-b - s) / (2 * a), (-b + s) / (2 * a), true
}

// EarliestBallBallTime estimates collision time for two balls with constant velocity.
func EarliestBallBallTime(p1, v1, p2, v2 vec.Vec3, sumR float64) (float64, bool) {
	d := p1.Sub(p2)
	dv := v1.Sub(v2)
	a := dv.LenSq()
	b := 2 * d.Dot(dv)
	c := d.LenSq() - sumR*sumR
	t1, t2, ok := SolveQuadratic(a, b, c)
	if !ok {
		return 0, false
	}
	best := math.Inf(1)
	for _, t := range []float64{t1, t2} {
		if t > Eps && t < best {
			best = t
		}
	}
	return best, best < math.Inf(1)
}
