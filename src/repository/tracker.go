package repository

import (
	"database/sql"
	"quit4real.today/logger"
	"quit4real.today/src/model"
)

type TrackerRepository struct {
	DatabaseImpl *DatabaseImpl
}

func (repository *TrackerRepository) Add(steamId string, gameId string) error {
	query := "INSERT INTO user_tracker (steam_id, game_id, played_amount) VALUES (?, ?, ?);"
	return repository.DatabaseImpl.ExecuteQuery(query, steamId, gameId, 0)
}

func (repository *TrackerRepository) Update(steamId string, gameId string, playedAmount int) error {
	query := "UPDATE user_tracker SET played_amount = ? WHERE steam_id = ? AND game_id = ?;"
	return repository.DatabaseImpl.ExecuteQuery(query, playedAmount, steamId, gameId)
}

func (repository *TrackerRepository) Get(steamId string, gameId string) (model.Tracker, error) {
	query := "SELECT * FROM user_tracker WHERE steam_id = ? AND game_id = ?;"
	rows, err := repository.DatabaseImpl.FetchRows(query, steamId, gameId)
	if err != nil {
		return model.Tracker{}, err
	}
	defer func(rows *sql.Rows) {
		err := closeRows(rows)
		if err != nil {
			logger.Fail("failed to close rows" + err.Error())
		}
	}(rows)

	if rows.Next() {
		return model.MapTracker(rows)
	}
	return model.Tracker{}, nil
}

func (repository *TrackerRepository) GetAll(steamId string) ([]model.Tracker, error) {
	query := "SELECT * FROM user_tracker WHERE steam_id = ?;"
	rows, err := repository.DatabaseImpl.FetchRows(query, steamId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := closeRows(rows)
		if err != nil {
			logger.Fail("failed to close rows" + err.Error())
		}
	}(rows)

	var trackers []model.Tracker
	for rows.Next() {
		tracker, err := model.MapTracker(rows)
		if err != nil {
			return nil, err
		}
		trackers = append(trackers, tracker)
	}
	return trackers, nil
}
