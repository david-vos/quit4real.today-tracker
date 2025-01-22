package model

import "database/sql"

// UserPlatform represents a user's subscription to a platform.
type UserPlatform struct {
	ID             int      `json:"id"`               // Auto-incrementing ID
	User           User     `json:"user"`             // Embedded User object
	Platform       Platform `json:"platform"`         // Embedded Platform object
	PlatformUserID string   `json:"platform_user_id"` // User's ID on the specific platform
}

// MapUserPlatform maps SQL rows to a UserPlatform struct.
func MapUserPlatform(rows *sql.Rows) (UserPlatform, error) {
	var userPlatform UserPlatform
	var platformID string
	if err := rows.Scan(
		&userPlatform.ID,
		&userPlatform.User.ID,
		&platformID,
		&userPlatform.PlatformUserID,
	); err != nil {
		return UserPlatform{}, err
	}

	// Create a Platform object based on the platform ID
	userPlatform.Platform = Platform{ID: platformID}
	return userPlatform, nil
}
