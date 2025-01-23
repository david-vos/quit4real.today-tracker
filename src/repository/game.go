package repository

import (
	"database/sql"
	"quit4real.today/logger"
)

type GameRepository struct {
	DatabaseImpl *DatabaseImpl
}

func (repository *GameRepository) Add(id string, name string, platformId string) error {
	query := `INSERT INTO games (id, name, platform_id) VALUES (?, ?, ?)`
	return repository.DatabaseImpl.ExecuteQuery(query, id, name, platformId)
}

func (repository *GameRepository) Exists(id string, platformId string) bool {
	query := `SELECT * FROM games WHERE id = ? AND platform_id = ?`
	rows, err := repository.DatabaseImpl.FetchRows(query, id, platformId)
	defer func(rows *sql.Rows) {
		if err := closeRows(rows); err != nil {
			logger.Fail("failed to close rows: " + err.Error())
		}
	}(rows)
	if err != nil {
		return false
	}
	return true
}
