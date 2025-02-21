package service

import "quit4real.today/src/model"

type TrackerService interface {
	UpdateSteamTracker(steamId string, steamApiResponse *model.SteamApiResponse)
}
