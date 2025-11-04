package migrations

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

// MigrateUp applies all up migrations from dir to databaseURL.
func MigrateUp(migrationsDir, databaseURL string, log *zap.Logger) error {
	if migrationsDir == "" {
		return fmt.Errorf("migrations dir is empty")
	}
	m, err := migrate.New("file://"+migrationsDir, databaseURL)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	log.Info("migrations.applied")
	return nil
}

// MigrateDown rolls back last migration.
func MigrateDown(migrationsDir, databaseURL string, log *zap.Logger) error {
	m, err := migrate.New("file://"+migrationsDir, databaseURL)
	if err != nil {
		return err
	}
	if err := m.Steps(-1); err != nil {
		return err
	}
	log.Info("migration.rolledback")
	return nil
}
