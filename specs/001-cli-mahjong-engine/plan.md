# Implementation Plan: CLI-Based Riichi Mahjong Game Engine

**Branch**: `001-cli-mahjong-engine` | **Date**: 2025-11-18 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-cli-mahjong-engine/spec.md`

## Summary

This plan details the implementation of a CLI-based Riichi Mahjong game engine with dual execution modes (interactive TUI and batch/AI training), supporting both JSON file-based persistence and high-performance Protobuf streaming for AI agents. The engine will run in Kubernetes container environments (headless-capable) and serve as the foundation for separate AI and web-integration modules.

**Core Capabilities:**
- Interactive 4-player Riichi Mahjong game with tcell-based TUI
- Batch mode for AI training with file-based state persistence (JSON primary, Protobuf optional)
- Majsoul paipu import/export for training data generation
- Read-only observer mode for live game watching
- Wall hash validation (SHA-256) with configurable tamper detection

**Technical Approach:**
- **Library-First Architecture**: Core game engine as standalone Go library (`pkg/mahjong/`)
- **Hybrid Serialization**: JSON for human-readability/debugging, Protobuf for network/performance
- **Headless Support**: Pure CLI operation, no graphical dependencies (tcell for terminal only)
- **Container-Ready**: Stateless design for horizontal scaling in Kubernetes

## Technical Context

**Language/Version**: Go 1.21+
**Primary Dependencies**:
- `alecthomas/kong` - CLI argument parsing
- `gdamore/tcell/v2` - Terminal UI (headless-compatible)
- `google/protobuf` - Binary serialization for network/fast mode
- Standard library: `encoding/json`, `crypto/sha256`, `math/rand`

**Storage**: Filesystem (JSON/Protobuf game state files)
**Testing**: Go standard testing framework (`go test`), table-driven tests for game rules
**Target Platform**: Linux/macOS/Windows CLI, Kubernetes containers (headless mode)
**Project Type**: Single Go binary with library-first architecture
**Performance Goals**:
- Interactive mode: <100ms response time per user input
- Batch mode: 1000+ game simulations per minute
- Fast agent mode (protobuf): 10,000+ actions per minute
- Furiten/yaku detection: 100% accuracy

**Constraints**:
- Headless operation (no OpenGL/graphics)
- Single-threaded game logic (concurrency at process level for AI training)
- JSON file format must be human-readable for debugging
- Protobuf schema must be forward/backward compatible
- Kubernetes-friendly (stateless, config via env/flags)

**Scale/Scope**:
- 52 standard yaku + scoring rules
- 136 tiles, 4 players per game
- Support 100+ concurrent games per container (batch mode)
- Game state file size: ~50KB JSON, ~10KB Protobuf

## Constitution Check

*GATE: Must pass before implementation begins*

### ✅ Compliance Status

| Principle | Requirement | Compliance | Notes |
|-----------|-------------|------------|-------|
| **I. Skill-First Development** | Review & activate skills before tasks | ✅ Pass | Will activate `clean-code`, `react-next-guide` (for future web), `doc-management` |
| **II. Documentation-First** | Create spec.md, plan.md before code | ✅ Pass | spec.md complete, plan.md in progress |
| **III. Test-Driven Development** | Write tests before implementation | ✅ Pass | Will create test files per phase (see Phase 0) |
| **IV. Library-First Architecture** | Core as standalone library | ✅ Pass | `pkg/mahjong/` as library, `cmd/majhong_cli/` as CLI wrapper |
| **V. CLI-First Interface** | stdin/stdout, JSON output | ✅ Pass | All modes support JSON output via `--json` flag |
| **VI. Evidence-Based Changes** | Verify before modifying | ✅ Pass | Will use git diff, file inspection for all changes |
| **VII. Minimal Complexity** | Simplest solution first | ✅ Pass | No ORM, no complex patterns, direct implementation |

### Complexity Justification

**None required** - Architecture follows constitution principles:
- Single Go project (not fragmented)
- Direct implementation (no unnecessary abstractions)
- Standard library preferred over external dependencies where possible

## Project Structure

### Documentation (this feature)

```text
specs/001-cli-mahjong-engine/
├── spec.md              # Requirements (completed)
├── plan.md              # This file - architecture & design
├── checklists/
│   └── requirements.md  # Validation checklist (completed)
└── tasks.md             # Task breakdown (created by /speckit.tasks)
```

### Source Code (repository root)

```text
majhong_cli/
├── cmd/
│   └── majhong_cli/         # CLI entry point (kong parser)
│       ├── main.go
│       ├── interactive.go   # Interactive TUI mode handler
│       ├── batch.go         # Batch/file mode handler
│       ├── observer.go      # Observer mode handler
│       └── fast_agent.go    # Protobuf streaming mode handler
│
├── pkg/
│   └── mahjong/             # Core game engine library
│       ├── engine/          # Game state & rule enforcement
│       │   ├── game.go
│       │   ├── round.go
│       │   ├── player.go
│       │   ├── wall.go
│       │   ├── tile.go
│       │   └── action.go
│       │
│       ├── rules/           # Rule implementations
│       │   ├── yaku.go          # Yaku recognition
│       │   ├── scoring.go       # Han/fu calculation
│       │   ├── furiten.go       # Furiten detection
│       │   ├── pao.go           # Responsibility payment
│       │   └── validation.go    # Action validation
│       │
│       ├── serialization/   # State persistence
│       │   ├── json.go          # JSON serializer
│       │   ├── protobuf.go      # Protobuf serializer
│       │   ├── converter.go     # JSON ↔ Protobuf
│       │   └── hash.go          # Wall hash validation
│       │
│       └── majsoul/         # Majsoul paipu integration
│           ├── importer.go      # Paipu → internal state
│           └── exporter.go      # Internal state → paipu
│
├── proto/                   # Protobuf schemas
│   ├── game_state.proto     # Game state message
│   ├── action.proto         # Action messages
│   └── majsoul.proto        # Majsoul paipu format
│
├── internal/
│   ├── tui/                 # Terminal UI (tcell)
│   │   ├── renderer.go          # Game state rendering
│   │   ├── input.go             # User input handling
│   │   └── layout.go            # TUI layout
│   │
│   └── config/              # Configuration
│       └── rules.go             # Rule set configuration
│
├── test/
│   ├── fixtures/            # Test data
│   │   ├── game_states/         # Sample game state JSON files
│   │   ├── majsoul_paipu/       # Sample Majsoul replays
│   │   └── yaku_examples/       # Known yaku patterns
│   │
│   ├── integration/         # Integration tests
│   │   ├── batch_mode_test.go
│   │   ├── interactive_test.go
│   │   └── majsoul_test.go
│   │
│   └── contract/            # Contract tests (TDD)
│       ├── yaku_test.go         # Test all 52 yaku
│       ├── scoring_test.go      # Score calculation
│       └── furiten_test.go      # Furiten scenarios
│
├── .specify/                # Speckit artifacts
│   └── memory/
│       ├── constitution.md
│       └── constitution_majhong_*.md
│
├── go.mod
├── go.sum
├── Makefile                 # Build, test, proto generation
├── Dockerfile               # Container image (headless)
└── README.md
```

**Structure Decision**: Single Go project with library-first architecture. The `pkg/mahjong/` library is the core engine (testable, reusable), while `cmd/majhong_cli/` provides CLI interface using kong. This enables:
1. Independent testing of game logic without CLI layer
2. Future reuse in web server (separate `majhong_server` binary)
3. AI module can import `pkg/mahjong/` as Go library
4. Clean separation: engine (domain) → serialization (infrastructure) → CLI (interface)

## Architecture Design

### System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLI Entry Point                          │
│                      (cmd/majhong_cli/)                          │
│                                                                   │
│  ┌──────────────┐  ┌──────────────┐  ┌─────────────────────┐   │
│  │ Interactive  │  │   Batch      │  │  Fast Agent Mode    │   │
│  │  TUI Mode    │  │  File Mode   │  │  (protobuf stream)  │   │
│  │  (tcell)     │  │  (JSON)      │  │                     │   │
│  └──────┬───────┘  └──────┬───────┘  └──────────┬──────────┘   │
│         │                  │                     │              │
│         └──────────────────┴─────────────────────┘              │
│                            │                                     │
└────────────────────────────┼─────────────────────────────────────┘
                             │
                    ┌────────▼────────┐
                    │ Serialization   │
                    │     Layer       │
                    │  JSON ↔ Proto   │
                    └────────┬────────┘
                             │
                    ┌────────▼────────────────────────────┐
                    │    Core Game Engine Library         │
                    │       (pkg/mahjong/)                │
                    │                                     │
                    │  ┌────────────┐  ┌─────────────┐   │
                    │  │   Game     │  │    Rules    │   │
                    │  │   State    │◄─┤  Yaku/Score│   │
                    │  │  Manager   │  │  Furiten    │   │
                    │  └────────────┘  └─────────────┘   │
                    │                                     │
                    │  ┌────────────┐  ┌─────────────┐   │
                    │  │   Wall     │  │  Majsoul    │   │
                    │  │   Hash     │  │  Importer   │   │
                    │  └────────────┘  └─────────────┘   │
                    └─────────────────────────────────────┘
```

### Module Breakdown

#### 1. Core Game Engine (`pkg/mahjong/engine/`)

**Responsibilities:**
- Maintain game state (tiles, players, scores, round)
- Enforce game rules (valid actions, turn order)
- Apply actions to state (immutable pattern)
- Detect game-ending conditions (win, draw, exhaustion)

**Key Components:**

```go
// game.go
type Game struct {
    State      GameState
    Config     RuleConfig
    History    []Action
    RNG        *rand.Rand  // Seeded for determinism
}

func (g *Game) ApplyAction(action Action) (*Game, error) {
    // Validate → Apply → Return new state
}

func (g *Game) ValidActions() []Action {
    // Calculate legal moves for current player
}

// wall.go
type Wall struct {
    LiveWall      []Tile    // 122 tiles
    DeadWall      DeadWall  // 14 tiles (dora + rinshan)
    Seed          int64     // For replay determinism
    HashSalt      string
    HashValue     string    // SHA-256(wall + salt)
}

func (w *Wall) Draw() (Tile, error) {
    // Draw from live wall, maintain integrity
}

func (w *Wall) ValidateHash() error {
    // Verify H(wall + salt) matches stored hash
}
```

**Design Patterns:**
- **Immutable State**: Actions return new `Game` struct (no mutation)
- **Value Objects**: Tile, Action are immutable structs
- **Builder Pattern**: `NewGame(config).WithSeed(seed).Build()`

#### 2. Rules Engine (`pkg/mahjong/rules/`)

**Responsibilities:**
- Yaku recognition (pattern matching)
- Scoring calculation (han/fu → points)
- Furiten detection (3 types)
- Pao responsibility determination

**Key Components:**

```go
// yaku.go
type YakuDetector struct {
    config RuleConfig
}

func (yd *YakuDetector) Detect(hand Hand, winTile Tile, context WinContext) []Yaku {
    // Pattern matching for all 52 yaku
    // Returns: [Riichi, Tanyao, ...] with han values
}

// scoring.go
func CalculateScore(yaku []Yaku, fu int, isDealer bool) ScoreResult {
    // Han/fu → base points → tsumo/ron distribution
    // Handles: honba, riichi sticks, pao
}

// furiten.go
func (g *Game) CheckFuriten(player Player) FuritenState {
    // Discard furiten: discarded a wait tile
    // Temporary furiten: declined ron
    // Riichi furiten: riichi + discard any wait
}
```

**Testing Strategy:**
- Table-driven tests with all 52 yaku examples from `test/fixtures/yaku_examples/`
- Property-based testing for score calculation consistency
- Furiten edge cases (e.g., discarding 3m when waiting on 3m/6m)

#### 3. Serialization Layer (`pkg/mahjong/serialization/`)

**Responsibilities:**
- Convert game state to/from JSON (human-readable)
- Convert game state to/from Protobuf (network/fast mode)
- Bidirectional JSON ↔ Protobuf conversion
- Wall hash generation and validation

**Key Components:**

```go
// json.go
type JSONSerializer struct{}

func (js *JSONSerializer) Marshal(game *Game) ([]byte, error) {
    // Game → JSON (pretty-printed)
    // Include: state, history, config, version
}

func (js *JSONSerializer) Unmarshal(data []byte) (*Game, error) {
    // JSON → Game (validate version)
}

// protobuf.go
type ProtoSerializer struct{}

func (ps *ProtoSerializer) Marshal(game *Game) ([]byte, error) {
    // Game → Protobuf binary
}

// converter.go
func JSONToProto(jsonData []byte) ([]byte, error) {
    // JSON → Game → Protobuf
}

func ProtoToJSON(protoData []byte) ([]byte, error) {
    // Protobuf → Game → JSON
}

// hash.go
func GenerateWallHash(wall Wall, salt string) string {
    // SHA-256(serialize(wall) + salt)
}

func ValidateWallHash(wall Wall, salt string, hash string) bool {
    // Compare computed hash with stored hash
}
```

**File Format Example (JSON):**

```json
{
  "version": "1.0.0",
  "seed": 1234567890,
  "round": "East-1",
  "dealer": "East",
  "honba": 0,
  "riichi_sticks": 0,
  "dora_indicators": ["5p"],
  "players": [
    {
      "seat": "East",
      "hand": ["1m", "2m", "3m", "7m", "7m", "1p", "2p", "3p", "5s", "6s", "7s", "1z", "1z"],
      "discards": ["9m", "8p"],
      "score": 25000,
      "riichi": false,
      "furiten": false
    },
    ...
  ],
  "wall": {
    "remaining": 68,
    "seed": 1234567890,
    "hash": "a3f2...",  // SHA-256(wall + salt)
    "salt_hash": "b4e1..."  // SHA-256(salt)
  },
  "history": [
    {"turn": 1, "player": "East", "action": "draw", "timestamp": "2025-11-18T10:00:00Z"},
    {"turn": 1, "player": "East", "action": "discard", "tile": "9m", "timestamp": "2025-11-18T10:00:05Z"}
  ]
}
```

#### 4. Majsoul Integration (`pkg/mahjong/majsoul/`)

**Responsibilities:**
- Parse Majsoul paipu JSON format
- Map Majsoul tile notation (0m/0p/0s for aka dora)
- Extract timeline events → Action sequence
- Export internal state to Majsoul format for replay viewers

**Key Components:**

```go
// importer.go
type PaipuImporter struct{}

func (pi *PaipuImporter) Import(paipuJSON []byte) (*Game, error) {
    // Majsoul paipu → internal Game state
    // Validates: wall hash, player names, final scores
}

func (pi *PaipuImporter) MapTile(majsoulTile string) Tile {
    // "0m" → Tile{Suit: Man, Rank: 5, IsAkaDora: true}
    // "5mr" → same as above (alternative notation)
}

// exporter.go
type PaipuExporter struct{}

func (pe *PaipuExporter) Export(game *Game) ([]byte, error) {
    // Internal state → Majsoul paipu JSON
    // Includes: wall hash, timeline, final scores
}
```

#### 5. CLI Layer (`cmd/majhong_cli/`)

**Responsibilities:**
- Parse command-line arguments (kong)
- Dispatch to appropriate mode (interactive/batch/observer/fast-agent)
- Format output (JSON or human-readable)
- Handle exit codes

**Kong CLI Structure:**

```go
// main.go
var CLI struct {
    // Interactive mode (default if no flags)
    Interactive struct {
        Load     string `arg optional name:"load" help:"Load saved game"`
        Config   string `flag name:"config" help:"Rule configuration file"`
    } `cmd help:"Start interactive TUI game"`

    // Batch mode (file-based)
    Batch struct {
        GameFile string `flag required name:"game-file" help:"Game state file path"`
        Action   string `flag name:"action" help:"Single action to apply"`
        Actions  string `flag name:"actions" help:"JSON array of actions (batch)"`
        Query    string `flag name:"query" help:"Query game state field"`
        Init     bool   `flag name:"init" help:"Initialize new game"`
        Seed     int64  `flag name:"seed" help:"RNG seed for determinism"`
        Save     string `flag name:"save" help:"Save game state to file"`
        JSON     bool   `flag name:"json" help:"Output JSON format"`
    } `cmd help:"Batch/AI training mode"`

    // Fast agent mode (protobuf streaming)
    FastAgent struct {
        Format string `flag default:"protobuf" help:"Stream format (protobuf)"`
    } `cmd help:"High-performance protobuf streaming for AI"`

    // Observer mode
    Observer struct {
        GameFile string `flag required name:"game-file" help:"Game to observe"`
        Follow   bool   `flag name:"follow" help:"Follow live updates"`
    } `cmd help:"Read-only observer mode"`

    // Majsoul import/export
    ImportPaipu struct {
        Input      string `arg name:"input" help:"Majsoul paipu JSON file"`
        Output     string `flag name:"output" help:"Output game state file"`
        ExitOnTamper bool `flag name:"exit-on-tamper" help:"Abort on hash validation failure"`
    } `cmd help:"Import Majsoul paipu"`

    ExportPaipu struct {
        Input  string `arg name:"input" help:"Game state file"`
        Output string `arg name:"output" help:"Output Majsoul paipu file"`
    } `cmd help:"Export to Majsoul paipu format"`
}
```

**Example Commands:**

```bash
# Interactive mode
majhong_cli interactive
majhong_cli interactive --load game_001.json

# Batch mode (single action)
majhong_cli batch --game-file state.json --action "discard 3m" --json

# Batch mode (multiple actions)
majhong_cli batch --game-file state.json --actions '["discard 3m", "pon", "discard 5s"]'

# Initialize new game
majhong_cli batch --init --game-file new_game.json --seed 42

# Query game state
majhong_cli batch --game-file state.json --query "valid-actions" --json

# Fast agent mode (protobuf stdin/stdout)
echo [protobuf bytes] | majhong_cli fast-agent > output.pb

# Observer mode
majhong_cli observer --game-file state.json --follow

# Import Majsoul paipu
majhong_cli import-paipu mjs_game.json --output state.json --exit-on-tamper

# Export to Majsoul format
majhong_cli export-paipu state.json mjs_output.json
```

#### 6. TUI Layer (`internal/tui/`)

**Responsibilities:**
- Render game state to terminal (tcell)
- Handle keyboard input
- Update display on state changes
- Support headless mode (no-op renderer)

**Key Components:**

```go
// renderer.go
type Renderer struct {
    screen tcell.Screen
    theme  Theme
}

func (r *Renderer) RenderGame(game *Game, currentPlayer Seat) {
    // Layout:
    // ┌─────────────────────────────────────┐
    // │ Round: East-1  Honba: 0  Riichi: 1  │
    // ├─────────────────────────────────────┤
    // │ Your Hand (East):                   │
    // │ 1m 2m 3m 7m 7m 1p 2p 3p 5s 6s 7s 1z 1z│
    // ├─────────────────────────────────────┤
    // │ Dora: 5p                            │
    // │ Discards:                           │
    // │  East: 9m 8p                        │
    // │  South: 2m 3s                       │
    // │  West: 7p 4m                        │
    // │  North: 1s 6p                       │
    // ├─────────────────────────────────────┤
    // │ Actions: (d)iscard (p)on (c)hi      │
    // │          (k)an (r)iichi (t)sumo     │
    // └─────────────────────────────────────┘
}

// input.go
func (r *Renderer) ReadAction() (Action, error) {
    // Read keyboard input, map to action
    // 'd' → prompt for tile, return DiscardAction
}
```

**Headless Mode:**
- Set `TERM=dumb` or `--headless` flag
- Renderer becomes no-op (returns immediately)
- Batch mode automatically detects and disables TUI

## Data Flow Diagrams

### Interactive Mode Flow

```
User → tcell TUI → Game.ValidActions()
                 ↓
       User selects action
                 ↓
       Game.ApplyAction(action) → New Game State
                 ↓
       Renderer.RenderGame(newState)
                 ↓
       Auto-save (optional)
```

### Batch Mode Flow (Single Action)

```
CLI Args → Parse action string
         ↓
   Load game state from JSON file
         ↓
   Game.ApplyAction(action)
         ↓
   Marshal new state to JSON
         ↓
   Atomic write to file (temp + rename)
         ↓
   Exit with code 0 (or error code)
```

### Fast Agent Mode Flow (Protobuf Streaming)

```
stdin (protobuf) → Unmarshal ActionRequest
                 ↓
   Load game state (embedded or from file)
                 ↓
   Game.ApplyAction(action)
                 ↓
   Marshal ActionResponse to protobuf
                 ↓
   Write to stdout
   (repeat for streaming)
```

### Majsoul Paipu Import Flow

```
Majsoul JSON → Parse paipu format
             ↓
   Extract: wall seed, timeline, players
             ↓
   Map tile notation (0m → aka dora)
             ↓
   Replay actions sequentially
             ↓
   Validate final scores match
             ↓
   Validate wall hash (SHA-256)
             ↓
   Export to internal JSON format
```

## Interface Contracts

### Core Engine API

```go
// pkg/mahjong/engine/game.go
package engine

// NewGame creates a new game with specified configuration
func NewGame(config RuleConfig) (*Game, error)

// ApplyAction applies an action and returns new state
// Returns error if action is invalid
func (g *Game) ApplyAction(action Action) (*Game, error)

// ValidActions returns all legal actions for current player
func (g *Game) ValidActions() []Action

// IsGameOver checks if game has ended (win/draw/exhaustion)
func (g *Game) IsGameOver() (bool, GameResult)

// Clone creates deep copy for simulation
func (g *Game) Clone() *Game
```

### Serialization API

```go
// pkg/mahjong/serialization/serializer.go
package serialization

type Serializer interface {
    Marshal(game *engine.Game) ([]byte, error)
    Unmarshal(data []byte) (*engine.Game, error)
}

// NewJSONSerializer creates JSON serializer
func NewJSONSerializer() Serializer

// NewProtoSerializer creates Protobuf serializer
func NewProtoSerializer() Serializer

// Convert between formats
func Convert(data []byte, from, to Format) ([]byte, error)
```

### Rules API

```go
// pkg/mahjong/rules/yaku.go
package rules

// DetectYaku identifies all yaku in winning hand
func DetectYaku(hand Hand, winTile Tile, context WinContext, config RuleConfig) []Yaku

// CalculateScore computes han/fu → final points
func CalculateScore(yaku []Yaku, fu int, context WinContext) ScoreResult

// CheckFuriten determines if player is in furiten state
func CheckFuriten(player Player, discards []Tile) FuritenState
```

### Majsoul API

```go
// pkg/mahjong/majsoul/importer.go
package majsoul

// ImportPaipu converts Majsoul paipu to internal game state
func ImportPaipu(paipuJSON []byte, validateHash bool) (*engine.Game, error)

// ExportPaipu converts internal state to Majsoul format
func ExportPaipu(game *engine.Game) ([]byte, error)
```

### CLI Exit Codes

```go
const (
    ExitSuccess         = 0  // Action successful
    ExitInvalidAction   = 1  // Invalid action submitted
    ExitIOError         = 2  // File I/O error
    ExitCorruptedState  = 3  // Corrupted game state file
    ExitTamperedData    = 4  // Wall hash validation failed (with --exit-on-tamper)
)
```

## Protobuf Schema

### File: `proto/game_state.proto`

```protobuf
syntax = "proto3";

package mahjong;

option go_package = "github.com/yourusername/majhong_cli/pkg/proto";

// GameState represents complete game state
message GameState {
  string version = 1;
  int64 seed = 2;
  Round round = 3;
  repeated Player players = 4;
  Wall wall = 5;
  repeated Action history = 6;
  RuleConfig config = 7;
}

message Round {
  Seat dealer = 1;
  Wind round_wind = 2;  // East/South/West/North
  int32 honba = 3;
  int32 riichi_sticks = 4;
  repeated Tile dora_indicators = 5;
  repeated Tile uradora_indicators = 6;
}

message Player {
  Seat seat = 1;
  repeated Tile hand = 2;
  repeated Tile discards = 3;
  repeated Meld melds = 4;
  int32 score = 5;
  bool riichi = 6;
  FuritenState furiten = 7;
  bool tenpai = 8;
}

message Wall {
  int32 remaining = 1;
  int64 seed = 2;
  string hash = 3;       // SHA-256(wall + salt)
  string salt_hash = 4;  // SHA-256(salt)
}

message Tile {
  Suit suit = 1;
  int32 rank = 2;  // 1-9 for suited, 1-7 for honors
  bool is_aka_dora = 3;
}

message Action {
  int32 turn = 1;
  Seat player = 2;
  ActionType type = 3;
  repeated Tile tiles = 4;
  int64 timestamp = 5;
}

enum Seat {
  EAST = 0;
  SOUTH = 1;
  WEST = 2;
  NORTH = 3;
}

enum Wind {
  WIND_EAST = 0;
  WIND_SOUTH = 1;
  WIND_WEST = 2;
  WIND_NORTH = 3;
}

enum Suit {
  MAN = 0;    // 萬子 (characters)
  PIN = 1;    // 筒子 (dots)
  SOU = 2;    // 索子 (bamboo)
  JI = 3;     // 字牌 (honors)
}

enum ActionType {
  DRAW = 0;
  DISCARD = 1;
  PON = 2;
  CHI = 3;
  KAN = 4;
  RIICHI = 5;
  TSUMO = 6;
  RON = 7;
}

message FuritenState {
  bool discard_furiten = 1;
  bool temporary_furiten = 2;
  bool riichi_furiten = 3;
}

message Meld {
  MeldType type = 1;
  repeated Tile tiles = 2;
  Seat from_player = 3;
}

enum MeldType {
  PON_MELD = 0;
  CHI_MELD = 1;
  ANKAN = 2;
  MINKAN = 3;
  KAKAN = 4;
}

message RuleConfig {
  bool riichi_mahjong = 1;
  int32 starting_points = 2;  // 25000 or 30000
  int32 return_points = 3;    // 30000 typically
  int32 aka_dora_count = 4;   // 0-4
  repeated string local_yaku = 5;  // e.g., ["renhou", "nagashimangan"]
}
```

### File: `proto/action.proto` (for fast agent streaming)

```protobuf
syntax = "proto3";

package mahjong;

option go_package = "github.com/yourusername/majhong_cli/pkg/proto";

// ActionRequest for fast agent mode
message ActionRequest {
  oneof request {
    InitGameRequest init = 1;
    ApplyActionRequest apply = 2;
    QueryStateRequest query = 3;
  }
}

message InitGameRequest {
  int64 seed = 1;
  RuleConfig config = 2;
}

message ApplyActionRequest {
  GameState current_state = 1;  // Embedded or reference
  Action action = 2;
}

message QueryStateRequest {
  string query_field = 1;  // "valid-actions", "scores", etc.
}

// ActionResponse from engine
message ActionResponse {
  bool success = 1;
  string error_message = 2;
  oneof response {
    GameState new_state = 3;
    QueryResult query_result = 4;
  }
}

message QueryResult {
  repeated Action valid_actions = 1;
  repeated int32 scores = 2;
  string round_info = 3;
  int32 wall_remaining = 4;
}
```

## Testing Strategy

### Phase 0: Contract Tests (TDD Foundation)

**Goal**: Define expected behavior for all 52 yaku and core rules before implementation

**Test Files**:
- `test/contract/yaku_test.go` - All yaku recognition tests
- `test/contract/scoring_test.go` - Han/fu calculation tests
- `test/contract/furiten_test.go` - Furiten detection tests
- `test/contract/wall_test.go` - Wall integrity tests

**Example Test Structure**:

```go
// test/contract/yaku_test.go
func TestYakuDetection(t *testing.T) {
    tests := []struct {
        name     string
        hand     []string
        winTile  string
        context  WinContext
        expected []Yaku
    }{
        {
            name:     "Riichi + Tanyao",
            hand:     []string{"2m", "3m", "4m", "5p", "5p", "5p", "6s", "7s", "8s", "3m", "3m", "3m", "6p"},
            winTile:  "6p",
            context:  WinContext{Riichi: true, Tsumo: false},
            expected: []Yaku{{Name: "Riichi", Han: 1}, {Name: "Tanyao", Han: 1}},
        },
        // ... 52 more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            hand := parseHand(tt.hand)
            winTile := parseTile(tt.winTile)
            result := rules.DetectYaku(hand, winTile, tt.context, defaultConfig)

            if !yakuListEqual(result, tt.expected) {
                t.Errorf("Expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

**Test Data Sources**:
- `test/fixtures/yaku_examples/` - One JSON file per yaku with 3-5 example hands
- Constitution files (`constitution_majhong_yaku.md`) - Reference for correct patterns

### Phase 1: Unit Tests

**Coverage Targets**:
- Game engine: 90%+ (core logic)
- Rules engine: 100% (all yaku, furiten, scoring)
- Serialization: 95% (JSON/Protobuf round-trip)
- CLI layer: 70% (interface code)

**Test Categories**:
1. **State Transitions**: Valid action application
2. **Rule Enforcement**: Invalid action rejection
3. **Edge Cases**: Multiple ron, kan interactions, wall exhaustion
4. **Determinism**: Same seed → same game flow
5. **Performance**: 1000 games/min benchmark

### Phase 2: Integration Tests

**Test Files**:
- `test/integration/batch_mode_test.go` - End-to-end batch mode workflows
- `test/integration/interactive_test.go` - Automated TUI interaction (scripted)
- `test/integration/majsoul_test.go` - Import real Majsoul paipu, verify replay

**Example Integration Test**:

```go
// test/integration/batch_mode_test.go
func TestBatchModeFullGame(t *testing.T) {
    // Create temporary game file
    tempFile := createTempFile(t)
    defer os.Remove(tempFile)

    // Initialize game
    runCLI(t, "batch", "--init", "--game-file", tempFile, "--seed", "42")

    // Play 10 turns
    for i := 0; i < 10; i++ {
        validActions := queryValidActions(t, tempFile)
        action := validActions[0]  // Pick first valid action
        runCLI(t, "batch", "--game-file", tempFile, "--action", action)
    }

    // Verify game state
    state := loadGameState(t, tempFile)
    assert.Equal(t, 10, len(state.History))
    assert.True(t, state.Wall.Remaining < 122)
}
```

### Phase 3: Contract Tests (External API)

**Goal**: Verify public API contracts remain stable

**Test Files**:
- `test/contract/cli_output_test.go` - JSON output schema validation
- `test/contract/protobuf_schema_test.go` - Protobuf backward compatibility

**Example**:

```go
// test/contract/cli_output_test.go
func TestBatchModeJSONOutput(t *testing.T) {
    output := runCLI(t, "batch", "--game-file", "test_game.json", "--query", "scores", "--json")

    var result map[string]interface{}
    err := json.Unmarshal([]byte(output), &result)
    require.NoError(t, err)

    // Validate schema
    assert.Contains(t, result, "scores")
    assert.IsType(t, []interface{}{}, result["scores"])
    assert.Len(t, result["scores"], 4)
}
```

## Implementation Phases

### Phase 0: Foundation (TDD Setup) - Week 1

**Deliverables**:
- [ ] Project structure setup (`go mod init`, directories)
- [ ] Protobuf schema definition (`proto/*.proto`)
- [ ] Contract test suite (52 yaku tests + scoring tests) - **MUST FAIL**
- [ ] CI/CD pipeline (GitHub Actions: lint, test, build)
- [ ] Makefile with targets: `test`, `build`, `proto`, `lint`

**Exit Criteria**:
- All contract tests written and **failing** (Red phase of TDD)
- `make test` runs (fails) and `make build` compiles empty main.go
- Documentation: architecture diagrams, API contracts finalized

**Commands**:
```bash
make proto       # Generate Go code from .proto files
make test        # Run all tests (should fail in Phase 0)
make build       # Compile CLI binary
make lint        # golangci-lint
```

### Phase 1: Core Engine - Week 2-3

**Goal**: Implement game engine to pass contract tests (Green phase)

**Implementation Order**:
1. **Tile & Wall** (`pkg/mahjong/engine/tile.go`, `wall.go`)
   - Tile struct, wall shuffling, draw/discard
   - Wall hash generation
   - Tests: `TestWallIntegrity`, `TestWallShuffle`

2. **Player & Hand** (`pkg/mahjong/engine/player.go`)
   - Hand management, meld tracking
   - Furiten state
   - Tests: `TestHandManagement`, `TestFuritenState`

3. **Game State** (`pkg/mahjong/engine/game.go`, `round.go`)
   - State struct, turn management
   - Action application (immutable pattern)
   - Tests: `TestGameStateTransitions`

4. **Action Validation** (`pkg/mahjong/engine/action.go`)
   - ValidActions() implementation
   - Pon/chi/kan legality
   - Tests: `TestValidActionGeneration`

**Exit Criteria**:
- Core engine tests passing (not yaku/scoring yet)
- Can play a full game programmatically (no TUI)
- Deterministic replay with same seed

### Phase 2: Rules Engine - Week 4-5

**Goal**: Implement all 52 yaku, scoring, furiten (Pass contract tests)

**Implementation Order**:
1. **Yaku Detection** (`pkg/mahjong/rules/yaku.go`)
   - Pattern matching for 37 regular yaku
   - Yakuman detection (13 orphans, etc.)
   - Tests: All 52 yaku contract tests **must pass**

2. **Scoring Calculation** (`pkg/mahjong/rules/scoring.go`)
   - Han/fu table lookup
   - Tsumo/ron distribution
   - Honba, riichi sticks
   - Tests: `TestScoringCalculation`

3. **Furiten Implementation** (`pkg/mahjong/rules/furiten.go`)
   - Discard furiten
   - Temporary furiten (declined ron)
   - Riichi furiten
   - Tests: `TestFuritenDetection`

4. **Pao Responsibility** (`pkg/mahjong/rules/pao.go`)
   - Daisangen, Daisuushii, Suukantsu
   - Payment splitting
   - Tests: `TestPaoPayment`

**Exit Criteria**:
- **ALL contract tests passing** (52 yaku + scoring + furiten)
- Rules engine code ≤ 500 lines (per constitution: minimal complexity)
- Performance: 1000 games/min benchmark

### Phase 3: Serialization & Persistence - Week 6

**Goal**: JSON/Protobuf serialization, file I/O

**Implementation Order**:
1. **JSON Serializer** (`pkg/mahjong/serialization/json.go`)
   - Marshal/Unmarshal
   - Pretty-printing
   - Tests: Round-trip test (Game → JSON → Game)

2. **Protobuf Serializer** (`pkg/mahjong/serialization/protobuf.go`)
   - Generate Go code: `make proto`
   - Marshal/Unmarshal
   - Tests: Round-trip, compatibility

3. **Converter** (`pkg/mahjong/serialization/converter.go`)
   - JSON ↔ Protobuf
   - Tests: Bidirectional conversion

4. **Wall Hash** (`pkg/mahjong/serialization/hash.go`)
   - SHA-256 generation
   - Validation with --exit-on-tamper
   - Tests: Tamper detection

**Exit Criteria**:
- JSON files human-readable
- Protobuf files <50% size of JSON
- Hash validation 100% tamper detection rate

### Phase 4: CLI Layer - Week 7

**Goal**: Kong-based CLI with all modes

**Implementation Order**:
1. **Kong CLI Parser** (`cmd/majhong_cli/main.go`)
   - Define command structure
   - Route to mode handlers
   - Tests: Argument parsing

2. **Batch Mode** (`cmd/majhong_cli/batch.go`)
   - File-based operations
   - Atomic writes
   - Tests: Integration tests

3. **Interactive Mode** (`cmd/majhong_cli/interactive.go`)
   - TUI skeleton (no rendering yet)
   - Action dispatch
   - Tests: Scripted interaction

4. **Fast Agent Mode** (`cmd/majhong_cli/fast_agent.go`)
   - Stdin/stdout protobuf streaming
   - Tests: Benchmark 10k actions/min

5. **Observer Mode** (`cmd/majhong_cli/observer.go`)
   - Read-only game state display
   - Tests: No state mutation

**Exit Criteria**:
- All CLI commands functional
- Exit codes correct
- JSON output validated

### Phase 5: TUI Implementation - Week 8

**Goal**: tcell-based terminal UI (headless-compatible)

**Implementation Order**:
1. **Renderer** (`internal/tui/renderer.go`)
   - Layout design (ASCII art)
   - Game state display
   - Tests: Render snapshots

2. **Input Handler** (`internal/tui/input.go`)
   - Keyboard mapping
   - Tile selection
   - Tests: Input simulation

3. **Headless Mode** (`internal/tui/headless.go`)
   - No-op renderer
   - Batch mode auto-detection
   - Tests: Container environment

**Exit Criteria**:
- Interactive mode playable
- Headless mode passes in Docker container
- Display updates <16ms (60 FPS equivalent)

### Phase 6: Majsoul Integration - Week 9

**Goal**: Import/export Majsoul paipu format

**Implementation Order**:
1. **Paipu Parser** (`pkg/mahjong/majsoul/parser.go`)
   - JSON schema parsing
   - Tile notation mapping
   - Tests: Sample paipu files

2. **Importer** (`pkg/mahjong/majsoul/importer.go`)
   - Timeline replay
   - Hash validation
   - Tests: Import 10 real games, verify scores

3. **Exporter** (`pkg/mahjong/majsoul/exporter.go`)
   - Internal → Majsoul format
   - Tests: Export + re-import = identical

**Exit Criteria**:
- 100% import accuracy (all test paipu files)
- Exported files viewable in Majsoul replay tools

### Phase 7: Performance & Polish - Week 10

**Goal**: Optimize, document, prepare for release

**Tasks**:
1. **Performance Optimization**
   - Profile with `go tool pprof`
   - Optimize hot paths (yaku detection)
   - Target: 1000+ games/min

2. **Documentation**
   - API docs (`godoc`)
   - User guide (README.md)
   - Architecture diagrams (Mermaid)

3. **Docker Image**
   - Multi-stage build
   - Headless mode verification
   - Size optimization

4. **Release Preparation**
   - Version tagging (v1.0.0)
   - Changelog
   - Binary releases (GitHub Actions)

**Exit Criteria**:
- Performance benchmarks met
- README complete
- Docker image <50MB

## Risk Mitigation

| Risk | Impact | Mitigation |
|------|--------|------------|
| **Yaku detection complexity** | High - Core functionality | TDD with 52 test cases before implementation, reference constitution docs |
| **Performance not meeting 1000 games/min** | Medium - AI training bottleneck | Profile early (Phase 2), consider Go's `pprof`, optimize critical paths |
| **Protobuf schema breaking changes** | High - AI module incompatibility | Versioning in proto files, backward compatibility tests |
| **tcell not truly headless** | Medium - Container deployment | Test in Docker early (Phase 5), fallback to pure text mode |
| **Majsoul format changes** | Low - External dependency | Version detection, graceful degradation, document supported versions |
| **File locking issues in K8s** | Medium - Concurrent games | Use advisory locks, test with 100 concurrent containers |
| **SHA-256 hash collision** | Very Low - Security concern | Standard crypto library, 256-bit is sufficient for gaming |

## Success Metrics

**Functional Completeness**:
- [ ] All 52 yaku recognized correctly (100% test pass rate)
- [ ] Score calculation matches reference games
- [ ] Furiten detection: 0 false positives/negatives
- [ ] Majsoul import: 100% accuracy on test set

**Performance**:
- [ ] Interactive mode: <100ms response time
- [ ] Batch mode: 1000+ games/min
- [ ] Fast agent mode: 10,000+ actions/min
- [ ] Container startup: <1 second

**Quality**:
- [ ] Test coverage: >90% overall, 100% rules engine
- [ ] Zero memory leaks (verified with `go test -race`)
- [ ] Lint-clean (`golangci-lint`)
- [ ] Documentation complete (godoc, README)

**Deployment**:
- [ ] Docker image builds and runs headless
- [ ] Kubernetes deployment tested (100 concurrent games)
- [ ] Binary size <20MB (single executable)

## Next Steps

1. **Immediate**: Create Phase 0 deliverables
   - Run `go mod init github.com/yourusername/majhong_cli`
   - Create project structure (directories)
   - Define Protobuf schemas
   - Write contract tests (failing)

2. **Week 1-2**: Implement core engine (Phase 1)
   - TDD: Make contract tests pass one by one
   - Focus on game state integrity

3. **Week 3-5**: Rules engine (Phase 2)
   - Most complex phase (52 yaku)
   - Frequent constitution doc reference

4. **Week 6+**: Follow phases sequentially
   - Serialization → CLI → TUI → Majsoul → Polish

**Command to start**:
```bash
/speckit.tasks  # Generate detailed task breakdown from this plan
```

This plan provides a complete technical architecture for the CLI Mahjong engine, with clear module boundaries, data flow, testing strategy, and phased implementation. The hybrid JSON/Protobuf approach balances human-readability with performance, and the library-first design enables future AI and web modules to build on this foundation.
