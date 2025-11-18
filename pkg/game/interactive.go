package game

import (
	"fmt"
	"time"

	"github.com/guny524/majhong_cli/pkg/mahjong/engine"
	"github.com/guny524/majhong_cli/pkg/mahjong/rules"
	"github.com/guny524/majhong_cli/pkg/ui/input"
	"github.com/guny524/majhong_cli/pkg/ui/terminal"
)

// InteractiveGame manages the interactive game loop
type InteractiveGame struct {
	game    *engine.GameState
	display *terminal.Display
	parser  *input.Parser
}

// NewInteractiveGame creates a new interactive game
func NewInteractiveGame(seed int64) *InteractiveGame {
	config := engine.DefaultRuleConfig()
	game := engine.NewGame(config, seed)

	return &InteractiveGame{
		game:    game,
		display: terminal.NewDisplay(),
		parser:  input.NewParser(),
	}
}

// Run runs the interactive game loop
func (ig *InteractiveGame) Run() error {
	// Initial display
	ig.display.RenderGameState(ig.game, engine.SeatEast)

	// Main game loop
	for !ig.game.IsGameOver {
		// Get current player
		currentPlayer := ig.game.CurrentPlayer

		// For AI players (South, West, North), simulate their turn
		if currentPlayer != engine.SeatEast {
			if err := ig.simulateAITurn(currentPlayer); err != nil {
				return err
			}
			continue
		}

		// Human player turn (East)
		if err := ig.handleHumanTurn(); err != nil {
			if err.Error() == "quit" {
				fmt.Println("\nGame ended by user.")
				return nil
			}
			// Display error and continue
			ig.display.PrintError(err)
			continue
		}

		// Render updated state
		ig.display.RenderGameState(ig.game, engine.SeatEast)
	}

	// Game over
	fmt.Println("\n" + "Game Over!")
	ig.display.RenderGameState(ig.game, engine.SeatEast)
	return nil
}

// handleHumanTurn handles the human player's turn
func (ig *InteractiveGame) handleHumanTurn() error {
	// Read input
	input, err := ig.parser.ReadInput()
	if err != nil {
		return err
	}

	// Parse action
	action, err := ig.parser.ParseAction(input, ig.game)
	if err != nil {
		return err
	}

	// Validate action
	if err := rules.ValidateAction(ig.game, action); err != nil {
		return fmt.Errorf("invalid action: %v", err)
	}

	// Apply action
	if err := ig.game.ApplyAction(action); err != nil {
		return fmt.Errorf("failed to apply action: %v", err)
	}

	return nil
}

// simulateAITurn simulates an AI player's turn
func (ig *InteractiveGame) simulateAITurn(seat engine.Seat) error {
	player := ig.game.Players[seat]

	// Brief pause for readability
	time.Sleep(500 * time.Millisecond)

	// Simple AI: draw and discard random tile
	if player.DrawnTile == nil {
		// Draw tile
		action := engine.NewAction(engine.ActionDraw, seat, nil, ig.game.TurnNumber)
		if err := ig.game.ApplyAction(action); err != nil {
			return err
		}
		fmt.Printf("\n%s draws a tile.\n", seat)
		time.Sleep(300 * time.Millisecond)
	}

	// Discard first tile (simple strategy)
	if player.DrawnTile != nil {
		// Discard the drawn tile
		action := engine.NewAction(engine.ActionDiscard, seat, []engine.Tile{*player.DrawnTile}, ig.game.TurnNumber)
		if err := ig.game.ApplyAction(action); err != nil {
			return err
		}
		fmt.Printf("%s discards %s.\n", seat, player.Discards[len(player.Discards)-1])
	} else if len(player.Hand) > 0 {
		// Fallback: discard first tile in hand
		tileToDiscard := player.Hand[0]
		action := engine.NewAction(engine.ActionDiscard, seat, []engine.Tile{tileToDiscard}, ig.game.TurnNumber)
		if err := ig.game.ApplyAction(action); err != nil {
			return err
		}
		fmt.Printf("%s discards %s.\n", seat, player.Discards[len(player.Discards)-1])
	}

	// Render updated state
	ig.display.RenderGameState(ig.game, engine.SeatEast)

	return nil
}

// GetGameState returns the current game state
func (ig *InteractiveGame) GetGameState() *engine.GameState {
	return ig.game
}
