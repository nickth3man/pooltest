// Package fx provides short-lived visual feedback: impact sparks, pocket
// puffs, and screen shake. It owns the feedback state (particles, shake) but
// leaves drawing to the renderer.
package fx

import (
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/vec"
)

// Particle is a short-lived spark or dust mote used for impact and pocket
// feedback.
type Particle struct {
	Pos, Vel vec.Vec2
	Life     float64 // counts down to 0
	Max      float64
	Radius   float64
	Clr      color.RGBA
}

// State holds the per-frame feedback that the renderer reads (shake) and that
// the physics hooks append to (particles).
type State struct {
	Particles []Particle
	Shake     float64
}

// SpawnImpact scatters sparks at a collision point, scaled by impact strength,
// and adds a touch of screen shake on hard hits.
func SpawnImpact(s *State, pos vec.Vec2, strength float64) {
	n := int(math.Min(10, 2+strength))
	for range n {
		ang := rand.Float64() * 2 * math.Pi
		spd := 0.6 + rand.Float64()*strength*0.35
		s.Particles = append(s.Particles, Particle{
			Pos:    pos,
			Vel:    vec.Vec2{X: math.Cos(ang), Y: math.Sin(ang)}.Scale(spd),
			Life:   8 + rand.Float64()*8,
			Max:    16,
			Radius: 1.2 + rand.Float64()*1.3,
			Clr:    color.RGBA{0xFF, 0xF2, 0xC0, 0xFF},
		})
	}
	s.Shake = math.Min(6, s.Shake+strength*0.35)
}

// SpawnPuff emits a small ring of felt dust when a ball drops into a pocket.
func SpawnPuff(s *State, pos vec.Vec2) {
	for range 8 {
		ang := rand.Float64() * 2 * math.Pi
		s.Particles = append(s.Particles, Particle{
			Pos:    pos,
			Vel:    vec.Vec2{X: math.Cos(ang), Y: math.Sin(ang)}.Scale(0.4 + rand.Float64()*0.8),
			Life:   10 + rand.Float64()*10,
			Max:    20,
			Radius: 1.5 + rand.Float64()*1.5,
			Clr:    color.RGBA{0x2A, 0x6B, 0x3A, 0xFF},
		})
	}
}

// Tick advances particles, screen shake, and pocket-drop animations. It runs
// every frame regardless of game state so feedback keeps animating between
// shots.
func Tick(s *State, balls []*ball.Ball) {
	s.Shake *= 0.85
	if s.Shake < 0.1 {
		s.Shake = 0
	}

	live := s.Particles[:0]
	for _, p := range s.Particles {
		p.Pos = p.Pos.Add(p.Vel)
		p.Vel = p.Vel.Scale(0.9)
		p.Life--
		if p.Life > 0 {
			live = append(live, p)
		}
	}
	s.Particles = live

	for _, b := range balls {
		if b.Sinking {
			b.SinkT += 0.12
			if b.SinkT >= 1 {
				b.Sinking = false
			}
		}
	}
}
