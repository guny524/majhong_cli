package engine

import (
	"testing"
)

func TestNewTile(t *testing.T) {
	tile := NewTile(Man, 3)
	if tile.Suit != Man {
		t.Errorf("Expected suit Man, got %v", tile.Suit)
	}
	if tile.Rank != 3 {
		t.Errorf("Expected rank 3, got %d", tile.Rank)
	}
	if tile.IsAkaDora {
		t.Error("Expected regular tile, got aka dora")
	}
}

func TestNewAkaDoraTile(t *testing.T) {
	tile := NewAkaDoraTile(Pin)
	if tile.Rank != 5 {
		t.Errorf("Expected rank 5 for aka dora, got %d", tile.Rank)
	}
	if !tile.IsAkaDora {
		t.Error("Expected aka dora flag to be true")
	}
}

func TestTileString(t *testing.T) {
	tests := []struct {
		tile     Tile
		expected string
	}{
		{NewTile(Man, 3), "3m"},
		{NewTile(Pin, 7), "7p"},
		{NewTile(Sou, 1), "1s"},
		{NewTile(Ji, 1), "1z"}, // East
		{NewTile(Ji, 5), "5z"}, // White
		{NewAkaDoraTile(Man), "0m"},
		{NewAkaDoraTile(Pin), "0p"},
		{NewAkaDoraTile(Sou), "0s"},
	}

	for _, tt := range tests {
		result := tt.tile.String()
		if result != tt.expected {
			t.Errorf("Tile %+v: expected %s, got %s", tt.tile, tt.expected, result)
		}
	}
}

func TestParseTile(t *testing.T) {
	tests := []struct {
		notation    string
		expectedErr bool
		expected    Tile
	}{
		{"3m", false, NewTile(Man, 3)},
		{"7p", false, NewTile(Pin, 7)},
		{"1s", false, NewTile(Sou, 1)},
		{"1z", false, NewTile(Ji, 1)},
		{"0m", false, NewAkaDoraTile(Man)},
		{"0p", false, NewAkaDoraTile(Pin)},
		{"0s", false, NewAkaDoraTile(Sou)},
		{"10m", true, Tile{}},  // Invalid rank
		{"0z", true, Tile{}},   // Aka dora cannot be honors
		{"5x", true, Tile{}},   // Invalid suit
		{"m", true, Tile{}},    // Missing rank
		{"", true, Tile{}},     // Empty
	}

	for _, tt := range tests {
		result, err := ParseTile(tt.notation)
		if tt.expectedErr {
			if err == nil {
				t.Errorf("ParseTile(%s): expected error, got nil", tt.notation)
			}
		} else {
			if err != nil {
				t.Errorf("ParseTile(%s): unexpected error: %v", tt.notation, err)
			}
			if !result.Equals(tt.expected) {
				t.Errorf("ParseTile(%s): expected %+v, got %+v", tt.notation, tt.expected, result)
			}
		}
	}
}

func TestTileEquals(t *testing.T) {
	tile1 := NewTile(Man, 3)
	tile2 := NewTile(Man, 3)
	tile3 := NewTile(Pin, 3)
	aka1 := NewAkaDoraTile(Man)
	aka2 := NewAkaDoraTile(Man)
	regular5m := NewTile(Man, 5)

	if !tile1.Equals(tile2) {
		t.Error("Same tiles should be equal")
	}
	if tile1.Equals(tile3) {
		t.Error("Different suits should not be equal")
	}
	if !aka1.Equals(aka2) {
		t.Error("Same aka dora should be equal")
	}
	if aka1.Equals(regular5m) {
		t.Error("Aka dora and regular 5m should not be equal")
	}
}

func TestTileSameType(t *testing.T) {
	aka5m := NewAkaDoraTile(Man)
	regular5m := NewTile(Man, 5)
	other3m := NewTile(Man, 3)

	if !aka5m.SameType(regular5m) {
		t.Error("Aka dora 5m and regular 5m should be same type")
	}
	if aka5m.SameType(other3m) {
		t.Error("5m and 3m should not be same type")
	}
}

func TestTileIsTerminal(t *testing.T) {
	tests := []struct {
		tile     Tile
		expected bool
	}{
		{NewTile(Man, 1), true},
		{NewTile(Man, 9), true},
		{NewTile(Pin, 5), false},
		{NewTile(Ji, 1), false}, // Honors are not terminals
	}

	for _, tt := range tests {
		result := tt.tile.IsTerminal()
		if result != tt.expected {
			t.Errorf("Tile %s: IsTerminal() expected %v, got %v", tt.tile, tt.expected, result)
		}
	}
}

func TestTileIsHonor(t *testing.T) {
	tests := []struct {
		tile     Tile
		expected bool
	}{
		{NewTile(Ji, 1), true},
		{NewTile(Ji, 5), true},
		{NewTile(Man, 5), false},
	}

	for _, tt := range tests {
		result := tt.tile.IsHonor()
		if result != tt.expected {
			t.Errorf("Tile %s: IsHonor() expected %v, got %v", tt.tile, tt.expected, result)
		}
	}
}

func TestTileIsSimple(t *testing.T) {
	tests := []struct {
		tile     Tile
		expected bool
	}{
		{NewTile(Man, 2), true},
		{NewTile(Pin, 8), true},
		{NewTile(Sou, 1), false},
		{NewTile(Man, 9), false},
		{NewTile(Ji, 5), false},
	}

	for _, tt := range tests {
		result := tt.tile.IsSimple()
		if result != tt.expected {
			t.Errorf("Tile %s: IsSimple() expected %v, got %v", tt.tile, tt.expected, result)
		}
	}
}

func TestTileNextDora(t *testing.T) {
	tests := []struct {
		indicator Tile
		expected  Tile
	}{
		{NewTile(Man, 5), NewTile(Man, 6)},
		{NewTile(Man, 9), NewTile(Man, 1)}, // Wraps around
		{NewTile(Ji, 1), NewTile(Ji, 2)},   // East→South
		{NewTile(Ji, 4), NewTile(Ji, 1)},   // North→East
		{NewTile(Ji, 5), NewTile(Ji, 6)},   // White→Green
		{NewTile(Ji, 7), NewTile(Ji, 5)},   // Red→White
	}

	for _, tt := range tests {
		result := tt.indicator.NextDora()
		if !result.Equals(tt.expected) {
			t.Errorf("NextDora(%s): expected %s, got %s", tt.indicator, tt.expected, result)
		}
	}
}

func TestTileWindAndDragonNames(t *testing.T) {
	winds := []string{"East", "South", "West", "North"}
	for i, name := range winds {
		tile := NewTile(Ji, i+1)
		if !tile.IsWind() {
			t.Errorf("Tile %s should be a wind", tile)
		}
		if tile.GetWindName() != name {
			t.Errorf("Expected wind name %s, got %s", name, tile.GetWindName())
		}
	}

	dragons := []string{"White", "Green", "Red"}
	for i, name := range dragons {
		tile := NewTile(Ji, i+5)
		if !tile.IsDragon() {
			t.Errorf("Tile %s should be a dragon", tile)
		}
		if tile.GetDragonName() != name {
			t.Errorf("Expected dragon name %s, got %s", name, tile.GetDragonName())
		}
	}
}

func TestParseTiles(t *testing.T) {
	notations := []string{"1m", "2m", "3m", "5p", "0s"}
	tiles, err := ParseTiles(notations)
	if err != nil {
		t.Fatalf("ParseTiles failed: %v", err)
	}
	if len(tiles) != 5 {
		t.Errorf("Expected 5 tiles, got %d", len(tiles))
	}
	if tiles[4].IsAkaDora != true {
		t.Error("Last tile should be aka dora")
	}
}

func TestTilesToString(t *testing.T) {
	tiles := []Tile{
		NewTile(Man, 1),
		NewTile(Pin, 5),
		NewAkaDoraTile(Sou),
	}
	strs := TilesToString(tiles)
	expected := []string{"1m", "5p", "0s"}

	for i, str := range strs {
		if str != expected[i] {
			t.Errorf("Index %d: expected %s, got %s", i, expected[i], str)
		}
	}
}
