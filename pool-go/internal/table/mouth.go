package table

import "math"

const (
	// CornerMouth is the pocket opening width at corner rails (pixels).
	CornerMouth = 30.0
	// SideMouth is half the side-pocket mouth width on top/bottom rails (pixels).
	SideMouth = 26.0
)

// InTopBottomMouth reports whether x lies in a top/bottom pocket opening.
func InTopBottomMouth(x float64) bool {
	return x <= float64(PlayLeft)+CornerMouth || x >= float64(PlayRight)-CornerMouth ||
		math.Abs(x-MidX) <= SideMouth
}

// InLeftRightMouth reports whether y lies in a side-rail pocket opening.
func InLeftRightMouth(y float64) bool {
	return y <= float64(PlayTop)+CornerMouth || y >= float64(PlayBottom)-CornerMouth
}
