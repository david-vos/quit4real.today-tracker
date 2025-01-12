package repository

import (
	"database/sql"
	"project/config"
)

type DBExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func executeQuery(db DBExecutor, query string, args ...interface{}) error {
	_, err := db.Exec(query, args...)
	return err
}

func fetchRowsWithClose(db DBExecutor, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.Query(query, args...)
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
