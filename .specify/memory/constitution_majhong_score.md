# Riichi Mahjong Scoring System Reference

**Purpose**: Complete scoring calculation system for majhong_cli implementation.
**Referenced by**: `.specify/memory/constitution.md`, game engine
**Last Updated**: 2025-11-17
**Sources**: MajSoul scoring tables (score1.png, score2.png, score3.png)

## Scoring Components

### Han / 판 / 翻 / 番 (Doubling Value)

**Definition**: Multiplier based on yaku combinations

**Sources**:
- Yaku base han (see `constitution_yaku.md`)
- Dora bonus tiles (+1 han each)
- Ura-dora (revealed after riichi, +1 han each)
- Aka-dora (red fives, +1 han each if enabled)

**Formula**: `Total Han = Yaku Han + Dora Han + Ura-dora Han + Aka-dora Han`

### Fu / 부 / 符 / 符 (Base Points)

**Definition**: Base point value from hand composition

**Base fu**: 20 fu (always)

**Additional fu sources**:

| Component | Closed (암/暗) | Open (명/明) |
|-----------|---------------|-------------|
| **Win method** |
| Ron (discard win) | +10 fu | +10 fu |
| Tsumo (self-draw, closed) | +2 fu | N/A |
| **Sequences** | 0 fu | 0 fu |
| **Triplets - Simples (2-8)** | +4 fu | +2 fu |
| **Triplets - Terminals/Honors (1,9,winds,dragons)** | +8 fu | +4 fu |
| **Quads - Simples** | +16 fu | +8 fu |
| **Quads - Terminals/Honors** | +32 fu | +16 fu |
| **Pair (head)** |
| Seat wind | +2 fu | +2 fu |
| Round wind | +2 fu | +2 fu |
| Dragons | +2 fu | +2 fu |
| Double wind (seat = round) | +4 fu | +4 fu |
| **Wait type** |
| Pair wait (tanki) | +2 fu | +2 fu |
| Edge wait (penchan) | +2 fu | +2 fu |
| Closed wait (kanchan) | +2 fu | +2 fu |
| Two-sided (ryanmen) | 0 fu | 0 fu |

**Fu rounding**: Round up to nearest 10 (e.g., 32 fu → 40 fu, 38 fu → 40 fu)

**Special cases**:
- **Chiitoitsu (七対子)**: Always 25 fu (fixed, no rounding)
- **Pinfu tsumo**: 20 fu (no rounding, stays at 20)
- **Open pinfu ron**: 30 fu minimum (20 base + 10 ron)

### Example Fu Calculation

**Example 1: Closed hand, ron win**
```
Hand: [2m-3m-4m] [5p-6p-7p] [White-White-White] [8s-8s-8s] [East-East] (pair)
Win: Ron on 4m
Seat: East, Round: East

Fu calculation:
- Base: 20 fu
- Ron: +10 fu
- White dragon triplet (closed, terminals/honors): +8 fu
- 8s triplet (closed, simples): +4 fu
- East pair (seat + round wind): +4 fu
Total: 46 fu → rounds to 50 fu
```

**Example 2: Pinfu tsumo**
```
Hand: All sequences, valueless pair, ryanmen wait
Win: Tsumo

Fu calculation:
- Base: 20 fu
- Tsumo (pinfu exception): +2 fu
Total: 22 fu → pinfu special rule: stays at 20 fu (not rounded)
```

## Scoring Tables

### Non-Dealer (자사 / 子 / 子) Scoring

**Table format**: Fu (rows) × Han (columns) = [Ron points / Tsumo from non-dealer, Tsumo from dealer]

| Fu | 1 Han | 2 Han | 3 Han | 4 Han |
|----|-------|-------|-------|-------|
| **20** | — | 700, 400 | 1300, 700 | 2600, 1300 |
| **25** | — | 1600, 800 | 3200, 1600 | 6400, 3200 |
| **30** | 1000, 300/500 | 2000, 500/1000 | 3900, 1000/2000 | 7700, 2000/3900 |
| **40** | 1300, 400/700 | 2600, 700/1300 | 5200, 1300/2600 | **Mangan** |
| **50** | 1600, 400/800 | 3200, 800/1600 | 6400, 1600/3200 | **Mangan** |
| **60** | 2000, 500/1000 | 3900, 1000/2000 | 7700, 2000/3900 | **Mangan** |
| **70** | 2300, 600/1200 | 4500, 1200/2300 | **Mangan** | **Mangan** |
| **80** | 2600, 700/1300 | 5200, 1300/2600 | **Mangan** | **Mangan** |
| **90** | 2900, 800/1500 | 5800, 1500/2900 | **Mangan** | **Mangan** |
| **100** | 3200, 800/1600 | 6400, 1600/3200 | **Mangan** | **Mangan** |
| **110** | 3600, 900/1800 | 7100, 1800/3600 | **Mangan** | **Mangan** |

**Note**: Tsumo payment format: `[from each non-dealer] / [from dealer]`

### Dealer (오야 / 親 / 莊) Scoring

**Table format**: Fu (rows) × Han (columns) = [Ron points / Tsumo from each player]

| Fu | 1 Han | 2 Han | 3 Han | 4 Han |
|----|-------|-------|-------|-------|
| **20** | — | 700 | 1300 | 2600 |
| **25** | — | 2400, 1600 | 4800, 3200 | 9600, 6400 |
| **30** | 1500, 500 | 2900, 1000 | 5800, 2000 | 11600, 3900 |
| **40** | 2000, 700 | 3900, 1300 | 7700, 2600 | **Mangan** |
| **50** | 2400, 800 | 4800, 1600 | 9600, 3200 | **Mangan** |
| **60** | 2900, 1000 | 5800, 2000 | 11600, 3900 | **Mangan** |
| **70** | 3400, 1200 | 6800, 2300 | **Mangan** | **Mangan** |
| **80** | 3900, 1300 | 7700, 2600 | **Mangan** | **Mangan** |
| **90** | 4400, 1500 | 8700, 2900 | **Mangan** | **Mangan** |
| **100** | 4800, 1600 | 9600, 3200 | **Mangan** | **Mangan** |
| **110** | 5300, 1800 | 10600, 3600 | **Mangan** | **Mangan** |

**Note**: Dealer receives 1.5× points compared to non-dealer

### Limit Hands (Mangan and Above)

| Limit | Han Range | Non-Dealer (Ron / Tsumo each) | Dealer (Ron / Tsumo each) |
|-------|-----------|-------------------------------|---------------------------|
| **Mangan (만관 / 満貫 / 滿貫)** | 5 han OR 4 han 40+ fu OR 3 han 70+ fu | 8,000 / 2,000, 4,000 | 12,000 / 4,000 |
| **Haneman (하네만 / 跳満 / 跳滿)** | 6-7 han | 12,000 / 3,000, 6,000 | 18,000 / 6,000 |
| **Baiman (배만 / 倍満 / 倍滿)** | 8-10 han | 16,000 / 4,000, 8,000 | 24,000 / 8,000 |
| **Sanbaiman (삼배만 / 三倍満 / 三倍滿)** | 11-12 han | 24,000 / 6,000, 12,000 | 36,000 / 12,000 |
| **Yakuman (역만 / 役満 / 役滿)** | 13+ han (counted) OR yakuman yaku | 32,000 / 8,000, 16,000 | 48,000 / 16,000 |
| **Double Yakuman (더블역만 / ダブル役満 / 雙倍役滿)** | Multiple yakuman OR special yakuman | 64,000 / 16,000, 32,000 | 96,000 / 32,000 |

**Tsumo payment format**:
- Non-dealer tsumo: [from each non-dealer, from dealer]
- Dealer tsumo: [from each non-dealer]

## Score Calculation Algorithm

### Step-by-Step Process

```python
def calculate_score(hand, win_tile, win_method, game_state):
    """
    Complete scoring calculation.

    Args:
        hand: Player's hand (tiles + melds)
        win_tile: Winning tile
        win_method: TSUMO or RON
        game_state: Current game state (dealer, round, etc.)

    Returns:
        ScoreResult with points, han, fu, yaku list
    """

    # Step 1: Detect yaku
    yaku_list = detect_yaku(hand, win_tile, win_method, game_state)

    if len(yaku_list) == 0:
        return ScoreResult.INVALID  # No yaku, cannot win

    # Step 2: Calculate han
    yaku_han = sum(y.han_value(hand.is_closed) for y in yaku_list)
    dora_han = count_dora(hand, game_state.dora_indicators)
    ura_dora_han = count_ura_dora(hand, game_state) if game_state.riichi_declared else 0
    aka_dora_han = count_aka_dora(hand) if game_state.red_fives_enabled else 0

    total_han = yaku_han + dora_han + ura_dora_han + aka_dora_han

    # Step 3: Calculate fu (skip if yakuman)
    if any(y.is_yakuman() for y in yaku_list) or total_han >= 13:
        fu = 0  # Fu irrelevant for yakuman
        limit = determine_yakuman_limit(yaku_list, total_han)
        points = get_yakuman_points(limit, game_state.is_dealer, win_method)
        return ScoreResult(points, total_han, fu, yaku_list, limit)

    # Calculate fu for non-yakuman hands
    fu = calculate_fu(hand, win_tile, win_method, yaku_list, game_state)

    # Step 4: Determine limit or use table
    limit = determine_limit(total_han, fu)

    if limit:
        points = get_limit_points(limit, game_state.is_dealer, win_method)
    else:
        points = get_table_points(total_han, fu, game_state.is_dealer, win_method)

    return ScoreResult(points, total_han, fu, yaku_list, limit)
```

### Fu Calculation Function

```python
def calculate_fu(hand, win_tile, win_method, yaku_list, game_state):
    """Calculate fu (base points) from hand composition."""

    # Special cases first
    if Yaku.CHIITOITSU in yaku_list:
        return 25  # Fixed 25 fu for seven pairs

    if Yaku.PINFU in yaku_list and win_method == WinMethod.TSUMO:
        return 20  # Pinfu tsumo exception

    # Standard fu calculation
    fu = 20  # Base

    # Win method
    if win_method == WinMethod.RON:
        fu += 10
    elif win_method == WinMethod.TSUMO and hand.is_closed:
        fu += 2  # Closed tsumo (unless pinfu)

    # Melds
    for meld in hand.melds:
        fu += calculate_meld_fu(meld, hand.is_concealed(meld))

    # Pair
    pair_tile = hand.pair_tile
    if is_dragon(pair_tile):
        fu += 2
    if is_seat_wind(pair_tile, game_state.seat):
        fu += 2
    if is_round_wind(pair_tile, game_state.round):
        fu += 2

    # Wait type (if not ryanmen)
    wait_type = determine_wait_type(hand, win_tile)
    if wait_type in [WaitType.TANKI, WaitType.KANCHAN, WaitType.PENCHAN]:
        fu += 2

    # Round up to nearest 10
    fu = ceil_to_10(fu)

    # Minimum 30 fu for open hands (unless pinfu/chiitoitsu)
    if not hand.is_closed and fu < 30:
        fu = 30

    return fu

def calculate_meld_fu(meld, is_concealed):
    """Calculate fu for single meld."""
    if meld.type == MeldType.SEQUENCE:
        return 0  # Sequences have no fu

    # Triplets and quads
    is_terminal_or_honor = is_terminal(meld.tiles[0]) or is_honor(meld.tiles[0])

    if meld.type == MeldType.TRIPLET:
        if is_concealed:
            return 8 if is_terminal_or_honor else 4
        else:
            return 4 if is_terminal_or_honor else 2

    elif meld.type == MeldType.QUAD:
        if is_concealed:
            return 32 if is_terminal_or_honor else 16
        else:
            return 16 if is_terminal_or_honor else 8

    return 0
```

### Limit Determination

```python
def determine_limit(han, fu):
    """Determine if hand reaches limit (mangan or higher)."""

    # Yakuman (13+ han)
    if han >= 13:
        return Limit.YAKUMAN

    # Sanbaiman (11-12 han)
    if han >= 11:
        return Limit.SANBAIMAN

    # Baiman (8-10 han)
    if han >= 8:
        return Limit.BAIMAN

    # Haneman (6-7 han)
    if han >= 6:
        return Limit.HANEMAN

    # Mangan (5 han OR 4 han 40+ fu OR 3 han 70+ fu)
    if han >= 5:
        return Limit.MANGAN
    if han == 4 and fu >= 40:
        return Limit.MANGAN
    if han == 3 and fu >= 70:
        return Limit.MANGAN

    # Below mangan, use table
    return None
```

## Payment Calculation

### Ron Payment (Discard Win)

**Single payer**: Discarding player pays full amount

```python
def calculate_ron_payment(points, discarding_player, winning_player):
    """Calculate ron payment."""
    return Payment(
        from_player=discarding_player,
        to_player=winning_player,
        amount=points
    )
```

### Pao / Responsibility Payment (책임지불/責任払い)

**Definition**: Special payment rule where a player who enables certain high-value yaku (typically yakuman) must bear responsibility for the score payment.

**When Pao applies**:
- Player makes a call (pon/kan) that confirms/locks in a yakuman for another player
- Typically applies to: Daisangen (大三元), Daisuushii (大四喜), Suukantsu (四槓子)

#### Daisangen Pao (大三元 責任払い)

**Trigger condition**:
- Two dragon types already called/revealed (open melds)
- Player discards the third dragon type → opponent calls pon/kan
- This **confirms** Daisangen (all 3 dragon triplets)

**Payment responsibility**:
- **Tsumo win**: Pao player pays **full amount** alone
- **Ron win**: Discarder and Pao player **split payment 50/50**

**Example**:
```
Player A has called: 白白白 (open), 發發發 (open)
Player B discards: 中 → Player A calls pon → 中中中 (Daisangen confirmed)
Player B is now "Pao responsible" for Player A's Daisangen

If Player A wins by tsumo: Player B pays full yakuman alone
If Player A wins by ron from Player C: Player B + Player C each pay half
```

**Implementation**:
```python
def check_daisangen_pao(player, called_tile, open_melds):
    """Check if calling this tile triggers Daisangen Pao."""
    dragon_types = ["白", "發", "中"]  # White, Green, Red

    # Count how many dragon types already open
    open_dragons = set()
    for meld in open_melds:
        if meld.tiles[0] in dragon_types:
            open_dragons.add(meld.tiles[0])

    # If 2 dragons open and calling 3rd → Pao
    if len(open_dragons) == 2 and called_tile in dragon_types:
        if called_tile not in open_dragons:
            return True  # Pao responsibility triggered

    return False
```

#### Daisuushii Pao (大四喜 責任払い)

**Trigger condition**:
- Three wind types already called/revealed (open melds)
- Player discards the fourth wind → opponent calls pon/kan
- This **confirms** Daisuushii (all 4 wind quads)

**Payment responsibility**: Same as Daisangen
- Tsumo: Pao player pays full
- Ron: Discarder + Pao player split 50/50

**Implementation**:
```python
def check_daisuushii_pao(player, called_tile, open_melds):
    """Check if calling this tile triggers Daisuushii Pao."""
    wind_types = ["東", "南", "西", "北"]  # E/S/W/N

    open_winds = set()
    for meld in open_melds:
        if meld.tiles[0] in wind_types:
            open_winds.add(meld.tiles[0])

    # If 3 winds open and calling 4th → Pao
    if len(open_winds) == 3 and called_tile in wind_types:
        if called_tile not in open_winds:
            return True

    return False
```

#### Suukantsu Pao (四槓子 責任払い)

**Trigger condition**:
- One player has already declared 3 kans
- Another player discards a tile that has **never appeared** (生패/fresh tile)
- First player calls daiminkan (open kan) → **confirms Suukantsu** (4 kans)

**Fresh tile definition** (生패):
- Tile that has not been discarded by anyone yet
- No copies visible in any discard pile

**Payment responsibility**: Same as above
- Tsumo: Pao player pays full
- Ron: Discarder + Pao player split 50/50

**Special note**: Some rulesets trigger abortive draw (四槓散了) instead of allowing 4th kan.

**Implementation**:
```python
def is_fresh_tile(tile, all_discards):
    """Check if tile has never been discarded (生패)."""
    for discard_pile in all_discards:
        if tile in discard_pile:
            return False
    return True

def check_suukantsu_pao(player, called_tile, all_discards):
    """Check if calling this tile triggers Suukantsu Pao."""
    # Count player's existing kans
    kan_count = sum(1 for meld in player.melds if meld.type == "KAN")

    if kan_count == 3 and is_fresh_tile(called_tile, all_discards):
        # Calling 4th kan with fresh tile → Pao
        return True

    return False
```

#### Optional Pao: Rinshan Kaihou from Daiminkan

**Some rulesets** apply Pao when:
- Player calls daiminkan (open kan from discard)
- Immediately wins on rinshan draw (嶺上開花)
- Discarder who enabled the kan becomes Pao responsible

**Payment**: Full payment from kan-enabling discarder

**Note**: Not universally adopted - check ruleset configuration.

#### Pao Payment Calculation

**Payment function**:
```python
def calculate_pao_payment(winner, pao_player, discarder, points, win_type):
    """
    Calculate payment when Pao responsibility applies.

    Args:
        winner: Player who won with yakuman
        pao_player: Player who triggered Pao
        discarder: Player who discarded winning tile (for ron)
        points: Total yakuman points
        win_type: "TSUMO" or "RON"

    Returns:
        List of Payment objects
    """
    if win_type == "TSUMO":
        # Pao player pays everything
        return [Payment(from_player=pao_player, to_player=winner, amount=points)]

    elif win_type == "RON":
        # Split 50/50 between discarder and Pao player
        half_payment = points // 2
        return [
            Payment(from_player=discarder, to_player=winner, amount=half_payment),
            Payment(from_player=pao_player, to_player=winner, amount=half_payment)
        ]
```

#### Pao Validation and Tracking

**Game state must track**:
```python
class GameState:
    def __init__(self):
        self.pao_responsibilities = {}  # {player_id: (responsible_for_player, yaku_type)}

    def on_call_declared(self, caller, called_tile, call_type):
        """Check if this call triggers Pao."""
        if check_daisangen_pao(caller, called_tile, caller.open_melds):
            # Find who discarded the tile
            discarder = self.last_discarder
            self.pao_responsibilities[caller.id] = (discarder.id, "DAISANGEN")

        elif check_daisuushii_pao(caller, called_tile, caller.open_melds):
            discarder = self.last_discarder
            self.pao_responsibilities[caller.id] = (discarder.id, "DAISUUSHII")

        elif check_suukantsu_pao(caller, called_tile, self.all_discards):
            discarder = self.last_discarder
            self.pao_responsibilities[caller.id] = (discarder.id, "SUUKANTSU")

    def calculate_win_payment(self, winner, win_type, points):
        """Calculate payment considering Pao if applicable."""
        if winner.id in self.pao_responsibilities:
            pao_player_id, yaku_type = self.pao_responsibilities[winner.id]
            pao_player = self.get_player(pao_player_id)

            if win_type == "RON":
                discarder = self.last_discarder
                return calculate_pao_payment(winner, pao_player, discarder, points, "RON")
            else:  # TSUMO
                return calculate_pao_payment(winner, pao_player, None, points, "TSUMO")

        # Normal payment if no Pao
        return calculate_normal_payment(winner, win_type, points)
```

### Tsumo Payment (Self-Draw Win)

**All pay**: Each player pays winner

```python
def calculate_tsumo_payment(base_points, is_dealer_win, winner, all_players):
    """
    Calculate tsumo payment distribution.

    Returns list of payments from each player to winner.
    """
    payments = []

    if is_dealer_win:
        # Dealer tsumo: Each non-dealer pays same amount
        amount_each = base_points  # From table
        for player in all_players:
            if player != winner:
                payments.append(Payment(
                    from_player=player,
                    to_player=winner,
                    amount=amount_each
                ))

    else:
        # Non-dealer tsumo: Different amounts from dealer vs non-dealers
        amount_from_non_dealer, amount_from_dealer = base_points  # Tuple from table

        for player in all_players:
            if player == winner:
                continue
            elif is_dealer(player):
                payments.append(Payment(
                    from_player=player,
                    to_player=winner,
                    amount=amount_from_dealer
                ))
            else:
                payments.append(Payment(
                    from_player=player,
                    to_player=winner,
                    amount=amount_from_non_dealer
                ))

    return payments
```

### Honba and Riichi Stick Bonuses

```python
def apply_bonuses(base_payments, game_state):
    """Add honba and riichi stick bonuses."""

    # Honba bonus: +300 per honba stick
    honba_bonus = game_state.honba_count * 300

    # Riichi stick bonus: All riichi deposits go to winner
    riichi_bonus = game_state.riichi_stick_count * 1000

    # Add bonuses to winner's total
    for payment in base_payments:
        if payment.to_player == winner:
            payment.amount += honba_bonus + riichi_bonus

    return base_payments
```

## Game End Conditions

### Round Progression

**Standard game structure**:
- **East Round (동장 / 東場 / 東場)**: East-1 through East-4
- **South Round (남장 / 南場 / 南場)**: South-1 through South-4

**Dealer continuation (연장 / 連荘 / 連莊)**:
- Dealer wins → Dealer continues, add 1 honba
- Dealer tenpai at exhaustive draw → Dealer continues, add 1 honba
- Dealer loses or not tenpai → Rotate dealer, reset honba

### Game Termination

**Standard end**: After South-4 (all last / オーラス)

**Early termination conditions**:
1. **Negative points**: Any player reaches negative points (bankruptcy)
2. **30,000+ points**: Some rules end if any player reaches 30,000+ after South-4
3. **West/North rounds**: Extended play if no one reaches target (optional)

### Final Ranking

**Rank determination**:
1. Highest points = 1st place
2. Ties broken by seat order (East > South > West > North)

**Uma (順位点 / placement bonus)**:
- Common distribution: +15,000 / +5,000 / -5,000 / -15,000 (1st/2nd/3rd/4th)
- Or: +20,000 / +10,000 / -10,000 / -20,000 (more aggressive)

**Oka (オカ / 오카 / top bonus)**:
- **Definition**: Bonus points for 1st place derived from starting point differential
- **Common setup**: Starting points = 25,000, Return points (反環点) = 30,000
- **Calculation**: (Return points - Starting points) × 4 players = Oka bonus
  - Example: (30,000 - 25,000) × 4 = 20,000 points → +20.0 uma to 1st place
- **Why it exists**: Compensates 1st place for the "missing" points when starting below return baseline

**Oka calculation examples**:

| Starting Points | Return Points | Oka Bonus (to 1st) |
|----------------|---------------|---------------------|
| 25,000 | 30,000 | +20,000 (+20.0 uma) |
| 30,000 | 30,000 | 0 (no oka) |
| 27,000 | 30,000 | +12,000 (+12.0 uma) |

**Final score calculation with Oka**:
```
Raw score = (Player's ending points - Return points) / 1000
Final score = Raw score + Uma + Oka (if 1st place)
```

**Example** (25k start, 30k return, oka +20):
```
Player A: 35,000 points (1st) → Raw: +5.0, Uma: +15.0, Oka: +20.0 → Final: +40.0
Player B: 28,000 points (2nd) → Raw: -2.0, Uma: +5.0, Oka: 0 → Final: +3.0
Player C: 22,000 points (3rd) → Raw: -8.0, Uma: -5.0, Oka: 0 → Final: -13.0
Player D: 15,000 points (4th) → Raw: -15.0, Uma: -15.0, Oka: 0 → Final: -30.0
Total: +40 +3 -13 -30 = 0 ✓ (must sum to zero)
```

**Extended play (南入/西入) with Oka**:
- **Condition**: After All Last (오라스), no player reaches return points
- **Rule**: Game extends to next round (West or North round) in sudden death mode
- **Termination**: First player to reach return points ends the game immediately

**Example of 西入 (West round extension)**:
```
After South-4 ends:
Player A: 29,000 (all below 30,000 return points)
Player B: 28,500
Player C: 21,500
Player D: 21,000

→ Game extends to West-1 (西1局)
→ Continue until someone reaches 30,000+
→ Game ends immediately when condition met
```

**Implementation**:
```python
def calculate_final_scores(players, starting_points=25000, return_points=30000, uma=[15, 5, -5, -15]):
    """
    Calculate final scores with Uma and Oka.

    Args:
        players: List of players sorted by ending points (descending)
        starting_points: Points each player starts with
        return_points: Baseline for score calculation
        uma: Uma distribution [1st, 2nd, 3rd, 4th]

    Returns:
        List of final scores
    """
    # Calculate Oka bonus (only for 1st place)
    oka_bonus = ((return_points - starting_points) * 4) / 1000  # Convert to uma scale

    final_scores = []
    for rank, player in enumerate(players):
        # Raw score: (ending points - return points) / 1000
        raw_score = (player.points - return_points) / 1000

        # Uma for this rank
        rank_uma = uma[rank]

        # Oka only for 1st place
        rank_oka = oka_bonus if rank == 0 else 0

        final_score = raw_score + rank_uma + rank_oka
        final_scores.append(final_score)

    # Validate: total must be zero
    assert abs(sum(final_scores)) < 0.01, "Final scores must sum to zero"

    return final_scores

def check_extended_play_needed(players, return_points=30000):
    """Check if game needs extension (南入/西入)."""
    return all(player.points < return_points for player in players)

def check_game_end_in_extension(players, return_points=30000):
    """In extension round, end immediately if anyone reaches return points."""
    return any(player.points >= return_points for player in players)
```

### Tenpai Settlement (Exhaustive Draw)

**Conditions**: Wall exhausted, no winner

**Tenpai payments**:
- Players in tenpai receive 1,000 points each
- Payments come from non-tenpai players
- Total pool: 3,000 points

**Examples**:
- 1 tenpai, 3 not: +3,000 to tenpai player
- 2 tenpai, 2 not: +1,500 each to tenpai players
- 3 tenpai, 1 not: +1,000 each to tenpai players, -3,000 to non-tenpai
- 0 or 4 tenpai: No payments

## Implementation Example

```python
class ScoringEngine:
    def __init__(self, rules=RiichiRules.STANDARD):
        self.rules = rules
        self.scoring_table = load_scoring_table()

    def score_hand(self, hand, win_tile, win_method, game_state):
        """Main scoring entry point."""

        # Validate minimum yaku
        yaku = detect_yaku(hand, win_tile, win_method, game_state)
        if not yaku:
            raise InvalidWinError("No yaku present")

        # Calculate han and fu
        han = calculate_total_han(hand, yaku, game_state)
        fu = calculate_fu(hand, win_tile, win_method, yaku, game_state)

        # Determine points
        limit = determine_limit(han, fu)
        if limit:
            points = self.get_limit_points(limit, game_state.is_dealer, win_method)
        else:
            points = self.scoring_table.lookup(han, fu, game_state.is_dealer, win_method)

        # Calculate payments
        if win_method == WinMethod.RON:
            payments = calculate_ron_payment(points, game_state)
        else:
            payments = calculate_tsumo_payment(points, game_state)

        # Apply bonuses
        payments = apply_bonuses(payments, game_state)

        return ScoreResult(
            points=points,
            han=han,
            fu=fu,
            yaku=yaku,
            limit=limit,
            payments=payments
        )
```

## Testing Requirements

```python
def test_pinfu_tsumo_20_fu():
    """Verify pinfu tsumo stays at 20 fu."""
    hand = create_pinfu_hand()
    result = score_hand(hand, win_tile, WinMethod.TSUMO, game_state)

    assert result.fu == 20
    assert result.han >= 1
    assert Yaku.PINFU in result.yaku

def test_mangan_threshold():
    """Verify 5 han reaches mangan."""
    hand = create_hand_with_5_han()
    result = score_hand(hand, win_tile, win_method, game_state)

    assert result.limit == Limit.MANGAN
    assert result.points == 8000  # Non-dealer ron

def test_dealer_payment_multiplier():
    """Verify dealer receives 1.5× points."""
    hand = create_1_han_30_fu_hand()

    # Non-dealer
    result_non_dealer = score_hand(hand, win_tile, WinMethod.RON, non_dealer_state)
    assert result_non_dealer.points == 1000

    # Dealer
    result_dealer = score_hand(hand, win_tile, WinMethod.RON, dealer_state)
    assert result_dealer.points == 1500  # 1.5× multiplier

def test_tsumo_payment_distribution():
    """Verify tsumo payments split correctly."""
    hand = create_3_han_30_fu_hand()
    result = score_hand(hand, win_tile, WinMethod.TSUMO, non_dealer_state)

    # Non-dealer tsumo: 3900 points (1000 from each non-dealer, 2000 from dealer)
    assert sum(p.amount for p in result.payments) == 3900
    assert len([p for p in result.payments if p.amount == 1000]) == 2  # Non-dealers
    assert len([p for p in result.payments if p.amount == 2000]) == 1  # Dealer
```

## References

- See `constitution_yaku.md` for yaku han values
- See `constitution_term.md` for terminology
- MajSoul scoring tables: score1.png, score2.png, score3.png
- Japanese Mahjong scoring guide: https://ja.wikipedia.org/wiki/麻雀の得点計算
- Riichi Mahjong scoring calculator: Various online tools for validation
