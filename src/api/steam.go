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
	DbTrack      model.Subscription
	SteamApiGame model.SteamGame
	failed       bool
}

func (api *SteamApi) GetSteamIdFromVanityName(VanityName string) (string, error) {
	apiKey := config.GetSteamApiKey()
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUser/ResolveVanityURL/v0001/?key=%s&vanityurl=%s", apiKey, VanityName)
	body, err := api.getAndValidateRequest(url)
	if err != nil {
		return "", err
	}

	var apiResponse model.SteamApiVanityResponse
	if err = json.Unmarshal(body, &apiResponse); err != nil {
		return "", fmt.Errorf("error unmarshalling JSON: %w", err)
	}
	if apiResponse.Response.Success != 1 {
		return "", fmt.Errorf("cannot find steamId linked to that name")
	}

	return apiResponse.Response.SteamId, nil
}

func (api *SteamApi) FetchApiGamesPlayer(steamId string) (*model.SteamAPIAllResponse, error) {
	apiKey := config.GetSteamApiKey()
	url := fmt.Sprintf("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&format=json", apiKey, steamId)
	body, err := api.getAndValidateRequest(url)
	if err != nil {
		return nil, err
	}

	var apiResponse model.SteamAPIResponseAllGames
	if err = json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return &apiResponse.Response, nil
}

func (api *SteamApi) GetRequestedGamePlayedTime(steamId string, gameId string) (int, error) {
	apiResponse, err := api.FetchApiGamesPlayer(steamId)
	if err != nil {
		return 0, fmt.Errorf("error fetching steam api response: %w", err)
	}
	if apiResponse.GameCount <= 0 {
		return 0, fmt.Errorf("it seems you don't own any games")
	}
	gameIdInt, err := strconv.Atoi(gameId)
	if err != nil {
		return 0, fmt.Errorf("cannot convert game id to int: %w", err)
	}
	for _, game := range apiResponse.Games {
		if gameIdInt != game.AppID {
			continue // early return to only handle the game that is requested
		}
		return game.PlaytimeForever, nil
	}
	return 0, fmt.Errorf("Requested game" + gameId + " not found in Player " + steamId + "Played games list")
}

// FetchRecentGames We use this only because we want to parse less data. In reality, it could be useful to use FetchApiGamesPlayer as it does not matter about downtime longer than 2 weeks.
func (api *SteamApi) FetchRecentGames(steamId string) (*model.SteamApiResponse, error) {
	apiKey := config.GetSteamApiKey()
	url := fmt.Sprintf("http://api.steampowered.com/IPlayerService/GetRecentlyPlayedGames/v0001/?key=%s&steamid=%s&format=json", apiKey, steamId)
	body, err := api.getAndValidateRequest(url)
	if err != nil {
		return nil, err
	}

	var apiResponse model.SteamApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	return &apiResponse, nil
}

func (api *SteamApi) getAndValidateRequest(url string) ([]byte, error) {
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

	return body, nil

}

func (api *SteamApi) GetOnlyFailed(
	apiResponse *model.SteamApiResponse,
	trackedGamesByUser []model.Subscription,
) []MatchedDbGameToSteamGameInfo {
	// Create a map for quick lookup of tracked game IDs
	trackedGameMap := make(map[string]model.Subscription)
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

func closeBody(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		logger.Fail("Error closing response body: " + err.Error())
	}
}
