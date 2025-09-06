package migrations

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const MIGRATION_DIR = "file://internal/database/migrations"

func RunMigrations(connectionString string) error {
	m, err := migrate.New(
		MIGRATION_DIR,
		connectionString,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Migrations applied successfully!")
	return nil
}

func RollbackMigrations(connectionString string) error {
	m, err := migrate.New(
		MIGRATION_DIR,
		connectionString,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Migrations rolled back successfully!")
	return nil
}

func ForceDBVersion(connectionString string, version int) error {
	m, err := migrate.New(
		MIGRATION_DIR,
		connectionString,
	)

	if err != nil {
		return fmt.Errorf("failed to force version: %w", err)
	}

	if err := m.Force(version); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("force version failed: %w", err)
	}

	log.Println("Version forced successfully!")
	return nil
}
