# Mahjong Wall Construction Procedures

**Purpose**: Defines physical and digital wall building procedures for offline/online Mahjong implementation.
**Referenced by**: `.specify/memory/constitution.md`, `constitution_majhong_rule.md`
**Last Updated**: 2025-11-17
**Video Reference**: https://www.youtube.com/watch?v=Q-IIkMKjR_U

## Physical Wall Construction (Offline Mahjong)

### 1. Tile Shuffling

**Manual shuffle procedure**:
1. Place all 136 tiles face-down on table
2. Mix tiles thoroughly by moving them in circular motions
3. Flip tiles multiple times during shuffle to ensure randomness
4. Continue until tiles are well-distributed

**Note**: Automatic mahjong tables perform this step mechanically.

### 2. Wall Building - Each Player's Segment

**Per-player wall specification**:
- **Tiles per segment**: 17 tiles × 2 layers = 34 tiles total
- **Layout**: 2 rows stacked vertically, 17 tiles long

**Stacking technique** (from video transcript):
1. Arrange 17 face-down tiles in a single row before you
2. Stack a second row of 17 tiles on top (aligned vertically)
3. Grasp wall segment carefully (like holding an egg gently with pinky finger support)
4. Lift and place as a standing wall segment in front of your position

**Result**: Four players create four wall segments → complete rectangular wall (136 tiles: 68 tiles × 2 layers)

### 3. Wall Break Point Determination

**Procedure** (performed by dealer after seat determination):

1. **Dealer rolls dice** (2 standard six-sided dice)
2. **Count players counter-clockwise** from dealer (including dealer = position 1)
   - Dice sum determines which player's wall to break
   - Example: Dealer rolls 8 → count 8 positions counter-clockwise from dealer
3. **Break point within player's wall**:
   - From selected player's wall, count (dice sum) × 2 tiles **from the RIGHT**
   - Remove this section, leaving 7 stacks (14 tiles) for the dead wall

**Example** (from video):
- Dealer rolls 8
- Count to 8th position counter-clockwise (lands on a specific player)
- From that player's wall RIGHT side, remove 8 stacks × 2 tiles = 16 tiles
- Attach removed section to right-side neighbor's wall (extends their wall)
- Remaining 7 stacks (14 tiles) become dead wall foundation

### 4. Dead Wall (王牌/Wanpai) Formation

**Dead wall composition** (14 tiles from left side of break point):

```
[Rinshan: 4 tiles]  [Dora indicators: 5]  [Ura-dora: 5]
[Top layer:    ●●●●●●●]
[Bottom layer: ●●●●●●●]
       ↑ 3rd tile from left flipped = Dora indicator
```

**Dora indicator reveal**:
- Flip the **3rd tile from the left** (top layer) face-up
- This tile indicates dora (next tile in sequence = actual dora)

**Rinshan tiles (嶺上牌/영상패)**:
- Leftmost 4 tiles (2 stacks)
- Drawn when declaring kan (槓)
- Winning on rinshan draw = "Rinshan Kaihou" (嶺上開花/영상개화) yaku

**Dead wall maintenance**:
- MUST always maintain exactly 14 tiles
- When rinshan tile drawn → append last tile from live wall to dead wall

### 5. Live Wall (活牌) Drawing Direction

**Initial drawing point**:
- Starts immediately after dead wall (to the left of the 14-tile dead wall section)
- Draw direction: **CLOCKWISE** around the table

**Drawing pattern**:
- All players draw from same continuous wall
- Not restricted to individual wall segments
- Wall functions as a circular queue

**Wall depletion marker**:
- When only 14 tiles remain (dead wall), game ends in exhaustive draw (流局/ryūkyoku)

## Tile Distribution (配牌/Haipai)

### Initial Hand Dealing

**Distribution sequence** (from wall break point, clockwise):

1. **First three rounds**: 4 tiles per player, **counter-clockwise** player order
   - Dealer takes 4 → South takes 4 → West takes 4 → North takes 4
   - Repeat 3 times → Each player has 12 tiles

2. **Final distribution**:
   - Dealer takes 2 tiles (1 middle tile + 1 jump tile, total 14 tiles)
   - Other players take 1 tile each (total 13 tiles each)

**Dealer advantage**: Dealer starts with 14 tiles (ready to discard), others have 13 tiles (must draw first)

## Hand Organization Techniques

### Efficient Tile Arrangement (from video)

**Method 1 - Scraping**:
1. Hold two tiles as "scrapers"
2. Slide tiles together to compress gaps
3. Reorganize by suit and number efficiently

**Method 2 - Attachment technique**:
- Instead of sorting in place (slow)
- Pick tiles and attach them to organized sections
- Attach → Attach → Scrape → faster organization

### Tsumo (Self-Draw) Tile Handling

**Critical rule**: Keep drawn tile SEPARATE from hand

**Reason**: Scoring and legality depend on which tile was just drawn
- Drawn tile placed at rightmost position (distinct from hand)
- When declaring tsumo win, drawn tile must be identifiable
- When discarding, move tile to discard area

**DO NOT**:
- Mix tsumo tile into hand immediately
- Place tsumo tile in center of hand
- Discard directly without showing drawn tile position

## Point Stick Distribution

### Starting Points: 25,000 per player

**Standard composition**:
- **Red (10,000 points)**: 1 stick
- **Yellow (5,000 points)**: 2 sticks
- **Blue (1,000 points)**: 4 sticks
- **Green (500 points)**: 1 stick (often substituted)
- **White (100 points)**: 5 sticks

### Practical 500-Point Handling

**Issue**: Physical sets often lack 500-point sticks

**Solution**:
- Mark a 100-point stick with "500" notation
- Use as substitute rather than using 5× 100-point sticks
- Reduces clutter and counting errors

## Digital Implementation Considerations

### Wall State Machine

```python
class MahjongWall:
    def __init__(self, tiles: List[Tile]):
        """Initialize wall with shuffled 136 tiles."""
        self.tiles = tiles  # Pre-shuffled
        self.break_point = 0
        self.dead_wall_start = 0
        self.live_wall_tiles = 122  # 136 - 14
        self.dead_wall = DeadWall()

    def determine_break(self, dice_sum: int, dealer_position: int) -> int:
        """Calculate wall break point from dice roll."""
        # Dice sum determines which player's wall (counter-clockwise)
        player_offset = (dice_sum - 1) % 4
        target_player = (dealer_position + player_offset) % 4

        # Within player's wall, dice_sum × 2 tiles from RIGHT
        tiles_from_right = dice_sum * 2
        segment_start = target_player * 34  # Each player has 34-tile segment
        break_index = segment_start + (34 - tiles_from_right)

        self.break_point = break_index
        self.dead_wall_start = break_index
        return break_index

    def distribute_initial_hands(self) -> Dict[int, List[Tile]]:
        """Distribute 14/13/13/13 tiles to dealer/players."""
        hands = {0: [], 1: [], 2: [], 3: []}
        draw_index = self.dead_wall_start + 14  # Skip dead wall

        # Three rounds of 4 tiles each
        for round in range(3):
            for player in range(4):
                for _ in range(4):
                    hands[player].append(self.tiles[draw_index])
                    draw_index += 1

        # Dealer +2, others +1
        hands[0].extend([self.tiles[draw_index], self.tiles[draw_index + 1]])
        draw_index += 2
        for player in [1, 2, 3]:
            hands[player].append(self.tiles[draw_index])
            draw_index += 1

        self.live_wall_tiles = 122 - sum(len(h) for h in hands.values()) + 14  # Dealer starts with 14
        return hands

    def draw_tile(self) -> Tile | None:
        """Draw one tile from live wall (clockwise)."""
        if self.live_wall_tiles <= 0:
            return None  # Wall exhausted → ryūkyoku
        tile = self.tiles[self.current_draw_index]
        self.current_draw_index += 1
        self.live_wall_tiles -= 1
        return tile

    def draw_rinshan(self) -> Tile:
        """Draw replacement tile for kan declaration."""
        rinshan_tile = self.dead_wall.pop_rinshan()

        # Replenish dead wall from live wall end
        last_live_tile = self.tiles[self.live_wall_end_index]
        self.dead_wall.append_to_end(last_live_tile)
        self.live_wall_end_index -= 1

        return rinshan_tile

    def flip_next_dora(self):
        """Reveal next dora indicator (on kan declaration)."""
        self.dead_wall.flip_dora_indicator()
```

### Dead Wall Management

```python
class DeadWall:
    def __init__(self, tiles: List[Tile]):
        """Initialize with 14 tiles from break point."""
        self.rinshan = tiles[0:4]          # Leftmost 4 tiles
        self.dora_indicators = tiles[4:9]   # Next 5 tiles (top layer)
        self.ura_indicators = tiles[9:14]   # Next 5 tiles (bottom layer)
        self.dora_visible = 1               # 3rd tile from left flipped initially

    def pop_rinshan(self) -> Tile:
        """Remove and return leftmost rinshan tile."""
        return self.rinshan.pop(0)

    def append_to_end(self, tile: Tile):
        """Add tile to rightmost position (replenish from live wall)."""
        self.rinshan.append(tile)  # Maintains 14-tile count

    def flip_dora_indicator(self):
        """Reveal next dora indicator (max 5)."""
        if self.dora_visible < 5:
            self.dora_visible += 1

    def get_visible_dora_indicators(self) -> List[Tile]:
        """Return currently revealed dora indicators."""
        return self.dora_indicators[:self.dora_visible]

    def get_ura_dora_indicators(self) -> List[Tile]:
        """Return ura-dora indicators (only after riichi win)."""
        return self.ura_indicators[:self.dora_visible]
```

### Validation Rules

**MUST enforce in digital implementation**:
1. Dead wall ALWAYS exactly 14 tiles
2. Live wall starts at 122 tiles (136 - 14)
3. Game ends when live wall reaches 0 (exhaustive draw)
4. Dora indicator count ≤ 5 (one initial + four possible kan declarations)
5. Rinshan draw only valid during kan declaration
6. Wall break calculation deterministic from dice + dealer position

## Testing Requirements

### Unit Tests for Wall Construction

```python
def test_wall_break_deterministic():
    """Verify same dice roll + dealer produces same break point."""
    wall = MahjongWall(shuffled_tiles)
    break1 = wall.determine_break(dice_sum=8, dealer_position=0)

    wall2 = MahjongWall(same_shuffled_tiles)
    break2 = wall2.determine_break(dice_sum=8, dealer_position=0)

    assert break1 == break2

def test_dead_wall_maintains_14_tiles():
    """Verify dead wall always 14 tiles after rinshan draws."""
    dead_wall = DeadWall(initial_14_tiles)
    assert len(dead_wall.all_tiles()) == 14

    rinshan = dead_wall.pop_rinshan()
    assert len(dead_wall.all_tiles()) == 13  # Temporarily

    dead_wall.append_to_end(replacement_tile)
    assert len(dead_wall.all_tiles()) == 14  # Restored

def test_initial_distribution_correct_counts():
    """Verify dealer gets 14, others get 13 tiles."""
    wall = MahjongWall(shuffled_tiles)
    wall.determine_break(dice_sum=7, dealer_position=0)
    hands = wall.distribute_initial_hands()

    assert len(hands[0]) == 14  # Dealer
    assert len(hands[1]) == 13
    assert len(hands[2]) == 13
    assert len(hands[3]) == 13
    assert wall.live_wall_tiles == 122 - (14 + 13*3)
```

### Integration Tests

```python
def test_full_game_wall_lifecycle():
    """Simulate complete game from wall construction to exhaustion."""
    game = MahjongGame()
    game.setup()  # Shuffle, break, distribute

    turns = 0
    while game.wall.live_wall_tiles > 0:
        current_player = game.current_player
        tile = game.wall.draw_tile()
        assert tile is not None

        game.players[current_player].add_tile(tile)
        discard = game.players[current_player].choose_discard()
        game.discard(discard)

        turns += 1

    # Game must end when wall exhausted
    assert game.wall.live_wall_tiles == 0
    assert game.state == GameState.EXHAUSTIVE_DRAW
```

## Physical Mahjong Best Practices

### From Video Transcript

1. **Tile handling efficiency**:
   - Use two-tile scraper technique for faster organization
   - Attach tiles to sorted sections rather than sorting in-place
   - Reduces setup time significantly

2. **Tsumo tile separation**:
   - ALWAYS keep drawn tile at rightmost position
   - Critical for scoring verification
   - Prevents disputes about winning tile source

3. **Point stick management**:
   - Use marked 500-point substitute to avoid 10× 100-point clutter
   - Reduces counting errors during payment

4. **Wall lifting technique**:
   - Support with pinky finger underneath
   - Gentle grip (like holding an egg)
   - Prevents wall collapse during placement

## References

- Video tutorial: https://www.youtube.com/watch?v=Q-IIkMKjR_U (1급 천재 channel)
- Riichi Mahjong rules: `constitution_majhong_rule.md`
- Physical mahjong set standards: Japanese Mahjong Federation
