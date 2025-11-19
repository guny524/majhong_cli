---
description: "Task list for CLI-Based Riichi Mahjong Game Engine"
---

# Tasks: CLI-Based Riichi Mahjong Game Engine

**Input**: Design documents from `/specs/001-cli-mahjong-engine/`
**Prerequisites**: plan.md ✓, spec.md ✓

**Tests**: Tests ARE included per plan.md TDD approach (Phase 0: Contract Tests)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

Project uses Go structure at repository root:
- `cmd/majhong_cli/` - CLI entry point
- `pkg/mahjong/` - Core game engine library
- `internal/` - Internal packages (TUI, config)
- `proto/` - Protobuf schemas
- `test/` - Test files

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [ ] T001 Initialize Go module with `go mod init` in repository root
- [ ] T002 [P] Create directory structure per plan.md: cmd/, pkg/, internal/, proto/, test/ directories
- [ ] T003 [P] Setup Makefile with targets: build, test, proto, lint, clean
- [ ] T004 [P] Configure golangci-lint with .golangci.yml configuration file
- [ ] T005 [P] Setup GitHub Actions CI/CD pipeline in .github/workflows/ci.yml
- [ ] T006 [P] Create Dockerfile for headless container deployment
- [ ] T007 [P] Add dependencies to go.mod: kong (CLI), tcell/v2 (TUI), protobuf
- [ ] T008 Create initial README.md with project overview and build instructions

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Protobuf Schema Definition

- [ ] T009 [P] Define proto/game_state.proto with GameState, Round, Player, Wall, Tile messages
- [ ] T010 [P] Define proto/action.proto with ActionRequest, ActionResponse, InitGameRequest messages
- [ ] T011 [P] Define proto/majsoul.proto for Majsoul paipu format compatibility
- [ ] T012 Generate Go code from protobuf schemas using `make proto` target

### Contract Tests (TDD Foundation) - MUST FAIL

- [ ] T013 [P] Create test/contract/yaku_test.go with 52 yaku test cases (MUST FAIL)
- [ ] T014 [P] Create test/contract/scoring_test.go with han/fu calculation tests (MUST FAIL)
- [ ] T015 [P] Create test/contract/furiten_test.go with furiten detection tests (MUST FAIL)
- [ ] T016 [P] Create test/contract/wall_test.go with wall integrity tests (MUST FAIL)
- [ ] T017 [P] Create test/fixtures/yaku_examples/ directory with sample hand JSON files

### Core Engine Foundation

- [ ] T018 [P] Create pkg/mahjong/engine/tile.go with Tile struct and constants (Suit, Rank enums)
- [ ] T019 [P] Create pkg/mahjong/engine/wall.go with Wall struct, shuffle, draw, hash validation
- [ ] T020 [P] Create pkg/mahjong/engine/player.go with Player struct, hand management
- [ ] T021 [P] Create pkg/mahjong/engine/action.go with Action types and validation
- [ ] T022 Create pkg/mahjong/engine/game.go with Game struct and ApplyAction method
- [ ] T023 Create pkg/mahjong/engine/round.go with Round state management
- [ ] T024 [P] Create internal/config/rules.go with RuleConfig struct for game configuration

**Checkpoint**: Foundation ready - contract tests defined (failing), core types implemented

---

## Phase 3: User Story 1 - Interactive Local Game Play (Priority: P1) 🎯 MVP

**Goal**: Complete 4-player Riichi Mahjong game in terminal with real-time input and visual display

**Independent Test**: Launch CLI in interactive mode, play through East-1 round completing all actions (draw, discard, pon, chi, kan, riichi, tsumo, ron), verify scoring and game state display

### Core Engine Implementation (Pass Contract Tests)

- [ ] T025 [US1] Implement Game.ValidActions() for legal action generation in pkg/mahjong/engine/game.go
- [ ] T026 [US1] Implement immutable state transitions in Game.ApplyAction() in pkg/mahjong/engine/game.go
- [ ] T027 [US1] Implement Wall.Draw() and Wall.ValidateHash() in pkg/mahjong/engine/wall.go
- [ ] T028 [US1] Implement turn order enforcement (East→South→West→North) in pkg/mahjong/engine/round.go
- [ ] T029 [US1] Add deterministic replay with seeded RNG in pkg/mahjong/engine/game.go

### Rules Engine Implementation

- [ ] T030 [P] [US1] Create pkg/mahjong/rules/yaku.go with YakuDetector for 37 regular yaku
- [ ] T031 [P] [US1] Implement yakuman detection (13 orphans, etc.) in pkg/mahjong/rules/yaku.go
- [ ] T032 [P] [US1] Create pkg/mahjong/rules/scoring.go with CalculateScore (han/fu table)
- [ ] T033 [P] [US1] Create pkg/mahjong/rules/furiten.go with CheckFuriten (discard/temporary/riichi)
- [ ] T034 [P] [US1] Create pkg/mahjong/rules/pao.go with responsibility payment logic
- [ ] T035 [US1] Create pkg/mahjong/rules/validation.go for action validation rules
- [ ] T036 [US1] Verify all 52 contract tests in test/contract/yaku_test.go PASS
- [ ] T037 [US1] Verify scoring contract tests in test/contract/scoring_test.go PASS
- [ ] T038 [US1] Verify furiten contract tests in test/contract/furiten_test.go PASS

### TUI Implementation

- [ ] T039 [P] [US1] Create internal/tui/renderer.go with tcell-based game state display
- [ ] T040 [P] [US1] Create internal/tui/input.go with keyboard input handling and tile selection
- [ ] T041 [P] [US1] Create internal/tui/layout.go with terminal layout design (hand, dora, discards, scores)
- [ ] T042 [US1] Implement headless mode detection in internal/tui/renderer.go (TERM=dumb)
- [ ] T043 [US1] Add display update performance optimization (<16ms per frame) in internal/tui/renderer.go

### CLI Entry Point for Interactive Mode

- [ ] T044 [US1] Create cmd/majhong_cli/main.go with kong CLI parser structure
- [ ] T045 [US1] Create cmd/majhong_cli/interactive.go with interactive mode handler
- [ ] T046 [US1] Implement game loop with TUI rendering and user input in cmd/majhong_cli/interactive.go
- [ ] T047 [US1] Add auto-save functionality (every turn) in cmd/majhong_cli/interactive.go
- [ ] T048 [US1] Add help command displaying rules and input syntax in cmd/majhong_cli/interactive.go

### Integration Tests for User Story 1

- [ ] T049 [P] [US1] Create test/integration/interactive_test.go with scripted TUI interaction tests
- [ ] T050 [US1] Add full game simulation test (East-1 through completion) in test/integration/interactive_test.go

**Checkpoint**: User Story 1 complete - can play full interactive game with correct yaku/scoring

---

## Phase 4: User Story 4 - Game State Persistence and Resume (Priority: P2)

**Goal**: Save/load game states to files for resume and sharing

**Independent Test**: Start game, play several turns, save to file, exit CLI, reload CLI with --load, verify all state (hands, wall, discards, dora, scores) restored exactly

**Note**: Prioritized before US2/US3 because both depend on serialization

### Serialization Layer

- [ ] T051 [P] [US4] Create pkg/mahjong/serialization/json.go with JSONSerializer (Marshal/Unmarshal)
- [ ] T052 [P] [US4] Create pkg/mahjong/serialization/protobuf.go with ProtoSerializer
- [ ] T053 [P] [US4] Create pkg/mahjong/serialization/converter.go with JSON↔Protobuf bidirectional conversion
- [ ] T054 [P] [US4] Create pkg/mahjong/serialization/hash.go with SHA-256 wall hash generation/validation
- [ ] T055 [US4] Implement pretty-printing for human-readable JSON in pkg/mahjong/serialization/json.go
- [ ] T056 [US4] Add game state versioning for forward/backward compatibility in pkg/mahjong/serialization/json.go

### File Operations

- [ ] T057 [US4] Add --save flag handler with atomic file writes in cmd/majhong_cli/interactive.go
- [ ] T058 [US4] Add --load flag handler with state restoration in cmd/majhong_cli/interactive.go
- [ ] T059 [US4] Implement auto-save to ~/.majhong_cli/autosave.json in cmd/majhong_cli/interactive.go
- [ ] T060 [US4] Add file format validation and migration in pkg/mahjong/serialization/json.go

### Tests for User Story 4

- [ ] T061 [P] [US4] Create test/integration/persistence_test.go with save/load round-trip tests
- [ ] T062 [P] [US4] Add JSON↔Protobuf conversion tests in pkg/mahjong/serialization/converter_test.go
- [ ] T063 [US4] Add wall hash validation tests (tamper detection) in test/contract/wall_test.go

**Checkpoint**: User Story 4 complete - game state fully serializable and restorable

---

## Phase 5: User Story 2 - Batch Mode for AI Integration (Priority: P2)

**Goal**: Programmatic game execution via CLI arguments with file-based state for AI training

**Independent Test**: Run `majhong_cli batch --init --game-file test.json --seed 42`, then execute multiple `--action` commands, verify game state file updates correctly and reaches completion

### Batch Mode CLI Implementation

- [ ] T064 [US2] Create cmd/majhong_cli/batch.go with batch mode command handler
- [ ] T065 [US2] Implement --init flag with seeded game initialization in cmd/majhong_cli/batch.go
- [ ] T066 [US2] Implement --action flag for single action execution in cmd/majhong_cli/batch.go
- [ ] T067 [US2] Implement --actions flag for batch action array execution in cmd/majhong_cli/batch.go
- [ ] T068 [US2] Implement --query flag for game state queries (valid-actions, scores, etc.) in cmd/majhong_cli/batch.go
- [ ] T069 [US2] Add --json flag for machine-readable output in cmd/majhong_cli/batch.go
- [ ] T070 [US2] Implement atomic file writes with temp+rename in cmd/majhong_cli/batch.go
- [ ] T071 [US2] Add proper exit codes (0=success, 1=invalid, 2=IO error, 3=corrupt, 4=tamper) in cmd/majhong_cli/batch.go

### Fast Agent Mode (Protobuf Streaming)

- [ ] T072 [P] [US2] Create cmd/majhong_cli/fast_agent.go with protobuf stdin/stdout streaming
- [ ] T073 [US2] Implement ActionRequest/ActionResponse handling in cmd/majhong_cli/fast_agent.go
- [ ] T074 [US2] Add performance optimization for 10,000+ actions/min in cmd/majhong_cli/fast_agent.go

### Tests for User Story 2

- [ ] T075 [P] [US2] Create test/integration/batch_mode_test.go with full game automation tests
- [ ] T076 [P] [US2] Add concurrent file access tests (10+ agents) in test/integration/batch_mode_test.go
- [ ] T077 [US2] Add fast agent performance benchmark test in test/integration/fast_agent_test.go
- [ ] T078 [US2] Add contract tests for CLI JSON output schema in test/contract/cli_output_test.go

**Checkpoint**: User Story 2 complete - batch mode functional for AI training with 1000+ games/min

---

## Phase 6: User Story 3 - Majsoul Paipu Compatibility (Priority: P3)

**Goal**: Import/export Majsoul replay format for training data and analysis

**Independent Test**: Download Majsoul paipu JSON, run `majhong_cli import-paipu mjs.json --output state.json`, verify final scores match, replay step-by-step, export back to paipu format

### Majsoul Integration

- [ ] T079 [P] [US3] Create pkg/mahjong/majsoul/parser.go with paipu JSON schema parsing
- [ ] T080 [P] [US3] Create pkg/mahjong/majsoul/importer.go with timeline replay and action extraction
- [ ] T081 [P] [US3] Create pkg/mahjong/majsoul/exporter.go with internal→Majsoul format conversion
- [ ] T082 [US3] Implement tile notation mapping (0m/0p/0s for aka dora) in pkg/mahjong/majsoul/parser.go
- [ ] T083 [US3] Add wall hash validation with --exit-on-tamper flag in pkg/mahjong/majsoul/importer.go
- [ ] T084 [US3] Implement metadata extraction (player names, room type, timestamp) in pkg/mahjong/majsoul/parser.go

### CLI Commands for Import/Export

- [ ] T085 [US3] Add import-paipu command in cmd/majhong_cli/main.go with input/output/exit-on-tamper flags
- [ ] T086 [US3] Add export-paipu command in cmd/majhong_cli/main.go

### Tests for User Story 3

- [ ] T087 [P] [US3] Create test/integration/majsoul_test.go with 10 real paipu import/replay tests
- [ ] T088 [P] [US3] Add test/fixtures/majsoul_paipu/ directory with sample replay files
- [ ] T089 [US3] Add export→import round-trip verification test in test/integration/majsoul_test.go
- [ ] T090 [US3] Add aka dora mapping accuracy tests in pkg/mahjong/majsoul/parser_test.go

**Checkpoint**: User Story 3 complete - 100% accurate Majsoul paipu import/export

---

## Phase 7: User Story 5 - Read-Only Observer Mode (Priority: P3)

**Goal**: Watch live games without participating for demonstration/learning

**Independent Test**: Start 4-player game, launch 5th CLI instance with `--observe`, verify observer sees updates but cannot input actions

### Observer Mode Implementation

- [ ] T091 [US5] Create cmd/majhong_cli/observer.go with read-only game state display
- [ ] T092 [US5] Implement --observe flag with game file polling in cmd/majhong_cli/observer.go
- [ ] T093 [US5] Add --follow flag for real-time updates in cmd/majhong_cli/observer.go
- [ ] T094 [US5] Implement input rejection with "read-only" error message in cmd/majhong_cli/observer.go
- [ ] T095 [US5] Add concealed hand filtering (show only public info) in cmd/majhong_cli/observer.go

### Tests for User Story 5

- [ ] T096 [P] [US5] Create test/integration/observer_test.go with multi-process observer tests
- [ ] T097 [US5] Add observer input rejection test in test/integration/observer_test.go

**Checkpoint**: User Story 5 complete - observer mode functional for demos

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, optimization, and deployment preparation

### Performance Optimization

- [ ] T098 [P] Profile yaku detection hot paths with `go tool pprof`
- [ ] T099 [P] Optimize batch mode for 1000+ games/min target
- [ ] T100 [P] Verify fast agent mode achieves 10,000+ actions/min
- [ ] T101 Reduce binary size to <20MB with build flags

### Documentation

- [ ] T102 [P] Generate API documentation with `godoc` for pkg/mahjong/
- [ ] T103 [P] Update README.md with complete usage examples and command reference
- [ ] T104 [P] Create architecture diagrams in docs/ using Mermaid syntax
- [ ] T105 Create CHANGELOG.md with version history

### Testing & Quality

- [ ] T106 [P] Run `go test -race` to verify zero memory leaks
- [ ] T107 [P] Verify test coverage >90% overall, 100% rules engine
- [ ] T108 Run `make lint` and fix all golangci-lint warnings
- [ ] T109 [P] Add unit tests for edge cases in test/unit/

### Container & Deployment

- [ ] T110 Optimize Dockerfile with multi-stage build
- [ ] T111 Test headless mode in Docker container with TERM=dumb
- [ ] T112 Verify container startup time <1 second
- [ ] T113 Test Kubernetes deployment with 100 concurrent game containers
- [ ] T114 Setup GitHub Actions release workflow for binary builds

### Release Preparation

- [ ] T115 Create version tag v1.0.0 in git
- [ ] T116 Generate release notes from CHANGELOG.md
- [ ] T117 Build and publish release binaries for linux/macos/windows

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational (Phase 2) - Core game engine
- **User Story 4 (Phase 4)**: Depends on US1 (game engine must work first)
- **User Story 2 (Phase 5)**: Depends on US4 (needs serialization layer)
- **User Story 3 (Phase 6)**: Depends on US4 (needs serialization layer)
- **User Story 5 (Phase 7)**: Depends on US1 and US4 (needs working game + file format)
- **Polish (Phase 8)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Foundation only → Independent
- **User Story 4 (P2)**: US1 → Extends with persistence
- **User Story 2 (P2)**: US1 + US4 → Uses engine + serialization
- **User Story 3 (P3)**: US1 + US4 → Import/export needs serialization
- **User Story 5 (P3)**: US1 + US4 → Observer needs game + files

### Within Each User Story

**US1 Flow:**
1. Core engine types (T025-T029) before rules
2. Rules implementation (T030-T038) to pass contract tests
3. TUI implementation (T039-T043) in parallel with CLI
4. CLI integration (T044-T048)
5. Integration tests (T049-T050)

**US4 Flow:**
1. All serializers in parallel (T051-T054)
2. Pretty-printing and versioning (T055-T056)
3. CLI file operations (T057-T060)
4. Tests (T061-T063)

**US2 Flow:**
1. Batch CLI (T064-T071) and fast agent (T072-T074) in parallel
2. Tests (T075-T078)

### Parallel Opportunities

**Phase 1 (Setup):** T002, T003, T004, T005, T006, T007, T008 can all run in parallel

**Phase 2 (Foundational):**
- Protobuf schemas (T009, T010, T011) in parallel
- Contract tests (T013, T014, T015, T016, T017) in parallel
- Core types (T018, T019, T020, T021, T024) in parallel

**Phase 3 (US1):**
- Rules modules (T030, T031, T032, T033, T034) in parallel
- TUI components (T039, T040, T041) in parallel
- Integration tests (T049) and unit tests parallel

**Phase 4 (US4):**
- Serializers (T051, T052, T053, T054) in parallel
- Tests (T061, T062) in parallel

**Phase 8 (Polish):**
- Performance profiling (T098, T099, T100) in parallel
- Documentation (T102, T103, T104) in parallel
- Testing (T106, T107, T109) in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all rules modules together:
Task T030: "Create pkg/mahjong/rules/yaku.go with YakuDetector"
Task T031: "Implement yakuman detection in pkg/mahjong/rules/yaku.go"
Task T032: "Create pkg/mahjong/rules/scoring.go"
Task T033: "Create pkg/mahjong/rules/furiten.go"
Task T034: "Create pkg/mahjong/rules/pao.go"

# Launch all TUI components together:
Task T039: "Create internal/tui/renderer.go"
Task T040: "Create internal/tui/input.go"
Task T041: "Create internal/tui/layout.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T008)
2. Complete Phase 2: Foundational (T009-T024) - CRITICAL
3. Complete Phase 3: User Story 1 (T025-T050)
4. **STOP and VALIDATE**: Play full interactive game, verify all yaku recognized
5. Deploy/demo if ready - **This is the MVP!**

### Incremental Delivery

1. **Foundation** (Phase 1-2) → Contract tests defined, core types ready
2. **MVP** (Phase 3) → Interactive game playable, all 52 yaku working
3. **Persistence** (Phase 4) → Add save/load for resume capability
4. **AI Training** (Phase 5) → Add batch mode for ML experimentation
5. **Data Integration** (Phase 6) → Add Majsoul import/export
6. **Teaching Tool** (Phase 7) → Add observer mode
7. **Production Ready** (Phase 8) → Optimize and document

Each increment adds value without breaking previous functionality.

### Parallel Team Strategy

With multiple developers:

1. **Week 1**: Team completes Setup + Foundational together (T001-T024)
2. **Week 2-3**: Team focuses on US1 together (complex - 52 yaku)
3. **Week 4+**: Split work:
   - Developer A: US4 (Persistence) → T051-T063
   - Developer B: US2 (Batch Mode) → T064-T078
   - Developer C: US3 (Majsoul) → T079-T090
4. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies, can run in parallel
- [Story] label maps task to specific user story (US1-US5) for traceability
- Each user story is independently completable and testable
- Contract tests (Phase 2) MUST FAIL before implementation begins (TDD Red phase)
- Stop at any checkpoint to validate story independently
- Target metrics: 1000+ games/min (batch), 100% yaku accuracy, <100ms interactive response
- Priority order: P1 (US1) → P2 (US4, US2) → P3 (US3, US5)
