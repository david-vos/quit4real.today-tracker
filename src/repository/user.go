package repository

import (
	"database/sql"
	"quit4real.today/logger"
	"quit4real.today/src/model"
)

type UserRepository struct {
	DatabaseImpl *DatabaseImpl
}

func (repository *UserRepository) GetAll() ([]model.User, error) {
	query := `SELECT * FROM users;`
	rows, err := repository.DatabaseImpl.FetchRows(query)
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
	query := "INSERT INTO users (name) VALUES (?)"
	err := repository.DatabaseImpl.ExecuteQuery(query, user.Name)
	return err
}
