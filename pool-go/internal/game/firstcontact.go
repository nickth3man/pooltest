package game

import (
	"math"

	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/physics"
	"github.com/user/pooltest/pool-go/internal/vec"
)

// aimDirection applies squirt from the spin dial to the raw pull direction.
func (g *Game) aimDirection(rawDir vec.Vec2) vec.Vec2 {
	dir := rawDir.Normalize()
	// Squirt from InstantaneousPoint model at representative power.
	preview := physics.PreviewCBPath(g.cue.Pos, dir, maxSpeed*0.5, g.spin.X, -g.spin.Y, 1)
	if len(preview) >= 2 {
		d := preview[1].Sub(preview[0])
		if d.LenSq() > 0 {
			return d.Normalize()
		}
	}
	return dir.Rotate(squirtGain * g.spin.X)
}

// firstContact walks the aim ray from the cue ball and returns the point at
// which it would first touch another ball, that ball, and whether a hit exists.
func (g *Game) firstContact(dir vec.Vec2) (vec.Vec2, *ball.Ball, bool) {
	best := math.Inf(1)
	var target *ball.Ball
	rad := 2 * ball.CollisionRadius()
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

// caromDirection estimates the object-ball exit after contact, including throw.
func (g *Game) caromDirection(ghost, targetCenter vec.Vec2, aimDir vec.Vec2, speedFrac float64) vec.Vec2 {
	n := targetCenter.Sub(ghost)
	if n.LenSq() == 0 {
		return vec.Vec2{}
	}
	n = n.Normalize()
	cutAngle := math.Acos(clampDot(math.Abs(aimDir.Dot(n))))
	throwDeg := physics.PredictThrowDeg(g.cue.Pos, targetCenter, aimDir, maxSpeed*speedFrac, g.spin.X, -g.spin.Y)
	if throwDeg == 0 {
		throwDeg = physics.EstimateThrowDeg(cutAngle, speedFrac, g.spin.X, -g.spin.Y)
	}
	cross := aimDir.X*n.Y - aimDir.Y*n.X
	sign := 1.0
	if cross < 0 {
		sign = -1
	}
	return n.Rotate(sign * throwDeg * math.Pi / 180)
}

func clampDot(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}
