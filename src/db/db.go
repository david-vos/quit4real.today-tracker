package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
	"quit4real.today/config"
	"quit4real.today/logger"
	"sort"
)

func Connect(dbPath string) (*sql.DB, error) {
	// Create the database file if it doesn't exist
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			return nil, err
		}
		err = file.Close()
		if err != nil {
			return nil, err
		}
		logger.Debug("Database file created:" + dbPath)
	}

	// Connect to SQLite
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	logger.Debug("Database connected successfully")
	return db, nil
}

func ApplyMigrations(db *sql.DB, migrationsPath string) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS migrations (id TEXT PRIMARY KEY);`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %v", err)
	}

	// Sort migration files by name to ensure they run in order
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		migrationName := file.Name()
		var exists string
		err := db.QueryRow(`SELECT id FROM migrations WHERE id = ?`, migrationName).Scan(&exists)
		if err == nil {
			log.Printf("Skipping already applied migration: %s", migrationName)
			continue
		}

		migrationPath := filepath.Join(migrationsPath, migrationName)
		content, err := os.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", migrationName, err)
		}

		log.Printf("Applying migration: %s", migrationName)
		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to apply migration %s: %v", migrationName, err)
		}

		_, err = db.Exec(`INSERT INTO migrations (id) VALUES (?);`, migrationName)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %v", migrationName, err)
		}
	}

	log.Println("Migrations applied successfully!")
	return nil
}

func Setup() *sql.DB {
	database, err := Connect(config.GetDBPath())
	if err != nil {
		logger.Fail("Failed to connect to the database: " + err.Error())
	}

	// Apply migrations
	err = ApplyMigrations(database, config.GetDBMigrationPath())
	if err != nil {
		logger.Fail("Failed to apply migrations: " + err.Error())
	}

	return database
}
