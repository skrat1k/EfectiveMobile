package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

func RunMigrations(connectionURL string) error {
	m, err := migrate.New("fille://migrations", connectionURL)
	if err != nil {
		return fmt.Errorf("failed to create migrations:%w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations:%w", err)
	}
	return nil
}
