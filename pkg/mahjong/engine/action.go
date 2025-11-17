package engine

import (
	"fmt"
	"time"
)

// ActionType represents the type of game action
type ActionType int

const (
	ActionDraw ActionType = iota
	ActionDiscard
	ActionPon
	ActionChi
	ActionKan
	ActionRiichi
	ActionTsumo
	ActionRon
)

func (a ActionType) String() string {
	types := []string{"Draw", "Discard", "Pon", "Chi", "Kan", "Riichi", "Tsumo", "Ron"}
	return types[a]
}

// Action represents a single game action
type Action struct {
	Type      ActionType
	Player    Seat
	Tiles     []Tile
	Timestamp time.Time
	Turn      int
}

// NewAction creates a new action
func NewAction(actionType ActionType, player Seat, tiles []Tile, turn int) Action {
	return Action{
		Type:      actionType,
		Player:    player,
		Tiles:     tiles,
		Timestamp: time.Now(),
		Turn:      turn,
	}
}

// String returns action description
func (a Action) String() string {
	return fmt.Sprintf("Turn %d: %s %s %v", a.Turn, a.Player, a.Type, a.Tiles)
}
