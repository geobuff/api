package main

import (
	"errors"
	"testing"

	"github.com/geobuff/geobuff-api/config"
	"github.com/geobuff/geobuff-api/database"
)

func TestMain(t *testing.T) {
	savedLoad := config.Load
	savedOpenConnection := database.OpenConnection
	savedServe := serve

	defer func() {
		config.Load = savedLoad
		database.OpenConnection = savedOpenConnection
		serve = savedServe
	}()

	tt := []struct {
		name           string
		load           func(fileName string) error
		openConnection func() error
		serve          func() error
	}{
		{
			name:           "error on config.Load",
			load:           func(fileName string) error { return errors.New("test") },
			openConnection: database.OpenConnection,
			serve:          serve,
		},
		{
			name:           "error on database.OpenConnection",
			load:           func(fileName string) error { return nil },
			openConnection: func() error { return errors.New("test") },
			serve:          serve,
		},
		{
			name:           "error on serve",
			load:           func(fileName string) error { return nil },
			openConnection: func() error { return nil },
			serve:          func() error { return errors.New("test") },
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			config.Load = tc.load
			database.OpenConnection = tc.openConnection
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
