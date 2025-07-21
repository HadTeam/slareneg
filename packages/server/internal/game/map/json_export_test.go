package gamemap

import (
	"encoding/json"
	"server/internal/game/block"
	"testing"
)

func TestExportToJSON(t *testing.T) {
	// Create a small map using the base generator
	size := Size{Width: 5, Height: 5}
	players := []Player{
		{Index: 0, Owner: 1, IsActive: true},
		{Index: 1, Owner: 2, IsActive: true},
	}
	config := GeneratorConfig{
		MountainDensity:   0.3,
		CastleDensity:     0.2,
		MinCastleDistance: 3,
		Seed:              42, // Fixed seed for deterministic test
	}

	// Generate map using base generator
	generator := NewBaseMapGenerator(config)
	gameMap, err := generator.Generate(size, players, config)
	if err != nil {
		t.Fatalf("Failed to generate map: %v", err)
	}

	// Export to JSON
	jsonData, err := ExportToJSON(gameMap)
	if err != nil {
		t.Fatalf("Failed to export map to JSON: %v", err)
	}

	// Verify it's valid JSON
	var exported ExportedMap
	err = json.Unmarshal(jsonData, &exported)
	if err != nil {
		t.Fatalf("Failed to unmarshal exported JSON: %v", err)
	}

	// Verify basic structure
	if exported.Size.Width != size.Width || exported.Size.Height != size.Height {
		t.Errorf("Expected size %v, got %v", size, exported.Size)
	}

	if exported.Info.Id == "" {
		t.Error("Expected non-empty map ID")
	}

	if len(exported.Blocks) != int(size.Height) {
		t.Errorf("Expected %d rows, got %d", size.Height, len(exported.Blocks))
	}

	for i, row := range exported.Blocks {
		if len(row) != int(size.Width) {
			t.Errorf("Row %d: expected %d columns, got %d", i, size.Width, len(row))
		}
	}

	// Check that we have at least some non-zero blocks (kings from players)
	nonZeroBlocks := 0
	for _, row := range exported.Blocks {
		for _, b := range row {
			if b.Num != 0 {
				nonZeroBlocks++
			}
		}
	}
	if nonZeroBlocks == 0 {
		t.Error("Expected at least some non-zero blocks in the map")
	}

	// Verify JSON formatting (should be indented)
	jsonStr := string(jsonData)
	if len(jsonStr) == 0 {
		t.Error("JSON output is empty")
	}
	
	// Check for proper indentation
	if !contains(jsonStr, "\n  ") {
		t.Error("JSON should be indented with 2 spaces")
	}
}

func TestExportToJSON_EmptyMap(t *testing.T) {
	// Create an empty map
	size := Size{Width: 3, Height: 3}
	info := Info{
		Id:   "test-empty",
		Name: "Empty Test Map",
		Desc: "A test map with no blocks",
	}
	
	emptyMap := NewEmptyBaseMap(size, info)
	if emptyMap == nil {
		t.Fatal("Failed to create empty map")
	}

	// Export to JSON
	jsonData, err := ExportToJSON(emptyMap)
	if err != nil {
		t.Fatalf("Failed to export empty map to JSON: %v", err)
	}

	// Verify it's valid JSON
	var exported ExportedMap
	err = json.Unmarshal(jsonData, &exported)
	if err != nil {
		t.Fatalf("Failed to unmarshal exported JSON: %v", err)
	}

	// Verify all blocks are zero-valued
	for y, row := range exported.Blocks {
		for x, b := range row {
			if b.Num != 0 || b.Owner != 0 {
				t.Errorf("Expected zero block at (%d,%d), got Num=%d, Owner=%d", 
					x, y, b.Num, b.Owner)
			}
		}
	}
}

func TestBlockDTO_JSONMarshaling(t *testing.T) {
	// Test individual BlockDTO marshaling
	dto := BlockDTO{
		Num:   block.Num(5),
		Owner: block.Owner(1),
		Meta: block.Meta{
			Name:        "TestBlock",
			Description: "A test block",
		},
	}

	jsonData, err := json.Marshal(dto)
	if err != nil {
		t.Fatalf("Failed to marshal BlockDTO: %v", err)
	}

	var unmarshaled BlockDTO
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal BlockDTO: %v", err)
	}

	if unmarshaled.Num != dto.Num {
		t.Errorf("Expected Num %d, got %d", dto.Num, unmarshaled.Num)
	}
	if unmarshaled.Owner != dto.Owner {
		t.Errorf("Expected Owner %d, got %d", dto.Owner, unmarshaled.Owner)
	}
	if unmarshaled.Meta.Name != dto.Meta.Name {
		t.Errorf("Expected Meta.Name %s, got %s", dto.Meta.Name, unmarshaled.Meta.Name)
	}
	if unmarshaled.Meta.Description != dto.Meta.Description {
		t.Errorf("Expected Meta.Description %s, got %s", dto.Meta.Description, unmarshaled.Meta.Description)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsAt(s, substr, 0)
}

func containsAt(s, substr string, start int) bool {
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
