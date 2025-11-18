package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/guny524/majhong_cli/pkg/mahjong/engine"
)

// Parser handles user input parsing
type Parser struct {
	reader *bufio.Reader
}

// NewParser creates a new input parser
func NewParser() *Parser {
	return &Parser{
		reader: bufio.NewReader(os.Stdin),
	}
}

// ReadInput reads a line of input from user
func (p *Parser) ReadInput() (string, error) {
	fmt.Print("> ")
	input, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// ParseAction parses user input into an action
func (p *Parser) ParseAction(input string, game *engine.GameState) (engine.Action, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return engine.Action{}, fmt.Errorf("empty input")
	}

	command := strings.ToLower(parts[0])
	player := game.CurrentPlayer

	switch command {
	case "d", "draw":
		return engine.NewAction(engine.ActionDraw, player, nil, game.TurnNumber), nil

	case "discard":
		if len(parts) < 2 {
			return engine.Action{}, fmt.Errorf("discard requires tile (e.g., 'discard 3m')")
		}
		tile, err := engine.ParseTile(parts[1])
		if err != nil {
			return engine.Action{}, fmt.Errorf("invalid tile: %v", err)
		}
		return engine.NewAction(engine.ActionDiscard, player, []engine.Tile{tile}, game.TurnNumber), nil

	case "riichi":
		return engine.NewAction(engine.ActionRiichi, player, nil, game.TurnNumber), nil

	case "tsumo":
		return engine.NewAction(engine.ActionTsumo, player, nil, game.TurnNumber), nil

	case "ron":
		return engine.NewAction(engine.ActionRon, player, nil, game.TurnNumber), nil

	case "pon":
		if len(parts) < 2 {
			return engine.Action{}, fmt.Errorf("pon requires tiles (e.g., 'pon 3m')")
		}
		// Parse 3 identical tiles for pon
		tile, err := engine.ParseTile(parts[1])
		if err != nil {
			return engine.Action{}, fmt.Errorf("invalid tile: %v", err)
		}
		tiles := []engine.Tile{tile, tile, tile}
		return engine.NewAction(engine.ActionPon, player, tiles, game.TurnNumber), nil

	case "chi":
		if len(parts) < 2 {
			return engine.Action{}, fmt.Errorf("chi requires tiles (e.g., 'chi 3m 4m 5m')")
		}
		if len(parts) < 4 {
			return engine.Action{}, fmt.Errorf("chi requires 3 tiles")
		}
		tiles := make([]engine.Tile, 3)
		for i := 0; i < 3; i++ {
			tile, err := engine.ParseTile(parts[i+1])
			if err != nil {
				return engine.Action{}, fmt.Errorf("invalid tile: %v", err)
			}
			tiles[i] = tile
		}
		return engine.NewAction(engine.ActionChi, player, tiles, game.TurnNumber), nil

	case "kan":
		if len(parts) < 2 {
			return engine.Action{}, fmt.Errorf("kan requires tiles (e.g., 'kan 3m')")
		}
		tile, err := engine.ParseTile(parts[1])
		if err != nil {
			return engine.Action{}, fmt.Errorf("invalid tile: %v", err)
		}
		tiles := []engine.Tile{tile, tile, tile, tile}
		return engine.NewAction(engine.ActionKan, player, tiles, game.TurnNumber), nil

	case "help", "h", "?":
		p.PrintHelp()
		return engine.Action{}, fmt.Errorf("help displayed")

	case "quit", "q", "exit":
		return engine.Action{}, fmt.Errorf("quit")

	default:
		return engine.Action{}, fmt.Errorf("unknown command: %s (type 'help' for commands)", command)
	}
}

// PrintHelp prints help information
func (p *Parser) PrintHelp() {
	fmt.Println("\nCOMMANDS:")
	fmt.Println("  d, draw          - Draw a tile from the wall")
	fmt.Println("  discard <tile>   - Discard a tile (e.g., 'discard 3m', 'discard 1z')")
	fmt.Println("  riichi           - Declare riichi")
	fmt.Println("  tsumo            - Declare tsumo (self-draw win)")
	fmt.Println("  ron              - Declare ron (win on discard)")
	fmt.Println("  pon <tile>       - Call pon (e.g., 'pon 3m')")
	fmt.Println("  chi <t1> <t2> <t3> - Call chi (e.g., 'chi 3m 4m 5m')")
	fmt.Println("  kan <tile>       - Call kan (e.g., 'kan 3m')")
	fmt.Println("  help, h, ?       - Show this help")
	fmt.Println("  quit, q, exit    - Quit game")
	fmt.Println("\nTILE NOTATION:")
	fmt.Println("  1-9m = Man (characters)")
	fmt.Println("  1-9p = Pin (circles)")
	fmt.Println("  1-9s = Sou (bamboo)")
	fmt.Println("  1-7z = Ji (honors: 1-4 winds, 5-7 dragons)")
	fmt.Println("  Example: 3m = 3 man, 5p = 5 pin, 1z = East wind")
	fmt.Println()
}
