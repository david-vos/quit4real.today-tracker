package config

import (
	"quit4real.today/logger"
)

func GetDBPath() string          { return "database.db" }
func GetDBMigrationPath() string { return "src/db/migrations" }
func GetSteamApiKey() string     { return "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" }
func InitLogger()                { logger.InitLogger("app.log") }
