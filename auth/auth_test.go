package auth

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geobuff/api/repo"
	"github.com/geobuff/api/validation"
)

func TestLogin(t *testing.T) {
	savedGetUserUsingEmail := repo.GetUserUsingEmail
	savedBuildToken := BuildToken

	defer func() {
		repo.GetUserUsingEmail = savedGetUserUsingEmail
		BuildToken = savedBuildToken
	}()

	tt := []struct {
		name              string
		getUserUsingEmail func(email string) (repo.User, error)
		buildToken        func(user repo.User) (string, error)
		body              string
		status            int
	}{
		{
			name:              "invalid body",
			getUserUsingEmail: repo.GetUserUsingEmail,
			buildToken:        BuildToken,
			body:              "testing",
			status:            http.StatusBadRequest,
		},
		{
			name:              "sql.ErrNoRows error on GetUserUsingEmail",
			getUserUsingEmail: func(email string) (repo.User, error) { return repo.User{}, sql.ErrNoRows },
			buildToken:        BuildToken,
			body:              `{"email": "scrub@gmail.com", "password": "Password1!"}`,
			status:            http.StatusBadRequest,
		},
		{
			name:              "other error on GetUserUsingEmail",
			getUserUsingEmail: func(email string) (repo.User, error) { return repo.User{}, errors.New("test") },
			buildToken:        BuildToken,
			body:              `{"email": "scrub@gmail.com", "password": "Password1!"}`,
			status:            http.StatusInternalServerError,
		},
		{
			name: "error on CompareHashAndPassword",
			getUserUsingEmail: func(email string) (repo.User, error) {
				return repo.User{
						PasswordHash: "testing",
					},
					nil
			},
			buildToken: BuildToken,
			body:       `{"email": "scrub@gmail.com", "password": "Password1!"}`,
			status:     http.StatusBadRequest,
		},
		{
			name: "error on BuildToken",
			getUserUsingEmail: func(email string) (repo.User, error) {
				return repo.User{
						PasswordHash: "$2a$04$EPhTOaXYzAqV366oEUzNQOCGnfUWwdnsxPMGmsATA4ikOxBi48buW",
					},
					nil
			},
			buildToken: func(user repo.User) (string, error) { return "", errors.New("test") },
			body:       `{"email": "scrub@gmail.com", "password": "Password1!"}`,
			status:     http.StatusInternalServerError,
		},
		{
			name: "happy path",
			getUserUsingEmail: func(email string) (repo.User, error) {
				return repo.User{
						PasswordHash: "$2a$04$EPhTOaXYzAqV366oEUzNQOCGnfUWwdnsxPMGmsATA4ikOxBi48buW",
					},
					nil
			},
			buildToken: func(user repo.User) (string, error) { return "test", nil },
			body:       `{"email": "scrub@gmail.com", "password": "Password1!"}`,
			status:     http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetUserUsingEmail = tc.getUserUsingEmail
			BuildToken = tc.buildToken

			request, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(tc.body)))
			if err != nil {
				t.Fatalf("could not create POST request: %v", err)
			}

			writer := httptest.NewRecorder()
			Login(writer, request)
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

				var parsed string
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestRegister(t *testing.T) {
	savedUsernameExists := repo.UsernameExists
	savedEmailExists := repo.EmailExists
	savedHashPassword := HashPassword
	savedInsertUser := repo.InsertUser
	savedBuildToken := BuildToken

	defer func() {
		repo.UsernameExists = savedUsernameExists
		repo.EmailExists = savedEmailExists
		HashPassword = savedHashPassword
		repo.InsertUser = savedInsertUser
		BuildToken = savedBuildToken
	}()

	validation.Init()

	tt := []struct {
		name           string
		usernameExists func(username string) (bool, error)
		emailExists    func(email string) (bool, error)
		hashPassword   func(password []byte) (string, error)
		insertUser     func(user repo.User) (repo.User, error)
		buildToken     func(user repo.User) (string, error)
		body           string
		status         int
	}{
		{
			name:           "invalid body",
			usernameExists: repo.UsernameExists,
			emailExists:    repo.EmailExists,
			hashPassword:   HashPassword,
			insertUser:     repo.InsertUser,
			buildToken:     BuildToken,
			body:           "testing",
			status:         http.StatusBadRequest,
		},
		{
			name:           "valid body, invalid struct",
			usernameExists: repo.UsernameExists,
			emailExists:    repo.EmailExists,
			hashPassword:   HashPassword,
			insertUser:     repo.InsertUser,
			buildToken:     BuildToken,
			body:           `{"username": "test"}`,
			status:         http.StatusBadRequest,
		},
		{
			name:           "error on UsernameExists",
			usernameExists: func(username string) (bool, error) { return false, errors.New("test") },
			emailExists:    repo.EmailExists,
			hashPassword:   HashPassword,
			insertUser:     repo.InsertUser,
			buildToken:     BuildToken,
			body:           `{"avatarId": 1, "username": "test", "email": "scrub@gmail.com", "countryCode": "nz", "password": "Password1!"}`,
			status:         http.StatusInternalServerError,
		},
		{
			name:           "username exists",
			usernameExists: func(username string) (bool, error) { return true, nil },
			emailExists:    repo.EmailExists,
			hashPassword:   HashPassword,
			insertUser:     repo.InsertUser,
			buildToken:     BuildToken,
			body:           `{"avatarId": 1, "username": "test", "email": "scrub@gmail.com", "countryCode": "nz", "password": "Password1!"}`,
			status:         http.StatusBadRequest,
		},
		{
			name:           "error on EmailExists",
			usernameExists: func(username string) (bool, error) { return false, nil },
			emailExists:    func(email string) (bool, error) { return false, errors.New("test") },
			hashPassword:   HashPassword,
			insertUser:     repo.InsertUser,
			buildToken:     BuildToken,
			body:           `{"avatarId": 1, "username": "test", "email": "scrub@gmail.com", "countryCode": "nz", "password": "Password1!"}`,
			status:         http.StatusInternalServerError,
		},
		{
			name:           "email exists",
			usernameExists: func(username string) (bool, error) { return false, nil },
			emailExists:    func(email string) (bool, error) { return true, nil },
			hashPassword:   HashPassword,
			insertUser:     repo.InsertUser,
			buildToken:     BuildToken,
			body:           `{"avatarId": 1, "username": "test", "email": "scrub@gmail.com", "countryCode": "nz", "password": "Password1!"}`,
			status:         http.StatusBadRequest,
		},
		{
			name:           "error on HashPassword",
			usernameExists: func(username string) (bool, error) { return false, nil },
			emailExists:    func(email string) (bool, error) { return false, nil },
			hashPassword:   func(password []byte) (string, error) { return "", errors.New("test") },
			insertUser:     repo.InsertUser,
			buildToken:     BuildToken,
			body:           `{"avatarId": 1, "username": "test", "email": "scrub@gmail.com", "countryCode": "nz", "password": "Password1!"}`,
			status:         http.StatusInternalServerError,
		},
		{
			name:           "error on InsertUser",
			usernameExists: func(username string) (bool, error) { return false, nil },
			emailExists:    func(email string) (bool, error) { return false, nil },
			hashPassword: func(password []byte) (string, error) {
				return "$2a$04$EPhTOaXYzAqV366oEUzNQOCGnfUWwdnsxPMGmsATA4ikOxBi48buW", nil
			},
			insertUser: func(user repo.User) (repo.User, error) { return repo.User{}, errors.New("test") },
			buildToken: BuildToken,
			body:       `{"avatarId": 1, "username": "test", "email": "scrub@gmail.com", "countryCode": "nz", "password": "Password1!"}`,
			status:     http.StatusInternalServerError,
		},
		{
			name:           "error on BuildToken",
			usernameExists: func(username string) (bool, error) { return false, nil },
			emailExists:    func(email string) (bool, error) { return false, nil },
			hashPassword:   func(password []byte) (string, error) { return "test", nil },
			insertUser:     func(user repo.User) (repo.User, error) { return repo.User{}, nil },
			buildToken:     func(user repo.User) (string, error) { return "", errors.New("test") },
			body:           `{"avatarId": 1, "username": "test", "email": "scrub@gmail.com", "countryCode": "nz", "password": "Password1!"}`,
			status:         http.StatusInternalServerError,
		},
		{
			name:           "happy path",
			usernameExists: func(username string) (bool, error) { return false, nil },
			emailExists:    func(email string) (bool, error) { return false, nil },
			hashPassword:   func(password []byte) (string, error) { return "test", nil },
			insertUser:     func(user repo.User) (repo.User, error) { return repo.User{}, nil },
			buildToken:     func(user repo.User) (string, error) { return "test", nil },
			body:           `{"avatarId": 1, "username": "test", "email": "scrub@gmail.com", "countryCode": "nz", "password": "Password1!"}`,
			status:         http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.UsernameExists = tc.usernameExists
			repo.EmailExists = tc.emailExists
			HashPassword = tc.hashPassword
			repo.InsertUser = tc.insertUser
			BuildToken = tc.buildToken

			request, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(tc.body)))
			if err != nil {
				t.Fatalf("could not create POST request: %v", err)
			}

			writer := httptest.NewRecorder()
			Register(writer, request)
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

				var parsed string
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}
