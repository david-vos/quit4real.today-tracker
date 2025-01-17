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

		migrationPath := filepath.Join(migrationsPath, file.Name())
		content, err := os.ReadFile(migrationPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", file.Name(), err)
		}

		log.Printf("Applying migration: %s", file.Name())
		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to apply migration %s: %v", file.Name(), err)
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
