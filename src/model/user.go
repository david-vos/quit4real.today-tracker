package model

import (
	"database/sql"
)

type User struct {
	ID            int    `json:"id"`   // Using string to match the database schema
	Name          string `json:"name"` // Internal Username
	SteamUserName string `json:"steamUsername"`
	SteamID       string `json:"steamid"`
	Password      string `json:"password"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func MapUser(rows *sql.Rows) (User, error) {
	var user User
	var steamID sql.NullString
	var steamUserName sql.NullString

	if err := rows.Scan(
		&user.ID,
		&user.Name,
		&user.Password,
		&steamID,
		&steamUserName,
	); err != nil {
		return User{}, err
	}

	// Map nullable strings to empty strings
	user.SteamID = mapNullString(steamID)
	user.SteamUserName = mapNullString(steamUserName)

	return user, nil
}

// Helper function to map sql.NullString to a string
func mapNullString(nullStr sql.NullString) string {
	if nullStr.Valid {
		return nullStr.String
	}
	return ""
}
