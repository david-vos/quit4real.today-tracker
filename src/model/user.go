package model

import "database/sql"

type User struct {
	ID   int    `json:"id"`   // Using string to match the database schema
	Name string `json:"name"` // Username
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
