package main

import (
	"errors"
	"testing"

	"github.com/geobuff/api/repo"
)

func TestMain(t *testing.T) {
	savedOpenConnection := repo.OpenConnection
	savedServe := serve

	defer func() {
		repo.OpenConnection = savedOpenConnection
		serve = savedServe
	}()

	tt := []struct {
		name           string
		openConnection func() error
		serve          func() error
	}{
		{
			name:           "error on repo.OpenConnection",
			openConnection: func() error { return errors.New("test") },
			serve:          serve,
		},
		{
			name:           "error on serve",
			openConnection: func() error { return nil },
			serve:          func() error { return errors.New("test") },
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.OpenConnection = tc.openConnection
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
