package model

import "database/sql"

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func MapUser(rows *sql.Rows) (User, error) {
	var user User
	if err := rows.Scan(
		&user.ID,
		&user.Name,
	); err != nil {
		return User{}, err
	}
	return user, nil
}
