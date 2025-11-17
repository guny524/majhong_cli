package engine

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	config := DefaultRuleConfig()
	game := NewGame(config, 123)

	// Check players initialized
	if len(game.Players) != 4 {
		t.Errorf("Expected 4 players, got %d", len(game.Players))
	}

	// Check starting scores
	for i, player := range game.Players {
		if player.Score != config.StartingPoints {
			t.Errorf("Player %d: expected %d points, got %d", i, config.StartingPoints, player.Score)
		}
	}

	// Check initial hands dealt
	for i, player := range game.Players {
		handSize := len(player.Hand)
		if player.DrawnTile != nil {
			handSize++
		}

		expected := 13
		if player.Seat == SeatEast {
			expected = 14 // Dealer gets extra tile
		}

		if handSize != expected {
			t.Errorf("Player %d: expected %d tiles, got %d", i, expected, handSize)
		}
	}

	// Check wall
	if game.Wall == nil {
		t.Error("Wall not initialized")
	}

	// Check round
	if game.Round.Dealer != SeatEast {
		t.Error("Expected dealer to be East")
	}
}

func TestGameApplyDrawDiscard(t *testing.T) {
	game := NewGame(DefaultRuleConfig(), 456)

	// Dealer (East) starts with drawn tile
	dealer := game.Players[SeatEast]

	// Discard the drawn tile
	discardTile := *dealer.DrawnTile
	action := NewAction(ActionDiscard, SeatEast, []Tile{discardTile}, 1)

	err := game.ApplyAction(action)
	if err != nil {
		t.Fatalf("Discard failed: %v", err)
	}

	// Check drawn tile removed
	if dealer.DrawnTile != nil {
		t.Error("Drawn tile should be nil after discard")
	}

	// Check tile in discards
	if len(dealer.Discards) != 1 {
		t.Errorf("Expected 1 discard, got %d", len(dealer.Discards))
	}

	// Check turn advanced to next player
	if game.CurrentPlayer != SeatSouth {
		t.Errorf("Expected current player to be South, got %v", game.CurrentPlayer)
	}
}

func TestGameRiichi(t *testing.T) {
	game := NewGame(DefaultRuleConfig(), 789)

	dealer := game.Players[SeatEast]
	initialScore := dealer.Score

	// Declare riichi
	action := NewAction(ActionRiichi, SeatEast, nil, 1)
	err := game.ApplyAction(action)
	if err != nil {
		t.Fatalf("Riichi failed: %v", err)
	}

	// Check riichi flag
	if !dealer.IsRiichi {
		t.Error("Player should be in riichi")
	}

	// Check score deducted
	if dealer.Score != initialScore-1000 {
		t.Errorf("Expected score %d, got %d", initialScore-1000, dealer.Score)
	}

	// Check riichi stick added
	if game.Round.RiichiSticks != 1 {
		t.Errorf("Expected 1 riichi stick, got %d", game.Round.RiichiSticks)
	}
}

func TestGameKan(t *testing.T) {
	game := NewGame(DefaultRuleConfig(), 111)

	player := game.Players[SeatEast]
	initialDoraCount := len(game.Wall.GetVisibleDoraIndicators())

	// Simulate kan with 4 tiles
	kanTiles := []Tile{
		NewTile(Man, 5),
		NewTile(Man, 5),
		NewTile(Man, 5),
		NewTile(Man, 5),
	}

	action := NewAction(ActionKan, SeatEast, kanTiles, 1)
	err := game.ApplyAction(action)
	if err != nil {
		t.Fatalf("Kan failed: %v", err)
	}

	// Check meld added
	if len(player.Melds) != 1 {
		t.Errorf("Expected 1 meld, got %d", len(player.Melds))
	}

	// Check rinshan tile drawn
	if player.DrawnTile == nil {
		t.Error("Should have rinshan tile after kan")
	}

	// Check dora indicator flipped
	newDoraCount := len(game.Wall.GetVisibleDoraIndicators())
	if newDoraCount != initialDoraCount+1 {
		t.Errorf("Expected %d dora indicators, got %d", initialDoraCount+1, newDoraCount)
	}
}

func TestGameValidActions(t *testing.T) {
	game := NewGame(DefaultRuleConfig(), 222)

	// Dealer should be able to discard (has drawn tile)
	actions := game.ValidActions()

	hasDiscard := false
	for _, action := range actions {
		if action.Type == ActionDiscard {
			hasDiscard = true
			break
		}
	}

	if !hasDiscard {
		t.Error("Dealer should have discard action available")
	}
}

func TestRoundNextDealer(t *testing.T) {
	round := NewRound(SeatEast, WindEast, 1)

	next := round.NextDealer()
	if next != SeatSouth {
		t.Errorf("Expected next dealer to be South, got %v", next)
	}

	round.Dealer = SeatNorth
	next = round.NextDealer()
	if next != SeatEast {
		t.Errorf("Expected next dealer to wrap to East, got %v", next)
	}
}

func TestRoundIsAllLast(t *testing.T) {
	// Not all last
	round := NewRound(SeatEast, WindEast, 1)
	if round.IsAllLast() {
		t.Error("East 1 should not be all last")
	}

	// All last (South 4)
	round = NewRound(SeatEast, WindSouth, 4)
	if !round.IsAllLast() {
		t.Error("South 4 should be all last")
	}
}
