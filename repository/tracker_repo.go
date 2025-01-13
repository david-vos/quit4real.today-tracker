package repository

import (
	"project/models"
)

type TrackerRepoController struct {
	DbContr *DatabaseController
}

func (c *TrackerRepoController) CreateTracker(steamId string, gameId string) error {
	query := "INSERT INTO user_tracker (steam_id, game_id, played_amount) VALUES (?, ?, ?);"
	return c.DbContr.ExecuteQuery(query, steamId, gameId, 0)
}

func (c *TrackerRepoController) UpdateTracker(steamId string, gameId string, playedAmount int) error {
	query := "UPDATE user_tracker SET played_amount = ? WHERE steam_id = ? AND game_id = ?;"
	return c.DbContr.ExecuteQuery(query, playedAmount, steamId, gameId)
}

func (c *TrackerRepoController) GetTracker(steamId string, gameId string) (models.Tracker, error) {
	query := "SELECT * FROM user_tracker WHERE steam_id = ? AND game_id = ?;"
	rows, err := c.DbContr.FetchRowsWithClose(query, steamId, gameId)
	if err != nil {
		return models.Tracker{}, err
	}
	defer closeRows(rows)

	if rows.Next() {
		return models.MapTracker(rows)
	}
	return models.Tracker{}, nil
}

func (c *TrackerRepoController) GetUserTracker(steamId string) ([]models.Tracker, error) {
	query := "SELECT * FROM user_tracker WHERE steam_id = ?;"
	rows, err := c.DbContr.FetchRowsWithClose(query, steamId)
	if err != nil {
		return nil, err
	}
	defer closeRows(rows)

	var trackers []models.Tracker
	for rows.Next() {
		tracker, err := models.MapTracker(rows)
		if err != nil {
			return nil, err
		}
		trackers = append(trackers, tracker)
	}
	return trackers, nil
}
