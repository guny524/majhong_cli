# CLI-Based Riichi Mahjong Game Engine

A terminal-based Riichi Mahjong game engine designed for interactive play and AI agent training.

## Features

**Interactive TUI Gameplay**
- Full 4-player Riichi Mahjong in terminal
- Real-time keyboard input with tcell
- Complete yaku detection (52 types including yakuman)

**AI Training Support**
- Batch mode for automated game simulation
- Fast agent mode (10,000+ actions/min via protobuf)
- Deterministic replay with seeded RNG

**Data Integration**
- Majsoul paipu import/export
- JSON/Protobuf serialization
- Wall hash validation

**Observer Mode**
- Watch live games without participation
- Read-only state display

## Prerequisites

- Go 1.21 or higher
- protoc (Protocol Buffers compiler)
- make

## Build

```bash
# Install dependencies
go mod download

# Build binary
make build

# Run tests
make test

# Generate protobuf code
make proto
```

## Quick Start

```bash
# Interactive mode (coming in Phase 3)
./bin/majhong_cli interactive

# Batch mode (coming in Phase 5)
./bin/majhong_cli batch --init --seed 42

# Import Majsoul replay (coming in Phase 6)
./bin/majhong_cli import-paipu replay.json
```

## Project Structure

```
majhong_cli/
├── cmd/majhong_cli/     # CLI entry point
├── pkg/mahjong/         # Core game engine library
│   ├── engine/          # Game state, actions, rules
│   ├── rules/           # Yaku detection, scoring
│   ├── serialization/   # JSON/Protobuf converters
│   └── majsoul/         # Majsoul paipu integration
├── internal/            # Internal packages
│   ├── tui/             # Terminal UI (tcell)
│   └── config/          # Configuration
├── proto/               # Protobuf schemas
└── test/                # Tests
    ├── contract/        # Contract tests (TDD)
    ├── integration/     # Integration tests
    └── fixtures/        # Test data
```

## Development Status

**Phase 1: Setup** ✅ (Current)
- [x] Project initialization
- [x] Directory structure
- [x] Build system (Makefile)
- [x] CI/CD pipeline

**Phase 2: Foundational** (Next)
- [ ] Protobuf schemas
- [ ] Contract tests (TDD)
- [ ] Core engine types

**Phase 3: MVP - Interactive Game** (Planned)
- [ ] Rules engine (52 yaku)
- [ ] TUI implementation
- [ ] CLI entry point

See [tasks.md](specs/001-cli-mahjong-engine/tasks.md) for complete roadmap.

## License

MIT
