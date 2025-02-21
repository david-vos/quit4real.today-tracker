package query

import (
	"quit4real.today/src/model"
	"quit4real.today/src/repository"
)

type UserQueryHandlerImpl struct {
	UserRepository repository.UserRepository
}

func (handler *UserQueryHandlerImpl) GetAll() ([]model.User, error) {
	users, err := handler.UserRepository.GetAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (handler *UserQueryHandlerImpl) GetAllSteamVerified() ([]model.User, error) {
	users, err := handler.UserRepository.GetAllSteamVerified()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (handler *UserQueryHandlerImpl) GetById(username string) (model.User, error) {
	user, err := handler.UserRepository.GetById(username)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}
