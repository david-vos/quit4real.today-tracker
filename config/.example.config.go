package config

import (
	"quit4real.today/logger"
)

func GetDBPath() string {
	return "./db/database.db"
}
func GetSteamApiKey() string {
	return "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
}
func InitLogger() {
	logger.InitLogger()
}
