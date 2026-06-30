package main

// Table geometry. The play rectangle (playLeft..playRight, playTop..playBottom)
// is the cushion-bounded felt; ball centers are later constrained to it inset by
// the ball radius. The long axis is horizontal, so the side pockets sit on the
// middle of the top and bottom rails.
const (
	screenW = 800
	screenH = 600

	playLeft   = 70
	playTop    = 130
	playRight  = 730
	playBottom = 460

	railWidth = 26 // visual rail frame thickness around the felt

	pocketRadius = 19 // capture radius of a pocket
)

// pockets returns the six pocket centers: four corners plus two side pockets.
func pockets() [6]Vec2 {
	midX := float64(playLeft+playRight) / 2
	return [6]Vec2{
		{playLeft, playTop},
		{midX, playTop},
		{playRight, playTop},
		{playLeft, playBottom},
		{midX, playBottom},
		{playRight, playBottom},
	}
}
