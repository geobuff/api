package users

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

func TestGetUsers(t *testing.T) {
	savedHasPermission := auth.HasPermission
	savedGetUsers := database.GetUsers
	savedGetUserID := database.GetFirstID

	defer func() {
		auth.HasPermission = savedHasPermission
		database.GetUsers = savedGetUsers
		database.GetFirstID = savedGetUserID
	}()

	tt := []struct {
		name          string
		hasPermission func(up auth.UserPermission) (bool, error)
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
			getUserID:     database.GetFirstID,
			page:          "testing",
			status:        http.StatusBadRequest,
			hasMore:       false,
		},
		{
			name:          "valid page parameter, error on HasPermission",
			hasPermission: auth.HasPermission,
			getUsers:      database.GetUsers,
			getUserID:     database.GetFirstID,
			page:          "0",
			status:        http.StatusInternalServerError,
			hasMore:       false,
		},
		{
			name:          "valid page parameter, invalid permissions",
			hasPermission: func(up auth.UserPermission) (bool, error) { return false, nil },
			getUsers:      database.GetUsers,
			getUserID:     database.GetFirstID,
			page:          "0",
			status:        http.StatusUnauthorized,
			hasMore:       false,
		},
		{
			name:          "valid page parameter, valid permissions, error on GetUsers",
			hasPermission: func(up auth.UserPermission) (bool, error) { return true, nil },
			getUsers:      func(limit int, offset int) ([]database.User, error) { return nil, errors.New("test") },
			getUserID:     database.GetFirstID,
			page:          "0",
			status:        http.StatusInternalServerError,
			hasMore:       false,
		},
		{
			name:          "valid page parameter, valid permissions, error on GetUserID",
			hasPermission: func(up auth.UserPermission) (bool, error) { return true, nil },
			getUsers:      func(limit int, offset int) ([]database.User, error) { return []database.User{}, nil },
			getUserID:     func(limit int, offset int) (int, error) { return 0, errors.New("test") },
			page:          "0",
			status:        http.StatusInternalServerError,
			hasMore:       false,
		},
		{
			name:          "happy path, has more is false",
			hasPermission: func(up auth.UserPermission) (bool, error) { return true, nil },
			getUsers:      func(limit int, offset int) ([]database.User, error) { return []database.User{}, nil },
			getUserID:     func(limit int, offset int) (int, error) { return 0, sql.ErrNoRows },
			page:          "0",
			status:        http.StatusOK,
			hasMore:       false,
		},
		{
			name:          "happy path, has more is true",
			hasPermission: func(up auth.UserPermission) (bool, error) { return true, nil },
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
			database.GetFirstID = tc.getUserID

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
	savedConfigValues := config.Values

	defer func() {
		auth.ValidUser = savedValidUser
		database.GetUser = savedGetUser
		config.Values = savedConfigValues
	}()

	tt := []struct {
		name      string
		validUser func(uv auth.UserValidation) (int, error)
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
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			getUser: database.GetUser,
			id:      "1",
			status:  http.StatusUnauthorized,
		},
		{
			name: "valid id, valid user, error on GetUser",
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			getUser: func(id int) (database.User, error) { return database.User{}, errors.New("test") },
			id:      "1",
			status:  http.StatusInternalServerError,
		},
		{
			name: "happy path",
			validUser: func(uv auth.UserValidation) (int, error) {
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
			config.Values = &config.Config{}

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

func TestGetUserID(t *testing.T) {
	savedGetUserID := database.GetUserID

	defer func() {
		database.GetUserID = savedGetUserID
	}()

	tt := []struct {
		name      string
		getUserID func(username string) (int, error)
		username  string
		status    int
	}{
		{
			name:      "error on GetUserID",
			getUserID: func(username string) (int, error) { return 0, errors.New("test") },
			username:  "testing",
			status:    http.StatusInternalServerError,
		},
		{
			name:      "error on GetUserID",
			getUserID: func(username string) (int, error) { return 0, sql.ErrNoRows },
			username:  "testing",
			status:    http.StatusNoContent,
		},
		{
			name:      "happy path",
			getUserID: func(username string) (int, error) { return 1, nil },
			username:  "testing",
			status:    http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetUserID = tc.getUserID

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"username": tc.username,
			})

			writer := httptest.NewRecorder()
			GetUserID(writer, request)
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

				var parsed int
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	savedGetKey := database.GetKey
	savedInsertUser := database.InsertUser

	defer func() {
		database.GetKey = savedGetKey
		database.InsertUser = savedInsertUser
	}()

	key := database.Key{
		ID:   1,
		Name: "valid key",
		Key:  "testkey123",
	}

	tt := []struct {
		name       string
		getKey     func(name string) (database.Key, error)
		insertUser func(user database.User) (int, error)
		key        string
		body       string
		status     int
	}{
		{
			name:       "error on GetKey",
			getKey:     func(name string) (database.Key, error) { return database.Key{}, errors.New("test") },
			insertUser: database.InsertUser,
			key:        "testkey123",
			body:       "",
			status:     http.StatusInternalServerError,
		},
		{
			name:       "invalid key",
			getKey:     func(name string) (database.Key, error) { return key, nil },
			insertUser: database.InsertUser,
			key:        "password1",
			body:       "",
			status:     http.StatusUnauthorized,
		},
		{
			name:       "valid key, invalid body",
			getKey:     func(name string) (database.Key, error) { return key, nil },
			insertUser: database.InsertUser,
			key:        "testkey123",
			body:       `testing`,
			status:     http.StatusBadRequest,
		},
		{
			name:       "valid key, valid body, error on InsertUser",
			getKey:     func(name string) (database.Key, error) { return key, nil },
			insertUser: func(user database.User) (int, error) { return 0, errors.New("test") },
			key:        "testkey123",
			body:       `{"username":"mrscrub"}`,
			status:     http.StatusInternalServerError,
		},
		{
			name:       "happy path",
			getKey:     func(name string) (database.Key, error) { return key, nil },
			insertUser: func(user database.User) (int, error) { return 1, nil },
			key:        "testkey123",
			body:       `{"username":"mrscrub"}`,
			status:     http.StatusCreated,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetKey = tc.getKey
			database.InsertUser = tc.insertUser

			request, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(tc.body)))
			if err != nil {
				t.Fatalf("could not create POST request: %v", err)
			}
			request.Header.Add("x-api-key", tc.key)

			writer := httptest.NewRecorder()
			CreateUser(writer, request)
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

				var parsed database.User
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	savedValidUser := auth.ValidUser
	savedUpdateUser := database.UpdateUser
	savedConfigValues := config.Values

	defer func() {
		auth.ValidUser = savedValidUser
		database.UpdateUser = savedUpdateUser
		config.Values = savedConfigValues
	}()

	tt := []struct {
		name       string
		validUser  func(uv auth.UserValidation) (int, error)
		updateUser func(user database.User) error
		id         string
		body       string
		status     int
	}{
		{
			name:       "invalid id",
			validUser:  auth.ValidUser,
			updateUser: database.UpdateUser,
			id:         "testing",
			body:       "",
			status:     http.StatusBadRequest,
		},
		{
			name:       "valid id, invalid user",
			validUser:  func(uv auth.UserValidation) (int, error) { return http.StatusUnauthorized, errors.New("test") },
			updateUser: database.UpdateUser,
			id:         "1",
			body:       "",
			status:     http.StatusUnauthorized,
		},
		{
			name:       "valid user, error on unmarshal",
			validUser:  func(uv auth.UserValidation) (int, error) { return http.StatusOK, nil },
			updateUser: database.UpdateUser,
			id:         "1",
			body:       "testing",
			status:     http.StatusBadRequest,
		},
		{
			name:       "valid user, valid body, error on UpdateUser",
			validUser:  func(uv auth.UserValidation) (int, error) { return http.StatusOK, nil },
			updateUser: func(user database.User) error { return errors.New("test") },
			id:         "1",
			body:       `{"username":"mrscrub"}`,
			status:     http.StatusInternalServerError,
		},
		{
			name:       "happy path",
			validUser:  func(uv auth.UserValidation) (int, error) { return http.StatusOK, nil },
			updateUser: func(user database.User) error { return nil },
			id:         "1",
			body:       `{"username":"mrscrub"}`,
			status:     http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			auth.ValidUser = tc.validUser
			database.UpdateUser = tc.updateUser
			config.Values = &config.Config{}

			request, err := http.NewRequest("PUT", "", bytes.NewBuffer([]byte(tc.body)))
			if err != nil {
				t.Fatalf("could not create PUT request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"id": tc.id,
			})

			writer := httptest.NewRecorder()
			UpdateUser(writer, request)
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

func TestDeleteUser(t *testing.T) {
	savedValidUser := auth.ValidUser
	savedGetUser := database.GetUser
	savedDeleteUser := database.DeleteUser
	savedConfigValues := config.Values

	defer func() {
		auth.ValidUser = savedValidUser
		database.GetUser = savedGetUser
		database.DeleteUser = savedDeleteUser
		config.Values = savedConfigValues
	}()

	tt := []struct {
		name       string
		validUser  func(uv auth.UserValidation) (int, error)
		getUser    func(id int) (database.User, error)
		deleteUser func(id int) error
		id         string
		status     int
	}{
		{
			name:       "invalid id",
			validUser:  auth.ValidUser,
			getUser:    database.GetUser,
			deleteUser: database.DeleteUser,
			id:         "testing",
			status:     http.StatusBadRequest,
		},
		{
			name: "valid id, invalid user",
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			getUser:    database.GetUser,
			deleteUser: database.DeleteUser,
			id:         "1",
			status:     http.StatusUnauthorized,
		},
		{
			name: "valid user, error on GetUser",
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			getUser:    func(id int) (database.User, error) { return database.User{}, errors.New("test") },
			deleteUser: database.DeleteUser,
			id:         "1",
			status:     http.StatusInternalServerError,
		},
		{
			name: "valid id, valid user, error on DeleteUser",
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			getUser:    func(id int) (database.User, error) { return database.User{}, nil },
			deleteUser: func(id int) error { return errors.New("test") },
			id:         "1",
			status:     http.StatusInternalServerError,
		},
		{
			name: "happy path",
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			getUser:    func(id int) (database.User, error) { return database.User{}, nil },
			deleteUser: func(id int) error { return nil },
			id:         "1",
			status:     http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			database.DeleteUser = tc.deleteUser
			config.Values = &config.Config{}

			request, err := http.NewRequest("DELETE", "", nil)
			if err != nil {
				t.Fatalf("could not create DELETE request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"id": tc.id,
			})

			writer := httptest.NewRecorder()
			DeleteUser(writer, request)
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
