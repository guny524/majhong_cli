package main

import (
	"fmt"
	"os"
	"time"

	"github.com/alecthomas/kong"
)

var version = "0.1.0"

// CLI represents the command-line interface
type CLI struct {
	// Batch mode flags
	GameFile     string `help:"Path to game state JSON file" short:"f" type:"path"`
	Action       string `help:"Single action to execute (e.g., 'draw', 'discard 3m')" short:"a"`
	Actions      string `help:"JSON array of actions to execute" short:"A"`
	Init         bool   `help:"Initialize new game" short:"i"`
	Query        string `help:"Query game state (current-player, valid-actions, scores, round, wall-remaining)" short:"q"`
	Seed         int64  `help:"Seed for random number generation (default: current time)" short:"s" default:"0"`
	JSON         bool   `help:"Output in JSON format" short:"j"`
	ExitOnTamper bool   `help:"Exit with code 4 on wall tampering" default:"false"`

	// Interactive mode
	Interactive bool `help:"Start interactive mode (TUI)" short:"I" default:"false"`

	// Global flags
	Version bool `help:"Print version information" short:"v"`
}

func main() {
	var cli CLI

	kong.Parse(&cli,
		kong.Name("majhong_cli"),
		kong.Description("Mahjong game engine CLI for AI training and gameplay"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)

	// Handle version flag
	if cli.Version {
		fmt.Printf("majhong_cli version %s\n", version)
		os.Exit(ExitSuccess)
	}

	// Use current time as seed if not specified
	if cli.Seed == 0 {
		cli.Seed = time.Now().UnixNano()
	}

	// Handle interactive mode
	if cli.Interactive {
		runInteractiveMode(cli.Seed)
		return
	}

	// Validate batch mode arguments
	if !cli.Init && cli.Query == "" && cli.Action == "" && cli.Actions == "" {
		// No batch mode flags specified - default to interactive mode
		fmt.Println("No batch mode flags specified. Starting interactive mode...")
		fmt.Println("Use --help to see all available options.")
		fmt.Println()
		runInteractiveMode(cli.Seed)
		return
	}

	// Create batch mode handler
	batch := &BatchMode{
		GameFile:     cli.GameFile,
		Action:       cli.Action,
		Actions:      cli.Actions,
		Init:         cli.Init,
		Query:        cli.Query,
		Seed:         cli.Seed,
		JSON:         cli.JSON,
		ExitOnTamper: cli.ExitOnTamper,
	}

	// Run batch mode
	exitCode := batch.Run()
	os.Exit(exitCode)
}

func runInteractiveMode(seed int64) {
	fmt.Println("================================================================================")
	fmt.Println("                    WELCOME TO RIICHI MAHJONG CLI")
	fmt.Println("================================================================================")
	fmt.Printf("\nInteractive mode will be available after merging branch 001-cli-mahjong-engine\n")
	fmt.Printf("\nFor now, use batch mode:\n")
	fmt.Printf("  --init --game-file game.json --seed %d\n", seed)
	fmt.Printf("  --query current-player --game-file game.json\n")
	fmt.Printf("  --action \"draw\" --game-file game.json\n")
	fmt.Printf("\nSee --help for all batch mode options.\n")
	os.Exit(0)
}
