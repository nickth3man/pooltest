package main

import "math"

// newRack builds the opening table: the cue ball on the head spot and the
// fifteen object balls racked in a triangle on the foot spot with the 8-ball in
// the center. balls[0] is always the cue ball.
//
// The corner solid/stripe alternation of tournament racks is not enforced; only
// the apex, the centered 8-ball, and a legal (overlap-free) triangle matter for
// the simulation.
func newRack() []*Ball {
	const spacing = 2*ballRadius + 1 // small gap so racked balls don't start overlapping
	rowDX := spacing * math.Sqrt(3) / 2

	centerY := float64(playTop+playBottom) / 2
	apexX := float64(playLeft) + footSpotFrac*float64(playRight-playLeft)
	headX := float64(playLeft) + headSpotFrac*float64(playRight-playLeft)

	rows := [5][]int{
		{1},
		{2, 3},
		{4, 8, 5},
		{6, 7, 9, 10},
		{11, 12, 13, 14, 15},
	}

	balls := make([]*Ball, 0, 16)
	balls = append(balls, newBall(0, Vec2{headX, centerY}))
	for r, row := range rows {
		x := apexX + float64(r)*rowDX
		for k, num := range row {
			y := centerY + (float64(k)-float64(r)/2)*spacing
			balls = append(balls, newBall(num, Vec2{x, y}))
		}
	}
	return balls
}

func newBall(number int, pos Vec2) *Ball {
	b := &Ball{Number: number, Pos: pos, Active: true}
	switch {
	case number == 0:
		b.Color = colorWhite
	case number >= 1 && number <= 8:
		b.Color = ballColors[number]
	default: // 9-15 reuse the 1-7 colors as stripes
		b.Color = ballColors[number-8]
		b.Stripe = true
	}
	return b
}
