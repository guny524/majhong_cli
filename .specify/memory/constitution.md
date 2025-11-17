<!--
Sync Impact Report:
- Version: 1.2.1 → 1.4.0 (minor version bump: completed PDF content integration + file structure reorganization)
- Modified principles: None (documentation expansion and reorganization only)
- New sections added:
  * Penalty and Violation System Reference - NEW document constitution_majhong_penalties.md
  * Etiquette and Sportsmanship Reference - NEW document constitution_majhong_etiquette.md
  * Pao/Responsibility Payment system - Added to constitution_majhong_score.md
  * Oka scoring bonus - Added to constitution_majhong_score.md
- Modified sections:
  * Domain-Specific References - Added 2 new document references (penalties, etiquette)
  * Mahjong Rules Reference (constitution_majhong_rule.md) - Added detailed Furiten section (3 types)
  * Scoring Calculation Reference (constitution_majhong_score.md) - Added Pao responsibility payment + Oka bonus
  * Terminology Reference (constitution_majhong_term.md) - Added 11 new terms (6 penalties + 5 Pao/Oka/Uma)
  * Compliance & Review - Added penalty handling and UI/UX design compliance checks
  * Agent Guidance - Added instructions for penalty validation and etiquette-informed design
- File structure reorganization:
  * constitution_term.md → constitution_majhong_term.md
  * constitution_yaku.md → constitution_majhong_yaku.md
  * constitution_score.md → constitution_majhong_score.md
  * constitution_penalties.md → constitution_majhong_penalties.md
  * constitution_etiquette.md → constitution_majhong_etiquette.md
  * Rationale: Consistent naming convention (constitution_majhong_* for domain docs, anti_tampering keeps original name)
- Content additions:
  * Penalty system: penalty points (벌칙 점수), win prohibition (화료 불가), chombo (촌보)
  * Furiten detailed rules: discard furiten, temporary furiten, riichi furiten with implementation code
  * Etiquette guidelines: physical mahjong behavior, digital UI/UX design principles
  * Pao/Responsibility payment: Daisangen, Daisuushii, Suukantsu payment responsibility rules
  * Oka: Starting/return points differential bonus for 1st place (+20.0 typical)
  * Extended play: West/North round continuation (남입/서입) when no one reaches return points
  * New terminology: 11 terms added (penalty 6 + Pao/Oka/Uma/return points/extended play 5)
- Gap analysis result: 100% complete coverage of "한 권으로 익히는 리치 마작0426.pdf"
  * Previously missing: Pao (책임지불) and Oka (오카) - NOW ADDED ✓
  * All major sections from PDF now documented in constitution
- Source: "한 권으로 익히는 리치 마작0426.pdf" (YOSTAR/Korean Mahjong League official guide, pages 23-24, 40-43, 462-473, 781-792)
- Templates requiring updates: spec.md, plan.md templates should reference new documents when relevant
- Follow-up TODOs:
  * Implement Pao responsibility tracking in game state
  * Implement Oka final score calculation
  * Implement penalty validation system in game engine
  * Add UI/UX etiquette guidelines to frontend design
  * Write unit tests for penalty detection scenarios
  * Write unit tests for furiten detection (3 types)
  * Write unit tests for Pao trigger detection (Daisangen/Daisuushii/Suukantsu)
  * Write unit tests for Oka bonus calculation

Previous Sync Reports:
- 1.2.0 → 1.2.1: Patch version bump - corrected Korean terminology to MajSoul standard (12 yaku names)
- 1.2.1 → 1.3.0: Minor version bump - added penalty system and etiquette references
-->

# Mahjong CLI Project Constitution

## Core Principles

### I. Skill-First Development

**Rule**: Before starting ANY task or implementation work, the agent MUST:
1. Review all available skills in the skill catalog
2. Identify skills relevant to the current task
3. Activate relevant skills proactively
4. Document which skills were activated and why

**Rationale**: Prevents context loss due to auto-compaction, ensures consistent use of project-specific workflows, and maximizes efficiency by leveraging specialized capabilities. This is NON-NEGOTIABLE as the project operates under token budget constraints.

### II. Documentation-First Workflow

**Rule**: Before implementing ANY feature:
- Document overall architecture and deployment structure FIRST
- Write reference documentation to persist knowledge across context resets
- Create design artifacts (spec.md, plan.md) before code
- NEVER jump directly to detailed implementation without documenting the big picture

**Rationale**: Given auto-compact context limitations, documentation serves as persistent memory. Prevents rework when context is lost and ensures alignment on architecture before code is written.

### III. Test-Driven Development (NON-NEGOTIABLE)

**Rule**: TDD is mandatory for all features:
1. Write tests FIRST based on user requirements
2. Get user approval on tests
3. Verify tests FAIL (Red)
4. Implement to make tests pass (Green)
5. Refactor while keeping tests green (Refactor)

**Rationale**: For a complex domain like Mahjong with intricate rules (yaku, scoring, wall distribution), tests are the specification. They prevent regression and serve as executable documentation.

### IV. Library-First Architecture

**Rule**: Every feature MUST be developed as a standalone library:
- Self-contained with clear boundaries
- Independently testable
- Well-documented with clear purpose
- CLI-exposable interface
- NO organizational-only libraries (must solve a real problem)

**Rationale**: Enables modular development of: (a) CLI game engine, (b) probability calculator, (c) AI trainer, (d) web integration, (e) replay downloader. Each module can be developed, tested, and deployed independently.

### V. CLI-First Interface

**Rule**: Every library MUST expose functionality via CLI:
- Text-based I/O protocol: stdin/args → stdout
- Errors/logs → stderr
- Support both JSON (machine-readable) and human-readable formats
- Composable with standard UNIX tools

**Rationale**: Supports the project goal of a CLI-based Mahjong system. Enables integration between modules (web parser → CLI engine → probability calculator) using standard text protocols.

### VI. Evidence-Based Changes

**Rule**: NEVER make changes based on assumptions or surface-level analysis:
1. Always verify with concrete evidence (git diff, file contents, logs)
2. Investigate root causes, not just symptoms
3. Question user statements - verify with direct inspection
4. Explain: what error, what's suspected, how to fix, what will be changed
5. Document reasoning for every significant change

**Rationale**: Prevents thrashing on incorrect fixes. For a probabilistic/AI system, bugs can be subtle - must trace actual code behavior rather than guessing.

### VII. Minimal Complexity & YAGNI

**Rule**:
- Start with the simplest solution
- Add complexity ONLY when proven necessary
- Justify any abstraction, pattern, or additional project
- Default to direct implementation unless complexity is unavoidable

**Rationale**: Early-stage project with evolving requirements (Riichi vs standard Mahjong, local rules, AI approaches). Premature abstraction adds cost. Follow "you aren't gonna need it" until you do.

## Development Workflow

### Code Quality Gates

**MUST** run before considering a task complete:
1. **Linting**: Fix all linter errors
2. **Build**: Ensure successful compilation
3. **Tests**: All tests must pass (if TDD is applied)

**MUST NOT** auto-fix without reasoning:
- Do NOT blindly remove "unused" code - trace why it was added, check global usage, propose alternatives
- Do NOT apply quick fixes to satisfy compiler without understanding implications
- Verify intent, provide alternatives (keep/refactor/remove), cite evidence (file:line/commit)

### Change Management

**For Git operations**:
- Always check `git status`, `git diff`, `git log` BEFORE any destructive operation
- NEVER use `git restore` without confirming what will be lost
- When user reports changes, verify with `git diff` rather than trusting verbally

**For refactoring**:
1. Explain: what is changing and why
2. Show: expected before/after behavior
3. Validate: run tests/build to confirm
4. Document: update relevant docs (README, spec)

### Documentation Management

**README consolidation**:
- Use a SINGLE README.md at project root
- Avoid fragmentation into multiple docs unless absolutely necessary
- Apply 5W1H (Who/When/Where/What/Why/How) validation

**TODO management**:
- Active TODOs live in `todos/*.md`
- Completed work moves to `done_TODO_YYYY-MM-DD_*.md`
- README contains high-level summary only, not detailed TODOs
- Do NOT delete TODOs without explanation
- Eliminate duplicates - keep info in ONE place

**On task completion**:
- Execute doc-management skill to update all documentation
- Ensure alignment between code, specs, and README

## Technical Constraints

### Language & Performance

**Primary language**: Go (current), migrating to Rust for single-binary distribution
- **Why Go now**: Faster prototyping, good concurrency for AI training
- **Why Rust later**: Zero-cost abstractions, single binary, memory safety for production

**Probability algorithm constraint**: ≤100 lines of core logic
- Forces algorithmic clarity
- Prevents over-engineering
- Enables human verification of correctness

**Performance targets** (to be refined):
- CLI game engine: <10ms move validation
- Probability calc: <100ms for full hand analysis
- AI inference: <500ms per move recommendation

### Domain-Specific Requirements

**Mahjong rule configurability**:
- Support Riichi Mahjong (primary)
- Support standard Mahjong
- Toggle local rules (yakuman variations, red fives, etc.)
- Explicitly document which ruleset is active

**Wall integrity verification**:
- Must support SHA-256 hash verification (as used by MajSoul)
- Validate wall composition: 136 tiles, proper distribution
- Document wall structure: dead wall (14), live wall (122), dora/ura-dora, rinshan

**AI/Probability requirements**:
- Explainable recommendations (natural language + probability %)
- Japanese Mahjong terminology for yaku names
- Decision tree output in markdown tables

## Domain-Specific References

**IMPORTANT**: The following reference documents provide detailed specifications for Mahjong domain implementation. Developers MUST consult these documents when implementing corresponding features.

### Terminology Reference

**Document**: `constitution_majhong_term.md`

**Contents**:
- Comprehensive glossary (100+ terms in EN/KR/JP/CN format)
- Core game concepts (dealer/non-dealer, hand states, winning methods)
- Tile categories and suits (number tiles, honor tiles, terminals/simples)
- Hand components (sequences, triplets, quads, pairs, melds)
- Calling actions (chi/pon/kan variations)
- Wall components (dead wall, live wall, dora, rinshan, haitei/houtei)
- Scoring terms (yaku, han, fu, mangan, yakuman)
- Game phases and end conditions
- Common confusion points clarified
- Code naming conventions and usage guidelines

**When to reference**:
- ANY time you encounter unfamiliar Mahjong term
- When naming variables, functions, or enums
- When writing documentation or comments
- When implementing UI with localized text

### Mahjong Rules Reference

**Document**: `constitution_majhong_rule.md`

**Contents**:
- Tile composition (136 tiles: man/pin/sou/honors)
- Tile notation standard (1m-9m, 1p-9p, 1s-9s, 1z-7z, 0m/0p/0s)
- Game setup procedures (seat determination, initial distribution)
- Wall structure (dead wall 14 tiles, live wall 122 tiles)
- Turn structure and calling actions (chi/pon/kan/ron)
- Riichi declaration requirements and effects
- Winning conditions (tenpai, minimum yaku requirement)
- Game end conditions (exhaustive draw, abortive draws)
- Implementation notes (state tracking, validation rules)
- Terminology mapping (English/Japanese/Korean/Code)

**When to reference**:
- Implementing game engine core logic
- Validating hand compositions
- Calculating probabilities for reaching tenpai
- Implementing calling action handlers

### Wall Construction Reference

**Document**: `constitution_majhong_wall.md`

**Contents**:
- Physical wall construction procedures (offline Mahjong)
- Wall break point determination (dice rolling algorithm)
- Dead wall formation (rinshan/dora/ura-dora layout)
- Initial hand distribution (haipai: 14/13/13/13 tiles)
- Hand organization techniques
- Point stick distribution (25,000 starting points)
- Digital implementation (wall state machine, dead wall management)
- Validation rules (dead wall always 14 tiles, deterministic break calculation)
- Testing requirements (unit tests, integration tests)

**When to reference**:
- Implementing wall shuffling and distribution
- Building dice roll mechanics
- Creating dead wall management system
- Validating wall integrity during game

### Anti-Tampering Verification Reference

**Document**: `constitution_anti_tampering.md`

**Contents**:
- SHA-256 cryptographic verification system
- Two-stage hash commitment scheme (H(wall_code + salt), H(salt))
- Wall code format (136 tiles × 2-char encoding)
- Server-side generation and hash computation
- Client-side verification procedures
- Game flow integration (start/mid-game/end)
- Security considerations (attack vectors, mitigations)
- Implementation best practices (enhanced security measures)
- Testing requirements (hash determinism, tampering detection)
- Migration notes (MD5 → SHA-256 evolution)

**When to reference**:
- Implementing multiplayer online mode
- Building wall generation for network games
- Creating verification UI (mid-game hash display)
- Implementing post-game replay verification
- Designing anti-cheat measures

### Yaku (Scoring Patterns) Reference

**Document**: `constitution_majhong_yaku.md`

**Contents**:
- Complete catalog of all yaku patterns (1-han through yakuman)
- 1-han yaku: riichi, ippatsu, menzen tsumo, tanyao, yakuhai, pinfu, iipeikou, haitei, houtei, rinshan, chankan
- 2-han yaku: chiitoitsu, toitoi, sanankou, sanshoku, ittsu, chanta, sankantsu, honroutou, shousangen, double riichi
- 3-han yaku: ryanpeikou, junchan, honitsu
- 6-han yaku: chinitsu
- Yakuman: kokushi, suuankou, daisangen, shousuushii, daisuushii, tsuuiisou, ryuuiisou, chinroutou, chuuren, suukantsu, tenhou, chiihou
- Local yaku: renhou, nagashi mangan, paarenchan (with MajSoul compatibility notes)
- Conditions, requirements, open/closed restrictions for each yaku
- Yaku combination rules (cannot combine, commonly combined)
- Detection algorithms and validation logic
- Testing requirements per yaku

**When to reference**:
- Implementing yaku detection system
- Validating winning hands (minimum yaku requirement)
- Calculating han values
- Building probability calculator (which yaku are achievable)
- Creating AI training features (yaku pattern recognition)

### Scoring Calculation Reference

**Document**: `constitution_majhong_score.md`

**Contents**:
- Complete han-fu scoring tables (dealer and non-dealer)
- Fu calculation rules (base 20 fu + hand composition bonuses)
- Fu sources: win method, melds (sequences/triplets/quads), pair, wait type
- Limit hands: mangan (5 han), haneman (6-7), baiman (8-10), sanbaiman (11-12), yakuman (13+)
- Payment calculations (ron vs tsumo, dealer vs non-dealer)
- Honba and riichi stick bonuses
- Tenpai settlement at exhaustive draw
- Game end conditions and round progression
- Dealer continuation rules
- Final ranking and uma (placement bonus)
- Complete scoring algorithm with code examples
- Testing requirements for score calculation

**When to reference**:
- Implementing scoring engine
- Calculating payments between players
- Determining game end conditions
- Building score display UI
- Validating winning hands (fu calculation)
- Creating score prediction features

### Penalty and Violation System Reference

**Document**: `constitution_majhong_penalties.md`

**Contents**:
- Penalty application period (after dealer's first discard until win/draw)
- Three penalty types: penalty points (벌칙 점수), win prohibition (화료 불가), chombo (촌보)
- Penalty points: 1000 point deposit for minor procedural errors
- Win prohibition scenarios: incorrect calls, wrong tile count, illegal swaps, false declarations
- Chombo scenarios: no-tenpai riichi at draw, winning without yaku, furiten ron, illegal kan after riichi
- Chombo payment structure: dealer pays 4000 to all (traditional) or 3000 to all (modern)
- Hand reset procedures after chombo
- Penalty severity decision tree
- Implementation validation checkpoints
- Testing requirements for penalty detection

**When to reference**:
- Implementing game rule validation system
- Validating riichi declarations (tenpai check)
- Validating win declarations (yaku check, furiten check, tile count)
- Implementing kan after riichi validation
- Handling rule violations and applying penalties
- Creating error handling for invalid actions
- Building replay audit system

### Etiquette and Sportsmanship Reference

**Document**: `constitution_majhong_etiquette.md`

**Contents**:
- Physical mahjong etiquette (tile handling, voice protocol, discard layout)
- Timing rules: calling windows, late calling prohibition, turn sequence
- Win declaration protocol (voice before reveal)
- Point exchange etiquette (proper announcement order, respectful transfer)
- Game end procedures (score confirmation, pre-final-hand verification)
- Digital implementation guidelines (UI/UX informed by physical etiquette)
- Communication features (allowed/restricted chat, emotes, report system)
- Timing fairness (action timers, AFK detection, progressive penalties)
- Sportsmanship enforcement (violations, warnings, penalties)
- Cultural context (respect, fairness, grace, enjoyment)

**When to reference**:
- Designing game UI/UX (discard layout, timer display, turn indicators)
- Implementing calling windows and action timing
- Building communication features (chat, emotes, reports)
- Creating win declaration flow (confirmation dialogs)
- Implementing AFK detection and timeout handling
- Designing point transfer animations
- Building score confirmation UI
- Implementing player behavior enforcement system

### Reference Document Usage Protocol

**MUST verify before implementation**:
1. Identify which reference document(s) apply to current feature
2. Read relevant sections thoroughly
3. Extract specific requirements (validation rules, data structures, algorithms)
4. Implement according to specifications
5. Write tests based on reference document's testing requirements
6. Validate implementation against reference document's examples

**MUST update when domain rules change**:
- If terminology needs clarification → update `constitution_majhong_term.md`
- If Mahjong rules interpretation changes → update `constitution_majhong_rule.md`
- If wall construction algorithm changes → update `constitution_majhong_wall.md`
- If security protocols change → update `constitution_anti_tampering.md`
- If yaku conditions or new yaku added → update `constitution_majhong_yaku.md`
- If scoring calculations or limits change → update `constitution_majhong_score.md`
- If penalty rules or violation handling changes → update `constitution_majhong_penalties.md`
- If etiquette guidelines or UI/UX best practices change → update `constitution_majhong_etiquette.md`
- Always increment Last Updated date
- Document rationale for changes

## Governance

### Amendment Process

**Constitution changes require**:
1. Document proposed change with rationale
2. Update version number per semantic versioning:
   - MAJOR: Breaking changes to principles or workflow
   - MINOR: New principles or expanded guidance
   - PATCH: Clarifications, typos, non-semantic fixes
3. Propagate changes to all dependent templates (spec, plan, tasks)
4. Create migration plan if existing work is affected

### Compliance & Review

**Every task/PR must verify**:
- Skill activation performed (Principle I)
- Documentation created before code (Principle II)
- Tests written and failing before implementation (Principle III)
- Library boundaries respected (Principle IV)
- CLI interface provided (Principle V)
- Evidence-based reasoning documented (Principle VI)
- Complexity justified if added (Principle VII)

**Domain-specific compliance**:
- Terminology used consistently per `constitution_majhong_term.md` (EN/KR/JP/CN format)
- Mahjong rules implemented per `constitution_majhong_rule.md` (includes Furiten detailed rules)
- Wall construction follows `constitution_majhong_wall.md` algorithms
- Online multiplayer uses anti-tampering from `constitution_anti_tampering.md`
- Yaku detection follows conditions in `constitution_majhong_yaku.md`
- Scoring calculations use tables from `constitution_majhong_score.md`
- Penalty handling follows severity levels in `constitution_majhong_penalties.md`
- UI/UX design informed by guidelines in `constitution_majhong_etiquette.md`

**Complexity justification table** (in plan.md):
| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [abstraction] | [specific need] | [why direct approach insufficient] |

### Agent Guidance

**For Claude and similar agents**:
- Follow global instructions in `~/.claude/CLAUDE.md`
- This constitution OVERRIDES default behaviors when conflicts arise
- Use project-specific skills (doc-management, clean-code, react-next-guide, etc.) as applicable
- Always think step-by-step and validate understanding before acting
- Ask clarifying questions when requirements are ambiguous
- Respond with minimal necessary sentences while including sufficient context

**When implementing Mahjong features**:
- ALWAYS consult relevant reference documents first
- Verify implementation against reference document specifications
- Include reference document citations in code comments (e.g., `// See constitution_majhong_rule.md: Tile Notation Standard`)
- If reference document is ambiguous, ask for clarification before implementing
- Use terminology consistently: `constitution_majhong_term.md` for all term definitions
- Reference yaku by code enum: `constitution_majhong_yaku.md` for yaku detection
- Use scoring tables: `constitution_majhong_score.md` for han-fu calculations
- Implement penalty validation: `constitution_majhong_penalties.md` for rule violation handling
- Design UI/UX with etiquette in mind: `constitution_majhong_etiquette.md` for player interaction guidelines

**Version**: 1.4.0 | **Ratified**: 2025-11-17 | **Last Amended**: 2025-11-18
