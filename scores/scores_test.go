package scores

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
	"github.com/geobuff/auth"
	"github.com/gorilla/mux"
)

func TestGetScores(t *testing.T) {
	savedGetUser := database.GetUser
	savedValidUser := auth.ValidUser
	savedGetScores := database.GetScores
	savedConfigValues := config.Values

	defer func() {
		database.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		database.GetScores = savedGetScores
		config.Values = savedConfigValues
	}()

	user := database.User{
		Username: "testing",
	}

	tt := []struct {
		name      string
		getUser   func(id int) (database.User, error)
		validUser func(uv auth.UserValidation) (int, error)
		getScores func(userID int) ([]database.ScoreDto, error)
		userID    string
		status    int
	}{
		{
			name:      "invalid id value",
			getUser:   database.GetUser,
			validUser: auth.ValidUser,
			getScores: database.GetScores,
			userID:    "testing",
			status:    http.StatusBadRequest,
		},
		{
			name:      "valid id, error on GetUser",
			getUser:   func(id int) (database.User, error) { return database.User{}, errors.New("test") },
			validUser: auth.ValidUser,
			getScores: database.GetScores,
			userID:    "1",
			status:    http.StatusInternalServerError,
		},
		{
			name:    "valid id, invalid user",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			getScores: database.GetScores,
			userID:    "1",
			status:    http.StatusUnauthorized,
		},
		{
			name:    "valid id, valid user, error on GetScores",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			getScores: func(userID int) ([]database.ScoreDto, error) { return nil, errors.New("test") },
			userID:    "1",
			status:    http.StatusInternalServerError,
		},
		{
			name:    "happy path",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			getScores: func(userID int) ([]database.ScoreDto, error) { return []database.ScoreDto{}, nil },
			userID:    "1",
			status:    http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			database.GetScores = tc.getScores
			config.Values = &config.Config{}

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"userId": tc.userID,
			})

			writer := httptest.NewRecorder()
			GetScores(writer, request)
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

				var parsed []database.Score
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestGetScore(t *testing.T) {
	savedValidUser := auth.ValidUser
	savedGetUserScore := database.GetUserScore
	savedConfigValues := config.Values

	defer func() {
		auth.ValidUser = savedValidUser
		database.GetUserScore = savedGetUserScore
		config.Values = savedConfigValues
	}()

	tt := []struct {
		name         string
		validUser    func(uv auth.UserValidation) (int, error)
		getUserScore func(userID, quizID int) (database.Score, error)
		userID       string
		quizID       string
		status       int
	}{
		{
			name:         "invalid user id value",
			validUser:    auth.ValidUser,
			getUserScore: database.GetUserScore,
			userID:       "testing",
			quizID:       "1",
			status:       http.StatusBadRequest,
		},
		{
			name:         "valid id, invalid quiz id value",
			validUser:    auth.ValidUser,
			getUserScore: database.GetUserScore,
			userID:       "1",
			quizID:       "testing",
			status:       http.StatusBadRequest,
		},
		{
			name: "valid id's, invalid user",
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			getUserScore: database.GetUserScore,
			userID:       "1",
			quizID:       "1",
			status:       http.StatusUnauthorized,
		},
		{
			name: "valid id's, valid user, error on GetUserScore",
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			getUserScore: func(userID, quizID int) (database.Score, error) { return database.Score{}, errors.New("test") },
			userID:       "1",
			quizID:       "1",
			status:       http.StatusInternalServerError,
		},
		{
			name: "valid id's, valid user, no rows on GetUserScore",
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			getUserScore: func(userID, quizID int) (database.Score, error) { return database.Score{}, sql.ErrNoRows },
			userID:       "1",
			quizID:       "1",
			status:       http.StatusNoContent,
		},
		{
			name: "happy path",
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			getUserScore: func(userID, quizID int) (database.Score, error) { return database.Score{}, nil },
			userID:       "1",
			quizID:       "1",
			status:       http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			auth.ValidUser = tc.validUser
			database.GetUserScore = tc.getUserScore
			config.Values = &config.Config{}

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"userId": tc.userID,
				"quizId": tc.quizID,
			})

			writer := httptest.NewRecorder()
			GetScore(writer, request)
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

				var parsed database.Score
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestCreateScore(t *testing.T) {
	savedGetUser := database.GetUser
	savedValidUser := auth.ValidUser
	savedInsertScore := database.InsertScore
	savedConfigValues := config.Values

	defer func() {
		database.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		database.InsertScore = savedInsertScore
		config.Values = savedConfigValues
	}()

	user := database.User{
		Username: "testing",
	}

	tt := []struct {
		name        string
		getUser     func(id int) (database.User, error)
		validUser   func(uv auth.UserValidation) (int, error)
		insertScore func(score database.Score) (int, error)
		body        string
		status      int
	}{
		{
			name:        "invalid body",
			getUser:     database.GetUser,
			validUser:   auth.ValidUser,
			insertScore: database.InsertScore,
			body:        "testing",
			status:      http.StatusBadRequest,
		},
		{
			name:        "valid body, error on GetUser",
			getUser:     func(id int) (database.User, error) { return database.User{}, errors.New("test") },
			validUser:   auth.ValidUser,
			insertScore: database.InsertScore,
			body:        `{"userId":1,"quizId":1,"score":1}`,
			status:      http.StatusInternalServerError,
		},
		{
			name:    "valid body, invalid user",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			insertScore: database.InsertScore,
			body:        `{"userId":1,"quizId":1,"score":1}`,
			status:      http.StatusUnauthorized,
		},
		{
			name:    "valid body, valid user, error on InsertScore",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			insertScore: func(score database.Score) (int, error) { return 0, errors.New("test") },
			body:        `{"userId":1,"quizId":1,"score":1}`,
			status:      http.StatusInternalServerError,
		},
		{
			name:    "happy path",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			insertScore: func(score database.Score) (int, error) { return 1, nil },
			body:        `{"userId":1,"quizId":1,"score":1}`,
			status:      http.StatusCreated,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			database.InsertScore = tc.insertScore
			config.Values = &config.Config{}

			request, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(tc.body)))
			if err != nil {
				t.Fatalf("could not create POST request: %v", err)
			}

			writer := httptest.NewRecorder()
			CreateScore(writer, request)
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

				var parsed database.Score
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestUpdateScore(t *testing.T) {
	savedGetUser := database.GetUser
	savedValidUser := auth.ValidUser
	savedUpdateScore := database.UpdateScore
	savedConfigValues := config.Values

	defer func() {
		database.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		database.UpdateScore = savedUpdateScore
		config.Values = savedConfigValues
	}()

	user := database.User{
		Username: "testing",
	}

	tt := []struct {
		name        string
		getUser     func(id int) (database.User, error)
		validUser   func(uv auth.UserValidation) (int, error)
		updateScore func(score database.Score) error
		id          string
		body        string
		status      int
	}{
		{
			name:        "invalid id",
			getUser:     database.GetUser,
			validUser:   auth.ValidUser,
			updateScore: database.UpdateScore,
			id:          "testing",
			body:        "",
			status:      http.StatusBadRequest,
		},
		{
			name:        "valid id, invalid body",
			getUser:     func(id int) (database.User, error) { return user, nil },
			validUser:   auth.ValidUser,
			updateScore: database.UpdateScore,
			id:          "1",
			body:        "testing",
			status:      http.StatusBadRequest,
		},
		{
			name:        "valid id, valid body, error on GetUser",
			getUser:     func(id int) (database.User, error) { return database.User{}, errors.New("test") },
			validUser:   auth.ValidUser,
			updateScore: database.UpdateScore,
			id:          "1",
			body:        `{"id":1,"userId":1,"quizId":1,"score":1}`,
			status:      http.StatusInternalServerError,
		},
		{
			name:    "valid id, valid body, invalid user",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			updateScore: database.UpdateScore,
			id:          "1",
			body:        `{"id":1,"userId":1,"quizId":1,"score":1}`,
			status:      http.StatusUnauthorized,
		},
		{
			name:    "valid id, valid body, valid user, error on UpdateScore",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			updateScore: func(score database.Score) error { return errors.New("test") },
			id:          "1",
			body:        `{"id":1,"userId":1,"quizId":1,"score":1}`,
			status:      http.StatusInternalServerError,
		},
		{
			name:    "happy path",
			getUser: func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			updateScore: func(score database.Score) error { return nil },
			id:          "1",
			body:        `{"id":1,"userId":1,"quizId":1,"score":1}`,
			status:      http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			database.UpdateScore = tc.updateScore
			config.Values = &config.Config{}

			request, err := http.NewRequest("PUT", "", bytes.NewBuffer([]byte(tc.body)))
			if err != nil {
				t.Fatalf("could not create PUT request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"id": tc.id,
			})

			writer := httptest.NewRecorder()
			UpdateScore(writer, request)
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

				var parsed database.Score
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestDeleteScore(t *testing.T) {
	savedGetScore := database.GetScore
	savedGetUser := database.GetUser
	savedValidUser := auth.ValidUser
	savedDeleteScore := database.DeleteScore
	savedConfigValues := config.Values

	defer func() {
		database.GetScore = savedGetScore
		database.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		database.DeleteScore = savedDeleteScore
		config.Values = savedConfigValues
	}()

	user := database.User{
		Username: "testing",
	}

	tt := []struct {
		name        string
		getScore    func(id int) (database.ScoreDto, error)
		getUser     func(id int) (database.User, error)
		validUser   func(uv auth.UserValidation) (int, error)
		deleteScore func(scoreID int) error
		id          string
		status      int
	}{
		{
			name:        "invalid id",
			getScore:    database.GetScore,
			getUser:     database.GetUser,
			validUser:   auth.ValidUser,
			deleteScore: database.DeleteScore,
			id:          "testing",
			status:      http.StatusBadRequest,
		},
		{
			name:        "valid id, error on GetScore",
			getScore:    func(id int) (database.ScoreDto, error) { return database.ScoreDto{}, errors.New("test") },
			getUser:     func(id int) (database.User, error) { return user, nil },
			validUser:   auth.ValidUser,
			deleteScore: database.DeleteScore,
			id:          "1",
			status:      http.StatusInternalServerError,
		},
		{
			name:        "valid id, error on GetUser",
			getScore:    func(id int) (database.ScoreDto, error) { return database.ScoreDto{}, nil },
			getUser:     func(id int) (database.User, error) { return database.User{}, errors.New("test") },
			validUser:   auth.ValidUser,
			deleteScore: database.DeleteScore,
			id:          "1",
			status:      http.StatusInternalServerError,
		},
		{
			name:     "valid id, score found, invalid user",
			getScore: func(id int) (database.ScoreDto, error) { return database.ScoreDto{}, nil },
			getUser:  func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			deleteScore: database.DeleteScore,
			id:          "1",
			status:      http.StatusUnauthorized,
		},
		{
			name:     "valid id, score found, valid user, error on DeleteScore",
			getScore: func(id int) (database.ScoreDto, error) { return database.ScoreDto{}, nil },
			getUser:  func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			deleteScore: func(scoreID int) error { return errors.New("test") },
			id:          "1",
			status:      http.StatusInternalServerError,
		},
		{
			name:     "happy path",
			getScore: func(id int) (database.ScoreDto, error) { return database.ScoreDto{}, nil },
			getUser:  func(id int) (database.User, error) { return user, nil },
			validUser: func(uv auth.UserValidation) (int, error) {
				return http.StatusOK, nil
			},
			deleteScore: func(scoreID int) error { return nil },
			id:          "1",
			status:      http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetScore = tc.getScore
			database.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			database.DeleteScore = tc.deleteScore
			config.Values = &config.Config{}

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"id": tc.id,
			})

			writer := httptest.NewRecorder()
			DeleteScore(writer, request)
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

				var parsed database.Score
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}
