package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
)

const DEFAULT_MIGRATION_DIR = "./db/migrate/"

func getMigrationDir() string {
	value := os.Getenv("PHTS_MIGRATION_DIR")
	if len(value) == 0 {
		value = DEFAULT_MIGRATION_DIR
	}

	value, err := filepath.Abs(value)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("file://%s", value)
}

func ApplyMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	migrationsPath := getMigrationDir()
	log.Printf("Loading migrations from %s", migrationsPath)
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		log.Printf("Error while creating migration: %s", err.Error())
		return err
	}

	err = m.Up()
	if err == migrate.ErrNoChange {
		log.Println("Database up to date!")
	} else if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Database migrated!")
	}

	return nil
}
