package repository

import (
	"database/sql"
	"log"
	"project/models"
)

func GetAllUsers(db *sql.DB) []models.User {
	query := `SELECT id, name, email FROM users;`
	rows, err := db.Query(query)
	if err != nil {
		return nil
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil
		}
		users = append(users, user)
	}

	return users
}

func CreateUser(db *sql.DB, user models.User) error {
	query := `INSERT INTO users (name, email) VALUES (?, ?);`
	_, err := db.Exec(query, user.Name, user.Email)
	if err != nil {
		return err
	}
	return nil
}
