package engine

import (
	"fmt"
)

// RuleConfig represents game rule configuration
type RuleConfig struct {
	RiichiMahjong  bool
	StartingPoints int
	ReturnPoints   int
	AkaDora        bool
	LocalYaku      []string
}

// DefaultRuleConfig returns standard Riichi Mahjong rules
func DefaultRuleConfig() RuleConfig {
	return RuleConfig{
		RiichiMahjong:  true,
		StartingPoints: 25000,
		ReturnPoints:   30000,
		AkaDora:        true,
		LocalYaku:      []string{},
	}
}

// GameState represents the complete game state
type GameState struct {
	Round         *Round
	Players       [4]*Player
	Wall          *Wall
	CurrentPlayer Seat
	TurnNumber    int
	History       []Action
	Config        RuleConfig
	IsGameOver    bool
}

// NewGame creates a new game with given configuration and seed
func NewGame(config RuleConfig, seed int64) *GameState {
	// Create wall
	wall := NewWallWithConfig(seed, config.AkaDora)

	// Create players
	var players [4]*Player
	for i := 0; i < 4; i++ {
		players[i] = NewPlayer(Seat(i))
		players[i].Score = config.StartingPoints
	}

	// Create initial round (East 1)
	round := NewRound(SeatEast, WindEast, 1)

	// Deal initial hands
	dealInitialHands(&players, wall)

	return &GameState{
		Round:         round,
		Players:       players,
		Wall:          wall,
		CurrentPlayer: SeatEast, // Dealer starts
		TurnNumber:    1,
		History:       make([]Action, 0),
		Config:        config,
		IsGameOver:    false,
	}
}

// dealInitialHands deals 13 tiles to each player, 14 to dealer
func dealInitialHands(players *[4]*Player, wall *Wall) {
	// Deal 12 tiles to each player (4 tiles × 3 rounds)
	for round := 0; round < 3; round++ {
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				tile, _ := wall.Draw()
				players[i].Hand = append(players[i].Hand, tile)
			}
		}
	}

	// Deal 1 tile to each player (total 13)
	for i := 0; i < 4; i++ {
		tile, _ := wall.Draw()
		players[i].Hand = append(players[i].Hand, tile)
	}

	// Dealer gets one extra tile (14 total)
	tile, _ := wall.Draw()
	players[SeatEast].Draw(tile)
}

// ApplyAction applies an action to the game state
func (g *GameState) ApplyAction(action Action) error {
	if g.IsGameOver {
		return fmt.Errorf("game is over")
	}

	// Validate action is from current player
	if action.Type != ActionRon && action.Type != ActionPon &&
	   action.Type != ActionChi && action.Player != g.CurrentPlayer {
		return fmt.Errorf("not player's turn")
	}

	switch action.Type {
	case ActionDraw:
		return g.applyDraw(action)
	case ActionDiscard:
		return g.applyDiscard(action)
	case ActionRiichi:
		return g.applyRiichi(action)
	case ActionTsumo:
		return g.applyTsumo(action)
	case ActionRon:
		return g.applyRon(action)
	case ActionPon:
		return g.applyPon(action)
	case ActionChi:
		return g.applyChi(action)
	case ActionKan:
		return g.applyKan(action)
	default:
		return fmt.Errorf("unknown action type: %v", action.Type)
	}
}

func (g *GameState) applyDraw(action Action) error {
	tile, err := g.Wall.Draw()
	if err != nil {
		// Wall exhausted - ryuukyoku
		g.IsGameOver = true
		return err
	}

	player := g.Players[action.Player]
	player.Draw(tile)

	g.History = append(g.History, action)
	return nil
}

func (g *GameState) applyDiscard(action Action) error {
	if len(action.Tiles) != 1 {
		return fmt.Errorf("discard must have exactly 1 tile")
	}

	player := g.Players[action.Player]
	err := player.Discard(action.Tiles[0])
	if err != nil {
		return err
	}

	g.History = append(g.History, action)

	// Advance to next player
	g.advanceTurn()
	return nil
}

func (g *GameState) applyRiichi(action Action) error {
	player := g.Players[action.Player]
	err := player.DeclareRiichi(g.TurnNumber)
	if err != nil {
		return err
	}

	g.Round.RiichiSticks++
	g.History = append(g.History, action)
	return nil
}

func (g *GameState) applyTsumo(action Action) error {
	// TODO: Validate winning hand and calculate score
	// For now, just mark game over
	g.IsGameOver = true
	g.History = append(g.History, action)
	return nil
}

func (g *GameState) applyRon(action Action) error {
	// TODO: Validate winning hand and calculate score
	// For now, just mark game over
	g.IsGameOver = true
	g.History = append(g.History, action)
	return nil
}

func (g *GameState) applyPon(action Action) error {
	if len(action.Tiles) != 3 {
		return fmt.Errorf("pon must have exactly 3 tiles")
	}

	player := g.Players[action.Player]

	// Remove 2 tiles from hand (3rd comes from discard)
	err := player.RemoveTilesFromHand(action.Tiles[:2])
	if err != nil {
		return err
	}

	// Add meld
	lastDiscard := g.getLastDiscard()
	meld := NewPon(action.Tiles, lastDiscard.Player)
	player.AddMeld(meld)

	g.History = append(g.History, action)
	g.CurrentPlayer = action.Player
	return nil
}

func (g *GameState) applyChi(action Action) error {
	if len(action.Tiles) != 3 {
		return fmt.Errorf("chi must have exactly 3 tiles")
	}

	player := g.Players[action.Player]

	// Remove 2 tiles from hand
	err := player.RemoveTilesFromHand(action.Tiles[:2])
	if err != nil {
		return err
	}

	// Add meld
	lastDiscard := g.getLastDiscard()
	meld := NewChi(action.Tiles, lastDiscard.Player)
	player.AddMeld(meld)

	g.History = append(g.History, action)
	g.CurrentPlayer = action.Player
	return nil
}

func (g *GameState) applyKan(action Action) error {
	if len(action.Tiles) != 4 {
		return fmt.Errorf("kan must have exactly 4 tiles")
	}

	player := g.Players[action.Player]

	// Determine kan type and create meld
	// For simplicity, assume ankan if all 4 tiles in hand
	meld := NewAnkan(action.Tiles)
	player.AddMeld(meld)

	// Draw rinshan tile
	rinshan, err := g.Wall.DrawRinshan()
	if err != nil {
		return err
	}
	player.Draw(rinshan)

	// Flip next dora indicator
	g.Wall.FlipNextDoraIndicator()

	g.History = append(g.History, action)
	return nil
}

func (g *GameState) advanceTurn() {
	g.CurrentPlayer = (g.CurrentPlayer + 1) % 4
	g.TurnNumber++
}

func (g *GameState) getLastDiscard() *Action {
	for i := len(g.History) - 1; i >= 0; i-- {
		if g.History[i].Type == ActionDiscard {
			return &g.History[i]
		}
	}
	return nil
}

// ValidActions returns all valid actions for current player
func (g *GameState) ValidActions() []Action {
	// TODO: Implement full action validation
	// For now, return basic actions
	actions := make([]Action, 0)

	player := g.Players[g.CurrentPlayer]

	// Can always draw (if not already drawn)
	if player.DrawnTile == nil && !g.Wall.IsExhausted() {
		actions = append(actions, NewAction(ActionDraw, g.CurrentPlayer, nil, g.TurnNumber))
	}

	// Can discard if tile was drawn
	if player.DrawnTile != nil {
		// Can discard any tile in hand or drawn tile
		for _, tile := range player.Hand {
			actions = append(actions, NewAction(ActionDiscard, g.CurrentPlayer, []Tile{tile}, g.TurnNumber))
		}
		actions = append(actions, NewAction(ActionDiscard, g.CurrentPlayer, []Tile{*player.DrawnTile}, g.TurnNumber))
	}

	return actions
}

// Clone creates a deep copy of game state
func (g *GameState) Clone() *GameState {
	// TODO: Implement full deep copy
	// For now, return shallow copy
	return g
}
