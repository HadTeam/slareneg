package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRandomMapHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/map/random", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RandomMapHandler)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response has correct content-type
	expected := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expected {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, expected)
	}

	// Check that the response is valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	// Check that the response has the expected structure
	if _, ok := result["size"]; !ok {
		t.Errorf("response missing 'size' field")
	}
	if _, ok := result["info"]; !ok {
		t.Errorf("response missing 'info' field")
	}
	if _, ok := result["blocks"]; !ok {
		t.Errorf("response missing 'blocks' field")
	}

	// Check size values
	if size, ok := result["size"].(map[string]interface{}); ok {
		if width, ok := size["width"].(float64); !ok || width != 20 {
			t.Errorf("expected width to be 20, got %v", size["width"])
		}
		if height, ok := size["height"].(float64); !ok || height != 20 {
			t.Errorf("expected height to be 20, got %v", size["height"])
		}
	}
}
