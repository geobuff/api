package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

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

			_, err = ioutil.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}
		})
	}
}
