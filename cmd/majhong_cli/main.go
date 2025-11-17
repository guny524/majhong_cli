package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/guny524/majhong_cli/pkg/game"
)

func main() {
	// Parse command-line flags
	var (
		seed        = flag.Int64("seed", 0, "Random seed for wall shuffling (0 for current time)")
		interactive = flag.Bool("interactive", true, "Run in interactive mode")
		help        = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		printHelp()
		return
	}

	// Use current time as seed if not specified
	if *seed == 0 {
		*seed = time.Now().UnixNano()
	}

	if *interactive {
		runInteractiveMode(*seed)
	} else {
		fmt.Println("Batch mode not yet implemented. Use --interactive flag.")
		os.Exit(1)
	}
}

func runInteractiveMode(seed int64) {
	fmt.Println("================================================================================")
	fmt.Println("                    WELCOME TO RIICHI MAHJONG CLI")
	fmt.Println("================================================================================")
	fmt.Printf("\nStarting new game with seed: %d\n", seed)
	fmt.Println("\nYou are playing as East (Dealer).")
	fmt.Println("The other three players are AI-controlled.")
	fmt.Println("\nType 'help' for available commands.")
	fmt.Println("Press Enter to continue...")
	fmt.Scanln()

	// Create and run interactive game
	ig := game.NewInteractiveGame(seed)
	if err := ig.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println("Majhong CLI - Riichi Mahjong Game Engine")
	fmt.Println("\nUSAGE:")
	fmt.Println("  majhong_cli [options]")
	fmt.Println("\nOPTIONS:")
	fmt.Println("  --interactive     Run in interactive mode (default: true)")
	fmt.Println("  --seed <number>   Random seed for wall shuffling (default: current time)")
	fmt.Println("  --help            Show this help")
	fmt.Println("\nEXAMPLES:")
	fmt.Println("  # Start interactive game with random seed")
	fmt.Println("  majhong_cli")
	fmt.Println()
	fmt.Println("  # Start game with specific seed for reproducibility")
	fmt.Println("  majhong_cli --seed 12345")
	fmt.Println()
	fmt.Println("\nIN-GAME COMMANDS:")
	fmt.Println("  Type 'help' during the game to see available commands.")
	fmt.Println()
}
