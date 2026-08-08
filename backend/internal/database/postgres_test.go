package database

import (
	"encoding/json"
	"testing"
)

func TestTownshipSeedSnapshot(t *testing.T) {
	var seeds []townshipSeed
	if err := json.Unmarshal(townshipsJSON, &seeds); err != nil {
		t.Fatal(err)
	}
	if len(seeds) != 368 {
		t.Fatalf("got %d township seeds, want 368", len(seeds))
	}

	codes := make(map[string]bool, len(seeds))
	parents := make(map[string]bool)
	for _, seed := range seeds {
		if seed.City == "" || seed.Name == "" || seed.Code == "" {
			t.Fatalf("incomplete seed: %+v", seed)
		}
		if codes[seed.Code] {
			t.Fatalf("duplicate external code %s", seed.Code)
		}
		codes[seed.Code] = true
		parents[seed.City] = true
	}
	if len(parents) != 20 {
		t.Fatalf("got %d website city parents, want 20 (excluding 其他)", len(parents))
	}
}
