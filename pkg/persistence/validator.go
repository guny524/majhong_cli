package persistence

import (
	"fmt"
)

// ValidationError represents a validation error with context
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error in %s: %s", e.Field, e.Message)
}

// Validate validates a game state structure
func Validate(state *GameState) error {
	if state == nil {
		return &ValidationError{Field: "game_state", Message: "game state is nil"}
	}

	// Validate players
	if err := validatePlayers(state.Players); err != nil {
		return err
	}

	// Validate round
	if err := validateRound(&state.Round); err != nil {
		return err
	}

	// Validate wall
	if err := validateWall(&state.Wall); err != nil {
		return err
	}

	// Validate config
	if err := validateConfig(&state.Config); err != nil {
		return err
	}

	// Validate current player
	if err := validateSeat(state.CurrentPlayer); err != nil {
		return &ValidationError{Field: "current_player", Message: err.Error()}
	}

	// Validate turn number
	if state.TurnNumber < 0 {
		return &ValidationError{Field: "turn_number", Message: "turn number cannot be negative"}
	}

	// Validate history
	if err := validateHistory(state.History); err != nil {
		return err
	}

	return nil
}

// validatePlayers validates the players array
func validatePlayers(players []Player) error {
	if len(players) != 4 {
		return &ValidationError{
			Field:   "players",
			Message: fmt.Sprintf("expected 4 players, got %d", len(players)),
		}
	}

	// Check for unique seats
	seats := make(map[string]bool)
	for i, player := range players {
		if err := validatePlayer(&player, i); err != nil {
			return err
		}

		if seats[player.Seat] {
			return &ValidationError{
				Field:   fmt.Sprintf("players[%d].seat", i),
				Message: fmt.Sprintf("duplicate seat: %s", player.Seat),
			}
		}
		seats[player.Seat] = true
	}

	// Ensure all seats are present
	requiredSeats := []string{"EAST", "SOUTH", "WEST", "NORTH"}
	for _, seat := range requiredSeats {
		if !seats[seat] {
			return &ValidationError{
				Field:   "players",
				Message: fmt.Sprintf("missing seat: %s", seat),
			}
		}
	}

	return nil
}

// validatePlayer validates a single player
func validatePlayer(player *Player, index int) error {
	prefix := fmt.Sprintf("players[%d]", index)

	// Validate seat
	if err := validateSeat(player.Seat); err != nil {
		return &ValidationError{Field: prefix + ".seat", Message: err.Error()}
	}

	// Validate hand tiles
	if err := validateTiles(player.Hand); err != nil {
		return &ValidationError{Field: prefix + ".hand", Message: err.Error()}
	}

	// Validate discards
	if err := validateTiles(player.Discards); err != nil {
		return &ValidationError{Field: prefix + ".discards", Message: err.Error()}
	}

	// Validate melds
	for i, meld := range player.Melds {
		if err := validateMeld(&meld); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("%s.melds[%d]", prefix, i),
				Message: err.Error(),
			}
		}
	}

	// Validate score
	if player.Score < 0 {
		return &ValidationError{
			Field:   prefix + ".score",
			Message: "score cannot be negative",
		}
	}

	// Validate hand size constraints
	totalTiles := len(player.Hand)
	for _, meld := range player.Melds {
		totalTiles += len(meld.Tiles)
	}

	// Maximum is 14 tiles (13 + 1 drawn, or special cases)
	if totalTiles > 18 {
		return &ValidationError{
			Field:   prefix,
			Message: fmt.Sprintf("too many tiles: %d (hand: %d, melds: %d)", totalTiles, len(player.Hand), totalTiles-len(player.Hand)),
		}
	}

	return nil
}

// validateRound validates the round information
func validateRound(round *Round) error {
	// Validate dealer
	if err := validateSeat(round.Dealer); err != nil {
		return &ValidationError{Field: "round.dealer", Message: err.Error()}
	}

	// Validate round wind
	if err := validateWind(round.RoundWind); err != nil {
		return &ValidationError{Field: "round.round_wind", Message: err.Error()}
	}

	// Validate round number
	if round.RoundNumber < 1 {
		return &ValidationError{
			Field:   "round.round_number",
			Message: "round number must be at least 1",
		}
	}

	// Validate honba
	if round.Honba < 0 {
		return &ValidationError{
			Field:   "round.honba",
			Message: "honba cannot be negative",
		}
	}

	// Validate riichi sticks
	if round.RiichiSticks < 0 || round.RiichiSticks > 4 {
		return &ValidationError{
			Field:   "round.riichi_sticks",
			Message: fmt.Sprintf("riichi sticks must be 0-4, got %d", round.RiichiSticks),
		}
	}

	// Validate dora indicators
	if len(round.DoraIndicators) < 1 || len(round.DoraIndicators) > 5 {
		return &ValidationError{
			Field:   "round.dora_indicators",
			Message: fmt.Sprintf("dora indicators must be 1-5, got %d", len(round.DoraIndicators)),
		}
	}

	if err := validateTiles(round.DoraIndicators); err != nil {
		return &ValidationError{Field: "round.dora_indicators", Message: err.Error()}
	}

	return nil
}

// validateWall validates the wall state
func validateWall(wall *Wall) error {
	// Validate remaining tiles
	if wall.Remaining < 0 || wall.Remaining > 122 {
		return &ValidationError{
			Field:   "wall.remaining",
			Message: fmt.Sprintf("remaining tiles must be 0-122, got %d", wall.Remaining),
		}
	}

	// Validate hash exists
	if wall.Hash == "" {
		return &ValidationError{
			Field:   "wall.hash",
			Message: "wall hash is required",
		}
	}

	// Validate salt hash exists
	if wall.SaltHash == "" {
		return &ValidationError{
			Field:   "wall.salt_hash",
			Message: "wall salt hash is required",
		}
	}

	return nil
}

// validateConfig validates the rule configuration
func validateConfig(config *RuleConfig) error {
	// Validate starting points
	validStartingPoints := []int{25000, 30000, 35000}
	valid := false
	for _, points := range validStartingPoints {
		if config.StartingPoints == points {
			valid = true
			break
		}
	}
	if !valid {
		return &ValidationError{
			Field:   "config.starting_points",
			Message: fmt.Sprintf("starting points must be 25000, 30000, or 35000, got %d", config.StartingPoints),
		}
	}

	// Validate return points
	if config.ReturnPoints < config.StartingPoints {
		return &ValidationError{
			Field:   "config.return_points",
			Message: fmt.Sprintf("return points (%d) must be >= starting points (%d)", config.ReturnPoints, config.StartingPoints),
		}
	}

	return nil
}

// validateHistory validates the action history
func validateHistory(history []Action) error {
	for i, action := range history {
		if err := validateAction(&action); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("history[%d]", i),
				Message: err.Error(),
			}
		}

		// Validate turn order
		if action.Turn != i+1 {
			return &ValidationError{
				Field:   fmt.Sprintf("history[%d].turn", i),
				Message: fmt.Sprintf("expected turn %d, got %d", i+1, action.Turn),
			}
		}
	}

	return nil
}

// validateAction validates a single action
func validateAction(action *Action) error {
	// Validate player
	if err := validateSeat(action.Player); err != nil {
		return fmt.Errorf("invalid player: %w", err)
	}

	// Validate action type
	validTypes := []string{"DRAW", "DISCARD", "PON", "CHI", "KAN", "RIICHI", "TSUMO", "RON"}
	valid := false
	for _, t := range validTypes {
		if action.Type == t {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid action type: %s", action.Type)
	}

	// Validate tiles
	if err := validateTiles(action.Tiles); err != nil {
		return fmt.Errorf("invalid tiles: %w", err)
	}

	// Validate tile count based on action type
	switch action.Type {
	case "DRAW":
		if len(action.Tiles) != 0 {
			return fmt.Errorf("draw action should have 0 tiles, got %d", len(action.Tiles))
		}
	case "DISCARD", "RIICHI":
		if len(action.Tiles) != 1 {
			return fmt.Errorf("%s action should have 1 tile, got %d", action.Type, len(action.Tiles))
		}
	case "PON", "CHI":
		if len(action.Tiles) != 3 {
			return fmt.Errorf("%s action should have 3 tiles, got %d", action.Type, len(action.Tiles))
		}
	case "KAN":
		if len(action.Tiles) != 4 {
			return fmt.Errorf("kan action should have 4 tiles, got %d", len(action.Tiles))
		}
	case "TSUMO", "RON":
		// Win actions can have varying tile counts
	}

	// Validate timestamp
	if action.Timestamp < 0 {
		return fmt.Errorf("timestamp cannot be negative")
	}

	return nil
}

// validateMeld validates a meld
func validateMeld(meld *Meld) error {
	// Validate type
	validTypes := []string{"PON", "CHI", "ANKAN", "MINKAN", "KAKAN"}
	valid := false
	for _, t := range validTypes {
		if meld.Type == t {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid meld type: %s", meld.Type)
	}

	// Validate tiles
	if err := validateTiles(meld.Tiles); err != nil {
		return fmt.Errorf("invalid tiles: %w", err)
	}

	// Validate tile count
	expectedCount := 3
	if meld.Type == "ANKAN" || meld.Type == "MINKAN" || meld.Type == "KAKAN" {
		expectedCount = 4
	}
	if len(meld.Tiles) != expectedCount {
		return fmt.Errorf("expected %d tiles for %s, got %d", expectedCount, meld.Type, len(meld.Tiles))
	}

	// Validate from player
	if err := validateSeat(meld.FromPlayer); err != nil {
		return fmt.Errorf("invalid from_player: %w", err)
	}

	return nil
}

// validateTiles validates a slice of tiles
func validateTiles(tiles []Tile) error {
	for i, tile := range tiles {
		if err := validateTile(&tile); err != nil {
			return fmt.Errorf("tile[%d]: %w", i, err)
		}
	}
	return nil
}

// validateTile validates a single tile
func validateTile(tile *Tile) error {
	// Validate suit
	if err := validateSuit(tile.Suit); err != nil {
		return err
	}

	// Validate rank based on suit
	switch tile.Suit {
	case "MAN", "PIN", "SOU":
		if tile.Rank < 1 || tile.Rank > 9 {
			return fmt.Errorf("rank must be 1-9 for suit %s, got %d", tile.Suit, tile.Rank)
		}
	case "JI":
		if tile.Rank < 1 || tile.Rank > 7 {
			return fmt.Errorf("rank must be 1-7 for honors, got %d", tile.Rank)
		}
	}

	// Validate aka dora (only on 5s of suited tiles)
	if tile.IsAkaDora {
		if tile.Suit == "JI" {
			return fmt.Errorf("aka dora cannot be honor tile")
		}
		if tile.Rank != 5 {
			return fmt.Errorf("aka dora must be rank 5, got %d", tile.Rank)
		}
	}

	return nil
}

// validateSeat validates a seat string
func validateSeat(seat string) error {
	validSeats := []string{"EAST", "SOUTH", "WEST", "NORTH"}
	for _, s := range validSeats {
		if seat == s {
			return nil
		}
	}
	return fmt.Errorf("invalid seat: %s", seat)
}

// validateWind validates a wind string
func validateWind(wind string) error {
	validWinds := []string{"EAST", "SOUTH", "WEST", "NORTH"}
	for _, w := range validWinds {
		if wind == w {
			return nil
		}
	}
	return fmt.Errorf("invalid wind: %s", wind)
}

// validateSuit validates a suit string
func validateSuit(suit string) error {
	validSuits := []string{"MAN", "PIN", "SOU", "JI"}
	for _, s := range validSuits {
		if suit == s {
			return nil
		}
	}
	return fmt.Errorf("invalid suit: %s", suit)
}
