package repository

import "quit4real.today/src/model"

type FailRepository interface {
	Get(userID string) ([]model.GameFailureRecord, error)
	Add(failure model.GameFailureRecord) error
	GetTopLeaderBoard() ([]model.FailResponse, error)
}
