package rules

import (
	"github.com/guny524/majhong_cli/pkg/mahjong/engine"
	"sort"
)

// YakuType represents a yaku pattern
type YakuType string

const (
	// 1-Han Yaku
	YakuRiichi       YakuType = "RIICHI"
	YakuMenzenTsumo  YakuType = "MENZEN_TSUMO"
	YakuTanyao       YakuType = "TANYAO"
	YakuPinfu        YakuType = "PINFU"
	YakuIipeikou     YakuType = "IIPEIKOU"
	YakuYakuhaiEast  YakuType = "YAKUHAI_WIND"
	YakuYakuhaiSouth YakuType = "YAKUHAI_WIND"
	YakuYakuhaiWest  YakuType = "YAKUHAI_WIND"
	YakuYakuhaiNorth YakuType = "YAKUHAI_WIND"
	YakuYakuhaiHaku  YakuType = "YAKUHAI_HAKU"
	YakuYakuhaiHatsu YakuType = "YAKUHAI_HATSU"
	YakuYakuhaiChun  YakuType = "YAKUHAI_CHUN"

	// 2-Han Yaku
	YakuChiitoitsu      YakuType = "CHIITOITSU"
	YakuToitoi          YakuType = "TOITOI"
	YakuSanankou        YakuType = "SANANKOU"
	YakuSanshokuDoujun  YakuType = "SANSHOKU_DOUJUN"
	YakuSanshokuDoukou  YakuType = "SANSHOKU_DOUKOU"
	YakuIttsu           YakuType = "ITTSU"
	YakuChanta          YakuType = "CHANTA"
	YakuHonroutou       YakuType = "HONROUTOU"
	YakuShousangen      YakuType = "SHOUSANGEN"
	YakuSankantsu       YakuType = "SANKANTSU"
	YakuDoubleRiichi    YakuType = "DOUBLE_RIICHI"

	// 3-Han Yaku
	YakuRyanpeikou YakuType = "RYANPEIKOU"
	YakuJunchan    YakuType = "JUNCHAN"
	YakuHonitsu    YakuType = "HONITSU"

	// 6-Han Yaku
	YakuChinitsu YakuType = "CHINITSU"

	// Yakuman (13 han)
	YakuKokushi    YakuType = "KOKUSHI"
	YakuSuuankou   YakuType = "SUUANKOU"
	YakuDaisangen  YakuType = "DAISANGEN"
	YakuShousuushii YakuType = "SHOUSUUSHII"
	YakuDaisuushii YakuType = "DAISUUSHII" // Double yakuman
	YakuTsuuiisou  YakuType = "TSUUIISOU"
	YakuRyuuiisou  YakuType = "RYUUIISOU"
	YakuChinroutou YakuType = "CHINROUTOU"
	YakuChuuren    YakuType = "CHUUREN"
	YakuSuukantsu  YakuType = "SUUKANTSU"
	YakuTenhou     YakuType = "TENHOU"
	YakuChiihou    YakuType = "CHIIHOU"
)

// Yaku represents a detected yaku with its han value
type Yaku struct {
	Type YakuType
	Han  int
}

// WinContext provides context for yaku detection
type WinContext struct {
	IsClosed  bool
	IsRiichi  bool
	IsTsumo   bool
	IsDealer  bool
	SeatWind  engine.Seat
	RoundWind engine.Wind
	Melds     []engine.Meld
	DoraCount int
	FirstTurn bool // For Tenhou/Chiihou detection
}

// HandAnalysis represents analyzed hand structure
type HandAnalysis struct {
	Tiles       []engine.Tile
	WinTile     engine.Tile
	Melds       []Meld
	Pair        Meld
	IsChiitoitsu bool
	IsKokushi   bool
}

// Meld represents a group of tiles (triplet or sequence)
type Meld struct {
	Type  MeldType
	Tiles []engine.Tile
}

// MeldType represents the type of meld
type MeldType int

const (
	MeldPair MeldType = iota
	MeldTriplet
	MeldSequence
	MeldQuad
)

// DetectYaku detects all yaku in a winning hand
func DetectYaku(hand []engine.Tile, winTile engine.Tile, context WinContext) []Yaku {
	yaku := []Yaku{}

	// Add melds from context to hand analysis
	allTiles := append([]engine.Tile{}, hand...)
	allTiles = append(allTiles, winTile)

	analysis := analyzeHand(allTiles, winTile, context.Melds)

	// Check for yakuman first (they override regular yaku)
	yakuman := detectYakuman(analysis, context)
	if len(yakuman) > 0 {
		return yakuman
	}

	// Detect regular yaku
	yaku = append(yaku, detectRegularYaku(analysis, context)...)

	return yaku
}

// detectYakuman checks for yakuman patterns
func detectYakuman(analysis HandAnalysis, context WinContext) []Yaku {
	yaku := []Yaku{}

	// Tenhou - Dealer wins on first draw
	if context.IsDealer && context.IsTsumo && context.FirstTurn && context.IsClosed {
		yaku = append(yaku, Yaku{Type: YakuTenhou, Han: 13})
		return yaku
	}

	// Chiihou - Non-dealer wins on first draw
	if !context.IsDealer && context.IsTsumo && context.FirstTurn && context.IsClosed {
		yaku = append(yaku, Yaku{Type: YakuChiihou, Han: 13})
		return yaku
	}

	// Kokushi Musou
	if isKokushi(analysis.Tiles) {
		yaku = append(yaku, Yaku{Type: YakuKokushi, Han: 13})
		return yaku
	}

	// Suuankou - Four concealed triplets
	if isSuuankou(analysis, context) {
		yaku = append(yaku, Yaku{Type: YakuSuuankou, Han: 13})
		return yaku
	}

	// Daisuushii - Four wind triplets (Double Yakuman)
	if isDaisuushii(analysis) {
		yaku = append(yaku, Yaku{Type: YakuDaisuushii, Han: 26})
		return yaku
	}

	// Shousuushii - Three wind triplets + wind pair
	if isShousuushii(analysis) {
		yaku = append(yaku, Yaku{Type: YakuShousuushii, Han: 13})
		return yaku
	}

	// Daisangen - Three dragon triplets
	if isDaisangen(analysis) {
		yaku = append(yaku, Yaku{Type: YakuDaisangen, Han: 13})
		return yaku
	}

	// Tsuuiisou - All honors
	if isTsuuiisou(analysis.Tiles) {
		yaku = append(yaku, Yaku{Type: YakuTsuuiisou, Han: 13})
		return yaku
	}

	// Ryuuiisou - All green (2,3,4,6,8 sou + green dragon)
	if isRyuuiisou(analysis.Tiles) {
		yaku = append(yaku, Yaku{Type: YakuRyuuiisou, Han: 13})
		return yaku
	}

	// Chinroutou - All terminals (no honors)
	if isChinroutou(analysis.Tiles) {
		yaku = append(yaku, Yaku{Type: YakuChinroutou, Han: 13})
		return yaku
	}

	// Chuuren Poutou - 1112345678999 + any tile (same suit, closed)
	if isChuuren(analysis.Tiles) && context.IsClosed {
		yaku = append(yaku, Yaku{Type: YakuChuuren, Han: 13})
		return yaku
	}

	// Suukantsu - Four quads
	if isSuukantsu(context.Melds) {
		yaku = append(yaku, Yaku{Type: YakuSuukantsu, Han: 13})
		return yaku
	}

	return yaku
}

// detectRegularYaku detects non-yakuman yaku
func detectRegularYaku(analysis HandAnalysis, context WinContext) []Yaku {
	yaku := []Yaku{}

	// Riichi
	if context.IsRiichi {
		// TODO: Check if double riichi (first turn riichi)
		if context.FirstTurn {
			yaku = append(yaku, Yaku{Type: YakuDoubleRiichi, Han: 2})
		} else {
			yaku = append(yaku, Yaku{Type: YakuRiichi, Han: 1})
		}
	}

	// Menzen Tsumo
	if context.IsClosed && context.IsTsumo {
		yaku = append(yaku, Yaku{Type: YakuMenzenTsumo, Han: 1})
	}

	// Chiitoitsu - Seven pairs (mutually exclusive with other patterns)
	if analysis.IsChiitoitsu {
		yaku = append(yaku, Yaku{Type: YakuChiitoitsu, Han: 2})
		return yaku // Stop here, cannot combine with other yaku
	}

	// Pinfu - All sequences, valueless pair, open wait
	if isPinfu(analysis, context) {
		yaku = append(yaku, Yaku{Type: YakuPinfu, Han: 1})
	}

	// Tanyao - All simples (2-8, no terminals or honors)
	if isTanyao(analysis.Tiles) {
		yaku = append(yaku, Yaku{Type: YakuTanyao, Han: 1})
	}

	// Yakuhai - Dragon or wind triplets
	yakuhaiYaku := detectYakuhai(analysis, context)
	yaku = append(yaku, yakuhaiYaku...)

	// Ryanpeikou - Two iipeikou (3 han, overrides iipeikou)
	if context.IsClosed && isRyanpeikou(analysis) {
		yaku = append(yaku, Yaku{Type: YakuRyanpeikou, Han: 3})
	} else if context.IsClosed && isIipeikou(analysis) {
		// Iipeikou - One pair of identical sequences
		yaku = append(yaku, Yaku{Type: YakuIipeikou, Han: 1})
	}

	// Toitoi - All triplets
	if isToitoi(analysis) {
		yaku = append(yaku, Yaku{Type: YakuToitoi, Han: 2})
	}

	// Sanankou - Three concealed triplets
	if isSanankou(analysis, context) {
		yaku = append(yaku, Yaku{Type: YakuSanankou, Han: 2})
	}

	// Sankantsu - Three quads
	if isSankantsu(context.Melds) {
		yaku = append(yaku, Yaku{Type: YakuSankantsu, Han: 2})
	}

	// Sanshoku Doujun - Same sequence in 3 suits (2 han closed, 1 han open)
	if isSanshokuDoujun(analysis) {
		if context.IsClosed {
			yaku = append(yaku, Yaku{Type: YakuSanshokuDoujun, Han: 2})
		} else {
			yaku = append(yaku, Yaku{Type: YakuSanshokuDoujun, Han: 1})
		}
	}

	// Sanshoku Doukou - Same triplet in 3 suits
	if isSanshokuDoukou(analysis) {
		yaku = append(yaku, Yaku{Type: YakuSanshokuDoukou, Han: 2})
	}

	// Ittsu - 123, 456, 789 in same suit (2 han closed, 1 han open)
	if isIttsu(analysis) {
		if context.IsClosed {
			yaku = append(yaku, Yaku{Type: YakuIttsu, Han: 2})
		} else {
			yaku = append(yaku, Yaku{Type: YakuIttsu, Han: 1})
		}
	}

	// Chinitsu/Honitsu - Flush (mutually exclusive)
	if isChinitsu(analysis.Tiles) {
		// Chinitsu - One suit only (6 han closed, 5 han open)
		if context.IsClosed {
			yaku = append(yaku, Yaku{Type: YakuChinitsu, Han: 6})
		} else {
			yaku = append(yaku, Yaku{Type: YakuChinitsu, Han: 5})
		}
	} else if isHonitsu(analysis.Tiles) {
		// Honitsu - One suit + honors (3 han closed, 2 han open)
		if context.IsClosed {
			yaku = append(yaku, Yaku{Type: YakuHonitsu, Han: 3})
		} else {
			yaku = append(yaku, Yaku{Type: YakuHonitsu, Han: 2})
		}
	}

	// Junchan/Chanta (mutually exclusive)
	if isJunchan(analysis) {
		// Junchan - All groups contain terminals (3 han closed, 2 han open)
		if context.IsClosed {
			yaku = append(yaku, Yaku{Type: YakuJunchan, Han: 3})
		} else {
			yaku = append(yaku, Yaku{Type: YakuJunchan, Han: 2})
		}
	} else if isChanta(analysis) {
		// Chanta - All groups contain terminals or honors (2 han closed, 1 han open)
		if context.IsClosed {
			yaku = append(yaku, Yaku{Type: YakuChanta, Han: 2})
		} else {
			yaku = append(yaku, Yaku{Type: YakuChanta, Han: 1})
		}
	}

	// Honroutou - All terminals and honors (no sequences)
	if isHonroutou(analysis.Tiles) {
		yaku = append(yaku, Yaku{Type: YakuHonroutou, Han: 2})
	}

	// Shousangen - Two dragon triplets + dragon pair
	if isShousangen(analysis) {
		yaku = append(yaku, Yaku{Type: YakuShousangen, Han: 2})
	}

	return yaku
}

// analyzeHand analyzes hand structure
func analyzeHand(tiles []engine.Tile, winTile engine.Tile, melds []engine.Meld) HandAnalysis {
	analysis := HandAnalysis{
		Tiles:   tiles,
		WinTile: winTile,
	}

	// Check for special patterns
	if isKokushi(tiles) {
		analysis.IsKokushi = true
		return analysis
	}

	// Try to find standard 4 groups + 1 pair pattern first
	sortedTiles := make([]engine.Tile, len(tiles))
	copy(sortedTiles, tiles)
	sortTiles(sortedTiles)

	// Find all possible meld decompositions
	analysis.Melds = findMelds(sortedTiles)

	// Only treat as chiitoitsu if standard decomposition failed
	if len(analysis.Melds) == 0 && isChiitoitsuPattern(tiles) {
		analysis.IsChiitoitsu = true
		return analysis
	}

	return analysis
}

// findMelds finds all melds in the hand using backtracking
func findMelds(tiles []engine.Tile) []Meld {
	if len(tiles) == 0 {
		return []Meld{}
	}

	// Try to find a valid 4 groups + 1 pair decomposition
	// This uses backtracking to try all possibilities
	var result []Meld
	findMeldsRecursive(tiles, []Meld{}, false, &result)
	return result
}

// findMeldsRecursive uses backtracking to find valid hand decomposition
func findMeldsRecursive(tiles []engine.Tile, current []Meld, foundPair bool, result *[]Meld) bool {
	if len(tiles) == 0 {
		// Successfully decomposed the hand
		if foundPair && len(current) == 5 {
			*result = make([]Meld, len(current))
			copy(*result, current)
			return true
		}
		return false
	}

	// Count tile occurrences
	counts := countTiles(tiles)

	// Try to form pair (if not found yet)
	if !foundPair && len(tiles) >= 2 {
		tile := tiles[0]
		if counts[tile.String()] >= 2 {
			// Try pair
			pair := Meld{
				Type:  MeldPair,
				Tiles: []engine.Tile{tile, tile},
			}
			newTiles := removeTiles(tiles, []engine.Tile{tile, tile})
			newCurrent := append(current, pair)
			if findMeldsRecursive(newTiles, newCurrent, true, result) {
				return true
			}
		}
	}

	// Try to form triplet
	if len(tiles) >= 3 {
		tile := tiles[0]
		if counts[tile.String()] >= 3 {
			triplet := Meld{
				Type:  MeldTriplet,
				Tiles: []engine.Tile{tile, tile, tile},
			}
			newTiles := removeTiles(tiles, []engine.Tile{tile, tile, tile})
			newCurrent := append(current, triplet)
			if findMeldsRecursive(newTiles, newCurrent, foundPair, result) {
				return true
			}
		}
	}

	// Try to form sequence
	if len(tiles) >= 3 && tiles[0].Suit != engine.Ji {
		tile := tiles[0]
		next1 := engine.NewTile(tile.Suit, tile.Rank+1)
		next2 := engine.NewTile(tile.Suit, tile.Rank+2)

		if tile.Rank <= 7 && hasTile(tiles, next1) && hasTile(tiles, next2) {
			sequence := Meld{
				Type:  MeldSequence,
				Tiles: []engine.Tile{tile, next1, next2},
			}
			newTiles := removeTiles(tiles, []engine.Tile{tile, next1, next2})
			newCurrent := append(current, sequence)
			if findMeldsRecursive(newTiles, newCurrent, foundPair, result) {
				return true
			}
		}
	}

	// Skip this tile and try next (dead-end in this branch)
	if len(tiles) > 1 {
		return findMeldsRecursive(tiles[1:], current, foundPair, result)
	}

	return false
}

// countTiles counts occurrences of each tile type
func countTiles(tiles []engine.Tile) map[string]int {
	counts := make(map[string]int)
	for _, tile := range tiles {
		counts[tile.String()]++
	}
	return counts
}

// hasTile checks if a tile exists in the slice
func hasTile(tiles []engine.Tile, target engine.Tile) bool {
	for _, tile := range tiles {
		if tile.SameType(target) {
			return true
		}
	}
	return false
}

// removeTiles removes specified tiles from the slice (one occurrence each)
func removeTiles(tiles []engine.Tile, toRemove []engine.Tile) []engine.Tile {
	result := make([]engine.Tile, len(tiles))
	copy(result, tiles)

	for _, remove := range toRemove {
		for i, tile := range result {
			if tile.SameType(remove) {
				result = append(result[:i], result[i+1:]...)
				break
			}
		}
	}

	return result
}

// isSequence checks if three tiles form a sequence
func isSequence(t1, t2, t3 engine.Tile) bool {
	if t1.Suit != t2.Suit || t2.Suit != t3.Suit {
		return false
	}
	if t1.Suit == engine.Ji {
		return false // No sequences in honors
	}
	ranks := []int{t1.Rank, t2.Rank, t3.Rank}
	sort.Ints(ranks)
	return ranks[1] == ranks[0]+1 && ranks[2] == ranks[1]+1
}

// sortTiles sorts tiles by suit then rank
func sortTiles(tiles []engine.Tile) {
	sort.Slice(tiles, func(i, j int) bool {
		if tiles[i].Suit != tiles[j].Suit {
			return tiles[i].Suit < tiles[j].Suit
		}
		return tiles[i].Rank < tiles[j].Rank
	})
}

// Yaku detection helper functions

func isKokushi(tiles []engine.Tile) bool {
	if len(tiles) != 14 {
		return false
	}

	// Kokushi requires: 1m, 9m, 1p, 9p, 1s, 9s, 1z-7z (13 types)
	required := map[string]int{
		"1m": 0, "9m": 0,
		"1p": 0, "9p": 0,
		"1s": 0, "9s": 0,
		"1z": 0, "2z": 0, "3z": 0, "4z": 0,
		"5z": 0, "6z": 0, "7z": 0,
	}

	for _, tile := range tiles {
		key := tile.String()
		if _, exists := required[key]; exists {
			required[key]++
		} else {
			return false
		}
	}

	// All 13 types must appear, one appears twice
	pairCount := 0
	for _, count := range required {
		if count == 0 {
			return false
		}
		if count == 2 {
			pairCount++
		} else if count != 1 {
			return false
		}
	}

	return pairCount == 1
}

func isChiitoitsuPattern(tiles []engine.Tile) bool {
	if len(tiles) != 14 {
		return false
	}

	sorted := make([]engine.Tile, len(tiles))
	copy(sorted, tiles)
	sortTiles(sorted)

	// Check for 7 pairs
	for i := 0; i < 14; i += 2 {
		if !sorted[i].SameType(sorted[i+1]) {
			return false
		}
	}

	return true
}

func isSuuankou(analysis HandAnalysis, context WinContext) bool {
	if !context.IsClosed || !context.IsTsumo {
		return false
	}

	tripletCount := 0
	for _, meld := range analysis.Melds {
		if meld.Type == MeldTriplet || meld.Type == MeldQuad {
			tripletCount++
		}
	}

	return tripletCount == 4
}

func isDaisuushii(analysis HandAnalysis) bool {
	windTriplets := 0
	for _, meld := range analysis.Melds {
		if meld.Type == MeldTriplet || meld.Type == MeldQuad {
			if len(meld.Tiles) > 0 && meld.Tiles[0].Suit == engine.Ji && meld.Tiles[0].Rank <= 4 {
				windTriplets++
			}
		}
	}
	return windTriplets == 4
}

func isShousuushii(analysis HandAnalysis) bool {
	windTriplets := 0
	windPair := false

	for _, meld := range analysis.Melds {
		if len(meld.Tiles) > 0 && meld.Tiles[0].Suit == engine.Ji && meld.Tiles[0].Rank <= 4 {
			if meld.Type == MeldTriplet || meld.Type == MeldQuad {
				windTriplets++
			} else if meld.Type == MeldPair {
				windPair = true
			}
		}
	}

	return windTriplets == 3 && windPair
}

func isDaisangen(analysis HandAnalysis) bool {
	dragonTriplets := 0
	for _, meld := range analysis.Melds {
		if meld.Type == MeldTriplet || meld.Type == MeldQuad {
			if len(meld.Tiles) > 0 && meld.Tiles[0].Suit == engine.Ji && meld.Tiles[0].Rank >= 5 && meld.Tiles[0].Rank <= 7 {
				dragonTriplets++
			}
		}
	}
	return dragonTriplets == 3
}

func isTsuuiisou(tiles []engine.Tile) bool {
	for _, tile := range tiles {
		if tile.Suit != engine.Ji {
			return false
		}
	}
	return true
}

func isRyuuiisou(tiles []engine.Tile) bool {
	// Green tiles: 2s, 3s, 4s, 6s, 8s, 6z (green dragon)
	for _, tile := range tiles {
		if tile.Suit == engine.Sou {
			if tile.Rank != 2 && tile.Rank != 3 && tile.Rank != 4 && tile.Rank != 6 && tile.Rank != 8 {
				return false
			}
		} else if tile.Suit == engine.Ji {
			if tile.Rank != 6 {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func isChinroutou(tiles []engine.Tile) bool {
	for _, tile := range tiles {
		if !tile.IsTerminal() {
			return false
		}
	}
	return true
}

func isChuuren(tiles []engine.Tile) bool {
	if len(tiles) != 14 {
		return false
	}

	// Check all tiles are same suit
	suit := tiles[0].Suit
	if suit == engine.Ji {
		return false
	}

	for _, tile := range tiles {
		if tile.Suit != suit {
			return false
		}
	}

	// Count each rank
	counts := make([]int, 10) // 1-9
	for _, tile := range tiles {
		counts[tile.Rank]++
	}

	// Should be 1112345678999 pattern (3x1, 1x2-8, 3x9, +1 extra)
	if counts[1] < 3 || counts[9] < 3 {
		return false
	}
	for i := 2; i <= 8; i++ {
		if counts[i] < 1 {
			return false
		}
	}

	return true
}

func isSuukantsu(melds []engine.Meld) bool {
	kanCount := 0
	for _, meld := range melds {
		if meld.Type == engine.MeldAnkan || meld.Type == engine.MeldMinkan || meld.Type == engine.MeldKakan {
			kanCount++
		}
	}
	return kanCount == 4
}

func isSankantsu(melds []engine.Meld) bool {
	kanCount := 0
	for _, meld := range melds {
		if meld.Type == engine.MeldAnkan || meld.Type == engine.MeldMinkan || meld.Type == engine.MeldKakan {
			kanCount++
		}
	}
	return kanCount == 3
}

func isPinfu(analysis HandAnalysis, context WinContext) bool {
	if !context.IsClosed {
		return false
	}

	// All melds must be sequences
	for _, meld := range analysis.Melds {
		if meld.Type == MeldTriplet || meld.Type == MeldQuad {
			return false
		}
	}

	// Pair must be valueless (not dragons, not seat/round wind)
	// TODO: Implement wind check

	// Wait must be ryanmen (open wait)
	// TODO: Implement wait type detection

	return true
}

func isTanyao(tiles []engine.Tile) bool {
	for _, tile := range tiles {
		if tile.IsTerminal() || tile.IsHonor() {
			return false
		}
	}
	return true
}

func isToitoi(analysis HandAnalysis) bool {
	for _, meld := range analysis.Melds {
		if meld.Type == MeldSequence {
			return false
		}
	}
	return true
}

func isSanankou(analysis HandAnalysis, context WinContext) bool {
	if context.IsTsumo {
		return false // Tsumo win = 4 concealed = suuankou
	}

	concealedCount := 0
	for _, meld := range analysis.Melds {
		if meld.Type == MeldTriplet || meld.Type == MeldQuad {
			// Check if meld contains win tile
			containsWin := false
			for _, tile := range meld.Tiles {
				if tile.Equals(analysis.WinTile) {
					containsWin = true
					break
				}
			}
			if !containsWin {
				concealedCount++
			}
		}
	}

	return concealedCount == 3
}

func isIipeikou(analysis HandAnalysis) bool {
	sequences := []Meld{}
	for _, meld := range analysis.Melds {
		if meld.Type == MeldSequence {
			sequences = append(sequences, meld)
		}
	}

	// Find matching sequences
	for i := 0; i < len(sequences); i++ {
		for j := i + 1; j < len(sequences); j++ {
			if sameMeld(sequences[i], sequences[j]) {
				return true
			}
		}
	}

	return false
}

func isRyanpeikou(analysis HandAnalysis) bool {
	sequences := []Meld{}
	for _, meld := range analysis.Melds {
		if meld.Type == MeldSequence {
			sequences = append(sequences, meld)
		}
	}

	if len(sequences) != 4 {
		return false
	}

	// Find two pairs of matching sequences
	matched := make([]bool, 4)
	pairCount := 0

	for i := 0; i < 4; i++ {
		if matched[i] {
			continue
		}
		for j := i + 1; j < 4; j++ {
			if matched[j] {
				continue
			}
			if sameMeld(sequences[i], sequences[j]) {
				matched[i] = true
				matched[j] = true
				pairCount++
				break
			}
		}
	}

	return pairCount == 2
}

func sameMeld(m1, m2 Meld) bool {
	if m1.Type != m2.Type || len(m1.Tiles) != len(m2.Tiles) {
		return false
	}
	for i := range m1.Tiles {
		if !m1.Tiles[i].SameType(m2.Tiles[i]) {
			return false
		}
	}
	return true
}

func isSanshokuDoujun(analysis HandAnalysis) bool {
	sequences := []Meld{}
	for _, meld := range analysis.Melds {
		if meld.Type == MeldSequence {
			sequences = append(sequences, meld)
		}
	}

	// Check for same sequence pattern in different suits
	for i := 0; i < len(sequences); i++ {
		for j := i + 1; j < len(sequences); j++ {
			for k := j + 1; k < len(sequences); k++ {
				if isSameSequencePattern(sequences[i], sequences[j], sequences[k]) {
					return true
				}
			}
		}
	}

	return false
}

func isSameSequencePattern(m1, m2, m3 Meld) bool {
	if m1.Tiles[0].Suit == m2.Tiles[0].Suit || m2.Tiles[0].Suit == m3.Tiles[0].Suit || m1.Tiles[0].Suit == m3.Tiles[0].Suit {
		return false // Must be different suits
	}

	return m1.Tiles[0].Rank == m2.Tiles[0].Rank && m2.Tiles[0].Rank == m3.Tiles[0].Rank
}

func isSanshokuDoukou(analysis HandAnalysis) bool {
	triplets := []Meld{}
	for _, meld := range analysis.Melds {
		if meld.Type == MeldTriplet {
			triplets = append(triplets, meld)
		}
	}

	// Check for same triplet rank in different suits
	for i := 0; i < len(triplets); i++ {
		for j := i + 1; j < len(triplets); j++ {
			for k := j + 1; k < len(triplets); k++ {
				if isSameTripletPattern(triplets[i], triplets[j], triplets[k]) {
					return true
				}
			}
		}
	}

	return false
}

func isSameTripletPattern(m1, m2, m3 Meld) bool {
	s1, s2, s3 := m1.Tiles[0].Suit, m2.Tiles[0].Suit, m3.Tiles[0].Suit

	if s1 == engine.Ji || s2 == engine.Ji || s3 == engine.Ji {
		return false // No honors
	}

	if s1 == s2 || s2 == s3 || s1 == s3 {
		return false // Must be different suits
	}

	return m1.Tiles[0].Rank == m2.Tiles[0].Rank && m2.Tiles[0].Rank == m3.Tiles[0].Rank
}

func isIttsu(analysis HandAnalysis) bool {
	sequences := []Meld{}
	for _, meld := range analysis.Melds {
		if meld.Type == MeldSequence {
			sequences = append(sequences, meld)
		}
	}

	// Group by suit
	bySuit := make(map[engine.Suit][]Meld)
	for _, seq := range sequences {
		suit := seq.Tiles[0].Suit
		bySuit[suit] = append(bySuit[suit], seq)
	}

	// Check each suit for 123, 456, 789
	for _, seqs := range bySuit {
		has123, has456, has789 := false, false, false
		for _, seq := range seqs {
			rank := seq.Tiles[0].Rank
			if rank == 1 {
				has123 = true
			} else if rank == 4 {
				has456 = true
			} else if rank == 7 {
				has789 = true
			}
		}
		if has123 && has456 && has789 {
			return true
		}
	}

	return false
}

func isChinitsu(tiles []engine.Tile) bool {
	if len(tiles) == 0 {
		return false
	}

	suit := tiles[0].Suit
	if suit == engine.Ji {
		return false
	}

	for _, tile := range tiles {
		if tile.Suit != suit {
			return false
		}
	}

	return true
}

func isHonitsu(tiles []engine.Tile) bool {
	suitCount := make(map[engine.Suit]int)
	for _, tile := range tiles {
		suitCount[tile.Suit]++
	}

	// Must have honors + exactly one number suit
	hasHonors := suitCount[engine.Ji] > 0
	numSuits := 0
	for suit := range suitCount {
		if suit != engine.Ji {
			numSuits++
		}
	}

	return hasHonors && numSuits == 1
}

func isJunchan(analysis HandAnalysis) bool {
	// All groups must contain terminals (1 or 9), no honors
	for _, meld := range analysis.Melds {
		if len(meld.Tiles) == 0 {
			return false
		}

		hasTerminal := false
		for _, tile := range meld.Tiles {
			if tile.IsHonor() {
				return false // No honors allowed
			}
			if tile.IsTerminal() {
				hasTerminal = true
			}
		}
		if !hasTerminal {
			return false
		}
	}

	return true
}

func isChanta(analysis HandAnalysis) bool {
	// All groups must contain terminals or honors
	for _, meld := range analysis.Melds {
		if len(meld.Tiles) == 0 {
			return false
		}

		hasTerminalOrHonor := false
		for _, tile := range meld.Tiles {
			if tile.IsTerminal() || tile.IsHonor() {
				hasTerminalOrHonor = true
				break
			}
		}
		if !hasTerminalOrHonor {
			return false
		}
	}

	return true
}

func isHonroutou(tiles []engine.Tile) bool {
	for _, tile := range tiles {
		if !tile.IsTerminal() && !tile.IsHonor() {
			return false
		}
	}
	return true
}

func isShousangen(analysis HandAnalysis) bool {
	dragonTriplets := 0
	dragonPair := false

	for _, meld := range analysis.Melds {
		if len(meld.Tiles) > 0 && meld.Tiles[0].Suit == engine.Ji && meld.Tiles[0].Rank >= 5 && meld.Tiles[0].Rank <= 7 {
			if meld.Type == MeldTriplet || meld.Type == MeldQuad {
				dragonTriplets++
			} else if meld.Type == MeldPair {
				dragonPair = true
			}
		}
	}

	return dragonTriplets == 2 && dragonPair
}

func detectYakuhai(analysis HandAnalysis, context WinContext) []Yaku {
	yaku := []Yaku{}

	for _, meld := range analysis.Melds {
		if meld.Type != MeldTriplet && meld.Type != MeldQuad {
			continue
		}

		if len(meld.Tiles) == 0 {
			continue
		}

		tile := meld.Tiles[0]
		if tile.Suit != engine.Ji {
			continue
		}

		// Dragons (5z-7z)
		if tile.Rank == 5 {
			yaku = append(yaku, Yaku{Type: YakuYakuhaiHaku, Han: 1})
		} else if tile.Rank == 6 {
			yaku = append(yaku, Yaku{Type: YakuYakuhaiHatsu, Han: 1})
		} else if tile.Rank == 7 {
			yaku = append(yaku, Yaku{Type: YakuYakuhaiChun, Han: 1})
		} else if tile.Rank >= 1 && tile.Rank <= 4 {
			// Winds (1z-4z)
			han := 0

			// Check if seat wind matches
			if int(context.SeatWind) == tile.Rank-1 {
				han++
			}

			// Check if round wind matches
			if int(context.RoundWind) == tile.Rank-1 {
				han++
			}

			if han > 0 {
				yaku = append(yaku, Yaku{Type: YakuYakuhaiEast, Han: han})
			}
		}
	}

	return yaku
}
