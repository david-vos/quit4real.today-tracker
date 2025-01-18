package model

type SteamGame struct {
	AppID           int    `json:"appid"`
	Name            string `json:"name"`
	Playtime2Weeks  int    `json:"playtime_2weeks"`
	PlaytimeForever int    `json:"playtime_forever"`
	ImgIconURL      string `json:"img_icon_url"`
}

type SteamApiGetLastPlayed struct {
	TotalCount int         `json:"total_count"`
	Games      []SteamGame `json:"games"`
}

type SteamApiResponse struct {
	Response SteamApiGetLastPlayed `json:"response"`
}

type SteamAPIAllGame struct {
	AppID                  int `json:"appid"`
	PlaytimeForever        int `json:"playtime_forever"`
	PlaytimeWindowsForever int `json:"playtime_windows_forever"`
	PlaytimeMacForever     int `json:"playtime_mac_forever"`
	PlaytimeLinuxForever   int `json:"playtime_linux_forever"`
	PlaytimeDeckForever    int `json:"playtime_deck_forever"`
	RTimeLastPlayed        int `json:"rtime_last_played"`
	PlaytimeDisconnected   int `json:"playtime_disconnected"`
}

type SteamAPIAllResponse struct {
	GameCount int               `json:"game_count"`
	Games     []SteamAPIAllGame `json:"games"`
}

type SteamAPIResponseAllGames struct {
	Response SteamAPIAllResponse `json:"response"`
}
