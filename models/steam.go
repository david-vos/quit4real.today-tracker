package models

// Define structs to match the JSON response structure
type Game struct {
	AppID           int    `json:"appid"`
	Name            string `json:"name"`
	Playtime2Weeks  int    `json:"playtime_2weeks"`
	PlaytimeForever int    `json:"playtime_forever"`
	ImgIconURL      string `json:"img_icon_url"`
}

type Response struct {
	TotalCount int    `json:"total_count"`
	Games      []Game `json:"games"`
}

type ApiResponse struct {
	Response Response `json:"response"`
}
