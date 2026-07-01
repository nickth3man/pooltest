package ptmath

import (
	"math"
)

// SolveQuartic finds real roots of a t⁴ + b t³ + c t² + d t + e = 0.
// Returns sorted real roots (may be empty).
func SolveQuartic(a, b, c, d, e float64) []float64 {
	if math.Abs(a) < Eps {
		return solveCubic(b, c, d, e)
	}
	// Depressed quartic via substitution t = x - b/(4a).
	bb := b / a
	cc := c / a
	dd := d / a
	ee := e / a
	p := cc - 3*bb*bb/8
	q := dd + bb*bb*bb/8 - bb*cc/2
	r := ee - 3*bb*bb*bb*bb/256 + bb*bb*cc/16 - bb*dd/4
	shift := bb / 4

	cubicRoots := solveCubic(1, 2*p, p*p-4*r, -q*q)
	if len(cubicRoots) == 0 {
		return nil
	}
	y := cubicRoots[0]
	for _, root := range cubicRoots[1:] {
		if root > y {
			y = root
		}
	}
	sqrtTerm := math.Sqrt(math.Max(0, 2*y))
	if sqrtTerm < Eps {
		return nil
	}
	var roots []float64
	for _, sign := range []float64{-1, 1} {
		num := sign*sqrtTerm - q/sqrtTerm
		disc := 2*y - num*num/4
		if disc < 0 {
			continue
		}
		s := math.Sqrt(disc)
		for _, tSign := range []float64{-1, 1} {
			x := -shift + num/2 + tSign*s/2
			if !math.IsNaN(x) && !math.IsInf(x, 0) {
				roots = append(roots, x)
			}
		}
	}
	return dedupeSorted(roots)
}

func solveCubic(a, b, c, d float64) []float64 {
	if math.Abs(a) < Eps {
		t1, t2, ok := SolveQuadratic(b, c, d)
		if !ok {
			return nil
		}
		return dedupeSorted([]float64{t1, t2})
	}
	// Cardano for x³ + px + q = 0 after depression.
	bb := b / a
	cc := c / a
	dd := d / a
	p := cc - bb*bb/3
	q := 2*bb*bb*bb/27 - bb*cc/3 + dd
	disc := q*q/4 + p*p*p/27
	if disc > Eps {
		u := math.Cbrt(-q/2 + math.Sqrt(disc))
		v := math.Cbrt(-q/2 - math.Sqrt(disc))
		return []float64{u + v - bb/3}
	}
	if math.Abs(disc) <= Eps {
		u := math.Cbrt(-q / 2)
		return dedupeSorted([]float64{2*u - bb/3, -u - bb/3})
	}
	r := math.Sqrt(-p * p * p / 27)
	phi := math.Acos(-q / (2 * r))
	m := 2 * math.Cbrt(r)
	return dedupeSorted([]float64{
		m*math.Cos(phi/3) - bb/3,
		m*math.Cos((phi+2*math.Pi)/3) - bb/3,
		m*math.Cos((phi+4*math.Pi)/3) - bb/3,
	})
}

func dedupeSorted(vals []float64) []float64 {
	if len(vals) == 0 {
		return nil
	}
	out := make([]float64, 0, len(vals))
	for _, v := range vals {
		dup := false
		for _, o := range out {
			if math.Abs(v-o) < 1e-7 {
				dup = true
				break
			}
		}
		if !dup {
			out = append(out, v)
		}
	}
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j] < out[i] {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}

// EarliestPositiveRoot returns the smallest root > eps.
func EarliestPositiveRoot(roots []float64) (float64, bool) {
	best := math.Inf(1)
	ok := false
	for _, t := range roots {
		if t > Eps && t < best {
			best, ok = t, true
		}
	}
	return best, ok
}
