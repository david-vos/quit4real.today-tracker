package model

import "database/sql"

type Platform struct {
	ID   string `json:"id"`   // Using string to match the database schema
	Name string `json:"name"` // Platform name
}

// MapPlatform maps SQL rows to a Platform struct.
func MapPlatform(rows *sql.Rows) (Platform, error) {
	var platform Platform
	if err := rows.Scan(
		&platform.ID,
		&platform.Name,
	); err != nil {
		return Platform{}, err
	}
	return platform, nil
}
