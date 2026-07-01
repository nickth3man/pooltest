// Package rules implements 8-ball rules: which group a ball belongs to and
// how a turn resolves once every ball has come to rest.
package rules

import (
	"fmt"

	"github.com/user/pooltest/pool-go/internal/ball"
)

// Group is a player's assigned ball set in 8-ball.
type Group int

const (
	// GroupNone means the player has not yet been assigned a ball set (open
	// table, the cue ball, or the 8-ball).
	GroupNone Group = iota
	// GroupSolids is balls 1-7.
	GroupSolids
	// GroupStripes is balls 9-15.
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

// Of returns the group a ball belongs to. The cue ball (0) and the 8-ball
// belong to no group.
func Of(number int) Group {
	switch {
	case number >= 1 && number <= 7:
		return GroupSolids
	case number >= 9 && number <= 15:
		return GroupStripes
	default:
		return GroupNone
	}
}

// Other returns the opposing group.
func Other(grp Group) Group {
	if grp == GroupSolids {
		return GroupStripes
	}
	return GroupSolids
}

// GroupRemaining counts the still-active balls of the given group.
func GroupRemaining(balls []*ball.Ball, grp Group) int {
	n := 0
	for _, b := range balls {
		if b.Active && Of(b.Number) == grp {
			n++
		}
	}
	return n
}

// MessageForLegal returns the "Player N continues" message for a player who
// legally potted at least one of their own balls.
func MessageForLegal(player int, grp Group) string {
	return fmt.Sprintf("Player %d continues (%s)", player+1, grp)
}
