package physics

import (
	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/physics/ptmath"
	"github.com/user/pooltest/pool-go/internal/physics/resolve/ballball"
	"github.com/user/pooltest/pool-go/internal/physics/resolve/ballcushion"
	"github.com/user/pooltest/pool-go/internal/table"
	"github.com/user/pooltest/pool-go/internal/vec"
)

// SideMouth exported for tests.
const SideMouth = table.SideMouth

// JawPegs holds cushion tip positions (pixels) for tests.
var JawPegs = []vec.Vec2{
	{X: table.PlayLeft + table.CornerMouth, Y: table.PlayTop},
	{X: table.PlayLeft, Y: table.PlayTop + table.CornerMouth},
	{X: table.PlayRight - table.CornerMouth, Y: table.PlayTop},
	{X: table.PlayRight, Y: table.PlayTop + table.CornerMouth},
	{X: table.PlayLeft + table.CornerMouth, Y: table.PlayBottom},
	{X: table.PlayLeft, Y: table.PlayBottom - table.CornerMouth},
	{X: table.PlayRight - table.CornerMouth, Y: table.PlayBottom},
	{X: table.PlayRight, Y: table.PlayBottom - table.CornerMouth},
	{X: table.MidX - table.SideMouth, Y: table.PlayTop},
	{X: table.MidX + table.SideMouth, Y: table.PlayTop},
	{X: table.MidX - table.SideMouth, Y: table.PlayBottom},
	{X: table.MidX + table.SideMouth, Y: table.PlayBottom},
}

// ResolveBallCollisions resolves overlapping active balls (test + legacy API).
func ResolveBallCollisions(balls []*ball.Ball) {
	for i := range balls {
		a := balls[i]
		if !a.Active {
			continue
		}
		for j := i + 1; j < len(balls); j++ {
			b := balls[j]
			if !b.Active {
				continue
			}
			sumR := a.Params.R + b.Params.R
			if a.RVW.R.Sub(b.RVW.R).LenSq() >= sumR*sumR {
				continue
			}
			vn := b.RVW.V.Sub(a.RVW.V)
			// approximate normal
			n := b.RVW.R.Sub(a.RVW.R)
			if n.LenSq() == 0 {
				continue
			}
			_ = vn
			ballball.Resolve(a, b)
		}
	}
}

// ResolveRailCollisions keeps balls inside play area (test helper).
func ResolveRailCollisions(balls []*ball.Ball) {
	cr := ball.CollisionRadius()
	minX := float64(table.PlayLeft) + cr
	maxX := float64(table.PlayRight) - cr
	minY := float64(table.PlayTop) + cr
	maxY := float64(table.PlayBottom) - cr
	for _, b := range balls {
		if !b.Active {
			continue
		}
		px := ball.MetersToPx(b.RVW.R.X)
		py := ball.MetersToPx(b.RVW.R.Y)
		if !table.InLeftRightMouth(py) {
			switch {
			case px < minX:
				b.SetPosPx(vec.Vec2{X: minX, Y: py})
				ballcushion.ResolveLinear(b, vec.Vec2{X: 1, Y: 0})
			case px > maxX:
				b.SetPosPx(vec.Vec2{X: maxX, Y: py})
				ballcushion.ResolveLinear(b, vec.Vec2{X: -1, Y: 0})
			}
		}
		if !table.InTopBottomMouth(px) {
			switch {
			case py < minY:
				b.SetPosPx(vec.Vec2{X: px, Y: minY})
				ballcushion.ResolveLinear(b, vec.Vec2{X: 0, Y: 1})
			case py > maxY:
				b.SetPosPx(vec.Vec2{X: px, Y: maxY})
				ballcushion.ResolveLinear(b, vec.Vec2{X: 0, Y: -1})
			}
		}
		b.SyncXY(1.0 / 60.0)
	}
}

// ResolveJawCollisions bounces balls off jaw pegs.
func ResolveJawCollisions(balls []*ball.Ball) {
	for _, b := range balls {
		if !b.Active {
			continue
		}
		for _, c := range JawPegs {
			ReflectOffPeg(b, c)
		}
	}
}

// ReflectOffPeg reflects ball off peg; returns impact speed or 0.
func ReflectOffPeg(b *ball.Ball, c vec.Vec2) float64 {
	before := ptmath.Norm3(b.RVW.V)
	ballcushion.ResolveCircular(b, c)
	after := ptmath.Norm3(b.RVW.V)
	b.SyncXY(1.0 / 60.0)
	if after > 0 && before > 0 {
		return ball.MetersToPx(before) / 60
	}
	return 0
}

// InTopBottomMouth reports open pocket mouth on top/bottom rail.
func InTopBottomMouth(x float64) bool { return table.InTopBottomMouth(x) }

// InLeftRightMouth reports open pocket mouth on side rails.
func InLeftRightMouth(y float64) bool { return table.InLeftRightMouth(y) }
