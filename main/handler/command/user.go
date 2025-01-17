package command

import (
	"project/main/model"
	"project/main/repository"
)

type UserCommandHandler struct {
	userRepository *repository.UserRepository
}

func (handler *UserCommandHandler) Add(user model.User) error {
	var err = handler.userRepository.Add(user)
	if err != nil {
		return err
	}
	return nil
}
