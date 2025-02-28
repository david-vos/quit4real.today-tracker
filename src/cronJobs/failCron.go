package cronJobs

import (
	"github.com/robfig/cron/v3"
	"quit4real.today/logger"
	"quit4real.today/src/handler/command"
	"quit4real.today/src/handler/query"
	"quit4real.today/src/handler/service"
	"quit4real.today/src/model"
)

type FailCronImpl struct {
	SubscriptionCommandHandler *command.SubscriptionCommandHandlerImpl
	SubscriptionQueryHandler   *query.SubscriptionQueryHandlerImpl
	UserQueryHandler           *query.UserQueryHandlerImpl
	SteamService               service.SteamService
	SubscriptionService        service.SubscriptionService
	TrackerService             service.TrackerService
}

type UserWithRecentPlayed struct {
	UserIdFrom        UserFrom
	RecentPlayedGames *model.SteamApiResponse
}

type UserFrom struct {
	UserId string
	From   []string
}

func (fc *FailCronImpl) Start() {
	// If in the future I will add another cron job then it's better to have a simulacra structure to how the Endpoints work
	cronJob := cron.New()
	// ((24*60)/10)*694 ~= 100.000 the STEAM API limit
	// 694 -> max amount of users :thinking per API key
	// Currently set to 1 minute for dev, I will increase this back to 10 minutes when going live
	_, err := cronJob.AddFunc("@every 1m", func() {
		fc.updateAndSendNotifySteam()
	})
	if err != nil {
		logger.Fail("Error adding cron job: " + err.Error())
		return
	}

	cronJob.Start()
}

func (fc *FailCronImpl) updateAndSendNotifySteam() {
	usersWithRecentPlayed, err := fc.getUserWithRecentPlayed()
	if err != nil {
		logger.Fail("Error getting users with the recent played game: " + err.Error())
		return
	}

	// Iterate over each unique platformUserId and update
	// This should all be abstracted away so you don't need a different service with a different platform
	for _, userWithRecentPlayed := range usersWithRecentPlayed {
		// Check if "subscription" is in the From list
		for _, from := range userWithRecentPlayed.UserIdFrom.From {
			if from == "subscription" {
				failedGamesByUser := fc.SubscriptionService.GetOnlyFailedSteam(
					userWithRecentPlayed.UserIdFrom.UserId, userWithRecentPlayed.RecentPlayedGames)
				fc.SubscriptionCommandHandler.UpdateSubscriptions(userWithRecentPlayed.UserIdFrom.UserId, failedGamesByUser)
			}

			if from == "user" {
				errs := fc.TrackerService.UpdateSteamTrackers(userWithRecentPlayed.UserIdFrom.UserId, userWithRecentPlayed.RecentPlayedGames)
				for _, err := range errs {
					if err != nil {
						logger.Fail("Error updating steam tracker: " + err.Error())
					}
				}
			}
		}
	}
}

// getUserWithRecentPlayed We do this in order to not have duplicate queries going out to the steamAPI
// I don't want to limit
func (fc *FailCronImpl) getUserWithRecentPlayed() ([]UserWithRecentPlayed, error) {
	allSteamSubscriptions, err := fc.SubscriptionQueryHandler.GetAllSteam()
	if err != nil {
		return nil, err
	}
	allSteamUsers, err := fc.UserQueryHandler.GetAllSteamVerified()
	if err != nil {
		return nil, err
	}
	userMap := make(map[string]*UserWithRecentPlayed)

	// Populate the map with all Steam subscriptions
	for _, subscription := range allSteamSubscriptions {
		fc.addOrUpdateUser(subscription.PlatFormUserId, "subscription", userMap)
	}

	// Populate the map with all Steam verified users
	for _, user := range allSteamUsers {
		fc.addOrUpdateUser(user.SteamID, "user", userMap)
	}

	// Convert map values to a slice
	var usersWithRecentPlayed []UserWithRecentPlayed
	for _, user := range userMap {
		usersWithRecentPlayed = append(usersWithRecentPlayed, *user)
	}

	// Fetch recent games
	for i := 0; i < len(usersWithRecentPlayed); i++ {
		apiResponse, err := fc.SteamService.FetchRecentGames(usersWithRecentPlayed[i].UserIdFrom.UserId)
		if err != nil {
			logger.Fail("failed to fetch recent games: " + err.Error())
			// Remove the failed one from the list
			usersWithRecentPlayed = append(usersWithRecentPlayed[:i], usersWithRecentPlayed[i+1:]...)
			i-- //  index to account for removed element this messed me up a decent bit...
			continue
		}
		usersWithRecentPlayed[i].RecentPlayedGames = apiResponse
	}

	return usersWithRecentPlayed, nil
}

// Helper function to add or update entries in the map
func (fc *FailCronImpl) addOrUpdateUser(userId, from string, userMap map[string]*UserWithRecentPlayed) {
	if entry, exists := userMap[userId]; exists {
		entry.UserIdFrom.From = append(entry.UserIdFrom.From, from)
	} else {
		userMap[userId] = &UserWithRecentPlayed{
			UserIdFrom: UserFrom{
				UserId: userId,
				From:   []string{from},
			},
		}
	}
}
