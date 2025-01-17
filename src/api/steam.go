package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"quit4real.today/config"
	"quit4real.today/logger"
	"quit4real.today/src/model"
	"strconv"
)

type SteamApi struct {
}

type MatchedDbGameToSteamGameInfo struct {
	DbTrack      model.Tracker
	SteamApiGame model.SteamGame
	failed       bool
}

func (api *SteamApi) FetchApiData(steamId string) (*model.SteamApiResponse, error) {
	apiKey := config.GetSteamApiKey()
	url := fmt.Sprintf("http://api.steampowered.com/IPlayerService/GetRecentlyPlayedGames/v0001/?key=%s&steamid=%s&format=json", apiKey, steamId)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}
	defer closeBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var apiResponse model.SteamApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	return &apiResponse, nil
}

func closeBody(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		logger.Fail("Error closing response body: " + err.Error())
	}
}

func (api *SteamApi) GetOnlyFailed(
	apiResponse *model.SteamApiResponse,
	trackedGamesByUser []model.Tracker,
) []MatchedDbGameToSteamGameInfo {
	// Create a map for quick lookup of tracked game IDs
	trackedGameMap := make(map[string]model.Tracker)
	for _, trackedGame := range trackedGamesByUser {
		trackedGameMap[trackedGame.GameId] = trackedGame
	}

	var response []MatchedDbGameToSteamGameInfo

	for _, game := range apiResponse.Response.Games {
		gameID := strconv.Itoa(game.AppID)
		trackedGame, exists := trackedGameMap[gameID]
		if !exists {
			continue // Skip if the game is not tracked
		}

		// Check if the played amount is greater than the current playtime
		// This will always fail on the first iteration when a new user is created it does not actually set the
		// PlayedAmount in the beginning making it playtimeForever > 0 which is always true
		if game.PlaytimeForever > trackedGame.PlayedAmount {
			info := MatchedDbGameToSteamGameInfo{
				DbTrack:      trackedGame,
				SteamApiGame: game,
				failed:       true,
			}
			response = append(response, info)
			return response
		}
	}

	return []MatchedDbGameToSteamGameInfo{}
}
