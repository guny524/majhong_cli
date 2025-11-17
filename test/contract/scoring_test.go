package contract

import (
	"testing"

	"github.com/guny524/majhong_cli/pkg/mahjong/rules"
)

// Contract tests for scoring calculation
// These tests MUST be written BEFORE implementation (TDD Red phase)
// Reference: constitution_majhong_score.md

// TestScoringCalculation tests han/fu calculation and scoring tables
func TestScoringCalculation(t *testing.T) {
	tests := []struct {
		name          string
		handTiles     []string
		winTile       string
		yaku          []string
		han           int
		fu            int
		isDealer      bool
		isTsumo       bool
		expectedRon   int
		expectedTsumo []int // [from non-dealer, from dealer] or [from each] for dealer
	}{
		// Basic scoring table tests
		{
			name:          "Non-dealer 1 han 30 fu - Ron",
			yaku:          []string{"TANYAO"},
			han:           1,
			fu:            30,
			isDealer:      false,
			isTsumo:       false,
			expectedRon:   1000,
			expectedTsumo: nil,
		},
		{
			name:          "Non-dealer 1 han 30 fu - Tsumo",
			yaku:          []string{"MENZEN_TSUMO"},
			han:           1,
			fu:            30,
			isDealer:      false,
			isTsumo:       true,
			expectedRon:   0,
			expectedTsumo: []int{300, 500}, // 300 from each non-dealer, 500 from dealer
		},
		{
			name:          "Dealer 1 han 30 fu - Ron",
			yaku:          []string{"TANYAO"},
			han:           1,
			fu:            30,
			isDealer:      true,
			isTsumo:       false,
			expectedRon:   1500,
			expectedTsumo: nil,
		},
		{
			name:          "Dealer 1 han 30 fu - Tsumo",
			yaku:          []string{"MENZEN_TSUMO"},
			han:           1,
			fu:            30,
			isDealer:      true,
			isTsumo:       true,
			expectedRon:   0,
			expectedTsumo: []int{500}, // 500 from each player
		},

		// 2 han tests
		{
			name:          "Non-dealer 2 han 30 fu - Ron",
			yaku:          []string{"TANYAO", "PINFU"},
			han:           2,
			fu:            30,
			isDealer:      false,
			isTsumo:       false,
			expectedRon:   2000,
			expectedTsumo: nil,
		},
		{
			name:          "Non-dealer 2 han 40 fu - Ron",
			yaku:          []string{"TOITOI"},
			han:           2,
			fu:            40,
			isDealer:      false,
			isTsumo:       false,
			expectedRon:   2600,
			expectedTsumo: nil,
		},

		// 3 han tests
		{
			name:          "Non-dealer 3 han 30 fu - Ron",
			yaku:          []string{"HONITSU"},
			han:           3,
			fu:            30,
			isDealer:      false,
			isTsumo:       false,
			expectedRon:   3900,
			expectedTsumo: nil,
		},
		{
			name:          "Non-dealer 3 han 60 fu - Ron",
			yaku:          []string{"HONITSU"},
			han:           3,
			fu:            60,
			isDealer:      false,
			isTsumo:       false,
			expectedRon:   7700,
			expectedTsumo: nil,
		},

		// 4 han tests (near mangan)
		{
			name:          "Non-dealer 4 han 30 fu - Ron",
			yaku:          []string{"RIICHI", "HONITSU"},
			han:           4,
			fu:            30,
			isDealer:      false,
			isTsumo:       false,
			expectedRon:   7700,
			expectedTsumo: nil,
		},
		{
			name:          "Non-dealer 4 han 40 fu - Mangan",
			yaku:          []string{"RIICHI", "HONITSU"},
			han:           4,
			fu:            40,
			isDealer:      false,
			isTsumo:       false,
			expectedRon:   8000, // Mangan
			expectedTsumo: nil,
		},

		// Limit hands
		{
			name:          "Non-dealer Mangan (5 han) - Ron",
			yaku:          []string{"CHINITSU"},
			han:           5,
			fu:            0, // Fu irrelevant for mangan
			isDealer:      false,
			isTsumo:       false,
			expectedRon:   8000,
			expectedTsumo: nil,
		},
		{
			name:          "Non-dealer Mangan (5 han) - Tsumo",
			yaku:          []string{"CHINITSU"},
			han:           5,
			fu:            0,
			isDealer:      false,
			isTsumo:       true,
			expectedRon:   0,
			expectedTsumo: []int{2000, 4000},
		},
		{
			name:          "Dealer Mangan - Ron",
			yaku:          []string{"CHINITSU"},
			han:           5,
			fu:            0,
			isDealer:      true,
			isTsumo:       false,
			expectedRon:   12000,
			expectedTsumo: nil,
		},
		{
			name:          "Dealer Mangan - Tsumo",
			yaku:          []string{"CHINITSU"},
			han:           5,
			fu:            0,
			isDealer:      true,
			isTsumo:       true,
			expectedRon:   0,
			expectedTsumo: []int{4000},
		},
		{
			name:          "Non-dealer Haneman (6 han) - Ron",
			han:           6,
			fu:            0,
			isDealer:      false,
			isTsumo:       false,
			expectedRon:   12000,
			expectedTsumo: nil,
		},
		{
			name:          "Non-dealer Haneman (7 han) - Tsumo",
			han:           7,
			fu:            0,
			isDealer:      false,
			isTsumo:       true,
			expectedRon:   0,
			expectedTsumo: []int{3000, 6000},
		},
		{
			name:          "Non-dealer Baiman (8 han) - Ron",
			han:           8,
			fu:            0,
			isDealer:      false,
			isTsumo:       false,
			expectedRon:   16000,
			expectedTsumo: nil,
		},
		{
			name:          "Non-dealer Baiman (10 han) - Tsumo",
			han:           10,
			fu:            0,
			isDealer:      false,
			isTsumo:       true,
			expectedRon:   0,
			expectedTsumo: []int{4000, 8000},
		},
		{
			name:          "Non-dealer Sanbaiman (11 han) - Ron",
			han:           11,
			fu:            0,
			isDealer:      false,
			isTsumo:       false,
			expectedRon:   24000,
			expectedTsumo: nil,
		},
		{
			name:          "Non-dealer Sanbaiman (12 han) - Tsumo",
			han:           12,
			fu:            0,
			isDealer:      false,
			isTsumo:       true,
			expectedRon:   0,
			expectedTsumo: []int{6000, 12000},
		},
		{
			name:          "Non-dealer Yakuman (13 han) - Ron",
			yaku:          []string{"KOKUSHI"},
			han:           13,
			fu:            0,
			isDealer:      false,
			isTsumo:       false,
			expectedRon:   32000,
			expectedTsumo: nil,
		},
		{
			name:          "Non-dealer Yakuman - Tsumo",
			yaku:          []string{"SUUANKOU"},
			han:           13,
			fu:            0,
			isDealer:      false,
			isTsumo:       true,
			expectedRon:   0,
			expectedTsumo: []int{8000, 16000},
		},
		{
			name:          "Dealer Yakuman - Ron",
			yaku:          []string{"DAISANGEN"},
			han:           13,
			fu:            0,
			isDealer:      true,
			isTsumo:       false,
			expectedRon:   48000,
			expectedTsumo: nil,
		},
		{
			name:          "Dealer Yakuman - Tsumo",
			yaku:          []string{"TSUUIISOU"},
			han:           13,
			fu:            0,
			isDealer:      true,
			isTsumo:       true,
			expectedRon:   0,
			expectedTsumo: []int{16000},
		},
		{
			name:          "Non-dealer Double Yakuman - Ron",
			yaku:          []string{"DAISUUSHII"},
			han:           26,
			fu:            0,
			isDealer:      false,
			isTsumo:       false,
			expectedRon:   64000,
			expectedTsumo: nil,
		},
		{
			name:          "Dealer Double Yakuman - Tsumo",
			yaku:          []string{"DAISUUSHII"},
			han:           26,
			fu:            0,
			isDealer:      true,
			isTsumo:       true,
			expectedRon:   0,
			expectedTsumo: []int{32000},
		},

		// Chiitoitsu special case (always 25 fu)
		{
			name:          "Chiitoitsu 2 han 25 fu - Ron",
			yaku:          []string{"CHIITOITSU"},
			han:           2,
			fu:            25,
			isDealer:      false,
			isTsumo:       false,
			expectedRon:   1600,
			expectedTsumo: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := rules.CalculateScore(tt.han, tt.fu, tt.isDealer, tt.isTsumo)

			if tt.isTsumo {
				if len(score.TsumoPayment) != len(tt.expectedTsumo) {
					t.Errorf("Expected tsumo payment length %d, got %d", len(tt.expectedTsumo), len(score.TsumoPayment))
					return
				}
				for i, expected := range tt.expectedTsumo {
					if score.TsumoPayment[i] != expected {
						t.Errorf("Tsumo payment[%d]: expected %d, got %d", i, expected, score.TsumoPayment[i])
					}
				}
			} else {
				if score.RonPayment != tt.expectedRon {
					t.Errorf("Expected ron %d, got %d", tt.expectedRon, score.RonPayment)
				}
			}
		})
	}
}

// TestFuCalculation tests fu (base points) calculation
func TestFuCalculation(t *testing.T) {
	tests := []struct {
		name        string
		handTiles   []string
		winTile     string
		melds       []string // "pon:1m", "chi:123m", "ankan:5p"
		pair        string
		waitType    string // "ryanmen", "kanchan", "penchan", "tanki"
		isTsumo     bool
		isClosed    bool
		seatWind    string
		roundWind   string
		expectedFu  int
	}{
		{
			name:       "Base fu only (open pinfu ron)",
			handTiles:  []string{"2m", "3m", "4m", "5p", "6p", "7p", "1s", "2s", "3s"},
			melds:      []string{"chi:456s"},
			pair:       "8m",
			waitType:   "ryanmen",
			isTsumo:    false,
			isClosed:   false,
			expectedFu: 30, // 20 base + 10 ron
		},
		{
			name:       "Pinfu tsumo special case",
			handTiles:  []string{"2m", "3m", "4m", "5p", "6p", "7p", "1s", "2s", "3s", "4s", "5s", "6s", "8m"},
			pair:       "8m",
			waitType:   "ryanmen",
			isTsumo:    true,
			isClosed:   true,
			expectedFu: 20, // Pinfu tsumo = 20 fu (no rounding)
		},
		{
			name:       "Chiitoitsu fixed 25 fu",
			handTiles:  []string{"1m", "1m", "3p", "3p", "5s", "5s", "2z", "2z", "4z", "4z", "6z", "6z", "7p"},
			pair:       "7p", // Actually all pairs
			isTsumo:    false,
			isClosed:   true,
			expectedFu: 25, // Fixed 25 fu
		},
		{
			name:       "Closed triplets - simples",
			handTiles:  []string{"2m", "2m", "2m", "5p", "5p", "5p", "8s", "8s", "8s"},
			melds:      []string{"anko:2m", "anko:5p", "anko:8s"},
			pair:       "3z",
			waitType:   "tanki",
			isTsumo:    false,
			isClosed:   true,
			expectedFu: 50, // 20 + 10(ron) + 4+4+4 (3 closed simple triplets) + 2(tanki) = 44 → 50
		},
		{
			name:       "Closed triplets - terminals/honors",
			handTiles:  []string{"1m", "1m", "1m", "9p", "9p", "9p", "5z", "5z", "5z"},
			melds:      []string{"anko:1m", "anko:9p", "anko:5z"},
			pair:       "2z",
			waitType:   "tanki",
			isTsumo:    false,
			isClosed:   true,
			expectedFu: 60, // 20 + 10(ron) + 8+8+8 (3 closed terminal/honor triplets) + 2(tanki) = 56 → 60
		},
		{
			name:       "Open triplets - simples",
			handTiles:  []string{"2m", "3m", "4m"},
			melds:      []string{"pon:5p", "pon:8s", "pon:3m"},
			pair:       "1z",
			waitType:   "tanki",
			isTsumo:    false,
			isClosed:   false,
			expectedFu: 40, // 20 + 10(ron) + 2+2+2 (3 open simple triplets) + 2(tanki) = 38 → 40
		},
		{
			name:       "Value pair (dragon)",
			handTiles:  []string{"2m", "3m", "4m", "5p", "6p", "7p", "1s", "2s", "3s", "4s", "5s", "6s"},
			pair:       "5z", // White dragon
			waitType:   "ryanmen",
			isTsumo:    false,
			isClosed:   true,
			expectedFu: 40, // 20 + 10(ron) + 2(dragon pair) = 32 → 40
		},
		{
			name:       "Double wind pair",
			handTiles:  []string{"2m", "3m", "4m", "5p", "6p", "7p", "1s", "2s", "3s", "4s", "5s", "6s"},
			pair:       "1z", // East
			waitType:   "ryanmen",
			isTsumo:    false,
			isClosed:   true,
			seatWind:   "East",
			roundWind:  "East",
			expectedFu: 40, // 20 + 10(ron) + 4(double east) = 34 → 40
		},
		{
			name:       "Kanchan wait",
			handTiles:  []string{"2m", "3m", "4m", "5p", "6p", "7p", "1s", "2s", "3s", "4s", "6s"},
			pair:       "8m",
			waitType:   "kanchan", // 4-6 waiting for 5
			isTsumo:    false,
			isClosed:   true,
			expectedFu: 40, // 20 + 10(ron) + 2(kanchan) = 32 → 40
		},
		{
			name:       "Penchan wait",
			handTiles:  []string{"2m", "3m", "4m", "5p", "6p", "7p", "1s", "2s", "3s", "1m", "2m"},
			pair:       "8m",
			waitType:   "penchan", // 1-2 waiting for 3
			isTsumo:    false,
			isClosed:   true,
			expectedFu: 40, // 20 + 10(ron) + 2(penchan) = 32 → 40
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("NOT IMPLEMENTED - TDD Red Phase: Implement fu calculation")

			// TODO: Implement fu calculation
			// fu := CalculateFu(hand, winTile, waitType, isTsumo, isClosed, seatWind, roundWind)
			//
			// if fu != tt.expectedFu {
			//     t.Errorf("Expected %d fu, got %d fu", tt.expectedFu, fu)
			// }
		})
	}
}

// TestHonbaAndRiichiSticks tests honba counters and riichi stick payments
func TestHonbaAndRiichiSticks(t *testing.T) {
	tests := []struct {
		name           string
		basePoints     int
		honbaCount     int
		riichiSticks   int
		isDealer       bool
		isTsumo        bool
		expectedTotal  int
		expectedTsumo  []int
	}{
		{
			name:          "Ron with 2 honba, 1 riichi stick",
			basePoints:    8000, // Mangan
			honbaCount:    2,
			riichiSticks:  1,
			isDealer:      false,
			isTsumo:       false,
			expectedTotal: 8000 + 600 + 1000, // Base + (300*2 honba) + 1000 riichi
		},
		{
			name:          "Tsumo with honba",
			basePoints:    8000, // Mangan
			honbaCount:    3,
			riichiSticks:  2,
			isDealer:      false,
			isTsumo:       true,
			expectedTotal: 8000 + 900 + 2000, // Base + (300*3) + (1000*2 riichi)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("NOT IMPLEMENTED - TDD Red Phase: Implement honba/riichi payments")
		})
	}
}

// TestPaoResponsibility tests pao (responsibility payment) rules
func TestPaoResponsibility(t *testing.T) {
	tests := []struct {
		name             string
		yakuman          string
		paoPlayer        string // Who fed the dangerous tile
		winMethod        string // "tsumo" or "ron"
		expectedPayments map[string]int
	}{
		{
			name:      "Daisangen pao - Tsumo",
			yakuman:   "DAISANGEN",
			paoPlayer: "South",
			winMethod: "tsumo",
			expectedPayments: map[string]int{
				"South": 32000, // Pao player pays full yakuman
			},
		},
		{
			name:      "Daisangen pao - Ron from pao player",
			yakuman:   "DAISANGEN",
			paoPlayer: "South",
			winMethod: "ron",
			expectedPayments: map[string]int{
				"South": 32000, // Pao player pays everything
			},
		},
		{
			name:      "Daisuushii pao - Tsumo",
			yakuman:   "DAISUUSHII",
			paoPlayer: "West",
			winMethod: "tsumo",
			expectedPayments: map[string]int{
				"West": 64000, // Double yakuman pao
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("NOT IMPLEMENTED - TDD Red Phase: Implement pao responsibility")
		})
	}
}

// TestUmaAndOka tests final score adjustments (placement bonuses)
func TestUmaAndOka(t *testing.T) {
	tests := []struct {
		name           string
		finalScores    []int // Raw scores for 4 players
		startingPoints int
		umaValues      []int // [1st, 2nd, 3rd, 4th]
		expectedFinal  []int
	}{
		{
			name:           "Standard 25000 start, 30000 return",
			finalScores:    []int{35000, 28000, 20000, 17000},
			startingPoints: 25000,
			umaValues:      []int{20000, 10000, -10000, -20000},
			expectedFinal:  []int{30000, 13000, -15000, -28000}, // After uma
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("NOT IMPLEMENTED - TDD Red Phase: Implement uma/oka calculation")
		})
	}
}
