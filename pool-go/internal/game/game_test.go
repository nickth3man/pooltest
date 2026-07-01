package game

import (
	"math"
	"testing"

	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/vec"
)

func TestFirstContactFindsGhostBall(t *testing.T) {
	g := &Game{}
	g.cue = &ball.Ball{Number: 0, Pos: vec.Vec2{X: 100, Y: 300}, Active: true}
	obj := &ball.Ball{Number: 1, Pos: vec.Vec2{X: 200, Y: 300}, Active: true}
	g.balls = []*ball.Ball{g.cue, obj}

	ghost, target, ok := g.firstContact(vec.Vec2{X: 1, Y: 0})
	if !ok || target != obj {
		t.Fatalf("firstContact missed the object ball (ok=%v)", ok)
	}
	if want := 200.0 - 2*ball.Radius; math.Abs(ghost.X-want) > 0.5 {
		t.Errorf("ghost ball at x=%.2f, want %.2f", ghost.X, want)
	}
}
