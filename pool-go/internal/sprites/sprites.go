// Package sprites builds the per-ball body sprites (shaded sphere with the
// printed number) and the shared specular / shadow sprites.
package sprites

import (
	"image"
	"image/color"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/user/pooltest/pool-go/internal/ball"
)

// Caches. Ball bodies are keyed by number (color/stripe derive from it), so
// they survive re-racks. All are built lazily on first Draw because creating
// ebiten.Images before the game loop starts is not allowed.
var (
	ballSprites  = map[int]*ebiten.Image{}
	specSprite   *ebiten.Image
	shadowSprite *ebiten.Image
)

// Clamp8 saturates v to the 0-255 range for image color channels.
func Clamp8(v float64) uint8 {
	if v <= 0 {
		return 0
	}
	if v >= 255 {
		return 255
	}
	return uint8(v)
}

// Ball returns the cached shaded body for a ball, building it on demand. The
// face is the font used to print the ball's number on the sprite.
func Ball(b *ball.Ball, face text.Face) *ebiten.Image {
	if img, ok := ballSprites[b.Number]; ok {
		return img
	}
	img := buildBallBody(b, face)
	ballSprites[b.Number] = img
	return img
}

// buildBallBody renders one pool ball as a lit sphere. The diffuse shading is
// radially symmetric (bright center → dark rim) so it is rotation-invariant —
// spinning the sprite then only visibly turns the number/stripe, which is the
// roll cue. Directional gloss is added separately by the fixed specular sprite.
func buildBallBody(b *ball.Ball, face text.Face) *ebiten.Image {
	r := ball.Radius
	d := int(math.Ceil(r*2)) + 2
	c := float64(d) / 2
	rgba := image.NewRGBA(image.Rect(0, 0, d, d))

	for py := range d {
		for px := range d {
			nx := (float64(px) + 0.5 - c) / r
			ny := (float64(py) + 0.5 - c) / r
			rr := nx*nx + ny*ny
			if rr > 1 {
				continue
			}
			shade := 1.0 - 0.45*rr
			if rr > 0.82 { // soft contact darkening at the very edge
				shade *= 1 - (rr-0.82)/0.18*0.45
			}
			base := b.Color
			switch {
			case b.Number == 0:
				base = color.RGBA{0xF4, 0xF1, 0xE8, 0xFF} // ivory cue ball
			case b.Stripe && math.Abs(ny) > 0.42:
				base = color.RGBA{0xF2, 0xF0, 0xE6, 0xFF} // white poles of a stripe
			}
			a := 1.0
			if rr > 0.92 { // antialias the silhouette
				a = (1 - rr) / 0.08
			}
			// Go's color.RGBA is alpha-premultiplied: fold shade and alpha in.
			rgba.SetRGBA(px, py, color.RGBA{
				Clamp8(float64(base.R) * shade * a),
				Clamp8(float64(base.G) * shade * a),
				Clamp8(float64(base.B) * shade * a),
				Clamp8(a * 255),
			})
		}
	}

	img := ebiten.NewImageFromImage(rgba)
	if b.Number == 0 {
		// A small spot so the cue ball's english is visible as it spins.
		vector.FillCircle(img, float32(c), float32(c-r*0.42), float32(r*0.16),
			color.RGBA{0xD0, 0x3A, 0x2A, 0xFF}, true)
		return img
	}

	// White number patch with the printed number.
	vector.FillCircle(img, float32(c), float32(c), float32(r*0.56),
		color.RGBA{0xF7, 0xF5, 0xEC, 0xFF}, true)
	op := &text.DrawOptions{}
	op.GeoM.Translate(c, c)
	op.PrimaryAlign = text.AlignCenter
	op.SecondaryAlign = text.AlignCenter
	op.ColorScale.ScaleWithColor(color.RGBA{0x20, 0x20, 0x20, 0xFF})
	text.Draw(img, strconv.Itoa(b.Number), face, op)
	return img
}

// EnsureSprites builds the shared specular and shadow sprites once.
func EnsureSprites() {
	if specSprite == nil {
		specSprite = buildSoftCircle(ball.Radius*0.62,
			color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}, 0.85, 1.6)
	}
	if shadowSprite == nil {
		shadowSprite = buildSoftCircle(ball.Radius*1.15,
			color.RGBA{0x00, 0x00, 0x00, 0xFF}, 0.40, 1.4)
	}
}

// SpecSprite returns the cached shared specular sprite (or nil if not yet
// built — callers should call EnsureSprites first).
func SpecSprite() *ebiten.Image { return specSprite }

// ShadowSprite returns the cached shared shadow sprite (or nil if not yet
// built — callers should call EnsureSprites first).
func ShadowSprite() *ebiten.Image { return shadowSprite }

// buildSoftCircle renders a radial gradient disc (opaque-ish center fading to
// transparent edge) in premultiplied alpha, used for highlights and shadows.
func buildSoftCircle(rad float64, col color.RGBA, maxA, pow float64) *ebiten.Image {
	d := int(math.Ceil(rad*2)) + 2
	c := float64(d) / 2
	rgba := image.NewRGBA(image.Rect(0, 0, d, d))
	for py := range d {
		for px := range d {
			nx := (float64(px) + 0.5 - c) / rad
			ny := (float64(py) + 0.5 - c) / rad
			rr := nx*nx + ny*ny
			if rr > 1 {
				continue
			}
			a := math.Pow(1-rr, pow) * maxA
			rgba.SetRGBA(px, py, color.RGBA{
				Clamp8(float64(col.R) * a),
				Clamp8(float64(col.G) * a),
				Clamp8(float64(col.B) * a),
				Clamp8(a * 255),
			})
		}
	}
	return ebiten.NewImageFromImage(rgba)
}
