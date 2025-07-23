package http

import (
	"log/slog"
	"net/http"
	gamemap "server/internal/game/map"
	"strconv"
	"time"
)

func RandomMapHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameters for width and height
	query := r.URL.Query()
	width := uint16(20)
	height := uint16(20)

	if w := query.Get("width"); w != "" {
		if parsedWidth, err := strconv.ParseUint(w, 10, 16); err == nil && parsedWidth >= 5 && parsedWidth <= 100 {
			width = uint16(parsedWidth)
		}
	}

	if h := query.Get("height"); h != "" {
		if parsedHeight, err := strconv.ParseUint(h, 10, 16); err == nil && parsedHeight >= 5 && parsedHeight <= 100 {
			height = uint16(parsedHeight)
		}
	}

	size := gamemap.Size{Width: width, Height: height}
	players := []gamemap.Player{{Index: 0, Owner: 1, IsActive: true}, {Index: 1, Owner: 2, IsActive: true}}

	// Use current time as seed for true randomness
	config := gamemap.DefaultGeneratorConfig()
	config.Seed = time.Now().UnixNano()

	m, err := gamemap.GenerateMap("new", size, players, config)
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
