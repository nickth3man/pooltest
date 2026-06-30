package main

import (
	"math"
	"testing"
)

func TestNewRackLayout(t *testing.T) {
	balls := newRack()
	if len(balls) != 16 {
		t.Fatalf("rack has %d balls, want 16", len(balls))
	}
	if balls[0].Number != 0 {
		t.Errorf("balls[0] is %d, want the cue ball (0)", balls[0].Number)
	}

	seen := map[int]bool{}
	for _, b := range balls {
		if !b.Active {
			t.Errorf("ball %d starts inactive", b.Number)
		}
		seen[b.Number] = true
	}
	for n := 0; n <= 15; n++ {
		if !seen[n] {
			t.Errorf("ball %d missing from rack", n)
		}
	}

	// No two racked balls may overlap.
	for i := range balls {
		for j := i + 1; j < len(balls); j++ {
			if d := balls[i].Pos.Sub(balls[j].Pos).Len(); d < 2*ballRadius {
				t.Errorf("balls %d and %d overlap (dist %.2f < %.2f)",
					balls[i].Number, balls[j].Number, d, 2*ballRadius)
			}
		}
	}
}

func TestGroupOf(t *testing.T) {
	cases := map[int]Group{
		0: GroupNone, 1: GroupSolids, 7: GroupSolids,
		8: GroupNone, 9: GroupStripes, 15: GroupStripes,
	}
	for n, want := range cases {
		if got := groupOf(n); got != want {
			t.Errorf("groupOf(%d) = %v, want %v", n, got, want)
		}
	}
}

func TestHeadOnCollisionTransfersMomentum(t *testing.T) {
	a := &Ball{Number: 0, Pos: Vec2{100, 100}, Vel: Vec2{1, 0}, Active: true}
	b := &Ball{Number: 1, Pos: Vec2{100 + 2*ballRadius - 0.5, 100}, Active: true}

	resolveBallCollisions([]*Ball{a, b})

	if a.Vel.X > 0.1 {
		t.Errorf("striking ball kept too much speed: %.3f", a.Vel.X)
	}
	if b.Vel.X < 0.8 {
		t.Errorf("struck ball gained too little speed: %.3f", b.Vel.X)
	}
}

func TestRailReflectsVelocity(t *testing.T) {
	b := &Ball{Number: 0, Pos: Vec2{playRight - ballRadius - 0.1, 300}, Vel: Vec2{5, 0}, Active: true}

	step([]*Ball{b})

	if b.Vel.X >= 0 {
		t.Errorf("ball did not bounce off the right rail: vx = %.3f", b.Vel.X)
	}
	if b.Pos.X > playRight-ballRadius {
		t.Errorf("ball escaped the cushion: x = %.3f > %.3f", b.Pos.X, playRight-ballRadius)
	}
}

func TestPocketCapture(t *testing.T) {
	p := pocketCenters[0]
	b := &Ball{Number: 3, Pos: p, Active: true}

	got := step([]*Ball{b})

	if b.Active {
		t.Error("ball over a pocket was not captured")
	}
	if len(got) != 1 || got[0] != 3 {
		t.Errorf("step returned %v, want [3]", got)
	}
}

func TestFrictionBringsBallToRest(t *testing.T) {
	b := &Ball{Number: 0, Pos: Vec2{400, 300}, Vel: Vec2{6, 0}, Active: true}
	balls := []*Ball{b}

	for i := 0; i < 2000 && !allStopped(balls); i++ {
		step(balls)
	}

	if b.moving() {
		t.Errorf("ball never came to rest: speed = %.4f", math.Sqrt(b.Vel.LenSq()))
	}
}

// A stationary ball carrying a backspin reservoir should crawl backward (draw)
// as the cloth grips it, then settle. This is the mechanic behind draw shots.
func TestSpinReservoirMovesAndSettles(t *testing.T) {
	b := &Ball{Number: 0, Pos: Vec2{400, 300}, Spin: Vec2{-4, 0}, Active: true}
	balls := []*Ball{b}

	for i := 0; i < 2000 && !allStopped(balls); i++ {
		step(balls)
	}

	if b.Pos.X >= 400 {
		t.Errorf("backspin did not draw the ball back: x = %.3f (want < 400)", b.Pos.X)
	}
	if b.moving() {
		t.Error("ball with spin never came to rest")
	}
}

// A ball driven into a cushion tip beside a pocket should rattle back out.
func TestJawPegRejectsBall(t *testing.T) {
	peg := jawPegs[0] // top-left corner tip on the top rail
	b := &Ball{Number: 1, Pos: Vec2{peg.X, peg.Y + ballRadius - 1}, Vel: Vec2{0, -3}, Active: true}

	if s := reflectOffPeg(b, peg); s <= 0 {
		t.Fatalf("expected a jaw collision, got impact %.3f", s)
	}
	if b.Vel.Y <= 0 {
		t.Errorf("ball did not bounce off the jaw: vy = %.3f (want > 0)", b.Vel.Y)
	}
}

func TestMouthPredicates(t *testing.T) {
	midX := float64(playLeft+playRight) / 2
	if !inTopBottomMouth(midX) {
		t.Error("side pocket mouth not detected at midX")
	}
	if inTopBottomMouth(midX + sideMouth + 10) {
		t.Error("rail wrongly reported open away from any pocket")
	}
	if !inLeftRightMouth(playTop + 1) {
		t.Error("corner mouth not detected near the top of a side rail")
	}
}

func TestFirstContactFindsGhostBall(t *testing.T) {
	g := &Game{}
	g.cue = &Ball{Number: 0, Pos: Vec2{100, 300}, Active: true}
	obj := &Ball{Number: 1, Pos: Vec2{200, 300}, Active: true}
	g.balls = []*Ball{g.cue, obj}

	ghost, target, ok := g.firstContact(Vec2{1, 0})
	if !ok || target != obj {
		t.Fatalf("firstContact missed the object ball (ok=%v)", ok)
	}
	if want := 200.0 - 2*ballRadius; math.Abs(ghost.X-want) > 0.5 {
		t.Errorf("ghost ball at x=%.2f, want %.2f", ghost.X, want)
	}
}
