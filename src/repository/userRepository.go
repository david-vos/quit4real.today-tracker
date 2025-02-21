package repository

import "quit4real.today/src/model"

type UserRepository interface {
	GetAll() ([]model.User, error)
	Add(user model.User) error
	Update(user model.User) error
	GetAllSteamVerified() ([]model.User, error)
	GetById(username string) (model.User, error)
}
