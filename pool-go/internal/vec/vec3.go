package vec

import "math"

// Vec3 is a 3D vector (table coords: x/y play plane, z height).
type Vec3 struct {
	X, Y, Z float64
}

// Len returns the Euclidean length.
func (v Vec3) Len() float64 { return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z) }

// LenSq returns squared length.
func (v Vec3) LenSq() float64 { return v.X*v.X + v.Y*v.Y + v.Z*v.Z }

// Add returns v + o.
func (v Vec3) Add(o Vec3) Vec3 { return Vec3{v.X + o.X, v.Y + o.Y, v.Z + o.Z} }

// Sub returns v - o.
func (v Vec3) Sub(o Vec3) Vec3 { return Vec3{v.X - o.X, v.Y - o.Y, v.Z - o.Z} }

// Scale returns v * s.
func (v Vec3) Scale(s float64) Vec3 { return Vec3{v.X * s, v.Y * s, v.Z * s} }

// Dot returns v · o.
func (v Vec3) Dot(o Vec3) float64 { return v.X*o.X + v.Y*o.Y + v.Z*o.Z }

// Cross returns v × o.
func (v Vec3) Cross(o Vec3) Vec3 {
	return Vec3{
		v.Y*o.Z - v.Z*o.Y,
		v.Z*o.X - v.X*o.Z,
		v.X*o.Y - v.Y*o.X,
	}
}

// XY returns the table-plane projection.
func (v Vec3) XY() Vec2 { return Vec2{X: v.X, Y: v.Y} }

// FromXY builds a Vec3 with z=0.
func FromXY(v Vec2) Vec3 { return Vec3{X: v.X, Y: v.Y} }

// WithZ returns v with the z component replaced.
func (v Vec3) WithZ(z float64) Vec3 { return Vec3{v.X, v.Y, z} }

// Normalize returns the unit vector; zero vector maps to zero.
func (v Vec3) Normalize() Vec3 {
	l := v.Len()
	if l == 0 {
		return Vec3{}
	}
	return v.Scale(1 / l)
}
