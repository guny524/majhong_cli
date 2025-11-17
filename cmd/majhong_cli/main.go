package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alecthomas/kong"
	"github.com/guny524/majhong_cli/pkg/mahjong/engine"
	"github.com/guny524/majhong_cli/pkg/persistence"
)

var CLI struct {
	Play PlayCmd `cmd:"" help:"Start a new game"`
	Save SaveCmd `cmd:"" help:"Save current game state"`
	Load LoadCmd `cmd:"" help:"Load a saved game"`
	List ListCmd `cmd:"" help:"List saved games"`
}

// PlayCmd starts a new game
type PlayCmd struct {
	Seed       int64  `help:"Random seed for game generation" default:"0"`
	Save       string `help:"Auto-save game to this file"`
	AkaDora    bool   `help:"Enable aka dora (red fives)" default:"true"`
	Points     int    `help:"Starting points" default:"25000" enum:"25000,30000,35000"`
	LoadFile   string `help:"Load game from file" optional:""`
}

func (cmd *PlayCmd) Run() error {
	storage := persistence.NewFileStorage(getSaveDir())

	var game *engine.GameState
	var err error

	if cmd.LoadFile != "" {
		// Load existing game
		fmt.Printf("Loading game from %s...\n", cmd.LoadFile)
		game, err = storage.Load(cmd.LoadFile)
		if err != nil {
			return fmt.Errorf("failed to load game: %w", err)
		}
		fmt.Println("Game loaded successfully!")
	} else {
		// Create new game
		seed := cmd.Seed
		if seed == 0 {
			seed = time.Now().UnixNano()
		}

		config := engine.RuleConfig{
			RiichiMahjong:  true,
			StartingPoints: cmd.Points,
			ReturnPoints:   30000,
			AkaDora:        cmd.AkaDora,
			LocalYaku:      []string{},
		}

		fmt.Printf("Starting new game with seed %d...\n", seed)
		game = engine.NewGame(config, seed)
		fmt.Println("Game created successfully!")
	}

	// Display game state
	displayGameState(game)

	// Auto-save if specified
	if cmd.Save != "" {
		fmt.Printf("\nSaving game to %s...\n", cmd.Save)
		if err := storage.Save(game, cmd.Save); err != nil {
			return fmt.Errorf("failed to save game: %w", err)
		}
		fmt.Println("Game saved successfully!")
	}

	return nil
}

// SaveCmd saves a game state
type SaveCmd struct {
	Input  string `arg:"" help:"Input game file to load"`
	Output string `arg:"" help:"Output file name"`
}

func (cmd *SaveCmd) Run() error {
	storage := persistence.NewFileStorage(getSaveDir())

	// Load the game
	fmt.Printf("Loading game from %s...\n", cmd.Input)
	game, err := storage.Load(cmd.Input)
	if err != nil {
		return fmt.Errorf("failed to load game: %w", err)
	}

	// Save to new file
	fmt.Printf("Saving game to %s...\n", cmd.Output)
	if err := storage.Save(game, cmd.Output); err != nil {
		return fmt.Errorf("failed to save game: %w", err)
	}

	fmt.Println("Game saved successfully!")
	return nil
}

// LoadCmd loads and displays a game state
type LoadCmd struct {
	File    string `arg:"" help:"File name to load"`
	Verbose bool   `help:"Show detailed information" short:"v"`
}

func (cmd *LoadCmd) Run() error {
	storage := persistence.NewFileStorage(getSaveDir())

	fmt.Printf("Loading game from %s...\n", cmd.File)
	game, err := storage.Load(cmd.File)
	if err != nil {
		return fmt.Errorf("failed to load game: %w", err)
	}

	fmt.Println("Game loaded successfully!")
	fmt.Println()

	displayGameState(game)

	if cmd.Verbose {
		fmt.Println("\n=== Action History ===")
		for i, action := range game.History {
			fmt.Printf("%d. %s: %s (tiles: %d)\n",
				i+1,
				seatToString(action.Player),
				actionTypeToString(action.Type),
				len(action.Tiles))
		}
	}

	return nil
}

// ListCmd lists all saved games
type ListCmd struct {
	Verbose bool `help:"Show detailed information" short:"v"`
}

func (cmd *ListCmd) Run() error {
	storage := persistence.NewFileStorage(getSaveDir())

	files, err := storage.List()
	if err != nil {
		return fmt.Errorf("failed to list saves: %w", err)
	}

	if len(files) == 0 {
		fmt.Println("No saved games found")
		return nil
	}

	fmt.Printf("Found %d saved game(s):\n\n", len(files))

	for _, file := range files {
		if cmd.Verbose {
			info, err := storage.GetInfo(file)
			if err != nil {
				fmt.Printf("  - %s (error: %v)\n", file, err)
				continue
			}
			fmt.Printf("  - %s\n", file)
			fmt.Printf("    Version: %s\n", info.Version)
			fmt.Printf("    Saved: %s\n", info.Timestamp.Format("2006-01-02 15:04:05"))
			fmt.Printf("    Size: %d bytes\n", info.Size)
			fmt.Println()
		} else {
			fmt.Printf("  - %s\n", file)
		}
	}

	return nil
}

func main() {
	ctx := kong.Parse(&CLI,
		kong.Name("majhong_cli"),
		kong.Description("Mahjong game state management CLI"),
		kong.UsageOnError(),
	)

	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}

// Helper functions

func getSaveDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".majhong_saves"
	}
	return filepath.Join(homeDir, ".majhong_saves")
}

func displayGameState(game *engine.GameState) {
	fmt.Println("\n=== Game State ===")
	fmt.Printf("Round: %s %d-%d\n",
		windToString(game.Round.RoundWind),
		game.Round.RoundNumber,
		game.Round.Honba)
	fmt.Printf("Dealer: %s\n", seatToString(game.Round.Dealer))
	fmt.Printf("Current Player: %s\n", seatToString(game.CurrentPlayer))
	fmt.Printf("Turn: %d\n", game.TurnNumber)
	fmt.Printf("Wall Remaining: %d tiles\n", len(game.Wall.LiveWall))
	fmt.Printf("Riichi Sticks: %d\n", game.Round.RiichiSticks)
	fmt.Printf("Game Over: %v\n", game.IsGameOver)

	fmt.Println("\n=== Players ===")
	for i, player := range game.Players {
		fmt.Printf("\n%s:", seatToString(engine.Seat(i)))
		fmt.Printf("\n  Score: %d", player.Score)
		fmt.Printf("\n  Hand: %d tiles", len(player.Hand))
		fmt.Printf("\n  Discards: %d tiles", len(player.Discards))
		fmt.Printf("\n  Melds: %d", len(player.Melds))
		fmt.Printf("\n  Riichi: %v", player.IsRiichi)
		fmt.Printf("\n  Tenpai: %v", player.IsTenpai)
		if player.Furiten.DiscardFuriten || player.Furiten.TemporaryFuriten || player.Furiten.RiichiFuriten {
			fmt.Printf("\n  Furiten: yes")
		}
	}
	fmt.Println()
}

func seatToString(seat engine.Seat) string {
	switch seat {
	case engine.SeatEast:
		return "East"
	case engine.SeatSouth:
		return "South"
	case engine.SeatWest:
		return "West"
	case engine.SeatNorth:
		return "North"
	default:
		return "Unknown"
	}
}

func windToString(wind engine.Wind) string {
	switch wind {
	case engine.WindEast:
		return "East"
	case engine.WindSouth:
		return "South"
	case engine.WindWest:
		return "West"
	case engine.WindNorth:
		return "North"
	default:
		return "Unknown"
	}
}

func actionTypeToString(actionType engine.ActionType) string {
	switch actionType {
	case engine.ActionDraw:
		return "Draw"
	case engine.ActionDiscard:
		return "Discard"
	case engine.ActionPon:
		return "Pon"
	case engine.ActionChi:
		return "Chi"
	case engine.ActionKan:
		return "Kan"
	case engine.ActionRiichi:
		return "Riichi"
	case engine.ActionTsumo:
		return "Tsumo"
	case engine.ActionRon:
		return "Ron"
	default:
		return "Unknown"
	}
}
