package model

type SteamGame struct {
	AppID           int    `json:"appid"`
	Name            string `json:"name"`
	Playtime2Weeks  int    `json:"playtime_2weeks"`
	PlaytimeForever int    `json:"playtime_forever"`
	ImgIconURL      string `json:"img_icon_url"`
}

type SteamResponse struct {
	TotalCount int         `json:"total_count"`
	Games      []SteamGame `json:"games"`
}

type SteamApiResponse struct {
	Response SteamResponse `json:"response"`
}
