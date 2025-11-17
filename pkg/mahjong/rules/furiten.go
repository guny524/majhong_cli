package rules

import (
	"github.com/guny524/majhong_cli/pkg/mahjong/engine"
)

// CheckDiscardFuriten checks if player is in discard furiten
// Discard furiten: if any wait tile is in player's discards, cannot ron
func CheckDiscardFuriten(waitTiles []string, discards []string) bool {
	discardMap := make(map[string]bool)
	for _, discard := range discards {
		discardMap[discard] = true
	}

	for _, wait := range waitTiles {
		if discardMap[wait] {
			return true // Any wait tile in discards = furiten
		}
	}

	return false
}

// CheckTemporaryFuriten checks if player is in temporary furiten
// Temporary furiten: if player skipped a ron opportunity, cannot ron until next turn
func CheckTemporaryFuriten(lastSkippedRonRound int, currentRound int, isRiichi bool) bool {
	// Temporary furiten doesn't apply during riichi (uses riichi furiten instead)
	if isRiichi {
		return false
	}

	// If skipped ron in current round, in temporary furiten
	return lastSkippedRonRound == currentRound
}

// CheckRiichiFuriten checks if player is in riichi furiten
// Riichi furiten: after riichi, if player skips ANY winning tile, permanent furiten
func CheckRiichiFuriten(isRiichi bool, skippedTilesAfterRiichi []string, waitTiles []string) bool {
	if !isRiichi {
		return false
	}

	// If player has skipped any tiles after riichi, in riichi furiten
	return len(skippedTilesAfterRiichi) > 0
}

// FuritenState represents the furiten status of a player
type FuritenState struct {
	DiscardFuriten   bool
	TemporaryFuriten bool
	RiichiFuriten    bool
}

// IsFuriten returns true if player is in ANY furiten state
func (f FuritenState) IsFuriten() bool {
	return f.DiscardFuriten || f.TemporaryFuriten || f.RiichiFuriten
}

// CanRon returns true if player can declare ron (not in furiten)
func (f FuritenState) CanRon() bool {
	return !f.IsFuriten()
}

// CanTsumo returns true if player can declare tsumo
// Furiten only affects ron, not tsumo
func (f FuritenState) CanTsumo() bool {
	return true
}

// CheckCompositeFuriten checks all three furiten types
func CheckCompositeFuriten(
	waitTiles []string,
	discards []string,
	lastSkippedRonRound int,
	currentRound int,
	isRiichi bool,
	skippedTilesAfterRiichi []string,
) FuritenState {
	return FuritenState{
		DiscardFuriten:   CheckDiscardFuriten(waitTiles, discards),
		TemporaryFuriten: CheckTemporaryFuriten(lastSkippedRonRound, currentRound, isRiichi),
		RiichiFuriten:    CheckRiichiFuriten(isRiichi, skippedTilesAfterRiichi, waitTiles),
	}
}

// DeclareRon attempts to declare ron (validates furiten)
func DeclareRon(furitenState FuritenState) error {
	if furitenState.IsFuriten() {
		return &engine.GameError{
			Code:    "FURITEN_VIOLATION",
			Message: "Cannot declare ron while in furiten",
		}
	}
	return nil
}

// DeclareTsumo attempts to declare tsumo (always allowed regardless of furiten)
func DeclareTsumo(furitenState FuritenState) error {
	// Tsumo is always allowed, even in furiten
	return nil
}
