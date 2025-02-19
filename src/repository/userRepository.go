package repository

import (
	"database/sql"
	"fmt"
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
	query := "INSERT INTO users (name, password) VALUES (?, ?)"
	err := repository.DatabaseImpl.ExecuteQuery(query, user.Name, user.Password)
	return err
}

func (repository *UserRepository) Update(user model.User) error {
	query := "UPDATE users SET name = ?, password = ?, steamid = ?, steam_display_name = ? WHERE id = ?"
	err := repository.DatabaseImpl.ExecuteQuery(query, user.Name, user.Password, user.ID, user.SteamID, user.SteamUserName)
	return err
}

func (repository *UserRepository) GetAllSteamVerified() ([]model.User, error) {
	query := `SELECT * FROM users WHERE steamid NOT NULL ;`
	rows, err := repository.DatabaseImpl.FetchRows(query)
	if err != nil {
		return nil, err
	}
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

func (repository *UserRepository) GetById(username string) (model.User, error) {
	query := `SELECT * FROM users WHERE name=?;`
	rows, err := repository.DatabaseImpl.FetchRows(query, username)
	if err != nil {
		return model.User{}, err
	}
	defer func(rows *sql.Rows) {
		err := closeRows(rows)
		if err != nil {
			logger.Fail(err.Error())
		}
	}(rows)
	for rows.Next() {
		user, err := model.MapUser(rows)
		if err != nil {
			return model.User{}, err
		}
		return user, nil
	}
	return model.User{}, fmt.Errorf("no user found with username %s", username)
}
