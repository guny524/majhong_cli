# Feature Specification: CLI-Based Riichi Mahjong Game Engine

**Feature Branch**: `001-cli-mahjong-engine`
**Created**: 2025-11-18
**Status**: Draft
**Input**: User description: "CLI-based Riichi Mahjong game engine with dual execution modes (interactive and batch), external game state file format for AI training compatibility, and Majsoul paipu format integration"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Interactive Local Game Play (Priority: P1)

A player wants to play a complete Riichi Mahjong game locally in their terminal, with real-time input for decisions and visual representation of the game state.

**Why this priority**: This is the foundational MVP that demonstrates the core game engine works correctly. All other features (AI training, web integration, probability calculation) depend on having a working game engine first.

**Independent Test**: Can be fully tested by launching the CLI, playing a complete round (East 1 through completion), making all standard actions (discard, pon, chi, kan, riichi, tsumo, ron), and verifying the game reaches a proper end state with correct scoring.

**Acceptance Scenarios**:

1. **Given** the CLI is launched in interactive mode, **When** the user starts a new game, **Then** 4 players receive 13/14 tiles, dora indicator is revealed, and the game state is displayed
2. **Given** it's a player's turn, **When** the player draws a tile and discards, **Then** the discard is added to the discard pond and turn passes to next player
3. **Given** a player has a valid winning hand, **When** they declare tsumo/ron, **Then** the game calculates score correctly according to yaku and han, updates player points, and displays final results
4. **Given** a round completes, **When** checking game state, **Then** the wall integrity is maintained (122 live wall + 14 dead wall = 136 tiles total)
5. **Given** a game session is active, **When** the user exits, **Then** the game state is persisted to a file for resume

---

### User Story 2 - Batch Mode for AI Integration (Priority: P2)

An AI training system needs to execute game actions programmatically by passing commands via CLI arguments and reading/writing game state from external files.

**Why this priority**: This enables AI training and experimentation, which is the second core objective of the project. It must be independent from interactive mode but share the same game engine logic.

**Independent Test**: Can be fully tested by running CLI commands like `majhong_cli --game-file game.json --action "discard 3m"`, verifying the game state file is updated correctly, and chaining multiple commands to complete a full game without user interaction.

**Acceptance Scenarios**:

1. **Given** a game state file exists, **When** CLI is executed with `--game-file path.json --action "draw"`, **Then** the action is applied, game state file is updated, and CLI exits with success code
2. **Given** a new game is needed, **When** CLI is executed with `--init --game-file path.json`, **Then** a new game state file is created with initial wall shuffling and tile distribution
3. **Given** multiple AI agents are training, **When** each agent calls CLI with their action in sequence, **Then** each receives updated game state reflecting all previous actions in proper turn order
4. **Given** an invalid action is submitted (e.g., discard a tile not in hand), **When** CLI executes, **Then** it returns an error code and error message without modifying game state
5. **Given** a game completes, **When** CLI is called with `--game-file path.json --query score`, **Then** it outputs final scores in machine-readable format (JSON)

---

### User Story 3 - Majsoul Paipu Compatibility (Priority: P3)

A user wants to load Majsoul game replays (paipu) into the CLI engine to analyze completed games or generate training data from real player matches.

**Why this priority**: This enables learning from real game data and validates the engine's correctness against a production Mahjong platform. It's lower priority because the engine must work standalone first.

**Independent Test**: Can be fully tested by downloading a Majsoul paipu JSON file, importing it via `majhong_cli --import-paipu mjs_game.json --export game_state.json`, and verifying the resulting game state file contains all actions, tiles, and final scores matching the original replay.

**Acceptance Scenarios**:

1. **Given** a Majsoul paipu JSON file, **When** CLI imports it, **Then** the wall seed, tile distribution, and all player actions are converted to the internal game state format
2. **Given** an imported game state, **When** replaying it step-by-step, **Then** the engine reproduces the exact same game flow and final scores as the original Majsoul game
3. **Given** a completed internal game, **When** exporting to Majsoul paipu format, **Then** the exported file is valid for viewing in Majsoul replay tools
4. **Given** a paipu with red dora (0m/0p/0s), **When** importing, **Then** red dora tiles are correctly mapped to the internal tile representation
5. **Given** a paipu includes wall hash validation (SHA-256), **When** importing, **Then** CLI validates hash integrity and rejects tampered replays

---

### User Story 4 - Game State Persistence and Resume (Priority: P2)

A user wants to save a game in progress and resume it later, or share game states with others for review and analysis.

**Why this priority**: This is critical for both AI training (saving/loading training states) and user experience. It's closely tied to P2 batch mode since they share the same file format.

**Independent Test**: Can be fully tested by starting a game, playing several turns, saving with `--save`, closing the CLI, reopening with `--load`, and verifying all game state (hands, wall, discards, dora, scores) is restored exactly.

**Acceptance Scenarios**:

1. **Given** a game in progress, **When** user saves to `game_001.json`, **Then** all game state including wall remaining tiles, player hands, discard ponds, dora indicators, riichi status, and current scores are serialized
2. **Given** a saved game file, **When** user loads it, **Then** the game resumes at the exact same state with correct player turn
3. **Given** a saved game file, **When** loading on a different machine, **Then** the game state is portable and loads identically
4. **Given** game state files from multiple points in a game, **When** analyzing them, **Then** it's possible to trace the complete history of tile draws, discards, and decisions
5. **Given** a corrupted game state file, **When** attempting to load, **Then** CLI validates the file format and reports specific errors without crashing

---

### User Story 5 - Read-Only Observer Mode (Priority: P3)

An instructor or student wants to watch a live game without participating, for demonstration or learning purposes.

**Why this priority**: Useful for teaching and demos, but not critical for core gameplay or AI training. Can be added after core features are stable.

**Independent Test**: Can be fully tested by starting a 4-player game, launching a 5th CLI instance with `--observe`, and verifying the observer sees real-time game state updates without ability to input actions.

**Acceptance Scenarios**:

1. **Given** a game in progress, **When** observer joins via `--observe` mode, **Then** they see current game state (all visible information: discards, dora, scores) but not concealed hands
2. **Given** an observer is watching, **When** players take actions, **Then** observer's display updates in real-time with each state change
3. **Given** an observer tries to input an action, **When** they attempt any command, **Then** the CLI rejects it with message "Observer mode is read-only"
4. **Given** an observer disconnects, **When** they leave, **Then** the game continues unaffected for the 4 players
5. **Given** multiple observers, **When** watching the same game, **Then** all observers see synchronized game state

---

### Edge Cases

- **What happens when the live wall is exhausted (ryuukyoku)?** The round ends in a draw, tenpai players receive points from noten players, and the game advances to the next round or ends based on East player's tenpai status
- **How does the system handle kan calls affecting the dead wall?** Each kan declaration flips a new dora indicator, the player draws a rinshan tile from the dead wall's first 4 tiles, and the last tile from live wall is moved to maintain 14 dead wall tiles
- **What happens when a player triggers furiten?** The player cannot declare ron on any discard until they change their hand (by drawing and discarding), even if the winning tile appears. This must be enforced by the engine
- **How is wall shuffling deterministic for reproducibility?** The wall uses a seeded random number generator, with the seed stored in game state files to ensure exact replay capability for debugging and AI training
- **What happens when multiple players can call the same discard?** Priority order is enforced: ron > pon/kan > chi. If multiple players declare ron simultaneously, multiple winners are handled (double/triple ron)
- **How does the system handle chombo (penalty violations)?** Serious rule violations result in penalty points, potential win prohibition, and the round may restart depending on the violation type
- **What happens when loading an old game state file after a format change?** The CLI includes versioning in the file format and provides migration utilities or clear error messages for incompatible versions
- **How does batch mode handle concurrent access to the same game file?** The CLI uses file locking or atomic writes to prevent corruption when multiple AI agents attempt simultaneous access
- **What happens when importing a Majsoul paipu with rules not yet implemented?** The CLI logs unsupported features (e.g., specific local rules) and either imports with warnings or gracefully fails with clear indication of what's missing
- **How does wall hash validation affect different game modes?** By default, validation failures show warnings and continue play for analysis. With `--exit-on-tamper`, validation failures abort immediately (exit code 4). Local AI training skips validation entirely. Future network play will disconnect from match on tamper detection
- **What happens when batch action execution fails midway?** When using `--actions` array for multiple actions, the CLI processes sequentially and stops at the first failed action, returning which specific action failed with its error message and leaving game state at the last successful action

## Requirements *(mandatory)*

### Functional Requirements

#### Core Game Engine

- **FR-001**: System MUST implement complete Riichi Mahjong rules including tile types (manzu/pinzu/souzu/jihai), yaku recognition, han/fu scoring, dora/uradora, and yakuman
- **FR-002**: System MUST maintain wall integrity with 136 tiles: 122 live wall + 14 dead wall (5 dora indicators, 5 uradora, 4 rinshan tiles)
- **FR-003**: System MUST support all standard player actions: discard, pon, chi, kan (ankan/minkan/kakan), riichi declaration, tsumo, and ron
- **FR-004**: System MUST enforce turn order (East → South → West → North) and validate actions based on current game state
- **FR-005**: System MUST implement furiten rules (discard furiten, temporary furiten, riichi furiten) to prevent invalid ron declarations
- **FR-006**: System MUST calculate scores correctly using the han/fu table including: base points, dealer bonus (1.5x), tsumo/ron distribution, honba counters, and riichi stick payments
- **FR-007**: System MUST handle special hands (kokushi musou, chuuren poutou, daisangen, etc.) with correct yakuman scoring
- **FR-008**: System MUST track game progression through rounds (East 1-4, South 1-4, optional West/North for extended play)
- **FR-009**: System MUST implement wall hash validation using SHA-256 for tamper detection following Majsoul's H(wall + salt) + H(salt) pattern
- **FR-010**: System MUST handle pao (responsibility payment) for daisangen, daisuushii, and suukantsu
- **FR-011**: System MUST calculate final scores with uma (placement bonuses) and oka (return point differential bonus)
- **FR-012**: System MUST detect and handle chombo (penalty violations) including: illegal riichi, incorrect ron declaration, and exposed hand errors

#### Interactive Mode

- **FR-013**: CLI MUST provide a text-based UI showing: current player's hand, dora indicators, discard ponds for all players, player scores, and current round/honba
- **FR-014**: CLI MUST accept user input for tile selection using standard notation (1-9m/p/s for numbered tiles, 1-7z for honors)
- **FR-015**: CLI MUST display available actions (discard, pon, chi, kan, riichi, tsumo, ron) based on current game state and prompt for user choice
- **FR-016**: CLI MUST provide a help command showing game rules, yaku list, and input syntax
- **FR-017**: CLI MUST validate user input and provide clear error messages for invalid actions (e.g., "Cannot pon: tile not in your hand")
- **FR-018**: CLI MUST support configurable rule sets via command-line flags: riichi mahjong vs. standard mahjong, local yaku inclusion, starting points (25000/30000), aka dora count

#### Observer Mode

- **FR-018a**: CLI MUST support read-only observer mode via `--observe` flag for watching games without participating
- **FR-018b**: Observer mode MUST display all publicly visible information (discard ponds, dora indicators, called melds, scores, round info) but NOT concealed hands
- **FR-018c**: Observer mode MUST update display in real-time as game state changes
- **FR-018d**: Observer mode MUST reject any action input with clear error message indicating read-only status
- **FR-018e**: Observer disconnection MUST NOT affect the ongoing game for active players

#### Batch Mode

- **FR-019**: CLI MUST accept game state file path via `--game-file <path>` argument
- **FR-020**: CLI MUST accept action commands via `--action "<action>"` argument for single action execution with syntax: "discard 3m", "pon", "chi 456m", "kan 1111z", "riichi", "tsumo", "ron"
- **FR-020a**: CLI MUST optionally support batch action execution via `--actions '<json_array>'` argument for AI training efficiency (e.g., `--actions '["discard 3m", "pon", "discard 5s"]'`), processing actions sequentially and reporting which action failed if any error occurs
- **FR-021**: CLI MUST initialize new game state via `--init` with optional `--seed <number>` for reproducible shuffling
- **FR-022**: CLI MUST atomically update game state files to prevent corruption from concurrent access
- **FR-023**: CLI MUST return exit codes: 0 for success, 1 for invalid action, 2 for file I/O error, 3 for corrupted state, 4 for tampered wall data (when `--exit-on-tamper` is set)
- **FR-024**: CLI MUST output action results to stdout in JSON format when `--json` flag is provided
- **FR-025**: CLI MUST query game state via `--query <field>` where field can be: "current-player", "valid-actions", "scores", "round", "wall-remaining"

#### Game State File Format

- **FR-026**: Game state files MUST be JSON format with schema version field for future compatibility
- **FR-027**: Game state MUST include: wall seed, tile positions (encoded to prevent cheating), player hands, discard ponds, dora indicators, call history, scores, round info, turn number
- **FR-028**: Game state MUST encode wall tiles as hashed or encrypted data that can only be revealed as tiles are drawn, preventing look-ahead cheating
- **FR-029**: Game state MUST include action history with timestamps for complete game replay
- **FR-030**: Game state MUST be deterministic: loading the same file and replaying actions produces identical results
- **FR-031**: Game state files MUST be human-readable (formatted JSON) for debugging but secure against wall manipulation

#### Majsoul Paipu Integration

- **FR-032**: CLI MUST import Majsoul paipu JSON format via `--import-paipu <path>` command
- **FR-033**: CLI MUST export internal game state to Majsoul paipu format via `--export-paipu <path>` command
- **FR-034**: CLI MUST map Majsoul tile notation (0m/0p/0s for aka dora, 5mr/5pr/5sr in some formats) to internal representation
- **FR-035**: CLI MUST validate imported paipu and loaded game state wall hashes (SHA-256) against H(wall + salt) and H(salt) for integrity verification
- **FR-035a**: When wall hash validation fails, CLI MUST by default display warning message and continue execution to allow forensic analysis
- **FR-035b**: CLI MUST support `--exit-on-tamper` flag to abort immediately on wall hash validation failure, exiting with error code 4 (tampered data)
- **FR-035c**: For local AI training mode, wall hash validation MUST be skipped or disabled as training uses locally-generated clean game states
- **FR-035d**: For network/multiplayer scenarios (future), `--exit-on-tamper` behavior MUST disconnect from the match rather than just exit the process
- **FR-036**: CLI MUST extract all player actions (draw, discard, call, declare) from paipu timeline format
- **FR-037**: CLI MUST handle paipu metadata including: player names, room type (Copper/Silver/Gold/Throne), timestamp, final scores

#### Data Persistence

- **FR-038**: CLI MUST save game state via `--save <path>` command during interactive mode
- **FR-039**: CLI MUST load saved game state via `--load <path>` command to resume play
- **FR-040**: CLI MUST support auto-save to `~/.majhong_cli/autosave.json` after every turn (configurable)
- **FR-041**: CLI MUST include versioning in saved files and provide migration path or clear errors for format incompatibility

### Key Entities

- **Game**: Represents a complete match with current round, honba counter, riichi stick count, player scores, rule configuration, and game history
- **Round**: Represents one round (e.g., East 1) with dealer position, wall state, dora/uradora indicators, and round wind
- **Player**: Has: seat position (East/South/West/North), point total, hand tiles (concealed + exposed), discard pond, riichi status, furiten state, tenpai status
- **Wall**: Contains 136 tiles with structure: live wall (122), dead wall (14 = 5 dora + 5 uradora + 4 rinshan), wall seed for reproducibility, remaining tile count, hash validation data
- **Tile**: Identified by suit (man/pin/sou/ji), rank (1-9 for suited, 1-7 for honors), aka dora flag (for red 5s), dora indicator status
- **Action**: Represents a game action with type (draw/discard/pon/chi/kan/riichi/tsumo/ron), actor (player seat), tiles involved, timestamp, validity flag
- **Meld (Call)**: Represents a pon/chi/kan with tile composition, call source (which player discarded), open/concealed status
- **Hand Analysis**: Contains: tenpai status, waits (machi tiles), potential yaku, han/fu count, shanten count (tiles away from tenpai), furiten state
- **Score Calculation**: Contains: winning hand, yaku list, han count, fu calculation breakdown, final point value, pao responsibility (if applicable), payment distribution

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A player can complete a full 4-player East-round game (East 1-4) in interactive mode with correct score calculation, verified against known game records
- **SC-002**: Batch mode supports 1000+ game simulations per minute on standard hardware for AI training throughput
- **SC-003**: Game state files are 100% deterministic - loading and replaying produces identical final scores across all test cases
- **SC-004**: Majsoul paipu import accuracy is 100% for all standard rules - all imported games produce identical final scores when replayed
- **SC-005**: Wall hash validation detects 100% of tampered wall data in test cases with modified game state files
- **SC-006**: Interactive mode responds to user input within 100ms for real-time playability
- **SC-007**: All 52 standard yaku (37 regular + 15 yakuman) are correctly recognized and scored according to official Riichi Mahjong rules
- **SC-008**: Furiten detection has 100% accuracy - no false positives or false negatives in comprehensive test suite
- **SC-009**: Concurrent batch mode file access handles 10+ simultaneous AI agents without corruption or deadlock
- **SC-010**: Game state file format is forward/backward compatible - version N+1 CLI can load version N files with migrations

## Assumptions

- **Wall Shuffling**: Fisher-Yates shuffle algorithm is sufficient for randomness; cryptographically secure random is not required for local play (only for competitive online play with hash validation)
- **Performance Target**: Single-threaded execution is acceptable; the engine does not need to support parallelism internally (AI training parallelism is handled at the process level)
- **Language Choice**: Implementation in Go as specified in constitution, with future Rust migration planned - specification remains language-agnostic
- **Tile Notation**: Using modern Mahjong notation (1-9m/p/s, 1-7z) as standard; conversion utilities for other notations (e.g., MPSZ format) can be added later
- **Network Protocol**: This specification covers local CLI only; network multiplayer is a future enhancement and out of scope
- **UI Complexity**: Text-based terminal UI is sufficient; graphical tile rendering is out of scope for the CLI engine
- **Rule Variants**: Default implementation is Japanese Riichi Mahjong per WRC (World Riichi Championship) rules; other variants (Chinese, Hong Kong) are configurable extensions
- **AI Agent Interface**: Batch mode provides sufficient interface for AI training; dedicated API endpoints (gRPC/REST) are future enhancements
- **File Format**: JSON is acceptable despite verbosity; binary formats (protobuf/msgpack) can be added for performance optimization later
- **Replay Storage**: Local filesystem is sufficient; cloud storage integration is out of scope
- **Red Dora Count**: Standard configuration includes 3 red fives (one each in 5m/5p/5s); this is configurable per game
- **Observer Mode Communication**: Read-only observer mode requires some form of IPC (inter-process communication) or shared game state file polling; exact mechanism is implementation detail
- **Tamper Detection Default**: Default behavior (warning without exit) prioritizes forensic analysis capability over strict security, appropriate for local/training scenarios; production network play would default to `--exit-on-tamper`

## Dependencies

- **Constitution Files**: Game logic must strictly follow rules defined in:
  - `constitution_majhong_rule.md` - Core game rules and wall structure
  - `constitution_majhong_yaku.md` - Complete yaku recognition patterns
  - `constitution_majhong_score.md` - Scoring tables including han/fu, pao, oka
  - `constitution_majhong_penalties.md` - Chombo detection and penalty enforcement
  - `constitution_majhong_term.md` - Standardized EN/KR/JP/CN terminology

- **External Libraries**: No external dependencies for core game logic (self-contained); JSON parsing can use standard library

- **Test Data**: Requires collection of known-correct game records for validation:
  - Majsoul paipu samples covering all standard yaku
  - Edge case games (furiten scenarios, multiple ron, pao triggers)
  - Benchmark games with verified score calculations

- **Development Tools**: Requires access to:
  - Go compiler (version 1.21+) as per constitution
  - Standard testing framework
  - JSON schema validation tools for game state format

## Out of Scope

- **Graphical User Interface**: No tile images, animations, or GUI - CLI text mode only
- **Online Multiplayer**: No client-server networking, matchmaking, or lobby systems
- **User Accounts**: No authentication, profiles, statistics tracking, or cloud save
- **Real-time Communication**: No chat, emoticons, or player interaction beyond game actions
- **Tournament Mode**: No bracket management, Swiss pairings, or competitive ranking
- **Replay Analysis Tools**: No built-in move analysis, probability visualization, or AI commentary (these are separate modules)
- **Web Integration**: No browser capture, screen scraping, or overlay rendering (covered by separate majhong_mjs_bridge repository)
- **AI Agents**: No probability calculator or reinforcement learning agents (covered by separate majhong_ai repository)
- **Mobile Support**: No iOS/Android builds - desktop CLI only
- **Internationalization**: Default English output; localization to other languages is a future enhancement
- **Sound Effects**: No audio - silent operation only
- **Configurable UI Themes**: Single monochrome text display; no color schemes or styling options
- **Undo/Redo**: No takeback of actions - all moves are final once submitted
- **Hints/Suggestions**: No move recommendations or tutorial mode (AI suggestion is separate module)
