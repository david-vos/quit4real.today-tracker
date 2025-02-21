package impl

import (
	"database/sql"
	"fmt"
)

type DatabaseImpl struct {
	DB *sql.DB
}

func (c *DatabaseImpl) ExecuteQuery(query string, args ...interface{}) error {
	_, err := c.DB.Exec(query, args...)
	return err
}

func (c *DatabaseImpl) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return c.DB.Query(query, args...)
}

func (c *DatabaseImpl) FetchRows(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rows: %w", err)
	}
	return rows, nil
}

func (c *DatabaseImpl) CloseRows(rows *sql.Rows) error {
	if err := rows.Close(); err != nil {
		return fmt.Errorf("error closing rows: %w", err)
	}
	return nil
}
