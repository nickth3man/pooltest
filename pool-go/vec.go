package main

import "math"

// Vec2 is a 2D vector used for positions (pixels) and velocities
// (pixels per frame).
type Vec2 struct {
	X, Y float64
}

func (v Vec2) Add(o Vec2) Vec2      { return Vec2{v.X + o.X, v.Y + o.Y} }
func (v Vec2) Sub(o Vec2) Vec2      { return Vec2{v.X - o.X, v.Y - o.Y} }
func (v Vec2) Scale(s float64) Vec2 { return Vec2{v.X * s, v.Y * s} }
func (v Vec2) Dot(o Vec2) float64   { return v.X*o.X + v.Y*o.Y }
func (v Vec2) Len() float64         { return math.Hypot(v.X, v.Y) }
func (v Vec2) LenSq() float64       { return v.X*v.X + v.Y*v.Y }

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
