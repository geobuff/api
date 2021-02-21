package countries

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geobuff/api/config"
	"github.com/geobuff/api/database"
	"github.com/geobuff/api/models"
	"github.com/geobuff/auth"
	"github.com/gorilla/mux"
)

func TestGetEntries(t *testing.T) {
	savedGetLeaderboardEntries := database.GetLeaderboardEntries
	savedGetLeaderboardEntryID := database.GetLeaderboardEntryID

	defer func() {
		database.GetLeaderboardEntries = savedGetLeaderboardEntries
		database.GetLeaderboardEntryID = savedGetLeaderboardEntryID
	}()

	tt := []struct {
		name                           string
		getCountriesLeaderboardEntries func(table string, filterParams models.GetEntriesFilterParams) ([]database.LeaderboardEntryDto, error)
		getLeaderboardEntryID          func(table string, filterParams models.GetEntriesFilterParams) (int, error)
		body                           string
		status                         int
		hasMore                        bool
	}{
		{
			name:                           "error on unmarshal",
			getCountriesLeaderboardEntries: database.GetLeaderboardEntries,
			getLeaderboardEntryID:          database.GetLeaderboardEntryID,
			body:                           "",
			status:                         http.StatusBadRequest,
			hasMore:                        false,
		},
		{
			name: "error on GetLeaderboardEntries",
			getCountriesLeaderboardEntries: func(table string, filterParams models.GetEntriesFilterParams) ([]database.LeaderboardEntryDto, error) {
				return nil, errors.New("test")
			},
			getLeaderboardEntryID: database.GetLeaderboardEntryID,
			body:                  `{"page": 0, "limit": 10, "range": "", "user": ""}`,
			status:                http.StatusInternalServerError,
			hasMore:               false,
		},
		{
			name: "error on GetLeaderboardEntryID",
			getCountriesLeaderboardEntries: func(table string, filterParams models.GetEntriesFilterParams) ([]database.LeaderboardEntryDto, error) {
				return []database.LeaderboardEntryDto{}, nil
			},
			getLeaderboardEntryID: func(table string, filterParams models.GetEntriesFilterParams) (int, error) {
				return 0, errors.New("test")
			},
			body:    `{"page": 0, "limit": 10, "range": "", "user": ""}`,
			status:  http.StatusInternalServerError,
			hasMore: false,
		},
		{
			name: "happy path, has more is false",
			getCountriesLeaderboardEntries: func(table string, filterParams models.GetEntriesFilterParams) ([]database.LeaderboardEntryDto, error) {
				return []database.LeaderboardEntryDto{}, nil
			},
			getLeaderboardEntryID: func(table string, filterParams models.GetEntriesFilterParams) (int, error) { return 0, sql.ErrNoRows },
			body:                  `{"page": 0, "limit": 10, "range": "", "user": ""}`,
			status:                http.StatusOK,
			hasMore:               false,
		},
		{
			name: "happy path, has more is true",
			getCountriesLeaderboardEntries: func(table string, filterParams models.GetEntriesFilterParams) ([]database.LeaderboardEntryDto, error) {
				return []database.LeaderboardEntryDto{}, nil
			},
			getLeaderboardEntryID: func(table string, filterParams models.GetEntriesFilterParams) (int, error) { return 1, nil },
			body:                  `{"page": 0, "limit": 10, "range": "", "user": ""}`,
			status:                http.StatusOK,
			hasMore:               true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetLeaderboardEntries = tc.getCountriesLeaderboardEntries
			database.GetLeaderboardEntryID = tc.getLeaderboardEntryID

			request, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(tc.body)))
			if err != nil {
				t.Fatalf("could not create POST request: %v", err)
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
	savedGetLeaderboardEntry := database.GetLeaderboardEntry

	defer func() {
		database.GetLeaderboardEntry = savedGetLeaderboardEntry
	}()

	tt := []struct {
		name                string
		getLeaderboardEntry func(table string, userID int) (database.LeaderboardEntry, error)
		userID              string
		status              int
	}{
		{
			name:                "invalid userId",
			getLeaderboardEntry: database.GetLeaderboardEntry,
			userID:              "testing",
			status:              http.StatusBadRequest,
		},
		{
			name: "valid userId, error on GetLeaderboardEntry",
			getLeaderboardEntry: func(table string, userID int) (database.LeaderboardEntry, error) {
				return database.LeaderboardEntry{}, errors.New("test")
			},
			userID: "1",
			status: http.StatusInternalServerError,
		},
		{
			name: "happy path",
			getLeaderboardEntry: func(table string, userID int) (database.LeaderboardEntry, error) {
				return database.LeaderboardEntry{}, nil
			},
			userID: "1",
			status: http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetLeaderboardEntry = tc.getLeaderboardEntry

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

				var parsed database.LeaderboardEntry
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
	savedInsertLeaderboardEntry := database.InsertLeaderboardEntry
	savedConfigValues := config.Values

	defer func() {
		database.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		database.InsertLeaderboardEntry = savedInsertLeaderboardEntry
		config.Values = savedConfigValues
	}()

	user := database.User{
		Username: "testing",
	}

	tt := []struct {
		name                   string
		getUser                func(id int) (database.User, error)
		validUser              func(uv auth.UserValidation) (int, error)
		insertLeaderboardEntry func(table string, entry database.LeaderboardEntry) (int, error)
		body                   string
		status                 int
	}{
		{
			name:                   "invalid body",
			getUser:                database.GetUser,
			validUser:              auth.ValidUser,
			insertLeaderboardEntry: database.InsertLeaderboardEntry,
			body:                   "testing",
			status:                 http.StatusBadRequest,
		},
		{
			name:                   "valid body, error on GetUser",
			getUser:                func(id int) (database.User, error) { return database.User{}, errors.New("test") },
			validUser:              auth.ValidUser,
			insertLeaderboardEntry: database.InsertLeaderboardEntry,
			body:                   `{"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                 http.StatusInternalServerError,
		},
		{
			name:    "valid body, invalid user",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			insertLeaderboardEntry: database.InsertLeaderboardEntry,
			body:                   `{"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                 http.StatusUnauthorized,
		},
		{
			name:    "valid body, valid user, error on InsertLeaderboardEntry",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			insertLeaderboardEntry: func(table string, entry database.LeaderboardEntry) (int, error) { return 0, errors.New("test") },
			body:                   `{"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                 http.StatusInternalServerError,
		},
		{
			name:    "happy path",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			insertLeaderboardEntry: func(table string, entry database.LeaderboardEntry) (int, error) { return 1, nil },
			body:                   `{"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                 http.StatusCreated,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			database.InsertLeaderboardEntry = tc.insertLeaderboardEntry
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

				var parsed database.LeaderboardEntry
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
	savedUpdateLeaderboardEntry := database.UpdateLeaderboardEntry
	savedConfigValues := config.Values

	defer func() {
		database.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		database.UpdateLeaderboardEntry = savedUpdateLeaderboardEntry
		config.Values = savedConfigValues
	}()

	user := database.User{
		Username: "testing",
	}

	tt := []struct {
		name                   string
		getUser                func(id int) (database.User, error)
		validUser              func(uv auth.UserValidation) (int, error)
		updateLeaderboardEntry func(table string, entry database.LeaderboardEntry) error
		id                     string
		body                   string
		status                 int
	}{
		{
			name:                   "invalid id",
			getUser:                database.GetUser,
			validUser:              auth.ValidUser,
			updateLeaderboardEntry: database.UpdateLeaderboardEntry,
			id:                     "testing",
			body:                   "",
			status:                 http.StatusBadRequest,
		},
		{
			name:                   "valid id, invalid body",
			getUser:                func(id int) (database.User, error) { return user, nil },
			validUser:              auth.ValidUser,
			updateLeaderboardEntry: database.UpdateLeaderboardEntry,
			id:                     "1",
			body:                   "testing",
			status:                 http.StatusBadRequest,
		},
		{
			name:                   "valid id, valid body, error on GetUser",
			getUser:                func(id int) (database.User, error) { return database.User{}, errors.New("test") },
			validUser:              auth.ValidUser,
			updateLeaderboardEntry: database.UpdateLeaderboardEntry,
			id:                     "1",
			body:                   `{"id": 1,"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                 http.StatusInternalServerError,
		},
		{
			name:    "valid id, valid body, invalid user",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			updateLeaderboardEntry: database.UpdateLeaderboardEntry,
			id:                     "1",
			body:                   `{"id": 1,"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                 http.StatusUnauthorized,
		},
		{
			name:    "valid id, valid body, valid user, error on UpdateLeaderboardEntry",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			updateLeaderboardEntry: func(table string, entry database.LeaderboardEntry) error { return errors.New("test") },
			id:                     "1",
			body:                   `{"id": 1,"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                 http.StatusInternalServerError,
		},
		{
			name:    "happy path",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			updateLeaderboardEntry: func(table string, entry database.LeaderboardEntry) error { return nil },
			id:                     "1",
			body:                   `{"id": 1,"userId": 1, "country": "New Zealand", "countries": 100, "time": 200}`,
			status:                 http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			database.UpdateLeaderboardEntry = tc.updateLeaderboardEntry
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

				var parsed database.LeaderboardEntry
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestDeleteEntry(t *testing.T) {
	savedGetLeaderboardEntry := database.GetLeaderboardEntry
	savedGetUser := database.GetUser
	savedValidUser := auth.ValidUser
	savedDeleteLeaderboardEntry := database.DeleteLeaderboardEntry
	savedConfigValues := config.Values

	defer func() {
		database.GetLeaderboardEntry = savedGetLeaderboardEntry
		database.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		database.DeleteLeaderboardEntry = savedDeleteLeaderboardEntry
		config.Values = savedConfigValues
	}()

	user := database.User{
		Username: "testing",
	}

	tt := []struct {
		name                   string
		getLeaderboardEntry    func(table string, userID int) (database.LeaderboardEntry, error)
		getUser                func(id int) (database.User, error)
		validUser              func(uv auth.UserValidation) (int, error)
		deleteLeaderboardEntry func(table string, entryID int) error
		id                     string
		status                 int
	}{
		{
			name:                   "invalid id",
			getLeaderboardEntry:    database.GetLeaderboardEntry,
			getUser:                database.GetUser,
			validUser:              auth.ValidUser,
			deleteLeaderboardEntry: database.DeleteLeaderboardEntry,
			id:                     "testing",
			status:                 http.StatusBadRequest,
		},
		{
			name: "valid id, error on GetLeaderboardEntry",
			getLeaderboardEntry: func(table string, userID int) (database.LeaderboardEntry, error) {
				return database.LeaderboardEntry{}, errors.New("test")
			},
			getUser:                func(id int) (database.User, error) { return user, nil },
			validUser:              auth.ValidUser,
			deleteLeaderboardEntry: database.DeleteLeaderboardEntry,
			id:                     "1",
			status:                 http.StatusInternalServerError,
		},
		{
			name: "valid id, entry found, error on GetUser",
			getLeaderboardEntry: func(table string, userID int) (database.LeaderboardEntry, error) {
				return database.LeaderboardEntry{}, nil
			},
			getUser:                func(id int) (database.User, error) { return database.User{}, errors.New("test") },
			validUser:              auth.ValidUser,
			deleteLeaderboardEntry: database.DeleteLeaderboardEntry,
			id:                     "1",
			status:                 http.StatusInternalServerError,
		},
		{
			name: "valid id, entry found, invalid user",
			getLeaderboardEntry: func(table string, userID int) (database.LeaderboardEntry, error) {
				return database.LeaderboardEntry{}, nil
			},
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			deleteLeaderboardEntry: database.DeleteLeaderboardEntry,
			id:                     "1",
			status:                 http.StatusUnauthorized,
		},
		{
			name: "valid id, entry found, valid user, error on DeleteLeaderboardEntry",
			getLeaderboardEntry: func(table string, userID int) (database.LeaderboardEntry, error) {
				return database.LeaderboardEntry{}, nil
			},
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			deleteLeaderboardEntry: func(table string, entryID int) error { return errors.New("test") },
			id:                     "1",
			status:                 http.StatusInternalServerError,
		},
		{
			name: "happy path",
			getLeaderboardEntry: func(table string, userID int) (database.LeaderboardEntry, error) {
				return database.LeaderboardEntry{}, nil
			},
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			deleteLeaderboardEntry: func(table string, entryID int) error { return nil },
			id:                     "1",
			status:                 http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetLeaderboardEntry = tc.getLeaderboardEntry
			database.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			database.DeleteLeaderboardEntry = tc.deleteLeaderboardEntry
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

				var parsed database.LeaderboardEntry
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}
