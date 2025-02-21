package model

import "database/sql"

// Game represents a game in the system.
type Game struct {
	ID       string   `json:"id"`       // Game ID
	Name     string   `json:"name"`     // Game name
	Platform Platform `json:"platform"` // Embedded Platform object
}

// MapGame maps SQL rows to a Game struct.
func MapGame(rows *sql.Rows) (Game, error) {
	var game Game
	var platformID string
	if err := rows.Scan(
		&game.ID,
		&game.Name,
		&platformID,
	); err != nil {
		return Game{}, err
	}

	// Create a Platform object based on the platform ID
	game.Platform = Platform{ID: platformID}
	return game, nil
}

// MatchedDbGameToSteamGameInfo holds the mapping between a database tracked game and a Steam game.
type MatchedDbGameToSteamGameInfo struct {
	DbTrack      Subscription
	SteamApiGame SteamGame
	Failed       bool
}
