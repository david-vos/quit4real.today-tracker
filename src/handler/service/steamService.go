package service

import "quit4real.today/src/model"

type SteamService interface {
	GetSteamIdFromVanityName(vanityName string) (string, error)
	FetchApiGamesPlayer(steamId string) (*model.SteamAPIAllResponse, error)
	GetRequestedGame(steamId string, gameId string) (model.SteamAPIAllGame, error)
	FetchRecentGames(steamId string) (*model.SteamApiResponse, error)
	FetchUserInfo(steamId string) (*model.SteamApiUserInfo, error)
	GetOnlyFailed(apiResponse *model.SteamApiResponse, trackedGamesByUser []model.Subscription) []model.MatchedDbGameToSteamGameInfo
}
