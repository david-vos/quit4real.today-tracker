package service

import "quit4real.today/src/model"

type SubscriptionService interface {
	GetOnlyFailedSteam(steamId string, apiResponse *model.SteamApiResponse) []model.MatchedDbGameToSteamGameInfo
}
