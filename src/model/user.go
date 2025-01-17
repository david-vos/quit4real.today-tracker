package model

import "database/sql"

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	SteamId string `json:"steam_id"`
	ApiKey  string `json:"api_key"`
}

func MapUser(rows *sql.Rows) (User, error) {
	var user User
	if err := rows.Scan(
		&user.ID,
		&user.Name,
		&user.SteamId,
		&user.ApiKey,
	); err != nil {
		return User{}, err
	}
	return user, nil
}
