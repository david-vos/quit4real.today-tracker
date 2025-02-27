package service

type UserService interface {
	UpdateUserTrackers(platformUserId string)
	CreateUserTrackers(platformUserId string)
}
