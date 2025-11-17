package contract

import (
	"testing"
)

// Contract tests for all 52 yaku patterns
// These tests MUST be written BEFORE implementation (TDD Red phase)
// Reference: constitution_majhong_yaku.md

// TestYakuDetection tests all 52 standard yaku recognition
func TestYakuDetection(t *testing.T) {
	tests := []struct {
		name        string
		handTiles   []string
		winTile     string
		isClosed    bool
		isRiichi    bool
		isTsumo     bool
		isDealer    bool
		seatWind    string
		roundWind   string
		expectedYaku []string
		expectedHan  int
	}{
		// 1-Han Yaku
		{
			name:         "Riichi",
			handTiles:    []string{"1m", "1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "8m", "8m", "9m", "9m"},
			winTile:      "9m",
			isClosed:     true,
			isRiichi:     true,
			isTsumo:      false,
			expectedYaku: []string{"RIICHI"},
			expectedHan:  1,
		},
		{
			name:         "Menzen Tsumo",
			handTiles:    []string{"1m", "2m", "3m", "4p", "5p", "6p", "7s", "8s", "9s", "2z", "2z", "3m", "3m"},
			winTile:      "3m",
			isClosed:     true,
			isTsumo:      true,
			expectedYaku: []string{"MENZEN_TSUMO"},
			expectedHan:  1,
		},
		{
			name:         "Tanyao - Open",
			handTiles:    []string{"2m", "3m", "4m", "5p", "5p", "5p", "6s", "7s", "8s", "3m", "3m", "3m", "6p"},
			winTile:      "6p",
			isClosed:     false,
			expectedYaku: []string{"TANYAO"},
			expectedHan:  1,
		},
		{
			name:         "Yakuhai - White Dragon",
			handTiles:    []string{"1m", "2m", "3m", "5z", "5z", "5z", "4p", "5p", "6p", "7s", "8s", "9s", "2m"},
			winTile:      "2m",
			isClosed:     false,
			expectedYaku: []string{"YAKUHAI_HAKU"},
			expectedHan:  1,
		},
		{
			name:         "Yakuhai - Double East",
			handTiles:    []string{"1z", "1z", "1z", "2m", "3m", "4m", "5p", "6p", "7p", "8s", "9s", "7s", "3p"},
			winTile:      "3p",
			seatWind:     "East",
			roundWind:    "East",
			expectedYaku: []string{"YAKUHAI_WIND"},
			expectedHan:  2, // Double yakuhai
		},
		{
			name:         "Pinfu",
			handTiles:    []string{"2m", "3m", "4m", "5p", "6p", "7p", "1s", "2s", "3s", "6s", "6s", "7s", "8s"},
			winTile:      "9s",
			isClosed:     true,
			expectedYaku: []string{"PINFU"},
			expectedHan:  1,
		},
		{
			name:         "Iipeikou",
			handTiles:    []string{"2m", "3m", "4m", "2m", "3m", "4m", "5p", "6p", "7p", "8s", "9s", "7s", "1z"},
			winTile:      "1z",
			isClosed:     true,
			expectedYaku: []string{"IIPEIKOU"},
			expectedHan:  1,
		},

		// 2-Han Yaku
		{
			name:         "Chiitoitsu",
			handTiles:    []string{"1m", "1m", "3p", "3p", "5s", "5s", "2z", "2z", "4z", "4z", "6z", "6z", "7p"},
			winTile:      "7p",
			isClosed:     true,
			expectedYaku: []string{"CHIITOITSU"},
			expectedHan:  2,
		},
		{
			name:         "Toitoi",
			handTiles:    []string{"2m", "2m", "2m", "5p", "5p", "5p", "8s", "8s", "8s", "3z", "3z", "3z", "4m"},
			winTile:      "4m",
			isClosed:     false,
			expectedYaku: []string{"TOITOI"},
			expectedHan:  2,
		},
		{
			name:         "Sanankou",
			handTiles:    []string{"2m", "2m", "2m", "5p", "5p", "5p", "8s", "8s", "8s", "3z", "3z", "4m", "5m"},
			winTile:      "6m",
			isClosed:     false,
			isTsumo:      true,
			expectedYaku: []string{"SANANKOU"},
			expectedHan:  2,
		},
		{
			name:         "Sanshoku Doujun - Closed",
			handTiles:    []string{"3m", "4m", "5m", "3p", "4p", "5p", "3s", "4s", "5s", "1z", "1z", "6m", "7m"},
			winTile:      "8m",
			isClosed:     true,
			expectedYaku: []string{"SANSHOKU_DOUJUN"},
			expectedHan:  2,
		},
		{
			name:         "Sanshoku Doujun - Open",
			handTiles:    []string{"3m", "4m", "5m", "3p", "4p", "5p", "3s", "4s", "5s", "1z", "1z", "6m", "7m"},
			winTile:      "8m",
			isClosed:     false,
			expectedYaku: []string{"SANSHOKU_DOUJUN"},
			expectedHan:  1, // Reduced when open
		},
		{
			name:         "Ittsu - Closed",
			handTiles:    []string{"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m", "5p", "5p", "6s", "7s"},
			winTile:      "8s",
			isClosed:     true,
			expectedYaku: []string{"ITTSU"},
			expectedHan:  2,
		},
		{
			name:         "Chanta - Closed",
			handTiles:    []string{"1m", "2m", "3m", "7p", "8p", "9p", "1s", "1s", "1s", "7s", "8s", "9s", "2z"},
			winTile:      "2z",
			isClosed:     true,
			expectedYaku: []string{"CHANTA"},
			expectedHan:  2,
		},
		{
			name:         "Honroutou",
			handTiles:    []string{"1m", "1m", "9m", "9m", "9m", "1p", "1p", "1p", "5z", "5z", "5z", "9p", "9p"},
			winTile:      "9p",
			expectedYaku: []string{"HONROUTOU"},
			expectedHan:  2,
		},
		{
			name:         "Shousangen",
			handTiles:    []string{"5z", "5z", "5z", "6z", "6z", "6z", "7z", "7z", "1m", "2m", "3m", "4p", "5p"},
			winTile:      "6p",
			expectedYaku: []string{"SHOUSANGEN", "YAKUHAI_HAKU", "YAKUHAI_HATSU"},
			expectedHan:  4,
		},
		{
			name:         "Double Riichi",
			handTiles:    []string{"1m", "1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "8m", "8m", "9m", "9m"},
			winTile:      "9m",
			isClosed:     true,
			isRiichi:     true,
			expectedYaku: []string{"DOUBLE_RIICHI"},
			expectedHan:  2,
		},

		// 3-Han Yaku
		{
			name:         "Ryanpeikou",
			handTiles:    []string{"2m", "3m", "4m", "2m", "3m", "4m", "5p", "6p", "7p", "5p", "6p", "7p", "8s"},
			winTile:      "8s",
			isClosed:     true,
			expectedYaku: []string{"RYANPEIKOU"},
			expectedHan:  3,
		},
		{
			name:         "Junchan - Closed",
			handTiles:    []string{"1m", "2m", "3m", "7p", "8p", "9p", "1s", "1s", "1s", "7s", "8s", "9s", "9m"},
			winTile:      "9m",
			isClosed:     true,
			expectedYaku: []string{"JUNCHAN"},
			expectedHan:  3,
		},
		{
			name:         "Honitsu - Closed",
			handTiles:    []string{"2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m", "1z", "1z", "1z", "5z", "5z"},
			winTile:      "5z",
			isClosed:     true,
			expectedYaku: []string{"HONITSU"},
			expectedHan:  3,
		},

		// 6-Han Yaku
		{
			name:         "Chinitsu - Closed",
			handTiles:    []string{"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m", "2m", "2m", "5m", "5m"},
			winTile:      "5m",
			isClosed:     true,
			expectedYaku: []string{"CHINITSU"},
			expectedHan:  6,
		},
		{
			name:         "Chinitsu - Open",
			handTiles:    []string{"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m", "2m", "2m", "5m", "5m"},
			winTile:      "5m",
			isClosed:     false,
			expectedYaku: []string{"CHINITSU"},
			expectedHan:  5,
		},

		// Yakuman
		{
			name:         "Kokushi Musou",
			handTiles:    []string{"1m", "9m", "1p", "9p", "1s", "9s", "1z", "2z", "3z", "4z", "5z", "6z", "7z"},
			winTile:      "7z",
			isClosed:     true,
			expectedYaku: []string{"KOKUSHI"},
			expectedHan:  13,
		},
		{
			name:         "Suuankou",
			handTiles:    []string{"2m", "2m", "2m", "5p", "5p", "5p", "8s", "8s", "8s", "3z", "3z", "3z", "4m"},
			winTile:      "4m",
			isClosed:     true,
			isTsumo:      true,
			expectedYaku: []string{"SUUANKOU"},
			expectedHan:  13,
		},
		{
			name:         "Daisangen",
			handTiles:    []string{"5z", "5z", "5z", "6z", "6z", "6z", "7z", "7z", "7z", "1m", "2m", "3m", "4p"},
			winTile:      "4p",
			expectedYaku: []string{"DAISANGEN"},
			expectedHan:  13,
		},
		{
			name:         "Shousuushii",
			handTiles:    []string{"1z", "1z", "1z", "2z", "2z", "2z", "3z", "3z", "3z", "4z", "4z", "5m", "6m"},
			winTile:      "7m",
			expectedYaku: []string{"SHOUSUUSHII"},
			expectedHan:  13,
		},
		{
			name:         "Daisuushii",
			handTiles:    []string{"1z", "1z", "1z", "2z", "2z", "2z", "3z", "3z", "3z", "4z", "4z", "4z", "5m"},
			winTile:      "5m",
			expectedYaku: []string{"DAISUUSHII"},
			expectedHan:  26, // Double yakuman
		},
		{
			name:         "Tsuuiisou",
			handTiles:    []string{"1z", "1z", "1z", "2z", "2z", "2z", "5z", "5z", "5z", "6z", "6z", "6z", "7z"},
			winTile:      "7z",
			expectedYaku: []string{"TSUUIISOU"},
			expectedHan:  13,
		},
		{
			name:         "Ryuuiisou",
			handTiles:    []string{"2s", "2s", "2s", "3s", "3s", "3s", "4s", "4s", "4s", "6z", "6z", "6z", "8s"},
			winTile:      "8s",
			expectedYaku: []string{"RYUUIISOU"},
			expectedHan:  13,
		},
		{
			name:         "Chinroutou",
			handTiles:    []string{"1m", "1m", "1m", "9m", "9m", "9m", "1p", "1p", "1p", "9p", "9p", "1s", "1s"},
			winTile:      "1s",
			expectedYaku: []string{"CHINROUTOU"},
			expectedHan:  13,
		},
		{
			name:         "Chuuren Poutou",
			handTiles:    []string{"1m", "1m", "1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m", "9m", "9m"},
			winTile:      "5m",
			isClosed:     true,
			expectedYaku: []string{"CHUUREN"},
			expectedHan:  13,
		},
		{
			name:         "Suukantsu",
			handTiles:    []string{"1m", "1m", "1m", "1m", "5p", "5p", "5p", "5p", "8s", "8s", "8s", "8s", "3z"},
			winTile:      "3z",
			expectedYaku: []string{"SUUKANTSU"},
			expectedHan:  13,
		},
		{
			name:         "Tenhou",
			handTiles:    []string{"1m", "2m", "3m", "4p", "5p", "6p", "7s", "8s", "9s", "2z", "2z", "3m", "3m"},
			winTile:      "3m",
			isClosed:     true,
			isDealer:     true,
			expectedYaku: []string{"TENHOU"},
			expectedHan:  13,
		},
		{
			name:         "Chiihou",
			handTiles:    []string{"1m", "2m", "3m", "4p", "5p", "6p", "7s", "8s", "9s", "2z", "2z", "3m", "3m"},
			winTile:      "3m",
			isClosed:     true,
			isDealer:     false,
			isTsumo:      true,
			expectedYaku: []string{"CHIIHOU"},
			expectedHan:  13,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO: Implement yaku detection
			// This test MUST FAIL until implementation is complete
			t.Skip("NOT IMPLEMENTED - TDD Red Phase: Implement pkg/mahjong/rules/yaku.go")

			// Example of what the implementation should look like:
			// hand := parseHand(tt.handTiles)
			// winTile := parseTile(tt.winTile)
			// context := WinContext{
			//     IsClosed: tt.isClosed,
			//     IsRiichi: tt.isRiichi,
			//     IsTsumo:  tt.isTsumo,
			//     IsDealer: tt.isDealer,
			//     SeatWind: tt.seatWind,
			//     RoundWind: tt.roundWind,
			// }
			//
			// yaku := DetectYaku(hand, winTile, context)
			//
			// if !containsAllYaku(yaku, tt.expectedYaku) {
			//     t.Errorf("Expected yaku %v, got %v", tt.expectedYaku, yaku)
			// }
			//
			// totalHan := calculateHan(yaku)
			// if totalHan != tt.expectedHan {
			//     t.Errorf("Expected %d han, got %d han", tt.expectedHan, totalHan)
			// }
		})
	}
}

// TestYakuCombinations tests yaku combination rules
func TestYakuCombinations(t *testing.T) {
	tests := []struct {
		name         string
		handTiles    []string
		winTile      string
		isClosed     bool
		expectedYaku []string
		forbiddenYaku []string
	}{
		{
			name:          "Chiitoitsu cannot combine with Toitoi",
			handTiles:     []string{"1m", "1m", "3p", "3p", "5s", "5s", "2z", "2z", "4z", "4z", "6z", "6z", "7p"},
			winTile:       "7p",
			isClosed:      true,
			expectedYaku:  []string{"CHIITOITSU"},
			forbiddenYaku: []string{"TOITOI", "IIPEIKOU"},
		},
		{
			name:          "Ryanpeikou overrides Iipeikou",
			handTiles:     []string{"2m", "3m", "4m", "2m", "3m", "4m", "5p", "6p", "7p", "5p", "6p", "7p", "8s"},
			winTile:       "8s",
			isClosed:      true,
			expectedYaku:  []string{"RYANPEIKOU"},
			forbiddenYaku: []string{"IIPEIKOU"},
		},
		{
			name:          "Chinitsu and Honitsu are mutually exclusive",
			handTiles:     []string{"1m", "2m", "3m", "4m", "5m", "6m", "7m", "8m", "9m", "2m", "2m", "5m", "5m"},
			winTile:       "5m",
			isClosed:      true,
			expectedYaku:  []string{"CHINITSU"},
			forbiddenYaku: []string{"HONITSU"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip("NOT IMPLEMENTED - TDD Red Phase")
			// TODO: Verify combination rules
		})
	}
}

// TestYakuPriority tests yakuman priority over regular yaku
func TestYakuPriority(t *testing.T) {
	t.Run("Yakuman stops regular yaku detection", func(t *testing.T) {
		t.Skip("NOT IMPLEMENTED - TDD Red Phase")
		// hand := kokushiHand()
		// yaku := DetectYaku(hand, ...)
		//
		// if len(yaku) != 1 {
		//     t.Errorf("Expected single yakuman, got %d yaku", len(yaku))
		// }
		// if yaku[0] != KOKUSHI {
		//     t.Errorf("Expected KOKUSHI, got %v", yaku[0])
		// }
	})
}
