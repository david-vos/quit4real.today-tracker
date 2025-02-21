package repository

import "quit4real.today/src/model"

// GameRepository defines an interface for game-related database operations.
type GameRepository interface {
	Add(id string, name string, platformId string) error
	Exists(id string, platformId string) bool
	Search(searchParam string, platform string) ([]model.Game, error)
}
