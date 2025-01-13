package repository

import (
	"project/models"
)

type UserRepoController struct {
	DbContr *DatabaseController
}

func (c *UserRepoController) GetAllUsers() ([]models.User, error) {
	query := `SELECT * FROM users;`
	rows, err := c.DbContr.FetchRowsWithClose(query)
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

func (c *UserRepoController) CreateUser(user models.User) error {
	query := "INSERT INTO users (id, name, steam_id, api_key) VALUES (?, ?, ?, ?)"
	err := c.DbContr.ExecuteQuery(query, user.ID, user.Name, user.SteamId, user.ApiKey)
	return err
}
