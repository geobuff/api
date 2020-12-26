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

	"github.com/geobuff/geobuff-api/auth"
	"github.com/geobuff/geobuff-api/database"
	"github.com/gorilla/mux"
)

func TestGetScores(t *testing.T) {
	savedValidUser := auth.ValidUser
	savedGetScores := database.GetScores

	defer func() {
		auth.ValidUser = savedValidUser
		database.GetScores = savedGetScores
	}()

	tt := []struct {
		name      string
		validUser func(request *http.Request, userID int, permission string) (int, error)
		getScores func(userID int) ([]database.Score, error)
		id        string
		status    int
	}{
		{
			name:      "invalid id value",
			validUser: auth.ValidUser,
			getScores: database.GetScores,
			id:        "testing",
			status:    http.StatusBadRequest,
		},
		{
			name: "valid id, invalid user",
			validUser: func(request *http.Request, userID int, permission string) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			getScores: database.GetScores,
			id:        "1",
			status:    http.StatusUnauthorized,
		},
		{
			name: "valid id, valid user, error on GetScores",
			validUser: func(request *http.Request, userID int, permission string) (int, error) {
				return http.StatusOK, nil
			},
			getScores: func(userID int) ([]database.Score, error) { return nil, errors.New("test") },
			id:        "1",
			status:    http.StatusInternalServerError,
		},
		{
			name: "happy path",
			validUser: func(request *http.Request, userID int, permission string) (int, error) {
				return http.StatusOK, nil
			},
			getScores: func(userID int) ([]database.Score, error) { return []database.Score{}, nil },
			id:        "1",
			status:    http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			auth.ValidUser = tc.validUser
			database.GetScores = tc.getScores

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"id": tc.id,
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

func TestCreateScore(t *testing.T) {
	savedValidUser := auth.ValidUser
	savedInsertScore := database.InsertScore

	defer func() {
		auth.ValidUser = savedValidUser
		database.InsertScore = savedInsertScore
	}()

	tt := []struct {
		name        string
		validUser   func(request *http.Request, userID int, permission string) (int, error)
		insertScore func(score database.Score) (int, error)
		body        string
		status      int
	}{
		{
			name:        "invalid body",
			validUser:   auth.ValidUser,
			insertScore: database.InsertScore,
			body:        "testing",
			status:      http.StatusBadRequest,
		},
		{
			name: "valid body, invalid user",
			validUser: func(request *http.Request, userID int, permission string) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			insertScore: database.InsertScore,
			body:        `{"id":1,"userId":1,"quizId":1,"score":1}`,
			status:      http.StatusUnauthorized,
		},
		{
			name: "valid body, valid user, error on InsertScore",
			validUser: func(request *http.Request, userID int, permission string) (int, error) {
				return http.StatusOK, nil
			},
			insertScore: func(score database.Score) (int, error) { return 0, errors.New("test") },
			body:        `{"id":1,"userId":1,"quizId":1,"score":1}`,
			status:      http.StatusInternalServerError,
		},
		{
			name: "happy path",
			validUser: func(request *http.Request, userID int, permission string) (int, error) {
				return http.StatusOK, nil
			},
			insertScore: func(score database.Score) (int, error) { return 1, nil },
			body:        `{"id":1,"userId":1,"quizId":1,"score":1}`,
			status:      http.StatusCreated,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			auth.ValidUser = tc.validUser
			database.InsertScore = tc.insertScore

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

func TestDeleteScore(t *testing.T) {
	savedGetScore := database.GetScore
	savedValidUser := auth.ValidUser
	savedDeleteScore := database.DeleteScore

	defer func() {
		database.GetScore = savedGetScore
		auth.ValidUser = savedValidUser
		database.DeleteScore = savedDeleteScore
	}()

	tt := []struct {
		name        string
		getScore    func(id int) (database.Score, error)
		validUser   func(request *http.Request, userID int, permission string) (int, error)
		deleteScore func(scoreID int) error
		id          string
		status      int
	}{
		{
			name:        "invalid id",
			getScore:    database.GetScore,
			validUser:   auth.ValidUser,
			deleteScore: database.DeleteScore,
			id:          "testing",
			status:      http.StatusBadRequest,
		},
		{
			name:        "valid id, score not found",
			getScore:    func(id int) (database.Score, error) { return database.Score{}, sql.ErrNoRows },
			validUser:   auth.ValidUser,
			deleteScore: database.DeleteScore,
			id:          "1",
			status:      http.StatusNotFound,
		},
		{
			name:        "valid id, unknown error on GetScore",
			getScore:    func(id int) (database.Score, error) { return database.Score{}, errors.New("test") },
			validUser:   auth.ValidUser,
			deleteScore: database.DeleteScore,
			id:          "1",
			status:      http.StatusInternalServerError,
		},
		{
			name:     "valid id, score found, invalid user",
			getScore: func(id int) (database.Score, error) { return database.Score{}, nil },
			validUser: func(request *http.Request, userID int, permission string) (int, error) {
				return http.StatusUnauthorized, errors.New("test")
			},
			deleteScore: database.DeleteScore,
			id:          "1",
			status:      http.StatusUnauthorized,
		},
		{
			name:     "valid id, score found, valid user, error on DeleteScore",
			getScore: func(id int) (database.Score, error) { return database.Score{}, nil },
			validUser: func(request *http.Request, userID int, permission string) (int, error) {
				return http.StatusOK, nil
			},
			deleteScore: func(scoreID int) error { return errors.New("test") },
			id:          "1",
			status:      http.StatusInternalServerError,
		},
		{
			name:     "happy path",
			getScore: func(id int) (database.Score, error) { return database.Score{}, nil },
			validUser: func(request *http.Request, userID int, permission string) (int, error) {
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
			auth.ValidUser = tc.validUser
			database.DeleteScore = tc.deleteScore

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
