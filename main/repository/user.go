package repository

import (
	"database/sql"
	"project/logger"
	"project/main/model"
)

type UserRepository struct {
	DatabaseImp *DatabaseImpl
}

func (repository *UserRepository) GetAll() ([]model.User, error) {
	query := `SELECT * FROM users;`
	rows, err := repository.DatabaseImp.FetchRows(query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := closeRows(rows)
		if err != nil {
			logger.Fail(err.Error())
		}
	}(rows)

	var users []model.User
	for rows.Next() {
		user, err := model.MapUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (repository *UserRepository) Add(user model.User) error {
	query := "INSERT INTO users (id, name, steam_id, api_key) VALUES (?, ?, ?, ?)"
	err := repository.DatabaseImp.ExecuteQuery(query, user.ID, user.Name, user.SteamId, user.ApiKey)
	return err
}
