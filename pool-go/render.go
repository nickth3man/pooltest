package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	colorBackground = color.RGBA{0x10, 0x18, 0x12, 0xFF}
	colorRail       = color.RGBA{0x5A, 0x37, 0x1E, 0xFF}
	colorRailLight  = color.RGBA{0x7A, 0x50, 0x30, 0xFF}
	colorRailDark   = color.RGBA{0x3C, 0x22, 0x12, 0xFF}
	colorFelt       = color.RGBA{0x0F, 0x6E, 0x2F, 0xFF}
	colorPocket     = color.RGBA{0x04, 0x07, 0x05, 0xFF}
	colorPocketRim  = color.RGBA{0x2A, 0x18, 0x0E, 0xFF}
	colorSpot       = color.RGBA{0xCF, 0xE0, 0xC8, 0x66}
	colorWhite      = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	colorHUD        = color.RGBA{0xEC, 0xEC, 0xEC, 0xFF}
	colorHUDDim     = color.RGBA{0x9A, 0xA6, 0x9E, 0xFF}
	colorActive     = color.RGBA{0xFF, 0xD7, 0x4A, 0xFF}
	colorGuide      = color.RGBA{0xFF, 0xFF, 0xFF, 0x55}
	colorGhost      = color.RGBA{0xFF, 0xFF, 0xFF, 0x99}
	colorCarom      = color.RGBA{0xBF, 0xE3, 0xFF, 0x77}
	colorCueWood    = color.RGBA{0xC8, 0x9B, 0x5A, 0xFF}
	colorCueButt    = color.RGBA{0x3A, 0x24, 0x16, 0xFF}
	colorCueTip     = color.RGBA{0x3C, 0x6E, 0xC8, 0xFF}
	colorMarker     = color.RGBA{0xD0, 0x3A, 0x2A, 0xFF}
)

const spinWidgetR = 30.0

var spinWidgetCenter = Vec2{56, screenH - 56}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.canvas == nil {
		g.canvas = ebiten.NewImage(screenW, screenH)
	}
	cv := g.canvas
	cv.Fill(colorBackground)

	g.drawTable(cv)
	g.drawBalls(cv)
	g.drawSinking(cv)
	g.drawParticles(cv)

	if g.state == stateAiming {
		g.drawSpinWidget(cv)
		if g.charging {
			g.drawAim(cv)
		}
	}
	if g.state == stateBallInHand {
		g.drawCuePreview(cv)
	}
	g.drawHUD(cv)

	// Blit the composed frame, jittered while the table is shaking.
	screen.Fill(colorBackground)
	op := &ebiten.DrawImageOptions{}
	if g.shake > 0 {
		op.GeoM.Translate((rand.Float64()*2-1)*g.shake, (rand.Float64()*2-1)*g.shake)
	}
	screen.DrawImage(cv, op)
}

func (g *Game) drawTable(dst *ebiten.Image) {
	w := float32(playRight - playLeft)
	h := float32(playBottom - playTop)

	// Rail frame with a lit top/left edge and shaded bottom/right edge for depth.
	vector.FillRect(dst, playLeft-railWidth, playTop-railWidth, w+2*railWidth, h+2*railWidth, colorRail, false)
	vector.FillRect(dst, playLeft-railWidth, playTop-railWidth, w+2*railWidth, 3, colorRailLight, false)
	vector.FillRect(dst, playLeft-railWidth, playTop-railWidth, 3, h+2*railWidth, colorRailLight, false)
	vector.FillRect(dst, playLeft-railWidth, playBottom+railWidth-3, w+2*railWidth, 3, colorRailDark, false)
	vector.FillRect(dst, playRight+railWidth-3, playTop-railWidth, 3, h+2*railWidth, colorRailDark, false)

	// Felt with a soft central sheen.
	vector.FillRect(dst, playLeft, playTop, w, h, colorFelt, false)
	ensureSprites()
	sw := float64(specSprite.Bounds().Dx())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-sw/2, -sw/2)
	op.GeoM.Scale(34, 22)
	op.GeoM.Translate(float64(playLeft+playRight)/2, float64(playTop+playBottom)/2)
	op.ColorScale.ScaleAlpha(0.06)
	dst.DrawImage(specSprite, op)

	// Table spots and rail sights.
	headX := float64(playLeft) + 0.25*float64(playRight-playLeft)
	footX := float64(playLeft) + 0.75*float64(playRight-playLeft)
	centerY := float64(playTop+playBottom) / 2
	drawSpot(dst, headX, centerY)
	drawSpot(dst, footX, centerY)
	g.drawSights(dst)

	// Deep pockets last so they sit on top of felt and sights.
	for _, p := range pockets() {
		drawPocket(dst, p)
	}
}

func drawSpot(dst *ebiten.Image, x, y float64) {
	vector.FillCircle(dst, float32(x), float32(y), 2.2, colorSpot, true)
}

func drawPocket(dst *ebiten.Image, p Vec2) {
	x, y := float32(p.X), float32(p.Y)
	vector.FillCircle(dst, x, y, pocketRadius+3, colorPocketRim, true)
	vector.FillCircle(dst, x, y, pocketRadius, colorPocket, true)
	vector.FillCircle(dst, x, y, pocketRadius*0.7, color.RGBA{0, 0, 0, 0xFF}, true)
}

func (g *Game) drawSights(dst *ebiten.Image) {
	topY := float64(playTop) - railWidth/2
	botY := float64(playBottom) + railWidth/2
	leftX := float64(playLeft) - railWidth/2
	rightX := float64(playRight) + railWidth/2
	for _, f := range []float64{0.125, 0.375, 0.625, 0.875} {
		x := float64(playLeft) + f*float64(playRight-playLeft)
		drawSpot(dst, x, topY)
		drawSpot(dst, x, botY)
	}
	for _, f := range []float64{0.25, 0.5, 0.75} {
		y := float64(playTop) + f*float64(playBottom-playTop)
		drawSpot(dst, leftX, y)
		drawSpot(dst, rightX, y)
	}
}

func (g *Game) drawBalls(dst *ebiten.Image) {
	ensureSprites()
	for _, b := range g.balls {
		if b.Active {
			g.drawBall(dst, b)
		}
	}
}

// drawBall composites a ball: a flattened drop shadow, the rotated shaded body,
// and a fixed glossy highlight (so the gloss stays put while the ball rolls).
func (g *Game) drawBall(dst *ebiten.Image, b *Ball) {
	cx, cy := b.Pos.X, b.Pos.Y

	shHalf := float64(shadowSprite.Bounds().Dx()) / 2
	sop := &ebiten.DrawImageOptions{}
	sop.GeoM.Translate(-shHalf, -shHalf)
	sop.GeoM.Scale(1.0, 0.6)
	sop.GeoM.Translate(cx+2.5, cy+3.5)
	dst.DrawImage(shadowSprite, sop)

	sprite := g.ballSprite(b)
	bHalf := float64(sprite.Bounds().Dx()) / 2
	bop := &ebiten.DrawImageOptions{}
	bop.GeoM.Translate(-bHalf, -bHalf)
	bop.GeoM.Rotate(b.Angle)
	bop.GeoM.Translate(cx, cy)
	dst.DrawImage(sprite, bop)

	spHalf := float64(specSprite.Bounds().Dx()) / 2
	pop := &ebiten.DrawImageOptions{}
	pop.GeoM.Translate(-spHalf, -spHalf)
	pop.GeoM.Translate(cx-ballRadius*0.32, cy-ballRadius*0.38)
	dst.DrawImage(specSprite, pop)
}

// drawAim draws the shot prediction (cue path, ghost ball at first contact, and
// the object ball's carom), the tapered cue stick, and the power meter.
func (g *Game) drawAim(dst *ebiten.Image) {
	mx, my := ebiten.CursorPosition()
	pull := g.cue.Pos.Sub(Vec2{float64(mx), float64(my)})
	dist := pull.Len()
	if dist < 1 {
		return
	}
	dir := pull.Normalize()
	power := math.Min(dist, maxPull)
	frac := power / maxPull

	end, target, hit := g.firstContact(dir)
	if !hit {
		end = rayToRail(g.cue.Pos, dir)
	}
	clr := powerColor(frac)
	vector.StrokeLine(dst, float32(g.cue.Pos.X), float32(g.cue.Pos.Y),
		float32(end.X), float32(end.Y), 2, clr, true)

	if hit {
		vector.StrokeCircle(dst, float32(end.X), float32(end.Y), ballRadius, 1.5, colorGhost, true)
		carom := target.Pos.Sub(end)
		if carom.LenSq() > 0 {
			c2 := target.Pos.Add(carom.Normalize().Scale(60))
			vector.StrokeLine(dst, float32(target.Pos.X), float32(target.Pos.Y),
				float32(c2.X), float32(c2.Y), 1.5, colorCarom, true)
		}
	}

	g.drawCueStick(dst, dir, power)
	g.drawPowerMeter(dst, frac)
}

// drawCueStick draws a tapered stick (tip, ferrule, shaft, wrapped butt) pulled
// back from the cue ball on the mouse side, the gap growing with shot power.
func (g *Game) drawCueStick(dst *ebiten.Image, dir Vec2, power float64) {
	back := dir.Scale(-1)
	gap := ballRadius + 6 + power*0.5
	at := func(d float64) Vec2 { return g.cue.Pos.Add(back.Scale(gap + d)) }
	line := func(a, b Vec2, wdt float32, c color.RGBA) {
		vector.StrokeLine(dst, float32(a.X), float32(a.Y), float32(b.X), float32(b.Y), wdt, c, true)
	}
	line(at(95), at(135), 6, colorCueButt) // wrapped butt
	line(at(8), at(95), 4, colorCueWood)   // shaft
	line(at(4), at(9), 4.5, colorWhite)    // ferrule
	line(at(0), at(4), 4.5, colorCueTip)   // tip
}

func (g *Game) drawPowerMeter(dst *ebiten.Image, frac float64) {
	x := float32(screenW - 24)
	y0 := float32(playTop)
	h := float32(playBottom - playTop)
	vector.StrokeRect(dst, x, y0, 12, h, 1, colorHUDDim, true)
	fh := h * float32(frac)
	vector.FillRect(dst, x, y0+h-fh, 12, fh, powerColor(frac), false)
}

func (g *Game) drawSpinWidget(dst *ebiten.Image) {
	cx, cy := float32(spinWidgetCenter.X), float32(spinWidgetCenter.Y)
	vector.FillCircle(dst, cx, cy, spinWidgetR, color.RGBA{0xEC, 0xEC, 0xEC, 0xCC}, true)
	vector.StrokeCircle(dst, cx, cy, spinWidgetR, 1.5, colorHUD, true)
	vector.StrokeLine(dst, cx-spinWidgetR, cy, cx+spinWidgetR, cy, 1, colorHUDDim, true)
	vector.StrokeLine(dst, cx, cy-spinWidgetR, cx, cy+spinWidgetR, 1, colorHUDDim, true)
	m := spinWidgetCenter.Add(g.spin.Scale(spinWidgetR))
	vector.FillCircle(dst, float32(m.X), float32(m.Y), 4.5, colorMarker, true)

	op := &text.DrawOptions{}
	op.GeoM.Translate(spinWidgetCenter.X, spinWidgetCenter.Y-spinWidgetR-8)
	op.PrimaryAlign = text.AlignCenter
	op.ColorScale.ScaleWithColor(colorHUDDim)
	text.Draw(dst, "English", g.smallFace, op)
}

func (g *Game) drawCuePreview(dst *ebiten.Image) {
	mx, my := ebiten.CursorPosition()
	pos := Vec2{float64(mx), float64(my)}
	c := color.RGBA{0xFF, 0xFF, 0xFF, 0x88}
	if !g.validCuePlacement(pos) {
		c = color.RGBA{0xFF, 0x44, 0x44, 0x88}
	}
	vector.FillCircle(dst, float32(pos.X), float32(pos.Y), ballRadius, c, true)
}

func (g *Game) drawHUD(dst *ebiten.Image) {
	g.drawPlayerLabel(dst, 0, 16, text.AlignStart)
	g.drawPlayerLabel(dst, 1, screenW-16, text.AlignEnd)

	msg := &text.DrawOptions{}
	msg.GeoM.Translate(screenW/2, screenH-30)
	msg.PrimaryAlign = text.AlignCenter
	msg.ColorScale.ScaleWithColor(colorHUD)
	text.Draw(dst, g.message, g.face, msg)

	hint := &text.DrawOptions{}
	hint.GeoM.Translate(screenW/2, screenH-12)
	hint.PrimaryAlign = text.AlignCenter
	hint.ColorScale.ScaleWithColor(colorGuide)
	text.Draw(dst, "Drag back from the cue ball to shoot · click the English dial for spin · right-click clears it", g.smallFace, hint)
}

func (g *Game) drawPlayerLabel(dst *ebiten.Image, idx int, x float64, align text.Align) {
	label := fmt.Sprintf("P%d: %s", idx+1, g.players[idx].Group)
	if grp := g.players[idx].Group; grp != GroupNone {
		label += fmt.Sprintf(" (%d)", g.groupRemaining(grp))
	}
	clr := color.Color(colorHUD)
	if g.current == idx && g.state != stateGameOver {
		clr = colorActive
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, 16)
	op.PrimaryAlign = align
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(dst, label, g.face, op)
}

// powerColor ramps green → yellow → red with shot power.
func powerColor(f float64) color.RGBA {
	f = math.Max(0, math.Min(1, f))
	var r, gr float64
	if f < 0.5 {
		r, gr = f*2, 1
	} else {
		r, gr = 1, 1-(f-0.5)*2
	}
	return color.RGBA{clamp8(r * 255), clamp8(gr * 255), 0x30, 0xFF}
}
