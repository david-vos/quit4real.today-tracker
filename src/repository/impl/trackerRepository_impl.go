package impl

import (
	"fmt"
	"quit4real.today/logger"
	"quit4real.today/src/model"
)

// trackerRepositoryImpl is the concrete implementation for TrackerRepository.
type trackerRepositoryImpl struct {
	DatabaseImpl *DatabaseImpl
}

// InsertTracker inserts a new row into the tracker table if it does not already exist (based on the composite key).
func (r *trackerRepositoryImpl) InsertTracker(tracker model.Tracker) error {
	query := `INSERT INTO tracker (user_id, game_id, day, time_played, new_total_time_played, amount_of_logins) 
              VALUES (?, ?, ?, ?, ?, ?)
              ON CONFLICT(user_id, game_id, day) DO NOTHING`
	// The `ON CONFLICT` clause ensures no duplicate records are inserted for the same (user_id, game_id, day).
	return r.DatabaseImpl.ExecuteQuery(query, tracker.UserID, tracker.GameID, tracker.Day, tracker.TimePlayed,
		tracker.NewTotalTimePlayed, tracker.AmountOfLogins)
}

// GetTrackerByID retrieves a tracker by its composite key (user_id, game_id, day).
func (r *trackerRepositoryImpl) GetTrackerByID(userID, gameID int, day string) (*model.Tracker, error) {
	query := `SELECT * FROM tracker WHERE user_id = ? AND game_id = ? AND day = ?`
	rows, err := r.DatabaseImpl.FetchRows(query, userID, gameID, day)
	if err != nil {
		return nil, fmt.Errorf("error fetching tracker rows: %w", err)
	}
	defer func() {
		if err := r.DatabaseImpl.CloseRows(rows); err != nil {
			logger.Fail("failed to close rows: " + err.Error())
		}
	}()

	if rows.Next() {
		tracker, err := model.MapTracker(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to map tracker row: %w", err)
		}
		return &tracker, nil
	}
	return nil, nil // Return nil if no rows are found
}

func (r *trackerRepositoryImpl) GetLatestTrackerByUserIdAndGameId(userID string, gameID int) (*model.Tracker, error) {
	query := `SELECT * FROM tracker WHERE user_id = ? AND game_id = ? ORDER BY day DESC LIMIT 1`
	rows, err := r.DatabaseImpl.FetchRows(query, userID, gameID)
	if err != nil {
		return nil, fmt.Errorf("error fetching tracker rows: %w", err)
	}
	defer func() {
		if err := r.DatabaseImpl.CloseRows(rows); err != nil {
			logger.Fail("failed to close rows: " + err.Error())
		}
	}()

	if rows.Next() {
		tracker, err := model.MapTracker(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to map tracker row: %w", err)
		}
		return &tracker, nil
	}
	return nil, nil // Return nil if no rows are found

}

// UpdateTracker updates an existing tracker by its composite key (user_id, game_id, day).
func (r *trackerRepositoryImpl) UpdateTracker(tracker *model.Tracker) error {
	query := `UPDATE tracker  SET time_played = ?, new_total_time_played = ?, amount_of_logins = ?
              WHERE user_id = ? AND game_id = ? AND day = ?`
	return r.DatabaseImpl.ExecuteQuery(query, tracker.TimePlayed, tracker.NewTotalTimePlayed, tracker.AmountOfLogins,
		tracker.UserID, tracker.GameID, tracker.Day)
}

// DeleteTracker deletes a tracker by its composite key (user_id, game_id, day).
func (r *trackerRepositoryImpl) DeleteTracker(userID, gameID int, day string) error {
	query := `DELETE FROM tracker WHERE user_id = ? AND game_id = ? AND day = ?`
	return r.DatabaseImpl.ExecuteQuery(query, userID, gameID, day)
}

// GetAllTrackers retrieves all tracker records.
func (r *trackerRepositoryImpl) GetAllTrackers() ([]model.Tracker, error) {
	query := `SELECT * FROM tracker`
	rows, err := r.DatabaseImpl.FetchRows(query)
	if err != nil {
		return nil, fmt.Errorf("error fetching tracker rows: %w", err)
	}
	defer func() {
		if err := r.DatabaseImpl.CloseRows(rows); err != nil {
			logger.Fail("failed to close rows: " + err.Error())
		}
	}()

	var trackers []model.Tracker
	for rows.Next() {
		tracker, err := model.MapTracker(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to map tracker row: %w", err)
		}
		trackers = append(trackers, tracker)
	}

	return trackers, nil
}
