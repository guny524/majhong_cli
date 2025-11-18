package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/guny524/majhong_cli/pkg/mahjong/engine"
	"github.com/guny524/majhong_cli/pkg/mahjong/serialization"
)

func TestBatchInit(t *testing.T) {
	tempDir := t.TempDir()
	gameFile := filepath.Join(tempDir, "test_game.json")

	batch := &BatchMode{
		GameFile: gameFile,
		Init:     true,
		Seed:     12345,
		JSON:     false,
	}

	exitCode := batch.runInit()
	if exitCode != ExitSuccess {
		t.Errorf("Expected exit code %d, got %d", ExitSuccess, exitCode)
	}

	// Verify file was created
	if _, err := os.Stat(gameFile); os.IsNotExist(err) {
		t.Fatal("Game file was not created")
	}

	// Load and verify
	game, err := serialization.LoadFromFile(gameFile)
	if err != nil {
		t.Fatalf("Failed to load game file: %v", err)
	}

	if game.Wall.Seed != 12345 {
		t.Errorf("Expected seed 12345, got %d", game.Wall.Seed)
	}
}

func TestBatchQuery(t *testing.T) {
	tempDir := t.TempDir()
	gameFile := filepath.Join(tempDir, "test_game.json")

	// Create and save a game
	config := engine.DefaultRuleConfig()
	game := engine.NewGame(config, 54321)
	if err := serialization.SaveToFile(game, gameFile); err != nil {
		t.Fatalf("Failed to save game: %v", err)
	}

	tests := []struct {
		query    string
		wantExit int
	}{
		{"current-player", ExitSuccess},
		{"valid-actions", ExitSuccess},
		{"scores", ExitSuccess},
		{"round", ExitSuccess},
		{"wall-remaining", ExitSuccess},
		{"invalid-query", ExitInvalidAction},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			batch := &BatchMode{
				GameFile: gameFile,
				Query:    tt.query,
				JSON:     true,
			}

			exitCode := batch.runQuery()
			if exitCode != tt.wantExit {
				t.Errorf("Query %s: expected exit code %d, got %d", tt.query, tt.wantExit, exitCode)
			}
		})
	}
}

func TestParseAction(t *testing.T) {
	config := engine.DefaultRuleConfig()
	game := engine.NewGame(config, 99999)

	tests := []struct {
		input       string
		wantType    engine.ActionType
		wantTiles   int
		wantErr     bool
	}{
		{"draw", engine.ActionDraw, 0, false},
		{"discard 3m", engine.ActionDiscard, 1, false},
		{"pon 5p", engine.ActionPon, 3, false},
		{"chi 345m", engine.ActionChi, 3, false},
		{"chi 3m4m5m", engine.ActionChi, 3, false},
		{"kan 7s", engine.ActionKan, 4, false},
		{"riichi", engine.ActionRiichi, 0, false},
		{"tsumo", engine.ActionTsumo, 0, false},
		{"ron", engine.ActionRon, 0, false},
		{"invalid", 0, 0, true},
		{"discard", 0, 0, true}, // Missing tile
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			action, err := parseAction(tt.input, game)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAction(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if action.Type != tt.wantType {
					t.Errorf("parseAction(%s) type = %v, want %v", tt.input, action.Type, tt.wantType)
				}
				if len(action.Tiles) != tt.wantTiles {
					t.Errorf("parseAction(%s) tiles count = %d, want %d", tt.input, len(action.Tiles), tt.wantTiles)
				}
			}
		})
	}
}

func TestParseSequence(t *testing.T) {
	tests := []struct {
		input   string
		want    int
		wantErr bool
	}{
		{"345m", 3, false},
		{"3m4m5m", 3, false},
		{"789s", 3, false},
		{"7s8s9s", 3, false},
		{"12m", 0, true}, // Invalid - rank must be 1-9
		{"invalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tiles, err := parseSequence(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSequence(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(tiles) != tt.want {
				t.Errorf("parseSequence(%s) length = %d, want %d", tt.input, len(tiles), tt.want)
			}
		})
	}
}

func TestBatchActionExecution(t *testing.T) {
	tempDir := t.TempDir()
	gameFile := filepath.Join(tempDir, "test_game.json")

	// Initialize game
	config := engine.DefaultRuleConfig()
	game := engine.NewGame(config, 11111)
	if err := serialization.SaveToFile(game, gameFile); err != nil {
		t.Fatalf("Failed to save initial game: %v", err)
	}

	// Get the drawn tile to discard
	drawnTile := game.Players[engine.SeatEast].DrawnTile.String()

	// Execute discard action
	batch := &BatchMode{
		GameFile: gameFile,
		Action:   "discard " + drawnTile,
		JSON:     false,
	}

	exitCode := batch.runAction()
	if exitCode != ExitSuccess {
		t.Errorf("Expected exit code %d, got %d", ExitSuccess, exitCode)
	}

	// Load and verify the game state was updated
	gameAfter, err := serialization.LoadFromFile(gameFile)
	if err != nil {
		t.Fatalf("Failed to load game after action: %v", err)
	}

	if len(gameAfter.History) != 1 {
		t.Errorf("Expected 1 action in history, got %d", len(gameAfter.History))
	}

	if gameAfter.CurrentPlayer != engine.SeatSouth {
		t.Errorf("Expected current player to be South, got %s", gameAfter.CurrentPlayer.String())
	}
}

func TestBatchMultipleActions(t *testing.T) {
	tempDir := t.TempDir()
	gameFile := filepath.Join(tempDir, "test_game.json")

	// Initialize game
	config := engine.DefaultRuleConfig()
	game := engine.NewGame(config, 22222)

	// Get the drawn tile to discard BEFORE saving
	drawnTile := game.Players[engine.SeatEast].DrawnTile.String()

	// Save the initial game
	if err := serialization.SaveToFile(game, gameFile); err != nil {
		t.Fatalf("Failed to save initial game: %v", err)
	}

	// Execute multiple actions
	batch := &BatchMode{
		GameFile: gameFile,
		Actions:  `["discard ` + drawnTile + `", "draw"]`,
		JSON:     false,
	}

	exitCode := batch.runAction()
	if exitCode != ExitSuccess {
		t.Errorf("Expected exit code %d, got %d", ExitSuccess, exitCode)
	}

	// Load and verify
	gameAfter, err := serialization.LoadFromFile(gameFile)
	if err != nil {
		t.Fatalf("Failed to load game after actions: %v", err)
	}

	if len(gameAfter.History) != 2 {
		t.Errorf("Expected 2 actions in history, got %d", len(gameAfter.History))
	}
}

func TestFormatAction(t *testing.T) {
	tile := engine.NewTile(engine.Man, 3)
	tests := []struct {
		action engine.Action
		want   string
	}{
		{engine.NewAction(engine.ActionDraw, engine.SeatEast, nil, 1), "draw"},
		{engine.NewAction(engine.ActionDiscard, engine.SeatEast, []engine.Tile{tile}, 1), "discard 3m"},
		{engine.NewAction(engine.ActionPon, engine.SeatEast, []engine.Tile{tile}, 1), "pon 3m"},
		{engine.NewAction(engine.ActionRiichi, engine.SeatEast, nil, 1), "riichi"},
		{engine.NewAction(engine.ActionTsumo, engine.SeatEast, nil, 1), "tsumo"},
		{engine.NewAction(engine.ActionRon, engine.SeatEast, nil, 1), "ron"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := formatAction(tt.action)
			if got != tt.want {
				t.Errorf("formatAction() = %s, want %s", got, tt.want)
			}
		})
	}
}
