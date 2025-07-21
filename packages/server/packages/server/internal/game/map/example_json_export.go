// +build ignore

package main

import (
	"fmt"
	"log"
	"server/internal/game/map"
)

func main() {
	// Create a small map
	size := gamemap.Size{Width: 3, Height: 3}
	players := []gamemap.Player{
		{Index: 0, Owner: 1, IsActive: true},
	}
	config := gamemap.GeneratorConfig{
		MountainDensity:   0.2,
		CastleDensity:     0.1,
		MinCastleDistance: 2,
		Seed:              123,
	}

	// Generate map
	generator := gamemap.NewBaseMapGenerator(config)
	gameMap, err := generator.Generate(size, players, config)
	if err != nil {
		log.Fatalf("Failed to generate map: %v", err)
	}

	// Export to JSON
	jsonData, err := gamemap.ExportToJSON(gameMap)
	if err != nil {
		log.Fatalf("Failed to export map: %v", err)
	}

	// Print the JSON
	fmt.Println(string(jsonData))
}
