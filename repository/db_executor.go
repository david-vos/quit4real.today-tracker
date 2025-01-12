package repository

import "database/sql"

// DBExecutor abstracts the methods needed for database interaction
type DBExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}
