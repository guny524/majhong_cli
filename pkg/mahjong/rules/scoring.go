package rules

import (
	"github.com/guny524/majhong_cli/pkg/mahjong/engine"
)

// Score represents the calculated score for a winning hand
type Score struct {
	Han          int
	Fu           int
	BasePoints   int
	RonPayment   int
	TsumoPayment []int // [from non-dealer, from dealer] or [from each] for dealer
	TotalPoints  int
	LimitHand    string // "Mangan", "Haneman", "Baiman", "Sanbaiman", "Yakuman", "Double Yakuman"
}

// CalculateScore calculates the final score based on han and fu
func CalculateScore(han int, fu int, isDealer bool, isTsumo bool) Score {
	score := Score{
		Han: han,
		Fu:  fu,
	}

	// Determine if this is a limit hand
	if han >= 13 {
		// Yakuman or higher
		if han >= 26 {
			score.LimitHand = "Double Yakuman"
			score.BasePoints = 16000
		} else if han >= 13 {
			score.LimitHand = "Yakuman"
			score.BasePoints = 8000
		}
	} else if han >= 11 {
		score.LimitHand = "Sanbaiman"
		score.BasePoints = 6000
	} else if han >= 8 {
		score.LimitHand = "Baiman"
		score.BasePoints = 4000
	} else if han >= 6 {
		score.LimitHand = "Haneman"
		score.BasePoints = 3000
	} else if han >= 5 {
		score.LimitHand = "Mangan"
		score.BasePoints = 2000
	} else if han == 4 && fu >= 40 {
		score.LimitHand = "Mangan"
		score.BasePoints = 2000
	} else if han == 3 && fu >= 70 {
		score.LimitHand = "Mangan"
		score.BasePoints = 2000
	} else {
		// Calculate base points from han and fu
		score.BasePoints = calculateBasePoints(han, fu)

		// Check if it reaches mangan
		if score.BasePoints >= 2000 {
			score.LimitHand = "Mangan"
			score.BasePoints = 2000
		}
	}

	// Calculate final payments
	if isDealer {
		if isTsumo {
			// Dealer tsumo: each player pays 2x base points
			payment := roundUp(score.BasePoints * 2)
			score.TsumoPayment = []int{payment}
			score.TotalPoints = payment * 3
		} else {
			// Dealer ron: loser pays 6x base points
			score.RonPayment = roundUp(score.BasePoints * 6)
			score.TotalPoints = score.RonPayment
		}
	} else {
		if isTsumo {
			// Non-dealer tsumo: non-dealers pay base points, dealer pays 2x
			nonDealerPayment := roundUp(score.BasePoints)
			dealerPayment := roundUp(score.BasePoints * 2)
			score.TsumoPayment = []int{nonDealerPayment, dealerPayment}
			score.TotalPoints = nonDealerPayment*2 + dealerPayment
		} else {
			// Non-dealer ron: loser pays 4x base points
			score.RonPayment = roundUp(score.BasePoints * 4)
			score.TotalPoints = score.RonPayment
		}
	}

	return score
}

// calculateBasePoints calculates base points from han and fu
func calculateBasePoints(han int, fu int) int {
	// Base formula: fu * 2^(han+2)
	base := fu * (1 << uint(han+2))

	// Cap at 2000 (mangan)
	if base > 2000 {
		return 2000
	}

	return base
}

// roundUp rounds up to the nearest 100
func roundUp(points int) int {
	return ((points + 99) / 100) * 100
}

// FuContext provides context for fu calculation
type FuContext struct {
	Hand         []engine.Tile
	WinTile      engine.Tile
	Melds        []Meld
	Pair         Meld
	WaitType     WaitType
	IsTsumo      bool
	IsClosed     bool
	IsPinfu      bool
	IsChiitoitsu bool
	SeatWind     engine.Seat
	RoundWind    engine.Wind
}

// WaitType represents the type of wait
type WaitType int

const (
	WaitRyanmen WaitType = iota // Open wait (23 waiting for 1 or 4)
	WaitKanchan                 // Closed wait (13 waiting for 2)
	WaitPenchan                 // Edge wait (12 waiting for 3, or 89 waiting for 7)
	WaitTanki                   // Pair wait (single tile wait)
	WaitShanpon                 // Double pon wait (two pairs, waiting for either to make triplet)
)

// CalculateFu calculates the fu (base points) for a hand
func CalculateFu(context FuContext) int {
	// Special cases
	if context.IsChiitoitsu {
		return 25 // Chiitoitsu is always 25 fu
	}

	if context.IsPinfu && context.IsTsumo {
		return 20 // Pinfu tsumo is always 20 fu (no rounding)
	}

	// Start with base fu
	fu := 20

	// Add for win method
	if context.IsTsumo {
		if context.IsClosed {
			fu += 2 // Closed tsumo
		}
		// Open tsumo has no fu bonus
	} else {
		fu += 10 // Ron (open or closed)
	}

	// Add for melds
	for _, meld := range context.Melds {
		if meld.Type == MeldPair {
			continue // Pairs are handled separately
		}

		fu += calculateMeldFu(meld, context.IsClosed)
	}

	// Add for pair
	fu += calculatePairFu(context.Pair, context.SeatWind, context.RoundWind)

	// Add for wait type
	fu += calculateWaitFu(context.WaitType)

	// Round up to nearest 10
	fu = ((fu + 9) / 10) * 10

	return fu
}

// calculateMeldFu calculates fu for a single meld
func calculateMeldFu(meld Meld, isClosed bool) int {
	if meld.Type == MeldSequence {
		return 0 // Sequences have no fu
	}

	if len(meld.Tiles) == 0 {
		return 0
	}

	tile := meld.Tiles[0]
	isTerminalOrHonor := tile.IsTerminal() || tile.IsHonor()

	if meld.Type == MeldTriplet {
		if isClosed {
			// Closed triplet (anko)
			if isTerminalOrHonor {
				return 8
			}
			return 4
		} else {
			// Open triplet (minko/pon)
			if isTerminalOrHonor {
				return 4
			}
			return 2
		}
	}

	if meld.Type == MeldQuad {
		if isClosed {
			// Closed quad (ankan)
			if isTerminalOrHonor {
				return 32
			}
			return 16
		} else {
			// Open quad (minkan/kakan)
			if isTerminalOrHonor {
				return 16
			}
			return 8
		}
	}

	return 0
}

// calculatePairFu calculates fu for the pair
func calculatePairFu(pair Meld, seatWind engine.Seat, roundWind engine.Wind) int {
	if len(pair.Tiles) == 0 {
		return 0
	}

	tile := pair.Tiles[0]

	// Only honors can be value pairs
	if tile.Suit != engine.Ji {
		return 0
	}

	fu := 0

	// Dragons (5z-7z) are always worth 2 fu
	if tile.Rank >= 5 && tile.Rank <= 7 {
		return 2
	}

	// Winds (1z-4z) - check if seat or round wind
	if tile.Rank >= 1 && tile.Rank <= 4 {
		// Convert tile rank to wind (1z=East=0, 2z=South=1, etc.)
		tileWind := tile.Rank - 1

		if int(seatWind) == tileWind {
			fu += 2
		}

		if int(roundWind) == tileWind {
			fu += 2
		}
	}

	return fu
}

// calculateWaitFu calculates fu for wait type
func calculateWaitFu(waitType WaitType) int {
	switch waitType {
	case WaitKanchan, WaitPenchan, WaitTanki:
		return 2
	default:
		return 0
	}
}

// AddHonbaAndRiichi adds honba counter and riichi stick bonuses to the score
func AddHonbaAndRiichi(score Score, honbaCount int, riichiSticks int) Score {
	// Honba: 300 points per counter
	honbaBonus := honbaCount * 300

	// Riichi sticks: 1000 points per stick
	riichiBonus := riichiSticks * 1000

	score.TotalPoints += honbaBonus + riichiBonus

	if score.RonPayment > 0 {
		score.RonPayment += honbaBonus + riichiBonus
	}

	return score
}

// PaoPayment calculates payment when pao (responsibility) applies
type PaoPayment struct {
	PaoPlayer string
	Payment   int
}

// CalculatePaoPayment calculates pao responsibility for yakuman
func CalculatePaoPayment(yakuman string, paoPlayer string, winMethod string, isDealer bool) PaoPayment {
	payment := PaoPayment{
		PaoPlayer: paoPlayer,
	}

	// Determine yakuman value
	baseYakuman := 32000 // Non-dealer yakuman
	if isDealer {
		baseYakuman = 48000 // Dealer yakuman
	}

	// Double yakuman
	if yakuman == "DAISUUSHII" {
		baseYakuman *= 2
	}

	// Pao player pays full amount regardless of win method
	payment.Payment = baseYakuman

	return payment
}

// UmaOkaCalculation calculates final scores with uma and oka
type UmaOkaCalculation struct {
	RawScores      []int
	StartingPoints int
	UmaValues      []int
	FinalScores    []int
}

// CalculateUmaOka applies uma and oka to final scores
func CalculateUmaOka(rawScores []int, startingPoints int, umaValues []int) []int {
	// Calculate raw score deltas from starting points
	deltas := make([]int, len(rawScores))
	for i, score := range rawScores {
		deltas[i] = score - startingPoints
	}

	// Rank players by score (highest to lowest)
	rankings := rankPlayers(rawScores)

	// Apply uma
	finalScores := make([]int, len(rawScores))
	for i, playerIdx := range rankings {
		finalScores[playerIdx] = deltas[playerIdx] + umaValues[i]
	}

	return finalScores
}

// rankPlayers returns player indices sorted by score (highest first)
func rankPlayers(scores []int) []int {
	rankings := make([]int, len(scores))
	for i := range rankings {
		rankings[i] = i
	}

	// Simple bubble sort by score
	for i := 0; i < len(rankings)-1; i++ {
		for j := i + 1; j < len(rankings); j++ {
			if scores[rankings[j]] > scores[rankings[i]] {
				rankings[i], rankings[j] = rankings[j], rankings[i]
			}
		}
	}

	return rankings
}
