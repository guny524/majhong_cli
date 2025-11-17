package engine

// Round represents a single round of mahjong
type Round struct {
	Dealer       Seat
	RoundWind    Wind
	RoundNumber  int
	Honba        int // 본장 counter
	RiichiSticks int // Riichi sticks on table
}

// NewRound creates a new round
func NewRound(dealer Seat, roundWind Wind, roundNumber int) *Round {
	return &Round{
		Dealer:       dealer,
		RoundWind:    roundWind,
		RoundNumber:  roundNumber,
		Honba:        0,
		RiichiSticks: 0,
	}
}

// NextDealer returns the next dealer after a round
func (r *Round) NextDealer() Seat {
	return (r.Dealer + 1) % 4
}

// IsAllLast checks if this is the last round (South 4 in East-South game)
func (r *Round) IsAllLast() bool {
	// South 4 is typically the last round
	return r.RoundWind == WindSouth && r.RoundNumber == 4
}
