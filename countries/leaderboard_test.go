package countries

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

func TestCreateEntry(t *testing.T) {
	savedGetUser := database.GetUser
	savedValidUser := auth.ValidUser
	savedInsertWorldLeaderboardEntry := database.InsertWorldLeaderboardEntry
	savedConfigValues := config.Values

	defer func() {
		database.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		database.InsertWorldLeaderboardEntry = savedInsertWorldLeaderboardEntry
		config.Values = savedConfigValues
	}()

	user := database.User{
		Username: "testing",
	}

	tt := []struct {
		name                        string
		getUser                     func(id int) (database.User, error)
		validUser                   func(uv auth.UserValidation) (int, error)
		insertWorldLeaderboardEntry func(entry database.WorldLeaderboardEntry) (int, error)
		body                        string
		status                      int
	}{
		{
			name:                        "invalid body",
			getUser:                     database.GetUser,
			validUser:                   auth.ValidUser,
			insertWorldLeaderboardEntry: database.InsertWorldLeaderboardEntry,
			body:                        "testing",
			status:                      http.StatusBadRequest,
		},
		{
			name:                        "valid body, error on GetUser",
			getUser:                     func(id int) (database.User, error) { return database.User{}, errors.New("test") },
			validUser:                   auth.ValidUser,
			insertWorldLeaderboardEntry: database.InsertWorldLeaderboardEntry,
			body:                        `{"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                      http.StatusInternalServerError,
		},
		{
			name:    "valid body, invalid user",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			insertWorldLeaderboardEntry: database.InsertWorldLeaderboardEntry,
			body:                        `{"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                      http.StatusUnauthorized,
		},
		{
			name:    "valid body, valid user, error on InsertWorldLeaderboardEntry",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			insertWorldLeaderboardEntry: func(entry database.WorldLeaderboardEntry) (int, error) { return 0, errors.New("test") },
			body:                        `{"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                      http.StatusInternalServerError,
		},
		{
			name:    "happy path",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			insertWorldLeaderboardEntry: func(entry database.WorldLeaderboardEntry) (int, error) { return 1, nil },
			body:                        `{"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                      http.StatusCreated,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			database.InsertWorldLeaderboardEntry = tc.insertWorldLeaderboardEntry
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

				var parsed database.WorldLeaderboardEntry
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
	savedUpdateWorldLeaderboardEntry := database.UpdateWorldLeaderboardEntry
	savedConfigValues := config.Values

	defer func() {
		database.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		database.UpdateWorldLeaderboardEntry = savedUpdateWorldLeaderboardEntry
		config.Values = savedConfigValues
	}()

	user := database.User{
		Username: "testing",
	}

	tt := []struct {
		name                        string
		getUser                     func(id int) (database.User, error)
		validUser                   func(uv auth.UserValidation) (int, error)
		updateWorldLeaderboardEntry func(entry database.WorldLeaderboardEntry) error
		id                          string
		body                        string
		status                      int
	}{
		{
			name:                        "invalid id",
			getUser:                     database.GetUser,
			validUser:                   auth.ValidUser,
			updateWorldLeaderboardEntry: database.UpdateWorldLeaderboardEntry,
			id:                          "testing",
			body:                        "",
			status:                      http.StatusBadRequest,
		},
		{
			name:                        "valid id, invalid body",
			getUser:                     func(id int) (database.User, error) { return user, nil },
			validUser:                   auth.ValidUser,
			updateWorldLeaderboardEntry: database.UpdateWorldLeaderboardEntry,
			id:                          "1",
			body:                        "testing",
			status:                      http.StatusBadRequest,
		},
		{
			name:                        "valid id, valid body, error on GetUser",
			getUser:                     func(id int) (database.User, error) { return database.User{}, errors.New("test") },
			validUser:                   auth.ValidUser,
			updateWorldLeaderboardEntry: database.UpdateWorldLeaderboardEntry,
			id:                          "1",
			body:                        `{"id": 1,"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                      http.StatusInternalServerError,
		},
		{
			name:    "valid id, valid body, invalid user",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			updateWorldLeaderboardEntry: database.UpdateWorldLeaderboardEntry,
			id:                          "1",
			body:                        `{"id": 1,"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                      http.StatusUnauthorized,
		},
		{
			name:    "valid id, valid body, valid user, error on UpdateWorldLeaderboardEntry",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			updateWorldLeaderboardEntry: func(entry database.WorldLeaderboardEntry) error { return errors.New("test") },
			id:                          "1",
			body:                        `{"id": 1,"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                      http.StatusInternalServerError,
		},
		{
			name:    "happy path",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			updateWorldLeaderboardEntry: func(entry database.WorldLeaderboardEntry) error { return nil },
			id:                          "1",
			body:                        `{"id": 1,"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                      http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			database.UpdateWorldLeaderboardEntry = tc.updateWorldLeaderboardEntry
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

				var parsed database.WorldLeaderboardEntry
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestDeleteEntry(t *testing.T) {
	savedGetWorldLeaderboardEntry := database.GetWorldLeaderboardEntry
	savedGetUser := database.GetUser
	savedValidUser := auth.ValidUser
	savedDeleteWorldLeaderboardEntry := database.DeleteWorldLeaderboardEntry
	savedConfigValues := config.Values

	defer func() {
		database.GetWorldLeaderboardEntry = savedGetWorldLeaderboardEntry
		database.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		database.DeleteWorldLeaderboardEntry = savedDeleteWorldLeaderboardEntry
		config.Values = savedConfigValues
	}()

	user := database.User{
		Username: "testing",
	}

	tt := []struct {
		name                        string
		getWorldLeaderboardEntry    func(userID int) (database.WorldLeaderboardEntry, error)
		getUser                     func(id int) (database.User, error)
		validUser                   func(uv auth.UserValidation) (int, error)
		deleteWorldLeaderboardEntry func(entryID int) error
		id                          string
		status                      int
	}{
		{
			name:                        "invalid id",
			getWorldLeaderboardEntry:    database.GetWorldLeaderboardEntry,
			getUser:                     database.GetUser,
			validUser:                   auth.ValidUser,
			deleteWorldLeaderboardEntry: database.DeleteWorldLeaderboardEntry,
			id:                          "testing",
			status:                      http.StatusBadRequest,
		},
		{
			name: "valid id, error on GetWorldLeaderboardEntry",
			getWorldLeaderboardEntry: func(userID int) (database.WorldLeaderboardEntry, error) {
				return database.WorldLeaderboardEntry{}, errors.New("test")
			},
			getUser:                     func(id int) (database.User, error) { return user, nil },
			validUser:                   auth.ValidUser,
			deleteWorldLeaderboardEntry: database.DeleteWorldLeaderboardEntry,
			id:                          "1",
			status:                      http.StatusInternalServerError,
		},
		{
			name: "valid id, entry found, error on GetUser",
			getWorldLeaderboardEntry: func(userID int) (database.WorldLeaderboardEntry, error) {
				return database.WorldLeaderboardEntry{}, nil
			},
			getUser:                     func(id int) (database.User, error) { return database.User{}, errors.New("test") },
			validUser:                   auth.ValidUser,
			deleteWorldLeaderboardEntry: database.DeleteWorldLeaderboardEntry,
			id:                          "1",
			status:                      http.StatusInternalServerError,
		},
		{
			name: "valid id, entry found, invalid user",
			getWorldLeaderboardEntry: func(userID int) (database.WorldLeaderboardEntry, error) {
				return database.WorldLeaderboardEntry{}, nil
			},
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			deleteWorldLeaderboardEntry: database.DeleteWorldLeaderboardEntry,
			id:                          "1",
			status:                      http.StatusUnauthorized,
		},
		{
			name: "valid id, entry found, valid user, error on DeleteWorldLeaderboardEntry",
			getWorldLeaderboardEntry: func(userID int) (database.WorldLeaderboardEntry, error) {
				return database.WorldLeaderboardEntry{}, nil
			},
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			deleteWorldLeaderboardEntry: func(entryID int) error { return errors.New("test") },
			id:                          "1",
			status:                      http.StatusInternalServerError,
		},
		{
			name: "happy path",
			getWorldLeaderboardEntry: func(userID int) (database.WorldLeaderboardEntry, error) {
				return database.WorldLeaderboardEntry{}, nil
			},
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			deleteWorldLeaderboardEntry: func(entryID int) error { return nil },
			id:                          "1",
			status:                      http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetWorldLeaderboardEntry = tc.getWorldLeaderboardEntry
			database.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			database.DeleteWorldLeaderboardEntry = tc.deleteWorldLeaderboardEntry
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

				var parsed database.WorldLeaderboardEntry
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}
