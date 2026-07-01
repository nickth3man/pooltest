package rack

import (
	"testing"

	"github.com/user/pooltest/pool-go/internal/ball"
)

func TestNewRackLayout(t *testing.T) {
	balls := NewRack()
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
	for n := range 16 {
		if !seen[n] {
			t.Errorf("ball %d missing from rack", n)
		}
	}

	// No two racked balls may overlap.
	for i := range balls {
		for j := i + 1; j < len(balls); j++ {
			if d := balls[i].Pos.Sub(balls[j].Pos).Len(); d < 2*ball.Radius {
				t.Errorf("balls %d and %d overlap (dist %.2f < %.2f)",
					balls[i].Number, balls[j].Number, d, 2*ball.Radius)
			}
		}
	}
}
