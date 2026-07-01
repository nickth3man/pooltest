package game

import (
	"math"

	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/vec"
)

// firstContact walks the aim ray from the cue ball and returns the point at
// which it would first touch another ball, that ball, and whether a hit exists.
func (g *Game) firstContact(dir vec.Vec2) (vec.Vec2, *ball.Ball, bool) {
	best := math.Inf(1)
	var target *ball.Ball
	rad := 2 * ball.Radius
	for _, b := range g.balls {
		if b == g.cue || !b.Active {
			continue
		}
		f := b.Pos.Sub(g.cue.Pos)
		proj := f.Dot(dir)
		if proj < 0 {
			continue // behind the shot
		}
		perpSq := f.LenSq() - proj*proj
		if perpSq > rad*rad {
			continue // ray misses this ball
		}
		t := proj - math.Sqrt(rad*rad-perpSq)
		if t < best {
			best, target = t, b
		}
	}
	if target == nil {
		return vec.Vec2{}, nil, false
	}
	return g.cue.Pos.Add(dir.Scale(best)), target, true
}
