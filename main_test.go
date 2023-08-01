package main

import (
	"errors"
	"testing"

	"github.com/geobuff/api/repo"
)

func TestMain(t *testing.T) {
	savedLoadConfig := loadConfig
	savedOpenConnection := repo.OpenConnection
	savedRunMigrations := runMigrations

	defer func() {
		loadConfig = savedLoadConfig
		repo.OpenConnection = savedOpenConnection
		runMigrations = savedRunMigrations
	}()

	tt := []struct {
		name           string
		loadConfig     func() error
		openConnection func() error
		runMigrations  func() error
	}{
		{
			name:           "error on loadConfig",
			loadConfig:     func() error { return errors.New("test") },
			openConnection: repo.OpenConnection,
			runMigrations:  runMigrations,
		},
		{
			name:           "error on repo.OpenConnection",
			loadConfig:     func() error { return nil },
			openConnection: func() error { return errors.New("test") },
			runMigrations:  runMigrations,
		},
		{
			name:           "error on runMigrations",
			loadConfig:     func() error { return nil },
			openConnection: func() error { return nil },
			runMigrations:  func() error { return errors.New("test") },
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			loadConfig = tc.loadConfig
			repo.OpenConnection = tc.openConnection
			runMigrations = tc.runMigrations

			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected panic; got nil")
				}
			}()

			main()
		})
	}
}
