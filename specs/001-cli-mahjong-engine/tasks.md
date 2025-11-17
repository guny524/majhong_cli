---
description: "Task list for CLI-Based Riichi Mahjong Game Engine implementation"
---

# Tasks: CLI-Based Riichi Mahjong Game Engine

**Input**: Design documents from `/specs/001-cli-mahjong-engine/`
**Prerequisites**: plan.md (required), spec.md (required for user stories)

**Tests**: Tests are NOT explicitly requested in the specification but are required per constitution (Principle III: Test-Driven Development). All contract tests MUST be written FIRST and FAIL before implementation.

**Organization**: Tasks are grouped by implementation phase and user story to enable independent implementation and testing.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3, US4, US5)
- Include exact file paths in descriptions

---

## Phase 0: Foundation (TDD Setup)

**Purpose**: Project initialization, protobuf schemas, and contract tests that MUST FAIL

**⚠️ CRITICAL**: All contract tests MUST be written and FAILING before any implementation begins (Red phase of TDD)

- [ ] T001 Initialize Go module: `go mod init github.com/guny524/majhong_cli` in repository root
- [ ] T002 [P] Create directory structure: `cmd/majhong_cli/`, `pkg/mahjong/engine/`, `pkg/mahjong/rules/`, `pkg/mahjong/serialization/`, `pkg/mahjong/majsoul/`, `internal/tui/`, `internal/config/`, `proto/`, `test/fixtures/`, `test/integration/`, `test/contract/`
- [ ] T003 [P] Define protobuf schema for game state in proto/game_state.proto
- [ ] T004 [P] Define protobuf schema for actions (fast agent mode) in proto/action.proto
- [ ] T005 [P] Define protobuf schema for Majsoul paipu format in proto/majsoul.proto
- [ ] T006 Create Makefile with targets: proto, test, build, lint, docker
- [ ] T007 [P] Setup GitHub Actions CI/CD in .github/workflows/ci.yml (lint, test, build)
- [ ] T008 [P] Configure golangci-lint in .golangci.yml

### Contract Tests (MUST FAIL - TDD Red Phase)

**Purpose**: Define expected behavior BEFORE implementation. These tests verify correctness against constitution docs.

- [ ] T009 [P] Write contract tests for all 52 yaku in test/contract/yaku_test.go (reference constitution_majhong_yaku.md) - MUST FAIL
- [ ] T010 [P] Write contract tests for scoring calculation in test/contract/scoring_test.go (reference constitution_majhong_score.md) - MUST FAIL
- [ ] T011 [P] Write contract tests for furiten detection in test/contract/furiten_test.go (reference constitution_majhong_rule.md) - MUST FAIL
- [ ] T012 [P] Write contract tests for wall integrity in test/contract/wall_test.go - MUST FAIL
- [ ] T013 [P] Create test fixtures in test/fixtures/yaku_examples/ (one JSON file per yaku with 3-5 example hands)
- [ ] T014 [P] Create test fixtures in test/fixtures/game_states/ (sample game state JSON files)
- [ ] T015 [P] Create test fixtures in test/fixtures/majsoul_paipu/ (sample Majsoul replay files)

**Checkpoint**: All contract tests written and FAILING. `make test` fails as expected (Red phase complete).

---

## Phase 1: Foundational (Blocking Prerequisites)

**Purpose**: Core game engine infrastructure that MUST be complete before ANY user story implementation

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Core Engine Foundation

- [ ] T016 [P] Implement Tile struct and constructors in pkg/mahjong/engine/tile.go
- [ ] T017 [P] Implement tile notation parser (1-9m/p/s, 1-7z) in pkg/mahjong/engine/tile.go
- [ ] T018 Implement Wall struct with Fisher-Yates shuffle in pkg/mahjong/engine/wall.go
- [ ] T019 Implement Wall.Draw() and wall integrity maintenance (122 live + 14 dead) in pkg/mahjong/engine/wall.go
- [ ] T020 Implement wall hash generation (SHA-256) in pkg/mahjong/serialization/hash.go
- [ ] T021 Implement wall hash validation with --exit-on-tamper support in pkg/mahjong/serialization/hash.go
- [ ] T022 [P] Implement Player struct with hand management in pkg/mahjong/engine/player.go
- [ ] T023 [P] Implement Meld struct (pon/chi/kan types) in pkg/mahjong/engine/player.go
- [ ] T024 Implement Round struct with dealer, wind, honba, riichi sticks in pkg/mahjong/engine/round.go
- [ ] T025 Implement Game struct with state management in pkg/mahjong/engine/game.go
- [ ] T026 Implement NewGame() with RuleConfig and seed in pkg/mahjong/engine/game.go
- [ ] T027 Implement Action struct and action types (draw, discard, pon, chi, kan, riichi, tsumo, ron) in pkg/mahjong/engine/action.go
- [ ] T028 Implement Game.ApplyAction() with immutable state pattern in pkg/mahjong/engine/game.go
- [ ] T029 Implement Game.ValidActions() for current player in pkg/mahjong/engine/game.go
- [ ] T030 Implement Game.IsGameOver() with win/draw/exhaustion detection in pkg/mahjong/engine/game.go
- [ ] T031 Implement action validation (turn order, tile ownership) in pkg/mahjong/rules/validation.go
- [ ] T032 Implement pon/chi/kan legality checking in pkg/mahjong/rules/validation.go
- [ ] T033 [P] Add unit tests for tile operations in pkg/mahjong/engine/tile_test.go
- [ ] T034 [P] Add unit tests for wall integrity in pkg/mahjong/engine/wall_test.go
- [ ] T035 [P] Add unit tests for game state transitions in pkg/mahjong/engine/game_test.go
- [ ] T036 Verify wall contract tests pass (test/contract/wall_test.go)

**Checkpoint**: Foundation ready - core engine functional, wall contract tests passing

---

## Phase 2: Rules Engine (Core Game Logic)

**Purpose**: Implement all Riichi Mahjong rules to pass ALL contract tests

**Goal**: Make ALL contract tests GREEN (TDD Green phase)

### Yaku Recognition (52 total)

- [ ] T037 Implement YakuDetector struct and base detection framework in pkg/mahjong/rules/yaku.go
- [ ] T038 [P] Implement 1-han yaku detection (Riichi, Tanyao, Pinfu, Iipeikou, etc.) in pkg/mahjong/rules/yaku.go
- [ ] T039 [P] Implement 2-han yaku detection (Chiitoitsu, Chanta, Sanshoku Doujun, etc.) in pkg/mahjong/rules/yaku.go
- [ ] T040 [P] Implement 3-han+ yaku detection (Honitsu, Junchan, Ryanpeikou) in pkg/mahjong/rules/yaku.go
- [ ] T041 [P] Implement yakuman detection (Kokushi, Suuankou, Daisangen, Shousuushii, Daisuushii, Tsuuiisou, Chinroutou, Ryuuiisou, Chuuren Poutou, Suukantsu, Tenhou, Chiihou, Renhou) in pkg/mahjong/rules/yaku.go
- [ ] T042 [P] Implement aka dora support in yaku detection in pkg/mahjong/rules/yaku.go
- [ ] T043 Verify ALL 52 yaku contract tests pass (test/contract/yaku_test.go)

### Scoring Calculation

- [ ] T044 Implement han/fu table lookup in pkg/mahjong/rules/scoring.go
- [ ] T045 Implement fu calculation (base 20, melds, wait patterns, pair) in pkg/mahjong/rules/scoring.go
- [ ] T046 Implement base points calculation (fu × 2^(han+2)) in pkg/mahjong/rules/scoring.go
- [ ] T047 Implement tsumo payment distribution in pkg/mahjong/rules/scoring.go
- [ ] T048 Implement ron payment calculation in pkg/mahjong/rules/scoring.go
- [ ] T049 Implement dealer bonus (1.5x multiplier) in pkg/mahjong/rules/scoring.go
- [ ] T050 Implement honba and riichi stick payments in pkg/mahjong/rules/scoring.go
- [ ] T051 Implement pao (responsibility payment) for daisangen/daisuushii/suukantsu in pkg/mahjong/rules/pao.go
- [ ] T052 Implement uma (placement bonuses) calculation in pkg/mahjong/rules/scoring.go
- [ ] T053 Implement oka (return point differential) in pkg/mahjong/rules/scoring.go
- [ ] T054 Verify ALL scoring contract tests pass (test/contract/scoring_test.go)

### Furiten Detection

- [ ] T055 Implement discard furiten detection in pkg/mahjong/rules/furiten.go
- [ ] T056 Implement temporary furiten (declined ron) in pkg/mahjong/rules/furiten.go
- [ ] T057 Implement riichi furiten in pkg/mahjong/rules/furiten.go
- [ ] T058 Implement furiten state management in Game struct in pkg/mahjong/engine/game.go
- [ ] T059 Verify ALL furiten contract tests pass (test/contract/furiten_test.go)

### Special Rules & Edge Cases

- [ ] T060 [P] Implement chombo (penalty violations) detection in pkg/mahjong/rules/validation.go (reference constitution_majhong_penalties.md)
- [ ] T061 [P] Implement multiple ron handling (double/triple ron) in pkg/mahjong/engine/game.go
- [ ] T062 [P] Implement ryuukyoku (exhaustive draw) handling in pkg/mahjong/engine/game.go
- [ ] T063 [P] Implement tenpai/noten point exchange on draw in pkg/mahjong/rules/scoring.go
- [ ] T064 [P] Implement kan interactions (new dora, rinshan draw) in pkg/mahjong/engine/game.go
- [ ] T065 Add performance benchmark test: 1000 games/min in test/integration/benchmark_test.go

**Checkpoint**: ALL contract tests passing (100% yaku + scoring + furiten). Rules engine complete.

---

## Phase 3: User Story 1 - Interactive Local Game Play (Priority: P1) 🎯 MVP

**Goal**: Enable a complete playable Riichi Mahjong game in the terminal

**Independent Test**: Launch CLI, play East 1 round to completion with all actions working

### Serialization Layer (Required for persistence)

- [ ] T066 [P] [US1] Implement JSON serializer Marshal in pkg/mahjong/serialization/json.go
- [ ] T067 [P] [US1] Implement JSON serializer Unmarshal in pkg/mahjong/serialization/json.go
- [ ] T068 [P] [US1] Add JSON round-trip test in pkg/mahjong/serialization/json_test.go
- [ ] T069 [P] [US1] Implement Protobuf serializer Marshal in pkg/mahjong/serialization/protobuf.go
- [ ] T070 [P] [US1] Implement Protobuf serializer Unmarshal in pkg/mahjong/serialization/protobuf.go
- [ ] T071 [P] [US1] Implement JSON ↔ Protobuf converter in pkg/mahjong/serialization/converter.go
- [ ] T072 [P] [US1] Add protobuf round-trip test in pkg/mahjong/serialization/protobuf_test.go

### TUI Implementation

- [ ] T073 [P] [US1] Setup tcell screen initialization in internal/tui/renderer.go
- [ ] T074 [P] [US1] Implement game state layout renderer in internal/tui/renderer.go
- [ ] T075 [P] [US1] Implement hand display (concealed tiles) in internal/tui/renderer.go
- [ ] T076 [P] [US1] Implement discard pond display for all 4 players in internal/tui/renderer.go
- [ ] T077 [P] [US1] Implement dora indicator display in internal/tui/renderer.go
- [ ] T078 [P] [US1] Implement score display for all players in internal/tui/renderer.go
- [ ] T079 [P] [US1] Implement action menu display (discard/pon/chi/kan/riichi/tsumo/ron) in internal/tui/layout.go
- [ ] T080 [US1] Implement keyboard input handler in internal/tui/input.go
- [ ] T081 [US1] Implement tile selection input in internal/tui/input.go
- [ ] T082 [US1] Implement action selection input in internal/tui/input.go
- [ ] T083 [P] [US1] Implement headless mode (no-op renderer) in internal/tui/headless.go
- [ ] T084 [P] [US1] Add TUI unit tests in internal/tui/renderer_test.go

### CLI Interactive Mode

- [ ] T085 [US1] Setup kong CLI parser structure in cmd/majhong_cli/main.go
- [ ] T086 [US1] Implement Interactive command handler in cmd/majhong_cli/interactive.go
- [ ] T087 [US1] Implement game loop (draw → action → render) in cmd/majhong_cli/interactive.go
- [ ] T088 [US1] Implement save game functionality (--save flag) in cmd/majhong_cli/interactive.go
- [ ] T089 [US1] Implement load game functionality (--load flag) in cmd/majhong_cli/interactive.go
- [ ] T090 [US1] Implement auto-save after each turn in cmd/majhong_cli/interactive.go
- [ ] T091 [US1] Implement help command showing rules and input syntax in cmd/majhong_cli/interactive.go
- [ ] T092 [P] [US1] Implement rule configuration loading from file in internal/config/rules.go
- [ ] T093 [US1] Add integration test for interactive mode in test/integration/interactive_test.go

### User Story 1 Validation

- [ ] T094 [US1] Verify acceptance scenario 1: Game initialization with 4 players
- [ ] T095 [US1] Verify acceptance scenario 2: Turn-based discard and turn passing
- [ ] T096 [US1] Verify acceptance scenario 3: Winning hand with correct scoring
- [ ] T097 [US1] Verify acceptance scenario 4: Wall integrity maintained throughout game
- [ ] T098 [US1] Verify acceptance scenario 5: Save/resume game state persistence

**Checkpoint**: User Story 1 complete - Full interactive game playable from terminal

---

## Phase 4: User Story 2 - Batch Mode for AI Integration (Priority: P2)

**Goal**: Enable programmatic game execution for AI training

**Independent Test**: Execute full game via CLI commands without user interaction, verify state file updates

### CLI Batch Mode

- [ ] T099 [US2] Implement Batch command structure in cmd/majhong_cli/main.go
- [ ] T100 [US2] Implement batch mode handler in cmd/majhong_cli/batch.go
- [ ] T101 [US2] Implement --init flag for new game creation in cmd/majhong_cli/batch.go
- [ ] T102 [US2] Implement --seed flag for deterministic wall shuffling in cmd/majhong_cli/batch.go
- [ ] T103 [US2] Implement --action flag for single action execution in cmd/majhong_cli/batch.go
- [ ] T104 [US2] Implement --actions flag for batch action array in cmd/majhong_cli/batch.go
- [ ] T105 [US2] Implement --query flag (current-player, valid-actions, scores, round, wall-remaining) in cmd/majhong_cli/batch.go
- [ ] T106 [US2] Implement --json flag for JSON output format in cmd/majhong_cli/batch.go
- [ ] T107 [US2] Implement atomic file writes (temp + rename) in cmd/majhong_cli/batch.go
- [ ] T108 [US2] Implement exit codes (0=success, 1=invalid action, 2=IO error, 3=corrupted state, 4=tampered wall) in cmd/majhong_cli/main.go
- [ ] T109 [P] [US2] Implement file locking for concurrent access in cmd/majhong_cli/batch.go
- [ ] T110 [US2] Add integration test for batch mode single action in test/integration/batch_mode_test.go
- [ ] T111 [US2] Add integration test for batch mode action array in test/integration/batch_mode_test.go
- [ ] T112 [US2] Add integration test for game state queries in test/integration/batch_mode_test.go

### Fast Agent Mode (Protobuf Streaming)

- [ ] T113 [P] [US2] Implement FastAgent command structure in cmd/majhong_cli/main.go
- [ ] T114 [P] [US2] Implement protobuf stdin/stdout streaming in cmd/majhong_cli/fast_agent.go
- [ ] T115 [P] [US2] Implement ActionRequest handling in cmd/majhong_cli/fast_agent.go
- [ ] T116 [P] [US2] Implement ActionResponse generation in cmd/majhong_cli/fast_agent.go
- [ ] T117 [P] [US2] Add benchmark test: 10,000+ actions/min in test/integration/fast_agent_test.go

### User Story 2 Validation

- [ ] T118 [US2] Verify acceptance scenario 1: File-based action execution updates state
- [ ] T119 [US2] Verify acceptance scenario 2: Init creates valid game state file
- [ ] T120 [US2] Verify acceptance scenario 3: Multiple AI agents in sequence
- [ ] T121 [US2] Verify acceptance scenario 4: Invalid actions return error codes without modifying state
- [ ] T122 [US2] Verify acceptance scenario 5: Query outputs machine-readable JSON scores

**Checkpoint**: User Story 2 complete - Batch mode enables AI training workflows

---

## Phase 5: User Story 4 - Game State Persistence and Resume (Priority: P2)

**Goal**: Enable save/load for long-running games and analysis

**Independent Test**: Save mid-game, close CLI, reload, verify exact state restoration

**Note**: Most implementation already complete in US1 (T088-T090), this phase adds validation and edge cases

### Enhanced Persistence Features

- [ ] T123 [P] [US4] Implement game state file versioning in pkg/mahjong/serialization/json.go
- [ ] T124 [P] [US4] Implement version compatibility checking in pkg/mahjong/serialization/json.go
- [ ] T125 [P] [US4] Implement corrupted file detection and error reporting in pkg/mahjong/serialization/json.go
- [ ] T126 [US4] Add integration test for save/load portability across machines in test/integration/persistence_test.go
- [ ] T127 [US4] Add integration test for multiple save points history tracing in test/integration/persistence_test.go

### User Story 4 Validation

- [ ] T128 [US4] Verify acceptance scenario 1: All game state serialized correctly
- [ ] T129 [US4] Verify acceptance scenario 2: Resume at exact same state with correct turn
- [ ] T130 [US4] Verify acceptance scenario 3: Portability across different machines
- [ ] T131 [US4] Verify acceptance scenario 4: History tracing through multiple save files
- [ ] T132 [US4] Verify acceptance scenario 5: Corrupted file validation with clear errors

**Checkpoint**: User Story 4 complete - Robust save/load system

---

## Phase 6: User Story 3 - Majsoul Paipu Compatibility (Priority: P3)

**Goal**: Import/export Majsoul replay format for training data generation

**Independent Test**: Import Majsoul paipu, verify game state matches original replay scores

### Majsoul Integration

- [ ] T133 [P] [US3] Implement Majsoul paipu JSON parser in pkg/mahjong/majsoul/parser.go
- [ ] T134 [P] [US3] Implement Majsoul tile notation mapper (0m/0p/0s, 5mr/5pr/5sr) in pkg/mahjong/majsoul/parser.go
- [ ] T135 [US3] Implement paipu importer (timeline → action sequence) in pkg/mahjong/majsoul/importer.go
- [ ] T136 [US3] Implement wall seed extraction from paipu in pkg/mahjong/majsoul/importer.go
- [ ] T137 [US3] Implement paipu wall hash validation in pkg/mahjong/majsoul/importer.go
- [ ] T138 [US3] Implement --exit-on-tamper flag for import validation in cmd/majhong_cli/main.go
- [ ] T139 [US3] Implement paipu exporter (internal state → Majsoul format) in pkg/mahjong/majsoul/exporter.go
- [ ] T140 [P] [US3] Implement ImportPaipu CLI command in cmd/majhong_cli/main.go
- [ ] T141 [P] [US3] Implement ExportPaipu CLI command in cmd/majhong_cli/main.go
- [ ] T142 [US3] Add integration test for import with 10 real Majsoul paipu files in test/integration/majsoul_test.go
- [ ] T143 [US3] Add integration test for export + re-import verification in test/integration/majsoul_test.go

### User Story 3 Validation

- [ ] T144 [US3] Verify acceptance scenario 1: Wall seed and tile distribution converted correctly
- [ ] T145 [US3] Verify acceptance scenario 2: Step-by-step replay matches original Majsoul scores
- [ ] T146 [US3] Verify acceptance scenario 3: Export produces valid Majsoul replay files
- [ ] T147 [US3] Verify acceptance scenario 4: Red dora (0m/0p/0s) mapping works correctly
- [ ] T148 [US3] Verify acceptance scenario 5: Wall hash validation rejects tampered replays

**Checkpoint**: User Story 3 complete - Majsoul compatibility enables real-game training data

---

## Phase 7: User Story 5 - Read-Only Observer Mode (Priority: P3)

**Goal**: Enable watching live games without participating

**Independent Test**: Start 4-player game, launch observer, verify read-only display updates

### Observer Mode Implementation

- [ ] T149 [P] [US5] Implement Observer command structure in cmd/majhong_cli/main.go
- [ ] T150 [US5] Implement observer mode handler in cmd/majhong_cli/observer.go
- [ ] T151 [US5] Implement game state file polling (--follow flag) in cmd/majhong_cli/observer.go
- [ ] T152 [US5] Implement read-only display (no concealed hands) in internal/tui/observer_view.go
- [ ] T153 [US5] Implement action input rejection with error message in cmd/majhong_cli/observer.go
- [ ] T154 [US5] Add integration test for observer disconnection in test/integration/observer_test.go

### User Story 5 Validation

- [ ] T155 [US5] Verify acceptance scenario 1: Observer sees public game state (no concealed hands)
- [ ] T156 [US5] Verify acceptance scenario 2: Display updates in real-time with game changes
- [ ] T157 [US5] Verify acceptance scenario 3: Action inputs rejected with read-only message
- [ ] T158 [US5] Verify acceptance scenario 4: Observer disconnect doesn't affect game
- [ ] T159 [US5] Verify acceptance scenario 5: Multiple observers see synchronized state

**Checkpoint**: User Story 5 complete - Observer mode enables teaching and demos

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Final optimization, documentation, and deployment preparation

### Performance Optimization

- [ ] T160 Profile with go tool pprof to identify hot paths
- [ ] T161 Optimize yaku detection algorithms for performance
- [ ] T162 Verify 1000+ games/min benchmark passes
- [ ] T163 Verify fast agent mode 10,000+ actions/min benchmark passes
- [ ] T164 Run race detector: `go test -race ./...`
- [ ] T165 Verify zero memory leaks

### Documentation

- [ ] T166 [P] Write API documentation with godoc in all public packages
- [ ] T167 [P] Write README.md with installation, usage, examples
- [ ] T168 [P] Write ARCHITECTURE.md with design diagrams
- [ ] T169 [P] Write CONTRIBUTING.md with development setup
- [ ] T170 [P] Create example game state files in examples/ directory
- [ ] T171 [P] Create example CLI command scripts in examples/scripts/

### Deployment & Release

- [ ] T172 Create Dockerfile with multi-stage build
- [ ] T173 Optimize Docker image size (<50MB target)
- [ ] T174 Test headless mode in Docker container
- [ ] T175 Test Kubernetes deployment with 100 concurrent games
- [ ] T176 Setup GitHub Actions release workflow (binary builds for linux/mac/windows)
- [ ] T177 Create CHANGELOG.md with v1.0.0 release notes
- [ ] T178 Tag version v1.0.0 in git
- [ ] T179 Run final acceptance test suite (all user stories)
- [ ] T180 Verify binary size <20MB

### Code Quality

- [ ] T181 Run golangci-lint and fix all issues
- [ ] T182 Verify test coverage >90% overall
- [ ] T183 Verify test coverage 100% for rules engine
- [ ] T184 Verify all contract tests still passing
- [ ] T185 Code cleanup and refactoring for clarity

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 0 (Foundation)**: No dependencies - start immediately
- **Phase 1 (Foundational)**: Depends on Phase 0 completion - BLOCKS all user stories
- **Phase 2 (Rules Engine)**: Depends on Phase 1 completion - BLOCKS all user stories
- **Phase 3 (User Story 1)**: Depends on Phase 2 completion - Can start independently
- **Phase 4 (User Story 2)**: Depends on Phase 2 completion - Can run in parallel with US1/US4/US3/US5
- **Phase 5 (User Story 4)**: Depends on Phase 3 completion (uses US1 persistence code)
- **Phase 6 (User Story 3)**: Depends on Phase 2 completion - Can run in parallel with US1/US2/US4/US5
- **Phase 7 (User Story 5)**: Depends on Phase 3 completion (uses US1 TUI code)
- **Phase 8 (Polish)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Requires Phase 2 (rules engine) - MVP foundation
- **User Story 2 (P2)**: Requires Phase 2 - Independent of US1 but can share serialization code
- **User Story 4 (P2)**: Requires US1 (extends persistence) - Sequential dependency
- **User Story 3 (P3)**: Requires Phase 2 - Independent of all other stories
- **User Story 5 (P3)**: Requires US1 (uses TUI) - Sequential dependency

### Within Each Phase

- Tasks marked [P] can run in parallel (different files, no dependencies)
- Contract tests MUST be written before implementation (Phase 0 before Phase 1)
- Core engine (Phase 1) MUST be complete before rules (Phase 2)
- All acceptance scenarios for a user story MUST pass before moving to next priority

### Parallel Opportunities

- **Phase 0**: T002-T008, T009-T015 (directory creation, schemas, contract tests)
- **Phase 1**: T016-T017, T022-T023, T033-T035 (different modules)
- **Phase 2**: T038-T041 (yaku categories), T060-T064 (special rules)
- **Phase 3 (US1)**: T066-T072 (serialization), T073-T084 (TUI components)
- **Phase 4 (US2)**: T113-T117 (fast agent mode independent of batch mode)
- **User Stories**: After Phase 2, US1/US2/US3 can theoretically be developed in parallel by different team members

---

## Implementation Strategy

### MVP First (Critical Path)

1. Complete Phase 0: Foundation (contract tests FAILING)
2. Complete Phase 1: Foundational (core engine working)
3. Complete Phase 2: Rules Engine (ALL contract tests PASSING)
4. Complete Phase 3: User Story 1 (Interactive mode - MVP!)
5. **STOP and VALIDATE**: Play complete game, verify all actions work
6. Demo/deploy MVP if ready

### Incremental Delivery After MVP

1. Phase 4: Add User Story 2 (Batch mode for AI) → Test independently → Deploy
2. Phase 5: Add User Story 4 (Enhanced persistence) → Test independently → Deploy
3. Phase 6: Add User Story 3 (Majsoul paipu) → Test independently → Deploy
4. Phase 7: Add User Story 5 (Observer mode) → Test independently → Deploy
5. Phase 8: Polish, optimize, document, release v1.0.0

### Testing Strategy Throughout

- **Red-Green-Refactor**: Contract tests fail (Red) → Implementation passes tests (Green) → Refactor for clarity
- **Continuous Integration**: Every commit runs lint, test, build in GitHub Actions
- **Integration Tests**: Add after each user story completes
- **Performance Benchmarks**: Run at end of Phase 2, Phase 4, Phase 8

---

## Success Metrics

**Functional Completeness**:
- ✅ All 52 yaku contract tests passing (100% accuracy)
- ✅ Scoring calculation matches constitution reference
- ✅ Furiten detection: 0 false positives/negatives
- ✅ Majsoul import: 100% accuracy on test fixtures
- ✅ All 5 user stories acceptance scenarios validated

**Performance**:
- ✅ Interactive mode: <100ms response time (T162)
- ✅ Batch mode: 1000+ games/min (T162)
- ✅ Fast agent mode: 10,000+ actions/min (T163)
- ✅ Container startup: <1 second (T174)

**Quality**:
- ✅ Test coverage: >90% overall (T182)
- ✅ Test coverage: 100% rules engine (T183)
- ✅ Zero memory leaks (T165)
- ✅ Lint-clean code (T181)
- ✅ All documentation complete (T166-T171)

**Deployment**:
- ✅ Docker image <50MB (T173)
- ✅ Kubernetes 100 concurrent games (T175)
- ✅ Binary size <20MB (T180)

---

## Notes

- [P] tasks = different files, no dependencies, can run in parallel
- [Story] label maps task to specific user story (US1-US5) for traceability
- Contract tests (Phase 0) MUST be written FIRST and FAIL before implementation
- Each user story has independent acceptance validation at end of phase
- Stop at any checkpoint to validate story independently before proceeding
- Commit after each task or logical group of tasks
- Reference constitution files for all game logic implementation:
  - constitution_majhong_rule.md (core rules, wall structure, furiten)
  - constitution_majhong_yaku.md (52 yaku patterns)
  - constitution_majhong_score.md (han/fu tables, pao, oka)
  - constitution_majhong_penalties.md (chombo rules)
  - constitution_majhong_term.md (terminology EN/KR/JP/CN)

---

**Total Tasks**: 185 tasks across 8 phases
**Estimated Timeline**: 10 weeks (per plan.md)
**MVP Milestone**: Phase 3 completion (User Story 1)
**Ready for AI Training**: Phase 4 completion (User Story 2)
**Full Feature Set**: Phase 7 completion (All 5 user stories)
