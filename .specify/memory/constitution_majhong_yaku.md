# Riichi Mahjong Yaku (Scoring Patterns) Reference

**Purpose**: Complete catalog of yaku (scoring patterns) for majhong_cli implementation.
**Referenced by**: `.specify/memory/constitution.md`, scoring calculation system
**Last Updated**: 2025-11-18

## Terminology

**Yaku / 역 / 役 / 役**: Scoring pattern required to win (minimum 1 yaku needed, dora doesn't count)

**Format**: English / Korean / Japanese / Chinese for all yaku names

**Open/Closed notation**:
- ✅ **Open OK**: Works with called melds (chi/pon/kan)
- 🔒 **Closed only**: Requires menzen (no called melds)
- ⬇️ **Reduced**: -1 han when open

## 1-Han Yaku (1판역)

### Riichi / 리치 / 立直 / 立直 🔒

**Han**: 1 (closed only)
**Conditions**:
- Declare riichi when tenpai (one tile from winning)
- Pay 1,000 point stake
- Hand must be closed (no chi/pon/open kan)

**Requirements**:
1. Announce "riichi" before discarding
2. Discard tile sideways to mark declaration
3. Cannot change hand composition after declaration
4. Must have 1,000+ points to declare
5. At least 4 tiles must remain in wall

**Effects**:
- Unlocks ippatsu (1-shot) opportunity
- Reveals ura-dora on winning
- Riichi stick goes to winner
- Furiten restrictions apply to all wait tiles

**Code**: `RIICHI`

---

### Ippatsu / 일발 / 一発 / 一發 🔒

**Han**: 1 (bonus to riichi, closed only)
**Conditions**:
- Win within one turn cycle after riichi declaration
- No interruptions (chi/pon/kan) between riichi and win

**Destroyed by**:
- Any player calling chi/pon/kan
- Hand reaching next turn without winning

**Code**: `IPPATSU`

---

### Menzen Tsumo / 멘젠쯔모 / 門前清自摸和 / 門前清自摸和 🔒

**Han**: 1 (closed only)
**Conditions**:
- Win by self-draw (tsumo)
- Hand completely closed (no called melds)

**Note**: Ankan (concealed kan) allowed, hand remains closed

**Code**: `MENZEN_TSUMO`

---

### Tanyao / 탕야오 / 断幺九 / 斷幺九 ✅

**Han**: 1 (open or closed)
**Conditions**:
- All tiles are simples (2-8 of man/pin/sou)
- No terminals (1, 9) or honors (winds, dragons)

**Tile requirements**:
- ✅ Allowed: 2m-8m, 2p-8p, 2s-8s
- ❌ Forbidden: 1m, 9m, 1p, 9p, 1s, 9s, all winds, all dragons

**Rule variant**: Some rules forbid open tanyao (kuitan), majhong_cli should make this configurable

**Code**: `TANYAO`

---

### Yakuhai / 역패 / 役牌 / 役牌 ✅

**Han**: 1 per triplet (open or closed)
**Conditions**:
- Triplet of value tiles:
  - Dragon tiles (white/green/red) - always valuable
  - Seat wind (your position)
  - Round wind (current round)

**Multiple yakuhai**: Stackable (e.g., East round, East player, East triplet = 2 han from yakuhai)

**Examples**:
- White dragon triplet: 1 han (always)
- East triplet in East round by East player: 2 han (seat + round wind)
- South triplet in East round by South player: 1 han (seat wind only)

**Code**: `YAKUHAI_HAKU` (white), `YAKUHAI_HATSU` (green), `YAKUHAI_CHUN` (red), `YAKUHAI_WIND`

---

### Pinfu / 핑후 / 平和 / 平和 🔒

**Han**: 1 (closed only)
**Conditions**:
1. All melds are sequences (no triplets)
2. Pair (head) is not value tiles (no dragons, no seat/round winds)
3. Wait is ryanmen (two-sided, e.g., 4-5 waiting 3 or 6)

**Forbidden waits**: Kanchan (middle), penchan (edge), shanpon (pair), tanki (single)

**Common combination**: "Mentanpin" = Menzen tsumo + Tanyao + Pinfu (3 han closed hand)

**Code**: `PINFU`

---

### Iipeikou / 이페코 / 一盃口 / 一盃口 🔒

**Han**: 1 (closed only)
**Conditions**:
- Two identical sequences in hand
- Example: 2-3-4 man and 2-3-4 man (exact duplicates)

**Note**: Cannot use tiles from called melds

**Code**: `IIPEIKOU`

---

### Haitei Raoyue / 해저로월 / 海底撈月 / 海底撈月 ✅

**Han**: 1 (open or closed)
**Conditions**:
- Win by self-draw (tsumo) on the last tile from live wall

**Note**: Last tile = when 14 tiles remain (dead wall), you draw the final live wall tile

**Code**: `HAITEI`

---

### Houtei Raoyui / 하저로어 / 河底撈魚 / 河底撈魚 ✅

**Han**: 1 (open or closed)
**Conditions**:
- Win by ron on another player's last discard tile

**Note**: Last discard = discard made when only dead wall (14 tiles) remains

**Code**: `HOUTEI`

---

### Rinshan Kaihou / 영상개화 / 嶺上開花 / 嶺上開花 ✅

**Han**: 1 (open or closed)
**Conditions**:
- Win by self-draw on rinshan tile (replacement tile after kan declaration)

**Note**: Cannot combine with ippatsu (kan interrupts riichi cycle)

**Code**: `RINSHAN`

---

### Chankan / 창깡 / 搶槓 / 搶槓 ✅

**Han**: 1 (open or closed)
**Conditions**:
- Win by ron when another player adds 4th tile to their pon (shouminkan)

**Special**: Only time you can ron mid-turn (not from discard pile)

**Kokushi exception**: Can rob ankan for kokushi musou yakuman

**Code**: `CHANKAN`

## 2-Han Yaku (2판역)

### Chiitoitsu / 치또이츠 / 七対子 / 七對子 🔒

**Han**: 2 (closed only)
**Conditions**:
- Seven different pairs (14 tiles total)
- All pairs must be distinct types

**Special**:
- Always 25 fu (fixed)
- Cannot be combined with iipeikou, ryanpeikou, toitoi

**Code**: `CHIITOITSU`

---

### Toitoi / 또이또이 / 対々和 / 對對和 ✅

**Han**: 2 (open or closed)
**Conditions**:
- All melds are triplets/quads (no sequences)
- 4 triplets + 1 pair

**Code**: `TOITOI`

---

### Sanankou / 산안커 / 三暗刻 / 三暗刻 ✅

**Han**: 2 (open or closed)
**Conditions**:
- Three concealed triplets
- Triplets formed without calling pon
- Ron win: winning tile cannot complete the 3rd triplet
- Tsumo win: all three triplets were concealed before draw

**Code**: `SANANKOU`

---

### Sanshoku Doujun / 삼색동순 / 三色同順 / 三色同順 ⬇️

**Han**: 2 closed, 1 open
**Conditions**:
- Same sequence in all three suits
- Example: 3-4-5 man, 3-4-5 pin, 3-4-5 sou

**Code**: `SANSHOKU_DOUJUN`

---

### Ittsu / 일기통관 / 一気通貫 / 一氣通貫 ⬇️

**Han**: 2 closed, 1 open
**Conditions**:
- Three sequences 1-2-3, 4-5-6, 7-8-9 in same suit

**Code**: `ITTSU`

---

### Chanta / 찬타 / 混全帯幺九 / 混全帶幺九 ⬇️

**Han**: 2 closed, 1 open
**Conditions**:
- Every meld and pair contains at least one terminal (1 or 9) or honor tile
- Must have at least one sequence (otherwise junchan)

**Code**: `CHANTA`

---

### Sankantsu / 산깡쯔 / 三槓子 / 三槓子 ✅

**Han**: 2 (open or closed)
**Conditions**:
- Three kan (quad) declarations

**Note**: Fourth kan usually causes abortive draw (suukantsu)

**Code**: `SANKANTSU`

---

### Honroutou / 혼노두 / 混老頭 / 混老頭 ✅

**Han**: 2 (open or closed)
**Conditions**:
- All tiles are terminals (1, 9) or honors
- Must have at least one honor (otherwise chinroutou yakuman)

**Note**: Usually combined with toitoi or chiitoitsu

**Code**: `HONROUTOU`

---

### Shousangen / 소삼원 / 小三元 / 小三元 ✅

**Han**: 2 (open or closed)
**Conditions**:
- Two dragon triplets + one dragon pair
- Automatically includes 2 han from yakuhai

**Total han**: 4 han minimum (2 shousangen + 2 yakuhai)

**Code**: `SHOUSANGEN`

---

### Double Riichi / 더블리치 / ダブル立直 / 雙立直 🔒

**Han**: 2 (closed only)
**Conditions**:
- Declare riichi on first turn (before any discards)
- No interruptions (chi/pon/kan) before your first discard

**Code**: `DOUBLE_RIICHI`

## 3-Han Yaku (3판역)

### Ryanpeikou / 량페코 / 両平和 / 兩盃口 🔒

**Han**: 3 (closed only)
**Conditions**:
- Two sets of iipeikou (four sequences, two pairs of identical sequences)
- Example: 2-3-4 man ×2, 5-6-7 pin ×2

**Note**: Overrides iipeikou, cannot stack with chiitoitsu

**Code**: `RYANPEIKOU`

---

### Junchan / 준찬타 / 純全帯幺九 / 純全帶幺九 ⬇️

**Han**: 3 closed, 2 open
**Conditions**:
- Every meld and pair contains terminal (1 or 9)
- NO honor tiles allowed
- Must have at least one sequence

**Code**: `JUNCHAN`

---

### Honitsu / 혼일색 / 混一色 / 混一色 ⬇️

**Han**: 3 closed, 2 open
**Conditions**:
- One number suit + honor tiles only
- Example: All man tiles + some wind/dragon tiles

**Code**: `HONITSU`

## 6-Han Yaku (6판역)

### Chinitsu / 청일색 / 清一色 / 清一色 ⬇️

**Han**: 6 closed, 5 open
**Conditions**:
- Single number suit only (man, pin, or sou)
- NO honor tiles

**Note**: Most common high-value yaku, often reaches mangan

**Code**: `CHINITSU`

## Yakuman (역만 / 役満 / 役滿) - 13+ Han

### Kokushi Musou / 국사무쌍 / 国士無双 / 國士無雙 🔒

**Han**: Yakuman (closed only)
**Conditions**:
- One each of all 13 terminal/honor types + one pair
- Tiles: 1m, 9m, 1p, 9p, 1s, 9s, East, South, West, North, White, Green, Red (one pair among these)

**Special**: Can ron on ankan (only exception)

**Code**: `KOKUSHI`

---

### Suuankou / 스안커 / 四暗刻 / 四暗刻 🔒

**Han**: Yakuman (closed only)
**Conditions**:
- Four concealed triplets
- Must win by tsumo OR ron on pair wait (tanki machi)

**Suuankou Tanki**: Double yakuman if winning tile completes pair (some rules)

**Code**: `SUUANKOU`

---

### Daisangen / 대삼원 / 大三元 / 大三元 ✅

**Han**: Yakuman (open or closed)
**Conditions**:
- Triplets of all three dragons (white, green, red)

**Code**: `DAISANGEN`

---

### Shousuushii / 소사희 / 小四喜 / 小四喜 ✅

**Han**: Yakuman (open or closed)
**Conditions**:
- Three wind triplets + one wind pair

**Code**: `SHOUSUUSHII`

---

### Daisuushii / 대사희 / 大四喜 / 大四喜 ✅

**Han**: Double yakuman (open or closed)
**Conditions**:
- Triplets of all four winds

**Code**: `DAISUUSHII`

---

### Tsuuiisou / 자일색 / 字一色 / 字一色 ✅

**Han**: Yakuman (open or closed)
**Conditions**:
- All tiles are honors (winds + dragons)

**Code**: `TSUUIISOU`

---

### Ryuuiisou / 녹일색 / 緑一色 / 綠一色 ✅

**Han**: Yakuman (open or closed)
**Conditions**:
- All tiles are "green": 2s, 3s, 4s, 6s, 8s, green dragon

**Note**: Some rules require green dragon, others don't

**Code**: `RYUUIISOU`

---

### Chinroutou / 청노두 / 清老頭 / 清老頭 ✅

**Han**: Yakuman (open or closed)
**Conditions**:
- All tiles are terminals (1 and 9 only)
- No honors

**Note**: Extremely rare, requires toitoi or chiitoitsu structure

**Code**: `CHINROUTOU`

---

### Chuuren Poutou / 구련보등 / 九蓮宝燈 / 九蓮寶燈 🔒

**Han**: Yakuman (closed only)
**Conditions**:
- 1112345678999 in single suit + one more tile of same suit
- Pure nine gates pattern

**Junsei Chuuren**: Double yakuman if waiting on any tile of the suit (9-sided wait)

**Code**: `CHUUREN`

---

### Suukantsu / 스깡쯔 / 四槓子 / 四槓子 ✅

**Han**: Yakuman (open or closed)
**Conditions**:
- Four kan declarations by one player

**Note**: Most difficult yakuman, usually game ends in draw before completion

**Code**: `SUUKANTSU`

---

### Tenhou / 천화 / 天和 / 天和 🔒

**Han**: Yakuman (dealer only)
**Conditions**:
- Dealer wins on initial deal (before any discards)

**Code**: `TENHOU`

---

### Chiihou / 지화 / 地和 / 地和 🔒

**Han**: Yakuman (non-dealer only)
**Conditions**:
- Non-dealer wins on first draw
- No interruptions (chi/pon/kan) before first draw

**Code**: `CHIIHOU`

## Local Yaku (로컬역 / ローカル役 / 地方役)

### Renhou / 인화 / 人和 / 人和

**Han**: Varies (mangan to yakuman depending on rules)
**Conditions**:
- Non-dealer wins by ron before their first draw
- No interruptions before ron

**MajSoul**: Treated as mangan (5 han)

**Code**: `RENHOU` (if enabled)

---

### Nagashi Mangan / 유국만관 / 流し満貫 / 流局滿貫

**Han**: Mangan (no yaku, special scoring)
**Conditions**:
- At exhaustive draw, all your discards are terminals/honors
- No one called any of your discards

**MajSoul**: Enabled in some room types

**Code**: `NAGASHI_MANGAN` (if enabled)

---

### Paarenchan / 파렌찬 / 八連荘 / 八連莊

**Han**: Yakuman (local)
**Conditions**:
- Dealer wins 8 consecutive hands

**MajSoul**: Not standard

---

### Sanbaikou / 삼배구 / 三倍口 / 三盃口

**Han**: Yakuman (local)
**Conditions**:
- Three identical sequences (same tiles ×3)

**Status**: Extremely rare, not standard in most rules

## Yaku Combination Rules

### Cannot Combine

- **Chiitoitsu** ↔ Iipeikou, Ryanpeikou, Toitoi, Sanankou (different hand structures)
- **Iipeikou** ↔ Ryanpeikou (ryanpeikou overrides)
- **Honitsu** ↔ Chinitsu (chinitsu is stricter, overrides)
- **Chanta** ↔ Junchan (junchan is stricter, overrides)
- **Rinshan** ↔ Ippatsu, Haitei (timing conflicts)
- **Chankan** ↔ Ippatsu, Houtei (timing/method conflicts)

### Commonly Combined

- **Riichi + Tanyao + Pinfu** = "Mentanpin" (3 han minimum, common strong hand)
- **Riichi + Ippatsu + Tsumo** = 3 han (aggressive fast win)
- **Chinitsu + Honitsu** = Invalid (mutually exclusive)
- **Yakuhai (multiple)** = Stackable (e.g., double east = 2 han)
- **Toitoi + Sanankou** = Common combination for triplet hands

## Implementation Guidelines

### Yaku Detection Order

```python
def detect_yaku(hand, win_tile, win_method, game_state):
    """Detect all applicable yaku in priority order."""
    yaku_list = []

    # 1. Check yakuman first (highest priority)
    if check_kokushi(hand):
        return [Yaku.KOKUSHI]  # Single yakuman, stop checking
    if check_suuankou(hand, win_method):
        return [Yaku.SUUANKOU]
    # ... other yakuman

    # 2. Check regular yaku
    if hand.is_closed:
        if check_riichi(game_state):
            yaku_list.append(Yaku.RIICHI)
            if check_ippatsu(game_state):
                yaku_list.append(Yaku.IPPATSU)

    # 3. Check structure-based yaku
    if check_pinfu(hand):
        yaku_list.append(Yaku.PINFU)
    if check_tanyao(hand):
        yaku_list.append(Yaku.TANYAO)

    # 4. Check situation yaku
    if win_method == WinMethod.HAITEI:
        yaku_list.append(Yaku.HAITEI)

    return yaku_list
```

### Minimum Yaku Validation

```python
def can_win(hand, yaku_list, dora_count):
    """Verify minimum yaku requirement."""
    # Must have at least 1 yaku (dora doesn't count)
    if len(yaku_list) == 0:
        return False

    # Yakuman always valid
    if any(y.is_yakuman() for y in yaku_list):
        return True

    # Regular yaku valid
    return len(yaku_list) >= 1
```

### Han Calculation

```python
def calculate_total_han(yaku_list, dora_count, ura_dora_count=0, aka_dora_count=0):
    """Calculate total han including bonuses."""
    base_han = sum(y.han_value(hand.is_closed) for y in yaku_list)
    bonus_han = dora_count + ura_dora_count + aka_dora_count

    total_han = base_han + bonus_han

    # Cap at yakuman level (13 han = counted yakuman)
    if total_han >= 13:
        return 13  # Kazoe yakuman

    return total_han
```

## Testing Requirements

### Unit Tests per Yaku

```python
def test_riichi_detection():
    """Verify riichi yaku detection."""
    hand = create_closed_hand([1m, 1m, 2m, 3m, 4m, 5m, 6m, 7m, 8m, 8m, 8m, 9m, 9m])
    game_state.riichi_declared = True

    yaku = detect_yaku(hand, win_tile=9m, win_method=WinMethod.RON, game_state)

    assert Yaku.RIICHI in yaku
    assert calculate_han(yaku) >= 1

def test_pinfu_detection():
    """Verify pinfu requirements."""
    # Valid pinfu: all sequences, valueless pair, ryanmen wait
    hand = create_closed_hand([2m, 3m, 4m, 5p, 6p, 7p, 1s, 2s, 3s, 6s, 6s, 7s, 8s])
    win_tile = 9s  # Ryanmen wait on 7-8

    yaku = detect_yaku(hand, win_tile, WinMethod.TSUMO, game_state)

    assert Yaku.PINFU in yaku

def test_yakuman_priority():
    """Verify yakuman stops regular yaku detection."""
    hand = create_kokushi_hand()

    yaku = detect_yaku(hand, win_tile, win_method, game_state)

    assert len(yaku) == 1
    assert yaku[0] == Yaku.KOKUSHI
    assert yaku[0].is_yakuman()
```

## References

- Namu Wiki Yaku List: https://namu.wiki/w/리치마작/역
- Japanese Wikipedia: https://ja.wikipedia.org/wiki/麻雀の役一覧
- MajSoul yaku reference (in-game help system)
- World Riichi Championship rules: http://wrc-rr.com/
- See `constitution_term.md` for terminology definitions
- See `constitution_score.md` for han-fu scoring calculations
