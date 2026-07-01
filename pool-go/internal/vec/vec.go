// Package vec provides a minimal 2D vector used for positions (pixels) and
// velocities (pixels per frame).
package vec

import "math"

// Vec2 is a 2D vector used for positions (pixels) and velocities
// (pixels per frame).
type Vec2 struct {
	X, Y float64
}

// Add returns the component-wise sum v + o.
func (v Vec2) Add(o Vec2) Vec2 { return Vec2{v.X + o.X, v.Y + o.Y} }

// Sub returns the component-wise difference v - o.
func (v Vec2) Sub(o Vec2) Vec2 { return Vec2{v.X - o.X, v.Y - o.Y} }

// Scale returns v multiplied by the scalar s.
func (v Vec2) Scale(s float64) Vec2 { return Vec2{v.X * s, v.Y * s} }

// Dot returns the dot product v · o.
func (v Vec2) Dot(o Vec2) float64 { return v.X*o.X + v.Y*o.Y }

// Len returns the Euclidean length |v|.
func (v Vec2) Len() float64 { return math.Sqrt(v.X*v.X + v.Y*v.Y) }

// LenSq returns the squared Euclidean length |v|² (avoids a sqrt for compares).
func (v Vec2) LenSq() float64 { return v.X*v.X + v.Y*v.Y }

// Perp returns the vector rotated 90° counter-clockwise (screen space).
func (v Vec2) Perp() Vec2 { return Vec2{-v.Y, v.X} }

// Normalize returns the unit vector with the same direction. The zero vector
// normalizes to the zero vector.
func (v Vec2) Normalize() Vec2 {
	l := v.Len()
	if l == 0 {
		return Vec2{}
	}
	return Vec2{v.X / l, v.Y / l}
}
