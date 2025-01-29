package repository

import (
	"database/sql"
	"fmt"
	"quit4real.today/logger"
	"quit4real.today/src/model"
)

type GameRepository struct {
	DatabaseImpl *DatabaseImpl
}

func (repository *GameRepository) Add(id string, name string, platformId string) error {
	query := `INSERT INTO games (id, name, platform_id) VALUES (?, ?, ?)`
	return repository.DatabaseImpl.ExecuteQuery(query, id, name, platformId)
}

func (repository *GameRepository) Exists(id string, platformId string) bool {
	query := `SELECT 1 FROM games WHERE id = ? AND platform_id = ? LIMIT 1`
	rows, err := repository.DatabaseImpl.FetchRows(query, id, platformId)
	if err != nil {
		logger.Fail("failed to fetch rows: " + err.Error())
		return false
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			logger.Fail("failed to close rows: " + err.Error())
		}
	}(rows)

	return rows.Next()
}

func (repository *GameRepository) Search(searchParam string, platform string) ([]model.Game, error) {
	query := `SELECT * FROM games WHERE name LIKE ? COLLATE NOCASE AND platform_id = ? LIMIT 20`
	searchParam = "%" + searchParam + "%"
	rows, err := repository.DatabaseImpl.FetchRows(query, searchParam, platform)
	defer func(rows *sql.Rows) {
		if err := closeRows(rows); err != nil {
			logger.Fail("failed to close rows: " + err.Error())
		}
	}(rows)
	if err != nil {
		return nil, err
	}

	var games []model.Game
	for rows.Next() {
		game, err := model.MapGame(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to map game row: %w", err)
		}
		games = append(games, game)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	return games, nil
}
