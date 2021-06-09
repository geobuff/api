package main

import (
	"errors"
	"testing"

	"github.com/geobuff/api/repo"
	"github.com/geobuff/api/validation"
)

func TestMain(t *testing.T) {
	savedLoadConfig := loadConfig
	savedOpenConnection := repo.OpenConnection
	savedRunMigrations := runMigrations
	savedInit := validation.Init
	savedServe := serve

	defer func() {
		loadConfig = savedLoadConfig
		repo.OpenConnection = savedOpenConnection
		runMigrations = savedRunMigrations
		validation.Init = savedInit
		serve = savedServe
	}()

	tt := []struct {
		name           string
		loadConfig     func() error
		openConnection func() error
		runMigrations  func() error
		init           func() error
		serve          func() error
	}{
		{
			name:           "error on loadConfig",
			loadConfig:     func() error { return errors.New("test") },
			openConnection: repo.OpenConnection,
			runMigrations:  runMigrations,
			init:           validation.Init,
			serve:          serve,
		},
		{
			name:           "error on repo.OpenConnection",
			loadConfig:     func() error { return nil },
			openConnection: func() error { return errors.New("test") },
			runMigrations:  runMigrations,
			init:           validation.Init,
			serve:          serve,
		},
		{
			name:           "error on runMigrations",
			loadConfig:     func() error { return nil },
			openConnection: func() error { return nil },
			runMigrations:  func() error { return errors.New("test") },
			init:           validation.Init,
			serve:          serve,
		},
		{
			name:           "error on validation init",
			loadConfig:     func() error { return nil },
			openConnection: func() error { return nil },
			runMigrations:  func() error { return nil },
			init:           func() error { return errors.New("test") },
			serve:          serve,
		},
		{
			name:           "error on serve",
			loadConfig:     func() error { return nil },
			openConnection: func() error { return nil },
			runMigrations:  func() error { return nil },
			init:           func() error { return nil },
			serve:          func() error { return errors.New("test") },
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			loadConfig = tc.loadConfig
			repo.OpenConnection = tc.openConnection
			runMigrations = tc.runMigrations
			validation.Init = tc.init
			serve = tc.serve

			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected panic; got nil")
				}
			}()

			main()
		})
	}
}
