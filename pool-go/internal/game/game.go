// Package game is the orchestrator. It owns the Game struct (the
// ebiten.Game implementation), the state machine, the input handling, and the
// rendering of the table, balls, HUD, and aim helpers.
package game

import (
	"bytes"
	"fmt"
	"log/slog"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/user/pooltest/pool-go/internal/audio"
	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/fx"
	"github.com/user/pooltest/pool-go/internal/physics"
	"github.com/user/pooltest/pool-go/internal/rack"
	"github.com/user/pooltest/pool-go/internal/rules"
	"github.com/user/pooltest/pool-go/internal/sprites"
	"github.com/user/pooltest/pool-go/internal/table"
	"github.com/user/pooltest/pool-go/internal/vec"
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

const (
	spinWidgetR  = 30.0
	spinWidgetCx = 56
	spinWidgetCy = table.ScreenH - 56
)

var spinWidgetCenter = vec.Vec2{X: float64(spinWidgetCx), Y: float64(spinWidgetCy)}

type player struct {
	Group rules.Group
}

// Game holds the entire match state and implements ebiten.Game.
type Game struct {
	balls   []*ball.Ball
	cue     *ball.Ball
	players [2]player
	current int

	state        gameState
	tableOpen    bool
	charging     bool
	shotPocketed []int

	spin vec.Vec2 // pending english from the dial: x = side, y = follow(-)/draw(+)

	face      text.Face
	smallFace text.Face

	message string

	// Presentation / feedback state.
	canvas *ebiten.Image
	fx     fx.State
}

// NewGame loads the font and racks a fresh match.
func NewGame() *Game {
	src, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		slog.Error("failed to load embedded font", "err", err)
		os.Exit(1)
	}
	g := &Game{
		face:      &text.GoTextFace{Source: src, Size: 16},
		smallFace: &text.GoTextFace{Source: src, Size: 11},
	}
	audio.Init()
	// Route physics events into feedback (particles, shake, sound).
	physics.OnBallImpact = func(pos vec.Vec2, strength float64) {
		fx.SpawnImpact(&g.fx, pos, strength)
		audio.PlayClick(strength)
	}
	physics.OnRailImpact = func(_ vec.Vec2, strength float64) {
		audio.PlayRail(strength)
		if strength > 6 {
			g.fx.Shake = math.Min(4, g.fx.Shake+1)
		}
	}
	physics.OnPocketDrop = func(pos vec.Vec2) {
		fx.SpawnPuff(&g.fx, pos)
		audio.PlayPocket()
	}
	g.reset()
	return g
}

func (g *Game) reset() {
	g.balls = rack.NewRack()
	g.cue = g.balls[0]
	g.players = [2]player{}
	g.current = 0
	g.state = stateAiming
	g.tableOpen = true
	g.charging = false
	g.shotPocketed = nil
	g.spin = vec.Vec2{}
	g.fx.Particles = nil
	g.fx.Shake = 0
	g.message = "Player 1 to break"
}

// Update advances the game one frame. It implements ebiten.Game.
func (g *Game) Update() error {
	fx.Tick(&g.fx, g.balls)
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

// Layout implements ebiten.Game.
func (g *Game) Layout(_, _ int) (int, int) {
	return table.ScreenW, table.ScreenH
}

// updateAiming implements a slingshot cue: hold the left button, drag back
// from the cue ball, and release to fire in the opposite direction with power
// proportional to the pull distance. The English dial in the corner intercepts
// clicks so the player can dial spin before taking the shot.
func (g *Game) updateAiming() {
	m := cursorVec()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		g.spin = vec.Vec2{}
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
func (g *Game) fire(pull vec.Vec2, power float64) {
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

	audio.PlayCueStrike(frac)
	g.spin = vec.Vec2{}
	g.shotPocketed = nil
	g.state = stateShooting
}

// inSpinWidget reports whether m lies within the English dial.
func (g *Game) inSpinWidget(m vec.Vec2) bool {
	return m.Sub(spinWidgetCenter).Len() <= spinWidgetR+4
}

// setSpinFromWidget records the dialed english as a unit offset (clamped to the
// dial face).
func (g *Game) setSpinFromWidget(m vec.Vec2) {
	off := m.Sub(spinWidgetCenter).Scale(1 / spinWidgetR)
	if l := off.Len(); l > 1 {
		off = off.Scale(1 / l)
	}
	g.spin = off
}

func cursorVec() vec.Vec2 {
	mx, my := ebiten.CursorPosition()
	return vec.Vec2{X: float64(mx), Y: float64(my)}
}

func (g *Game) updateShooting() {
	g.shotPocketed = append(g.shotPocketed, physics.Step(g.balls)...)
	if physics.AllStopped(g.balls) {
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
	g.cue.Vel = vec.Vec2{}
	g.cue.Spin = vec.Vec2{}
	g.cue.SideSpin = 0
	g.cue.Angle = 0
	g.cue.Active = true
	g.state = stateAiming
	g.message = fmt.Sprintf("Player %d's turn (%s)", g.current+1, g.players[g.current].Group)
}

// validCuePlacement reports whether the cue ball may be placed at p: inside the
// cushions and clear of every other active ball.
func (g *Game) validCuePlacement(p vec.Vec2) bool {
	if p.X < table.PlayLeft+ball.Radius || p.X > table.PlayRight-ball.Radius ||
		p.Y < table.PlayTop+ball.Radius || p.Y > table.PlayBottom-ball.Radius {
		return false
	}
	for _, b := range g.balls {
		if b == g.cue || !b.Active {
			continue
		}
		if p.Sub(b.Pos).Len() < 2*ball.Radius {
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

// resolveTurn applies 8-ball rules once every ball has come to rest, using
// the balls pocketed during the shot (g.shotPocketed; the cue ball is
// reported as 0).
//
// Simplifications vs. tournament 8-ball: shots are not "called", the only
// foul is scratching the cue ball, and clearing your group then sinking the 8
// on the same stroke still counts as a win.
func (g *Game) resolveTurn() {
	cueScratch := false
	eight := false
	for _, n := range g.shotPocketed {
		switch n {
		case 0:
			cueScratch = true
		case 8:
			eight = true
		}
	}

	if eight {
		cur := g.players[g.current].Group
		cleared := cur != rules.GroupNone && rules.GroupRemaining(g.balls, cur) == 0
		if cleared && !cueScratch {
			g.endGame(g.current)
		} else {
			g.endGame(1 - g.current)
		}
		return
	}

	if cueScratch {
		g.foulBallInHand()
		return
	}

	legal := false
	if g.tableOpen {
		// The first ball pocketed legally assigns groups to both players.
		for _, n := range g.shotPocketed {
			if grp := rules.Of(n); grp != rules.GroupNone {
				g.players[g.current].Group = grp
				g.players[1-g.current].Group = rules.Other(grp)
				g.tableOpen = false
				legal = true
				break
			}
		}
	} else {
		cur := g.players[g.current].Group
		for _, n := range g.shotPocketed {
			if rules.Of(n) == cur {
				legal = true
				break
			}
		}
	}

	if legal {
		g.message = rules.MessageForLegal(g.current, g.players[g.current].Group)
		g.state = stateAiming
	} else {
		g.switchTurn()
	}
}

// SpriteFor returns the cached shaded body sprite for a ball. Used by render
// and by tests that need to inspect sprite data.
func (g *Game) SpriteFor(b *ball.Ball) *ebiten.Image {
	return sprites.Ball(b, g.smallFace)
}
