package models

type Tracker struct {
	ID           int    `json:"id"`
	SteamId      string `json:"steam_id"`
	GameId       string `json:"api_key"`
	PlayedAmount int    `json:"played_amount"`
}
