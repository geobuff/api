package leaderboard

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geobuff/api/auth"
	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

func TestGetEntries(t *testing.T) {
	savedGetLeaderboardEntries := repo.GetLeaderboardEntries
	savedGetLeaderboardEntryID := repo.GetLeaderboardEntryID

	defer func() {
		repo.GetLeaderboardEntries = savedGetLeaderboardEntries
		repo.GetLeaderboardEntryID = savedGetLeaderboardEntryID
	}()

	tt := []struct {
		name                  string
		getLeaderboardEntries func(quizID int, filterParams repo.GetLeaderboardEntriesFilterParams) ([]repo.LeaderboardEntryDto, error)
		getLeaderboardEntryID func(quizID int, filterParams repo.GetLeaderboardEntriesFilterParams) (int, error)
		quizID                string
		body                  string
		status                int
		hasMore               bool
	}{
		{
			name:                  "invalid quizId value",
			getLeaderboardEntries: repo.GetLeaderboardEntries,
			getLeaderboardEntryID: repo.GetLeaderboardEntryID,
			quizID:                "test",
			body:                  "",
			status:                http.StatusBadRequest,
			hasMore:               false,
		},
		{
			name:                  "error on unmarshal",
			getLeaderboardEntries: repo.GetLeaderboardEntries,
			getLeaderboardEntryID: repo.GetLeaderboardEntryID,
			quizID:                "1",
			body:                  "",
			status:                http.StatusBadRequest,
			hasMore:               false,
		},
		{
			name: "error on GetLeaderboardEntries",
			getLeaderboardEntries: func(quizID int, filterParams repo.GetLeaderboardEntriesFilterParams) ([]repo.LeaderboardEntryDto, error) {
				return nil, errors.New("test")
			},
			getLeaderboardEntryID: repo.GetLeaderboardEntryID,
			quizID:                "1",
			body:                  `{"page": 0, "limit": 10, "range": "", "user": ""}`,
			status:                http.StatusInternalServerError,
			hasMore:               false,
		},
		{
			name: "error on GetLeaderboardEntryID",
			getLeaderboardEntries: func(quizID int, filterParams repo.GetLeaderboardEntriesFilterParams) ([]repo.LeaderboardEntryDto, error) {
				return []repo.LeaderboardEntryDto{}, nil
			},
			getLeaderboardEntryID: func(quizID int, filterParams repo.GetLeaderboardEntriesFilterParams) (int, error) {
				return 0, errors.New("test")
			},
			quizID:  "1",
			body:    `{"page": 0, "limit": 10, "range": "", "user": ""}`,
			status:  http.StatusInternalServerError,
			hasMore: false,
		},
		{
			name: "happy path, has more is false",
			getLeaderboardEntries: func(quizID int, filterParams repo.GetLeaderboardEntriesFilterParams) ([]repo.LeaderboardEntryDto, error) {
				return []repo.LeaderboardEntryDto{}, nil
			},
			getLeaderboardEntryID: func(quizID int, filterParams repo.GetLeaderboardEntriesFilterParams) (int, error) {
				return 0, sql.ErrNoRows
			},
			quizID:  "1",
			body:    `{"page": 0, "limit": 10, "range": "", "user": ""}`,
			status:  http.StatusOK,
			hasMore: false,
		},
		{
			name: "happy path, has more is true",
			getLeaderboardEntries: func(quizID int, filterParams repo.GetLeaderboardEntriesFilterParams) ([]repo.LeaderboardEntryDto, error) {
				return []repo.LeaderboardEntryDto{}, nil
			},
			getLeaderboardEntryID: func(quizID int, filterParams repo.GetLeaderboardEntriesFilterParams) (int, error) { return 1, nil },
			quizID:                "1",
			body:                  `{"page": 0, "limit": 10, "range": "", "user": ""}`,
			status:                http.StatusOK,
			hasMore:               true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetLeaderboardEntries = tc.getLeaderboardEntries
			repo.GetLeaderboardEntryID = tc.getLeaderboardEntryID

			request, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(tc.body)))
			if err != nil {
				t.Fatalf("could not create POST request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"quizId": tc.quizID,
			})

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

func TestGetUserEntries(t *testing.T) {
	savedGetUserLeaderboardEntries := repo.GetUserLeaderboardEntries

	defer func() {
		repo.GetUserLeaderboardEntries = savedGetUserLeaderboardEntries
	}()

	tt := []struct {
		name                      string
		getUserLeaderboardEntries func(userID int) ([]repo.UserLeaderboardEntryDto, error)
		userID                    string
		status                    int
	}{
		{
			name:                      "invalid userId",
			getUserLeaderboardEntries: repo.GetUserLeaderboardEntries,
			userID:                    "testing",
			status:                    http.StatusBadRequest,
		},
		{
			name: "valid userId, no rows found",
			getUserLeaderboardEntries: func(userID int) ([]repo.UserLeaderboardEntryDto, error) {
				return nil, sql.ErrNoRows
			},
			userID: "1",
			status: http.StatusNoContent,
		},
		{
			name: "valid userId, unknown error on GetUserLeaderboardEntries",
			getUserLeaderboardEntries: func(userID int) ([]repo.UserLeaderboardEntryDto, error) {
				return nil, errors.New("test")
			},
			userID: "1",
			status: http.StatusInternalServerError,
		},
		{
			name: "happy path",
			getUserLeaderboardEntries: func(userID int) ([]repo.UserLeaderboardEntryDto, error) {
				return []repo.UserLeaderboardEntryDto{}, nil
			},
			userID: "1",
			status: http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetUserLeaderboardEntries = tc.getUserLeaderboardEntries

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"userId": tc.userID,
			})

			writer := httptest.NewRecorder()
			GetUserEntries(writer, request)
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

				var parsed []repo.LeaderboardEntryDto
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestGetEntry(t *testing.T) {
	savedGetLeaderboardEntry := repo.GetLeaderboardEntry

	defer func() {
		repo.GetLeaderboardEntry = savedGetLeaderboardEntry
	}()

	tt := []struct {
		name                string
		getLeaderboardEntry func(quizID int, userID int) (repo.LeaderboardEntryDto, error)
		quizID              string
		userID              string
		status              int
	}{
		{
			name:                "invalid quizId",
			getLeaderboardEntry: repo.GetLeaderboardEntry,
			quizID:              "testing",
			userID:              "",
			status:              http.StatusBadRequest,
		},
		{
			name:                "invalid userId",
			getLeaderboardEntry: repo.GetLeaderboardEntry,
			quizID:              "1",
			userID:              "testing",
			status:              http.StatusBadRequest,
		},
		{
			name: "valid userId, no rows found",
			getLeaderboardEntry: func(quizID int, userID int) (repo.LeaderboardEntryDto, error) {
				return repo.LeaderboardEntryDto{}, sql.ErrNoRows
			},
			quizID: "1",
			userID: "1",
			status: http.StatusNoContent,
		},
		{
			name: "valid userId, unknown error on GetLeaderboardEntry",
			getLeaderboardEntry: func(quizID int, userID int) (repo.LeaderboardEntryDto, error) {
				return repo.LeaderboardEntryDto{}, errors.New("test")
			},
			quizID: "1",
			userID: "1",
			status: http.StatusInternalServerError,
		},
		{
			name: "happy path",
			getLeaderboardEntry: func(quizID int, userID int) (repo.LeaderboardEntryDto, error) {
				return repo.LeaderboardEntryDto{}, nil
			},
			quizID: "1",
			userID: "1",
			status: http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetLeaderboardEntry = tc.getLeaderboardEntry

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"quizId": tc.quizID,
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

				var parsed repo.LeaderboardEntry
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestCreateEntry(t *testing.T) {
	savedGetUser := repo.GetUser
	savedValidUser := auth.ValidUser
	savedScoreExceedsMax := repo.ScoreExceedsMax
	savedInsertLeaderboardEntry := repo.InsertLeaderboardEntry

	defer func() {
		repo.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		repo.ScoreExceedsMax = savedScoreExceedsMax
		repo.InsertLeaderboardEntry = savedInsertLeaderboardEntry
	}()

	user := repo.UserDto{
		Username: "testing",
	}

	tt := []struct {
		name                   string
		getUser                func(id int) (repo.UserDto, error)
		validUser              func(request *http.Request, id int) (int, error)
		scoreExceedsMax        func(quizID, score int) (bool, error)
		insertLeaderboardEntry func(entry repo.LeaderboardEntry) (int, error)
		body                   string
		status                 int
	}{
		{
			name:                   "invalid body",
			getUser:                repo.GetUser,
			validUser:              auth.ValidUser,
			scoreExceedsMax:        repo.ScoreExceedsMax,
			insertLeaderboardEntry: repo.InsertLeaderboardEntry,
			body:                   "testing",
			status:                 http.StatusBadRequest,
		},
		{
			name:                   "valid body, error on GetUser",
			getUser:                func(id int) (repo.UserDto, error) { return repo.UserDto{}, errors.New("test") },
			validUser:              auth.ValidUser,
			scoreExceedsMax:        repo.ScoreExceedsMax,
			insertLeaderboardEntry: repo.InsertLeaderboardEntry,
			body:                   `{"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                 http.StatusInternalServerError,
		},
		{
			name:    "valid body, invalid user",
			getUser: func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			scoreExceedsMax:        repo.ScoreExceedsMax,
			insertLeaderboardEntry: repo.InsertLeaderboardEntry,
			body:                   `{"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                 http.StatusUnauthorized,
		},
		{
			name:    "valid body, valid user, error on ScoreExceedsMax",
			getUser: func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			scoreExceedsMax:        func(quizID, score int) (bool, error) { return false, errors.New("test") },
			insertLeaderboardEntry: repo.InsertLeaderboardEntry,
			body:                   `{"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                 http.StatusInternalServerError,
		},
		{
			name:    "valid body, valid user, score exceeds max",
			getUser: func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			scoreExceedsMax:        func(quizID, score int) (bool, error) { return true, nil },
			insertLeaderboardEntry: repo.InsertLeaderboardEntry,
			body:                   `{"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                 http.StatusBadRequest,
		},
		{
			name:    "valid body, valid user, error on InsertLeaderboardEntry",
			getUser: func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			scoreExceedsMax:        func(quizID, score int) (bool, error) { return false, nil },
			insertLeaderboardEntry: func(entry repo.LeaderboardEntry) (int, error) { return 0, errors.New("test") },
			body:                   `{"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                 http.StatusInternalServerError,
		},
		{
			name:    "happy path",
			getUser: func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			scoreExceedsMax:        func(quizID, score int) (bool, error) { return false, nil },
			insertLeaderboardEntry: func(entry repo.LeaderboardEntry) (int, error) { return 1, nil },
			body:                   `{"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                 http.StatusCreated,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			repo.ScoreExceedsMax = tc.scoreExceedsMax
			repo.InsertLeaderboardEntry = tc.insertLeaderboardEntry

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

				var parsed repo.LeaderboardEntry
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestUpdateEntry(t *testing.T) {
	savedGetUser := repo.GetUser
	savedValidUser := auth.ValidUser
	savedScoreExceedsMax := repo.ScoreExceedsMax
	savedGetLeaderboardEntry := repo.GetLeaderboardEntry
	savedUpdateLeaderboardEntry := repo.UpdateLeaderboardEntry

	defer func() {
		repo.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		repo.ScoreExceedsMax = savedScoreExceedsMax
		repo.GetLeaderboardEntry = savedGetLeaderboardEntry
		repo.UpdateLeaderboardEntry = savedUpdateLeaderboardEntry
	}()

	user := repo.UserDto{
		Username: "testing",
	}

	tt := []struct {
		name                   string
		getUser                func(id int) (repo.UserDto, error)
		validUser              func(request *http.Request, id int) (int, error)
		scoreExceedsMax        func(quizID, score int) (bool, error)
		getLeaderboardEntry    func(quizID, userID int) (repo.LeaderboardEntryDto, error)
		updateLeaderboardEntry func(entry repo.LeaderboardEntry) error
		id                     string
		body                   string
		status                 int
	}{
		{
			name:                   "invalid id",
			getUser:                repo.GetUser,
			validUser:              auth.ValidUser,
			scoreExceedsMax:        repo.ScoreExceedsMax,
			getLeaderboardEntry:    repo.GetLeaderboardEntry,
			updateLeaderboardEntry: repo.UpdateLeaderboardEntry,
			id:                     "testing",
			body:                   "",
			status:                 http.StatusBadRequest,
		},
		{
			name:                   "valid id, invalid body",
			getUser:                func(id int) (repo.UserDto, error) { return user, nil },
			validUser:              auth.ValidUser,
			scoreExceedsMax:        repo.ScoreExceedsMax,
			getLeaderboardEntry:    repo.GetLeaderboardEntry,
			updateLeaderboardEntry: repo.UpdateLeaderboardEntry,
			id:                     "1",
			body:                   "testing",
			status:                 http.StatusBadRequest,
		},
		{
			name:                   "valid id, valid body, error on GetUser",
			getUser:                func(id int) (repo.UserDto, error) { return repo.UserDto{}, errors.New("test") },
			validUser:              auth.ValidUser,
			scoreExceedsMax:        repo.ScoreExceedsMax,
			getLeaderboardEntry:    repo.GetLeaderboardEntry,
			updateLeaderboardEntry: repo.UpdateLeaderboardEntry,
			id:                     "1",
			body:                   `{"id": 1,"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                 http.StatusInternalServerError,
		},
		{
			name:    "valid id, valid body, invalid user",
			getUser: func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			scoreExceedsMax:        repo.ScoreExceedsMax,
			getLeaderboardEntry:    repo.GetLeaderboardEntry,
			updateLeaderboardEntry: repo.UpdateLeaderboardEntry,
			id:                     "1",
			body:                   `{"id": 1,"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                 http.StatusUnauthorized,
		},
		{
			name:    "valid id, valid body, valid user, error on ScoreExceedsMax",
			getUser: func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			scoreExceedsMax:        func(quizID, score int) (bool, error) { return false, errors.New("test") },
			getLeaderboardEntry:    repo.GetLeaderboardEntry,
			updateLeaderboardEntry: repo.UpdateLeaderboardEntry,
			id:                     "1",
			body:                   `{"id": 1,"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                 http.StatusInternalServerError,
		},
		{
			name:    "valid id, valid body, valid user, score exceeds max",
			getUser: func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			scoreExceedsMax:        func(quizID, score int) (bool, error) { return true, nil },
			getLeaderboardEntry:    repo.GetLeaderboardEntry,
			updateLeaderboardEntry: repo.UpdateLeaderboardEntry,
			id:                     "1",
			body:                   `{"id": 1,"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                 http.StatusBadRequest,
		},
		{
			name:    "valid id, valid body, valid user, no rows error on GetLeaderboardEntry",
			getUser: func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			scoreExceedsMax: func(quizID, score int) (bool, error) { return false, nil },
			getLeaderboardEntry: func(quizID, userID int) (repo.LeaderboardEntryDto, error) {
				return repo.LeaderboardEntryDto{}, sql.ErrNoRows
			},
			updateLeaderboardEntry: repo.UpdateLeaderboardEntry,
			id:                     "1",
			body:                   `{"id": 1,"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                 http.StatusBadRequest,
		},
		{
			name:    "valid id, valid body, valid user, other error on GetLeaderboardEntry",
			getUser: func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			scoreExceedsMax: func(quizID, score int) (bool, error) { return false, nil },
			getLeaderboardEntry: func(quizID, userID int) (repo.LeaderboardEntryDto, error) {
				return repo.LeaderboardEntryDto{}, errors.New("test")
			},
			updateLeaderboardEntry: repo.UpdateLeaderboardEntry,
			id:                     "1",
			body:                   `{"id": 1,"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                 http.StatusInternalServerError,
		},
		{
			name:    "valid id, valid body, valid user, error on UpdateLeaderboardEntry",
			getUser: func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			scoreExceedsMax: func(quizID, score int) (bool, error) { return false, nil },
			getLeaderboardEntry: func(quizID, userID int) (repo.LeaderboardEntryDto, error) {
				return repo.LeaderboardEntryDto{}, nil
			},
			updateLeaderboardEntry: func(entry repo.LeaderboardEntry) error { return errors.New("test") },
			id:                     "1",
			body:                   `{"id": 1,"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                 http.StatusInternalServerError,
		},
		{
			name:    "happy path",
			getUser: func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			scoreExceedsMax: func(quizID, score int) (bool, error) { return false, nil },
			getLeaderboardEntry: func(quizID, userID int) (repo.LeaderboardEntryDto, error) {
				return repo.LeaderboardEntryDto{}, nil
			},
			updateLeaderboardEntry: func(entry repo.LeaderboardEntry) error { return nil },
			id:                     "1",
			body:                   `{"id": 1,"userId": 1, "country": "New Zealand", "capitals": 100, "time": 200}`,
			status:                 http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			repo.ScoreExceedsMax = tc.scoreExceedsMax
			repo.GetLeaderboardEntry = tc.getLeaderboardEntry
			repo.UpdateLeaderboardEntry = tc.updateLeaderboardEntry

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

				var parsed repo.LeaderboardEntry
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestDeleteEntry(t *testing.T) {
	savedGetLeaderboardEntryById := repo.GetLeaderboardEntryById
	savedGetUser := repo.GetUser
	savedValidUser := auth.ValidUser
	savedDeleteLeaderboardEntry := repo.DeleteLeaderboardEntry

	defer func() {
		repo.GetLeaderboardEntryById = savedGetLeaderboardEntryById
		repo.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		repo.DeleteLeaderboardEntry = savedDeleteLeaderboardEntry
	}()

	user := repo.UserDto{
		Username: "testing",
	}

	tt := []struct {
		name                    string
		getLeaderboardEntryById func(entryID int) (repo.LeaderboardEntry, error)
		getUser                 func(id int) (repo.UserDto, error)
		validUser               func(request *http.Request, id int) (int, error)
		deleteLeaderboardEntry  func(entryID int) error
		id                      string
		status                  int
	}{
		{
			name:                    "invalid id",
			getLeaderboardEntryById: repo.GetLeaderboardEntryById,
			getUser:                 repo.GetUser,
			validUser:               auth.ValidUser,
			deleteLeaderboardEntry:  repo.DeleteLeaderboardEntry,
			id:                      "testing",
			status:                  http.StatusBadRequest,
		},
		{
			name: "valid id, error on GetLeaderboardEntry",
			getLeaderboardEntryById: func(entryID int) (repo.LeaderboardEntry, error) {
				return repo.LeaderboardEntry{}, errors.New("test")
			},
			getUser:                func(id int) (repo.UserDto, error) { return user, nil },
			validUser:              auth.ValidUser,
			deleteLeaderboardEntry: repo.DeleteLeaderboardEntry,
			id:                     "1",
			status:                 http.StatusInternalServerError,
		},
		{
			name: "valid id, entry found, error on GetUser",
			getLeaderboardEntryById: func(entryID int) (repo.LeaderboardEntry, error) {
				return repo.LeaderboardEntry{}, nil
			},
			getUser:                func(id int) (repo.UserDto, error) { return repo.UserDto{}, errors.New("test") },
			validUser:              auth.ValidUser,
			deleteLeaderboardEntry: repo.DeleteLeaderboardEntry,
			id:                     "1",
			status:                 http.StatusInternalServerError,
		},
		{
			name: "valid id, entry found, invalid user",
			getLeaderboardEntryById: func(entryID int) (repo.LeaderboardEntry, error) {
				return repo.LeaderboardEntry{}, nil
			},
			getUser: func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			deleteLeaderboardEntry: repo.DeleteLeaderboardEntry,
			id:                     "1",
			status:                 http.StatusUnauthorized,
		},
		{
			name: "valid id, entry found, valid user, error on DeleteLeaderboardEntry",
			getLeaderboardEntryById: func(entryID int) (repo.LeaderboardEntry, error) {
				return repo.LeaderboardEntry{}, nil
			},
			getUser: func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			deleteLeaderboardEntry: func(entryID int) error { return errors.New("test") },
			id:                     "1",
			status:                 http.StatusInternalServerError,
		},
		{
			name: "happy path",
			getLeaderboardEntryById: func(entryID int) (repo.LeaderboardEntry, error) {
				return repo.LeaderboardEntry{}, nil
			},
			getUser: func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			deleteLeaderboardEntry: func(entryID int) error { return nil },
			id:                     "1",
			status:                 http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetLeaderboardEntryById = tc.getLeaderboardEntryById
			repo.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			repo.DeleteLeaderboardEntry = tc.deleteLeaderboardEntry

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

				var parsed repo.LeaderboardEntry
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}
