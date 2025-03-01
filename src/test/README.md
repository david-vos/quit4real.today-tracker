# Test Suite for Quit4Real.Today Tracker

This directory contains the test suite for the Quit4Real.Today Tracker application.

## Test Structure

The test suite is organized as follows:

- `main_test.go`: Contains the test setup and teardown logic, including the test database initialization, test server setup, and mock Steam API.
- `user_test.go`: Tests for user-related endpoints (registration, login, Steam ID).
- `fails_test.go`: Tests for the fails leaderboard endpoint.
- `games_test.go`: Tests for the game search endpoint.
- `subscription_test.go`: Tests for the subscription-related endpoints.

## Test Environment

The test environment uses:

- A temporary SQLite database (`test_database.db`) that is created before tests and deleted after tests.
- A test HTTP server using `httptest.NewServer` instead of the real server, which avoids port conflicts and ensures proper cleanup.
- A mock Steam API server that intercepts and responds to Steam API requests with predefined data.
- Test logs are written to `test_logs.log`.

## Mock Steam API

The test suite includes a mock Steam API server that intercepts requests to the Steam API and returns predefined responses. This allows testing of functionality that depends on the Steam API without making actual API calls. The mock server handles the following endpoints:

- `/ISteamUser/ResolveVanityURL`: Returns a predefined Steam ID for vanity URL resolution.
- `/ISteamUser/GetPlayerSummaries`: Returns player profile information.
- `/IPlayerService/GetOwnedGames`: Returns a list of owned games.
- `/IPlayerService/GetRecentlyPlayedGames`: Returns recently played games.

## Running Tests

To run the tests, use the following command from the project root:

```bash
go test -v ./src/test
```

## Adding New Tests

To add new tests:

1. Create a new test file or add to an existing one based on the functionality being tested.
2. Use the `makeRequest` helper function to make HTTP requests to the test server.
3. Follow the existing test patterns for consistency.
4. If your test requires Steam API interaction, ensure the mock API server handles the necessary endpoints.

## Test Helper Functions

- `makeRequest(method, path, body, headers)`: Helper function to make HTTP requests to the test server.

## Notes

- The test server is automatically started and stopped for each test run.
- The test database is automatically created and destroyed for each test run.
- All tests use the test server URL as the base URL, which is set up in `main_test.go`.
- The mock Steam API server is set up in `main_test.go` and is used for all tests that interact with the Steam API.

## Test Coverage

Currently, the tests cover the "happy path" for each endpoint:

- User registration
- User login
- Getting a Steam ID
- Getting the fails leaderboard
- Searching for games
- Adding a subscription
- Getting subscriptions

## Adding New Tests

To add new tests:

1. Identify the endpoint you want to test
2. Create a new test function in the appropriate test file
3. Use the `makeRequest` helper function to make HTTP requests
4. Assert that the response is as expected

Example:

```go
func TestNewEndpoint(t *testing.T) {
    // Make the request
    resp, err := makeRequest(
        "GET",
        "/new/endpoint",
        nil,
        nil,
    )
    if err != nil {
        t.Fatalf("Failed to make request: %v", err)
    }
    defer resp.Body.Close()

    // Check the response status code
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
    }

    // Check the response body
    var response map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        t.Fatalf("Failed to decode response: %v", err)
    }

    // Assert on the response
    if _, ok := response["expectedField"]; !ok {
        t.Error("Response does not contain expectedField")
    }
}
``` 