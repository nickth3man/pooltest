package table

import (
	"math"

	"github.com/user/pooltest/pool-go/internal/vec"
)

// SegmentKind identifies cushion geometry.
type SegmentKind int

const (
	SegmentLinear SegmentKind = iota
	SegmentCircular
)

// CushionSegment is one rail or jaw piece for collision detection.
type CushionSegment struct {
	Kind SegmentKind
	// Linear: from P1 to P2; normal points inward to play area.
	P1, P2 vec.Vec2
	Normal vec.Vec2
	// Circular jaw: center and radius (pixels).
	Center vec.Vec2
	Radius float64
}

// PhysicsTable holds segment-based cushion geometry.
type PhysicsTable struct {
	Linear    []CushionSegment
	Circular  []CushionSegment
	PocketPts [6]vec.Vec2
	PocketR   float64
	PlayMinX  float64
	PlayMaxX  float64
	PlayMinY  float64
	PlayMaxY  float64
}

// DefaultPhysicsTable builds segments matching the visual table.
func DefaultPhysicsTable() *PhysicsTable {
	cr := CollisionRadiusPx()
	t := &PhysicsTable{
		PocketR:  PocketRadius,
		PlayMinX: float64(PlayLeft) + cr,
		PlayMaxX: float64(PlayRight) - cr,
		PlayMinY: float64(PlayTop) + cr,
		PlayMaxY: float64(PlayBottom) - cr,
	}
	t.PocketPts = PocketCenters

	cornerMouth := CornerMouth
	sideMouth := SideMouth
	jawR := 5.0

	addLinear := func(p1, p2, n vec.Vec2) {
		t.Linear = append(t.Linear, CushionSegment{Kind: SegmentLinear, P1: p1, P2: p2, Normal: n})
	}

	addLinear(vec.Vec2{X: float64(PlayLeft) + cornerMouth, Y: float64(PlayTop)}, vec.Vec2{X: MidX - sideMouth, Y: float64(PlayTop)}, vec.Vec2{Y: 1})
	addLinear(vec.Vec2{X: MidX + sideMouth, Y: float64(PlayTop)}, vec.Vec2{X: float64(PlayRight) - cornerMouth, Y: float64(PlayTop)}, vec.Vec2{Y: 1})
	addLinear(vec.Vec2{X: float64(PlayLeft) + cornerMouth, Y: float64(PlayBottom)}, vec.Vec2{X: MidX - sideMouth, Y: float64(PlayBottom)}, vec.Vec2{Y: -1})
	addLinear(vec.Vec2{X: MidX + sideMouth, Y: float64(PlayBottom)}, vec.Vec2{X: float64(PlayRight) - cornerMouth, Y: float64(PlayBottom)}, vec.Vec2{Y: -1})
	addLinear(vec.Vec2{X: float64(PlayLeft), Y: float64(PlayTop) + cornerMouth}, vec.Vec2{X: float64(PlayLeft), Y: float64(PlayBottom) - cornerMouth}, vec.Vec2{X: 1})
	addLinear(vec.Vec2{X: float64(PlayRight), Y: float64(PlayTop) + cornerMouth}, vec.Vec2{X: float64(PlayRight), Y: float64(PlayBottom) - cornerMouth}, vec.Vec2{X: -1})

	jaws := []vec.Vec2{
		{X: float64(PlayLeft) + cornerMouth, Y: float64(PlayTop)},
		{X: float64(PlayLeft), Y: float64(PlayTop) + cornerMouth},
		{X: float64(PlayRight) - cornerMouth, Y: float64(PlayTop)},
		{X: float64(PlayRight), Y: float64(PlayTop) + cornerMouth},
		{X: float64(PlayLeft) + cornerMouth, Y: float64(PlayBottom)},
		{X: float64(PlayLeft), Y: float64(PlayBottom) - cornerMouth},
		{X: float64(PlayRight) - cornerMouth, Y: float64(PlayBottom)},
		{X: float64(PlayRight), Y: float64(PlayBottom) - cornerMouth},
		{X: MidX - sideMouth, Y: float64(PlayTop)},
		{X: MidX + sideMouth, Y: float64(PlayTop)},
		{X: MidX - sideMouth, Y: float64(PlayBottom)},
		{X: MidX + sideMouth, Y: float64(PlayBottom)},
	}
	for _, c := range jaws {
		t.Circular = append(t.Circular, CushionSegment{Kind: SegmentCircular, Center: c, Radius: jawR})
	}
	return t
}

// CollisionRadiusPx is physics ball radius in pixels.
func CollisionRadiusPx() float64 {
	const rM = 0.028575
	const playW = float64(PlayRight - PlayLeft)
	pxPerM := playW / 1.9812
	return rM * pxPerM
}

// RayToRail returns where a ray from start in direction dir meets the play bounds.
func (t *PhysicsTable) RayToRail(start, dir vec.Vec2) vec.Vec2 {
	inv := func(a, d, lo, hi float64) float64 {
		switch {
		case d > 1e-9:
			return (hi - a) / d
		case d < -1e-9:
			return (lo - a) / d
		default:
			return math.Inf(1)
		}
	}
	dt := minF(inv(start.X, dir.X, t.PlayMinX, t.PlayMaxX), inv(start.Y, dir.Y, t.PlayMinY, t.PlayMaxY))
	if math.IsInf(dt, 1) {
		return start
	}
	return start.Add(dir.Scale(dt))
}

func minF(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
