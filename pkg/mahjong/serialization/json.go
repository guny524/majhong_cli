package serialization

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/guny524/majhong_cli/pkg/mahjong/engine"
)

// JSONGameState represents the game state in JSON format
type JSONGameState struct {
	Version string      `json:"version"`
	Seed    int64       `json:"seed"`
	Round   JSONRound   `json:"round"`
	Players []JSONPlayer `json:"players"`
	Wall    JSONWall    `json:"wall"`
	History []JSONAction `json:"history"`
	Config  JSONConfig  `json:"config"`
	Current JSONCurrent `json:"current,omitempty"`
}

// JSONRound represents round information
type JSONRound struct {
	Dealer           string   `json:"dealer"`
	RoundWind        string   `json:"roundWind"`
	RoundNumber      int      `json:"roundNumber"`
	Honba            int      `json:"honba"`
	RiichiSticks     int      `json:"riichiSticks"`
	DoraIndicators   []string `json:"doraIndicators"`
	UraDoraIndicators []string `json:"uraDoraIndicators"`
}

// JSONPlayer represents player state
type JSONPlayer struct {
	Seat      string      `json:"seat"`
	Hand      []string    `json:"hand"`
	Discards  []string    `json:"discards"`
	Melds     []JSONMeld  `json:"melds"`
	Score     int         `json:"score"`
	Riichi    bool        `json:"riichi"`
	Furiten   JSONFuriten `json:"furiten"`
	Tenpai    bool        `json:"tenpai"`
	DrawnTile *string     `json:"drawnTile,omitempty"`
}

// JSONMeld represents a meld
type JSONMeld struct {
	Type       string   `json:"type"`
	Tiles      []string `json:"tiles"`
	FromPlayer string   `json:"fromPlayer,omitempty"`
	Concealed  bool     `json:"concealed,omitempty"`
}

// JSONFuriten represents furiten state
type JSONFuriten struct {
	DiscardFuriten   bool `json:"discardFuriten"`
	TemporaryFuriten bool `json:"temporaryFuriten"`
	RiichiFuriten    bool `json:"riichiFuriten"`
}

// JSONWall represents wall state
type JSONWall struct {
	Remaining int    `json:"remaining"`
	Seed      int64  `json:"seed"`
	Hash      string `json:"hash"`
	SaltHash  string `json:"saltHash"`
}

// JSONAction represents a game action
type JSONAction struct {
	Turn      int      `json:"turn"`
	Player    string   `json:"player"`
	Action    string   `json:"action"`
	Tiles     []string `json:"tiles,omitempty"`
	Timestamp string   `json:"timestamp"`
}

// JSONConfig represents game configuration
type JSONConfig struct {
	RiichiMahjong  bool     `json:"riichiMahjong"`
	StartingPoints int      `json:"startingPoints"`
	ReturnPoints   int      `json:"returnPoints"`
	AkaDora        int      `json:"akaDora"`
	LocalYaku      []string `json:"localYaku"`
}

// JSONCurrent represents current game state (for queries)
type JSONCurrent struct {
	Player      string `json:"player,omitempty"`
	TurnNumber  int    `json:"turnNumber,omitempty"`
	IsGameOver  bool   `json:"isGameOver,omitempty"`
}

// ToJSON converts GameState to JSON format
func ToJSON(gs *engine.GameState) (*JSONGameState, error) {
	// Convert players
	players := make([]JSONPlayer, 4)
	for i, p := range gs.Players {
		hand := make([]string, len(p.Hand))
		for j, tile := range p.Hand {
			hand[j] = tile.String()
		}

		discards := make([]string, len(p.Discards))
		for j, tile := range p.Discards {
			discards[j] = tile.String()
		}

		melds := make([]JSONMeld, len(p.Melds))
		for j, meld := range p.Melds {
			tiles := make([]string, len(meld.Tiles))
			for k, tile := range meld.Tiles {
				tiles[k] = tile.String()
			}
			melds[j] = JSONMeld{
				Type:       meldTypeToString(meld.Type),
				Tiles:      tiles,
				FromPlayer: meld.FromPlayer.String(),
				Concealed:  meld.IsConcealed,
			}
		}

		var drawnTile *string
		if p.DrawnTile != nil {
			dt := p.DrawnTile.String()
			drawnTile = &dt
		}

		players[i] = JSONPlayer{
			Seat:      p.Seat.String(),
			Hand:      hand,
			Discards:  discards,
			Melds:     melds,
			Score:     p.Score,
			Riichi:    p.IsRiichi,
			Furiten: JSONFuriten{
				DiscardFuriten:   p.Furiten.DiscardFuriten,
				TemporaryFuriten: p.Furiten.TemporaryFuriten,
				RiichiFuriten:    p.Furiten.RiichiFuriten,
			},
			Tenpai:    false, // TODO: Calculate tenpai
			DrawnTile: drawnTile,
		}
	}

	// Convert dora indicators
	doraIndicators := make([]string, len(gs.Wall.DeadWall.DoraIndicators))
	for i, tile := range gs.Wall.DeadWall.DoraIndicators {
		doraIndicators[i] = tile.String()
	}

	uraDoraIndicators := make([]string, len(gs.Wall.DeadWall.UraDoraIndicators))
	for i, tile := range gs.Wall.DeadWall.UraDoraIndicators {
		uraDoraIndicators[i] = tile.String()
	}

	// Convert history
	history := make([]JSONAction, len(gs.History))
	for i, action := range gs.History {
		tiles := make([]string, len(action.Tiles))
		for j, tile := range action.Tiles {
			tiles[j] = tile.String()
		}
		history[i] = JSONAction{
			Turn:      action.Turn,
			Player:    action.Player.String(),
			Action:    actionTypeToString(action.Type),
			Tiles:     tiles,
			Timestamp: action.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &JSONGameState{
		Version: "1.0.0",
		Seed:    gs.Wall.Seed,
		Round: JSONRound{
			Dealer:            gs.Round.Dealer.String(),
			RoundWind:         gs.Round.RoundWind.String(),
			RoundNumber:       gs.Round.RoundNumber,
			Honba:             gs.Round.Honba,
			RiichiSticks:      gs.Round.RiichiSticks,
			DoraIndicators:    doraIndicators,
			UraDoraIndicators: uraDoraIndicators,
		},
		Players: players,
		Wall: JSONWall{
			Remaining: len(gs.Wall.LiveWall),
			Seed:      gs.Wall.Seed,
			Hash:      gs.Wall.Hash,
			SaltHash:  gs.Wall.SaltHash,
		},
		History: history,
		Config: JSONConfig{
			RiichiMahjong:  gs.Config.RiichiMahjong,
			StartingPoints: gs.Config.StartingPoints,
			ReturnPoints:   gs.Config.ReturnPoints,
			AkaDora:        akaDoraBoolToInt(gs.Config.AkaDora),
			LocalYaku:      gs.Config.LocalYaku,
		},
		Current: JSONCurrent{
			Player:     gs.CurrentPlayer.String(),
			TurnNumber: gs.TurnNumber,
			IsGameOver: gs.IsGameOver,
		},
	}, nil
}

// FromJSON converts JSON format to GameState
func FromJSON(jgs *JSONGameState) (*engine.GameState, error) {
	// Parse config
	config := engine.RuleConfig{
		RiichiMahjong:  jgs.Config.RiichiMahjong,
		StartingPoints: jgs.Config.StartingPoints,
		ReturnPoints:   jgs.Config.ReturnPoints,
		AkaDora:        jgs.Config.AkaDora > 0,
		LocalYaku:      jgs.Config.LocalYaku,
	}

	// Create wall from seed (for validation)
	wall := engine.NewWallWithConfig(jgs.Seed, config.AkaDora)

	// Validate wall hash
	if wall.Hash != jgs.Wall.Hash {
		return nil, fmt.Errorf("wall hash mismatch: expected %s, got %s", jgs.Wall.Hash, wall.Hash)
	}

	// Parse round
	dealer, err := parseSeat(jgs.Round.Dealer)
	if err != nil {
		return nil, fmt.Errorf("invalid dealer: %w", err)
	}

	roundWind, err := parseWind(jgs.Round.RoundWind)
	if err != nil {
		return nil, fmt.Errorf("invalid round wind: %w", err)
	}

	round := &engine.Round{
		Dealer:       dealer,
		RoundWind:    roundWind,
		RoundNumber:  jgs.Round.RoundNumber,
		Honba:        jgs.Round.Honba,
		RiichiSticks: jgs.Round.RiichiSticks,
	}

	// Parse players
	var players [4]*engine.Player
	for i, jp := range jgs.Players {
		seat, err := parseSeat(jp.Seat)
		if err != nil {
			return nil, fmt.Errorf("invalid seat for player %d: %w", i, err)
		}

		hand, err := engine.ParseTiles(jp.Hand)
		if err != nil {
			return nil, fmt.Errorf("invalid hand for player %d: %w", i, err)
		}

		discards, err := engine.ParseTiles(jp.Discards)
		if err != nil {
			return nil, fmt.Errorf("invalid discards for player %d: %w", i, err)
		}

		melds := make([]engine.Meld, len(jp.Melds))
		for j, jm := range jp.Melds {
			tiles, err := engine.ParseTiles(jm.Tiles)
			if err != nil {
				return nil, fmt.Errorf("invalid meld tiles for player %d, meld %d: %w", i, j, err)
			}

			meldType, err := parseMeldType(jm.Type)
			if err != nil {
				return nil, fmt.Errorf("invalid meld type for player %d, meld %d: %w", i, j, err)
			}

			fromPlayer := engine.SeatEast // default
			if jm.FromPlayer != "" {
				fromPlayer, err = parseSeat(jm.FromPlayer)
				if err != nil {
					return nil, fmt.Errorf("invalid from player for meld %d: %w", j, err)
				}
			}

			melds[j] = engine.Meld{
				Type:        meldType,
				Tiles:       tiles,
				FromPlayer:  fromPlayer,
				IsConcealed: jm.Concealed,
			}
		}

		player := engine.NewPlayer(seat)
		player.Hand = hand
		player.Discards = discards
		player.Melds = melds
		player.Score = jp.Score
		player.IsRiichi = jp.Riichi
		player.Furiten = engine.FuritenState{
			DiscardFuriten:   jp.Furiten.DiscardFuriten,
			TemporaryFuriten: jp.Furiten.TemporaryFuriten,
			RiichiFuriten:    jp.Furiten.RiichiFuriten,
		}

		// Restore drawn tile if present
		if jp.DrawnTile != nil {
			tile, err := engine.ParseTile(*jp.DrawnTile)
			if err != nil {
				return nil, fmt.Errorf("invalid drawn tile for player %d: %w", i, err)
			}
			player.DrawnTile = &tile
		}

		players[int(seat)] = player
	}

	// Parse history
	history := make([]engine.Action, len(jgs.History))
	for i, ja := range jgs.History {
		player, err := parseSeat(ja.Player)
		if err != nil {
			return nil, fmt.Errorf("invalid player for action %d: %w", i, err)
		}

		actionType, err := parseActionType(ja.Action)
		if err != nil {
			return nil, fmt.Errorf("invalid action type for action %d: %w", i, err)
		}

		tiles, err := engine.ParseTiles(ja.Tiles)
		if err != nil {
			return nil, fmt.Errorf("invalid tiles for action %d: %w", i, err)
		}

		history[i] = engine.Action{
			Type:   actionType,
			Player: player,
			Tiles:  tiles,
			Turn:   ja.Turn,
		}
	}

	// Parse current player
	currentPlayer, err := parseSeat(jgs.Current.Player)
	if err != nil {
		return nil, fmt.Errorf("invalid current player: %w", err)
	}

	// Reconstruct wall state by replaying history
	// We need to draw the same tiles that were drawn in the original game
	// For now, we'll trust the wall state and just update the remaining count
	drawnCount := (14 * 4) + len(jgs.History) // Initial deal + history actions
	for i := 0; i < drawnCount && i < len(wall.LiveWall); i++ {
		wall.LiveWall = wall.LiveWall[1:] // Remove drawn tiles
	}

	return &engine.GameState{
		Round:         round,
		Players:       players,
		Wall:          wall,
		CurrentPlayer: currentPlayer,
		TurnNumber:    jgs.Current.TurnNumber,
		History:       history,
		Config:        config,
		IsGameOver:    jgs.Current.IsGameOver,
	}, nil
}

// LoadFromFile loads game state from JSON file
func LoadFromFile(path string) (*engine.GameState, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var jgs JSONGameState
	if err := json.Unmarshal(data, &jgs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return FromJSON(&jgs)
}

// SaveToFile saves game state to JSON file (atomic write)
func SaveToFile(gs *engine.GameState, path string) error {
	jgs, err := ToJSON(gs)
	if err != nil {
		return fmt.Errorf("failed to convert to JSON: %w", err)
	}

	data, err := json.MarshalIndent(jgs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Atomic write: write to temp file, then rename
	dir := filepath.Dir(path)
	if dir == "" {
		dir = "."
	}

	tempFile, err := os.CreateTemp(dir, ".majhong_tmp_*.json")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tempPath := tempFile.Name()
	defer os.Remove(tempPath) // Clean up temp file if rename fails

	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// Helper functions for parsing

func parseSeat(s string) (engine.Seat, error) {
	switch s {
	case "East":
		return engine.SeatEast, nil
	case "South":
		return engine.SeatSouth, nil
	case "West":
		return engine.SeatWest, nil
	case "North":
		return engine.SeatNorth, nil
	default:
		return 0, fmt.Errorf("invalid seat: %s", s)
	}
}

func parseWind(s string) (engine.Wind, error) {
	switch s {
	case "East":
		return engine.WindEast, nil
	case "South":
		return engine.WindSouth, nil
	case "West":
		return engine.WindWest, nil
	case "North":
		return engine.WindNorth, nil
	default:
		return 0, fmt.Errorf("invalid wind: %s", s)
	}
}

func parseMeldType(s string) (engine.MeldType, error) {
	switch s {
	case "chi":
		return engine.MeldChi, nil
	case "pon":
		return engine.MeldPon, nil
	case "ankan":
		return engine.MeldAnkan, nil
	case "minkan":
		return engine.MeldMinkan, nil
	case "kakan":
		return engine.MeldKakan, nil
	default:
		return 0, fmt.Errorf("invalid meld type: %s", s)
	}
}

func parseActionType(s string) (engine.ActionType, error) {
	switch s {
	case "draw":
		return engine.ActionDraw, nil
	case "discard":
		return engine.ActionDiscard, nil
	case "pon":
		return engine.ActionPon, nil
	case "chi":
		return engine.ActionChi, nil
	case "kan":
		return engine.ActionKan, nil
	case "riichi":
		return engine.ActionRiichi, nil
	case "tsumo":
		return engine.ActionTsumo, nil
	case "ron":
		return engine.ActionRon, nil
	default:
		return 0, fmt.Errorf("invalid action type: %s", s)
	}
}

func meldTypeToString(mt engine.MeldType) string {
	switch mt {
	case engine.MeldChi:
		return "chi"
	case engine.MeldPon:
		return "pon"
	case engine.MeldAnkan:
		return "ankan"
	case engine.MeldMinkan:
		return "minkan"
	case engine.MeldKakan:
		return "kakan"
	default:
		return "unknown"
	}
}

func actionTypeToString(at engine.ActionType) string {
	switch at {
	case engine.ActionDraw:
		return "draw"
	case engine.ActionDiscard:
		return "discard"
	case engine.ActionPon:
		return "pon"
	case engine.ActionChi:
		return "chi"
	case engine.ActionKan:
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

func akaDoraBoolToInt(b bool) int {
	if b {
		return 3 // Standard 3 aka dora
	}
	return 0
}
