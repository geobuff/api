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
	savedGetCountriesLeaderboardEntries := database.GetCountriesLeaderboardEntries
	savedGetCountriesLeaderboardEntryID := database.GetCountriesLeaderboardEntryID

	defer func() {
		database.GetCountriesLeaderboardEntries = savedGetCountriesLeaderboardEntries
		database.GetCountriesLeaderboardEntryID = savedGetCountriesLeaderboardEntryID
	}()

	tt := []struct {
		name                           string
		getCountriesLeaderboardEntries func(limit int, offset int) ([]database.CountriesLeaderboardEntry, error)
		getCountriesLeaderboardEntryID func(limit int, offset int) (int, error)
		page                           string
		status                         int
		hasMore                        bool
	}{
		{
			name:                           "invalid page parameter",
			getCountriesLeaderboardEntries: database.GetCountriesLeaderboardEntries,
			getCountriesLeaderboardEntryID: database.GetCountriesLeaderboardEntryID,
			page:                           "testing",
			status:                         http.StatusBadRequest,
			hasMore:                        false,
		},
		{
			name: "valid page, error on GetCountriesLeaderboardEntries",
			getCountriesLeaderboardEntries: func(limit int, offset int) ([]database.CountriesLeaderboardEntry, error) {
				return nil, errors.New("test")
			},
			getCountriesLeaderboardEntryID: database.GetCountriesLeaderboardEntryID,
			page:                           "0",
			status:                         http.StatusInternalServerError,
			hasMore:                        false,
		},
		{
			name: "valid page, error on GetCountriesLeaderboardEntryID",
			getCountriesLeaderboardEntries: func(limit int, offset int) ([]database.CountriesLeaderboardEntry, error) {
				return []database.CountriesLeaderboardEntry{}, nil
			},
			getCountriesLeaderboardEntryID: func(limit int, offset int) (int, error) { return 0, errors.New("test") },
			page:                           "0",
			status:                         http.StatusInternalServerError,
			hasMore:                        false,
		},
		{
			name: "happy path, has more is false",
			getCountriesLeaderboardEntries: func(limit int, offset int) ([]database.CountriesLeaderboardEntry, error) {
				return []database.CountriesLeaderboardEntry{}, nil
			},
			getCountriesLeaderboardEntryID: func(limit int, offset int) (int, error) { return 0, sql.ErrNoRows },
			page:                           "0",
			status:                         http.StatusOK,
			hasMore:                        false,
		},
		{
			name: "happy path, has more is true",
			getCountriesLeaderboardEntries: func(limit int, offset int) ([]database.CountriesLeaderboardEntry, error) {
				return []database.CountriesLeaderboardEntry{}, nil
			},
			getCountriesLeaderboardEntryID: func(limit int, offset int) (int, error) { return 1, nil },
			page:                           "0",
			status:                         http.StatusOK,
			hasMore:                        true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetCountriesLeaderboardEntries = tc.getCountriesLeaderboardEntries
			database.GetCountriesLeaderboardEntryID = tc.getCountriesLeaderboardEntryID

			request, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/api/Countries/leaderboard?page=%v", tc.page), nil)
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
	savedGetCountriesLeaderboardEntry := database.GetCountriesLeaderboardEntry

	defer func() {
		database.GetCountriesLeaderboardEntry = savedGetCountriesLeaderboardEntry
	}()

	tt := []struct {
		name                         string
		getCountriesLeaderboardEntry func(userID int) (database.CountriesLeaderboardEntry, error)
		userID                       string
		status                       int
	}{
		{
			name:                         "invalid userId",
			getCountriesLeaderboardEntry: database.GetCountriesLeaderboardEntry,
			userID:                       "testing",
			status:                       http.StatusBadRequest,
		},
		{
			name: "valid userId, error on GetCountriesLeaderboardEntry",
			getCountriesLeaderboardEntry: func(userID int) (database.CountriesLeaderboardEntry, error) {
				return database.CountriesLeaderboardEntry{}, errors.New("test")
			},
			userID: "1",
			status: http.StatusInternalServerError,
		},
		{
			name: "happy path",
			getCountriesLeaderboardEntry: func(userID int) (database.CountriesLeaderboardEntry, error) {
				return database.CountriesLeaderboardEntry{}, nil
			},
			userID: "1",
			status: http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetCountriesLeaderboardEntry = tc.getCountriesLeaderboardEntry

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

				var parsed database.CountriesLeaderboardEntry
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
	savedInsertCountriesLeaderboardEntry := database.InsertCountriesLeaderboardEntry
	savedConfigValues := config.Values

	defer func() {
		database.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		database.InsertCountriesLeaderboardEntry = savedInsertCountriesLeaderboardEntry
		config.Values = savedConfigValues
	}()

	user := database.User{
		Username: "testing",
	}

	tt := []struct {
		name                            string
		getUser                         func(id int) (database.User, error)
		validUser                       func(uv auth.UserValidation) (int, error)
		insertCountriesLeaderboardEntry func(entry database.CountriesLeaderboardEntry) (int, error)
		body                            string
		status                          int
	}{
		{
			name:                            "invalid body",
			getUser:                         database.GetUser,
			validUser:                       auth.ValidUser,
			insertCountriesLeaderboardEntry: database.InsertCountriesLeaderboardEntry,
			body:                            "testing",
			status:                          http.StatusBadRequest,
		},
		{
			name:                            "valid body, error on GetUser",
			getUser:                         func(id int) (database.User, error) { return database.User{}, errors.New("test") },
			validUser:                       auth.ValidUser,
			insertCountriesLeaderboardEntry: database.InsertCountriesLeaderboardEntry,
			body:                            `{"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                          http.StatusInternalServerError,
		},
		{
			name:    "valid body, invalid user",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			insertCountriesLeaderboardEntry: database.InsertCountriesLeaderboardEntry,
			body:                            `{"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                          http.StatusUnauthorized,
		},
		{
			name:    "valid body, valid user, error on InsertCountriesLeaderboardEntry",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			insertCountriesLeaderboardEntry: func(entry database.CountriesLeaderboardEntry) (int, error) { return 0, errors.New("test") },
			body:                            `{"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                          http.StatusInternalServerError,
		},
		{
			name:    "happy path",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			insertCountriesLeaderboardEntry: func(entry database.CountriesLeaderboardEntry) (int, error) { return 1, nil },
			body:                            `{"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                          http.StatusCreated,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			database.InsertCountriesLeaderboardEntry = tc.insertCountriesLeaderboardEntry
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

				var parsed database.CountriesLeaderboardEntry
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
	savedUpdateCountriesLeaderboardEntry := database.UpdateCountriesLeaderboardEntry
	savedConfigValues := config.Values

	defer func() {
		database.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		database.UpdateCountriesLeaderboardEntry = savedUpdateCountriesLeaderboardEntry
		config.Values = savedConfigValues
	}()

	user := database.User{
		Username: "testing",
	}

	tt := []struct {
		name                            string
		getUser                         func(id int) (database.User, error)
		validUser                       func(uv auth.UserValidation) (int, error)
		updateCountriesLeaderboardEntry func(entry database.CountriesLeaderboardEntry) error
		id                              string
		body                            string
		status                          int
	}{
		{
			name:                            "invalid id",
			getUser:                         database.GetUser,
			validUser:                       auth.ValidUser,
			updateCountriesLeaderboardEntry: database.UpdateCountriesLeaderboardEntry,
			id:                              "testing",
			body:                            "",
			status:                          http.StatusBadRequest,
		},
		{
			name:                            "valid id, invalid body",
			getUser:                         func(id int) (database.User, error) { return user, nil },
			validUser:                       auth.ValidUser,
			updateCountriesLeaderboardEntry: database.UpdateCountriesLeaderboardEntry,
			id:                              "1",
			body:                            "testing",
			status:                          http.StatusBadRequest,
		},
		{
			name:                            "valid id, valid body, error on GetUser",
			getUser:                         func(id int) (database.User, error) { return database.User{}, errors.New("test") },
			validUser:                       auth.ValidUser,
			updateCountriesLeaderboardEntry: database.UpdateCountriesLeaderboardEntry,
			id:                              "1",
			body:                            `{"id": 1,"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                          http.StatusInternalServerError,
		},
		{
			name:    "valid id, valid body, invalid user",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			updateCountriesLeaderboardEntry: database.UpdateCountriesLeaderboardEntry,
			id:                              "1",
			body:                            `{"id": 1,"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                          http.StatusUnauthorized,
		},
		{
			name:    "valid id, valid body, valid user, error on UpdateCountriesLeaderboardEntry",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			updateCountriesLeaderboardEntry: func(entry database.CountriesLeaderboardEntry) error { return errors.New("test") },
			id:                              "1",
			body:                            `{"id": 1,"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                          http.StatusInternalServerError,
		},
		{
			name:    "happy path",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			updateCountriesLeaderboardEntry: func(entry database.CountriesLeaderboardEntry) error { return nil },
			id:                              "1",
			body:                            `{"id": 1,"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                          http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			database.UpdateCountriesLeaderboardEntry = tc.updateCountriesLeaderboardEntry
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

				var parsed database.CountriesLeaderboardEntry
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestDeleteEntry(t *testing.T) {
	savedGetCountriesLeaderboardEntry := database.GetCountriesLeaderboardEntry
	savedGetUser := database.GetUser
	savedValidUser := auth.ValidUser
	savedDeleteCountriesLeaderboardEntry := database.DeleteCountriesLeaderboardEntry
	savedConfigValues := config.Values

	defer func() {
		database.GetCountriesLeaderboardEntry = savedGetCountriesLeaderboardEntry
		database.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		database.DeleteCountriesLeaderboardEntry = savedDeleteCountriesLeaderboardEntry
		config.Values = savedConfigValues
	}()

	user := database.User{
		Username: "testing",
	}

	tt := []struct {
		name                            string
		getCountriesLeaderboardEntry    func(userID int) (database.CountriesLeaderboardEntry, error)
		getUser                         func(id int) (database.User, error)
		validUser                       func(uv auth.UserValidation) (int, error)
		deleteCountriesLeaderboardEntry func(entryID int) error
		id                              string
		status                          int
	}{
		{
			name:                            "invalid id",
			getCountriesLeaderboardEntry:    database.GetCountriesLeaderboardEntry,
			getUser:                         database.GetUser,
			validUser:                       auth.ValidUser,
			deleteCountriesLeaderboardEntry: database.DeleteCountriesLeaderboardEntry,
			id:                              "testing",
			status:                          http.StatusBadRequest,
		},
		{
			name: "valid id, error on GetCountriesLeaderboardEntry",
			getCountriesLeaderboardEntry: func(userID int) (database.CountriesLeaderboardEntry, error) {
				return database.CountriesLeaderboardEntry{}, errors.New("test")
			},
			getUser:                         func(id int) (database.User, error) { return user, nil },
			validUser:                       auth.ValidUser,
			deleteCountriesLeaderboardEntry: database.DeleteCountriesLeaderboardEntry,
			id:                              "1",
			status:                          http.StatusInternalServerError,
		},
		{
			name: "valid id, entry found, error on GetUser",
			getCountriesLeaderboardEntry: func(userID int) (database.CountriesLeaderboardEntry, error) {
				return database.CountriesLeaderboardEntry{}, nil
			},
			getUser:                         func(id int) (database.User, error) { return database.User{}, errors.New("test") },
			validUser:                       auth.ValidUser,
			deleteCountriesLeaderboardEntry: database.DeleteCountriesLeaderboardEntry,
			id:                              "1",
			status:                          http.StatusInternalServerError,
		},
		{
			name: "valid id, entry found, invalid user",
			getCountriesLeaderboardEntry: func(userID int) (database.CountriesLeaderboardEntry, error) {
				return database.CountriesLeaderboardEntry{}, nil
			},
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			deleteCountriesLeaderboardEntry: database.DeleteCountriesLeaderboardEntry,
			id:                              "1",
			status:                          http.StatusUnauthorized,
		},
		{
			name: "valid id, entry found, valid user, error on DeleteCountriesLeaderboardEntry",
			getCountriesLeaderboardEntry: func(userID int) (database.CountriesLeaderboardEntry, error) {
				return database.CountriesLeaderboardEntry{}, nil
			},
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			deleteCountriesLeaderboardEntry: func(entryID int) error { return errors.New("test") },
			id:                              "1",
			status:                          http.StatusInternalServerError,
		},
		{
			name: "happy path",
			getCountriesLeaderboardEntry: func(userID int) (database.CountriesLeaderboardEntry, error) {
				return database.CountriesLeaderboardEntry{}, nil
			},
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			deleteCountriesLeaderboardEntry: func(entryID int) error { return nil },
			id:                              "1",
			status:                          http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetCountriesLeaderboardEntry = tc.getCountriesLeaderboardEntry
			database.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			database.DeleteCountriesLeaderboardEntry = tc.deleteCountriesLeaderboardEntry
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

				var parsed database.CountriesLeaderboardEntry
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}
