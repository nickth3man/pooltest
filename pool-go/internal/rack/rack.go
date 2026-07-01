// Package rack builds the opening table: cue ball on the head spot plus the
// fifteen object balls racked in a triangle on the foot spot.
package rack

import (
	"math"

	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/table"
	"github.com/user/pooltest/pool-go/internal/vec"
)

// NewRack builds the opening table: the cue ball on the head spot and the
// fifteen object balls racked in a triangle on the foot spot with the 8-ball in
// the center. balls[0] is always the cue ball.
//
// The corner solid/stripe alternation of tournament racks is not enforced; only
// the apex, the centered 8-ball, and a legal (overlap-free) triangle matter for
// the simulation.
func NewRack() []*ball.Ball {
	const spacing = 2*ball.Radius + 1 // small gap so racked balls don't start overlapping
	rowDX := spacing * math.Sqrt(3) / 2

	centerY := float64(table.PlayTop+table.PlayBottom) / 2
	apexX := float64(table.PlayLeft) + table.FootSpotFrac*float64(table.PlayRight-table.PlayLeft)
	headX := float64(table.PlayLeft) + table.HeadSpotFrac*float64(table.PlayRight-table.PlayLeft)

	rows := [5][]int{
		{1},
		{2, 3},
		{4, 8, 5},
		{6, 7, 9, 10},
		{11, 12, 13, 14, 15},
	}

	balls := make([]*ball.Ball, 0, 16)
	balls = append(balls, NewBall(0, vec.Vec2{X: headX, Y: centerY}))
	for r, row := range rows {
		x := apexX + float64(r)*rowDX
		for k, num := range row {
			y := centerY + (float64(k)-float64(r)/2)*spacing
			balls = append(balls, NewBall(num, vec.Vec2{X: x, Y: y}))
		}
	}
	return balls
}

// NewBall returns a ball with the right color and stripe flag for the given
// number, located at pos.
func NewBall(number int, pos vec.Vec2) *ball.Ball {
	b := &ball.Ball{Number: number, Pos: pos, Active: true}
	switch {
	case number == 0:
		b.Color = ball.ColorWhite
	case number >= 1 && number <= 8:
		b.Color = ball.Colors[number]
	default: // 9-15 reuse the 1-7 colors as stripes
		b.Color = ball.Colors[number-8]
		b.Stripe = true
	}
	return b
}
