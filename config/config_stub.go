//go:build dependabot

package config

import (
	"os"
	"quit4real.today/logger"
)

func GetDBPath() string {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "/app/data/database.db" // Default if not set
	}
	return dbPath
}

func GetDBMigrationPath() string {
	return "/app/src/db/migrations" // Migrations remain unchanged
}
func GetSteamApiKey() string { return "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" }
func InitLogger()            { logger.InitLogger("/app/logs/app.log") }

func JwtSecret() []byte { return []byte("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx") }

func BackendUrl() string { return "https://tracker.quit4real.today/" }

func FrontendUrl() string { return "https://quit4real.today/" }
