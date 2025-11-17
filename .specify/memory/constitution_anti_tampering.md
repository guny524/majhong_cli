# Anti-Tampering System: Wall Integrity Verification

**Purpose**: Defines SHA-256 based cryptographic verification system for preventing wall manipulation in online Mahjong multiplayer games.
**Referenced by**: `.specify/memory/constitution.md`
**Last Updated**: 2025-11-17
**Reference Implementation**: MajSoul (雀魂/작혼) online Mahjong platform

## Problem Statement

### Trust Issues in Online Multiplayer Mahjong

**Security concerns**:
1. **Server manipulation**: Server could generate favorable walls for specific players
2. **Man-in-the-middle attacks**: Wall data could be intercepted and modified during transmission
3. **Post-game disputes**: Players cannot verify wall fairness after game completion
4. **Cheating detection**: Need proof that wall was pre-determined, not generated adaptively

**What SHA-256 verification CANNOT prevent**:
- Server choosing specific favorable wall at game start (pre-knowledge attack)
- Server or client compromised before game begins
- Does NOT prove fairness of wall generation, ONLY proves wall integrity (no mid-game tampering)

**What SHA-256 verification CAN prove**:
- Wall composition did not change between game start and game end
- Salt value was not modified during game
- All players/observers can independently verify wall integrity post-game

## Cryptographic Verification System

### Overview

**Two-stage hash commitment scheme**:
1. **H(wall_code + salt)**: Commits to specific wall + salt combination
2. **H(salt)**: Commits to specific salt value independently

**Verification flow**:
```
[Game Start]
  Server generates: wall_code (136 tiles), salt (random string)
  Server publishes: H(wall_code + salt), H(salt)
  Server keeps secret: wall_code, salt

[Game in Progress]
  Players can verify hashes remain unchanged (click dora indicator)
  No knowledge of actual wall composition or salt value

[Game End]
  Server reveals: wall_code, salt
  All parties verify:
    - H(wall_code + salt) matches published hash
    - H(salt) matches published hash
    - Wall_code contains exactly 136 valid tiles
```

### Why Two Hashes?

**Single hash H(wall_code) is insufficient**:
- Problem: Attacker could pre-compute rainbow table for all possible wall permutations
- 136! / (4!^34) ≈ 10^157 combinations, but many combinations are functionally equivalent
- Reduced search space makes brute-force attacks theoretically feasible with massive computation

**H(wall_code + salt) solves rainbow table attack**:
- Salt adds unpredictable random data → makes pre-computation infeasible
- Even if attacker knows wall_code candidates, cannot verify without salt

**H(salt) provides salt integrity proof**:
- Prevents server from changing salt after seeing game outcome
- Without H(salt), server could claim different salt to justify observed wall
- H(salt) commitment binds server to specific salt before game starts

### Hash Algorithm: SHA-256

**Properties**:
- **Pre-image resistance**: Given H(x), computationally infeasible to find x
- **Collision resistance**: Infeasible to find x ≠ y where H(x) = H(y)
- **Avalanche effect**: Small change in input causes completely different hash output

**Example from MajSoul** (from sha1.png and sha2.png):
```
Game wall hash: H(wall_code + salt)
  f3708433572ed2f15009e55ca4af4c09410170e94e63b33062c4feed50b3e7b3

Salt hash: H(salt)
  c2741b03516377bb7f276840c179844e4ad321c41503e5678cdc116a18a39f36
```

## Wall Code Format (패산코드열)

### Tile Notation (from README2.md and sha2.png)

**Standard encoding**:
```
Man (萬/만):    1m 2m 3m 4m 5m 6m 7m 8m 9m
Pin (筒/통):    1p 2p 3p 4p 5p 6p 7p 8p 9p
Sou (索/삭):    1s 2s 3s 4s 5s 6s 7s 8s 9s
Winds:          1z (East/동), 2z (South/남), 3z (West/서), 4z (North/북)
Dragons:        5z (White/백), 6z (Green/발), 7z (Red/중)
Red fives:      0m 0p 0s (optional red-five variants)
```

**Wall code string format**:
- Concatenate all 136 tile codes in wall order
- Example: `"1m1m1m1m2m2m2m2m3m3m...7z7z7z"`
- Length: 272 characters (136 tiles × 2 chars each)

**Wall composition validation**:
```python
def validate_wall_code(wall_code: str) -> bool:
    """Verify wall contains exactly correct tile counts."""
    tiles = [wall_code[i:i+2] for i in range(0, len(wall_code), 2)]

    if len(tiles) != 136:
        return False

    # Count each tile type
    counts = Counter(tiles)

    # Standard tiles: 4 of each
    standard_tiles = (
        [f"{n}m" for n in range(1, 10)] +
        [f"{n}p" for n in range(1, 10)] +
        [f"{n}s" for n in range(1, 10)] +
        [f"{n}z" for n in range(1, 8)]
    )

    for tile in standard_tiles:
        # Special case: if red fives enabled, 5m/5p/5s only 3× (one is 0m/0p/0s)
        if tile in ["5m", "5p", "5s"] and "0" + tile[1] in counts:
            expected = 3
        else:
            expected = 4

        if counts.get(tile, 0) != expected:
            return False

    # Verify red fives (if present) are exactly 1 each
    for red_five in ["0m", "0p", "0s"]:
        if red_five in counts and counts[red_five] != 1:
            return False

    return True
```

## Implementation Specification

### Server-Side: Game Initialization

```python
import hashlib
import secrets
from typing import List, Tuple

class WallIntegrityManager:
    def __init__(self):
        self.wall_code: str = ""
        self.salt: str = ""
        self.wall_hash: str = ""
        self.salt_hash: str = ""

    def generate_game_wall(self) -> Tuple[str, str]:
        """
        Generate wall and compute integrity hashes.

        Returns:
            (wall_hash, salt_hash): Published at game start
        """
        # 1. Generate random wall (shuffle tiles)
        tiles = self._generate_tile_set()
        random.shuffle(tiles)
        self.wall_code = "".join(tiles)

        # 2. Generate cryptographically secure random salt
        self.salt = secrets.token_hex(32)  # 64 hex characters = 256 bits

        # 3. Compute SHA-256 hashes
        self.wall_hash = hashlib.sha256(
            (self.wall_code + self.salt).encode('utf-8')
        ).hexdigest()

        self.salt_hash = hashlib.sha256(
            self.salt.encode('utf-8')
        ).hexdigest()

        # 4. Return public hashes (keep wall_code and salt secret)
        return (self.wall_hash, self.salt_hash)

    def _generate_tile_set(self) -> List[str]:
        """Generate standard 136-tile set."""
        tiles = []

        # Number tiles: 1-9 man/pin/sou × 4 each
        for suit in ['m', 'p', 's']:
            for num in range(1, 10):
                tiles.extend([f"{num}{suit}"] * 4)

        # Honor tiles: winds + dragons × 4 each
        for z_tile in range(1, 8):
            tiles.extend([f"{z_tile}z"] * 4)

        # Optional: Replace one 5 of each suit with red five
        if RED_FIVES_ENABLED:
            for suit in ['m', 'p', 's']:
                idx = tiles.index(f"5{suit}")
                tiles[idx] = f"0{suit}"

        return tiles

    def reveal_wall_data(self) -> Tuple[str, str]:
        """
        Reveal wall and salt at game end.

        Returns:
            (wall_code, salt): Secret data for verification
        """
        return (self.wall_code, self.salt)
```

### Client-Side: Verification

```python
class WallVerifier:
    @staticmethod
    def verify_integrity(
        wall_code: str,
        salt: str,
        published_wall_hash: str,
        published_salt_hash: str
    ) -> bool:
        """
        Verify wall integrity after game end.

        Args:
            wall_code: Revealed 136-tile wall string
            salt: Revealed salt value
            published_wall_hash: H(wall_code + salt) from game start
            published_salt_hash: H(salt) from game start

        Returns:
            True if wall integrity verified, False if tampering detected
        """
        # 1. Recompute hashes
        computed_wall_hash = hashlib.sha256(
            (wall_code + salt).encode('utf-8')
        ).hexdigest()

        computed_salt_hash = hashlib.sha256(
            salt.encode('utf-8')
        ).hexdigest()

        # 2. Compare with published hashes
        wall_match = computed_wall_hash == published_wall_hash
        salt_match = computed_salt_hash == published_salt_hash

        # 3. Validate wall composition
        valid_composition = validate_wall_code(wall_code)

        return wall_match and salt_match and valid_composition

    @staticmethod
    def check_hash_during_game(
        published_wall_hash: str,
        published_salt_hash: str,
        server_current_wall_hash: str,
        server_current_salt_hash: str
    ) -> bool:
        """
        Mid-game verification: Ensure server hasn't changed hashes.

        Players can query server during game to verify hashes remain unchanged.
        """
        return (
            published_wall_hash == server_current_wall_hash and
            published_salt_hash == server_current_salt_hash
        )
```

### Game Flow Integration

```python
class MahjongGame:
    def __init__(self):
        self.integrity_mgr = WallIntegrityManager()
        self.published_hashes = None

    def start_game(self):
        """Initialize game with integrity verification."""
        # 1. Generate wall and compute hashes
        wall_hash, salt_hash = self.integrity_mgr.generate_game_wall()
        self.published_hashes = (wall_hash, salt_hash)

        # 2. Broadcast hashes to all players
        self.broadcast_to_players({
            "type": "wall_integrity_hashes",
            "wall_hash": wall_hash,
            "salt_hash": salt_hash,
            "timestamp": time.time()
        })

        # 3. Log hashes to immutable audit log
        self.audit_log.record({
            "event": "game_start",
            "game_id": self.game_id,
            "wall_hash": wall_hash,
            "salt_hash": salt_hash
        })

        # 4. Continue with normal game setup
        self.distribute_tiles()

    def handle_mid_game_verification_request(self, player_id: int):
        """
        Handle player clicking dora indicator to check hashes.

        Corresponds to MajSoul feature where players can verify hashes mid-game.
        """
        return {
            "type": "wall_integrity_status",
            "wall_hash": self.published_hashes[0],
            "salt_hash": self.published_hashes[1],
            "message": "Wall integrity hashes unchanged since game start"
        }

    def end_game(self):
        """Reveal wall data for post-game verification."""
        wall_code, salt = self.integrity_mgr.reveal_wall_data()

        # 1. Broadcast revealed data
        self.broadcast_to_players({
            "type": "wall_data_reveal",
            "wall_code": wall_code,
            "salt": salt,
            "published_wall_hash": self.published_hashes[0],
            "published_salt_hash": self.published_hashes[1]
        })

        # 2. Log to audit trail
        self.audit_log.record({
            "event": "game_end",
            "game_id": self.game_id,
            "wall_code": wall_code,
            "salt": salt
        })

        # 3. All clients verify independently
        # (handled client-side with WallVerifier.verify_integrity)
```

## User Interface Elements

### In-Game Hash Display (from sha1.png)

**MajSoul implementation**: Click on dora indicator (left side of dead wall) during game

**Display shows**:
```
솔트값 추가 SHA256 검증 방식:
1. 대국 중, 좌측 상단 도라패의 위치를 클릭하여 본 대국의 '패산 코드열+솔트값' 및
   '솔트값'과 대응하는 SHA256코드를 확인할 수 있으며, 솔트 코드열은 매번 고정되어
   생성되는 랜덤 코드열입니다.

2. 대국 종료 후, 패보에서 본 대국의 '패산 코드열+솔트값' 및 '솔트값'이 SHA256코드
   해당 대국의 것과 일치하는지 확인하여 패산에 변조의 사실이 없음을 확인하실 수 있습니다.

*2023년 5월 24일 점검 전, 모든 패보는 MD5 암호화를 사용합니다.

MD5 검증 방식:
[...MD5 legacy information...]
```

**Key information displayed**:
- H(wall_code + salt): `f3708433572ed2f15009e55ca4af4c09410170e94e63b33062c4feed50b3e7b3`
- H(salt): `c2741b03516377bb7f276840c179844e4ad321c41503e5678cdc116a18a39f36`

### Post-Game Verification UI (from sha2.png)

**Display after game completion**:
```
[작혼: 리치 마작]은 100% 랜덤 시스템으로 운영되며, 대국 시작 시, 패산의 순서는
이미 고정되어 있어 특정 패를 나오게 하거나, 중간에 패산을 조작할 수 없습니다.

이를 증명하기 위해 암호화 검증 시스템을 도입하였습니다.
(2023년 5월 24일까지 MD5 암호화 사용, 2023년 5월 24일 SHA256 암호화로
업데이트 되었으며, 2024년 2월 28일 솔트값이 추가된 SHA256 암호화로
업데이트 되었습니다.)

각각의 패마다 두 자릿수 인코딩을 사용하고 있으며, 1만~9만은 1m~9m, 1통~9통은
1p~9p, 1삭~9삭은 1s~9s, 동남서북백발중은 1z~7z, 적색 5만, 적색 5통, 적색 5삭은
각각 0m, 0p, 0s로 인코딩되어있습니다.

패산은 순서에 따라 코드열을 확인하실 수 있습니다. 본 암호화 검증 시스템은 코드
열과역에 기초하여 생성되었습니다.

*2024년 2월 28일 후, 모든 패보는 솔트값이 추가된 SHA256 암호화를 사용하며,
암호화 알고리즘의 보안성을 향상하였습니다.
```

**Revealed data**:
- Wall code string (136 tiles × 2 chars)
- Salt value
- Original published hashes for comparison

## Security Considerations

### Attack Vectors and Mitigations

#### 1. Rainbow Table Attack

**Attack**: Pre-compute H(wall_code) for all possible walls
**Mitigation**: Salt concatenation makes pre-computation infeasible
**Status**: ✅ Mitigated

#### 2. Server Pre-Knowledge Attack

**Attack**: Server generates many walls, chooses favorable one before publishing hash
**Mitigation**: ❌ NOT PREVENTED by hash system
**Note**: This is a fundamental limitation. Hash proves integrity, not fairness of initial generation

**Potential mitigations** (not in MajSoul):
- Client-contributed randomness (each player provides seed)
- Commit-reveal protocol with time delays
- Trusted third-party random number generation

#### 3. Salt Manipulation Attack

**Attack**: Server changes salt after seeing game outcome to justify observed wall
**Mitigation**: ✅ H(salt) commitment prevents this
**Verification**: H(salt) must match published value

#### 4. Hash Substitution Mid-Game

**Attack**: Server silently changes published hashes during game
**Mitigation**: ✅ Players can verify hashes mid-game (click dora indicator)
**Best practice**: Clients should cache published hashes locally, not trust server's current claim

#### 5. Replay Attack

**Attack**: Server reuses same wall/salt across multiple games
**Mitigation**: Include game_id or timestamp in hash computation
**Enhancement**: `H(game_id || timestamp || wall_code || salt)`

### Implementation Best Practices

```python
class SecureWallIntegrityManager(WallIntegrityManager):
    """Enhanced version with additional security measures."""

    def generate_game_wall(self, game_id: str, timestamp: int) -> Tuple[str, str]:
        """Include game_id and timestamp in hash to prevent replay attacks."""
        tiles = self._generate_tile_set()
        random.shuffle(tiles)
        self.wall_code = "".join(tiles)

        self.salt = secrets.token_hex(32)

        # Enhanced hash: include game context
        wall_data = f"{game_id}|{timestamp}|{self.wall_code}|{self.salt}"
        self.wall_hash = hashlib.sha256(wall_data.encode('utf-8')).hexdigest()

        salt_data = f"{game_id}|{timestamp}|{self.salt}"
        self.salt_hash = hashlib.sha256(salt_data.encode('utf-8')).hexdigest()

        return (self.wall_hash, self.salt_hash)

    def verify_with_context(
        self, game_id: str, timestamp: int, wall_code: str, salt: str,
        published_wall_hash: str, published_salt_hash: str
    ) -> bool:
        """Verify including game context."""
        wall_data = f"{game_id}|{timestamp}|{wall_code}|{salt}"
        computed_wall_hash = hashlib.sha256(wall_data.encode('utf-8')).hexdigest()

        salt_data = f"{game_id}|{timestamp}|{salt}"
        computed_salt_hash = hashlib.sha256(salt_data.encode('utf-8')).hexdigest()

        return (
            computed_wall_hash == published_wall_hash and
            computed_salt_hash == published_salt_hash and
            validate_wall_code(wall_code)
        )
```

## Testing Requirements

### Unit Tests

```python
def test_hash_determinism():
    """Same input always produces same hash."""
    wall_code = "1m1m1m1m2m2m..."
    salt = "abc123"

    hash1 = hashlib.sha256((wall_code + salt).encode()).hexdigest()
    hash2 = hashlib.sha256((wall_code + salt).encode()).hexdigest()

    assert hash1 == hash2

def test_hash_collision_resistance():
    """Different inputs produce different hashes."""
    wall_code1 = "1m1m1m1m2m2m..."
    wall_code2 = "1m1m1m2m1m2m..."  # Different order
    salt = "abc123"

    hash1 = hashlib.sha256((wall_code1 + salt).encode()).hexdigest()
    hash2 = hashlib.sha256((wall_code2 + salt).encode()).hexdigest()

    assert hash1 != hash2

def test_salt_prevents_precomputation():
    """Same wall_code with different salts produces different hashes."""
    wall_code = "1m1m1m1m2m2m..."

    hash1 = hashlib.sha256((wall_code + "salt1").encode()).hexdigest()
    hash2 = hashlib.sha256((wall_code + "salt2").encode()).hexdigest()

    assert hash1 != hash2

def test_wall_code_validation():
    """Validate wall contains exactly correct tiles."""
    valid_wall = generate_valid_wall_code()
    assert validate_wall_code(valid_wall) == True

    invalid_wall = valid_wall[:-2] + "XX"  # Replace last tile with invalid
    assert validate_wall_code(invalid_wall) == False

def test_end_to_end_verification():
    """Full workflow: generate → publish → reveal → verify."""
    mgr = WallIntegrityManager()

    # Game start: generate and publish hashes
    wall_hash, salt_hash = mgr.generate_game_wall()

    # Game end: reveal data
    wall_code, salt = mgr.reveal_wall_data()

    # Client verification
    verifier = WallVerifier()
    result = verifier.verify_integrity(wall_code, salt, wall_hash, salt_hash)

    assert result == True
```

### Integration Tests

```python
def test_tampering_detection():
    """Verify tampering is detected."""
    mgr = WallIntegrityManager()
    wall_hash, salt_hash = mgr.generate_game_wall()

    # Attacker modifies wall_code
    original_wall, salt = mgr.reveal_wall_data()
    tampered_wall = original_wall[:10] + "9z9z" + original_wall[14:]  # Invalid modification

    verifier = WallVerifier()
    result = verifier.verify_integrity(tampered_wall, salt, wall_hash, salt_hash)

    assert result == False  # Tampering detected

def test_mid_game_hash_verification():
    """Players can verify hashes during game."""
    game = MahjongGame()
    game.start_game()

    # Mid-game verification request
    response = game.handle_mid_game_verification_request(player_id=0)

    assert response["wall_hash"] == game.published_hashes[0]
    assert response["salt_hash"] == game.published_hashes[1]
```

## Migration from MD5 to SHA-256

**Historical note** (from sha1.png and sha2.png):
- **Before 2023-05-24**: MD5 hashing (deprecated due to collision vulnerabilities)
- **2023-05-24 to 2024-02-27**: SHA-256 without salt (vulnerable to rainbow tables)
- **After 2024-02-28**: SHA-256 with salt (current standard)

**For majhong_cli implementation**:
- Use SHA-256 with salt from day one
- Do NOT implement MD5 mode (obsolete and insecure)
- Consider including game_id and timestamp for enhanced security

## References

- MajSoul (雀魂/작혼) integrity verification system: https://mahjongsoul.com/
- SHA-256 specification: FIPS PUB 180-4
- Cryptographic commitment schemes: "Secure Hash Standard (SHS)", NIST
- Images: `sha1.png`, `sha2.png` from MajSoul in-game verification UI
