package rules

import "testing"

func TestOf(t *testing.T) {
	cases := map[int]Group{
		0: GroupNone, 1: GroupSolids, 7: GroupSolids,
		8: GroupNone, 9: GroupStripes, 15: GroupStripes,
	}
	for n, want := range cases {
		if got := Of(n); got != want {
			t.Errorf("Of(%d) = %v, want %v", n, got, want)
		}
	}
}
