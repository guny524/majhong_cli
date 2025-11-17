package engine

import "fmt"

// Seat represents player position
type Seat int

const (
	SeatEast Seat = iota
	SeatSouth
	SeatWest
	SeatNorth
)

func (s Seat) String() string {
	seats := []string{"East", "South", "West", "North"}
	return seats[s]
}

// Wind represents round wind
type Wind int

const (
	WindEast Wind = iota
	WindSouth
	WindWest
	WindNorth
)

func (w Wind) String() string {
	winds := []string{"East", "South", "West", "North"}
	return winds[w]
}

// MeldType represents the type of meld (called tiles)
type MeldType int

const (
	MeldChi MeldType = iota   // 치 - sequence (e.g., 3-4-5)
	MeldPon                   // 퐁 - triplet (e.g., 3-3-3)
	MeldAnkan                 // 암깡 - concealed quad (4 tiles, closed)
	MeldMinkan                // 명깡 - open quad from discard
	MeldKakan                 // 가깡 - added quad (add 4th to pon)
)

func (m MeldType) String() string {
	types := []string{"Chi", "Pon", "Ankan", "Minkan", "Kakan"}
	return types[m]
}

// Meld represents a called set of tiles
type Meld struct {
	Type       MeldType
	Tiles      []Tile
	FromPlayer Seat // Which player discarded the tile (for pon/chi/minkan)
	IsConcealed bool // Ankan is concealed even though it's a kan
}

// NewChi creates a chi meld
func NewChi(tiles []Tile, fromPlayer Seat) Meld {
	return Meld{
		Type:        MeldChi,
		Tiles:       tiles,
		FromPlayer:  fromPlayer,
		IsConcealed: false,
	}
}

// NewPon creates a pon meld
func NewPon(tiles []Tile, fromPlayer Seat) Meld {
	return Meld{
		Type:        MeldPon,
		Tiles:       tiles,
		FromPlayer:  fromPlayer,
		IsConcealed: false,
	}
}

// NewAnkan creates a concealed kan
func NewAnkan(tiles []Tile) Meld {
	return Meld{
		Type:        MeldAnkan,
		Tiles:       tiles,
		FromPlayer:  SeatEast, // Not used for ankan
		IsConcealed: true,
	}
}

// NewMinkan creates an open kan from discard
func NewMinkan(tiles []Tile, fromPlayer Seat) Meld {
	return Meld{
		Type:        MeldMinkan,
		Tiles:       tiles,
		FromPlayer:  fromPlayer,
		IsConcealed: false,
	}
}

// NewKakan creates an added kan (4th tile to pon)
func NewKakan(tiles []Tile, fromPlayer Seat) Meld {
	return Meld{
		Type:        MeldKakan,
		Tiles:       tiles,
		FromPlayer:  fromPlayer,
		IsConcealed: false,
	}
}

// FuritenState tracks three types of furiten
type FuritenState struct {
	DiscardFuriten   bool // Permanent - wait tile in own discards
	TemporaryFuriten bool // Until next turn - skipped ron this round
	RiichiFuriten    bool // Permanent - skipped any win after riichi
}

// IsFuriten checks if any furiten condition is active
func (f FuritenState) IsFuriten() bool {
	return f.DiscardFuriten || f.TemporaryFuriten || f.RiichiFuriten
}

// Player represents a mahjong player
type Player struct {
	Seat              Seat
	Hand              []Tile // Concealed tiles in hand
	Discards          []Tile // Tiles discarded by this player
	Melds             []Meld // Called melds (chi/pon/kan)
	Score             int
	IsRiichi          bool
	Furiten           FuritenState
	IsTenpai          bool
	DrawnTile         *Tile // Most recently drawn tile (kept separate)
	RiichiTurn        int   // Turn when riichi was declared (-1 if not in riichi)
	SkippedWinTiles   []Tile // Tiles skipped after riichi (for riichi furiten)
}

// NewPlayer creates a new player
func NewPlayer(seat Seat) *Player {
	return &Player{
		Seat:            seat,
		Hand:            make([]Tile, 0, 14),
		Discards:        make([]Tile, 0),
		Melds:           make([]Meld, 0),
		Score:           25000, // Default starting score
		IsRiichi:        false,
		Furiten:         FuritenState{},
		IsTenpai:        false,
		DrawnTile:       nil,
		RiichiTurn:      -1,
		SkippedWinTiles: make([]Tile, 0),
	}
}

// IsClosed checks if hand is closed (menzen)
func (p *Player) IsClosed() bool {
	if p.IsRiichi {
		return true // Riichi is always closed
	}

	// Check for open melds
	for _, meld := range p.Melds {
		if !meld.IsConcealed {
			return false
		}
	}
	return true
}

// TileCount returns total number of tiles (hand + drawn)
func (p *Player) TileCount() int {
	count := len(p.Hand)
	if p.DrawnTile != nil {
		count++
	}
	// Melds are already removed from hand, so don't count them separately
	// (they're shown separately but part of the 14-tile hand)
	return count
}

// AllTiles returns all tiles including hand, drawn, and melds
func (p *Player) AllTiles() []Tile {
	tiles := make([]Tile, 0, 14)
	tiles = append(tiles, p.Hand...)
	if p.DrawnTile != nil {
		tiles = append(tiles, *p.DrawnTile)
	}
	for _, meld := range p.Melds {
		tiles = append(tiles, meld.Tiles...)
	}
	return tiles
}

// Draw adds a tile to the drawn tile slot
func (p *Player) Draw(tile Tile) {
	p.DrawnTile = &tile
}

// Discard removes a tile and adds it to discard pile
func (p *Player) Discard(tile Tile) error {
	// Try to discard from drawn tile first
	if p.DrawnTile != nil && p.DrawnTile.Equals(tile) {
		p.Discards = append(p.Discards, tile)
		p.DrawnTile = nil
		return nil
	}

	// Otherwise discard from hand
	for i, t := range p.Hand {
		if t.Equals(tile) {
			p.Hand = append(p.Hand[:i], p.Hand[i+1:]...)
			p.Discards = append(p.Discards, tile)
			return nil
		}
	}

	return fmt.Errorf("tile %s not in hand", tile)
}

// AddMeld adds a meld to the player's melds
func (p *Player) AddMeld(meld Meld) {
	p.Melds = append(p.Melds, meld)
}

// RemoveTilesFromHand removes specified tiles from hand (for calling melds)
func (p *Player) RemoveTilesFromHand(tiles []Tile) error {
	handCopy := make([]Tile, len(p.Hand))
	copy(handCopy, p.Hand)

	for _, tileToRemove := range tiles {
		found := false
		for i, handTile := range handCopy {
			if handTile.Equals(tileToRemove) {
				handCopy = append(handCopy[:i], handCopy[i+1:]...)
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("tile %s not found in hand", tileToRemove)
		}
	}

	p.Hand = handCopy
	return nil
}

// DeclareRiichi declares riichi
func (p *Player) DeclareRiichi(turn int) error {
	if !p.IsClosed() {
		return fmt.Errorf("cannot riichi: hand is open")
	}
	if p.Score < 1000 {
		return fmt.Errorf("cannot riichi: insufficient points (need 1000, have %d)", p.Score)
	}

	p.IsRiichi = true
	p.RiichiTurn = turn
	p.Score -= 1000 // Pay riichi stick
	return nil
}

// String returns player description
func (p *Player) String() string {
	return fmt.Sprintf("Player %s (Score: %d, Hand: %d tiles, Melds: %d, Riichi: %v)",
		p.Seat, p.Score, p.TileCount(), len(p.Melds), p.IsRiichi)
}
