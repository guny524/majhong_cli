package rules

import (
	"fmt"
	"github.com/guny524/majhong_cli/pkg/mahjong/engine"
)

// ValidateAction checks if an action is legal in the current game state
func ValidateAction(game *engine.GameState, action engine.Action) error {
	switch action.Type {
	case engine.ActionDraw:
		return validateDraw(game, action)
	case engine.ActionDiscard:
		return validateDiscard(game, action)
	case engine.ActionPon:
		return validatePon(game, action)
	case engine.ActionChi:
		return validateChi(game, action)
	case engine.ActionKan:
		return validateKan(game, action)
	case engine.ActionRiichi:
		return validateRiichi(game, action)
	case engine.ActionTsumo:
		return validateTsumo(game, action)
	case engine.ActionRon:
		return validateRon(game, action)
	default:
		return fmt.Errorf("unknown action type: %v", action.Type)
	}
}

func validateDraw(game *engine.GameState, action engine.Action) error {
	if game.Wall.IsExhausted() {
		return fmt.Errorf("wall exhausted, cannot draw")
	}

	player := game.Players[action.Player]
	if player.DrawnTile != nil {
		return fmt.Errorf("player already has drawn tile")
	}

	return nil
}

func validateDiscard(game *engine.GameState, action engine.Action) error {
	if len(action.Tiles) != 1 {
		return fmt.Errorf("must discard exactly 1 tile")
	}

	player := game.Players[action.Player]
	tile := action.Tiles[0]

	// Check if tile is in hand or is drawn tile
	if player.DrawnTile != nil && player.DrawnTile.Equals(tile) {
		return nil
	}

	for _, handTile := range player.Hand {
		if handTile.Equals(tile) {
			return nil
		}
	}

	return fmt.Errorf("tile %s not in hand", tile)
}

func validatePon(game *engine.GameState, action engine.Action) error {
	if len(action.Tiles) != 3 {
		return fmt.Errorf("pon must have exactly 3 tiles")
	}

	// All tiles must be same type
	for i := 1; i < 3; i++ {
		if !action.Tiles[i].SameType(action.Tiles[0]) {
			return fmt.Errorf("pon tiles must be same type")
		}
	}

	player := game.Players[action.Player]

	// Player must have 2 matching tiles in hand
	count := 0
	for _, handTile := range player.Hand {
		if handTile.SameType(action.Tiles[0]) {
			count++
		}
	}
	if count < 2 {
		return fmt.Errorf("insufficient tiles for pon")
	}

	return nil
}

func validateChi(game *engine.GameState, action engine.Action) error {
	if len(action.Tiles) != 3 {
		return fmt.Errorf("chi must have exactly 3 tiles")
	}

	// Can only chi from player to the left
	lastDiscard := getLastDiscard(game)
	if lastDiscard == nil {
		return fmt.Errorf("no tile to chi")
	}

	expectedPlayer := (lastDiscard.Player + 1) % 4
	if action.Player != expectedPlayer {
		return fmt.Errorf("can only chi from player to the left")
	}

	// Must be a sequence in same suit
	tiles := action.Tiles
	if tiles[0].Suit != tiles[1].Suit || tiles[1].Suit != tiles[2].Suit {
		return fmt.Errorf("chi must be same suit")
	}

	// Cannot chi honors
	if tiles[0].Suit == engine.Ji {
		return fmt.Errorf("cannot chi honor tiles")
	}

	// Check sequence (sorted)
	ranks := []int{tiles[0].Rank, tiles[1].Rank, tiles[2].Rank}
	// Sort ranks
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 3; j++ {
			if ranks[i] > ranks[j] {
				ranks[i], ranks[j] = ranks[j], ranks[i]
			}
		}
	}
	if ranks[1] != ranks[0]+1 || ranks[2] != ranks[1]+1 {
		return fmt.Errorf("chi must be consecutive sequence")
	}

	return nil
}

func validateKan(game *engine.GameState, action engine.Action) error {
	if len(action.Tiles) != 4 {
		return fmt.Errorf("kan must have exactly 4 tiles")
	}

	// All tiles must be same type
	for i := 1; i < 4; i++ {
		if !action.Tiles[i].SameType(action.Tiles[0]) {
			return fmt.Errorf("kan tiles must be same type")
		}
	}

	return nil
}

func validateRiichi(game *engine.GameState, action engine.Action) error {
	player := game.Players[action.Player]

	if !player.IsClosed() {
		return fmt.Errorf("cannot riichi with open hand")
	}

	if player.Score < 1000 {
		return fmt.Errorf("insufficient points for riichi")
	}

	if game.Wall.Remaining() < 4 {
		return fmt.Errorf("not enough tiles remaining for riichi")
	}

	// TODO: Validate tenpai status

	return nil
}

func validateTsumo(game *engine.GameState, action engine.Action) error {
	player := game.Players[action.Player]

	if player.DrawnTile == nil {
		return fmt.Errorf("no drawn tile to tsumo")
	}

	// TODO: Validate winning hand and yaku

	return nil
}

func validateRon(game *engine.GameState, action engine.Action) error {
	player := game.Players[action.Player]

	// Check furiten
	if player.Furiten.IsFuriten() {
		return fmt.Errorf("cannot ron in furiten")
	}

	// TODO: Validate winning hand and yaku

	return nil
}

func getLastDiscard(game *engine.GameState) *engine.Action {
	for i := len(game.History) - 1; i >= 0; i-- {
		if game.History[i].Type == engine.ActionDiscard {
			return &game.History[i]
		}
	}
	return nil
}
