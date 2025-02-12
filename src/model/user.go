package model

import "database/sql"

type User struct {
	ID       int    `json:"id"`   // Using string to match the database schema
	Name     string `json:"name"` // Username
	SteamID  string `json:"steamid"`
	Password string `json:"password"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func MapUser(rows *sql.Rows) (User, error) {
	var user User
	if err := rows.Scan(
		&user.ID,
		&user.Name,
		&user.Password,
		&user.SteamID,
	); err != nil {
		return User{}, err
	}
	return user, nil
}
