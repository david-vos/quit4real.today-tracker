package repository

import (
	"database/sql"
	"log"
	"project/models"
)

func GetAllUsers(db *sql.DB) []models.User {
	query := `SELECT * FROM users;`
	rows, err := db.Query(query)
	if err != nil {
		return nil
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.SteamId, &user.ApiKey)
		if err != nil {
			return nil
		}
		users = append(users, user)
	}

	return users
}

func CreateUser(db DBExecutor, user models.User) error {
	query := "INSERT INTO users (id, name, steam_id, api_key) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, user.ID, user.Name, user.SteamId, user.ApiKey)
	return err
}

func GetTracker(db *sql.DB, steamId string, gameId string) (models.Tracker, error) {
	query := "SELECT * FROM user_tracker WHERE steam_id = ? AND game_id = ?;"
	rows, err := db.Query(query, steamId, gameId)
	if err != nil {
		log.Fatal(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)
	var tracker models.Tracker
	err = rows.Scan(
		&tracker.ID,
		&tracker.SteamId,
		&tracker.GameId,
		&tracker.PlayedAmount)
	if err != nil {
		return models.Tracker{}, err
	}
	return tracker, err
}

func CreateTracker(db DBExecutor, steamId string, gameId string) error {
	query := "INSERT INTO user_tracker (steam_id, game_id, played_amount) VALUES (?, ?, ?);"
	_, err := db.Exec(query, steamId, gameId, 0)
	return err
}

func UpdateTracker(db DBExecutor, steamId string, gameId string, playedAmount int) error {
	query := "UPDATE user_tracker SET played_amount = ? WHERE steam_id = ? AND game_id = ?;"
	_, err := db.Exec(query, playedAmount, steamId, gameId)
	return err
}

func GetUserTracker(db *sql.DB, steamId string) ([]models.Tracker, error) {
	query := "SELECT * FROM user_tracker WHERE steam_id = ?"
	rows, err := db.Query(query, steamId)
	if err != nil {
		return nil, err
	}
	var trackers []models.Tracker
	for rows.Next() {
		var tracker models.Tracker
		var err = rows.Scan(
			&tracker.ID,
			&tracker.SteamId,
			&tracker.GameId,
			&tracker.PlayedAmount)
		if err != nil {
			return nil, err
		}
		trackers = append(trackers, tracker)
	}
	return trackers, err
}
