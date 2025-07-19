package http

import (
	"log/slog"
	"net/http"
	gamemap "server/internal/game/map"
)

func RandomMapHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	size := gamemap.Size{Width: 20, Height: 20}
	players := []gamemap.Player{{Index: 0, Owner: 1, IsActive: true}, {Index: 1, Owner: 2, IsActive: true}}
	
	m, err := gamemap.GenerateMap("base", size, players)
	if err != nil {
		slog.Error("failed to generate map", "error", err)
		http.Error(w, "Failed to generate map", http.StatusInternalServerError)
		return
	}
	
	data, err := gamemap.ExportToJSON(m)
	if err != nil {
		slog.Error("failed to export map to JSON", "error", err)
		http.Error(w, "Failed to export map", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
