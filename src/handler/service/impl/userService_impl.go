package impl

import (
	"quit4real.today/logger"
	"quit4real.today/src/handler/service"
)

type UserServiceImpl struct {
	TrackerService service.TrackerService
	SteamService   service.SteamService
}

func (service *UserServiceImpl) UpdateUserTrackers(platformUserId string) {
	// Still need to think about how I want to know which platform a userId is from
	// So there kind of needs to be logic that handles this, but for now I assume it is just steam.
	// The other services/endpoints already do this anyway :shrug:
	service.UpdateUserTrackersSteam(platformUserId)
}

func (service *UserServiceImpl) UpdateUserTrackersSteam(steamId string) {
	apiResponse, err := service.SteamService.FetchRecentGames(steamId)
	if err != nil {
		logger.Fail("Failed to fetch Steam API games for user: " + steamId + " | ERROR: " + err.Error())
	}

	errs := service.TrackerService.UpdateSteamTrackers(steamId, apiResponse)
	if errs != nil && len(errs) > 0 {
		for _, err := range errs {
			logger.Fail("Failed to update trackers for user: " + steamId + " | ERROR: " + err.Error())
		}
	}
}

func (service *UserServiceImpl) CreateUserTrackers(platformUserId string) {
	// Still need to think about how I want to know which platform a userId is from
	// So there kind of needs to be logic that handles this, but for now I assume it is just steam.
	// The other services/endpoints already do this anyway :shrug:
	service.CreateUserTrackersSteam(platformUserId)
}

func (service *UserServiceImpl) CreateUserTrackersSteam(steamId string) {
	SteamService, err := service.SteamService.FetchApiGamesPlayer(steamId)
	if err != nil {
		logger.Fail("Failed to fetch Steam API games for user: " + steamId + " | ERROR: " + err.Error())
	}

	err = service.TrackerService.CreateSteamTrackers(steamId, SteamService)
	if err != nil {
		logger.Fail("Failed to create trackers for user: " + steamId + " | ERROR: " + err.Error())
	}
}
