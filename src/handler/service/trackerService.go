package service

import "quit4real.today/src/model"

type TrackerService interface {
	UpdateSteamTrackers(steamId string, steamApiResponse *model.SteamApiResponse)
	CreateSteamTrackers(steamId string, allGames *model.SteamAPIAllResponse) error
}
