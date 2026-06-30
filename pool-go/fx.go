package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// particle is a short-lived spark or dust mote used for impact and pocket feedback.
type particle struct {
	pos, vel Vec2
	life     float64 // counts down to 0
	max      float64
	radius   float64
	clr      color.RGBA
}

// spawnImpact scatters sparks at a collision point, scaled by impact strength,
// and adds a touch of screen shake on hard hits.
func (g *Game) spawnImpact(pos Vec2, strength float64) {
	n := int(math.Min(10, 2+strength))
	for i := 0; i < n; i++ {
		ang := rand.Float64() * 2 * math.Pi
		spd := 0.6 + rand.Float64()*strength*0.35
		g.particles = append(g.particles, particle{
			pos:    pos,
			vel:    Vec2{math.Cos(ang), math.Sin(ang)}.Scale(spd),
			life:   8 + rand.Float64()*8,
			max:    16,
			radius: 1.2 + rand.Float64()*1.3,
			clr:    color.RGBA{0xFF, 0xF2, 0xC0, 0xFF},
		})
	}
	g.shake = math.Min(6, g.shake+strength*0.35)
}

// spawnPuff emits a small ring of felt dust when a ball drops into a pocket.
func (g *Game) spawnPuff(pos Vec2) {
	for i := 0; i < 8; i++ {
		ang := rand.Float64() * 2 * math.Pi
		g.particles = append(g.particles, particle{
			pos:    pos,
			vel:    Vec2{math.Cos(ang), math.Sin(ang)}.Scale(0.4 + rand.Float64()*0.8),
			life:   10 + rand.Float64()*10,
			max:    20,
			radius: 1.5 + rand.Float64()*1.5,
			clr:    color.RGBA{0x2A, 0x6B, 0x3A, 0xFF},
		})
	}
}

// tickFX advances particles, screen shake, and pocket-drop animations. It runs
// every frame regardless of game state so feedback keeps animating between shots.
func (g *Game) tickFX() {
	g.shake *= 0.85
	if g.shake < 0.1 {
		g.shake = 0
	}

	live := g.particles[:0]
	for _, p := range g.particles {
		p.pos = p.pos.Add(p.vel)
		p.vel = p.vel.Scale(0.9)
		p.life--
		if p.life > 0 {
			live = append(live, p)
		}
	}
	g.particles = live

	for _, b := range g.balls {
		if b.Sinking {
			b.SinkT += 0.12
			if b.SinkT >= 1 {
				b.Sinking = false
			}
		}
	}
}

func (g *Game) drawParticles(dst *ebiten.Image) {
	for _, p := range g.particles {
		a := p.life / p.max
		c := p.clr
		c.A = clamp8(float64(c.A) * a)
		vector.FillCircle(dst, float32(p.pos.X), float32(p.pos.Y), float32(p.radius), c, true)
	}
}

// drawSinking renders balls mid-drop, shrinking and easing toward the pocket center.
func (g *Game) drawSinking(dst *ebiten.Image) {
	ensureSprites()
	for _, b := range g.balls {
		if !b.Sinking {
			continue
		}
		t := b.SinkT
		pos := b.SinkFrom.Add(b.SinkPos.Sub(b.SinkFrom).Scale(t))
		scale := 1 - t
		sprite := g.ballSprite(b)
		w, _ := sprite.Bounds().Dx(), sprite.Bounds().Dy()
		half := float64(w) / 2
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-half, -half)
		op.GeoM.Scale(scale, scale)
		op.GeoM.Rotate(b.Angle)
		op.GeoM.Translate(pos.X, pos.Y)
		op.ColorScale.ScaleAlpha(float32(scale))
		dst.DrawImage(sprite, op)
	}
}
