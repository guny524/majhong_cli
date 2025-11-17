package serialization

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/guny524/majhong_cli/pkg/mahjong/engine"
)

func TestToJSONFromJSON(t *testing.T) {
	// Create a new game
	config := engine.DefaultRuleConfig()
	game := engine.NewGame(config, 12345)

	// Convert to JSON
	jgs, err := ToJSON(game)
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	// Verify JSON structure
	if jgs.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", jgs.Version)
	}

	if jgs.Seed != 12345 {
		t.Errorf("Expected seed 12345, got %d", jgs.Seed)
	}

	if len(jgs.Players) != 4 {
		t.Errorf("Expected 4 players, got %d", len(jgs.Players))
	}

	// Convert back from JSON
	gameRestored, err := FromJSON(jgs)
	if err != nil {
		t.Fatalf("FromJSON failed: %v", err)
	}

	// Verify restored game
	if gameRestored.Wall.Seed != game.Wall.Seed {
		t.Errorf("Seed mismatch: expected %d, got %d", game.Wall.Seed, gameRestored.Wall.Seed)
	}

	if gameRestored.CurrentPlayer != game.CurrentPlayer {
		t.Errorf("Current player mismatch: expected %s, got %s",
			game.CurrentPlayer.String(), gameRestored.CurrentPlayer.String())
	}

	// Verify players
	for i := 0; i < 4; i++ {
		if gameRestored.Players[i].Score != game.Players[i].Score {
			t.Errorf("Player %d score mismatch: expected %d, got %d",
				i, game.Players[i].Score, gameRestored.Players[i].Score)
		}
	}
}

func TestSaveLoadFile(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_game.json")

	// Create a new game
	config := engine.DefaultRuleConfig()
	game := engine.NewGame(config, 54321)

	// Apply some actions
	action := engine.NewAction(engine.ActionDiscard, engine.SeatEast, []engine.Tile{*game.Players[engine.SeatEast].DrawnTile}, 1)
	if err := game.ApplyAction(action); err != nil {
		t.Fatalf("Failed to apply action: %v", err)
	}

	// Save to file
	if err := SaveToFile(game, filePath); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("File was not created: %s", filePath)
	}

	// Load from file
	gameLoaded, err := LoadFromFile(filePath)
	if err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}

	// Verify loaded game
	if gameLoaded.Wall.Seed != game.Wall.Seed {
		t.Errorf("Seed mismatch after load: expected %d, got %d", game.Wall.Seed, gameLoaded.Wall.Seed)
	}

	if gameLoaded.TurnNumber != game.TurnNumber {
		t.Errorf("Turn number mismatch: expected %d, got %d", game.TurnNumber, gameLoaded.TurnNumber)
	}

	if len(gameLoaded.History) != len(game.History) {
		t.Errorf("History length mismatch: expected %d, got %d", len(game.History), len(gameLoaded.History))
	}
}

func TestLoadFromFileNotFound(t *testing.T) {
	_, err := LoadFromFile("/nonexistent/file.json")
	if err == nil {
		t.Fatal("Expected error for nonexistent file")
	}
}

func TestParseSeat(t *testing.T) {
	tests := []struct {
		input    string
		expected engine.Seat
		wantErr  bool
	}{
		{"East", engine.SeatEast, false},
		{"South", engine.SeatSouth, false},
		{"West", engine.SeatWest, false},
		{"North", engine.SeatNorth, false},
		{"Invalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseSeat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSeat(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("parseSeat(%s) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseActionType(t *testing.T) {
	tests := []struct {
		input    string
		expected engine.ActionType
		wantErr  bool
	}{
		{"draw", engine.ActionDraw, false},
		{"discard", engine.ActionDiscard, false},
		{"pon", engine.ActionPon, false},
		{"chi", engine.ActionChi, false},
		{"kan", engine.ActionKan, false},
		{"riichi", engine.ActionRiichi, false},
		{"tsumo", engine.ActionTsumo, false},
		{"ron", engine.ActionRon, false},
		{"invalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseActionType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseActionType(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("parseActionType(%s) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseMeldType(t *testing.T) {
	tests := []struct {
		input    string
		expected engine.MeldType
		wantErr  bool
	}{
		{"chi", engine.MeldChi, false},
		{"pon", engine.MeldPon, false},
		{"ankan", engine.MeldAnkan, false},
		{"minkan", engine.MeldMinkan, false},
		{"kakan", engine.MeldKakan, false},
		{"invalid", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseMeldType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMeldType(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("parseMeldType(%s) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
