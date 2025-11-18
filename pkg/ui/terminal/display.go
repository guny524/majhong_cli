package terminal

import (
	"fmt"
	"strings"

	"github.com/guny524/majhong_cli/pkg/mahjong/engine"
)

// Display renders the game state to terminal
type Display struct {
	width  int
	height int
}

// NewDisplay creates a new terminal display
func NewDisplay() *Display {
	return &Display{
		width:  80,
		height: 40,
	}
}

// RenderGameState renders the complete game state
func (d *Display) RenderGameState(game *engine.GameState, playerSeat engine.Seat) {
	d.Clear()
	d.PrintHeader(game)
	d.PrintScores(game)
	d.PrintDora(game)
	d.PrintOtherPlayersDiscards(game, playerSeat)
	d.PrintPlayerHand(game, playerSeat)
	d.PrintActions(game, playerSeat)
}

// Clear clears the terminal screen
func (d *Display) Clear() {
	fmt.Print("\033[H\033[2J") // ANSI escape code to clear screen
}

// PrintHeader prints game header information
func (d *Display) PrintHeader(game *engine.GameState) {
	fmt.Println(strings.Repeat("=", d.width))
	fmt.Printf("  RIICHI MAHJONG - %s Round %d\n", game.Round.RoundWind, game.Round.RoundNumber)
	fmt.Printf("  Turn: %d | Wall Remaining: %d | Riichi Sticks: %d | Honba: %d\n",
		game.TurnNumber, game.Wall.Remaining(), game.Round.RiichiSticks, game.Round.Honba)
	fmt.Println(strings.Repeat("=", d.width))
	fmt.Println()
}

// PrintScores prints all player scores
func (d *Display) PrintScores(game *engine.GameState) {
	fmt.Println("SCORES:")
	for i, player := range game.Players {
		seat := engine.Seat(i)
		marker := " "
		if seat == game.CurrentPlayer {
			marker = "►"
		}
		riichiStatus := ""
		if player.IsRiichi {
			riichiStatus = " [RIICHI]"
		}
		fmt.Printf("  %s %s: %6d%s\n", marker, seat, player.Score, riichiStatus)
	}
	fmt.Println()
}

// PrintDora prints dora indicators
func (d *Display) PrintDora(game *engine.GameState) {
	indicators := game.Wall.DoraIndicators()
	fmt.Printf("DORA INDICATORS: %s\n", formatTiles(indicators))
	fmt.Println()
}

// PrintOtherPlayersDiscards prints discards for other players
func (d *Display) PrintOtherPlayersDiscards(game *engine.GameState, playerSeat engine.Seat) {
	fmt.Println("OTHER PLAYERS' DISCARDS:")
	for i := range game.Players {
		seat := engine.Seat(i)
		if seat == playerSeat {
			continue
		}
		player := game.Players[i]
		fmt.Printf("  %s: %s\n", seat, formatTiles(player.Discards))
	}
	fmt.Println()
}

// PrintPlayerHand prints the current player's hand
func (d *Display) PrintPlayerHand(game *engine.GameState, playerSeat engine.Seat) {
	player := game.Players[playerSeat]

	fmt.Println("YOUR HAND:")
	fmt.Printf("  Seat: %s\n", playerSeat)

	// Print melds
	if len(player.Melds) > 0 {
		fmt.Print("  Melds: ")
		for _, meld := range player.Melds {
			fmt.Printf("[%s: %s] ", meld.Type, formatTiles(meld.Tiles))
		}
		fmt.Println()
	}

	// Print concealed tiles
	fmt.Printf("  Hand: %s", formatTiles(player.Hand))

	// Print drawn tile separately
	if player.DrawnTile != nil {
		fmt.Printf(" | Drew: %s", player.DrawnTile.String())
	}
	fmt.Println()

	// Print discards
	fmt.Printf("  Discards: %s\n", formatTiles(player.Discards))
	fmt.Println()
}

// PrintActions prints available actions
func (d *Display) PrintActions(game *engine.GameState, playerSeat engine.Seat) {
	if game.CurrentPlayer != playerSeat {
		fmt.Println("Waiting for other player's turn...")
		return
	}

	player := game.Players[playerSeat]

	fmt.Println("AVAILABLE ACTIONS:")
	if player.DrawnTile == nil {
		fmt.Println("  1. Draw tile (d)")
	} else {
		fmt.Println("  1. Discard tile (discard <tile>, e.g., 'discard 3m')")
		if player.IsClosed() && !player.IsRiichi {
			fmt.Println("  2. Declare riichi (riichi)")
		}
		fmt.Println("  3. Declare tsumo (tsumo) - if you have winning hand")
	}

	// TODO: Add call actions (pon, chi, kan, ron) based on game state
	fmt.Println()
}

// PrintWinResult prints win result
func (d *Display) PrintWinResult(game *engine.GameState, winner engine.Seat, yakuList []string, han int, fu int, score int) {
	fmt.Println(strings.Repeat("=", d.width))
	fmt.Printf("  %s WINS!\n", winner)
	fmt.Println(strings.Repeat("=", d.width))
	fmt.Println()
	fmt.Println("YAKU:")
	for _, yaku := range yakuList {
		fmt.Printf("  - %s\n", yaku)
	}
	fmt.Printf("\nHan: %d, Fu: %d\n", han, fu)
	fmt.Printf("Score: %d points\n", score)
	fmt.Println()
	d.PrintScores(game)
}

// PrintError prints error message
func (d *Display) PrintError(err error) {
	fmt.Printf("ERROR: %v\n", err)
}

// Helper functions

// formatTiles formats tiles for display
func formatTiles(tiles []engine.Tile) string {
	if len(tiles) == 0 {
		return "(none)"
	}

	parts := make([]string, len(tiles))
	for i, tile := range tiles {
		parts[i] = tile.String()
	}
	return strings.Join(parts, " ")
}
