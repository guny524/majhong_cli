package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/guny524/majhong_cli/pkg/mahjong/engine"
	"github.com/guny524/majhong_cli/pkg/mahjong/serialization"
)

// Exit codes
const (
	ExitSuccess       = 0
	ExitInvalidAction = 1
	ExitFileIOError   = 2
	ExitCorruptedFile = 3
	ExitTamperedWall  = 4
)

// BatchMode handles batch processing mode
type BatchMode struct {
	GameFile      string
	Action        string
	Actions       string
	Init          bool
	Query         string
	Seed          int64
	JSON          bool
	ExitOnTamper  bool
}

// Run executes the batch mode operation
func (b *BatchMode) Run() int {
	// Handle --init mode
	if b.Init {
		return b.runInit()
	}

	// Handle --query mode
	if b.Query != "" {
		return b.runQuery()
	}

	// Handle --action or --actions mode
	if b.Action != "" || b.Actions != "" {
		return b.runAction()
	}

	fmt.Fprintf(os.Stderr, "Error: must specify --init, --query, or --action\n")
	return ExitInvalidAction
}

// runInit initializes a new game and saves to file
func (b *BatchMode) runInit() int {
	if b.GameFile == "" {
		fmt.Fprintf(os.Stderr, "Error: --game-file required for --init\n")
		return ExitInvalidAction
	}

	// Create new game with seed
	seed := b.Seed
	if seed == 0 {
		seed = 1234567890 // Default seed
	}

	config := engine.DefaultRuleConfig()
	game := engine.NewGame(config, seed)

	// Save to file
	if err := serialization.SaveToFile(game, b.GameFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to save game file: %v\n", err)
		return ExitFileIOError
	}

	if b.JSON {
		b.printJSON(map[string]interface{}{
			"status": "initialized",
			"seed":   seed,
			"file":   b.GameFile,
		})
	} else {
		fmt.Printf("Game initialized with seed %d and saved to %s\n", seed, b.GameFile)
	}

	return ExitSuccess
}

// runQuery queries game state information
func (b *BatchMode) runQuery() int {
	if b.GameFile == "" {
		fmt.Fprintf(os.Stderr, "Error: --game-file required for --query\n")
		return ExitInvalidAction
	}

	// Load game state
	game, exitCode := b.loadGameState()
	if game == nil {
		return exitCode
	}

	// Execute query
	switch b.Query {
	case "current-player":
		if b.JSON {
			b.printJSON(map[string]interface{}{
				"currentPlayer": game.CurrentPlayer.String(),
			})
		} else {
			fmt.Println(game.CurrentPlayer.String())
		}

	case "valid-actions":
		actions := game.ValidActions()
		if b.JSON {
			actionStrs := make([]string, len(actions))
			for i, action := range actions {
				actionStrs[i] = formatAction(action)
			}
			b.printJSON(map[string]interface{}{
				"validActions": actionStrs,
			})
		} else {
			for _, action := range actions {
				fmt.Println(formatAction(action))
			}
		}

	case "scores":
		scores := make(map[string]int)
		for _, player := range game.Players {
			scores[player.Seat.String()] = player.Score
		}
		if b.JSON {
			b.printJSON(map[string]interface{}{
				"scores": scores,
			})
		} else {
			for seat, score := range scores {
				fmt.Printf("%s: %d\n", seat, score)
			}
		}

	case "round":
		if b.JSON {
			b.printJSON(map[string]interface{}{
				"dealer":      game.Round.Dealer.String(),
				"roundWind":   game.Round.RoundWind.String(),
				"roundNumber": game.Round.RoundNumber,
				"honba":       game.Round.Honba,
			})
		} else {
			fmt.Printf("Round: %s %d, Dealer: %s, Honba: %d\n",
				game.Round.RoundWind.String(),
				game.Round.RoundNumber,
				game.Round.Dealer.String(),
				game.Round.Honba)
		}

	case "wall-remaining":
		if b.JSON {
			b.printJSON(map[string]interface{}{
				"remaining": len(game.Wall.LiveWall),
			})
		} else {
			fmt.Printf("%d\n", len(game.Wall.LiveWall))
		}

	default:
		fmt.Fprintf(os.Stderr, "Error: unknown query type: %s\n", b.Query)
		fmt.Fprintf(os.Stderr, "Valid queries: current-player, valid-actions, scores, round, wall-remaining\n")
		return ExitInvalidAction
	}

	return ExitSuccess
}

// runAction executes one or more actions and updates game state
func (b *BatchMode) runAction() int {
	if b.GameFile == "" {
		fmt.Fprintf(os.Stderr, "Error: --game-file required for --action\n")
		return ExitInvalidAction
	}

	// Load game state
	game, exitCode := b.loadGameState()
	if game == nil {
		return exitCode
	}

	// Parse actions
	var actionStrs []string
	if b.Actions != "" {
		// Parse JSON array of actions
		if err := json.Unmarshal([]byte(b.Actions), &actionStrs); err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid actions JSON: %v\n", err)
			return ExitInvalidAction
		}
	} else {
		actionStrs = []string{b.Action}
	}

	// Execute actions
	for i, actionStr := range actionStrs {
		action, err := parseAction(actionStr, game)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: invalid action %d '%s': %v\n", i+1, actionStr, err)
			return ExitInvalidAction
		}

		if err := game.ApplyAction(action); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to apply action %d '%s': %v\n", i+1, actionStr, err)
			return ExitInvalidAction
		}

		if b.JSON {
			b.printJSON(map[string]interface{}{
				"status":      "action_applied",
				"action":      formatAction(action),
				"actionIndex": i + 1,
			})
		} else {
			fmt.Printf("Applied: %s\n", formatAction(action))
		}
	}

	// Save updated game state
	if err := serialization.SaveToFile(game, b.GameFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to save game file: %v\n", err)
		return ExitFileIOError
	}

	if b.JSON {
		b.printJSON(map[string]interface{}{
			"status":        "success",
			"actionsApplied": len(actionStrs),
			"currentPlayer": game.CurrentPlayer.String(),
			"turnNumber":    game.TurnNumber,
		})
	} else {
		fmt.Printf("Game state updated successfully\n")
	}

	return ExitSuccess
}

// loadGameState loads game state from file with error handling
func (b *BatchMode) loadGameState() (*engine.GameState, int) {
	game, err := serialization.LoadFromFile(b.GameFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: game file not found: %s\n", b.GameFile)
			return nil, ExitFileIOError
		}

		// Check if it's a wall hash mismatch (tampering)
		if strings.Contains(err.Error(), "wall hash mismatch") {
			fmt.Fprintf(os.Stderr, "Error: wall tampering detected: %v\n", err)
			if b.ExitOnTamper {
				return nil, ExitTamperedWall
			}
		}

		// Check if it's a corrupted file
		if strings.Contains(err.Error(), "unmarshal") || strings.Contains(err.Error(), "invalid") {
			fmt.Fprintf(os.Stderr, "Error: corrupted game file: %v\n", err)
			return nil, ExitCorruptedFile
		}

		fmt.Fprintf(os.Stderr, "Error: failed to load game file: %v\n", err)
		return nil, ExitFileIOError
	}

	return game, ExitSuccess
}

// parseAction parses action string into Action struct
func parseAction(actionStr string, game *engine.GameState) (engine.Action, error) {
	parts := strings.Fields(actionStr)
	if len(parts) == 0 {
		return engine.Action{}, fmt.Errorf("empty action")
	}

	actionType := strings.ToLower(parts[0])
	currentPlayer := game.CurrentPlayer

	switch actionType {
	case "draw":
		return engine.NewAction(engine.ActionDraw, currentPlayer, nil, game.TurnNumber), nil

	case "discard":
		if len(parts) != 2 {
			return engine.Action{}, fmt.Errorf("discard requires tile argument (e.g., 'discard 3m')")
		}
		tile, err := engine.ParseTile(parts[1])
		if err != nil {
			return engine.Action{}, fmt.Errorf("invalid tile: %w", err)
		}
		return engine.NewAction(engine.ActionDiscard, currentPlayer, []engine.Tile{tile}, game.TurnNumber), nil

	case "pon":
		if len(parts) != 2 {
			return engine.Action{}, fmt.Errorf("pon requires tile argument (e.g., 'pon 3m')")
		}
		tile, err := engine.ParseTile(parts[1])
		if err != nil {
			return engine.Action{}, fmt.Errorf("invalid tile: %w", err)
		}
		// Pon uses 3 of the same tile
		tiles := []engine.Tile{tile, tile, tile}
		return engine.NewAction(engine.ActionPon, currentPlayer, tiles, game.TurnNumber), nil

	case "chi":
		if len(parts) != 2 {
			return engine.Action{}, fmt.Errorf("chi requires sequence argument (e.g., 'chi 345m')")
		}
		// Parse sequence like "345m" or "3m4m5m"
		tiles, err := parseSequence(parts[1])
		if err != nil {
			return engine.Action{}, fmt.Errorf("invalid sequence: %w", err)
		}
		if len(tiles) != 3 {
			return engine.Action{}, fmt.Errorf("chi requires exactly 3 tiles")
		}
		return engine.NewAction(engine.ActionChi, currentPlayer, tiles, game.TurnNumber), nil

	case "kan":
		if len(parts) != 2 {
			return engine.Action{}, fmt.Errorf("kan requires tile argument (e.g., 'kan 3m')")
		}
		tile, err := engine.ParseTile(parts[1])
		if err != nil {
			return engine.Action{}, fmt.Errorf("invalid tile: %w", err)
		}
		// Kan uses 4 of the same tile
		tiles := []engine.Tile{tile, tile, tile, tile}
		return engine.NewAction(engine.ActionKan, currentPlayer, tiles, game.TurnNumber), nil

	case "riichi":
		return engine.NewAction(engine.ActionRiichi, currentPlayer, nil, game.TurnNumber), nil

	case "tsumo":
		return engine.NewAction(engine.ActionTsumo, currentPlayer, nil, game.TurnNumber), nil

	case "ron":
		return engine.NewAction(engine.ActionRon, currentPlayer, nil, game.TurnNumber), nil

	default:
		return engine.Action{}, fmt.Errorf("unknown action type: %s", actionType)
	}
}

// parseSequence parses tile sequence like "345m" or "3m4m5m"
func parseSequence(seq string) ([]engine.Tile, error) {
	// Try parsing as compact notation first (e.g., "345m")
	if len(seq) >= 4 {
		suit := seq[len(seq)-1:]
		numbers := seq[:len(seq)-1]

		// Check if all characters before suit are digits
		allDigits := true
		for _, c := range numbers {
			if c < '0' || c > '9' {
				allDigits = false
				break
			}
		}

		if allDigits && len(numbers) == 3 {
			tiles := make([]engine.Tile, 3)
			for i, c := range numbers {
				tileStr := string(c) + suit
				tile, err := engine.ParseTile(tileStr)
				if err != nil {
					return nil, err
				}
				tiles[i] = tile
			}
			return tiles, nil
		}
	}

	// Try parsing as full notation (e.g., "3m4m5m")
	parts := make([]string, 0)
	current := ""
	for _, c := range seq {
		current += string(c)
		if c == 'm' || c == 'p' || c == 's' || c == 'z' {
			parts = append(parts, current)
			current = ""
		}
	}

	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid sequence format")
	}

	return engine.ParseTiles(parts)
}

// formatAction formats an action for display
func formatAction(action engine.Action) string {
	switch action.Type {
	case engine.ActionDraw:
		return "draw"
	case engine.ActionDiscard:
		if len(action.Tiles) > 0 {
			return fmt.Sprintf("discard %s", action.Tiles[0].String())
		}
		return "discard"
	case engine.ActionPon:
		if len(action.Tiles) > 0 {
			return fmt.Sprintf("pon %s", action.Tiles[0].String())
		}
		return "pon"
	case engine.ActionChi:
		if len(action.Tiles) > 0 {
			return fmt.Sprintf("chi %s%s%s",
				action.Tiles[0].String(),
				action.Tiles[1].String(),
				action.Tiles[2].String())
		}
		return "chi"
	case engine.ActionKan:
		if len(action.Tiles) > 0 {
			return fmt.Sprintf("kan %s", action.Tiles[0].String())
		}
		return "kan"
	case engine.ActionRiichi:
		return "riichi"
	case engine.ActionTsumo:
		return "tsumo"
	case engine.ActionRon:
		return "ron"
	default:
		return "unknown"
	}
}

// printJSON prints JSON output
func (b *BatchMode) printJSON(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to marshal JSON: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}
