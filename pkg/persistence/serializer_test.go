package persistence

import (
	"testing"

	"github.com/guny524/majhong_cli/pkg/mahjong/engine"
)

func TestSerializeDeserialize(t *testing.T) {
	// Create a test game state
	config := engine.DefaultRuleConfig()
	game := engine.NewGame(config, 12345)

	// Serialize the game state
	data, err := Serialize(game)
	if err != nil {
		t.Fatalf("Failed to serialize game state: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("Serialized data is empty")
	}

	// Deserialize the game state
	loadedGame, err := Deserialize(data)
	if err != nil {
		t.Fatalf("Failed to deserialize game state: %v", err)
	}

	// Verify key fields
	if loadedGame.Wall.Seed != game.Wall.Seed {
		t.Errorf("Seed mismatch: expected %d, got %d", game.Wall.Seed, loadedGame.Wall.Seed)
	}

	if loadedGame.Round.Dealer != game.Round.Dealer {
		t.Errorf("Dealer mismatch: expected %v, got %v", game.Round.Dealer, loadedGame.Round.Dealer)
	}

	if loadedGame.Round.RoundWind != game.Round.RoundWind {
		t.Errorf("Round wind mismatch: expected %v, got %v", game.Round.RoundWind, loadedGame.Round.RoundWind)
	}

	if loadedGame.CurrentPlayer != game.CurrentPlayer {
		t.Errorf("Current player mismatch: expected %v, got %v", game.CurrentPlayer, loadedGame.CurrentPlayer)
	}

	if loadedGame.TurnNumber != game.TurnNumber {
		t.Errorf("Turn number mismatch: expected %d, got %d", game.TurnNumber, loadedGame.TurnNumber)
	}

	if loadedGame.IsGameOver != game.IsGameOver {
		t.Errorf("Game over mismatch: expected %v, got %v", game.IsGameOver, loadedGame.IsGameOver)
	}

	// Verify players
	for i := 0; i < 4; i++ {
		if loadedGame.Players[i].Score != game.Players[i].Score {
			t.Errorf("Player %d score mismatch: expected %d, got %d",
				i, game.Players[i].Score, loadedGame.Players[i].Score)
		}

		if len(loadedGame.Players[i].Hand) != len(game.Players[i].Hand) {
			t.Errorf("Player %d hand size mismatch: expected %d, got %d",
				i, len(game.Players[i].Hand), len(loadedGame.Players[i].Hand))
		}
	}
}

func TestSerializeNilGameState(t *testing.T) {
	_, err := Serialize(nil)
	if err == nil {
		t.Error("Expected error when serializing nil game state")
	}
}

func TestDeserializeEmptyData(t *testing.T) {
	_, err := Deserialize([]byte{})
	if err == nil {
		t.Error("Expected error when deserializing empty data")
	}
}

func TestDeserializeInvalidJSON(t *testing.T) {
	invalidJSON := []byte(`{"invalid": json}`)
	_, err := Deserialize(invalidJSON)
	if err == nil {
		t.Error("Expected error when deserializing invalid JSON")
	}
}

func TestTileConversion(t *testing.T) {
	tests := []struct {
		name string
		tile engine.Tile
	}{
		{
			name: "Man tile",
			tile: engine.NewTile(engine.Man, 5),
		},
		{
			name: "Pin tile",
			tile: engine.NewTile(engine.Pin, 3),
		},
		{
			name: "Sou tile",
			tile: engine.NewTile(engine.Sou, 7),
		},
		{
			name: "Ji tile (honors)",
			tile: engine.NewTile(engine.Ji, 1),
		},
		{
			name: "Aka dora tile",
			tile: engine.NewAkaDoraTile(engine.Man),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert to serializable format
			serialized := convertTile(tt.tile)

			// Convert back to engine format
			deserialized, err := convertToEngineTile(serialized)
			if err != nil {
				t.Fatalf("Failed to convert tile: %v", err)
			}

			// Verify
			if deserialized.Suit != tt.tile.Suit {
				t.Errorf("Suit mismatch: expected %v, got %v", tt.tile.Suit, deserialized.Suit)
			}
			if deserialized.Rank != tt.tile.Rank {
				t.Errorf("Rank mismatch: expected %d, got %d", tt.tile.Rank, deserialized.Rank)
			}
			if deserialized.IsAkaDora != tt.tile.IsAkaDora {
				t.Errorf("IsAkaDora mismatch: expected %v, got %v", tt.tile.IsAkaDora, deserialized.IsAkaDora)
			}
		})
	}
}

func TestSeatConversion(t *testing.T) {
	tests := []struct {
		seat     engine.Seat
		expected string
	}{
		{engine.SeatEast, "EAST"},
		{engine.SeatSouth, "SOUTH"},
		{engine.SeatWest, "WEST"},
		{engine.SeatNorth, "NORTH"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			// Convert to string
			str := seatToString(tt.seat)
			if str != tt.expected {
				t.Errorf("Seat to string: expected %s, got %s", tt.expected, str)
			}

			// Convert back
			seat, err := stringToSeat(str)
			if err != nil {
				t.Fatalf("String to seat failed: %v", err)
			}
			if seat != tt.seat {
				t.Errorf("String to seat: expected %v, got %v", tt.seat, seat)
			}
		})
	}
}

func TestWindConversion(t *testing.T) {
	tests := []struct {
		wind     engine.Wind
		expected string
	}{
		{engine.WindEast, "EAST"},
		{engine.WindSouth, "SOUTH"},
		{engine.WindWest, "WEST"},
		{engine.WindNorth, "NORTH"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			// Convert to string
			str := windToString(tt.wind)
			if str != tt.expected {
				t.Errorf("Wind to string: expected %s, got %s", tt.expected, str)
			}

			// Convert back
			wind, err := stringToWind(str)
			if err != nil {
				t.Fatalf("String to wind failed: %v", err)
			}
			if wind != tt.wind {
				t.Errorf("String to wind: expected %v, got %v", tt.wind, wind)
			}
		})
	}
}

func TestActionTypeConversion(t *testing.T) {
	tests := []struct {
		actionType engine.ActionType
		expected   string
	}{
		{engine.ActionDraw, "DRAW"},
		{engine.ActionDiscard, "DISCARD"},
		{engine.ActionPon, "PON"},
		{engine.ActionChi, "CHI"},
		{engine.ActionKan, "KAN"},
		{engine.ActionRiichi, "RIICHI"},
		{engine.ActionTsumo, "TSUMO"},
		{engine.ActionRon, "RON"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			// Convert to string
			str := actionTypeToString(tt.actionType)
			if str != tt.expected {
				t.Errorf("Action type to string: expected %s, got %s", tt.expected, str)
			}

			// Convert back
			actionType, err := stringToActionType(str)
			if err != nil {
				t.Fatalf("String to action type failed: %v", err)
			}
			if actionType != tt.actionType {
				t.Errorf("String to action type: expected %v, got %v", tt.actionType, actionType)
			}
		})
	}
}
