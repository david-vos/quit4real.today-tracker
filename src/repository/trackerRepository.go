package repository

import "quit4real.today/src/model"

type TrackerRepository interface {
	InsertTracker(tracker model.Tracker) error
	GetTrackerByID(userID, gameID int, day string) (*model.Tracker, error)
	UpdateTracker(tracker model.Tracker) error
	DeleteTracker(userID, gameID int, day string) error
	GetAllTrackers() ([]model.Tracker, error)
}
