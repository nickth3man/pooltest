package main

import "fmt"

// Group is a player's assigned ball set in 8-ball.
type Group int

const (
	GroupNone Group = iota
	GroupSolids
	GroupStripes
)

func (grp Group) String() string {
	switch grp {
	case GroupSolids:
		return "Solids"
	case GroupStripes:
		return "Stripes"
	default:
		return "Open"
	}
}

// groupOf returns the group a ball belongs to. The cue ball (0) and the
// 8-ball belong to no group.
func groupOf(number int) Group {
	switch {
	case number >= 1 && number <= 7:
		return GroupSolids
	case number >= 9 && number <= 15:
		return GroupStripes
	default:
		return GroupNone
	}
}

func otherGroup(grp Group) Group {
	if grp == GroupSolids {
		return GroupStripes
	}
	return GroupSolids
}

// groupRemaining counts the still-active balls of the given group.
func (g *Game) groupRemaining(grp Group) int {
	n := 0
	for _, b := range g.balls {
		if b.Active && groupOf(b.Number) == grp {
			n++
		}
	}
	return n
}

// resolveTurn applies 8-ball rules once every ball has come to rest, using the
// balls pocketed during the shot (g.shotPocketed; the cue ball is reported
// as 0).
//
// Simplifications vs. tournament 8-ball: shots are not "called", the only foul
// is scratching the cue ball, and clearing your group then sinking the 8 on the
// same stroke still counts as a win.
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
		cleared := cur != GroupNone && g.groupRemaining(cur) == 0
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
			if grp := groupOf(n); grp != GroupNone {
				g.players[g.current].Group = grp
				g.players[1-g.current].Group = otherGroup(grp)
				g.tableOpen = false
				legal = true
				break
			}
		}
	} else {
		cur := g.players[g.current].Group
		for _, n := range g.shotPocketed {
			if groupOf(n) == cur {
				legal = true
				break
			}
		}
	}

	if legal {
		g.message = fmt.Sprintf("Player %d continues (%s)", g.current+1, g.players[g.current].Group)
		g.state = stateAiming
	} else {
		g.switchTurn()
	}
}
