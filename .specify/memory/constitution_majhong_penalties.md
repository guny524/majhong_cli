# Mahjong Penalty and Violation Rules Reference

**Purpose**: Defines penalty system, chombo rules, and violation handling for Riichi Mahjong implementation.
**Referenced by**: `.specify/memory/constitution.md`
**Last Updated**: 2025-11-18
**Source**: "한 권으로 익히는 리치 마작0426.pdf" (YOSTAR/Korean Mahjong League official guide, pages 40-41)

## 1. Penalty Application Period

**When penalties apply** (벌칙의 적용):
- **Start point**: After dealer's first discard (친의 첫 타패 이후)
- **End point**: Until win (화료) or exhaustive draw (유국)
- **Conflict resolution**: Win takes precedence over penalties if they occur simultaneously

**Implementation note**: Track game phase to determine if penalty rules are active.

## 2. Penalty Types Overview

Mahjong violations are classified into three severity levels:

| Type | Severity | Impact Scope | Example |
|------|----------|--------------|---------|
| **Penalty Points** (벌칙 점수) | Minor | No major disruption | Incorrect call declaration |
| **Win Prohibition** (화료 불가) | Medium | Individual player only | Wrong number of tiles |
| **Chombo** (촌보) | Critical | Entire game/hand | No-tenpai riichi at exhaustive draw |

## 3. Penalty Points (벌칙 점수)

### 3.1 Definition

**벌칙 점수** (Penalty points) applies to minor mistakes that do not significantly affect game flow.

**Common scenario**: Incorrect call declaration (발성을 잘못 한 경우)
- Example: Saying "pon" but intending "chi"
- Note: Verbal error alone (not taking tiles) may not trigger penalty

### 3.2 Penalty Amount

**1,000 points deposited** by the offending player:
- Points placed in the center as deposit (公托/공탁)
- Winner of the current or next hand collects all deposits
- Similar to riichi deposit mechanics

### 3.3 Insufficient Points Handling (Bust Rule - 들통)

**If player has <1,000 points when penalty occurs**:
1. Borrow 1,000 points from other players to deposit
2. Continue playing the hand
3. At end of hand:
   - If still negative → bust (들통) → game ends
   - If penalty recovered → continue normally

**Implementation consideration**: Only applies if bust/negative points rule is enabled in game configuration.

## 4. Win Prohibition (화료 불가)

### 4.1 Definition

**화료 불가** (Win prohibition / No-win status) occurs when a player makes a mistake affecting only themselves during gameplay.

**Consequences**:
- Treated as **no-tenpai** (노텐) immediately
- **Cannot tsumo or ron**
- **Cannot call chi/pon/kan**
- Violating these restrictions → escalates to **chombo**

### 4.2 Riichi with Win Prohibition

**If win prohibition occurs after riichi declaration**:
- Still treated as no-tenpai
- At exhaustive draw → revealed as no-tenpai riichi → **chombo**

### 4.3 Win Prohibition Scenarios

#### (1) Incorrect Chi/Pon/Kan (잘못된 울기)

**When it applies**:
- After declaring call and revealing hand tiles
- Taking discard tile from river
- Confirmed as invalid meld formation

**When it does NOT apply**:
- Verbal error only (말실수만 한 경우) - corrected before revealing tiles
- Example: Saying "pon" but correcting to "chi" before taking tile

**Validation logic**:
```python
def validate_meld_formation(called_tile, hand_tiles, meld_type):
    """Validate if meld can be legally formed."""
    if meld_type == "CHI":
        # Must form sequence in same suit
        if not forms_valid_sequence(called_tile, hand_tiles):
            return "WIN_PROHIBITION"
    elif meld_type == "PON":
        # Must have 2 matching tiles
        if hand_tiles.count(called_tile) < 2:
            return "WIN_PROHIBITION"
    return "VALID"
```

#### (2) Wrong Tile Count (소패/다패)

**소패 (Shortage / 少牌)**: Fewer tiles than required
**다패 (Excess / 多牌)**: More tiles than required

**Standard tile count** (after discarding):
- **Base**: 13 tiles
- **Per kan**: +1 tile (14 tiles with one kan, 15 with two kans, etc.)

**Common causes**:
- 소패: Forgetting to draw tile on your turn
- 다패: Accidentally drawing tile twice

**Detection timing**:
- Checked when player attempts to declare tsumo/ron
- Checked during furiten verification
- May be noticed by other players during gameplay

**Implementation**:
```python
def validate_tile_count(hand_tiles, open_melds):
    """Verify player has correct number of tiles."""
    kan_count = sum(1 for meld in open_melds if meld.type == "KAN")
    expected_count = 13 + kan_count

    if len(hand_tiles) != expected_count:
        if len(hand_tiles) < expected_count:
            return "SHORTAGE"  # 소패
        else:
            return "EXCESS"    # 다패
    return "VALID"
```

#### (3) Swap Cheating (울어바꿈)

**Types of swap cheating**:
- **스지 울어바꿈** (Suji swap): Swapping tiles related to wait patterns
- **현물 울어바꿈** (Genbutsu swap): Swapping safe tiles

**What this means**: Illegally exchanging tiles between hand and called melds after calling chi/pon/kan.

**Implementation note**: This is difficult to detect in digital implementations but important for physical mahjong rules.

#### (4) Incorrect Tsumo/Ron Declaration (쯔모, 론의 선언을 잘못한 경우)

**When it applies**:
- Declaring tsumo or ron
- But **NOT revealing tiles** yet
- Example: Player says "tsumo" then realizes mistake and retracts before showing hand

**Key distinction**: If tiles are revealed → escalates to chombo (see 5.4.3, 5.4.4)

#### (5) Incorrect Riichi Declaration (리치의 선언을 잘못한 경우)

**When it applies**:
- Declaring riichi by mistake
- Attempting to cancel riichi before depositing stick or discarding sideways

**Consequences**:
- Win prohibition for current hand
- Cannot riichi again in same hand

## 5. Chombo (촌보)

### 5.1 Definition

**촌보** (Chombo / 錯和 / cuòhé) occurs when a player makes a critical mistake that makes it impossible to continue the hand normally, affecting all players.

**Etymology**:
- Korean: 촌보 (phonetic from Japanese)
- Japanese: チョンボ (Chombo)
- Chinese: 錯和 (cuòhé) - "incorrect win"

### 5.2 Penalty Amount

**Payment structure** (만관 상당 / Mangan equivalent):

**Traditional payment**:
- **Dealer (親/oya)**: Pays 4,000 points to each player (12,000 total)
- **Non-dealer (子/ko)**: Pays 2,000 to non-dealers, 4,000 to dealer (8,000 total)

**Modern equal payment** (공평한 방식):
- All offenders pay **3,000 points to each player** (9,000 total)
- Increasingly common in modern rulesets

**Implementation configuration**:
```yaml
chombo_payment_style: "traditional"  # traditional | modern_equal
chombo_traditional:
  dealer: [4000, 4000, 4000]      # pays to each player
  non_dealer_to_dealer: 4000
  non_dealer_to_others: 2000
chombo_modern:
  all_players: 3000                # equal payment to all
```

### 5.3 Hand Reset Procedure

**When chombo occurs**:
1. Current hand becomes void (무효)
2. Hand restarted from beginning
3. **Honba count does NOT increase** (연장 횟수 쌓이지 않음)
4. **All deposits returned to original owners**:
   - Riichi sticks (리치 선언봉)
   - Penalty point deposits (벌칙 점수 공탁)
5. Dealer position does NOT change (dealer continues)

**Implementation note**: Save initial game state before hand starts to enable clean reset.

### 5.4 Chombo Scenarios

#### (1) Exposing Wall Tiles (패산을 무너트려 패가 보인 경우)

**Threshold**: Cumulatively 5+ tiles from wall revealed
- Knocking over wall accidentally
- Exposing dead wall tiles
- Revealing opponent's draw tiles

**Digital implementation**: Less relevant for online/CLI games, but important for physical mahjong simulation.

#### (2) No-Tenpai Riichi at Exhaustive Draw (노텐 리치를 하고 유국이 되었을 경우)

**Critical rule**: Riichi must be declared with valid tenpai hand

**When chombo occurs**:
- Player declares riichi (no-tenpai or invalid wait)
- Hand proceeds to exhaustive draw (유국)
- Riichi declarer must reveal hand
- Hand revealed as not tenpai → chombo

**Prevention validation**:
```python
def validate_riichi_declaration(hand, riichi_tile):
    """Verify riichi is declared with valid tenpai."""
    # Check if hand is 1-away from winning
    if not is_tenpai(hand):
        return "INVALID_RIICHI"

    # Check if waiting tiles exist in remaining wall
    wait_tiles = calculate_wait_tiles(hand)
    if not any(tile in remaining_tiles for tile in wait_tiles):
        return "INVALID_WAIT"  # No winnable tiles remain

    return "VALID"
```

**Note**: If another player wins before exhaustive draw, no-tenpai riichi is not exposed → no penalty.

#### (3) Winning Without Yaku (역이 없는데 화료를 선언하고 손패를 열어 공개한 경우)

**Critical requirement**: Must have at least 1 yaku (excluding dora) to win

**When chombo occurs**:
- Player declares tsumo or ron
- **Reveals hand tiles**
- Hand has no valid yaku (only dora)

**Validation code**:
```python
def validate_winning_hand(hand, win_tile, game_state):
    """Verify hand has minimum yaku requirement before allowing win."""
    yaku_list = detect_yaku(hand, win_tile, game_state)

    # Filter out dora (not counted as yaku)
    actual_yaku = [y for y in yaku_list if y.type != "DORA"]

    if len(actual_yaku) == 0:
        if hand_already_revealed:
            return "CHOMBO_NO_YAKU"
        else:
            return "WIN_PROHIBITION"  # If not yet revealed

    return "VALID"
```

**Key distinction**: If caught BEFORE revealing tiles → win prohibition, not chombo.

#### (4) Furiten Ron (후리텐 상황에서 론 화료를 선언하고 손패를 열어 공개한 경우)

**Furiten rule**: Cannot ron when in furiten status (see constitution_majhong_rule.md)

**When chombo occurs**:
- Player is in furiten (own discard, temporary, or riichi furiten)
- Declares ron (not tsumo - tsumo is allowed in furiten)
- **Reveals hand tiles**

**Validation**:
```python
def validate_ron_declaration(player, win_tile):
    """Check if ron is allowed (not in furiten)."""
    furiten_status = check_furiten_status(player)

    if furiten_status.is_furiten:
        if hand_already_revealed:
            return "CHOMBO_FURITEN_RON"
        else:
            return "WIN_PROHIBITION"

    return "VALID"
```

**Three types of furiten** (see detailed explanation in constitution_majhong_rule.md):
1. Discard furiten (자신이 버린 패)
2. Temporary furiten (동순 후리텐)
3. Riichi furiten (리치 후 넘긴 패)

#### (5) Exposing Opponent's Hand (다른 사람의 손패를 건드리거나 공개해버린 경우)

**Two scenarios**:
- Touching another player's hand tiles (건드리거나)
- Exposing them by impact/shock (충격으로 공개)

**Intent irrelevant**: Accidental exposure still results in chombo

**Digital implementation**: Not applicable to online/CLI games, but document for completeness.

#### (6) Looking at Opponent's Hand (다른 사람의 손패를 들여다보는 경우)

**Intentional peeking** at another player's concealed tiles

**Digital implementation**: Prevented by design in online games, but relevant for:
- Replay validation (ensuring no information leakage)
- Spectator mode restrictions
- Anti-cheat mechanisms

#### (7) Illegal Kan After Riichi (리치 이후에 허용되지 않는 깡을 한 경우)

**Riichi restriction**: After riichi, kan is only allowed for **pure, independent ankans** (순수하게 독립된 안커)

**Rule**: "패의 구성이 변하는 깡" (kan that changes wait pattern) is forbidden

**Definition of pure independent ankan**:
- Quad is completely separate from waiting tiles
- Declaring kan does NOT change:
  - Wait tiles
  - Wait patterns
  - Number of winning combinations

**Example 1 - ALLOWED**:
```
Hand: 東東東 3455678999 man
Wait: 2-4-5-7-8 man (multi-sided wait)
Riichi declared.
Draw: 東 (East wind)

Kan 東東東東 → ALLOWED
Reason: East quad is independent of 3455678999 man wait structure
```

**Example 2 - FORBIDDEN**:
```
Hand: 東東東 3455678999 man
Wait: 2-4-5-7-8 man
Riichi declared.
Draw: 9 man

Kan 9999 man → FORBIDDEN (Chombo if declared)
Reason: Making 9999 man quad eliminates the 89 man wait (7 man edge wait disappears)
Original wait: 2m/4m/5m/7m/8m (includes 89m edge → 7m)
After kan: 2m/4m/5m/8m only (89m pattern destroyed)
```

**Validation algorithm**:
```python
def validate_riichi_kan(hand, kan_tiles, original_wait):
    """Check if kan is allowed after riichi."""
    # Simulate hand after kan
    simulated_hand = hand.copy()
    simulated_hand.remove_tiles(kan_tiles)
    simulated_hand.add_open_meld(Kan(kan_tiles))

    # Calculate new wait tiles and patterns
    new_wait = calculate_wait_tiles(simulated_hand)

    # Compare wait patterns
    if original_wait.tiles != new_wait.tiles:
        return "FORBIDDEN_KAN_CHANGES_WAIT"

    if original_wait.patterns != new_wait.patterns:
        return "FORBIDDEN_KAN_CHANGES_PATTERNS"

    return "ALLOWED"
```

**Key insight**: Even if the wait *tiles* remain the same, if the wait *patterns* change (e.g., losing a edge wait option), the kan is forbidden.

## 6. Penalty Severity Decision Tree

```
Mistake occurs
    ↓
Is it during active play period? (after dealer's first discard)
    ↓ YES
Does it affect all players / make hand unplayable?
    ↓ YES → CHOMBO (촌보)
        - Pay mangan equivalent to all players
        - Hand restarted, deposits returned
    ↓ NO
Does it prevent player from winning?
    ↓ YES → WIN PROHIBITION (화료 불가)
        - Treated as no-tenpai
        - Cannot tsumo/ron/call
        - If violated → escalates to chombo
    ↓ NO
Is it a minor procedural error?
    ↓ YES → PENALTY POINTS (벌칙 점수)
        - 1,000 point deposit
        - Game continues normally
    ↓ NO
    → NO PENALTY (continue game)
```

## 7. Implementation Guidelines

### 7.1 Validation Checkpoints

**Before allowing player actions**:
```python
# Before riichi declaration
validate_riichi_tenpai(hand)
validate_sufficient_points(player, 1000)
validate_remaining_wall_tiles(wall)

# Before declaring win
validate_tile_count(hand, open_melds)
validate_yaku_exists(hand, win_tile)
validate_not_furiten(player, win_tile)

# Before calling chi/pon/kan
validate_meld_formation(called_tile, hand_tiles, meld_type)

# Before kan after riichi
validate_riichi_kan(hand, kan_tiles, original_wait)
```

### 7.2 Error Handling

**When validation fails**:
1. Determine penalty type (penalty points / win prohibition / chombo)
2. Apply penalty immediately
3. Log violation for replay/audit
4. Update game state accordingly

**Example error handler**:
```python
def handle_invalid_action(player, action, violation_type):
    """Process penalty for invalid action."""
    if violation_type == "CHOMBO":
        # Critical violation
        apply_chombo_penalty(player)
        reset_hand(restore_deposits=True)
        log_violation(player, action, "CHOMBO")

    elif violation_type == "WIN_PROHIBITION":
        # Medium violation
        player.status = "NO_WIN"
        player.is_tenpai = False
        log_violation(player, action, "WIN_PROHIBITION")

    elif violation_type == "PENALTY_POINTS":
        # Minor violation
        player.points -= 1000
        game.deposits += 1000
        log_violation(player, action, "PENALTY_POINTS")
```

### 7.3 Testing Requirements

**Unit tests for each penalty scenario**:
```python
def test_no_tenpai_riichi_becomes_chombo():
    """Verify no-tenpai riichi at exhaustive draw triggers chombo."""
    game = MahjongGame()
    player = game.players[0]

    # Declare riichi with invalid hand (not tenpai)
    player.declare_riichi(hand=not_tenpai_hand)

    # Proceed to exhaustive draw
    game.play_until_exhaustive_draw()

    # Verify chombo applied
    assert game.last_penalty_type == "CHOMBO"
    assert player.points == initial_points - chombo_payment
    assert game.honba == initial_honba  # No increase

def test_wrong_tile_count_win_prohibition():
    """Verify incorrect tile count triggers win prohibition."""
    player = create_player_with_tiles([...])  # 12 tiles only

    draw_result = player.draw_tile(tile)
    # Player now has 13 tiles, attempts to discard and declare tsumo

    result = validate_winning_hand(player)
    assert result == "WIN_PROHIBITION"
    assert player.can_tsumo == False
    assert player.can_ron == False

def test_illegal_kan_after_riichi():
    """Verify kan that changes wait pattern triggers chombo."""
    hand = "東東東3455678999m"
    player.declare_riichi(wait=[2,4,5,7,8])

    player.draw_tile("9m")
    result = player.declare_kan("9999m")

    assert result == "CHOMBO"
    assert game.state == "HAND_RESET"
```

## 8. Terminology Mapping

| English | Korean | Japanese | Chinese | Code |
|---------|--------|----------|---------|------|
| Penalty period | 벌칙 적용 시점 | 罰則の適用 | 罚则适用期 | `PENALTY_PERIOD` |
| Penalty points | 벌칙 점수 | 罰符 | 罚分 | `PENALTY_POINTS` |
| Win prohibition | 화료 불가 | 和了不可 | 不可和 | `WIN_PROHIBITION` |
| Chombo | 촌보 | チョンボ | 錯和 | `CHOMBO` |
| No-tenpai | 노텐 | 不聴 | 不听 | `NOT_TENPAI` |
| Tile shortage | 소패 | 少牌 | 少牌 | `TILE_SHORTAGE` |
| Tile excess | 다패 | 多牌 | 多牌 | `TILE_EXCESS` |
| Swap cheating | 울어바꿈 | 食い替え | 吃换 | `ILLEGAL_SWAP` |
| Bust (negative points) | 들통 | トビ | 飞 | `BUST` |
| Mangan payment | 만관 지불 | 満貫払い | 满贯支付 | `MANGAN_PAYMENT` |

## References

- Korean Mahjong guide: "한 권으로 익히는 리치 마작0426.pdf" (YOSTAR/Korean Mahjong League, pages 40-41)
- Related constitution documents:
  - `constitution_majhong_rule.md` - Core game rules and furiten
  - `constitution_term.md` - Terminology reference
  - `constitution_yaku.md` - Yaku requirements for win validation
