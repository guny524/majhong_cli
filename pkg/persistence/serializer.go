package persistence

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/guny524/majhong_cli/pkg/mahjong/engine"
)

const (
	// FormatVersion is the current save format version
	FormatVersion = "1.0.0"
)

// SaveFormat represents the serializable game state format
type SaveFormat struct {
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	GameState GameState `json:"game_state"`
}

// GameState represents the serializable game state
type GameState struct {
	Seed          int64        `json:"seed"`
	Round         Round        `json:"round"`
	Players       []Player     `json:"players"`
	Wall          Wall         `json:"wall"`
	CurrentPlayer string       `json:"current_player"`
	TurnNumber    int          `json:"turn_number"`
	History       []Action     `json:"history"`
	Config        RuleConfig   `json:"config"`
	IsGameOver    bool         `json:"is_game_over"`
}

// Round represents serializable round information
type Round struct {
	Dealer         string   `json:"dealer"`
	RoundWind      string   `json:"round_wind"`
	RoundNumber    int      `json:"round_number"`
	Honba          int      `json:"honba"`
	RiichiSticks   int      `json:"riichi_sticks"`
	DoraIndicators []Tile   `json:"dora_indicators"`
}

// Player represents serializable player state
type Player struct {
	Seat     string        `json:"seat"`
	Hand     []Tile        `json:"hand"`
	Discards []Tile        `json:"discards"`
	Melds    []Meld        `json:"melds"`
	Score    int           `json:"score"`
	IsRiichi bool          `json:"is_riichi"`
	Furiten  FuritenState  `json:"furiten"`
	IsTenpai bool          `json:"is_tenpai"`
}

// Wall represents serializable wall state
type Wall struct {
	Remaining int    `json:"remaining"`
	Seed      int64  `json:"seed"`
	Hash      string `json:"hash"`
	SaltHash  string `json:"salt_hash"`
}

// Tile represents a serializable tile
type Tile struct {
	Suit      string `json:"suit"`
	Rank      int    `json:"rank"`
	IsAkaDora bool   `json:"is_aka_dora"`
}

// Action represents a serializable action
type Action struct {
	Turn      int      `json:"turn"`
	Player    string   `json:"player"`
	Type      string   `json:"type"`
	Tiles     []Tile   `json:"tiles"`
	Timestamp int64    `json:"timestamp"`
}

// Meld represents a serializable meld
type Meld struct {
	Type       string `json:"type"`
	Tiles      []Tile `json:"tiles"`
	FromPlayer string `json:"from_player"`
}

// FuritenState represents serializable furiten state
type FuritenState struct {
	DiscardFuriten   bool `json:"discard_furiten"`
	TemporaryFuriten bool `json:"temporary_furiten"`
	RiichiFuriten    bool `json:"riichi_furiten"`
}

// RuleConfig represents serializable rule configuration
type RuleConfig struct {
	RiichiMahjong  bool     `json:"riichi_mahjong"`
	StartingPoints int      `json:"starting_points"`
	ReturnPoints   int      `json:"return_points"`
	AkaDora        bool     `json:"aka_dora"`
	LocalYaku      []string `json:"local_yaku"`
}

// Serialize converts an engine.GameState to JSON bytes
func Serialize(state *engine.GameState) ([]byte, error) {
	if state == nil {
		return nil, fmt.Errorf("game state is nil")
	}

	saveFormat := SaveFormat{
		Version:   FormatVersion,
		Timestamp: time.Now(),
		GameState: convertGameState(state),
	}

	data, err := json.MarshalIndent(saveFormat, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal game state: %w", err)
	}

	return data, nil
}

// Deserialize converts JSON bytes to an engine.GameState
func Deserialize(data []byte) (*engine.GameState, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	var saveFormat SaveFormat
	if err := json.Unmarshal(data, &saveFormat); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game state: %w", err)
	}

	// Verify version compatibility
	if saveFormat.Version != FormatVersion {
		return nil, fmt.Errorf("incompatible save format version: %s (expected %s)", saveFormat.Version, FormatVersion)
	}

	state, err := convertToEngineGameState(&saveFormat.GameState)
	if err != nil {
		return nil, fmt.Errorf("failed to convert game state: %w", err)
	}

	return state, nil
}

// convertGameState converts engine.GameState to serializable format
func convertGameState(state *engine.GameState) GameState {
	round := convertRound(state.Round)

	// Add dora indicators from wall
	doraIndicators := state.Wall.GetVisibleDoraIndicators()
	round.DoraIndicators = convertTiles(doraIndicators)

	return GameState{
		Seed:          state.Wall.Seed,
		Round:         round,
		Players:       convertPlayers(state.Players),
		Wall:          convertWall(state.Wall),
		CurrentPlayer: seatToString(state.CurrentPlayer),
		TurnNumber:    state.TurnNumber,
		History:       convertHistory(state.History),
		Config:        convertRuleConfig(state.Config),
		IsGameOver:    state.IsGameOver,
	}
}

// convertRound converts engine.Round to serializable format
func convertRound(round *engine.Round) Round {
	return Round{
		Dealer:         seatToString(round.Dealer),
		RoundWind:      windToString(round.RoundWind),
		RoundNumber:    round.RoundNumber,
		Honba:          round.Honba,
		RiichiSticks:   round.RiichiSticks,
		DoraIndicators: []Tile{}, // Will be populated from Wall
	}
}

// convertPlayers converts engine.Player array to serializable format
func convertPlayers(players [4]*engine.Player) []Player {
	result := make([]Player, 4)
	for i, p := range players {
		result[i] = convertPlayer(p)
	}
	return result
}

// convertPlayer converts engine.Player to serializable format
func convertPlayer(player *engine.Player) Player {
	return Player{
		Seat:     seatToString(player.Seat),
		Hand:     convertTiles(player.Hand),
		Discards: convertTiles(player.Discards),
		Melds:    convertMelds(player.Melds),
		Score:    player.Score,
		IsRiichi: player.IsRiichi,
		Furiten: FuritenState{
			DiscardFuriten:   player.Furiten.DiscardFuriten,
			TemporaryFuriten: player.Furiten.TemporaryFuriten,
			RiichiFuriten:    player.Furiten.RiichiFuriten,
		},
		IsTenpai: player.IsTenpai,
	}
}

// convertWall converts engine.Wall to serializable format
func convertWall(wall *engine.Wall) Wall {
	return Wall{
		Remaining: len(wall.LiveWall),
		Seed:      wall.Seed,
		Hash:      wall.Hash,
		SaltHash:  wall.SaltHash,
	}
}

// convertTiles converts slice of engine.Tile to serializable format
func convertTiles(tiles []engine.Tile) []Tile {
	result := make([]Tile, len(tiles))
	for i, t := range tiles {
		result[i] = convertTile(t)
	}
	return result
}

// convertTile converts engine.Tile to serializable format
func convertTile(tile engine.Tile) Tile {
	return Tile{
		Suit:      suitToString(tile.Suit),
		Rank:      tile.Rank,
		IsAkaDora: tile.IsAkaDora,
	}
}

// convertHistory converts action history to serializable format
func convertHistory(history []engine.Action) []Action {
	result := make([]Action, len(history))
	for i, a := range history {
		result[i] = convertAction(a)
	}
	return result
}

// convertAction converts engine.Action to serializable format
func convertAction(action engine.Action) Action {
	return Action{
		Turn:      action.Turn,
		Player:    seatToString(action.Player),
		Type:      actionTypeToString(action.Type),
		Tiles:     convertTiles(action.Tiles),
		Timestamp: action.Timestamp.Unix(),
	}
}

// convertMelds converts slice of engine.Meld to serializable format
func convertMelds(melds []engine.Meld) []Meld {
	result := make([]Meld, len(melds))
	for i, m := range melds {
		result[i] = convertMeld(m)
	}
	return result
}

// convertMeld converts engine.Meld to serializable format
func convertMeld(meld engine.Meld) Meld {
	return Meld{
		Type:       meldTypeToString(meld.Type),
		Tiles:      convertTiles(meld.Tiles),
		FromPlayer: seatToString(meld.FromPlayer),
	}
}

// convertRuleConfig converts engine.RuleConfig to serializable format
func convertRuleConfig(config engine.RuleConfig) RuleConfig {
	return RuleConfig{
		RiichiMahjong:  config.RiichiMahjong,
		StartingPoints: config.StartingPoints,
		ReturnPoints:   config.ReturnPoints,
		AkaDora:        config.AkaDora,
		LocalYaku:      config.LocalYaku,
	}
}

// Helper functions for string conversion
func seatToString(seat engine.Seat) string {
	switch seat {
	case engine.SeatEast:
		return "EAST"
	case engine.SeatSouth:
		return "SOUTH"
	case engine.SeatWest:
		return "WEST"
	case engine.SeatNorth:
		return "NORTH"
	default:
		return "UNKNOWN"
	}
}

func windToString(wind engine.Wind) string {
	switch wind {
	case engine.WindEast:
		return "EAST"
	case engine.WindSouth:
		return "SOUTH"
	case engine.WindWest:
		return "WEST"
	case engine.WindNorth:
		return "NORTH"
	default:
		return "UNKNOWN"
	}
}

func suitToString(suit engine.Suit) string {
	switch suit {
	case engine.Man:
		return "MAN"
	case engine.Pin:
		return "PIN"
	case engine.Sou:
		return "SOU"
	case engine.Ji:
		return "JI"
	default:
		return "UNKNOWN"
	}
}

func actionTypeToString(actionType engine.ActionType) string {
	switch actionType {
	case engine.ActionDraw:
		return "DRAW"
	case engine.ActionDiscard:
		return "DISCARD"
	case engine.ActionPon:
		return "PON"
	case engine.ActionChi:
		return "CHI"
	case engine.ActionKan:
		return "KAN"
	case engine.ActionRiichi:
		return "RIICHI"
	case engine.ActionTsumo:
		return "TSUMO"
	case engine.ActionRon:
		return "RON"
	default:
		return "UNKNOWN"
	}
}

func meldTypeToString(meldType engine.MeldType) string {
	switch meldType {
	case engine.MeldPon:
		return "PON"
	case engine.MeldChi:
		return "CHI"
	case engine.MeldAnkan:
		return "ANKAN"
	case engine.MeldMinkan:
		return "MINKAN"
	case engine.MeldKakan:
		return "KAKAN"
	default:
		return "UNKNOWN"
	}
}

// Conversion back to engine types

// convertToEngineGameState converts serializable format to engine.GameState
func convertToEngineGameState(state *GameState) (*engine.GameState, error) {
	// Convert players
	var players [4]*engine.Player
	for i := 0; i < 4; i++ {
		player, err := convertToEnginePlayer(&state.Players[i])
		if err != nil {
			return nil, fmt.Errorf("failed to convert player %d: %w", i, err)
		}
		players[i] = player
	}

	// Convert round
	round, err := convertToEngineRound(&state.Round)
	if err != nil {
		return nil, fmt.Errorf("failed to convert round: %w", err)
	}

	// Convert wall - recreate from seed
	wall := engine.NewWallWithConfig(state.Seed, state.Config.AkaDora)

	// Convert history
	history, err := convertToEngineHistory(state.History)
	if err != nil {
		return nil, fmt.Errorf("failed to convert history: %w", err)
	}

	// Convert current player
	currentPlayer, err := stringToSeat(state.CurrentPlayer)
	if err != nil {
		return nil, fmt.Errorf("failed to convert current player: %w", err)
	}

	return &engine.GameState{
		Round:         round,
		Players:       players,
		Wall:          wall,
		CurrentPlayer: currentPlayer,
		TurnNumber:    state.TurnNumber,
		History:       history,
		Config: engine.RuleConfig{
			RiichiMahjong:  state.Config.RiichiMahjong,
			StartingPoints: state.Config.StartingPoints,
			ReturnPoints:   state.Config.ReturnPoints,
			AkaDora:        state.Config.AkaDora,
			LocalYaku:      state.Config.LocalYaku,
		},
		IsGameOver: state.IsGameOver,
	}, nil
}

// convertToEnginePlayer converts serializable Player to engine.Player
func convertToEnginePlayer(player *Player) (*engine.Player, error) {
	seat, err := stringToSeat(player.Seat)
	if err != nil {
		return nil, err
	}

	hand, err := convertToEngineTiles(player.Hand)
	if err != nil {
		return nil, fmt.Errorf("failed to convert hand: %w", err)
	}

	discards, err := convertToEngineTiles(player.Discards)
	if err != nil {
		return nil, fmt.Errorf("failed to convert discards: %w", err)
	}

	melds, err := convertToEngineMelds(player.Melds)
	if err != nil {
		return nil, fmt.Errorf("failed to convert melds: %w", err)
	}

	return &engine.Player{
		Seat:     seat,
		Hand:     hand,
		Discards: discards,
		Melds:    melds,
		Score:    player.Score,
		IsRiichi: player.IsRiichi,
		Furiten: engine.FuritenState{
			DiscardFuriten:   player.Furiten.DiscardFuriten,
			TemporaryFuriten: player.Furiten.TemporaryFuriten,
			RiichiFuriten:    player.Furiten.RiichiFuriten,
		},
		IsTenpai: player.IsTenpai,
	}, nil
}

// convertToEngineRound converts serializable Round to engine.Round
func convertToEngineRound(round *Round) (*engine.Round, error) {
	dealer, err := stringToSeat(round.Dealer)
	if err != nil {
		return nil, err
	}

	wind, err := stringToWind(round.RoundWind)
	if err != nil {
		return nil, err
	}

	return &engine.Round{
		Dealer:       dealer,
		RoundWind:    wind,
		RoundNumber:  round.RoundNumber,
		Honba:        round.Honba,
		RiichiSticks: round.RiichiSticks,
	}, nil
}

// convertToEngineTiles converts slice of serializable Tile to engine.Tile
func convertToEngineTiles(tiles []Tile) ([]engine.Tile, error) {
	result := make([]engine.Tile, len(tiles))
	for i, t := range tiles {
		tile, err := convertToEngineTile(t)
		if err != nil {
			return nil, fmt.Errorf("failed to convert tile %d: %w", i, err)
		}
		result[i] = tile
	}
	return result, nil
}

// convertToEngineTile converts serializable Tile to engine.Tile
func convertToEngineTile(tile Tile) (engine.Tile, error) {
	suit, err := stringToSuit(tile.Suit)
	if err != nil {
		return engine.Tile{}, err
	}

	return engine.Tile{
		Suit:      suit,
		Rank:      tile.Rank,
		IsAkaDora: tile.IsAkaDora,
	}, nil
}

// convertToEngineHistory converts serializable action history to engine.Action slice
func convertToEngineHistory(history []Action) ([]engine.Action, error) {
	result := make([]engine.Action, len(history))
	for i, a := range history {
		action, err := convertToEngineAction(a)
		if err != nil {
			return nil, fmt.Errorf("failed to convert action %d: %w", i, err)
		}
		result[i] = action
	}
	return result, nil
}

// convertToEngineAction converts serializable Action to engine.Action
func convertToEngineAction(action Action) (engine.Action, error) {
	player, err := stringToSeat(action.Player)
	if err != nil {
		return engine.Action{}, err
	}

	actionType, err := stringToActionType(action.Type)
	if err != nil {
		return engine.Action{}, err
	}

	tiles, err := convertToEngineTiles(action.Tiles)
	if err != nil {
		return engine.Action{}, fmt.Errorf("failed to convert tiles: %w", err)
	}

	return engine.Action{
		Turn:      action.Turn,
		Player:    player,
		Type:      actionType,
		Tiles:     tiles,
		Timestamp: time.Unix(action.Timestamp, 0),
	}, nil
}

// convertToEngineMelds converts slice of serializable Meld to engine.Meld
func convertToEngineMelds(melds []Meld) ([]engine.Meld, error) {
	result := make([]engine.Meld, len(melds))
	for i, m := range melds {
		meld, err := convertToEngineMeld(m)
		if err != nil {
			return nil, fmt.Errorf("failed to convert meld %d: %w", i, err)
		}
		result[i] = meld
	}
	return result, nil
}

// convertToEngineMeld converts serializable Meld to engine.Meld
func convertToEngineMeld(meld Meld) (engine.Meld, error) {
	meldType, err := stringToMeldType(meld.Type)
	if err != nil {
		return engine.Meld{}, err
	}

	tiles, err := convertToEngineTiles(meld.Tiles)
	if err != nil {
		return engine.Meld{}, fmt.Errorf("failed to convert tiles: %w", err)
	}

	fromPlayer, err := stringToSeat(meld.FromPlayer)
	if err != nil {
		return engine.Meld{}, err
	}

	return engine.Meld{
		Type:       meldType,
		Tiles:      tiles,
		FromPlayer: fromPlayer,
	}, nil
}

// String to enum conversion helpers
func stringToSeat(s string) (engine.Seat, error) {
	switch s {
	case "EAST":
		return engine.SeatEast, nil
	case "SOUTH":
		return engine.SeatSouth, nil
	case "WEST":
		return engine.SeatWest, nil
	case "NORTH":
		return engine.SeatNorth, nil
	default:
		return 0, fmt.Errorf("invalid seat: %s", s)
	}
}

func stringToWind(s string) (engine.Wind, error) {
	switch s {
	case "EAST":
		return engine.WindEast, nil
	case "SOUTH":
		return engine.WindSouth, nil
	case "WEST":
		return engine.WindWest, nil
	case "NORTH":
		return engine.WindNorth, nil
	default:
		return 0, fmt.Errorf("invalid wind: %s", s)
	}
}

func stringToSuit(s string) (engine.Suit, error) {
	switch s {
	case "MAN":
		return engine.Man, nil
	case "PIN":
		return engine.Pin, nil
	case "SOU":
		return engine.Sou, nil
	case "JI":
		return engine.Ji, nil
	default:
		return 0, fmt.Errorf("invalid suit: %s", s)
	}
}

func stringToActionType(s string) (engine.ActionType, error) {
	switch s {
	case "DRAW":
		return engine.ActionDraw, nil
	case "DISCARD":
		return engine.ActionDiscard, nil
	case "PON":
		return engine.ActionPon, nil
	case "CHI":
		return engine.ActionChi, nil
	case "KAN":
		return engine.ActionKan, nil
	case "RIICHI":
		return engine.ActionRiichi, nil
	case "TSUMO":
		return engine.ActionTsumo, nil
	case "RON":
		return engine.ActionRon, nil
	default:
		return 0, fmt.Errorf("invalid action type: %s", s)
	}
}

func stringToMeldType(s string) (engine.MeldType, error) {
	switch s {
	case "PON":
		return engine.MeldPon, nil
	case "CHI":
		return engine.MeldChi, nil
	case "ANKAN":
		return engine.MeldAnkan, nil
	case "MINKAN":
		return engine.MeldMinkan, nil
	case "KAKAN":
		return engine.MeldKakan, nil
	default:
		return 0, fmt.Errorf("invalid meld type: %s", s)
	}
}
