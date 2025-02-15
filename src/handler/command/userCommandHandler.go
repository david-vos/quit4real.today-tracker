package command

import (
	"quit4real.today/src/model"
	"quit4real.today/src/repository"
)

type UserCommandHandler struct {
	UserRepository *repository.UserRepository
}

func (handler *UserCommandHandler) Add(user model.User) error {
	var err = handler.UserRepository.Add(user)
	if err != nil {
		return err
	}
	return nil
}

func (handler *UserCommandHandler) Update(user model.User) error {
	var err = handler.UserRepository.Update(user)
	if err != nil {
		return err
	}
	return nil
}
