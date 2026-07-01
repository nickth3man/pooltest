package physics

import (
	"math"
	"testing"

	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/table"
	"github.com/user/pooltest/pool-go/internal/vec"
)

func TestHeadOnCollisionTransfersMomentum(t *testing.T) {
	a := &ball.Ball{Number: 0, Pos: vec.Vec2{X: 100, Y: 100}, Vel: vec.Vec2{X: 1, Y: 0}, Active: true}
	b := &ball.Ball{Number: 1, Pos: vec.Vec2{X: 100 + 2*ball.Radius - 0.5, Y: 100}, Active: true}

	ResolveBallCollisions([]*ball.Ball{a, b})

	if a.Vel.X > 0.1 {
		t.Errorf("striking ball kept too much speed: %.3f", a.Vel.X)
	}
	if b.Vel.X < 0.8 {
		t.Errorf("struck ball gained too little speed: %.3f", b.Vel.X)
	}
}

func TestRailReflectsVelocity(t *testing.T) {
	b := &ball.Ball{Number: 0, Pos: vec.Vec2{X: table.PlayRight - ball.Radius - 0.1, Y: 300}, Vel: vec.Vec2{X: 5, Y: 0}, Active: true}

	Step([]*ball.Ball{b})

	if b.Vel.X >= 0 {
		t.Errorf("ball did not bounce off the right rail: vx = %.3f", b.Vel.X)
	}
	if b.Pos.X > table.PlayRight-ball.Radius {
		t.Errorf("ball escaped the cushion: x = %.3f > %.3f", b.Pos.X, table.PlayRight-ball.Radius)
	}
}

func TestPocketCapture(t *testing.T) {
	p := table.PocketCenters[0]
	b := &ball.Ball{Number: 3, Pos: p, Active: true}

	got := Step([]*ball.Ball{b})

	if b.Active {
		t.Error("ball over a pocket was not captured")
	}
	if len(got) != 1 || got[0] != 3 {
		t.Errorf("step returned %v, want [3]", got)
	}
}

func TestFrictionBringsBallToRest(t *testing.T) {
	b := &ball.Ball{Number: 0, Pos: vec.Vec2{X: 400, Y: 300}, Vel: vec.Vec2{X: 6, Y: 0}, Active: true}
	balls := []*ball.Ball{b}

	for i := 0; i < 2000 && !AllStopped(balls); i++ {
		Step(balls)
	}

	if b.Moving() {
		t.Errorf("ball never came to rest: speed = %.4f", math.Sqrt(b.Vel.LenSq()))
	}
}

// A stationary ball carrying a backspin reservoir should crawl backward (draw)
// as the cloth grips it, then settle. This is the mechanic behind draw shots.
func TestSpinReservoirMovesAndSettles(t *testing.T) {
	b := &ball.Ball{Number: 0, Pos: vec.Vec2{X: 400, Y: 300}, Spin: vec.Vec2{X: -4, Y: 0}, Active: true}
	balls := []*ball.Ball{b}

	for i := 0; i < 2000 && !AllStopped(balls); i++ {
		Step(balls)
	}

	if b.Pos.X >= 400 {
		t.Errorf("backspin did not draw the ball back: x = %.3f (want < 400)", b.Pos.X)
	}
	if b.Moving() {
		t.Error("ball with spin never came to rest")
	}
}

// A ball driven into a cushion tip beside a pocket should rattle back out.
func TestJawPegRejectsBall(t *testing.T) {
	peg := JawPegs[0] // top-left corner tip on the top rail
	b := &ball.Ball{Number: 1, Pos: vec.Vec2{X: peg.X, Y: peg.Y + ball.Radius - 1}, Vel: vec.Vec2{X: 0, Y: -3}, Active: true}

	if s := ReflectOffPeg(b, peg); s <= 0 {
		t.Fatalf("expected a jaw collision, got impact %.3f", s)
	}
	if b.Vel.Y <= 0 {
		t.Errorf("ball did not bounce off the jaw: vy = %.3f (want > 0)", b.Vel.Y)
	}
}

func TestMouthPredicates(t *testing.T) {
	midX := float64(table.PlayLeft+table.PlayRight) / 2
	if !InTopBottomMouth(midX) {
		t.Error("side pocket mouth not detected at midX")
	}
	if InTopBottomMouth(midX + SideMouth + 10) {
		t.Error("rail wrongly reported open away from any pocket")
	}
	if !InLeftRightMouth(table.PlayTop + 1) {
		t.Error("corner mouth not detected near the top of a side rail")
	}
}
