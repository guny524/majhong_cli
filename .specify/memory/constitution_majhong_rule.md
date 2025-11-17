# Riichi Mahjong Game Rules Reference

**Purpose**: Defines official Riichi Mahjong rules for implementation in majhong_cli project.
**Referenced by**: `.specify/memory/constitution.md`
**Last Updated**: 2025-11-17

## Tile Composition (136 tiles)

### Number Tiles (108 tiles)
- **Man (萬/만, characters)**: 1m~9m × 4 each = 36 tiles
- **Pin (筒/통, circles)**: 1p~9p × 4 each = 36 tiles
- **Sou (索/삭, bamboo)**: 1s~9s × 4 each = 36 tiles

### Honor Tiles (28 tiles)
- **Wind tiles (風牌)**: East/South/West/North (동/남/서/북, 1z~4z) × 4 each = 16 tiles
- **Dragon tiles (三元牌)**: White/Green/Red (백/발/중, 5z~7z) × 4 each = 12 tiles

### Red Five Tiles (Optional)
- **Aka-dora (赤ドラ)**: 0m, 0p, 0s (red five variants, used in some rulesets)

## Tile Notation Standard

**For majhong_cli implementation**:
```
1m 2m 3m 4m 5m 6m 7m 8m 9m  # Man/characters
1p 2p 3p 4p 5p 6p 7p 8p 9p  # Pin/circles
1s 2s 3s 4s 5s 6s 7s 8s 9s  # Sou/bamboo
1z 2z 3z 4z                 # East/South/West/North
5z 6z 7z                    # White/Green/Red
0m 0p 0s                    # Red fives (optional)
```

## Game Setup

### 1. Seat Determination

**Procedure**:
1. Each player draws one wind tile (East/South/West/North)
2. Player who draws East chooses their seat position
3. Other players sit counter-clockwise: South/West/North relative to East
4. East player rolls dice to determine "temporary dealer" (假東/가동)
   - Count counter-clockwise from East (including East = 1)
   - Position matching dice sum becomes temporary dealer
5. Temporary dealer rolls again to determine "actual dealer" (眞東/진동, Oya/親)
   - Count counter-clockwise from temporary dealer (including self = 1)
   - This player receives the dealer marker (親マーク/선 마커)

### 2. Initial Point Sticks (25,000 points)

**Standard distribution**:
- Red (10,000): ×1
- Yellow (5,000): ×2
- Blue (1,000): ×4
- Green (500): ×1 (often substituted with marked 100-point stick)
- White (100): ×5

**Note**: 500-point sticks often don't exist in physical sets. Mark a 100-point stick with "500" for practical use.

### 3. Wall Construction

See `constitution_majhong_wall.md` for detailed wall building procedures.

**Summary**:
- Each player builds 17 tiles × 2 layers = 34 tiles wall segment
- Dealer rolls dice to determine wall break point
- Separate 14 tiles as dead wall (王牌/wanpai)
- Remaining 122 tiles form live wall (活牌)

### 4. Initial Hand Distribution (Haipai/配牌)

**Distribution sequence**:
1. Starting from dealer, distribute 4 tiles per player × 3 rounds = 12 tiles each
2. Dealer takes 2 additional tiles (14 total)
3. Other players take 1 additional tile each (13 total)

**Direction**: Tiles distributed counter-clockwise (dealer → South → West → North)
**Drawing from wall**: Clockwise direction from the break point

## Wall Structure

### Dead Wall (王牌/Wanpai) - 14 tiles

**Composition** (from left to right):
1. **Rinshan tiles (嶺上牌/영상패)**: First 4 tiles (leftmost)
   - Drawn when declaring kan (槓)
   - Winning with rinshan draw = "Rinshan Kaihou" (嶺上開花/영상개화) yaku
2. **Dora indicators (ドラ表示牌/도라표시패)**: Next 5 tiles
   - 3rd tile from left is flipped face-up at game start
   - Additional indicators flip upon kan declarations (max 5 dora indicators)
3. **Ura-dora indicators (裏ドラ表示牌/뒷도라표시패)**: Last 5 tiles (rightmost)
   - Only revealed after successful riichi and winning
   - Under dora indicators (in the bottom layer)

**Dead wall maintenance**:
- Always maintains exactly 14 tiles
- When rinshan tile is drawn (during kan), append last tile from live wall to maintain count

### Live Wall (活牌) - 122 tiles

**Drawing order**: Clockwise from the wall break point
**Special tiles**:
- **Haitei tile (海底牌/해저패)**: Last drawable tile from live wall
  - Self-draw win = "Haitei Raoyue" (海底撈月/해저로월) yaku
- **Houtei tile (河底牌/하저패)**: Last discarded tile
  - Ron win = "Houtei Raoyui" (河底撈魚/하저로어) yaku

## Turn Structure

### Basic Turn Flow

1. **Tsumo (自摸/쯔모)**: Draw one tile from live wall (dealer skips first draw)
2. **Discard**: Discard one tile to the center discard area
3. **Calling window**: Other players may call chi/pon/kan/ron before next player's turn
4. **Turn advance**: Next player (counter-clockwise) begins their turn

### Hand Organization

**Physical placement**:
- Drawn tile (tsumo) kept separate (rightmost) until decision made
- Distinguishes between hand tiles and just-drawn tile (affects scoring in some situations)
- Called tiles (chi/pon/kan) placed horizontally to the right of hand, face-up

### Discard Rules

**Visible discards**:
- Place discarded tiles face-up in personal discard area (河/kawa)
- Arrange in chronological order (6 tiles per row is standard)
- **Furiten rule**: Cannot declare ron on tiles you previously discarded

**Riichi discard**:
- When declaring riichi, place discard tile sideways (90° rotation)
- Indicates riichi declaration point

## Calling Actions (鳴き/Naki)

### Chi (チー/치) - Sequential Meld

**Requirements**:
- Only from the immediately previous player's discard (player to your left)
- Form a 3-tile sequence (e.g., 3p 4p 5p)
- Must be within same suit (man/pin/sou only, not honors)

**Effect**: Hand becomes open (門前/menten status lost)

### Pon (ポン/퐁) - Triplet Meld

**Requirements**:
- From any player's discard
- Requires 2 matching tiles in hand

**Effect**: Hand becomes open, priority over chi

### Kan (槓/깡) - Quad Meld

**Four types**:

1. **Daiminkan (大明槓/대명깡)**: Open kan from another's discard
   - Requires 3 matching tiles in hand
   - Draw rinshan tile immediately
   - Hand becomes open

2. **Shouminkan (小明槓/소명깡)**: Added kan
   - Add 4th tile to previously called pon
   - Draw rinshan tile immediately
   - Can be robbed (槍槓/chankan) for specific yaku

3. **Ankan (暗槓/암깡)**: Concealed kan
   - All 4 tiles in hand, declared during your turn
   - Tiles placed face-down (concealed status maintained)
   - Draw rinshan tile immediately
   - Hand remains closed for menten yaku eligibility

**Kan effects**:
- Flip next dora indicator (max 5 total)
- Draw rinshan tile from dead wall
- Replenish dead wall from live wall's end

### Ron (ロン/론) - Win on Discard

**Requirements**:
- Complete winning hand with another player's discard
- Must have at least 1 yaku (excluding dora)
- Not in furiten status

**Priority**: Ron has priority over chi/pon (except for double/triple ron)

## Riichi Declaration (立直/リーチ)

### Requirements

1. **Tenpai (聴牌/텐파이)**: Hand is 1 tile away from winning
2. **Closed hand**: No called melds (chi/pon/daiminkan/shouminkan)
3. **Point requirement**: Must have 1,000+ points for riichi deposit
4. **Remaining tiles**: At least 4 tiles left in live wall

### Declaration Procedure

1. Announce "riichi" before discarding
2. Place riichi stick (1,000 points) on table
3. Discard tile sideways (90° rotation) to mark riichi point
4. Cannot change hand composition after riichi (except for ankan in specific situations)

### Riichi Effects

**Advantages**:
- +1 han (riichi yaku)
- +1 han if win on first self-draw after riichi (一発/ippatsu)
- Ura-dora revealed upon winning
- Riichi stick deposits go to winner

**Restrictions**:
- Cannot change wait tiles (furiten applies to all waiting tiles)
- Must declare tsumo/ron if winning tile appears
- Hand is "autopilot" - minimal decision making

## Winning Conditions (和了/Hōra/화료)

### Hand Structure Requirements

**Standard winning hand** (14 tiles):
- 4 melds (面子/mentsu) + 1 pair (対子/toitsu)
- **Meld**: Either sequence (順子/shuntsu) or triplet/quad (刻子/koutsu)
- **Pair**: 2 identical tiles (雀頭/jantou)

**Special patterns**: Some yaku override standard structure (e.g., Seven Pairs, Thirteen Orphans)

### Minimum Yaku Requirement

**Critical rule**: Must have at least 1 yaku (役), **excluding dora**
- Dora alone CANNOT enable winning
- Dora adds han but doesn't count as yaku

### Victory Types

1. **Tsumo (ツモ/쯔모)**: Self-draw from wall
   - All players pay dealer/non-dealer rates

2. **Ron (ロン/론)**: Win from discard
   - Discarding player pays full amount

### Furiten (振聴/후리텐) - Forbidden Win Status

**Definition**: A state where a player in tenpai is prohibited from declaring ron (but can still win by tsumo).

**Critical rule**: Furiten applies to **ALL waiting tiles** simultaneously. If even one wait tile is in furiten, the player cannot ron on ANY of their wait tiles.

**Example of multi-wait furiten**:
```
Hand: 2-5-8 man wait (multi-sided wait)
If 5 man is in your discard: FURITEN on all waits (2m, 5m, 8m)
Cannot ron on 2m or 8m either, even though they're not in discards
```

**Example of shanpon (double pon) wait**:
```
Hand: 東東 + 南南 (waiting for 東 or 南 to complete)
If 東 is in your discard: FURITEN applies
Cannot ron on 南 either
```

#### Three Types of Furiten

**1. Discard Furiten (自分의 버림패로 후리텐)**

**Rule**: Cannot ron on any tile that matches your own discards.

**Duration**: Permanent for the entire hand (해당 국이 끝날 때까지)

**Applies to**: All waiting tiles if any wait tile is in your discard pile

**Example**:
```
Early game: Discard 3 pin
Later: Reach tenpai waiting on 3-6 pin
Result: FURITEN - cannot ron on 3p or 6p for rest of hand
Can only win by tsumo
```

**Implementation**:
```python
def check_discard_furiten(player, wait_tiles):
    """Check if any wait tile is in player's discard pile."""
    for wait_tile in wait_tiles:
        if wait_tile in player.discards:
            return True  # Furiten on ALL waits
    return False
```

**2. Temporary Furiten (同巡後리텐 / 同順後리텐)**

**Rule**: If you skip a winning tile (don't call ron), you cannot ron until your next turn.

**Duration**: Until your next draw (내 차례가 다시 돌아오면 해소)

**Condition**: Only applies when NOT in riichi

**Example**:
```
Player position: East (your turn order: E → S → W → N → E)
Current player: South discards your wait tile (3 man)
You: DON'T call ron (choose to skip)
Result: TEMPORARY FURITEN
- West discards another wait tile (6 man) → Cannot ron
- North discards another wait tile → Cannot ron
- Your turn comes (draw tile) → Furiten CLEARED

Now if South/West/North discard wait tiles again, you CAN ron
```

**Why this rule exists**: Prevents selective ron to manipulate payment (choosing to ron from specific players only).

**Implementation**:
```python
def check_temporary_furiten(player, game_state):
    """Check if player skipped ron in current round."""
    if player.is_riichi:
        return False  # Riichi has different furiten rules

    current_round = game_state.current_round_number

    if player.last_skipped_ron_round == current_round:
        return True  # Still in same round since skip

    return False

def on_player_turn_start(player):
    """Clear temporary furiten when player's turn comes."""
    player.temporary_furiten = False
```

**3. Riichi Furiten (리치 후 넘긴 패로 후리텐)**

**Rule**: After declaring riichi, if you skip ANY winning tile (tsumo or ron), permanent furiten for the rest of the hand.

**Duration**: Permanent until hand ends (해당 국이 끝날 때까지)

**Critical difference from temporary furiten**:
- Applies even to **tsumo tiles you skip**
- NOT cleared on your next turn (unlike temporary furiten)

**Example**:
```
Declare riichi waiting on 3-6 man
Turn 1: Draw 3 man (winning tile) → Don't declare tsumo → PERMANENT FURITEN
Turn 2: Opponent discards 6 man → Cannot ron (furiten)
Turn 3+: Can only tsumo (if you didn't skip any more winning draws)
```

**Special rule**: You can still declare furiten riichi (후리텐 리치), but only tsumo is possible.

**Why this rule exists**:
- Riichi declares your hand publicly → you commit to winning immediately
- Skipping a win suggests strategic waiting → violates riichi spirit
- Penalty: lose ron option, tsumo-only

**Implementation**:
```python
def check_riichi_furiten(player):
    """Check if player skipped winning tile after riichi."""
    if not player.is_riichi:
        return False

    # Check if player skipped any winning tile since riichi
    if player.skipped_win_after_riichi:
        return True  # PERMANENT furiten

    return False

def on_tile_draw_or_discard(player, tile):
    """Track if player skips winning tile after riichi."""
    if player.is_riichi and tile in player.wait_tiles:
        # Player has option to win but doesn't
        if not player.declared_win:
            player.skipped_win_after_riichi = True
            player.ron_allowed = False  # Only tsumo from now on
```

#### Furiten Decision Tree

```
Is player in tenpai?
    ↓ NO → Not furiten (can't be furiten without tenpai)
    ↓ YES
Are any wait tiles in player's own discards?
    ↓ YES → DISCARD FURITEN (permanent, cannot ron)
    ↓ NO
Has player declared riichi?
    ↓ YES
    → Did player skip any winning tile (tsumo or ron) after riichi?
        ↓ YES → RIICHI FURITEN (permanent, cannot ron)
        ↓ NO → Not furiten, can ron
    ↓ NO (not riichi)
    → Did player skip ron in current round (since last turn)?
        ↓ YES → TEMPORARY FURITEN (until next turn)
        ↓ NO → Not furiten, can ron
```

#### Implementation Validation

**Complete furiten check**:
```python
def is_furiten(player, game_state):
    """Comprehensive furiten status check."""
    if not player.is_tenpai:
        return False

    wait_tiles = calculate_wait_tiles(player.hand)

    # Check discard furiten (highest priority)
    if check_discard_furiten(player, wait_tiles):
        return True

    # Check riichi furiten (if riichi declared)
    if check_riichi_furiten(player):
        return True

    # Check temporary furiten (if not riichi)
    if check_temporary_furiten(player, game_state):
        return True

    return False

def validate_ron_declaration(player, win_tile, game_state):
    """Verify ron is allowed before accepting declaration."""
    if is_furiten(player, game_state):
        # Player attempts to ron while furiten → chombo if hand revealed
        return "FURITEN_VIOLATION"

    # Other validations (yaku check, tile count, etc.)
    return "VALID"
```

**Test requirements**:
```python
def test_discard_furiten_applies_to_all_waits():
    """Verify discarding one wait tile makes all waits furiten."""
    player = create_player(hand="258m wait")
    player.discard_tile("5m")

    assert is_furiten(player) == True
    assert can_ron_on(player, "2m") == False
    assert can_ron_on(player, "5m") == False
    assert can_ron_on(player, "8m") == False
    assert can_tsumo_on(player, "2m") == True  # Tsumo still allowed

def test_temporary_furiten_clears_on_next_turn():
    """Verify temporary furiten clears when player's turn comes."""
    player = create_player(riichi=False)
    opponent_discards("3m")  # Player's wait tile
    player.skip_ron()

    assert is_furiten(player) == True  # Temporary furiten active

    # Other players discard
    assert can_ron_on(player, "6m") == False

    # Player's turn arrives
    player.draw_tile()

    assert is_furiten(player) == False  # Furiten cleared

def test_riichi_furiten_permanent_after_skip():
    """Verify riichi furiten is permanent after skipping win."""
    player = create_player(riichi=True)
    player.draw_tile("3m")  # Winning tile
    player.skip_tsumo()  # Skip winning

    assert is_furiten(player) == True
    assert player.can_ron == False

    # Several turns later
    opponent_discards("3m")
    assert can_ron_on(player, "3m") == False  # Still furiten
    assert can_tsumo_on(player, "3m") == True  # Tsumo allowed
```

#### Common Furiten Scenarios

**Scenario 1: Changing wait strategy**
```
Early game: Building 3-4-5 man sequence, discard 3m
Mid game: Hand transforms to 1-2-3m + 6-7-8m, now waiting on 3m
Result: Discard furiten on 3m → cannot ron, only tsumo
```

**Scenario 2: Shanpon wait trap**
```
Hand: 東東 + 南南 (waiting on 東 or 南)
Earlier: Discarded 東 when building hand
Now: Reach tenpai with 東-南 wait
Result: Furiten on both 東 and 南 → cannot ron
```

**Scenario 3: Riichi furiten by accident**
```
Declare riichi: waiting on 4-7m
Draw: 7m → Player hesitates, forgets to call tsumo, discards
Result: PERMANENT riichi furiten → cannot ron for rest of hand
```

**Scenario 4: Temporary furiten strategy**
```
Not riichi, waiting on 3-6m
Dealer discards 3m → Skip ron (want larger payment from non-dealer)
Next player discards 6m → Cannot ron (temporary furiten)
Your turn comes → Draw tile → Furiten cleared
Later: Same non-dealer discards 3m → Can ron now (worth the wait)
```

#### Related Penalties

**Furiten ron attempt**:
- Player in furiten declares ron
- If hand NOT revealed yet → **Win prohibition** (see constitution_penalties.md 4.4)
- If hand already revealed → **Chombo** (see constitution_penalties.md 5.4.4)

**No-tenpai riichi at exhaustive draw**:
- Declare riichi (actually not tenpai or furiten)
- Hand goes to exhaustive draw
- Must reveal hand → Exposed as no-tenpai → **Chombo** (see constitution_penalties.md 5.4.2)

## Game End Conditions

### Exhaustive Draw (流局/Ryūkyoku/류국)

**Occurs when**:
- Live wall depletes (14 tiles remain in dead wall)
- No one declares win

**Tenpai settlement (聴牌料)**:
- Players in tenpai reveal hands
- Non-tenpai players pay 1,000 points each → split among tenpai players
- If dealer is tenpai: dealer continues (連莊/renchan)
- If dealer not tenpai: rotate dealer

### Abortive Draws (途中流局)

**Nine Terminals (九種九牌/Kyūshu Kyūhai)**:
- Dealer's first turn with 9+ different terminal/honor tiles
- Optional declaration (can continue playing)

**Four Wind Discards (四風連打)**:
- All four players discard same wind tile on first round

**Four Kan Declared (四槓算了)**:
- 4 total kan declarations across all players (or 4 by one player in some rulesets)

**Triple Ron (三家和)**:
- Three players declare ron simultaneously on same discard

### End of Round (局/Kyoku)

Game proceeds through wind rounds:
- **East Round** (東場): East-1 through East-4
- **South Round** (南場): South-1 through South-4
- Each sub-round advances unless dealer wins (continues) or is in tenpai at exhaustive draw

### Game Termination

**Standard rule**: Game ends after South-4
**Extended play**: Can continue to West/North rounds if no player reaches 30,000+ points

## Implementation Notes for majhong_cli

### Tile State Tracking

```
- hand_tiles: [tile_list]           # Concealed tiles in hand
- called_melds: [(type, tiles)]     # Chi/pon/kan melds
- discards: [tile_list]             # Chronological discard order
- riichi_discard_index: int | null  # Position of riichi declaration
- tsumo_tile: tile | null           # Just-drawn tile (kept separate)
```

### Wall State Management

```
- live_wall: [tile_list]            # Remaining drawable tiles (max 122)
- dead_wall: {
    rinshan: [4 tiles],
    dora_indicators: [5 tiles],
    ura_indicators: [5 tiles]
  }
- dora_visible_count: int           # Number of flipped dora (1~5)
- wall_tiles_remaining: int         # For haitei/houtei detection
```

### Rule Configurability

```yaml
# Example config structure
mahjong_variant: "riichi"           # riichi | standard
red_fives_enabled: true
local_yaku_enabled: []              # Additional yakuman or local rules
min_han_requirement: 1              # Some rulesets require 2+ han
```

### Game State Validation

**MUST verify at each action**:
- Furiten status before allowing ron
- Minimum yaku requirement before declaring win
- Remaining wall tiles before riichi declaration
- Dead wall always exactly 14 tiles after kan
- Valid meld formations (sequence same-suit, triplet exact match)

## Terminology Mapping

| English | Japanese | Korean | Code |
|---------|----------|--------|------|
| Self-draw | 自摸 (Tsumo) | 쯔모 | `TSUMO` |
| Win on discard | 栄和 (Ron) | 론 | `RON` |
| Dealer | 親 (Oya) | 오야/선 | `DEALER` |
| Non-dealer | 子 (Ko) | 코 | `NON_DEALER` |
| Sequence | 順子 (Shuntsu) | 순자 | `SEQUENCE` |
| Triplet | 刻子 (Koutsu) | 커쯔 | `TRIPLET` |
| Pair | 対子 (Toitsu) | 토이쯔 | `PAIR` |
| Dead wall | 王牌 (Wanpai) | 왕패 | `DEAD_WALL` |
| Live wall | 活牌 (Huopai) | 활패 | `LIVE_WALL` |
| Dora | ドラ (Dora) | 도라 | `DORA` |
| Ura-dora | 裏ドラ (Ura-dora) | 뒷도라/우라도라 | `URA_DORA` |
| Ready hand | 聴牌 (Tenpai) | 텐파이 | `TENPAI` |
| Forbidden win | 振聴 (Furiten) | 후리텐 | `FURITEN` |

## References

- Namu Wiki (Korean): https://namu.wiki/w/리치마작/진행
- Official Riichi rules: Japan Mahjong Federation
- MajSoul implementation: https://mahjongsoul.com/
