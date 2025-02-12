package service

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

type SteamService struct{}

// MatchedDbGameToSteamGameInfo holds the mapping between a database tracked game and a Steam game.
type MatchedDbGameToSteamGameInfo struct {
	DbTrack      model.Subscription
	SteamApiGame model.SteamGame
	Failed       bool
}

// GetSteamIdFromVanityName resolves a vanity name to a Steam ID.
func (api *SteamService) GetSteamIdFromVanityName(vanityName string) (string, error) {
	apiKey := config.GetSteamApiKey()
	url := fmt.Sprintf("https://api.steampowered.com/ISteamUser/ResolveVanityURL/v0001/?key=%s&vanityurl=%s", apiKey, vanityName)
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

// FetchApiGamesPlayer retrieves all games owned by a player.
func (api *SteamService) FetchApiGamesPlayer(steamId string) (*model.SteamAPIAllResponse, error) {
	apiKey := config.GetSteamApiKey()
	url := fmt.Sprintf("http://api.steampowered.com/IPlayerService/GetOwnedGames/v0001/?key=%s&steamid=%s&include_appinfo=true&format=json", apiKey, steamId)
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

// GetRequestedGame  retrieves the game owned by a specific player.
func (api *SteamService) GetRequestedGame(steamId string, gameId string) (model.SteamAPIAllGame, error) {
	apiResponse, err := api.FetchApiGamesPlayer(steamId)
	if err != nil {
		return model.SteamAPIAllGame{}, fmt.Errorf("error fetching steam api response: %w", err)
	}
	if apiResponse.GameCount <= 0 {
		return model.SteamAPIAllGame{}, fmt.Errorf("it seems you don't own any games")
	}

	gameIdInt, err := strconv.Atoi(gameId)
	if err != nil {
		return model.SteamAPIAllGame{}, fmt.Errorf("cannot convert game id to int: %w", err)
	}

	for _, game := range apiResponse.Games {
		if gameIdInt == game.Appid {
			return game, nil
		}
	}
	return model.SteamAPIAllGame{}, fmt.Errorf("requested game %s not found in player %s's played games list", gameId, steamId)
}

// FetchRecentGames retrieves the recently played games for a player.
func (api *SteamService) FetchRecentGames(steamId string) (*model.SteamApiResponse, error) {
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

// getAndValidateRequest performs an HTTP GET request and validates the response.
func (api *SteamService) getAndValidateRequest(url string) ([]byte, error) {
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

// GetOnlyFailed filters out the games that have failed based on the API response and tracked games.
func (api *SteamService) GetOnlyFailed(
	apiResponse *model.SteamApiResponse,
	trackedGamesByUser []model.Subscription,
) []MatchedDbGameToSteamGameInfo {
	// Create a map for quick lookup of tracked game IDs
	trackedGameMap := make(map[string]model.Subscription)
	for _, trackedGame := range trackedGamesByUser {
		trackedGameMap[trackedGame.PlatformGameId] = trackedGame
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
				Failed:       true,
			}
			response = append(response, info)
		}
	}

	return response
}

// closeBody closes the response body and logs any errors.
func closeBody(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		logger.Fail("Error closing response body: " + err.Error())
	}
}
