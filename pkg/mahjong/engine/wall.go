package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
)

// Wall represents the mahjong wall (136 tiles total)
type Wall struct {
	LiveWall  []Tile // 122 tiles initially
	DeadWall  DeadWall
	Seed      int64
	Salt      string
	Hash      string // SHA-256(wall + salt)
	SaltHash  string // SHA-256(salt)
	rng       *rand.Rand
}

// DeadWall represents the 14-tile dead wall (王牌)
type DeadWall struct {
	RinshanTiles      []Tile // 4 rinshan tiles (嶺上牌)
	DoraIndicators    []Tile // Up to 5 dora indicators
	UraDoraIndicators []Tile // Up to 5 ura-dora indicators (revealed after riichi win)
	visibleDoraCount  int    // Number of dora indicators currently visible
}

// NewWall creates a new shuffled wall with given seed
func NewWall(seed int64) *Wall {
	return NewWallWithConfig(seed, true)
}

// NewWallWithConfig creates a new wall with optional aka dora
func NewWallWithConfig(seed int64, akaDora bool) *Wall {
	// Create all 136 tiles
	tiles := generateAllTiles(akaDora)

	// Shuffle using Fisher-Yates with seeded RNG
	rng := rand.New(rand.NewSource(seed))
	shuffle(tiles, rng)

	// Split into live wall (122) and dead wall (14)
	liveWall := tiles[:122]
	deadWallTiles := tiles[122:]

	// Dead wall structure:
	// [0-3]: Rinshan tiles (4 tiles)
	// [4-8]: Dora indicators (5 tiles)
	// [9-13]: Ura-dora indicators (5 tiles)
	deadWall := DeadWall{
		RinshanTiles:      deadWallTiles[0:4],
		DoraIndicators:    deadWallTiles[4:9],
		UraDoraIndicators: deadWallTiles[9:14],
		visibleDoraCount:  1, // Initial dora indicator visible
	}

	wall := &Wall{
		LiveWall: liveWall,
		DeadWall: deadWall,
		Seed:     seed,
		rng:      rng,
	}

	// Generate hash
	salt := generateSalt(seed)
	wall.Salt = salt
	wall.Hash = wall.GenerateHash(salt)
	wall.SaltHash = hashString(salt)

	return wall
}

// generateAllTiles creates all 136 mahjong tiles
func generateAllTiles(akaDora bool) []Tile {
	tiles := make([]Tile, 0, 136)

	// Number tiles: Man/Pin/Sou (1-9), 4 of each = 108 tiles
	suits := []Suit{Man, Pin, Sou}
	for _, suit := range suits {
		for rank := 1; rank <= 9; rank++ {
			// Add 4 of each tile
			for count := 0; count < 4; count++ {
				// Handle aka dora (red fives)
				if akaDora && rank == 5 && count == 0 {
					tiles = append(tiles, NewAkaDoraTile(suit))
				} else {
					tiles = append(tiles, NewTile(suit, rank))
				}
			}
		}
	}

	// Honor tiles: Winds (1z-4z) + Dragons (5z-7z), 4 of each = 28 tiles
	for rank := 1; rank <= 7; rank++ {
		for count := 0; count < 4; count++ {
			tiles = append(tiles, NewTile(Ji, rank))
		}
	}

	return tiles
}

// shuffle performs Fisher-Yates shuffle
func shuffle(tiles []Tile, rng *rand.Rand) {
	for i := len(tiles) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		tiles[i], tiles[j] = tiles[j], tiles[i]
	}
}

// Draw draws one tile from the live wall
func (w *Wall) Draw() (Tile, error) {
	if len(w.LiveWall) == 0 {
		return Tile{}, fmt.Errorf("wall exhausted (ryuukyoku)")
	}

	// Draw from front of live wall
	tile := w.LiveWall[0]
	w.LiveWall = w.LiveWall[1:]
	return tile, nil
}

// DrawRinshan draws a rinshan tile (after kan) and replenishes dead wall
func (w *Wall) DrawRinshan() (Tile, error) {
	if len(w.DeadWall.RinshanTiles) == 0 {
		return Tile{}, fmt.Errorf("no rinshan tiles remaining")
	}

	// Draw from rinshan
	tile := w.DeadWall.RinshanTiles[0]
	w.DeadWall.RinshanTiles = w.DeadWall.RinshanTiles[1:]

	// Replenish dead wall from live wall end
	if len(w.LiveWall) > 0 {
		lastTile := w.LiveWall[len(w.LiveWall)-1]
		w.LiveWall = w.LiveWall[:len(w.LiveWall)-1]
		w.DeadWall.RinshanTiles = append(w.DeadWall.RinshanTiles, lastTile)
	}

	return tile, nil
}

// FlipNextDoraIndicator flips the next dora indicator (called after kan)
func (w *Wall) FlipNextDoraIndicator() error {
	if w.DeadWall.visibleDoraCount >= 5 {
		return fmt.Errorf("all dora indicators already visible")
	}
	w.DeadWall.visibleDoraCount++
	return nil
}

// GetVisibleDoraIndicators returns currently visible dora indicators
func (w *Wall) GetVisibleDoraIndicators() []Tile {
	return w.DeadWall.DoraIndicators[:w.DeadWall.visibleDoraCount]
}

// GetVisibleUraDoraIndicators returns ura-dora indicators (only after riichi win)
func (w *Wall) GetVisibleUraDoraIndicators() []Tile {
	// Ura-dora not visible by default
	return []Tile{}
}

// RevealUraDoraForRiichiWin reveals ura-dora indicators after riichi win
func (w *Wall) RevealUraDoraForRiichiWin() []Tile {
	return w.DeadWall.UraDoraIndicators[:w.DeadWall.visibleDoraCount]
}

// Remaining returns number of tiles left in live wall
func (w *Wall) Remaining() int {
	return len(w.LiveWall)
}

// IsExhausted checks if wall is exhausted (ryuukyoku condition)
func (w *Wall) IsExhausted() bool {
	return len(w.LiveWall) == 0
}

// GetAllTiles returns all tiles in the wall (for testing/validation)
func (w *Wall) GetAllTiles() []Tile {
	tiles := make([]Tile, 0, 136)
	tiles = append(tiles, w.LiveWall...)
	tiles = append(tiles, w.DeadWall.RinshanTiles...)
	tiles = append(tiles, w.DeadWall.DoraIndicators...)
	tiles = append(tiles, w.DeadWall.UraDoraIndicators...)
	return tiles
}

// GenerateHash generates SHA-256 hash of wall + salt
func (w *Wall) GenerateHash(salt string) string {
	allTiles := w.GetAllTiles()
	data := fmt.Sprintf("%v%s", allTiles, salt)
	return hashString(data)
}

// ValidateHash validates wall integrity using stored hash
func (w *Wall) ValidateHash(hash string, salt string) bool {
	computed := w.GenerateHash(salt)
	return computed == hash
}

// ValidateIntegrity validates wall and returns error if tampered
func (w *Wall) ValidateIntegrity(exitOnTamper bool) error {
	if !w.ValidateHash(w.Hash, w.Salt) {
		if exitOnTamper {
			return fmt.Errorf("wall tampered: hash validation failed")
		}
		// Warning only
		fmt.Println("WARNING: Wall hash validation failed (possible tampering)")
	}

	// Validate salt hash
	if hashString(w.Salt) != w.SaltHash {
		if exitOnTamper {
			return fmt.Errorf("salt tampered: hash validation failed")
		}
		fmt.Println("WARNING: Salt hash validation failed")
	}

	return nil
}

// CountTileType counts occurrences of a tile type in the wall
func (w *Wall) CountTileType(tile Tile) int {
	count := 0
	for _, t := range w.GetAllTiles() {
		if t.SameType(tile) {
			count++
		}
	}
	return count
}

// Helper functions

func hashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}

func generateSalt(seed int64) string {
	return fmt.Sprintf("majhong_salt_%d", seed)
}

// DoraIndicators returns all currently visible dora indicators
func (w *Wall) DoraIndicators() []Tile {
	return w.GetVisibleDoraIndicators()
}

// IsDoraForIndicator checks if a tile is dora for given indicator
func (w *Wall) IsDoraForIndicator(tile Tile, indicator Tile) bool {
	// Get the next tile in sequence after indicator
	nextTile := indicator.NextDora()
	return tile.Equals(nextTile)
}
