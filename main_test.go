package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func TestGetRoutes(t *testing.T) {
	tt := []struct {
		name   string
		route  string
		status int
	}{
		{name: "get users", route: "users", status: http.StatusUnauthorized},
		{name: "get user 1", route: "users/1", status: http.StatusUnauthorized},
		{name: "get scores for user 1", route: "scores/1", status: http.StatusUnauthorized},
		{name: "get isocodes", route: "isocodes", status: http.StatusOK},
		{name: "get world countries", route: "world/countries", status: http.StatusOK},
		{name: "get world country alternatives", route: "world/countries/alternatives", status: http.StatusOK},
		{name: "get world country prefixes", route: "world/countries/prefixes", status: http.StatusOK},
		{name: "get world country map", route: "world/countries/map", status: http.StatusOK},
		{name: "get world leaderboard missing page", route: "world/leaderboard", status: http.StatusBadRequest},
	}

	server := httptest.NewServer(router())
	defer server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			response, err := http.Get(fmt.Sprintf("%s/api/%s", server.URL, tc.route))
			if err != nil {
				t.Fatalf("could not send GET request: %v", err)
			}
			defer response.Body.Close()

			if response.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, response.Status)
			}

			_, err = ioutil.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}
		})
	}
}

func TestPostRoutes(t *testing.T) {
	tt := []struct {
		name   string
		route  string
		body   string
		status int
	}{
		{name: "create world leaderboard entry", route: "world/leaderboard", body: "", status: http.StatusUnauthorized},
	}

	server := httptest.NewServer(router())
	defer server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			jsonString, err := json.Marshal(tc.body)
			if err != nil {
				t.Error("failed to parse test body")
			}

			response, err := http.Post(fmt.Sprintf("%s/api/%s", server.URL, tc.route), "application/json", bytes.NewReader(jsonString))
			if err != nil {
				t.Fatalf("could not send POST request: %v", err)
			}
			defer response.Body.Close()

			if response.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, response.Status)
			}
		})
	}
}

func TestPutRoutes(t *testing.T) {
	tt := []struct {
		name   string
		route  string
		body   string
		status int
	}{
		{name: "upsert score", route: "scores", body: "", status: http.StatusUnauthorized},
		{name: "update world leaderboard entry", route: "world/leaderboard/1", body: "", status: http.StatusUnauthorized},
	}

	server := httptest.NewServer(router())
	defer server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			jsonString, err := json.Marshal(tc.body)
			if err != nil {
				t.Error("failed to parse test body")
			}

			request, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/%s", server.URL, tc.route), bytes.NewReader(jsonString))
			if err != nil {
				t.Fatalf("could not create PUT request: %v", err)
			}

			response, err := server.Client().Do(request)
			if err != nil {
				t.Fatalf("could not send PUT request: %v", err)
			}
			defer response.Body.Close()

			if response.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, response.Status)
			}
		})
	}
}

func TestDeleteRoutes(t *testing.T) {
	tt := []struct {
		name   string
		route  string
		status int
	}{
		{name: "delete user", route: "users/1", status: http.StatusUnauthorized},
		{name: "delete score", route: "scores/1", status: http.StatusUnauthorized},
		{name: "delete world leaderboard entry", route: "world/leaderboard/1", status: http.StatusUnauthorized},
	}

	server := httptest.NewServer(router())
	defer server.Close()

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/%s", server.URL, tc.route), nil)
			if err != nil {
				t.Fatalf("could not create DELETE request: %v", err)
			}

			response, err := server.Client().Do(request)
			if err != nil {
				t.Fatalf("could not send DELETE request: %v", err)
			}
			defer response.Body.Close()

			if response.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, response.Status)
			}
		})
	}
}
