package repository

import (
	"project/models"
)

func GetAllUsers(db DBExecutor) ([]models.User, error) {
	query := `SELECT * FROM users;`
	rows, err := fetchRowsWithClose(db, query)
	if err != nil {
		return nil, err
	}
	defer closeRows(rows)

	var users []models.User
	for rows.Next() {
		user, err := models.MapUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func CreateUser(db DBExecutor, user models.User) error {
	query := "INSERT INTO users (id, name, steam_id, api_key) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, user.ID, user.Name, user.SteamId, user.ApiKey)
	return err
}
