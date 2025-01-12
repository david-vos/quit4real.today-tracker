package models

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	SteamId string `json:"steam_id"`
	ApiKey  string `json:"api_key"`
}
