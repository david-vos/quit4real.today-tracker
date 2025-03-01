package test

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUserRegistration tests the user registration endpoint
func TestUserRegistration(t *testing.T) {
	// Test data
	requestBody := `{"username":"testuser2","password":"testpassword"}`

	// Make the request
	resp, err := makeRequest("POST", "/users", []byte(requestBody), map[string]string{
		"Content-Type": "application/json",
	})

	// Check for errors
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// For registration, we may not get a JSON response, so just check the status code
	resp.Body.Close()
}

// TestUserLogin tests the user login endpoint
func TestUserLogin(t *testing.T) {
	// First register a user
	registerBody := `{"username":"loginuser","password":"loginpassword"}`
	resp, err := makeRequest("POST", "/users", []byte(registerBody), map[string]string{
		"Content-Type": "application/json",
	})
	assert.NoError(t, err)
	resp.Body.Close()

	// Now test login
	loginBody := `{"username":"loginuser","password":"loginpassword"}`
	resp, err = makeRequest("POST", "/users/login", []byte(loginBody), map[string]string{
		"Content-Type": "application/json",
	})

	// Check for errors
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	assert.NoError(t, err)

	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	assert.NoError(t, err)

	// Verify the response contains the expected fields
	assert.Contains(t, response, "token")
}

// TestGetSteamID tests the get steam ID endpoint
func TestGetSteamID(t *testing.T) {
	// Make the request to get Steam ID for a vanity name
	resp, err := makeRequest("GET", "/users/testuser/steamId", nil, nil)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check the response status code - allow either 200 (OK) or 204 (No Content)
	statusOK := resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent
	if !statusOK {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusNoContent, resp.StatusCode)
	}

	// Only try to parse the response body if we got a 200 OK
	if resp.StatusCode == http.StatusOK {
		// Read and parse the response
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		var response map[string]interface{}
		err = json.Unmarshal(body, &response)
		assert.NoError(t, err)

		// Verify the response contains the expected fields
		assert.Contains(t, response, "steamId")
	}
}
