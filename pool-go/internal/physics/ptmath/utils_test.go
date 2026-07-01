package ptmath

import (
	"math"
	"testing"

	"github.com/user/pooltest/pool-go/internal/vec"
)

func TestGetSlideRollSpinTimes(t *testing.T) {
	rvw := ballRVW{
		V: vec.Vec3{X: 1},
		W: vec.Vec3{Z: 10},
	}
	R, us, ur, usp, g := 0.028575, 0.2, 0.01, 0.063, 9.81
	ts := GetSlideTime(rvw, R, us, g)
	tr := GetRollTime(rvw, ur, g)
	tsp := GetSpinTime(rvw, R, usp, g)
	if ts <= 0 || math.IsInf(ts, 1) {
		t.Fatalf("slide time = %v", ts)
	}
	if tr <= 0 {
		t.Fatalf("roll time = %v", tr)
	}
	if tsp <= 0 {
		t.Fatalf("spin time = %v", tsp)
	}
}

func TestSolveQuadratic(t *testing.T) {
	t1, t2, ok := SolveQuadratic(1, -3, 2)
	if !ok {
		t.Fatal("expected roots")
	}
	if math.Abs(t1-1) > 1e-9 || math.Abs(t2-2) > 1e-9 {
		t.Fatalf("roots = %v %v", t1, t2)
	}
}

func TestSolveQuarticPositiveRoot(t *testing.T) {
	// t² - 5t + 6 = 0 → t=2,3
	roots := SolveQuartic(0, 1, -5, 6, 0)
	if len(roots) == 0 {
		t.Fatal("expected roots")
	}
	tBest, ok := EarliestPositiveRoot(roots)
	if !ok || tBest < 1.9 || tBest > 2.1 {
		t.Fatalf("earliest root = %v want ~2", tBest)
	}
}

func TestEarliestBallBallTime(t *testing.T) {
	p1 := vec.Vec3{X: 0}
	p2 := vec.Vec3{X: 0.1}
	v1 := vec.Vec3{X: 1}
	v2 := vec.Vec3{}
	sumR := 0.05715
	tColl, ok := EarliestBallBallTime(p1, v1, p2, v2, sumR)
	if !ok || tColl <= 0 {
		t.Fatalf("no collision time: %v ok=%v", tColl, ok)
	}
}
