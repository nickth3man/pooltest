package physics

import (
	"math"
	"testing"

	"github.com/user/pooltest/pool-go/internal/ball"
	"github.com/user/pooltest/pool-go/internal/table"
	"github.com/user/pooltest/pool-go/internal/vec"
)

const sitPeakSpin = 0.6 // rad/s scale for ~50% english tests

func clampDot(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}

func pxVel(v vec.Vec2) vec.Vec3 {
	return vec.Vec3{
		X: ball.PxToMeters(v.X * 60),
		Y: ball.PxToMeters(v.Y * 60),
	}
}

func mkBall(num int, pos vec.Vec2, vel vec.Vec2, wz float64) *ball.Ball {
	b := &ball.Ball{Number: num, Active: true, Params: ball.DefaultParams}
	b.SetPosPx(pos)
	b.RVW.V = pxVel(vel)
	b.RVW.W.Z = wz
	if vel.LenSq() > 0 || wz != 0 {
		b.Motion = ball.MotionSliding
	}
	return b
}

func angleBetween(a, b vec.Vec2) float64 {
	la, lb := a.Len(), b.Len()
	if la == 0 || lb == 0 {
		return 0
	}
	return math.Acos(clampDot(a.Dot(b) / (la * lb)))
}

func TestHeadOnCollisionTransfersMomentum(t *testing.T) {
	a := mkBall(0, vec.Vec2{X: 100, Y: 100}, vec.Vec2{X: 1, Y: 0}, 0)
	b := mkBall(1, vec.Vec2{X: 100 + 2*ball.CollisionRadius() - 0.5, Y: 100}, vec.Vec2{}, 0)
	ResolveBallCollisions([]*ball.Ball{a, b})
	a.SyncXY(1.0 / 60.0)
	b.SyncXY(1.0 / 60.0)
	if a.Vel.X > 0.1 {
		t.Errorf("striking ball kept too much speed: %.3f", a.Vel.X)
	}
	if b.Vel.X < 0.5 {
		t.Errorf("struck ball gained too little speed: %.3f", b.Vel.X)
	}
}

func TestHeadOnCollisionHasNoThrow(t *testing.T) {
	a := mkBall(0, vec.Vec2{X: 100, Y: 100}, vec.Vec2{X: 10, Y: 0}, 0)
	b := mkBall(1, vec.Vec2{X: 100 + 2*ball.CollisionRadius() - 0.5, Y: 100}, vec.Vec2{}, 0)
	ResolveBallCollisions([]*ball.Ball{a, b})
	b.SyncXY(1.0 / 60.0)
	if math.Abs(b.Vel.Y) > 0.05 {
		t.Errorf("head-on collision threw object ball sideways: vy = %.4f", b.Vel.Y)
	}
}

func TestCutShotThrowsObjectBall(t *testing.T) {
	const cutAngle = math.Pi / 6
	ax, ay := 100.0, 100.0
	overlap := 2*ball.CollisionRadius() - 0.5
	bx := ax + overlap*math.Cos(cutAngle)
	by := ay + overlap*math.Sin(cutAngle)
	a := mkBall(0, vec.Vec2{X: ax, Y: ay}, vec.Vec2{X: 10, Y: 0}, 0)
	b := mkBall(1, vec.Vec2{X: bx, Y: by}, vec.Vec2{}, 0)
	ResolveBallCollisions([]*ball.Ball{a, b})
	b.SyncXY(1.0 / 60.0)
	n := b.Pos.Sub(a.Pos).Normalize()
	throwDeg := angleBetween(b.Vel, n) * 180 / math.Pi
	if throwDeg <= 0.01 {
		t.Errorf("cut shot had no throw: %.3f°", throwDeg)
	}
	if throwDeg > 8 {
		t.Errorf("cut throw too large: %.3f° (want < 8°)", throwDeg)
	}
}

func TestSideSpinCausesSpinInducedThrow(t *testing.T) {
	a := mkBall(0, vec.Vec2{X: 100, Y: 100}, vec.Vec2{X: 10, Y: 0}, sitPeakSpin)
	b := mkBall(1, vec.Vec2{X: 100 + 2*ball.CollisionRadius() - 0.5, Y: 100}, vec.Vec2{}, 0)
	ResolveBallCollisions([]*ball.Ball{a, b})
	b.SyncXY(1.0 / 60.0)
	if math.Abs(b.Vel.Y) < 0.02 {
		t.Errorf("sidespin did not throw object ball: vy = %.4f", b.Vel.Y)
	}
}

func TestRailReflectsVelocity(t *testing.T) {
	b := mkBall(0, vec.Vec2{X: table.PlayRight - ball.CollisionRadius() - 0.1, Y: 300}, vec.Vec2{X: 5, Y: 0}, 0)
	Step([]*ball.Ball{b})
	if b.Vel.X >= 0 {
		t.Errorf("ball did not bounce off the right rail: vx = %.3f", b.Vel.X)
	}
}

func TestPocketCapture(t *testing.T) {
	p := table.PocketCenters[0]
	b := mkBall(3, p, vec.Vec2{}, 0)
	got := Step([]*ball.Ball{b})
	if b.Active {
		t.Error("ball over a pocket was not captured")
	}
	if len(got) != 1 || got[0] != 3 {
		t.Errorf("step returned %v, want [3]", got)
	}
}

func TestFrictionBringsBallToRest(t *testing.T) {
	b := mkBall(0, vec.Vec2{X: 400, Y: 300}, vec.Vec2{X: 6, Y: 0}, 0)
	balls := []*ball.Ball{b}
	for range 3000 {
		if AllStopped(balls) {
			break
		}
		Step(balls)
	}
	if b.Moving() {
		t.Errorf("ball never came to rest")
	}
}

func TestDrawPathAfterContact(t *testing.T) {
	b := mkBall(0, vec.Vec2{X: 400, Y: 300}, vec.Vec2{X: 1.5, Y: 0}, 0)
	b.RVW.W.Y = -50 // draw spin (backspin along +x shot)
	b.Motion = ball.MotionSliding
	startX := b.Pos.X
	balls := []*ball.Ball{b}
	for range 4000 {
		if AllStopped(balls) {
			break
		}
		Step(balls)
	}
	if b.Pos.X >= startX-2 {
		t.Errorf("backspin did not draw the ball back: x = %.3f (want < %.3f)", b.Pos.X, startX)
	}
}

func TestJawPegRejectsBall(t *testing.T) {
	peg := JawPegs[0]
	b := mkBall(1, vec.Vec2{X: peg.X, Y: peg.Y + ball.CollisionRadius() - 1}, vec.Vec2{X: 0, Y: -3}, 0)
	if s := ReflectOffPeg(b, peg); s <= 0 {
		t.Fatalf("expected a jaw collision, got impact %.3f", s)
	}
	b.SyncXY(1.0 / 60.0)
	if b.Vel.Y <= 0 {
		t.Errorf("ball did not bounce off the jaw: vy = %.3f", b.Vel.Y)
	}
}

func TestMouthPredicates(t *testing.T) {
	midX := float64(table.PlayLeft+table.PlayRight) / 2
	if !InTopBottomMouth(midX) {
		t.Error("side pocket mouth not detected at midX")
	}
	if InTopBottomMouth(midX + SideMouth + 10) {
		t.Error("rail wrongly reported open away from any pocket")
	}
}

func cutThrowDeg(speed float64, wz float64, wySpin float64) float64 {
	const cutAngle = math.Pi / 6
	ax, ay := 100.0, 100.0
	overlap := 2*ball.CollisionRadius() - 0.5
	bx := ax + overlap*math.Cos(cutAngle)
	by := ay + overlap*math.Sin(cutAngle)
	a := mkBall(0, vec.Vec2{X: ax, Y: ay}, vec.Vec2{X: speed, Y: 0}, wz)
	a.RVW.W.Y = wySpin
	b := mkBall(1, vec.Vec2{X: bx, Y: by}, vec.Vec2{}, 0)
	ResolveBallCollisions([]*ball.Ball{a, b})
	b.SyncXY(1.0 / 60.0)
	n := b.Pos.Sub(a.Pos).Normalize()
	return angleBetween(b.Vel, n) * 180 / math.Pi
}

func TestSlowCutThrowsMoreThanFast(t *testing.T) {
	slow := cutThrowDeg(4, 0, 0)
	fast := cutThrowDeg(14, 0, 0)
	if slow <= fast {
		t.Errorf("slow cut throw %.3f° should exceed fast cut %.3f°", slow, fast)
	}
}

func TestDrawReducesThrow(t *testing.T) {
	stun := cutThrowDeg(5, 0, 0)
	draw := cutThrowDeg(5, 0, -15)
	if draw >= stun*0.75 {
		t.Errorf("draw throw %.3f° should be well below stun throw %.3f°", draw, stun)
	}
}

func TestSITSaturatesAtHighSideSpin(t *testing.T) {
	moderate := cutThrowDeg(5, sitPeakSpin*0.5, 0)
	peak := cutThrowDeg(5, sitPeakSpin, 0)
	if peak <= moderate {
		t.Errorf("peak sidespin throw %.3f° should exceed moderate %.3f°", peak, moderate)
	}
}

func TestGearingOEReducesThrow(t *testing.T) {
	// TP A.26: outside gearing at ~0.8×(1−hit_fraction) for 30° cut.
	const hitFraction = 0.5 // half-ball cut
	gearing := -sitPeakSpin * 0.8 * (1 - hitFraction)
	stun := cutThrowDeg(5, 0, 0)
	geared := cutThrowDeg(5, gearing, 0)
	if math.Abs(geared) >= math.Abs(stun) {
		t.Errorf("gearing OE throw %.3f° should be less than stun %.3f°", geared, stun)
	}
}

func TestPocketClearsSpin(t *testing.T) {
	p := table.PocketCenters[0]
	b := mkBall(3, p, vec.Vec2{}, 1.5)
	Step([]*ball.Ball{b})
	if b.RVW.W.Z != 0 {
		t.Errorf("pocketed ball spin = %.3f, want 0", b.RVW.W.Z)
	}
}

func TestMovingIncludesSpin(t *testing.T) {
	b := &ball.Ball{Number: 0, Active: true, Params: ball.DefaultParams}
	b.RVW.W.Z = 0.05
	b.Motion = ball.MotionSpinning
	if !b.Moving() {
		t.Error("ball with spin should report Moving")
	}
	b.RVW.W.Z = 0
	b.Motion = ball.MotionStationary
	if b.Moving() {
		t.Error("ball with no motion should not report Moving")
	}
}
