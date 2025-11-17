# Mahjong Etiquette and Sportsmanship Guidelines

**Purpose**: Defines player behavior expectations, etiquette guidelines, and best practices for physical and digital Riichi Mahjong implementation.
**Referenced by**: `.specify/memory/constitution.md`
**Last Updated**: 2025-11-18
**Source**: "한 권으로 익히는 리치 마작0426.pdf" (YOSTAR/Korean Mahjong League official guide, pages 42-43)

## Purpose and Scope

This document provides etiquette guidelines for:
- **Physical mahjong games**: Tile handling, voice protocol, sportsmanship
- **Digital implementations**: UI/UX design informed by physical etiquette norms
- **Online multiplayer**: Player behavior expectations and chat conduct

**Note**: While some guidelines are specific to physical play, understanding them informs better digital game design (e.g., clear turn indicators, proper discard display, timing fairness).

## 1. Game Start and Automatic Table Operation

**Applicable to**: Physical mahjong with automatic shuffling tables

### 1.1 Table Control Authority

**Rule**: The dealer (親/oya) of the starting hand operates the automatic table controls.

**Rationale**: Centralized control prevents confusion and accidental button presses.

**Digital equivalent**: Dealer initiates game start UI (countdown, shuffle animation).

### 1.2 Wall Positioning

**Rule**: When wall rises from table, push it slightly diagonally to the right for easier drawing.

**Visual guide** (from PDF):
```
    ┌─────────┐
    │  WALL   │ ← Push diagonally right →
    └─────────┘
    Player's position
```

**Reason**: Ergonomic positioning allows smooth tile drawing without knocking over wall.

**Digital equivalent**: Wall animation should visually indicate "active draw zone" for clarity.

### 1.3 Initial Distribution (Haipai)

**Rule**: Dealer takes final **two tiles together** (not one at a time) during initial distribution.

**Reason**: Traditional protocol matching the 14-tile dealer starting hand.

**Caution for automatic tables**: If table auto-distributes tiles, dealer must remember to **draw first tile** at game start (some tables don't auto-draw dealer's first tsumo).

**Digital implementation note**:
```python
def start_hand(game):
    """Initialize hand with proper dealer tile count."""
    distribute_initial_tiles(game)  # Dealer gets 14, others 13
    # Dealer's turn starts immediately (already has 14 tiles)
    # No need to draw first tile
```

## 2. In-Game Behavior

### 2.1 Tile Handling - Avoiding Damage

**Physical handling guidelines**:

#### (1) Avoid Violent Tile Handling

**Don't**:
- Dragging tiles on table surface when declaring tsumo
- Slamming tiles against table edge
- Forcefully slapping down tiles during discard (강타 행위)

**Reason**: Repeated impact damages mahjong tiles (most common cause of tile wear).

**Do instead**:
- Lift tiles cleanly when declaring win
- Place discards gently on table
- Handle tiles with care (they are precision-molded game pieces)

**Digital equivalent**: Animations should reflect respectful tile handling (smooth transitions, not violent throws).

#### (2) Proper Discard Arrangement

**Standard layout**: 6 tiles per row, arranged chronologically.

**Example discard area**:
```
Row 1: [1m] [5p] [7s] [東] [3m] [9s]
Row 2: [2p] [6m] [4s] [南] [1p] [8m]
                                 ↑ newest discard
```

**Reason**:
- Clear chronological order for furiten checking
- Easy visual parsing for all players
- Standard convention across physical mahjong

**Digital implementation**: Display discards in grid format (6 per row) with newest tile highlighted.

#### (3) No Strategic Table Talk

**Forbidden during gameplay**:
- Complaining about your tiles ("My hand is terrible!")
- Comments suggesting your hand state ("I'm so close!")
- External distractions intended to confuse opponents
- Psychological manipulation through voice/behavior

**Reason**: Maintains fair play and prevents information leakage.

**Acceptable communication**:
- Clear call declarations (chi/pon/kan/riichi/tsumo/ron)
- Point announcements
- Polite game-related questions ("How many tiles remain?")

**Digital equivalent**:
- Disable chat during active gameplay (or restrict to preset phrases)
- No emoji spam or distracting animations
- Mute/report options for disruptive players

#### (4) One-Handed Tile Operations

**Rule**: Draw (tsumo) and discard using **only one hand**.

**Reason**:
- Prevents "invisible switching" (using two hands to manipulate tiles covertly)
- Traditional protocol enforcing fair play
- Easier to observe for rule compliance

**Exception**: Organizing hand tiles during setup can use both hands.

**Digital equivalent**: Not directly applicable, but informs timing design (no actions allowed mid-discard animation).

#### (5) No Post-Game Wall Inspection

**Forbidden after hand ends**:
- Flipping over dead wall (rinshan, ura-dora indicators)
- Checking remaining wall tiles ("What if I had drawn...?")
- Inspecting what tiles would have come next

**Reason**:
- Slows down game flow
- Leads to regret/frustration rather than forward play
- Some consider it bad luck superstition

**Acceptable**: Reviewing **already-revealed** tiles (discards, called melds, dora indicators).

**Digital equivalent**:
- Don't auto-reveal unrevealed tiles after hand
- Offer "View Full Hand Replay" as opt-in feature (not default)
- Maintain mystery for better game pacing

#### (6) Device Usage During Game

**Rule**: Minimize phone usage and other distractions during active gameplay.

**If urgent matter arises**:
- Politely request brief pause
- Wait for appropriate moment (not mid-discard sequence)
- Communicate estimated time ("One minute, please")

**Reason**: Respects other players' time and maintains game flow.

**Digital equivalent**:
- Implement AFK detection (auto-discard if player doesn't act in time)
- Notify players when someone is taking excessive time
- Balance: Allow reasonable thinking time, but prevent stalling

## 3. Voice Protocol (Calling Declarations)

### 3.1 Late Calling is Prohibited

**Rule**: Once you touch the wall to draw, you **cannot** declare chi/pon/kan on the previous discard.

**Timing sequence**:
```
Opponent discards tile → [Calling window: 2-3 seconds]
                      ↓
Your hand touches wall → Calling window CLOSED
```

**Reason**: Clear timing rules prevent disputes and manipulation.

**Digital implementation**:
```python
def on_tile_discarded(tile, discarding_player, game_state):
    """Open calling window for other players."""
    calling_deadline = current_time() + CALLING_WINDOW_SECONDS

    for player in other_players:
        enable_calling_buttons(player, tile, deadline=calling_deadline)

def on_player_touches_wall(player):
    """Close calling window when next player draws."""
    disable_calling_buttons_for_previous_discard()
    player.draw_tile()
```

**Exception**: Ron always has priority (can be declared even if next player started drawing, with resolution rollback).

### 3.2 Win Declaration Protocol

**Proper sequence for declaring win**:

1. **First**: Clearly announce "쯔모" (tsumo) or "론" (ron)
2. **Second**: Reveal hand tiles face-up
3. **Third**: Announce point value ("천, 이천" = 1000/2000)

**Critical**: Voice declaration MUST come before revealing tiles.

**Reason**:
- Prevents "trial reveals" (checking if hand is valid before committing)
- Clear communication that win is being declared
- Follows traditional protocol

**Common mistake**: Revealing tiles first, then saying "tsumo" → May be penalized as unclear declaration.

**Special note on "멘젠쯔모" (Menzen Tsumo) yaku**:
- This is a **yaku** (scoring pattern for closed-hand tsumo)
- The voice call "쯔모" is **always required** regardless of whether you have this yaku
- Don't confuse the yaku name with the voice declaration requirement

**Digital implementation**:
- Require explicit button press for "Tsumo" or "Ron"
- Display "Declare Tsumo?" confirmation before revealing hand
- Auto-validate hand before allowing declaration (prevent invalid wins)

### 3.3 Covering Tiles = No-Tenpai Declaration

**Implied meaning**: Placing hand face-down implicitly declares **no-tenpai**.

**When this matters**: At exhaustive draw (유국), players in tenpai reveal hands to claim settlement.

**Rule**: Don't cover tiles unless you are genuinely not in tenpai.

**Consequence of false no-tenpai**:
- If later discovered you were tenpai → Penalty points
- Forfeits tenpai settlement payment (1000 points from non-tenpai players)

**Digital equivalent**:
- Explicit "Tenpai" / "No-Tenpai" button at exhaustive draw
- No ambiguous gestures

## 4. Point Exchange Protocol

### 4.1 Tsumo Point Announcement Order

**Standard format**: "자에게 받을 점수 → 친에게 받을 점수"
- First: Points from non-dealers (子/ko)
- Second: Points from dealer (親/oya)

**Examples**:

| Hand Value | Announcement (Non-dealer wins) |
|------------|-------------------------------|
| 30 fu 3 han | "천, 이천" (1000 from ko, 2000 from oya) |
| 40 fu 3 han | "천삼백, 이천육백" (1300 from ko, 2600 from oya) |

**With honba (連荘 counters)**:
1. Announce base points first
2. Then announce total with honba added

**Example**: 3 honba, 20 fu 3 han tsumo:
- Base: 700 from ko, 1300 from oya
- Announcement: "칠백, 천삼백은 천, 천육백"
  - Translation: "700, 1300... (with honba:) 1000, 1600"

**Digital equivalent**: Display point breakdown clearly:
```
Non-dealer pays: 1000 (700 base + 300 honba)
Dealer pays: 1600 (1300 base + 300 honba)
```

### 4.2 Physical Point Stick Exchange

**Proper procedure**:
1. Gently place point sticks in **front corner** of recipient's area
2. Do NOT throw or toss sticks
3. Avoid aggressive gestures that show frustration

**Reason**: Sportsmanship and respect for opponents, regardless of game outcome.

**Digital equivalent**: Point transfer animations should be neutral and respectful (no mocking or excessive celebration).

## 5. Game End and Scorekeeping

### 5.1 Score Confirmation Obligations

**General rule**: If playing on automatic table (shows all scores), less need for manual confirmation.

**If NOT on automatic table**: Any player can request score confirmation, and others must comply.

**Reason**: Transparency and fairness in score tracking.

**Critical timing**: Before final hand (オーラス/오라스), scores MUST be confirmed clearly.

**What to check**:
- Total points for each player
- Honba count
- Riichi stick deposits on table
- Calculation accuracy

**Digital implementation**:
- Always display live scores (no hidden information)
- Highlight score changes with clear +/- indicators
- Confirmation dialog before final hand: "Current scores: [list]. Begin final hand?"

### 5.2 Pre-Final-Hand Verification

**Extra diligence needed**:
- Verify total points add up to 100,000 (4 players × 25,000 starting points)
- Check for scoring errors accumulated during game
- Confirm honba and deposit counts
- Even on automatic tables, double-check display matches reality

**Reason**: Final hand outcomes determine placement bonuses (uma). Scoring errors here are most costly.

**Example error scenario**:
```
Automatic table display: Player A = 28,000
Physical sticks: Player A actually has 27,000 (miscounted earlier)
→ Affects final placement if score difference is close
```

**Digital implementation**:
- Periodic validation: Total points = 100,000 (flag discrepancies immediately)
- Transaction log: All point transfers recorded for audit
- Score dispute resolution: Replay hand history with point calculations

## 6. Digital Implementation: Etiquette-Informed Design

### 6.1 Timing and Fairness

**Principle**: Design timers that are fair but prevent stalling.

**Recommended timing**:
- **Discard decision**: 5-10 seconds (normal turn)
- **Calling decision (chi/pon/kan)**: 2-3 seconds (quick reflexes expected)
- **Riichi declaration**: 10-15 seconds (requires hand verification)
- **Win declaration**: 15-20 seconds (score calculation needed)

**Progressive penalties**:
- First timeout: Warning + auto-discard (random tile)
- Repeated timeouts: AFK flag + game forfeit (after 3 consecutive)

### 6.2 Communication Features

**Allowed**:
- Preset phrases: "Hello", "Good game", "Sorry", "Thank you"
- Emote reactions: Celebrate win, acknowledge good play
- Post-game chat: Open discussion after hand ends

**Restricted during active gameplay**:
- Free-form text chat (enable only between hands)
- Voice chat (too distracting, use push-to-talk if allowed)
- Emoji spam (rate limit: max 3 per minute)

### 6.3 Sportsmanship Enforcement

**Report system**:
- Abusive language
- Intentional slow play (stalling)
- AFK/disconnection abuse
- Collusion with other players

**Penalties**:
- Warning (first offense)
- Temporary chat restriction (repeat offenses)
- Temporary matchmaking ban (severe violations)
- Permanent ban (extreme cases: cheating, harassment)

### 6.4 Visual Clarity and Accessibility

**Design principles from physical etiquette**:

| Physical Etiquette | Digital Equivalent |
|--------------------|-------------------|
| 6 tiles per row discard layout | Grid display (6 per row) in discard area |
| Clear voice declarations | Prominent "TSUMO/RON" button with confirmation |
| Dealer marker (親マーク) visible | Persistent dealer indicator on UI |
| Point sticks placed in front | Point transfer animation to recipient's area |
| Revealed melds on right side | Called melds displayed separately from hand |
| Riichi discard sideways | Riichi discard highlighted/rotated in UI |

## 7. Etiquette Violations and Enforcement

### 7.1 Minor Violations (Warnings)

**Examples**:
- Forgetting to announce tsumo/ron before revealing
- Accidentally touching another player's tiles (physical)
- Excessive celebration/gloating
- Asking for score confirmation too frequently

**Penalty**: Verbal warning, no point loss.

### 7.2 Serious Violations (Penalties)

**Examples**:
- Intentional slow play to stall game
- Revealing tiles before win declaration
- Strategic table talk to mislead opponents
- Checking wall tiles after hand ends

**Penalty**: May result in **penalty points** (see constitution_penalties.md 3.2) or **win prohibition** depending on severity.

### 7.3 Critical Violations (Chombo)

**Examples**:
- Looking at opponent's concealed tiles
- Exposing opponent's tiles by impact/force
- Intentional cheating (swap tiles, false declarations)

**Penalty**: **Chombo** (see constitution_penalties.md 5.4.5, 5.4.6).

## 8. Cultural Context and Sportsmanship

### 8.1 Mahjong as Social Game

**Philosophy**: Mahjong is traditionally a social game emphasizing:
- **Respect** (尊重/존중): Honor opponents regardless of skill level
- **Fairness** (公平/공평): Follow rules consistently
- **Grace** (優雅/우아): Win humbly, lose gracefully
- **Enjoyment** (楽しみ/즐거움): Prioritize fun over winning

**Practical application**:
- Don't berate opponents for mistakes
- Congratulate winners sincerely
- Accept losses without excuses
- Help beginners learn rules

### 8.2 Balancing Competition and Courtesy

**Competitive play is encouraged**, but within etiquette bounds:

**Acceptable competitive behavior**:
- Analyzing discards to deduce opponent hands (reading)
- Strategic tile selection to maximize winning chances
- Defensive play (betaori/베타오리) when opponents are threatening

**Unacceptable behavior**:
- Mocking opponents for poor play
- Excessive celebration that humiliates losers
- Deliberately slowing game when losing to frustrate winners
- Verbal/physical intimidation

## 9. Implementation Testing

### 9.1 UI/UX Validation

**Checklist for digital implementation**:
- [ ] Discard area displays 6 tiles per row in chronological order
- [ ] Dealer indicator prominently visible throughout hand
- [ ] Point transfer animations clear and respectful (no mocking tone)
- [ ] Tsumo/Ron buttons require confirmation before revealing hand
- [ ] Called melds (chi/pon/kan) displayed separately from hand
- [ ] Riichi discard visually distinct (highlighted/rotated)
- [ ] Timer warnings give sufficient notice before auto-action
- [ ] Score breakdown clear (base points + honba + deposits)

### 9.2 Etiquette Enforcement Testing

**Test scenarios**:
```python
def test_late_calling_prevented():
    """Verify calling disabled after next player draws."""
    game = MahjongGame()
    player_a.discard("3m")
    time.sleep(CALLING_WINDOW_SECONDS + 0.1)

    player_b.draw_tile()  # Next player starts turn
    result = player_c.call_pon("3m")  # Try to call on previous discard

    assert result == "CALLING_WINDOW_CLOSED"

def test_win_declaration_before_reveal():
    """Verify hand reveal only after explicit tsumo/ron declaration."""
    player = create_player(winning_hand)
    player.draw_tile(winning_tile)

    # Attempt to reveal without declaration
    result = player.reveal_hand()
    assert result == "MUST_DECLARE_WIN_FIRST"

    # Proper sequence
    player.declare_tsumo()
    result = player.reveal_hand()
    assert result == "SUCCESS"
```

## 10. Summary: Etiquette Principles for Implementation

1. **Clarity**: Clear turn indicators, timers, and action availability
2. **Fairness**: Consistent timing, no information leakage, anti-cheat measures
3. **Respect**: Neutral animations, no mocking, sportsmanship enforcement
4. **Accessibility**: Visual clarity inspired by physical layout conventions
5. **Transparency**: Open scorekeeping, audit logs, dispute resolution
6. **Enjoyment**: Balance competitive play with courtesy and fun

**Goal**: Digital implementation should capture the **spirit of physical mahjong etiquette** while leveraging digital advantages (automated scoring, timing enforcement, cheat prevention).

## References

- Korean Mahjong guide: "한 권으로 익히는 리치 마작0426.pdf" (YOSTAR/Korean Mahjong League, pages 42-43)
- Related constitution documents:
  - `constitution_majhong_rule.md` - Core game rules
  - `constitution_penalties.md` - Violations and penalties
  - `constitution_term.md` - Terminology reference
- Traditional mahjong etiquette: Japanese Mahjong Federation guidelines
