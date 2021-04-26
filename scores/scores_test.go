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

	"github.com/geobuff/api/auth"
	"github.com/geobuff/api/config"
	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

func TestGetScores(t *testing.T) {
	savedGetScores := repo.GetScores

	defer func() {
		repo.GetScores = savedGetScores
	}()

	tt := []struct {
		name      string
		getScores func(userID int) ([]repo.ScoreDto, error)
		userID    string
		status    int
	}{
		{
			name:      "invalid id value",
			getScores: repo.GetScores,
			userID:    "testing",
			status:    http.StatusBadRequest,
		},
		{
			name:      "valid id, error on GetScores",
			getScores: func(userID int) ([]repo.ScoreDto, error) { return nil, errors.New("test") },
			userID:    "1",
			status:    http.StatusInternalServerError,
		},
		{
			name:      "happy path",
			getScores: func(userID int) ([]repo.ScoreDto, error) { return []repo.ScoreDto{}, nil },
			userID:    "1",
			status:    http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetScores = tc.getScores

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

				var parsed []repo.Score
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestGetScore(t *testing.T) {
	savedGetUserScore := repo.GetUserScore

	defer func() {
		repo.GetUserScore = savedGetUserScore
	}()

	tt := []struct {
		name         string
		getUserScore func(userID, quizID int) (repo.Score, error)
		userID       string
		quizID       string
		status       int
	}{
		{
			name:         "invalid user id value",
			getUserScore: repo.GetUserScore,
			userID:       "testing",
			quizID:       "1",
			status:       http.StatusBadRequest,
		},
		{
			name:         "valid id, invalid quiz id value",
			getUserScore: repo.GetUserScore,
			userID:       "1",
			quizID:       "testing",
			status:       http.StatusBadRequest,
		},
		{
			name:         "valid id's, error on GetUserScore",
			getUserScore: func(userID, quizID int) (repo.Score, error) { return repo.Score{}, errors.New("test") },
			userID:       "1",
			quizID:       "1",
			status:       http.StatusInternalServerError,
		},
		{
			name:         "valid id's, no rows on GetUserScore",
			getUserScore: func(userID, quizID int) (repo.Score, error) { return repo.Score{}, sql.ErrNoRows },
			userID:       "1",
			quizID:       "1",
			status:       http.StatusNoContent,
		},
		{
			name:         "happy path",
			getUserScore: func(userID, quizID int) (repo.Score, error) { return repo.Score{}, nil },
			userID:       "1",
			quizID:       "1",
			status:       http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetUserScore = tc.getUserScore

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

				var parsed repo.Score
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestCreateScore(t *testing.T) {
	savedGetUser := repo.GetUser
	savedValidUser := auth.ValidUser
	savedInsertScore := repo.InsertScore
	savedConfigValues := config.Values

	defer func() {
		repo.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		repo.InsertScore = savedInsertScore
		config.Values = savedConfigValues
	}()

	user := repo.UserDto{
		Username: "testing",
	}

	tt := []struct {
		name        string
		getUser     func(id int) (repo.UserDto, error)
		validUser   func(request *http.Request, id int) (int, error)
		insertScore func(score repo.Score) (int, error)
		body        string
		status      int
	}{
		{
			name:        "invalid body",
			getUser:     repo.GetUser,
			validUser:   auth.ValidUser,
			insertScore: repo.InsertScore,
			body:        "testing",
			status:      http.StatusBadRequest,
		},
		{
			name:        "valid body, error on GetUser",
			getUser:     func(id int) (repo.UserDto, error) { return repo.UserDto{}, errors.New("test") },
			validUser:   auth.ValidUser,
			insertScore: repo.InsertScore,
			body:        `{"userId":1,"quizId":1,"score":1}`,
			status:      http.StatusInternalServerError,
		},
		{
			name:    "valid body, invalid user",
			getUser: func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			insertScore: repo.InsertScore,
			body:        `{"userId":1,"quizId":1,"score":1}`,
			status:      http.StatusUnauthorized,
		},
		{
			name:    "valid body, valid user, error on InsertScore",
			getUser: func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			insertScore: func(score repo.Score) (int, error) { return 0, errors.New("test") },
			body:        `{"userId":1,"quizId":1,"score":1}`,
			status:      http.StatusInternalServerError,
		},
		{
			name:    "happy path",
			getUser: func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			insertScore: func(score repo.Score) (int, error) { return 1, nil },
			body:        `{"userId":1,"quizId":1,"score":1}`,
			status:      http.StatusCreated,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			repo.InsertScore = tc.insertScore
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

				var parsed repo.Score
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestUpdateScore(t *testing.T) {
	savedGetUser := repo.GetUser
	savedValidUser := auth.ValidUser
	savedUpdateScore := repo.UpdateScore
	savedConfigValues := config.Values

	defer func() {
		repo.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		repo.UpdateScore = savedUpdateScore
		config.Values = savedConfigValues
	}()

	user := repo.UserDto{
		Username: "testing",
	}

	tt := []struct {
		name        string
		getUser     func(id int) (repo.UserDto, error)
		validUser   func(request *http.Request, id int) (int, error)
		updateScore func(score repo.Score) error
		id          string
		body        string
		status      int
	}{
		{
			name:        "invalid id",
			getUser:     repo.GetUser,
			validUser:   auth.ValidUser,
			updateScore: repo.UpdateScore,
			id:          "testing",
			body:        "",
			status:      http.StatusBadRequest,
		},
		{
			name:        "valid id, invalid body",
			getUser:     func(id int) (repo.UserDto, error) { return user, nil },
			validUser:   auth.ValidUser,
			updateScore: repo.UpdateScore,
			id:          "1",
			body:        "testing",
			status:      http.StatusBadRequest,
		},
		{
			name:        "valid id, valid body, error on GetUser",
			getUser:     func(id int) (repo.UserDto, error) { return repo.UserDto{}, errors.New("test") },
			validUser:   auth.ValidUser,
			updateScore: repo.UpdateScore,
			id:          "1",
			body:        `{"id":1,"userId":1,"quizId":1,"score":1}`,
			status:      http.StatusInternalServerError,
		},
		{
			name:    "valid id, valid body, invalid user",
			getUser: func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			updateScore: repo.UpdateScore,
			id:          "1",
			body:        `{"id":1,"userId":1,"quizId":1,"score":1}`,
			status:      http.StatusUnauthorized,
		},
		{
			name:    "valid id, valid body, valid user, error on UpdateScore",
			getUser: func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			updateScore: func(score repo.Score) error { return errors.New("test") },
			id:          "1",
			body:        `{"id":1,"userId":1,"quizId":1,"score":1}`,
			status:      http.StatusInternalServerError,
		},
		{
			name:    "happy path",
			getUser: func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			updateScore: func(score repo.Score) error { return nil },
			id:          "1",
			body:        `{"id":1,"userId":1,"quizId":1,"score":1}`,
			status:      http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			repo.UpdateScore = tc.updateScore
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

				var parsed repo.Score
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestDeleteScore(t *testing.T) {
	savedGetScore := repo.GetScore
	savedGetUser := repo.GetUser
	savedValidUser := auth.ValidUser
	savedDeleteScore := repo.DeleteScore
	savedConfigValues := config.Values

	defer func() {
		repo.GetScore = savedGetScore
		repo.GetUser = savedGetUser
		auth.ValidUser = savedValidUser
		repo.DeleteScore = savedDeleteScore
		config.Values = savedConfigValues
	}()

	user := repo.UserDto{
		Username: "testing",
	}

	tt := []struct {
		name        string
		getScore    func(id int) (repo.ScoreDto, error)
		getUser     func(id int) (repo.UserDto, error)
		validUser   func(request *http.Request, id int) (int, error)
		deleteScore func(scoreID int) error
		id          string
		status      int
	}{
		{
			name:        "invalid id",
			getScore:    repo.GetScore,
			getUser:     repo.GetUser,
			validUser:   auth.ValidUser,
			deleteScore: repo.DeleteScore,
			id:          "testing",
			status:      http.StatusBadRequest,
		},
		{
			name:        "valid id, error on GetScore",
			getScore:    func(id int) (repo.ScoreDto, error) { return repo.ScoreDto{}, errors.New("test") },
			getUser:     func(id int) (repo.UserDto, error) { return user, nil },
			validUser:   auth.ValidUser,
			deleteScore: repo.DeleteScore,
			id:          "1",
			status:      http.StatusInternalServerError,
		},
		{
			name:        "valid id, error on GetUser",
			getScore:    func(id int) (repo.ScoreDto, error) { return repo.ScoreDto{}, nil },
			getUser:     func(id int) (repo.UserDto, error) { return repo.UserDto{}, errors.New("test") },
			validUser:   auth.ValidUser,
			deleteScore: repo.DeleteScore,
			id:          "1",
			status:      http.StatusInternalServerError,
		},
		{
			name:     "valid id, score found, invalid user",
			getScore: func(id int) (repo.ScoreDto, error) { return repo.ScoreDto{}, nil },
			getUser:  func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			deleteScore: repo.DeleteScore,
			id:          "1",
			status:      http.StatusUnauthorized,
		},
		{
			name:     "valid id, score found, valid user, error on DeleteScore",
			getScore: func(id int) (repo.ScoreDto, error) { return repo.ScoreDto{}, nil },
			getUser:  func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			deleteScore: func(scoreID int) error { return errors.New("test") },
			id:          "1",
			status:      http.StatusInternalServerError,
		},
		{
			name:     "happy path",
			getScore: func(id int) (repo.ScoreDto, error) { return repo.ScoreDto{}, nil },
			getUser:  func(id int) (repo.UserDto, error) { return user, nil },
			validUser: func(request *http.Request, id int) (int, error) {
				return http.StatusOK, nil
			},
			deleteScore: func(scoreID int) error { return nil },
			id:          "1",
			status:      http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetScore = tc.getScore
			repo.GetUser = tc.getUser
			auth.ValidUser = tc.validUser
			repo.DeleteScore = tc.deleteScore
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

				var parsed repo.Score
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}
