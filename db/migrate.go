package migrate

import (
	"fmt"

	gomigrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(dbURL, migrationFilesPath string) error {
	m, err := gomigrate.New(migrationFilesPath, dbURL)
	if err != nil {
		return fmt.Errorf("error while initializing migration: %s", err)
	}
	if err := m.Up(); err != nil {
		if err != gomigrate.ErrNoChange {
			return fmt.Errorf("error while migrating: %s", err)
		}
	}
	return nil
}
