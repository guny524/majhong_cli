package engine

import (
	"fmt"
	"testing"
)

func TestNewWall(t *testing.T) {
	seed := int64(12345)
	wall := NewWall(seed)

	// Check total tiles
	allTiles := wall.GetAllTiles()
	if len(allTiles) != 136 {
		t.Errorf("Expected 136 tiles, got %d", len(allTiles))
	}

	// Check live wall
	if len(wall.LiveWall) != 122 {
		t.Errorf("Expected 122 live wall tiles, got %d", len(wall.LiveWall))
	}

	// Check dead wall components
	if len(wall.DeadWall.RinshanTiles) != 4 {
		t.Errorf("Expected 4 rinshan tiles, got %d", len(wall.DeadWall.RinshanTiles))
	}
	if len(wall.DeadWall.DoraIndicators) != 5 {
		t.Errorf("Expected 5 dora indicators, got %d", len(wall.DeadWall.DoraIndicators))
	}
	if len(wall.DeadWall.UraDoraIndicators) != 5 {
		t.Errorf("Expected 5 ura-dora indicators, got %d", len(wall.DeadWall.UraDoraIndicators))
	}

	// Check initial visible dora
	visibleDora := wall.GetVisibleDoraIndicators()
	if len(visibleDora) != 1 {
		t.Errorf("Expected 1 visible dora indicator initially, got %d", len(visibleDora))
	}
}

func TestWallTileComposition(t *testing.T) {
	wall := NewWall(123)
	allTiles := wall.GetAllTiles()

	// Count each tile type
	tileCounts := make(map[string]int)
	for _, tile := range allTiles {
		// For aka dora, count as regular 5
		key := tile.String()
		if tile.IsAkaDora {
			key = fmt.Sprintf("5%s", tile.Suit)
		}
		tileCounts[key]++
	}

	// Each tile type should appear 4 times
	for tileType, count := range tileCounts {
		if count != 4 {
			t.Errorf("Tile %s appears %d times, expected 4", tileType, count)
		}
	}

	// Should have 34 unique tile types
	if len(tileCounts) != 34 {
		t.Errorf("Expected 34 unique tile types, got %d", len(tileCounts))
	}
}

func TestWallDraw(t *testing.T) {
	wall := NewWall(456)
	initialCount := len(wall.LiveWall)

	// Draw one tile
	tile, err := wall.Draw()
	if err != nil {
		t.Fatalf("Draw failed: %v", err)
	}
	if tile.Suit == 0 && tile.Rank == 0 {
		t.Error("Drew empty tile")
	}

	// Check live wall decreased
	if len(wall.LiveWall) != initialCount-1 {
		t.Errorf("Expected live wall to decrease by 1, got %d tiles", len(wall.LiveWall))
	}
}

func TestWallDrawExhaustion(t *testing.T) {
	wall := NewWall(789)

	// Draw all 122 tiles
	for i := 0; i < 122; i++ {
		_, err := wall.Draw()
		if err != nil {
			t.Fatalf("Draw %d failed: %v", i+1, err)
		}
	}

	// Wall should be exhausted
	if !wall.IsExhausted() {
		t.Error("Wall should be exhausted after drawing 122 tiles")
	}

	// Next draw should fail
	_, err := wall.Draw()
	if err == nil {
		t.Error("Expected error when drawing from exhausted wall")
	}
}

func TestWallDrawRinshan(t *testing.T) {
	wall := NewWall(111)
	initialRinshan := len(wall.DeadWall.RinshanTiles)
	initialLive := len(wall.LiveWall)

	// Draw rinshan tile
	tile, err := wall.DrawRinshan()
	if err != nil {
		t.Fatalf("DrawRinshan failed: %v", err)
	}
	if tile.Suit == 0 && tile.Rank == 0 {
		t.Error("Drew empty rinshan tile")
	}

	// Rinshan should be replenished (still 4 tiles)
	if len(wall.DeadWall.RinshanTiles) != initialRinshan {
		t.Errorf("Expected %d rinshan tiles after draw (replenished), got %d",
			initialRinshan, len(wall.DeadWall.RinshanTiles))
	}

	// Live wall should decrease by 1 (used for replenishment)
	if len(wall.LiveWall) != initialLive-1 {
		t.Errorf("Expected live wall to decrease by 1, got %d tiles", len(wall.LiveWall))
	}

	// Total tiles in wall decreases by 1 (rinshan tile taken by player)
	// But dead wall is replenished, so dead wall stays at 14
	allTiles := wall.GetAllTiles()
	expectedTotal := 135 // 121 live + 14 dead
	if len(allTiles) != expectedTotal {
		t.Errorf("Expected %d total tiles after rinshan draw, got %d", expectedTotal, len(allTiles))
	}
}

func TestWallDoraIndicators(t *testing.T) {
	wall := NewWall(222)

	// Initial: 1 dora visible
	if len(wall.GetVisibleDoraIndicators()) != 1 {
		t.Errorf("Expected 1 visible dora initially, got %d", len(wall.GetVisibleDoraIndicators()))
	}

	// Flip next dora (after kan)
	err := wall.FlipNextDoraIndicator()
	if err != nil {
		t.Errorf("FlipNextDoraIndicator failed: %v", err)
	}
	if len(wall.GetVisibleDoraIndicators()) != 2 {
		t.Errorf("Expected 2 visible dora after flip, got %d", len(wall.GetVisibleDoraIndicators()))
	}

	// Flip 3 more times (total 5)
	for i := 0; i < 3; i++ {
		wall.FlipNextDoraIndicator()
	}
	if len(wall.GetVisibleDoraIndicators()) != 5 {
		t.Errorf("Expected 5 visible dora (max), got %d", len(wall.GetVisibleDoraIndicators()))
	}

	// Try to flip beyond 5 - should error
	err = wall.FlipNextDoraIndicator()
	if err == nil {
		t.Error("Expected error when flipping beyond 5 dora indicators")
	}
}

func TestWallUraDoraReveal(t *testing.T) {
	wall := NewWall(333)

	// Initially no ura-dora visible
	uraDora := wall.GetVisibleUraDoraIndicators()
	if len(uraDora) != 0 {
		t.Errorf("Expected 0 ura-dora initially, got %d", len(uraDora))
	}

	// After riichi win, reveal ura-dora
	revealed := wall.RevealUraDoraForRiichiWin()
	if len(revealed) != 1 {
		t.Errorf("Expected 1 ura-dora after riichi win, got %d", len(revealed))
	}

	// Flip more dora indicators
	wall.FlipNextDoraIndicator()
	wall.FlipNextDoraIndicator()

	// Reveal should match visible dora count
	revealed = wall.RevealUraDoraForRiichiWin()
	if len(revealed) != 3 {
		t.Errorf("Expected 3 ura-dora after 2 kans, got %d", len(revealed))
	}
}

func TestWallDeterminism(t *testing.T) {
	seed := int64(42)

	wall1 := NewWall(seed)
	wall2 := NewWall(seed)

	// Same seed should produce same wall
	for i := 0; i < 50; i++ {
		tile1, _ := wall1.Draw()
		tile2, _ := wall2.Draw()

		if !tile1.Equals(tile2) {
			t.Errorf("Draw %d mismatch: wall1=%s, wall2=%s", i+1, tile1, tile2)
		}
	}
}

func TestWallHashGeneration(t *testing.T) {
	wall := NewWall(555)

	// Hash should be set
	if wall.Hash == "" {
		t.Error("Wall hash not generated")
	}
	if wall.SaltHash == "" {
		t.Error("Salt hash not generated")
	}

	// Hash should be valid
	if !wall.ValidateHash(wall.Hash, wall.Salt) {
		t.Error("Generated hash failed validation")
	}
}

func TestWallHashValidation(t *testing.T) {
	wall := NewWall(666)
	originalHash := wall.Hash

	// Valid hash should pass
	if !wall.ValidateHash(originalHash, wall.Salt) {
		t.Error("Valid hash failed validation")
	}

	// Tamper with wall
	if len(wall.LiveWall) > 0 {
		wall.LiveWall[0] = NewTile(Man, 9)
	}

	// Tampered wall should fail validation
	if wall.ValidateHash(originalHash, wall.Salt) {
		t.Error("Tampered wall passed validation - hash check failed!")
	}
}

func TestWallIntegrityCheck(t *testing.T) {
	wall := NewWall(777)

	// Valid wall should pass integrity check
	err := wall.ValidateIntegrity(false)
	if err != nil {
		t.Errorf("Valid wall failed integrity check: %v", err)
	}

	// Tamper with wall
	originalTile := wall.LiveWall[0]
	wall.LiveWall[0] = NewTile(Man, 9)

	// With exitOnTamper=false: should warn but not error
	err = wall.ValidateIntegrity(false)
	if err != nil {
		t.Errorf("Integrity check with exitOnTamper=false should not error: %v", err)
	}

	// With exitOnTamper=true: should error
	err = wall.ValidateIntegrity(true)
	if err == nil {
		t.Error("Expected error with exitOnTamper=true for tampered wall")
	}

	// Restore
	wall.LiveWall[0] = originalTile
}

func TestWallCountTileType(t *testing.T) {
	wall := NewWall(888)

	// Each tile type should appear exactly 4 times
	testTile := NewTile(Man, 3)
	count := wall.CountTileType(testTile)
	if count != 4 {
		t.Errorf("Expected 4 occurrences of 3m, got %d", count)
	}

	// Test aka dora counting
	regular5m := NewTile(Man, 5)
	totalFives := wall.CountTileType(regular5m)
	if totalFives != 4 {
		t.Errorf("Expected 4 total 5m (including aka), got %d", totalFives)
	}
}

func TestWallWithoutAkaDora(t *testing.T) {
	wall := NewWallWithConfig(999, false)
	allTiles := wall.GetAllTiles()

	// Count aka dora tiles
	akaCount := 0
	for _, tile := range allTiles {
		if tile.IsAkaDora {
			akaCount++
		}
	}

	if akaCount != 0 {
		t.Errorf("Expected 0 aka dora tiles, got %d", akaCount)
	}
}

func TestWallWithAkaDora(t *testing.T) {
	wall := NewWallWithConfig(1000, true)
	allTiles := wall.GetAllTiles()

	// Count aka dora tiles
	akaCount := 0
	for _, tile := range allTiles {
		if tile.IsAkaDora {
			akaCount++
		}
	}

	// Should have 3 aka dora (0m, 0p, 0s)
	if akaCount != 3 {
		t.Errorf("Expected 3 aka dora tiles, got %d", akaCount)
	}
}
