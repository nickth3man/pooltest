package main

import (
	"bytes"
	"fmt"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

type gameState int

const (
	stateAiming gameState = iota
	stateShooting
	stateBallInHand
	stateGameOver
)

const (
	maxPull  = 140.0 // pull-back distance (px) that maps to full power
	maxSpeed = 18.0  // cue-ball launch speed (px/frame) at full power
	minPull  = 6.0   // shorter drags are ignored

	// Spin gains translate a unit dial offset into the cue ball's spin reservoir.
	followGain  = 0.55 // top/back-spin strength (draw and follow)
	swerveGain  = 0.30 // sideways carry that bends the cue ball's path
	spinVisGain = 0.60 // english fed into the visible roll
)

type player struct {
	Group Group
}

// Game holds the entire match state and implements ebiten.Game.
type Game struct {
	balls   []*Ball
	cue     *Ball
	players [2]player
	current int

	state        gameState
	tableOpen    bool
	charging     bool
	shotPocketed []int

	spin Vec2 // pending english from the dial: x = side, y = follow(-)/draw(+)

	face      *text.GoTextFace
	smallFace *text.GoTextFace

	message string

	// Presentation / feedback state.
	canvas    *ebiten.Image
	particles []particle
	shake     float64
}

// NewGame loads the font and racks a fresh match.
func NewGame() *Game {
	src, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
	}
	g := &Game{
		face:      &text.GoTextFace{Source: src, Size: 16},
		smallFace: &text.GoTextFace{Source: src, Size: 11},
	}
	initAudio()
	// Route physics events into feedback (particles, shake, sound).
	onBallImpact = func(pos Vec2, strength float64) {
		g.spawnImpact(pos, strength)
		playClick(strength)
	}
	onRailImpact = func(_ Vec2, strength float64) {
		playRail(strength)
		if strength > 6 {
			g.shake = math.Min(4, g.shake+1)
		}
	}
	onPocketDrop = func(pos Vec2) {
		g.spawnPuff(pos)
		playPocket()
	}
	g.reset()
	return g
}

func (g *Game) reset() {
	g.balls = newRack()
	g.cue = g.balls[0]
	g.players = [2]player{}
	g.current = 0
	g.state = stateAiming
	g.tableOpen = true
	g.charging = false
	g.shotPocketed = nil
	g.spin = Vec2{}
	g.particles = nil
	g.shake = 0
	g.message = "Player 1 to break"
}

func (g *Game) Update() error {
	g.tickFX()
	switch g.state {
	case stateAiming:
		g.updateAiming()
	case stateShooting:
		g.updateShooting()
	case stateBallInHand:
		g.updateBallInHand()
	case stateGameOver:
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.reset()
		}
	}
	return nil
}

// updateAiming implements a slingshot cue: hold the left button, drag back
// from the cue ball, and release to fire in the opposite direction with power
// proportional to the pull distance. The English dial in the corner intercepts
// clicks so the player can dial spin before taking the shot.
func (g *Game) updateAiming() {
	m := cursorVec()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.spin = Vec2{}
	}
	if g.inSpinWidget(m) {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			g.setSpinFromWidget(m)
		}
		return
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.charging = true
	}
	if !g.charging || !inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		return
	}

	g.charging = false
	pull := g.cue.Pos.Sub(m)
	dist := pull.Len()
	if dist < minPull {
		return
	}
	g.fire(pull, math.Min(dist, maxPull))
}

// fire launches the cue ball along pull with the given power and converts the
// dialed english into the cue ball's spin reservoir.
func (g *Game) fire(pull Vec2, power float64) {
	dir := pull.Normalize()
	perp := dir.Perp()
	frac := power / maxPull

	g.cue.Vel = dir.Scale(frac * maxSpeed)
	fb := -g.spin.Y // dial up = follow, down = draw
	sd := g.spin.X  // dial right = right english
	g.cue.Spin = dir.Scale(fb * followGain * frac * maxSpeed).
		Add(perp.Scale(sd * swerveGain * frac * maxSpeed))
	g.cue.SideSpin = sd * spinVisGain * frac
	g.cue.Angle = 0

	playCueStrike(frac)
	g.spin = Vec2{}
	g.shotPocketed = nil
	g.state = stateShooting
}

// inSpinWidget reports whether m lies within the English dial.
func (g *Game) inSpinWidget(m Vec2) bool {
	return m.Sub(spinWidgetCenter).Len() <= spinWidgetR+4
}

// setSpinFromWidget records the dialed english as a unit offset (clamped to the
// dial face).
func (g *Game) setSpinFromWidget(m Vec2) {
	off := m.Sub(spinWidgetCenter).Scale(1 / spinWidgetR)
	if l := off.Len(); l > 1 {
		off = off.Scale(1 / l)
	}
	g.spin = off
}

func cursorVec() Vec2 {
	mx, my := ebiten.CursorPosition()
	return Vec2{float64(mx), float64(my)}
}

func (g *Game) updateShooting() {
	g.shotPocketed = append(g.shotPocketed, step(g.balls)...)
	if allStopped(g.balls) {
		g.resolveTurn()
	}
}

func (g *Game) updateBallInHand() {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}
	pos := cursorVec()
	if !g.validCuePlacement(pos) {
		return
	}
	g.cue.Pos = pos
	g.cue.Vel = Vec2{}
	g.cue.Spin = Vec2{}
	g.cue.SideSpin = 0
	g.cue.Angle = 0
	g.cue.Active = true
	g.state = stateAiming
	g.message = fmt.Sprintf("Player %d's turn (%s)", g.current+1, g.players[g.current].Group)
}

// validCuePlacement reports whether the cue ball may be placed at p: inside the
// cushions and clear of every other active ball.
func (g *Game) validCuePlacement(p Vec2) bool {
	if p.X < playLeft+ballRadius || p.X > playRight-ballRadius ||
		p.Y < playTop+ballRadius || p.Y > playBottom-ballRadius {
		return false
	}
	for _, b := range g.balls {
		if b == g.cue || !b.Active {
			continue
		}
		if p.Sub(b.Pos).Len() < 2*ballRadius {
			return false
		}
	}
	return true
}

func (g *Game) switchTurn() {
	g.current = 1 - g.current
	g.state = stateAiming
	g.message = fmt.Sprintf("Player %d's turn (%s)", g.current+1, g.players[g.current].Group)
}

func (g *Game) foulBallInHand() {
	g.current = 1 - g.current
	g.state = stateBallInHand
	g.message = fmt.Sprintf("Scratch! Player %d: click to place the cue ball", g.current+1)
}

func (g *Game) endGame(winner int) {
	g.state = stateGameOver
	g.message = fmt.Sprintf("Player %d wins!  Press R to play again", winner+1)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenW, screenH
}
