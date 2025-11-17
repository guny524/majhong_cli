<!--
Sync Impact Report:
- Version: 0.0.0 → 1.0.0 (initial ratification)
- Modified principles: N/A (initial version)
- Added sections: All core principles and governance
- Removed sections: N/A
- Templates requiring updates:
  ✅ plan-template.md - Constitution Check section already present
  ✅ spec-template.md - No constitution-specific constraints needed
  ✅ tasks-template.md - Task categorization aligns with principles
- Follow-up TODOs: None
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

**Version**: 1.0.0 | **Ratified**: 2025-11-17 | **Last Amended**: 2025-11-17
