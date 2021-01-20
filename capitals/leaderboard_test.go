package capitals

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geobuff/api/config"
	"github.com/geobuff/api/database"
	"github.com/geobuff/auth"
	"github.com/gorilla/mux"
)

func TestGetEntries(t *testing.T) {
	savedGetCapitalsLeaderboardEntries := database.GetCapitalsLeaderboardEntries
	savedGetCapitalsLeaderboardEntryID := database.GetCapitalsLeaderboardEntryID

	defer func() {
		database.GetCapitalsLeaderboardEntries = savedGetCapitalsLeaderboardEntries
		database.GetCapitalsLeaderboardEntryID = savedGetCapitalsLeaderboardEntryID
	}()

	tt := []struct {
		name                          string
		getCapitalsLeaderboardEntries func(limit int, offset int) ([]database.CapitalsLeaderboardEntry, error)
		getCapitalsLeaderboardEntryID func(limit int, offset int) (int, error)
		page                          string
		status                        int
		hasMore                       bool
	}{
		{
			name:                          "invalid page parameter",
			getCapitalsLeaderboardEntries: database.GetCapitalsLeaderboardEntries,
			getCapitalsLeaderboardEntryID: database.GetCapitalsLeaderboardEntryID,
			page:                          "testing",
			status:                        http.StatusBadRequest,
			hasMore:                       false,
		},
		{
			name: "valid page, error on GetCapitalsLeaderboardEntries",
			getCapitalsLeaderboardEntries: func(limit int, offset int) ([]database.CapitalsLeaderboardEntry, error) {
				return nil, errors.New("test")
			},
			getCapitalsLeaderboardEntryID: database.GetCapitalsLeaderboardEntryID,
			page:                          "0",
			status:                        http.StatusInternalServerError,
			hasMore:                       false,
		},
		{
			name: "valid page, error on GetCapitalsLeaderboardEntryID",
			getCapitalsLeaderboardEntries: func(limit int, offset int) ([]database.CapitalsLeaderboardEntry, error) {
				return []database.CapitalsLeaderboardEntry{}, nil
			},
			getCapitalsLeaderboardEntryID: func(limit int, offset int) (int, error) { return 0, errors.New("test") },
			page:                          "0",
			status:                        http.StatusInternalServerError,
			hasMore:                       false,
		},
		{
			name: "happy path, has more is false",
			getCapitalsLeaderboardEntries: func(limit int, offset int) ([]database.CapitalsLeaderboardEntry, error) {
				return []database.CapitalsLeaderboardEntry{}, nil
			},
			getCapitalsLeaderboardEntryID: func(limit int, offset int) (int, error) { return 0, sql.ErrNoRows },
			page:                          "0",
			status:                        http.StatusOK,
			hasMore:                       false,
		},
		{
			name: "happy path, has more is true",
			getCapitalsLeaderboardEntries: func(limit int, offset int) ([]database.CapitalsLeaderboardEntry, error) {
				return []database.CapitalsLeaderboardEntry{}, nil
			},
			getCapitalsLeaderboardEntryID: func(limit int, offset int) (int, error) { return 1, nil },
			page:                          "0",
			status:                        http.StatusOK,
			hasMore:                       true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetCapitalsLeaderboardEntries = tc.getCapitalsLeaderboardEntries
			database.GetCapitalsLeaderboardEntryID = tc.getCapitalsLeaderboardEntryID

			request, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/api/capitals/leaderboard?page=%v", tc.page), nil)
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
	savedGetCapitalsLeaderboardEntry := database.GetCapitalsLeaderboardEntry

	defer func() {
		database.GetCapitalsLeaderboardEntry = savedGetCapitalsLeaderboardEntry
	}()

	tt := []struct {
		name                        string
		getCapitalsLeaderboardEntry func(userID int) (database.CapitalsLeaderboardEntry, error)
		userID                      string
		status                      int
	}{
		{
			name:                        "invalid userId",
			getCapitalsLeaderboardEntry: database.GetCapitalsLeaderboardEntry,
			userID:                      "testing",
			status:                      http.StatusBadRequest,
		},
		{
			name: "valid userId, error on GetCapitalsLeaderboardEntry",
			getCapitalsLeaderboardEntry: func(userID int) (database.CapitalsLeaderboardEntry, error) {
				return database.CapitalsLeaderboardEntry{}, errors.New("test")
			},
			userID: "1",
			status: http.StatusInternalServerError,
		},
		{
			name: "happy path",
			getCapitalsLeaderboardEntry: func(userID int) (database.CapitalsLeaderboardEntry, error) {
				return database.CapitalsLeaderboardEntry{}, nil
			},
			userID: "1",
			status: http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetCapitalsLeaderboardEntry = tc.getCapitalsLeaderboardEntry

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

				var parsed database.CapitalsLeaderboardEntry
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestCreateEntry(t *testing.T) {
	savedGetUser := database.GetUser
	savedValidUser := auth.ValidUser
	savedInsertCapitalsLeaderboardEntry := database.InsertCapitalsLeaderboardEntry
	savedConfigValues := config.Values

	defer func() {
		database.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		database.InsertCapitalsLeaderboardEntry = savedInsertCapitalsLeaderboardEntry
		config.Values = savedConfigValues
	}()

	user := database.User{
		Username: "testing",
	}

	tt := []struct {
		name                           string
		getUser                        func(id int) (database.User, error)
		validUser                      func(uv auth.UserValidation) (int, error)
		insertCapitalsLeaderboardEntry func(entry database.CapitalsLeaderboardEntry) (int, error)
		body                           string
		status                         int
	}{
		{
			name:                           "invalid body",
			getUser:                        database.GetUser,
			validUser:                      auth.ValidUser,
			insertCapitalsLeaderboardEntry: database.InsertCapitalsLeaderboardEntry,
			body:                           "testing",
			status:                         http.StatusBadRequest,
		},
		{
			name:                           "valid body, error on GetUser",
			getUser:                        func(id int) (database.User, error) { return database.User{}, errors.New("test") },
			validUser:                      auth.ValidUser,
			insertCapitalsLeaderboardEntry: database.InsertCapitalsLeaderboardEntry,
			body:                           `{"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                         http.StatusInternalServerError,
		},
		{
			name:    "valid body, invalid user",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			insertCapitalsLeaderboardEntry: database.InsertCapitalsLeaderboardEntry,
			body:                           `{"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                         http.StatusUnauthorized,
		},
		{
			name:    "valid body, valid user, error on InsertCapitalsLeaderboardEntry",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			insertCapitalsLeaderboardEntry: func(entry database.CapitalsLeaderboardEntry) (int, error) { return 0, errors.New("test") },
			body:                           `{"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                         http.StatusInternalServerError,
		},
		{
			name:    "happy path",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			insertCapitalsLeaderboardEntry: func(entry database.CapitalsLeaderboardEntry) (int, error) { return 1, nil },
			body:                           `{"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                         http.StatusCreated,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			database.InsertCapitalsLeaderboardEntry = tc.insertCapitalsLeaderboardEntry
			config.Values = &config.Config{}

			request, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(tc.body)))
			if err != nil {
				t.Fatalf("could not create POST request: %v", err)
			}

			writer := httptest.NewRecorder()
			CreateEntry(writer, request)
			result := writer.Result()
			defer result.Body.Close()

			if result.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, result.StatusCode)
			}

			if tc.status == http.StatusCreated {
				body, err := ioutil.ReadAll(result.Body)
				if err != nil {
					t.Fatalf("could not read response: %v", err)
				}

				var parsed database.CapitalsLeaderboardEntry
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestUpdateEntry(t *testing.T) {
	savedGetUser := database.GetUser
	savedValidUser := auth.ValidUser
	savedUpdateCapitalsLeaderboardEntry := database.UpdateCapitalsLeaderboardEntry
	savedConfigValues := config.Values

	defer func() {
		database.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		database.UpdateCapitalsLeaderboardEntry = savedUpdateCapitalsLeaderboardEntry
		config.Values = savedConfigValues
	}()

	user := database.User{
		Username: "testing",
	}

	tt := []struct {
		name                           string
		getUser                        func(id int) (database.User, error)
		validUser                      func(uv auth.UserValidation) (int, error)
		updateCapitalsLeaderboardEntry func(entry database.CapitalsLeaderboardEntry) error
		id                             string
		body                           string
		status                         int
	}{
		{
			name:                           "invalid id",
			getUser:                        database.GetUser,
			validUser:                      auth.ValidUser,
			updateCapitalsLeaderboardEntry: database.UpdateCapitalsLeaderboardEntry,
			id:                             "testing",
			body:                           "",
			status:                         http.StatusBadRequest,
		},
		{
			name:                           "valid id, invalid body",
			getUser:                        func(id int) (database.User, error) { return user, nil },
			validUser:                      auth.ValidUser,
			updateCapitalsLeaderboardEntry: database.UpdateCapitalsLeaderboardEntry,
			id:                             "1",
			body:                           "testing",
			status:                         http.StatusBadRequest,
		},
		{
			name:                           "valid id, valid body, error on GetUser",
			getUser:                        func(id int) (database.User, error) { return database.User{}, errors.New("test") },
			validUser:                      auth.ValidUser,
			updateCapitalsLeaderboardEntry: database.UpdateCapitalsLeaderboardEntry,
			id:                             "1",
			body:                           `{"id": 1,"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                         http.StatusInternalServerError,
		},
		{
			name:    "valid id, valid body, invalid user",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			updateCapitalsLeaderboardEntry: database.UpdateCapitalsLeaderboardEntry,
			id:                             "1",
			body:                           `{"id": 1,"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                         http.StatusUnauthorized,
		},
		{
			name:    "valid id, valid body, valid user, error on UpdateCapitalsLeaderboardEntry",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			updateCapitalsLeaderboardEntry: func(entry database.CapitalsLeaderboardEntry) error { return errors.New("test") },
			id:                             "1",
			body:                           `{"id": 1,"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                         http.StatusInternalServerError,
		},
		{
			name:    "happy path",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			updateCapitalsLeaderboardEntry: func(entry database.CapitalsLeaderboardEntry) error { return nil },
			id:                             "1",
			body:                           `{"id": 1,"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                         http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			database.UpdateCapitalsLeaderboardEntry = tc.updateCapitalsLeaderboardEntry
			config.Values = &config.Config{}

			request, err := http.NewRequest("PUT", "", bytes.NewBuffer([]byte(tc.body)))
			if err != nil {
				t.Fatalf("could not create PUT request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"id": tc.id,
			})

			writer := httptest.NewRecorder()
			UpdateEntry(writer, request)
			result := writer.Result()
			defer result.Body.Close()

			if result.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, result.StatusCode)
			}

			if tc.status == http.StatusCreated {
				body, err := ioutil.ReadAll(result.Body)
				if err != nil {
					t.Fatalf("could not read response: %v", err)
				}

				var parsed database.CapitalsLeaderboardEntry
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestDeleteEntry(t *testing.T) {
	savedGetCapitalsLeaderboardEntry := database.GetCapitalsLeaderboardEntry
	savedGetUser := database.GetUser
	savedValidUser := auth.ValidUser
	savedDeleteCapitalsLeaderboardEntry := database.DeleteCapitalsLeaderboardEntry
	savedConfigValues := config.Values

	defer func() {
		database.GetCapitalsLeaderboardEntry = savedGetCapitalsLeaderboardEntry
		database.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		database.DeleteCapitalsLeaderboardEntry = savedDeleteCapitalsLeaderboardEntry
		config.Values = savedConfigValues
	}()

	user := database.User{
		Username: "testing",
	}

	tt := []struct {
		name                           string
		getCapitalsLeaderboardEntry    func(userID int) (database.CapitalsLeaderboardEntry, error)
		getUser                        func(id int) (database.User, error)
		validUser                      func(uv auth.UserValidation) (int, error)
		deleteCapitalsLeaderboardEntry func(entryID int) error
		id                             string
		status                         int
	}{
		{
			name:                           "invalid id",
			getCapitalsLeaderboardEntry:    database.GetCapitalsLeaderboardEntry,
			getUser:                        database.GetUser,
			validUser:                      auth.ValidUser,
			deleteCapitalsLeaderboardEntry: database.DeleteCapitalsLeaderboardEntry,
			id:                             "testing",
			status:                         http.StatusBadRequest,
		},
		{
			name: "valid id, error on GetCapitalsLeaderboardEntry",
			getCapitalsLeaderboardEntry: func(userID int) (database.CapitalsLeaderboardEntry, error) {
				return database.CapitalsLeaderboardEntry{}, errors.New("test")
			},
			getUser:                        func(id int) (database.User, error) { return user, nil },
			validUser:                      auth.ValidUser,
			deleteCapitalsLeaderboardEntry: database.DeleteCapitalsLeaderboardEntry,
			id:                             "1",
			status:                         http.StatusInternalServerError,
		},
		{
			name: "valid id, entry found, error on GetUser",
			getCapitalsLeaderboardEntry: func(userID int) (database.CapitalsLeaderboardEntry, error) {
				return database.CapitalsLeaderboardEntry{}, nil
			},
			getUser:                        func(id int) (database.User, error) { return database.User{}, errors.New("test") },
			validUser:                      auth.ValidUser,
			deleteCapitalsLeaderboardEntry: database.DeleteCapitalsLeaderboardEntry,
			id:                             "1",
			status:                         http.StatusInternalServerError,
		},
		{
			name: "valid id, entry found, invalid user",
			getCapitalsLeaderboardEntry: func(userID int) (database.CapitalsLeaderboardEntry, error) {
				return database.CapitalsLeaderboardEntry{}, nil
			},
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			deleteCapitalsLeaderboardEntry: database.DeleteCapitalsLeaderboardEntry,
			id:                             "1",
			status:                         http.StatusUnauthorized,
		},
		{
			name: "valid id, entry found, valid user, error on DeleteCapitalsLeaderboardEntry",
			getCapitalsLeaderboardEntry: func(userID int) (database.CapitalsLeaderboardEntry, error) {
				return database.CapitalsLeaderboardEntry{}, nil
			},
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			deleteCapitalsLeaderboardEntry: func(entryID int) error { return errors.New("test") },
			id:                             "1",
			status:                         http.StatusInternalServerError,
		},
		{
			name: "happy path",
			getCapitalsLeaderboardEntry: func(userID int) (database.CapitalsLeaderboardEntry, error) {
				return database.CapitalsLeaderboardEntry{}, nil
			},
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			deleteCapitalsLeaderboardEntry: func(entryID int) error { return nil },
			id:                             "1",
			status:                         http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetCapitalsLeaderboardEntry = tc.getCapitalsLeaderboardEntry
			database.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			database.DeleteCapitalsLeaderboardEntry = tc.deleteCapitalsLeaderboardEntry
			config.Values = &config.Config{}

			request, err := http.NewRequest("DELETE", "", nil)
			if err != nil {
				t.Fatalf("could not create DELETE request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"id": tc.id,
			})

			writer := httptest.NewRecorder()
			DeleteEntry(writer, request)
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

				var parsed database.CapitalsLeaderboardEntry
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}
