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

	"github.com/geobuff/api/auth"
	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

func TestGetUsers(t *testing.T) {
	savedIsAdmin := auth.IsAdmin
	savedGetUsers := repo.GetUsers
	savedGetFirstID := repo.GetFirstID

	defer func() {
		auth.IsAdmin = savedIsAdmin
		repo.GetUsers = savedGetUsers
		repo.GetFirstID = savedGetFirstID
	}()

	tt := []struct {
		name       string
		isAdmin    func(request *http.Request) (int, error)
		getUsers   func(limit int, offset int) ([]repo.UserDto, error)
		getFirstID func(limit int, offset int) (int, error)
		page       string
		status     int
		hasMore    bool
	}{
		{
			name:       "invalid page parameter",
			isAdmin:    auth.IsAdmin,
			getUsers:   repo.GetUsers,
			getFirstID: repo.GetFirstID,
			page:       "testing",
			status:     http.StatusBadRequest,
			hasMore:    false,
		},
		{
			name:       "valid page parameter, error on HasPermission",
			isAdmin:    auth.IsAdmin,
			getUsers:   repo.GetUsers,
			getFirstID: repo.GetFirstID,
			page:       "0",
			status:     http.StatusInternalServerError,
			hasMore:    false,
		},
		{
			name:       "valid page parameter, invalid permissions",
			isAdmin:    func(request *http.Request) (int, error) { return http.StatusUnauthorized, errors.New("test") },
			getUsers:   repo.GetUsers,
			getFirstID: repo.GetFirstID,
			page:       "0",
			status:     http.StatusUnauthorized,
			hasMore:    false,
		},
		{
			name:       "valid page parameter, valid permissions, error on GetUsers",
			isAdmin:    func(request *http.Request) (int, error) { return http.StatusOK, nil },
			getUsers:   func(limit int, offset int) ([]repo.UserDto, error) { return nil, errors.New("test") },
			getFirstID: repo.GetFirstID,
			page:       "0",
			status:     http.StatusInternalServerError,
			hasMore:    false,
		},
		{
			name:       "valid page parameter, valid permissions, error on GetFirstID",
			isAdmin:    func(request *http.Request) (int, error) { return http.StatusOK, nil },
			getUsers:   func(limit int, offset int) ([]repo.UserDto, error) { return []repo.UserDto{}, nil },
			getFirstID: func(limit int, offset int) (int, error) { return 0, errors.New("test") },
			page:       "0",
			status:     http.StatusInternalServerError,
			hasMore:    false,
		},
		{
			name:       "happy path, has more is false",
			isAdmin:    func(request *http.Request) (int, error) { return http.StatusOK, nil },
			getUsers:   func(limit int, offset int) ([]repo.UserDto, error) { return []repo.UserDto{}, nil },
			getFirstID: func(limit int, offset int) (int, error) { return 0, sql.ErrNoRows },
			page:       "0",
			status:     http.StatusOK,
			hasMore:    false,
		},
		{
			name:       "happy path, has more is true",
			isAdmin:    func(request *http.Request) (int, error) { return http.StatusOK, nil },
			getUsers:   func(limit int, offset int) ([]repo.UserDto, error) { return []repo.UserDto{}, nil },
			getFirstID: func(limit int, offset int) (int, error) { return 1, nil },
			page:       "0",
			status:     http.StatusOK,
			hasMore:    true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			auth.IsAdmin = tc.isAdmin
			repo.GetUsers = tc.getUsers
			repo.GetFirstID = tc.getFirstID

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
	savedGetUser := repo.GetUser

	defer func() {
		auth.ValidUser = savedValidUser
		repo.GetUser = savedGetUser
	}()

	tt := []struct {
		name      string
		validUser func(request *http.Request, id int) (int, error)
		getUser   func(id int) (repo.UserDto, error)
		id        string
		status    int
	}{
		{
			name:      "invalid id",
			validUser: auth.ValidUser,
			getUser:   repo.GetUser,
			id:        "testing",
			status:    http.StatusBadRequest,
		},
		{
			name: "valid id, invalid user",
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			getUser: repo.GetUser,
			id:      "1",
			status:  http.StatusUnauthorized,
		},
		{
			name: "valid id, valid user, error on GetUser",
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			getUser: func(id int) (repo.UserDto, error) { return repo.UserDto{}, errors.New("test") },
			id:      "1",
			status:  http.StatusInternalServerError,
		},
		{
			name: "happy path",
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			getUser: func(id int) (repo.UserDto, error) { return repo.UserDto{}, nil },
			id:      "1",
			status:  http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			auth.ValidUser = tc.validUser
			repo.GetUser = tc.getUser

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

				var parsed repo.User
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
	savedUpdateUser := repo.UpdateUser

	defer func() {
		auth.ValidUser = savedValidUser
		repo.UpdateUser = savedUpdateUser
	}()

	tt := []struct {
		name       string
		validUser  func(request *http.Request, id int) (int, error)
		updateUser func(user repo.User) error
		id         string
		body       string
		status     int
	}{
		{
			name:       "invalid id",
			validUser:  auth.ValidUser,
			updateUser: repo.UpdateUser,
			id:         "testing",
			body:       "",
			status:     http.StatusBadRequest,
		},
		{
			name:       "valid id, invalid user",
			validUser:  func(request *http.Request, id int) (int, error) { return http.StatusUnauthorized, errors.New("test") },
			updateUser: repo.UpdateUser,
			id:         "1",
			body:       "",
			status:     http.StatusUnauthorized,
		},
		{
			name:       "valid user, error on unmarshal",
			validUser:  func(request *http.Request, id int) (int, error) { return http.StatusOK, nil },
			updateUser: repo.UpdateUser,
			id:         "1",
			body:       "testing",
			status:     http.StatusBadRequest,
		},
		{
			name:       "valid user, valid body, error on UpdateUser",
			validUser:  func(request *http.Request, id int) (int, error) { return http.StatusOK, nil },
			updateUser: func(user repo.User) error { return errors.New("test") },
			id:         "1",
			body:       `{"username":"mrscrub"}`,
			status:     http.StatusInternalServerError,
		},
		{
			name:       "happy path",
			validUser:  func(request *http.Request, id int) (int, error) { return http.StatusOK, nil },
			updateUser: func(user repo.User) error { return nil },
			id:         "1",
			body:       `{"username":"mrscrub"}`,
			status:     http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			auth.ValidUser = tc.validUser
			repo.UpdateUser = tc.updateUser

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

				var parsed repo.User
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
	savedGetUser := repo.GetUser
	savedDeleteUser := repo.DeleteUser

	defer func() {
		auth.ValidUser = savedValidUser
		repo.GetUser = savedGetUser
		repo.DeleteUser = savedDeleteUser
	}()

	tt := []struct {
		name       string
		validUser  func(request *http.Request, id int) (int, error)
		getUser    func(id int) (repo.UserDto, error)
		deleteUser func(id int) error
		id         string
		status     int
	}{
		{
			name:       "invalid id",
			validUser:  auth.ValidUser,
			getUser:    repo.GetUser,
			deleteUser: repo.DeleteUser,
			id:         "testing",
			status:     http.StatusBadRequest,
		},
		{
			name: "valid id, invalid user",
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			getUser:    repo.GetUser,
			deleteUser: repo.DeleteUser,
			id:         "1",
			status:     http.StatusUnauthorized,
		},
		{
			name: "valid user, error on GetUser",
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			getUser:    func(id int) (repo.UserDto, error) { return repo.UserDto{}, errors.New("test") },
			deleteUser: repo.DeleteUser,
			id:         "1",
			status:     http.StatusInternalServerError,
		},
		{
			name: "valid id, valid user, error on DeleteUser",
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			getUser:    func(id int) (repo.UserDto, error) { return repo.UserDto{}, nil },
			deleteUser: func(id int) error { return errors.New("test") },
			id:         "1",
			status:     http.StatusInternalServerError,
		},
		{
			name: "happy path",
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			getUser:    func(id int) (repo.UserDto, error) { return repo.UserDto{}, nil },
			deleteUser: func(id int) error { return nil },
			id:         "1",
			status:     http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			repo.DeleteUser = tc.deleteUser

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

				var parsed repo.User
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}
