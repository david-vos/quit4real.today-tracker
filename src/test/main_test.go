package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"quit4real.today/logger"
	"quit4real.today/src"
	databaseImpl "quit4real.today/src/db"
	"quit4real.today/src/model"
)

var (
	testApp      *src.App
	testServer   *httptest.Server
	mockSteamAPI *httptest.Server
	testDBPath   = "test_database.db"
	baseURL      string
)

// TestMain sets up the test environment before running tests and tears it down afterward
func TestMain(m *testing.M) {
	logger.InitLogger("test_logs.log")
	setup()

	code := m.Run()
	teardown()
	os.Exit(code)
}

// setup initializes the test environment
func setup() {
	// Remove any existing test database
	os.Remove(testDBPath)

	// Create a test database
	db, err := databaseImpl.Connect(testDBPath)
	if err != nil {
		log.Fatalf("Failed to open test database: %v", err)
	}

	// Use a relative path to the migrations directory that will work both locally and in CI
	migrationsPath := getDynamicMigrationsPath()
	log.Printf("Using migrations path: %s", migrationsPath)

	// Initialize the database schema
	err = databaseImpl.ApplyMigrations(db, migrationsPath)
	if err != nil {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	// Insert test data
	insertTestData(db)

	// Set up mock Steam API server
	setupMockSteamAPI()

	// Initialize the application with the test database
	testApp = src.AppInit(db)

	// Initialize the endpoints without starting the real server
	testApp.Endpoints.UserEndpoint.User()
	testApp.Endpoints.FailEndpoint.Fail()
	testApp.Endpoints.SubscriptionEndpoint.Subscription()
	testApp.Endpoints.GamesEndpoint.Games()

	// Create a test server using the router from the app
	testServer = httptest.NewServer(testApp.Endpoints.Router)
	baseURL = testServer.URL

	// Wait a moment for the server to start
	time.Sleep(100 * time.Millisecond)
}

// setupMockSteamAPI creates a mock Steam API server
func setupMockSteamAPI() {
	mockSteamAPI = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Mock Steam API received request: %s", r.URL.String())

		// Parse the URL to determine which API endpoint is being called
		if strings.Contains(r.URL.Path, "/ISteamUser/ResolveVanityURL") {
			// Always return success regardless of authentication
			vanityResponse := model.SteamApiVanityResponse{
				Response: model.SteamApiVanity{
					SteamId: "76561198012345678",
					Success: 1,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(vanityResponse)
			return
		}

		if strings.Contains(r.URL.Path, "/ISteamUser/GetPlayerSummaries") {
			// Mock response for player summaries
			userInfoResponse := model.SteamApiUserInfoResponse{
				Response: struct {
					Players []model.SteamApiUserInfo `json:"players"`
				}{
					Players: []model.SteamApiUserInfo{
						{
							SteamID:                  "76561198012345678",
							CommunityVisibilityState: 3,
							ProfileState:             1,
							PersonaName:              "TestUser",
							ProfileURL:               "https://steamcommunity.com/id/testuser/",
							Avatar:                   "https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/fe/fef49e7fa7e1997310d705b2a6158ff8dc1cdfeb.jpg",
							AvatarMedium:             "https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/fe/fef49e7fa7e1997310d705b2a6158ff8dc1cdfeb_medium.jpg",
							AvatarFull:               "https://steamcdn-a.akamaihd.net/steamcommunity/public/images/avatars/fe/fef49e7fa7e1997310d705b2a6158ff8dc1cdfeb_full.jpg",
							PersonaState:             1,
							RealName:                 "Test User",
							TimeCreated:              1234567890,
						},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(userInfoResponse)
			return
		}

		if strings.Contains(r.URL.Path, "/IPlayerService/GetOwnedGames") {
			// Mock response for owned games
			gamesResponse := model.SteamAPIResponseAllGames{
				Response: model.SteamAPIAllResponse{
					GameCount: 1,
					Games: []model.SteamAPIAllGame{
						{
							Appid:           730,
							Name:            "Counter-Strike 2",
							PlaytimeForever: 1000,
							ImgIconUrl:      "69f7ebe2735c366c65c0b33dae00e12dc40edbe4",
							RtimeLastPlayed: int(time.Now().Unix()),
						},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(gamesResponse)
			return
		}

		if strings.Contains(r.URL.Path, "/IPlayerService/GetRecentlyPlayedGames") {
			// Mock response for recently played games
			recentGamesResponse := model.SteamApiResponse{
				Response: model.SteamApiGetLastPlayed{
					TotalCount: 1,
					Games: []model.SteamGame{
						{
							AppID:           730,
							Name:            "Counter-Strike 2",
							Playtime2Weeks:  120,
							PlaytimeForever: 1000,
							ImgIconURL:      "69f7ebe2735c366c65c0b33dae00e12dc40edbe4",
						},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(recentGamesResponse)
			return
		}

		// Default response for unknown endpoints
		log.Printf("Unknown Steam API endpoint: %s", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response": {"success": true}}`))
	}))

	// Override the Steam API URL in the app configuration
	// This is a bit of a hack, but it allows us to intercept Steam API calls
	os.Setenv("STEAM_API_URL", mockSteamAPI.URL)
	log.Printf("Set Steam API URL to: %s", mockSteamAPI.URL)
}

// insertTestData inserts required data for tests
func insertTestData(db *sql.DB) {
	// Check if platform exists
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM platforms WHERE id = 'steam'").Scan(&count)
	if err != nil {
		log.Printf("Error checking platform: %v", err)
	}

	// Insert platform if it doesn't exist
	if count == 0 {
		_, err := db.Exec("INSERT INTO platforms (id, name) VALUES ('steam', 'Steam')")
		if err != nil {
			log.Printf("Failed to insert platform: %v", err)
		}
	}

	// Check if game exists
	err = db.QueryRow("SELECT COUNT(*) FROM games WHERE id = '730' AND platform_id = 'steam'").Scan(&count)
	if err != nil {
		log.Printf("Error checking game: %v", err)
	}

	// Insert game if it doesn't exist
	if count == 0 {
		_, err := db.Exec("INSERT INTO games (id, name, platform_id) VALUES ('730', 'Counter-Strike 2', 'steam')")
		if err != nil {
			log.Printf("Failed to insert game: %v", err)
		}
	}

	// Check if user exists
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE name = 'testuser'").Scan(&count)
	if err != nil {
		log.Printf("Error checking user: %v", err)
	}

	// Insert user if it doesn't exist
	if count == 0 {
		_, err := db.Exec("INSERT INTO users (name) VALUES ('testuser')")
		if err != nil {
			log.Printf("Failed to insert user: %v", err)
		}
	}

	// Insert a test failure record for the leaderboard test
	_, err = db.Exec(`
		INSERT INTO game_failure_records 
		(display_name, platform_id, platform_game_id, platform_user_id, duration_minutes, reason, timestamp) 
		VALUES ('Test User', 'steam', '730', 'testuser', 60, 'Test Reason', datetime('now'))
	`)
	if err != nil {
		log.Printf("Failed to insert failure record: %v", err)
	}
}

// teardown cleans up the test environment
func teardown() {
	// Close the test server
	if testServer != nil {
		testServer.Close()
	}

	// Close the mock Steam API server
	if mockSteamAPI != nil {
		mockSteamAPI.Close()
	}

	// Close the database connection
	if testApp != nil && testApp.DatabaseImpl != nil && testApp.DatabaseImpl.DB != nil {
		testApp.DatabaseImpl.DB.Close()
	}

	os.Remove(testDBPath)
	os.Remove("test_logs.log")
}

// Helper function to make HTTP requests
func makeRequest(method, path string, body []byte, headers map[string]string) (*http.Response, error) {
	var req *http.Request
	var err error

	if body != nil {
		req, err = http.NewRequest(method, baseURL+path, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequest(method, baseURL+path, nil)
	}

	if err != nil {
		return nil, err
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Make the request
	client := &http.Client{}
	return client.Do(req)
}

// getDynamicMigrationsPath returns a path to the migrations directory that works in any environment
func getDynamicMigrationsPath() string {
	// First try the current directory
	if _, err := os.Stat("src/db/migrations"); err == nil {
		return "src/db/migrations"
	}

	// Then try parent directory (in case we're running from the src/test directory)
	if _, err := os.Stat("../db/migrations"); err == nil {
		return "../db/migrations"
	}

	// As a fallback, try the absolute path (useful for local development)
	workDir, err := os.Getwd()
	if err != nil {
		log.Printf("Warning: couldn't get working directory: %v", err)
		return "src/db/migrations" // Default to the relative path
	}

	// Check if we're in the root directory or a subdirectory
	possiblePaths := []string{
		filepath.Join(workDir, "src/db/migrations"),
		filepath.Join(workDir, "../db/migrations"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// If all else fails, return the basic path and let it fail with a clear error
	return "src/db/migrations"
}
