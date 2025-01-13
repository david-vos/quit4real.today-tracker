package repository

import (
	"database/sql"
	"project/config"
)

type DatabaseController struct {
	DB *sql.DB
}

func (c *DatabaseController) ExecuteQuery(query string, args ...interface{}) error {
	_, err := c.DB.Exec(query, args...)
	return err
}

func (c *DatabaseController) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return c.DB.Query(query, args...)
}

func (c *DatabaseController) FetchRowsWithClose(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func closeRows(rows *sql.Rows) error {
	if err := rows.Close(); err != nil {
		config.HandleError("Error closing rows", err)
		return err
	}
	return nil
}
