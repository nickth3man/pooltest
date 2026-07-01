package ball

import (
	"github.com/user/pooltest/pool-go/internal/table"
)

const (
	// SpriteRadius is the visual ball radius in pixels.
	SpriteRadius = 11.0

	// PlayLengthM is pooltool SEVEN_FOOT_SHOWOOD play length (m).
	PlayLengthM = 1.9812
	// PlayWidthM is play width (m).
	PlayWidthM = 0.9906
)

// PxPerM converts meters to pixels from the play rectangle.
func PxPerM() float64 {
	playW := float64(table.PlayRight - table.PlayLeft)
	return playW / PlayLengthM
}

// CollisionRadiusPx returns physics collision radius in pixels for params.
func CollisionRadiusPx(p Params) float64 {
	return p.R * PxPerM()
}

// CollisionRadius returns the default collision radius in pixels.
func CollisionRadius() float64 {
	return CollisionRadiusPx(DefaultParams)
}

// Radius is kept for compatibility: sprite radius for rendering/placement UI.
const Radius = SpriteRadius
