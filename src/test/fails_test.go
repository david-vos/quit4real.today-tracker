package test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

// TestGetFailsLeaderboard tests the fails leaderboard endpoint
func TestGetFailsLeaderboard(t *testing.T) {
	// Make the request
	resp, err := makeRequest("GET", "/fail/leaderboard", nil, nil)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var response []map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify the response is an array (even if empty)
	if response == nil {
		t.Error("Response should be an array, got nil")
	}
}
