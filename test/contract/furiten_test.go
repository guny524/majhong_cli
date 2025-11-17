package contract

import (
	"testing"
)

// Contract tests for furiten detection
// These tests MUST be written BEFORE implementation (TDD Red phase)
// Reference: constitution_majhong_rule.md (Furiten section)

// TestDiscardFuriten tests discard furiten (permanent for the hand)
func TestDiscardFuriten(t *testing.T) {
	tests := []struct {
		name            string
		handTiles       []string
		waitTiles       []string
		discards        []string
		expectedFuriten bool
		description     string
	}{
		{
			name:            "Single wait - tile in discards",
			handTiles:       []string{"2m", "3m", "4m", "5p", "6p", "7p", "1s", "2s", "3s", "4s", "5s", "6s", "8m"},
			waitTiles:       []string{"8m"},
			discards:        []string{"9m", "1p", "8m", "2z"},
			expectedFuriten: true,
			description:     "Waiting on 8m, but 8m is in discards → furiten",
		},
		{
			name:            "Single wait - tile NOT in discards",
			handTiles:       []string{"2m", "3m", "4m", "5p", "6p", "7p", "1s", "2s", "3s", "4s", "5s", "6s", "8m"},
			waitTiles:       []string{"8m"},
			discards:        []string{"9m", "1p", "7m", "2z"},
			expectedFuriten: false,
			description:     "Waiting on 8m, 8m not in discards → no furiten",
		},
		{
			name:            "Multi-wait ryanmen - one wait tile in discards",
			handTiles:       []string{"2m", "3m", "4m", "5p", "6p", "7p", "1s", "2s", "3s", "4s", "6s"},
			waitTiles:       []string{"5s", "7s"}, // 4-6s waiting for 5s or 7s
			discards:        []string{"9m", "1p", "5s", "2z"},
			expectedFuriten: true,
			description:     "Waiting on 5s/7s, 5s is in discards → furiten on ALL waits",
		},
		{
			name:            "Multi-wait ryanmen - can't ron on other wait",
			handTiles:       []string{"2m", "3m", "4m", "5p", "6p", "7p", "1s", "2s", "3s", "4s", "6s"},
			waitTiles:       []string{"5s", "7s"},
			discards:        []string{"9m", "1p", "5s", "2z"},
			expectedFuriten: true,
			description:     "Cannot ron on 7s either, even though 7s not in discards",
		},
		{
			name:            "Shanpon (double pon) wait - one wait in discards",
			handTiles:       []string{"1m", "2m", "3m", "4p", "5p", "6p", "1z", "1z", "2z", "2z"},
			waitTiles:       []string{"1z", "2z"}, // Waiting for 1z or 2z to complete
			discards:        []string{"9m", "1z", "5s"},
			expectedFuriten: true,
			description:     "Shanpon wait (東東 + 南南), 東 in discards → furiten on both",
		},
		{
			name:            "Three-sided wait - middle tile in discards",
			handTiles:       []string{"1m", "2m", "3m", "4p", "5p", "6p", "2s", "3s", "4s", "5s"},
			waitTiles:       []string{"1s", "4s", "6s"}, // 2-3-4-5 waiting for 1/4/6
			discards:        []string{"9m", "4s", "7p"},
			expectedFuriten: true,
			description:     "Three-sided wait, one wait tile in discards → furiten on all",
		},
		{
			name:            "Thirteen orphans wait - multiple terminals",
			handTiles:       []string{"1m", "9m", "1p", "9p", "1s", "9s", "1z", "2z", "3z", "4z", "5z", "6z", "7z"},
			waitTiles:       []string{"1m", "9m", "1p", "9p", "1s", "9s", "1z", "2z", "3z", "4z", "5z", "6z", "7z"},
			discards:        []string{"2m", "3p", "1m"}, // 1m in discards
			expectedFuriten: true,
			description:     "Kokushi 13-way wait, one terminal in discards → furiten",
		},
		{
			name:            "No discards yet - no furiten",
			handTiles:       []string{"2m", "3m", "4m", "5p", "6p", "7p", "1s", "2s", "3s", "4s", "5s", "6s", "8m"},
			waitTiles:       []string{"8m"},
			discards:        []string{},
			expectedFuriten: false,
			description:     "No discards yet → no discard furiten",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("NOT IMPLEMENTED - TDD Red Phase: Implement pkg/mahjong/rules/furiten.go")

			// TODO: Implement discard furiten check
			// furiten := CheckDiscardFuriten(tt.waitTiles, tt.discards)
			//
			// if furiten != tt.expectedFuriten {
			//     t.Errorf("%s\nExpected furiten=%v, got %v", tt.description, tt.expectedFuriten, furiten)
			// }
		})
	}
}

// TestTemporaryFuriten tests temporary furiten (cleared on next turn)
func TestTemporaryFuriten(t *testing.T) {
	tests := []struct {
		name            string
		waitTiles       []string
		skippedTile     string
		currentRound    int
		skipRound       int
		isRiichi        bool
		expectedFuriten bool
		description     string
	}{
		{
			name:            "Skip ron - same round",
			waitTiles:       []string{"3m"},
			skippedTile:     "3m",
			currentRound:    5,
			skipRound:       5,
			isRiichi:        false,
			expectedFuriten: true,
			description:     "Skipped ron on 3m, still in same turn cycle → temporary furiten",
		},
		{
			name:            "Skip ron - next turn arrived",
			waitTiles:       []string{"3m"},
			skippedTile:     "3m",
			currentRound:    6,
			skipRound:       5,
			isRiichi:        false,
			expectedFuriten: false,
			description:     "Skipped ron on 3m, but next turn arrived → furiten cleared",
		},
		{
			name:            "Not applicable during riichi",
			waitTiles:       []string{"3m"},
			skippedTile:     "3m",
			currentRound:    5,
			skipRound:       5,
			isRiichi:        true,
			expectedFuriten: false,
			description:     "Riichi uses different furiten rule (riichi furiten), not temporary",
		},
		{
			name:            "Multi-wait - skip one affects all",
			waitTiles:       []string{"3m", "6m"},
			skippedTile:     "3m",
			currentRound:    5,
			skipRound:       5,
			isRiichi:        false,
			expectedFuriten: true,
			description:     "Skipped 3m → cannot ron on 6m either until next turn",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("NOT IMPLEMENTED - TDD Red Phase: Implement temporary furiten")

			// TODO: Implement temporary furiten check
			// player := &Player{
			//     IsRiichi: tt.isRiichi,
			//     LastSkippedRonRound: tt.skipRound,
			// }
			// gameState := &GameState{
			//     CurrentRoundNumber: tt.currentRound,
			// }
			//
			// furiten := CheckTemporaryFuriten(player, gameState)
			//
			// if furiten != tt.expectedFuriten {
			//     t.Errorf("%s\nExpected furiten=%v, got %v", tt.description, tt.expectedFuriten, furiten)
			// }
		})
	}
}

// TestRiichiFuriten tests riichi furiten (permanent after skipping any win)
func TestRiichiFuriten(t *testing.T) {
	tests := []struct {
		name            string
		waitTiles       []string
		skippedTiles    []string // Tiles skipped after riichi (tsumo or ron)
		isRiichi        bool
		expectedFuriten bool
		description     string
	}{
		{
			name:            "Riichi - skip tsumo → permanent furiten",
			waitTiles:       []string{"3m", "6m"},
			skippedTiles:    []string{"3m"}, // Drew 3m, didn't declare tsumo
			isRiichi:        true,
			expectedFuriten: true,
			description:     "After riichi, drew 3m (winning tile) and didn't tsumo → permanent furiten",
		},
		{
			name:            "Riichi - skip ron → permanent furiten",
			waitTiles:       []string{"3m", "6m"},
			skippedTiles:    []string{"6m"}, // Opponent discarded 6m, didn't ron
			isRiichi:        true,
			expectedFuriten: true,
			description:     "After riichi, opponent discarded 6m and didn't ron → permanent furiten",
		},
		{
			name:            "Riichi - no skips → no riichi furiten",
			waitTiles:       []string{"3m", "6m"},
			skippedTiles:    []string{},
			isRiichi:        true,
			expectedFuriten: false,
			description:     "After riichi, haven't skipped any winning tiles → no riichi furiten",
		},
		{
			name:            "Not in riichi → riichi furiten doesn't apply",
			waitTiles:       []string{"3m", "6m"},
			skippedTiles:    []string{"3m"},
			isRiichi:        false,
			expectedFuriten: false,
			description:     "Not in riichi → riichi furiten rule doesn't apply",
		},
		{
			name:            "Riichi - skip multiple tiles",
			waitTiles:       []string{"3m", "6m"},
			skippedTiles:    []string{"3m", "6m"}, // Skipped both wait tiles
			isRiichi:        true,
			expectedFuriten: true,
			description:     "Skipped multiple winning tiles after riichi → permanent furiten",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("NOT IMPLEMENTED - TDD Red Phase: Implement riichi furiten")

			// TODO: Implement riichi furiten check
			// player := &Player{
			//     IsRiichi: tt.isRiichi,
			//     SkippedTilesAfterRiichi: tt.skippedTiles,
			// }
			//
			// furiten := CheckRiichiFuriten(player, tt.waitTiles)
			//
			// if furiten != tt.expectedFuriten {
			//     t.Errorf("%s\nExpected furiten=%v, got %v", tt.description, tt.expectedFuriten, furiten)
			// }
		})
	}
}

// TestCompositeFuriten tests combination of multiple furiten types
func TestCompositeFuriten(t *testing.T) {
	tests := []struct {
		name                 string
		waitTiles            []string
		discards             []string
		skippedThisRound     bool
		skippedAfterRiichi   bool
		isRiichi             bool
		expectedDiscardFuriten   bool
		expectedTemporaryFuriten bool
		expectedRiichiFuriten    bool
		canRon               bool
		canTsumo             bool
		description          string
	}{
		{
			name:                 "No furiten - can win both ways",
			waitTiles:            []string{"3m"},
			discards:             []string{"1p", "2s"},
			skippedThisRound:     false,
			skippedAfterRiichi:   false,
			isRiichi:             false,
			expectedDiscardFuriten:   false,
			expectedTemporaryFuriten: false,
			expectedRiichiFuriten:    false,
			canRon:               true,
			canTsumo:             true,
			description:          "No furiten conditions → can win by ron or tsumo",
		},
		{
			name:                 "Discard furiten only - can tsumo",
			waitTiles:            []string{"3m"},
			discards:             []string{"1p", "3m", "2s"},
			skippedThisRound:     false,
			skippedAfterRiichi:   false,
			isRiichi:             false,
			expectedDiscardFuriten:   true,
			expectedTemporaryFuriten: false,
			expectedRiichiFuriten:    false,
			canRon:               false,
			canTsumo:             true,
			description:          "Discard furiten → cannot ron, can tsumo",
		},
		{
			name:                 "Temporary furiten - cleared on next turn",
			waitTiles:            []string{"3m"},
			discards:             []string{"1p", "2s"},
			skippedThisRound:     true,
			skippedAfterRiichi:   false,
			isRiichi:             false,
			expectedDiscardFuriten:   false,
			expectedTemporaryFuriten: true,
			expectedRiichiFuriten:    false,
			canRon:               false,
			canTsumo:             true,
			description:          "Temporary furiten → cannot ron this round, can tsumo",
		},
		{
			name:                 "Riichi furiten - permanent",
			waitTiles:            []string{"3m"},
			discards:             []string{"1p", "2s"},
			skippedThisRound:     false,
			skippedAfterRiichi:   true,
			isRiichi:             true,
			expectedDiscardFuriten:   false,
			expectedTemporaryFuriten: false,
			expectedRiichiFuriten:    true,
			canRon:               false,
			canTsumo:             true,
			description:          "Riichi furiten → cannot ron permanently, can tsumo",
		},
		{
			name:                 "All three furiten types active",
			waitTiles:            []string{"3m"},
			discards:             []string{"1p", "3m", "2s"},
			skippedThisRound:     true,
			skippedAfterRiichi:   true,
			isRiichi:             true,
			expectedDiscardFuriten:   true,
			expectedTemporaryFuriten: true,
			expectedRiichiFuriten:    true,
			canRon:               false,
			canTsumo:             true,
			description:          "All furiten types active → definitely cannot ron, can tsumo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("NOT IMPLEMENTED - TDD Red Phase: Implement composite furiten")

			// TODO: Implement composite furiten check
			// player := &Player{
			//     Discards: tt.discards,
			//     IsRiichi: tt.isRiichi,
			//     SkippedThisRound: tt.skippedThisRound,
			//     SkippedAfterRiichi: tt.skippedAfterRiichi,
			// }
			//
			// discardFuriten := CheckDiscardFuriten(tt.waitTiles, tt.discards)
			// temporaryFuriten := CheckTemporaryFuriten(player, gameState)
			// riichiFuriten := CheckRiichiFuriten(player, tt.waitTiles)
			//
			// if discardFuriten != tt.expectedDiscardFuriten {
			//     t.Errorf("Discard furiten: expected %v, got %v", tt.expectedDiscardFuriten, discardFuriten)
			// }
			// if temporaryFuriten != tt.expectedTemporaryFuriten {
			//     t.Errorf("Temporary furiten: expected %v, got %v", tt.expectedTemporaryFuriten, temporaryFuriten)
			// }
			// if riichiFuriten != tt.expectedRiichiFuriten {
			//     t.Errorf("Riichi furiten: expected %v, got %v", tt.expectedRiichiFuriten, riichiFuriten)
			// }
			//
			// canRon := !discardFuriten && !temporaryFuriten && !riichiFuriten
			// if canRon != tt.canRon {
			//     t.Errorf("Can ron: expected %v, got %v", tt.canRon, canRon)
			// }
		})
	}
}

// TestFuritenRonPrevention tests that furiten prevents ron but allows tsumo
func TestFuritenRonPrevention(t *testing.T) {
	t.Run("Furiten player declares ron - should be rejected", func(t *testing.T) {
		t.Skip("NOT IMPLEMENTED - TDD Red Phase")

		// player := &Player{
		//     Hand: parseHand("2m 3m 4m 5p 6p 7p 1s 2s 3s 4s 5s 6s 8m"),
		//     Discards: []string{"8m"}, // Wait tile in discards
		// }
		// winTile := parseTile("8m")
		//
		// // Try to declare ron
		// err := DeclareRon(player, winTile)
		//
		// if err == nil {
		//     t.Error("Expected error when declaring ron in furiten, got nil")
		// }
		// if !strings.Contains(err.Error(), "furiten") {
		//     t.Errorf("Expected error message to mention furiten, got: %v", err)
		// }
	})

	t.Run("Furiten player declares tsumo - should be allowed", func(t *testing.T) {
		t.Skip("NOT IMPLEMENTED - TDD Red Phase")

		// player := &Player{
		//     Hand: parseHand("2m 3m 4m 5p 6p 7p 1s 2s 3s 4s 5s 6s 8m"),
		//     Discards: []string{"8m"}, // Wait tile in discards
		// }
		// drawnTile := parseTile("8m")
		//
		// // Try to declare tsumo
		// err := DeclareTsumo(player, drawnTile)
		//
		// if err != nil {
		//     t.Errorf("Tsumo should be allowed even in furiten, got error: %v", err)
		// }
	})
}

// TestFuritenEdgeCases tests edge cases and special scenarios
func TestFuritenEdgeCases(t *testing.T) {
	t.Run("Furiten riichi declaration - legal but risky", func(t *testing.T) {
		t.Skip("NOT IMPLEMENTED - TDD Red Phase")

		// Player can declare riichi even in furiten state (후리텐 리치)
		// But can only win by tsumo, not ron

		// player := &Player{
		//     Hand: parseHand("2m 3m 4m 5p 6p 7p 1s 2s 3s 4s 5s 6s 8m"),
		//     Discards: []string{"8m"}, // In furiten
		// }
		//
		// // Riichi declaration should succeed
		// err := DeclareRiichi(player)
		// if err != nil {
		//     t.Errorf("Furiten riichi should be legal, got error: %v", err)
		// }
		//
		// // But ron should still be blocked
		// err = DeclareRon(player, parseTile("8m"))
		// if err == nil {
		//     t.Error("Ron should still be blocked in furiten riichi")
		// }
	})

	t.Run("Changing wait after furiten - furiten persists", func(t *testing.T) {
		t.Skip("NOT IMPLEMENTED - TDD Red Phase")

		// If you discard a wait tile and later change your hand structure,
		// furiten still applies if new wait includes that tile

		// Initial: waiting on 3m (discard 3m → furiten)
		// Later: hand changes, now waiting on 3m and 6m
		// Result: Still in furiten because 3m is in discards
	})
}
