package engine

import (
	"fmt"
	"strconv"
	"strings"
)

// Suit represents the suit of a mahjong tile
type Suit int

const (
	Man Suit = iota // 萬子 (characters)
	Pin             // 筒子 (dots)
	Sou             // 索子 (bamboo)
	Ji              // 字牌 (honors)
)

func (s Suit) String() string {
	switch s {
	case Man:
		return "m"
	case Pin:
		return "p"
	case Sou:
		return "s"
	case Ji:
		return "z"
	default:
		return "?"
	}
}

// Tile represents a single mahjong tile
type Tile struct {
	Suit      Suit
	Rank      int  // 1-9 for suited tiles, 1-7 for honors
	IsAkaDora bool // Red five (0m, 0p, 0s)
}

// NewTile creates a new tile
func NewTile(suit Suit, rank int) Tile {
	return Tile{
		Suit:      suit,
		Rank:      rank,
		IsAkaDora: false,
	}
}

// NewAkaDoraTile creates a red five tile
func NewAkaDoraTile(suit Suit) Tile {
	if suit == Ji {
		panic("aka dora cannot be honor tiles")
	}
	return Tile{
		Suit:      suit,
		Rank:      5,
		IsAkaDora: true,
	}
}

// String returns the standard notation (e.g., "3m", "5p", "1z", "0m")
func (t Tile) String() string {
	if t.IsAkaDora {
		return fmt.Sprintf("0%s", t.Suit)
	}
	return fmt.Sprintf("%d%s", t.Rank, t.Suit)
}

// Equals checks if two tiles are identical (including aka dora status)
func (t Tile) Equals(other Tile) bool {
	return t.Suit == other.Suit &&
		t.Rank == other.Rank &&
		t.IsAkaDora == other.IsAkaDora
}

// SameType checks if two tiles are the same type (ignoring aka dora)
func (t Tile) SameType(other Tile) bool {
	return t.Suit == other.Suit && t.Rank == other.Rank
}

// IsTerminal checks if tile is a terminal (1 or 9)
func (t Tile) IsTerminal() bool {
	if t.Suit == Ji {
		return false
	}
	return t.Rank == 1 || t.Rank == 9
}

// IsHonor checks if tile is an honor tile
func (t Tile) IsHonor() bool {
	return t.Suit == Ji
}

// IsSimple checks if tile is a simple (2-8)
func (t Tile) IsSimple() bool {
	if t.Suit == Ji {
		return false
	}
	return t.Rank >= 2 && t.Rank <= 8
}

// IsWind checks if tile is a wind tile (East/South/West/North)
func (t Tile) IsWind() bool {
	return t.Suit == Ji && t.Rank >= 1 && t.Rank <= 4
}

// IsDragon checks if tile is a dragon tile (White/Green/Red)
func (t Tile) IsDragon() bool {
	return t.Suit == Ji && t.Rank >= 5 && t.Rank <= 7
}

// GetWindName returns the wind name for wind tiles
func (t Tile) GetWindName() string {
	if !t.IsWind() {
		return ""
	}
	winds := []string{"East", "South", "West", "North"}
	return winds[t.Rank-1]
}

// GetDragonName returns the dragon name for dragon tiles
func (t Tile) GetDragonName() string {
	if !t.IsDragon() {
		return ""
	}
	dragons := []string{"White", "Green", "Red"}
	return dragons[t.Rank-5]
}

// NextDora returns the next tile in sequence for dora calculation
// e.g., 5m indicator → 6m is dora, 9m indicator → 1m is dora
func (t Tile) NextDora() Tile {
	if t.Suit == Ji {
		// Honor tiles cycle within their group
		if t.IsWind() {
			// 1z→2z→3z→4z→1z (East→South→West→North→East)
			nextRank := (t.Rank % 4) + 1
			return NewTile(Ji, nextRank)
		}
		// 5z→6z→7z→5z (White→Green→Red→White)
		nextRank := ((t.Rank - 5) % 3) + 5
		return NewTile(Ji, nextRank)
	}

	// Number tiles: 1→2→...→9→1
	nextRank := (t.Rank % 9) + 1
	return NewTile(t.Suit, nextRank)
}

// ParseTile parses tile notation (e.g., "3m", "5p", "1z", "0m")
func ParseTile(notation string) (Tile, error) {
	notation = strings.TrimSpace(notation)
	if len(notation) < 2 {
		return Tile{}, fmt.Errorf("invalid tile notation: %s", notation)
	}

	// Extract suit (last character)
	suitChar := notation[len(notation)-1]
	var suit Suit
	switch suitChar {
	case 'm':
		suit = Man
	case 'p':
		suit = Pin
	case 's':
		suit = Sou
	case 'z':
		suit = Ji
	default:
		return Tile{}, fmt.Errorf("invalid suit: %c", suitChar)
	}

	// Extract rank
	rankStr := notation[:len(notation)-1]
	rank, err := strconv.Atoi(rankStr)
	if err != nil {
		return Tile{}, fmt.Errorf("invalid rank: %s", rankStr)
	}

	// Handle aka dora (0m, 0p, 0s)
	if rank == 0 {
		if suit == Ji {
			return Tile{}, fmt.Errorf("aka dora cannot be honor tiles")
		}
		return NewAkaDoraTile(suit), nil
	}

	// Validate rank ranges
	if suit == Ji {
		if rank < 1 || rank > 7 {
			return Tile{}, fmt.Errorf("honor tile rank must be 1-7, got: %d", rank)
		}
	} else {
		if rank < 1 || rank > 9 {
			return Tile{}, fmt.Errorf("number tile rank must be 1-9, got: %d", rank)
		}
	}

	return NewTile(suit, rank), nil
}

// ParseTiles parses multiple tile notations
func ParseTiles(notations []string) ([]Tile, error) {
	tiles := make([]Tile, 0, len(notations))
	for _, notation := range notations {
		tile, err := ParseTile(notation)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", notation, err)
		}
		tiles = append(tiles, tile)
	}
	return tiles, nil
}

// TilesToString converts tiles to notation strings
func TilesToString(tiles []Tile) []string {
	strs := make([]string, len(tiles))
	for i, tile := range tiles {
		strs[i] = tile.String()
	}
	return strs
}
