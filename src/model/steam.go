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
	Appid                    int    `json:"appid"`
	Name                     string `json:"name"`
	PlaytimeForever          int    `json:"playtime_forever"`
	ImgIconUrl               string `json:"img_icon_url"`
	HasCommunityVisibleStats bool   `json:"has_community_visible_stats"`
	PlaytimeWindowsForever   int    `json:"playtime_windows_forever"`
	PlaytimeMacForever       int    `json:"playtime_mac_forever"`
	PlaytimeLinuxForever     int    `json:"playtime_linux_forever"`
	PlaytimeDeckForever      int    `json:"playtime_deck_forever"`
	RtimeLastPlayed          int    `json:"rtime_last_played"`
	HasLeaderboards          bool   `json:"has_leaderboards"`
	PlaytimeDisconnected     int    `json:"playtime_disconnected"`
}

type SteamAPIAllResponse struct {
	GameCount int               `json:"game_count"`
	Games     []SteamAPIAllGame `json:"games"`
}

type SteamAPIResponseAllGames struct {
	Response SteamAPIAllResponse `json:"response"`
}

type SteamApiVanityResponse struct {
	Response SteamApiVanity `json:"response"`
}

type SteamApiVanity struct {
	SteamId string `json:"steamid"`
	Success int    `json:"success"`
}
