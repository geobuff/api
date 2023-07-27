package src

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
	"time"

	"github.com/geobuff/api/repo"
	"github.com/geobuff/api/utils"
	"github.com/gorilla/mux"
	"github.com/sendgrid/rest"
)

func TestLogin(t *testing.T) {
	savedGetUserUsingEmail := repo.GetAuthUserUsingEmail
	savedBuildToken := buildToken

	defer func() {
		repo.GetAuthUserUsingEmail = savedGetUserUsingEmail
		buildToken = savedBuildToken
	}()

	tt := []struct {
		name              string
		getUserUsingEmail func(email string) (repo.AuthUserDto, error)
		buildToken        func(user repo.AuthUserDto) (string, error)
		body              string
		status            int
	}{
		{
			name:              "invalid body",
			getUserUsingEmail: repo.GetAuthUserUsingEmail,
			buildToken:        buildToken,
			body:              "testing",
			status:            http.StatusBadRequest,
		},
		{
			name:              "sql.ErrNoRows error on GetUserUsingEmail",
			getUserUsingEmail: func(email string) (repo.AuthUserDto, error) { return repo.AuthUserDto{}, sql.ErrNoRows },
			buildToken:        buildToken,
			body:              `{"email": "scrub@gmail.com", "password": "Password1!"}`,
			status:            http.StatusBadRequest,
		},
		{
			name:              "other error on GetUserUsingEmail",
			getUserUsingEmail: func(email string) (repo.AuthUserDto, error) { return repo.AuthUserDto{}, errors.New("test") },
			buildToken:        buildToken,
			body:              `{"email": "scrub@gmail.com", "password": "Password1!"}`,
			status:            http.StatusInternalServerError,
		},
		{
			name: "error on CompareHashAndPassword",
			getUserUsingEmail: func(email string) (repo.AuthUserDto, error) {
				return repo.AuthUserDto{
						PasswordHash: "testing",
					},
					nil
			},
			buildToken: buildToken,
			body:       `{"email": "scrub@gmail.com", "password": "Password1!"}`,
			status:     http.StatusBadRequest,
		},
		{
			name: "error on BuildToken",
			getUserUsingEmail: func(email string) (repo.AuthUserDto, error) {
				return repo.AuthUserDto{
						PasswordHash: "$2a$04$EPhTOaXYzAqV366oEUzNQOCGnfUWwdnsxPMGmsATA4ikOxBi48buW",
					},
					nil
			},
			buildToken: func(user repo.AuthUserDto) (string, error) { return "", errors.New("test") },
			body:       `{"email": "scrub@gmail.com", "password": "Password1!"}`,
			status:     http.StatusInternalServerError,
		},
		{
			name: "happy path",
			getUserUsingEmail: func(email string) (repo.AuthUserDto, error) {
				return repo.AuthUserDto{
						PasswordHash: "$2a$04$EPhTOaXYzAqV366oEUzNQOCGnfUWwdnsxPMGmsATA4ikOxBi48buW",
					},
					nil
			},
			buildToken: func(user repo.AuthUserDto) (string, error) { return "test", nil },
			body:       `{"email": "scrub@gmail.com", "password": "Password1!"}`,
			status:     http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetAuthUserUsingEmail = tc.getUserUsingEmail
			buildToken = tc.buildToken

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
	savedHashPassword := hashPassword
	savedInsertUser := repo.InsertUser
	savedGetUser := repo.GetUser

	defer func() {
		repo.UsernameExists = savedUsernameExists
		repo.EmailExists = savedEmailExists
		hashPassword = savedHashPassword
		repo.InsertUser = savedInsertUser
		repo.GetUser = savedGetUser
	}()

	utils.Init()

	tt := []struct {
		name           string
		usernameExists func(username string) (bool, error)
		emailExists    func(email string) (bool, error)
		hashPassword   func(password []byte) (string, error)
		insertUser     func(user repo.User) (int, error)
		getUser        func(id int) (repo.UserDto, error)
		body           string
		status         int
	}{
		{
			name:           "invalid body",
			usernameExists: repo.UsernameExists,
			emailExists:    repo.EmailExists,
			hashPassword:   hashPassword,
			insertUser:     repo.InsertUser,
			getUser:        repo.GetUser,
			body:           "testing",
			status:         http.StatusBadRequest,
		},
		{
			name:           "valid body, invalid struct",
			usernameExists: repo.UsernameExists,
			emailExists:    repo.EmailExists,
			hashPassword:   hashPassword,
			insertUser:     repo.InsertUser,
			getUser:        repo.GetUser,
			body:           `{"username": "test"}`,
			status:         http.StatusBadRequest,
		},
		{
			name:           "error on UsernameExists",
			usernameExists: func(username string) (bool, error) { return false, errors.New("test") },
			emailExists:    repo.EmailExists,
			hashPassword:   hashPassword,
			insertUser:     repo.InsertUser,
			getUser:        repo.GetUser,
			body:           `{"avatarId": 1, "username": "test", "email": "scrub@gmail.com", "countryCode": "nz", "password": "Password1!"}`,
			status:         http.StatusInternalServerError,
		},
		{
			name:           "username exists",
			usernameExists: func(username string) (bool, error) { return true, nil },
			emailExists:    repo.EmailExists,
			hashPassword:   hashPassword,
			insertUser:     repo.InsertUser,
			getUser:        repo.GetUser,
			body:           `{"avatarId": 1, "username": "test", "email": "scrub@gmail.com", "countryCode": "nz", "password": "Password1!"}`,
			status:         http.StatusBadRequest,
		},
		{
			name:           "error on EmailExists",
			usernameExists: func(username string) (bool, error) { return false, nil },
			emailExists:    func(email string) (bool, error) { return false, errors.New("test") },
			hashPassword:   hashPassword,
			insertUser:     repo.InsertUser,
			getUser:        repo.GetUser,
			body:           `{"avatarId": 1, "username": "test", "email": "scrub@gmail.com", "countryCode": "nz", "password": "Password1!"}`,
			status:         http.StatusInternalServerError,
		},
		{
			name:           "email exists",
			usernameExists: func(username string) (bool, error) { return false, nil },
			emailExists:    func(email string) (bool, error) { return true, nil },
			hashPassword:   hashPassword,
			insertUser:     repo.InsertUser,
			getUser:        repo.GetUser,
			body:           `{"avatarId": 1, "username": "test", "email": "scrub@gmail.com", "countryCode": "nz", "password": "Password1!"}`,
			status:         http.StatusBadRequest,
		},
		{
			name:           "error on HashPassword",
			usernameExists: func(username string) (bool, error) { return false, nil },
			emailExists:    func(email string) (bool, error) { return false, nil },
			hashPassword:   func(password []byte) (string, error) { return "", errors.New("test") },
			insertUser:     repo.InsertUser,
			getUser:        repo.GetUser,
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
			insertUser: func(user repo.User) (int, error) { return 0, errors.New("test") },
			getUser:    repo.GetUser,
			body:       `{"avatarId": 1, "username": "test", "email": "scrub@gmail.com", "countryCode": "nz", "password": "Password1!"}`,
			status:     http.StatusInternalServerError,
		},
		{
			name:           "error on GetUser",
			usernameExists: func(username string) (bool, error) { return false, nil },
			emailExists:    func(email string) (bool, error) { return false, nil },
			hashPassword:   func(password []byte) (string, error) { return "test", nil },
			insertUser:     func(user repo.User) (int, error) { return 1, nil },
			getUser:        func(id int) (repo.UserDto, error) { return repo.UserDto{}, errors.New("test") },
			body:           `{"avatarId": 1, "username": "test", "email": "scrub@gmail.com", "countryCode": "nz", "password": "Password1!"}`,
			status:         http.StatusInternalServerError,
		},
		{
			name:           "happy path",
			usernameExists: func(username string) (bool, error) { return false, nil },
			emailExists:    func(email string) (bool, error) { return false, nil },
			hashPassword:   func(password []byte) (string, error) { return "test", nil },
			insertUser:     func(user repo.User) (int, error) { return 1, nil },
			getUser:        func(id int) (repo.UserDto, error) { return repo.UserDto{}, nil },
			body:           `{"avatarId": 1, "username": "test", "email": "scrub@gmail.com", "countryCode": "nz", "password": "Password1!"}`,
			status:         http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.UsernameExists = tc.usernameExists
			repo.EmailExists = tc.emailExists
			hashPassword = tc.hashPassword
			repo.InsertUser = tc.insertUser
			repo.GetUser = tc.getUser

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
		})
	}
}

func TestSendResetToken(t *testing.T) {
	savedGetUserUsingEmail := repo.GetAuthUserUsingEmail
	savedSetPasswordResetValues := repo.SetPasswordResetValues
	savedSendResetToken := utils.SendResetToken

	defer func() {
		repo.GetAuthUserUsingEmail = savedGetUserUsingEmail
		repo.SetPasswordResetValues = savedSetPasswordResetValues
		utils.SendResetToken = savedSendResetToken
	}()

	tt := []struct {
		name                   string
		getUserUsingEmail      func(email string) (repo.AuthUserDto, error)
		setPasswordResetValues func(userID int, resetToken string, expiryDate time.Time) error
		sendResetToken         func(email, resetLink string) (*rest.Response, error)
		body                   string
		status                 int
	}{
		{
			name:                   "invalid body",
			getUserUsingEmail:      repo.GetAuthUserUsingEmail,
			setPasswordResetValues: repo.SetPasswordResetValues,
			sendResetToken:         utils.SendResetToken,
			body:                   "testing",
			status:                 http.StatusBadRequest,
		},
		{
			name:                   "sql.ErrNoRows on GetUserUsingEmail",
			getUserUsingEmail:      func(email string) (repo.AuthUserDto, error) { return repo.AuthUserDto{}, sql.ErrNoRows },
			setPasswordResetValues: repo.SetPasswordResetValues,
			sendResetToken:         utils.SendResetToken,
			body:                   `{"email": "scrub@gmail.com"}`,
			status:                 http.StatusBadRequest,
		},
		{
			name:                   "other error on GetUserUsingEmail",
			getUserUsingEmail:      func(email string) (repo.AuthUserDto, error) { return repo.AuthUserDto{}, errors.New("test") },
			setPasswordResetValues: repo.SetPasswordResetValues,
			sendResetToken:         utils.SendResetToken,
			body:                   `{"email": "scrub@gmail.com"}`,
			status:                 http.StatusInternalServerError,
		},
		{
			name:                   "error on SetPasswordResetValues",
			getUserUsingEmail:      func(email string) (repo.AuthUserDto, error) { return repo.AuthUserDto{}, nil },
			setPasswordResetValues: func(userID int, resetToken string, expiryDate time.Time) error { return errors.New("test") },
			sendResetToken:         utils.SendResetToken,
			body:                   `{"email": "scrub@gmail.com"}`,
			status:                 http.StatusInternalServerError,
		},
		{
			name:                   "error on SendResetToken",
			getUserUsingEmail:      func(email string) (repo.AuthUserDto, error) { return repo.AuthUserDto{}, nil },
			setPasswordResetValues: func(userID int, resetToken string, expiryDate time.Time) error { return nil },
			sendResetToken:         func(email, resetLink string) (*rest.Response, error) { return nil, errors.New("test") },
			body:                   `{"email": "scrub@gmail.com"}`,
			status:                 http.StatusInternalServerError,
		},
		{
			name:                   "happy path",
			getUserUsingEmail:      func(email string) (repo.AuthUserDto, error) { return repo.AuthUserDto{}, nil },
			setPasswordResetValues: func(userID int, resetToken string, expiryDate time.Time) error { return nil },
			sendResetToken:         func(email, resetLink string) (*rest.Response, error) { return nil, nil },
			body:                   `{"email": "scrub@gmail.com"}`,
			status:                 http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetAuthUserUsingEmail = tc.getUserUsingEmail
			repo.SetPasswordResetValues = tc.setPasswordResetValues
			utils.SendResetToken = tc.sendResetToken

			request, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(tc.body)))
			if err != nil {
				t.Fatalf("could not create POST request: %v", err)
			}

			writer := httptest.NewRecorder()
			SendResetToken(writer, request)
			result := writer.Result()
			defer result.Body.Close()

			if result.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, result.StatusCode)
			}
		})
	}
}

func TestResetTokenValid(t *testing.T) {
	savedGetUser := repo.GetAuthUser
	savedIsResetTokenValid := isResetTokenValid

	defer func() {
		repo.GetAuthUser = savedGetUser
		isResetTokenValid = savedIsResetTokenValid
	}()

	tt := []struct {
		name              string
		getUser           func(id int) (repo.AuthUserDto, error)
		isResetTokenValid func(userToken sql.NullString, requestToken string, expiry sql.NullTime) bool
		userId            string
		status            int
	}{
		{
			name:              "invalid userId",
			getUser:           repo.GetAuthUser,
			isResetTokenValid: isResetTokenValid,
			userId:            "test",
			status:            http.StatusBadRequest,
		},
		{
			name:              "sql.ErrNoRows error on GetUser",
			getUser:           func(id int) (repo.AuthUserDto, error) { return repo.AuthUserDto{}, sql.ErrNoRows },
			isResetTokenValid: isResetTokenValid,
			userId:            "1",
			status:            http.StatusBadRequest,
		},
		{
			name:              "other error on GetUser",
			getUser:           func(id int) (repo.AuthUserDto, error) { return repo.AuthUserDto{}, errors.New("test") },
			isResetTokenValid: isResetTokenValid,
			userId:            "1",
			status:            http.StatusInternalServerError,
		},
		{
			name:              "happy path",
			getUser:           func(id int) (repo.AuthUserDto, error) { return repo.AuthUserDto{}, nil },
			isResetTokenValid: func(userToken sql.NullString, requestToken string, expiry sql.NullTime) bool { return true },
			userId:            "1",
			status:            http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetAuthUser = tc.getUser
			isResetTokenValid = tc.isResetTokenValid

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"userId": tc.userId,
			})

			writer := httptest.NewRecorder()
			ResetTokenValid(writer, request)
			result := writer.Result()
			defer result.Body.Close()

			if result.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, result.StatusCode)
			}

			if result.StatusCode == http.StatusOK {
				body, err := ioutil.ReadAll(result.Body)
				if err != nil {
					t.Fatalf("could not read response: %v", err)
				}

				var parsed bool
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshall response body: %v", err)
				}
			}
		})
	}
}

func TestUpdatePasswordUsingToken(t *testing.T) {
	savedGetuser := repo.GetAuthUser
	savedIsResetTokenValid := isResetTokenValid
	savedHashPassword := hashPassword
	savedResetPassword := repo.ResetPassword

	defer func() {
		repo.GetAuthUser = savedGetuser
		isResetTokenValid = savedIsResetTokenValid
		hashPassword = savedHashPassword
		repo.ResetPassword = savedResetPassword
	}()

	tt := []struct {
		name              string
		getUser           func(id int) (repo.AuthUserDto, error)
		isResetTokenValid func(userToken sql.NullString, requestToken string, expiry sql.NullTime) bool
		hashPassword      func(password []byte) (string, error)
		resetPassword     func(userID int, passwordHash string) error
		body              string
		status            int
	}{
		{
			name:              "invalid body",
			getUser:           repo.GetAuthUser,
			isResetTokenValid: isResetTokenValid,
			hashPassword:      hashPassword,
			resetPassword:     repo.ResetPassword,
			body:              "testing",
			status:            http.StatusBadRequest,
		},
		{
			name:              "sql.ErrNoRows on GetUser",
			getUser:           func(id int) (repo.AuthUserDto, error) { return repo.AuthUserDto{}, sql.ErrNoRows },
			isResetTokenValid: isResetTokenValid,
			hashPassword:      hashPassword,
			resetPassword:     repo.ResetPassword,
			body:              `{"userId": 1, "token": "test", "password": "Password1!"}`,
			status:            http.StatusBadRequest,
		},
		{
			name:              "other error on GetUser",
			getUser:           func(id int) (repo.AuthUserDto, error) { return repo.AuthUserDto{}, errors.New("test") },
			isResetTokenValid: isResetTokenValid,
			hashPassword:      hashPassword,
			resetPassword:     repo.ResetPassword,
			body:              `{"userId": 1, "token": "test", "password": "Password1!"}`,
			status:            http.StatusInternalServerError,
		},
		{
			name:              "reset token invalid",
			getUser:           func(id int) (repo.AuthUserDto, error) { return repo.AuthUserDto{}, nil },
			isResetTokenValid: func(userToken sql.NullString, requestToken string, expiry sql.NullTime) bool { return false },
			hashPassword:      hashPassword,
			resetPassword:     repo.ResetPassword,
			body:              `{"userId": 1, "token": "test", "password": "Password1!"}`,
			status:            http.StatusBadRequest,
		},
		{
			name:              "error on HashPassword",
			getUser:           func(id int) (repo.AuthUserDto, error) { return repo.AuthUserDto{}, nil },
			isResetTokenValid: func(userToken sql.NullString, requestToken string, expiry sql.NullTime) bool { return true },
			hashPassword:      func(password []byte) (string, error) { return "", errors.New("test") },
			resetPassword:     repo.ResetPassword,
			body:              `{"userId": 1, "token": "test", "password": "Password1!"}`,
			status:            http.StatusInternalServerError,
		},
		{
			name:              "error on ResetPassword",
			getUser:           func(id int) (repo.AuthUserDto, error) { return repo.AuthUserDto{}, nil },
			isResetTokenValid: func(userToken sql.NullString, requestToken string, expiry sql.NullTime) bool { return true },
			hashPassword:      func(password []byte) (string, error) { return "test", nil },
			resetPassword:     func(userID int, passwordHash string) error { return errors.New("test") },
			body:              `{"userId": 1, "token": "test", "password": "Password1!"}`,
			status:            http.StatusInternalServerError,
		},
		{
			name:              "happy path",
			getUser:           func(id int) (repo.AuthUserDto, error) { return repo.AuthUserDto{}, nil },
			isResetTokenValid: func(userToken sql.NullString, requestToken string, expiry sql.NullTime) bool { return true },
			hashPassword:      func(password []byte) (string, error) { return "test", nil },
			resetPassword:     func(userID int, passwordHash string) error { return nil },
			body:              `{"userId": 1, "token": "test", "password": "Password1!"}`,
			status:            http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetAuthUser = tc.getUser
			isResetTokenValid = tc.isResetTokenValid
			hashPassword = tc.hashPassword
			repo.ResetPassword = tc.resetPassword

			request, err := http.NewRequest("PUT", "", bytes.NewBuffer([]byte(tc.body)))
			if err != nil {
				t.Fatalf("could not create PUT request: %v", err)
			}

			writer := httptest.NewRecorder()
			UpdatePasswordUsingToken(writer, request)
			result := writer.Result()
			defer result.Body.Close()

			if result.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, result.StatusCode)
			}

			if result.StatusCode == http.StatusOK {
				body, err := ioutil.ReadAll(result.Body)
				if err != nil {
					t.Fatalf("could not read response: %v", err)
				}

				var parsed repo.UserDto
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshall response body: %v", err)
				}
			}
		})
	}
}

func TestIsResetTokenValid(t *testing.T) {
	tt := []struct {
		name         string
		userToken    sql.NullString
		requestToken string
		expiry       sql.NullTime
		expected     bool
	}{
		{
			name: "userToken invalid",
			userToken: sql.NullString{
				String: "",
				Valid:  false,
			},
			requestToken: "",
			expiry:       sql.NullTime{},
			expected:     false,
		},
		{
			name: "expiry invalid",
			userToken: sql.NullString{
				String: "test",
				Valid:  true,
			},
			requestToken: "",
			expiry: sql.NullTime{
				Time:  time.Now(),
				Valid: false,
			},
			expected: false,
		},
		{
			name: "userToken different to requestToken",
			userToken: sql.NullString{
				String: "test",
				Valid:  true,
			},
			requestToken: "different",
			expiry: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			expected: false,
		},
		{
			name: "expiry less than current time",
			userToken: sql.NullString{
				String: "test",
				Valid:  true,
			},
			requestToken: "test",
			expiry: sql.NullTime{
				Time:  time.Now().AddDate(0, 0, -1),
				Valid: true,
			},
			expected: false,
		},
		{
			name: "happy path",
			userToken: sql.NullString{
				String: "test",
				Valid:  true,
			},
			requestToken: "test",
			expiry: sql.NullTime{
				Time:  time.Now().AddDate(0, 0, 1),
				Valid: true,
			},
			expected: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := isResetTokenValid(tc.userToken, tc.requestToken, tc.expiry)
			if result != tc.expected {
				t.Errorf("expected %v; got %v", tc.expected, result)
			}
		})
	}
}

func TestValidUser(t *testing.T) {
	savedGetToken := getToken
	savedGetClaims := getClaims

	defer func() {
		getToken = savedGetToken
		getClaims = savedGetClaims
	}()

	tt := []struct {
		name      string
		getToken  func(request *http.Request) (string, error)
		getClaims func(tokenString string) (*CustomClaims, error)
		id        int
		expected  int
	}{
		{
			name:      "error on getToken",
			getToken:  func(request *http.Request) (string, error) { return "", errors.New("test") },
			getClaims: getClaims,
			id:        0,
			expected:  http.StatusInternalServerError,
		},
		{
			name:      "error on getClaims",
			getToken:  func(request *http.Request) (string, error) { return "", nil },
			getClaims: func(tokenString string) (*CustomClaims, error) { return nil, errors.New("test") },
			id:        0,
			expected:  http.StatusInternalServerError,
		},
		{
			name:      "claims userId does not equal id",
			getToken:  func(request *http.Request) (string, error) { return "", nil },
			getClaims: func(tokenString string) (*CustomClaims, error) { return &CustomClaims{UserID: 2}, nil },
			id:        1,
			expected:  http.StatusUnauthorized,
		},
		{
			name:     "claims userId does not equal id but user is admin",
			getToken: func(request *http.Request) (string, error) { return "", nil },
			getClaims: func(tokenString string) (*CustomClaims, error) {
				return &CustomClaims{UserID: 2, IsAdmin: true}, nil
			},
			id:       1,
			expected: http.StatusOK,
		},
		{
			name:      "claims userId equals id",
			getToken:  func(request *http.Request) (string, error) { return "", nil },
			getClaims: func(tokenString string) (*CustomClaims, error) { return &CustomClaims{UserID: 1}, nil },
			id:        1,
			expected:  http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			getToken = tc.getToken
			getClaims = tc.getClaims

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			status, _ := ValidUser(request, tc.id)

			if status != tc.expected {
				t.Errorf("expected status %v; got %v", tc.expected, status)
			}
		})
	}
}

func TestIsAdmin(t *testing.T) {
	savedGetToken := getToken
	savedGetClaims := getClaims

	defer func() {
		getToken = savedGetToken
		getClaims = savedGetClaims
	}()

	tt := []struct {
		name      string
		getToken  func(request *http.Request) (string, error)
		getClaims func(tokenString string) (*CustomClaims, error)
		expected  int
	}{
		{
			name:      "error on getToken",
			getToken:  func(request *http.Request) (string, error) { return "", errors.New("test") },
			getClaims: getClaims,
			expected:  http.StatusInternalServerError,
		},
		{
			name:      "error on getClaims",
			getToken:  func(request *http.Request) (string, error) { return "", nil },
			getClaims: func(tokenString string) (*CustomClaims, error) { return &CustomClaims{}, errors.New("test") },
			expected:  http.StatusInternalServerError,
		},
		{
			name:      "user not admin",
			getToken:  func(request *http.Request) (string, error) { return "", nil },
			getClaims: func(tokenString string) (*CustomClaims, error) { return &CustomClaims{IsAdmin: false}, nil },
			expected:  http.StatusUnauthorized,
		},
		{
			name:      "user is admin",
			getToken:  func(request *http.Request) (string, error) { return "", nil },
			getClaims: func(tokenString string) (*CustomClaims, error) { return &CustomClaims{IsAdmin: true}, nil },
			expected:  http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			getToken = tc.getToken
			getClaims = tc.getClaims

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			status, _ := IsAdmin(request)
			if status != tc.expected {
				t.Errorf("expected status %v; got %v", tc.expected, status)
			}
		})
	}
}

func TestGetToken(t *testing.T) {
	tt := []struct {
		name          string
		token         string
		expectedToken string
		expectedError string
	}{
		{
			name:          "token invalid length",
			token:         "",
			expectedToken: "",
			expectedError: "token missing or invalid length",
		},
		{
			name:          "token valid",
			token:         "abcdefgh",
			expectedToken: "abcdefgh",
			expectedError: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tc.token))

			token, err := getToken(request)
			if token != tc.expectedToken {
				t.Errorf("expected token %s; got %s", tc.expectedToken, token)
			}

			if err != nil && err.Error() != tc.expectedError {
				t.Errorf("expected err %v; got %v", tc.expectedError, err)
			}
		})
	}

}