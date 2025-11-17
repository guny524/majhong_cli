package persistence

import (
	"testing"

	"github.com/guny524/majhong_cli/pkg/mahjong/engine"
)

func TestValidateValidGameState(t *testing.T) {
	// Create a valid game state
	config := engine.DefaultRuleConfig()
	game := engine.NewGame(config, 12345)

	// Convert to serializable format
	state := convertGameState(game)

	// Validate
	err := Validate(&state)
	if err != nil {
		t.Errorf("Valid game state failed validation: %v", err)
	}
}

func TestValidateNilGameState(t *testing.T) {
	err := Validate(nil)
	if err == nil {
		t.Error("Expected error when validating nil game state")
	}
}

func TestValidateInvalidPlayerCount(t *testing.T) {
	state := &GameState{
		Players: []Player{}, // Empty players
		Round: Round{
			Dealer:         "EAST",
			RoundWind:      "EAST",
			RoundNumber:    1,
			Honba:          0,
			RiichiSticks:   0,
			DoraIndicators: []Tile{{Suit: "MAN", Rank: 1}},
		},
		Wall: Wall{
			Remaining: 70,
			Seed:      12345,
			Hash:      "test",
			SaltHash:  "test",
		},
		CurrentPlayer: "EAST",
		TurnNumber:    1,
		Config: RuleConfig{
			RiichiMahjong:  true,
			StartingPoints: 25000,
			ReturnPoints:   30000,
		},
	}

	err := Validate(state)
	if err == nil {
		t.Error("Expected error for invalid player count")
	}
}

func TestValidateDuplicateSeats(t *testing.T) {
	state := &GameState{
		Players: []Player{
			{Seat: "EAST", Score: 25000},
			{Seat: "EAST", Score: 25000}, // Duplicate
			{Seat: "WEST", Score: 25000},
			{Seat: "NORTH", Score: 25000},
		},
		Round: Round{
			Dealer:         "EAST",
			RoundWind:      "EAST",
			RoundNumber:    1,
			Honba:          0,
			RiichiSticks:   0,
			DoraIndicators: []Tile{{Suit: "MAN", Rank: 1}},
		},
		Wall: Wall{
			Remaining: 70,
			Seed:      12345,
			Hash:      "test",
			SaltHash:  "test",
		},
		CurrentPlayer: "EAST",
		TurnNumber:    1,
		Config: RuleConfig{
			RiichiMahjong:  true,
			StartingPoints: 25000,
			ReturnPoints:   30000,
		},
	}

	err := Validate(state)
	if err == nil {
		t.Error("Expected error for duplicate seats")
	}
}

func TestValidateInvalidSeat(t *testing.T) {
	state := &GameState{
		Players: []Player{
			{Seat: "INVALID", Score: 25000},
			{Seat: "SOUTH", Score: 25000},
			{Seat: "WEST", Score: 25000},
			{Seat: "NORTH", Score: 25000},
		},
		Round: Round{
			Dealer:         "EAST",
			RoundWind:      "EAST",
			RoundNumber:    1,
			Honba:          0,
			RiichiSticks:   0,
			DoraIndicators: []Tile{{Suit: "MAN", Rank: 1}},
		},
		Wall: Wall{
			Remaining: 70,
			Seed:      12345,
			Hash:      "test",
			SaltHash:  "test",
		},
		CurrentPlayer: "EAST",
		TurnNumber:    1,
		Config: RuleConfig{
			RiichiMahjong:  true,
			StartingPoints: 25000,
			ReturnPoints:   30000,
		},
	}

	err := Validate(state)
	if err == nil {
		t.Error("Expected error for invalid seat")
	}
}

func TestValidateNegativeScore(t *testing.T) {
	state := &GameState{
		Players: []Player{
			{Seat: "EAST", Score: -1000}, // Negative score
			{Seat: "SOUTH", Score: 25000},
			{Seat: "WEST", Score: 25000},
			{Seat: "NORTH", Score: 25000},
		},
		Round: Round{
			Dealer:         "EAST",
			RoundWind:      "EAST",
			RoundNumber:    1,
			Honba:          0,
			RiichiSticks:   0,
			DoraIndicators: []Tile{{Suit: "MAN", Rank: 1}},
		},
		Wall: Wall{
			Remaining: 70,
			Seed:      12345,
			Hash:      "test",
			SaltHash:  "test",
		},
		CurrentPlayer: "EAST",
		TurnNumber:    1,
		Config: RuleConfig{
			RiichiMahjong:  true,
			StartingPoints: 25000,
			ReturnPoints:   30000,
		},
	}

	err := Validate(state)
	if err == nil {
		t.Error("Expected error for negative score")
	}
}

func TestValidateInvalidTile(t *testing.T) {
	tests := []struct {
		name string
		tile Tile
	}{
		{
			name: "Invalid suit",
			tile: Tile{Suit: "INVALID", Rank: 5},
		},
		{
			name: "Invalid rank for suited tile",
			tile: Tile{Suit: "MAN", Rank: 10},
		},
		{
			name: "Invalid rank for honor tile",
			tile: Tile{Suit: "JI", Rank: 8},
		},
		{
			name: "Aka dora on honor tile",
			tile: Tile{Suit: "JI", Rank: 1, IsAkaDora: true},
		},
		{
			name: "Aka dora on non-5 rank",
			tile: Tile{Suit: "MAN", Rank: 3, IsAkaDora: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTile(&tt.tile)
			if err == nil {
				t.Errorf("Expected error for %s", tt.name)
			}
		})
	}
}

func TestValidateRoundNumber(t *testing.T) {
	state := &GameState{
		Players: []Player{
			{Seat: "EAST", Score: 25000},
			{Seat: "SOUTH", Score: 25000},
			{Seat: "WEST", Score: 25000},
			{Seat: "NORTH", Score: 25000},
		},
		Round: Round{
			Dealer:         "EAST",
			RoundWind:      "EAST",
			RoundNumber:    0, // Invalid
			Honba:          0,
			RiichiSticks:   0,
			DoraIndicators: []Tile{{Suit: "MAN", Rank: 1}},
		},
		Wall: Wall{
			Remaining: 70,
			Seed:      12345,
			Hash:      "test",
			SaltHash:  "test",
		},
		CurrentPlayer: "EAST",
		TurnNumber:    1,
		Config: RuleConfig{
			RiichiMahjong:  true,
			StartingPoints: 25000,
			ReturnPoints:   30000,
		},
	}

	err := Validate(state)
	if err == nil {
		t.Error("Expected error for invalid round number")
	}
}

func TestValidateDoraIndicators(t *testing.T) {
	tests := []struct {
		name       string
		indicators []Tile
		wantError  bool
	}{
		{
			name:       "Too few dora indicators",
			indicators: []Tile{},
			wantError:  true,
		},
		{
			name: "Valid single dora",
			indicators: []Tile{
				{Suit: "MAN", Rank: 1},
			},
			wantError: false,
		},
		{
			name: "Valid multiple dora",
			indicators: []Tile{
				{Suit: "MAN", Rank: 1},
				{Suit: "PIN", Rank: 2},
				{Suit: "SOU", Rank: 3},
			},
			wantError: false,
		},
		{
			name: "Too many dora indicators",
			indicators: []Tile{
				{Suit: "MAN", Rank: 1},
				{Suit: "PIN", Rank: 2},
				{Suit: "SOU", Rank: 3},
				{Suit: "MAN", Rank: 4},
				{Suit: "PIN", Rank: 5},
				{Suit: "SOU", Rank: 6}, // 6 dora
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &GameState{
				Players: []Player{
					{Seat: "EAST", Score: 25000},
					{Seat: "SOUTH", Score: 25000},
					{Seat: "WEST", Score: 25000},
					{Seat: "NORTH", Score: 25000},
				},
				Round: Round{
					Dealer:         "EAST",
					RoundWind:      "EAST",
					RoundNumber:    1,
					Honba:          0,
					RiichiSticks:   0,
					DoraIndicators: tt.indicators,
				},
				Wall: Wall{
					Remaining: 70,
					Seed:      12345,
					Hash:      "test",
					SaltHash:  "test",
				},
				CurrentPlayer: "EAST",
				TurnNumber:    1,
				Config: RuleConfig{
					RiichiMahjong:  true,
					StartingPoints: 25000,
					ReturnPoints:   30000,
				},
			}

			err := Validate(state)
			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestValidateStartingPoints(t *testing.T) {
	tests := []struct {
		points       int
		returnPoints int
		wantError    bool
	}{
		{25000, 30000, false},
		{30000, 30000, false},
		{35000, 35000, false},
		{20000, 30000, true},
		{40000, 40000, true},
	}

	for _, tt := range tests {
		state := &GameState{
			Players: []Player{
				{Seat: "EAST", Score: tt.points},
				{Seat: "SOUTH", Score: tt.points},
				{Seat: "WEST", Score: tt.points},
				{Seat: "NORTH", Score: tt.points},
			},
			Round: Round{
				Dealer:         "EAST",
				RoundWind:      "EAST",
				RoundNumber:    1,
				Honba:          0,
				RiichiSticks:   0,
				DoraIndicators: []Tile{{Suit: "MAN", Rank: 1}},
			},
			Wall: Wall{
				Remaining: 70,
				Seed:      12345,
				Hash:      "test",
				SaltHash:  "test",
			},
			CurrentPlayer: "EAST",
			TurnNumber:    1,
			Config: RuleConfig{
				RiichiMahjong:  true,
				StartingPoints: tt.points,
				ReturnPoints:   tt.returnPoints,
			},
		}

		err := Validate(state)
		if tt.wantError && err == nil {
			t.Errorf("Expected error for starting points %d", tt.points)
		}
		if !tt.wantError && err != nil {
			t.Errorf("Unexpected error for starting points %d: %v", tt.points, err)
		}
	}
}
