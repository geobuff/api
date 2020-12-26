package users

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geobuff/geobuff-api/auth"
	"github.com/geobuff/geobuff-api/database"
	"github.com/gorilla/mux"
)

func TestGetUsers(t *testing.T) {
	savedHasPermission := auth.HasPermission
	savedGetUsers := database.GetUsers
	savedGetUserID := database.GetUserID

	defer func() {
		auth.HasPermission = savedHasPermission
		database.GetUsers = savedGetUsers
		database.GetUserID = savedGetUserID
	}()

	tt := []struct {
		name          string
		hasPermission func(request *http.Request, permission string) (bool, error)
		getUsers      func(limit int, offset int) ([]database.User, error)
		getUserID     func(limit int, offset int) (int, error)
		page          string
		status        int
		hasMore       bool
	}{
		{
			name:          "invalid page parameter",
			hasPermission: auth.HasPermission,
			getUsers:      database.GetUsers,
			getUserID:     database.GetUserID,
			page:          "testing",
			status:        http.StatusBadRequest,
			hasMore:       false,
		},
		{
			name:          "valid page parameter, error on HasPermission",
			hasPermission: auth.HasPermission,
			getUsers:      database.GetUsers,
			getUserID:     database.GetUserID,
			page:          "0",
			status:        http.StatusInternalServerError,
			hasMore:       false,
		},
		{
			name:          "valid page parameter, invalid permissions",
			hasPermission: func(request *http.Request, permission string) (bool, error) { return false, nil },
			getUsers:      database.GetUsers,
			getUserID:     database.GetUserID,
			page:          "0",
			status:        http.StatusUnauthorized,
			hasMore:       false,
		},
		{
			name:          "valid page parameter, valid permissions, error on GetUsers",
			hasPermission: func(request *http.Request, permission string) (bool, error) { return true, nil },
			getUsers:      func(limit int, offset int) ([]database.User, error) { return nil, errors.New("test") },
			getUserID:     database.GetUserID,
			page:          "0",
			status:        http.StatusInternalServerError,
			hasMore:       false,
		},
		{
			name:          "valid page parameter, valid permissions, error on GetUserID",
			hasPermission: func(request *http.Request, permission string) (bool, error) { return true, nil },
			getUsers:      func(limit int, offset int) ([]database.User, error) { return []database.User{}, nil },
			getUserID:     func(limit int, offset int) (int, error) { return 0, errors.New("test") },
			page:          "0",
			status:        http.StatusInternalServerError,
			hasMore:       false,
		},
		{
			name:          "happy path, has more is false",
			hasPermission: func(request *http.Request, permission string) (bool, error) { return true, nil },
			getUsers:      func(limit int, offset int) ([]database.User, error) { return []database.User{}, nil },
			getUserID:     func(limit int, offset int) (int, error) { return 0, sql.ErrNoRows },
			page:          "0",
			status:        http.StatusOK,
			hasMore:       false,
		},
		{
			name:          "happy path, has more is true",
			hasPermission: func(request *http.Request, permission string) (bool, error) { return true, nil },
			getUsers:      func(limit int, offset int) ([]database.User, error) { return []database.User{}, nil },
			getUserID:     func(limit int, offset int) (int, error) { return 1, nil },
			page:          "0",
			status:        http.StatusOK,
			hasMore:       true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			auth.HasPermission = tc.hasPermission
			database.GetUsers = tc.getUsers
			database.GetUserID = tc.getUserID

			request, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/users?page=%v", tc.page), nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			writer := httptest.NewRecorder()
			GetUsers(writer, request)
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

				var parsed PageDto
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

func TestGetUser(t *testing.T) {
	savedValidUser := auth.ValidUser
	savedGetUser := database.GetUser

	defer func() {
		auth.ValidUser = savedValidUser
		database.GetUser = savedGetUser
	}()

	tt := []struct {
		name      string
		validUser func(request *http.Request, userID int, permission string) (int, error)
		getUser   func(id int) (database.User, error)
		id        string
		status    int
	}{
		{
			name:      "invalid id",
			validUser: auth.ValidUser,
			getUser:   database.GetUser,
			id:        "testing",
			status:    http.StatusBadRequest,
		},
		{
			name: "valid id, invalid user",
			validUser: func(request *http.Request, userID int, permission string) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			getUser: database.GetUser,
			id:      "1",
			status:  http.StatusUnauthorized,
		},
		{
			name: "valid id, valid user, user not found",
			validUser: func(request *http.Request, userID int, permission string) (int, error) {
				return http.StatusOK, nil
			},
			getUser: func(id int) (database.User, error) { return database.User{}, sql.ErrNoRows },
			id:      "1",
			status:  http.StatusNotFound,
		},
		{
			name: "valid id, valid user, error on GetUser",
			validUser: func(request *http.Request, userID int, permission string) (int, error) {
				return http.StatusOK, nil
			},
			getUser: func(id int) (database.User, error) { return database.User{}, errors.New("test") },
			id:      "1",
			status:  http.StatusInternalServerError,
		},
		{
			name: "happy path",
			validUser: func(request *http.Request, userID int, permission string) (int, error) {
				return http.StatusOK, nil
			},
			getUser: func(id int) (database.User, error) { return database.User{}, nil },
			id:      "1",
			status:  http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			auth.ValidUser = tc.validUser
			database.GetUser = tc.getUser

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"id": tc.id,
			})

			writer := httptest.NewRecorder()
			GetUser(writer, request)
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

				var parsed database.User
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}
