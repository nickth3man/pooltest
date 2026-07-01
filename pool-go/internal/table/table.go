// Package table holds the table geometry constants and the derived pocket
// centers / jaw peg positions.
package table

import "github.com/user/pooltest/pool-go/internal/vec"

const (
	// ScreenW is the fixed logical screen width returned by Layout.
	ScreenW = 800
	// ScreenH is the fixed logical screen height returned by Layout.
	ScreenH = 600

	// PlayLeft is the left cushion x of the felt rectangle.
	PlayLeft = 70
	// PlayTop is the top cushion y of the felt rectangle.
	PlayTop = 130
	// PlayRight is the right cushion x of the felt rectangle.
	PlayRight = 730
	// PlayBottom is the bottom cushion y of the felt rectangle.
	PlayBottom = 460

	// RailWidth is the visual rail frame thickness around the felt.
	RailWidth = 26

	// PocketRadius is the capture radius of a pocket.
	PocketRadius = 18

	// MidX is the horizontal center of the play rectangle, where the side
	// pockets sit.
	MidX = float64(PlayLeft+PlayRight) / 2

	// HeadSpotFrac places the cue ball (head spot) along the long axis; the
	// drawn table spot uses the same fraction.
	HeadSpotFrac = 0.25
	// FootSpotFrac places the rack apex (foot spot) along the long axis; the
	// drawn table spot uses the same fraction.
	FootSpotFrac = 0.75
)

// PocketCenters holds the six pocket centers: four corners plus two side
// pockets. The table geometry is fixed, so it is computed once.
var PocketCenters = [6]vec.Vec2{
	{X: PlayLeft, Y: PlayTop},
	{X: MidX, Y: PlayTop},
	{X: PlayRight, Y: PlayTop},
	{X: PlayLeft, Y: PlayBottom},
	{X: MidX, Y: PlayBottom},
	{X: PlayRight, Y: PlayBottom},
}
