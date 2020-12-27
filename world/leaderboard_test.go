package world

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geobuff/geobuff-api/database"
	"github.com/gorilla/mux"
)

func TestGetEntries(t *testing.T) {
	savedGetWorldLeaderboardEntries := database.GetWorldLeaderboardEntries
	savedGetWorldLeaderboardEntryID := database.GetWorldLeaderboardEntryID

	defer func() {
		database.GetWorldLeaderboardEntries = savedGetWorldLeaderboardEntries
		database.GetWorldLeaderboardEntryID = savedGetWorldLeaderboardEntryID
	}()

	tt := []struct {
		name                       string
		getWorldLeaderboardEntries func(limit int, offset int) ([]database.WorldLeaderboardEntry, error)
		getWorldLeaderboardEntryID func(limit int, offset int) (int, error)
		page                       string
		status                     int
		hasMore                    bool
	}{
		{
			name:                       "invalid page parameter",
			getWorldLeaderboardEntries: database.GetWorldLeaderboardEntries,
			getWorldLeaderboardEntryID: database.GetWorldLeaderboardEntryID,
			page:                       "testing",
			status:                     http.StatusBadRequest,
			hasMore:                    false,
		},
		{
			name:                       "valid page, error on GetWorldLeaderboardEntries",
			getWorldLeaderboardEntries: func(limit int, offset int) ([]database.WorldLeaderboardEntry, error) { return nil, errors.New("test") },
			getWorldLeaderboardEntryID: database.GetWorldLeaderboardEntryID,
			page:                       "0",
			status:                     http.StatusInternalServerError,
			hasMore:                    false,
		},
		{
			name: "valid page, error on GetWorldLeaderboardEntryID",
			getWorldLeaderboardEntries: func(limit int, offset int) ([]database.WorldLeaderboardEntry, error) {
				return []database.WorldLeaderboardEntry{}, nil
			},
			getWorldLeaderboardEntryID: func(limit int, offset int) (int, error) { return 0, errors.New("test") },
			page:                       "0",
			status:                     http.StatusInternalServerError,
			hasMore:                    false,
		},
		{
			name: "happy path, has more is false",
			getWorldLeaderboardEntries: func(limit int, offset int) ([]database.WorldLeaderboardEntry, error) {
				return []database.WorldLeaderboardEntry{}, nil
			},
			getWorldLeaderboardEntryID: func(limit int, offset int) (int, error) { return 0, sql.ErrNoRows },
			page:                       "0",
			status:                     http.StatusOK,
			hasMore:                    false,
		},
		{
			name: "happy path, has more is true",
			getWorldLeaderboardEntries: func(limit int, offset int) ([]database.WorldLeaderboardEntry, error) {
				return []database.WorldLeaderboardEntry{}, nil
			},
			getWorldLeaderboardEntryID: func(limit int, offset int) (int, error) { return 1, nil },
			page:                       "0",
			status:                     http.StatusOK,
			hasMore:                    true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetWorldLeaderboardEntries = tc.getWorldLeaderboardEntries
			database.GetWorldLeaderboardEntryID = tc.getWorldLeaderboardEntryID

			request, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/api/world/leaderboard?page=%v", tc.page), nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			writer := httptest.NewRecorder()
			GetEntries(writer, request)
			result := writer.Result()
			defer result.Body.Close()

			if result.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, result.StatusCode)
			}

			if tc.status == http.StatusOK {
				body, err := ioutil.ReadAll(result.Body)
				if err != nil {
					t.Fatalf("could not read response: %v", err)
				}

				var parsed EntriesDto
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}

				if parsed.HasMore != tc.hasMore {
					t.Errorf("expected hasMore = %v; got: %v", tc.hasMore, parsed.HasMore)
				}
			}
		})
	}
}

func TestGetEntry(t *testing.T) {
	savedGetWorldLeaderboardEntry := database.GetWorldLeaderboardEntry

	defer func() {
		database.GetWorldLeaderboardEntry = savedGetWorldLeaderboardEntry
	}()

	tt := []struct {
		name                     string
		getWorldLeaderboardEntry func(userID int) (database.WorldLeaderboardEntry, error)
		userID                   string
		status                   int
	}{
		{
			name:                     "invalid userId",
			getWorldLeaderboardEntry: database.GetWorldLeaderboardEntry,
			userID:                   "testing",
			status:                   http.StatusBadRequest,
		},
		{
			name: "valid userId, entry not found",
			getWorldLeaderboardEntry: func(userID int) (database.WorldLeaderboardEntry, error) {
				return database.WorldLeaderboardEntry{}, sql.ErrNoRows
			},
			userID: "1",
			status: http.StatusNotFound,
		},
		{
			name: "valid userId, error on GetWorldLeaderboardEntry",
			getWorldLeaderboardEntry: func(userID int) (database.WorldLeaderboardEntry, error) {
				return database.WorldLeaderboardEntry{}, errors.New("test")
			},
			userID: "1",
			status: http.StatusInternalServerError,
		},
		{
			name: "happy path",
			getWorldLeaderboardEntry: func(userID int) (database.WorldLeaderboardEntry, error) {
				return database.WorldLeaderboardEntry{}, nil
			},
			userID: "1",
			status: http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetWorldLeaderboardEntry = tc.getWorldLeaderboardEntry

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"userId": tc.userID,
			})

			writer := httptest.NewRecorder()
			GetEntry(writer, request)
			result := writer.Result()
			defer result.Body.Close()

			if result.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, result.StatusCode)
			}

			if tc.status == http.StatusOK {
				body, err := ioutil.ReadAll(result.Body)
				if err != nil {
					t.Fatalf("could not read response: %v", err)
				}

				var parsed database.WorldLeaderboardEntry
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}
