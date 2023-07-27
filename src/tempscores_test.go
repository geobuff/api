package src

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

func TestGetTempScores(t *testing.T) {
	savedGetTempScore := repo.GetTempScore
	savedDeleteTempScore := repo.DeleteTempScore

	defer func() {
		repo.GetTempScore = savedGetTempScore
		repo.DeleteTempScore = savedDeleteTempScore
	}()

	tt := []struct {
		name            string
		getTempScore    func(id int) (repo.TempScore, error)
		deleteTempScore func(id int) error
		id              string
		status          int
	}{
		{
			name:            "invalid id",
			getTempScore:    repo.GetTempScore,
			deleteTempScore: repo.DeleteTempScore,
			id:              "testing",
			status:          http.StatusBadRequest,
		},
		{
			name:            "valid id, error on GetTempScore",
			getTempScore:    func(id int) (repo.TempScore, error) { return repo.TempScore{}, errors.New("test") },
			deleteTempScore: repo.DeleteTempScore,
			id:              "1",
			status:          http.StatusInternalServerError,
		},
		{
			name:            "valid id, error on DeleteTempScore",
			getTempScore:    func(id int) (repo.TempScore, error) { return repo.TempScore{}, nil },
			deleteTempScore: func(id int) error { return errors.New("test") },
			id:              "1",
			status:          http.StatusInternalServerError,
		},
		{
			name:            "happy path",
			getTempScore:    func(id int) (repo.TempScore, error) { return repo.TempScore{}, nil },
			deleteTempScore: func(id int) error { return nil },
			id:              "1",
			status:          http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetTempScore = tc.getTempScore
			repo.DeleteTempScore = tc.deleteTempScore

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"id": tc.id,
			})

			writer := httptest.NewRecorder()
			GetTempScore(writer, request)
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

				var parsed repo.TempScore
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestCreateTempScore(t *testing.T) {
	savedInsertTempScore := repo.InsertTempScore

	defer func() {
		repo.InsertTempScore = savedInsertTempScore
	}()

	tt := []struct {
		name            string
		insertTempScore func(score repo.TempScore) (int, error)
		body            string
		status          int
	}{
		{
			name:            "invalid body",
			insertTempScore: repo.InsertTempScore,
			body:            "testing",
			status:          http.StatusBadRequest,
		},
		{
			name:            "valid body, error on InsertTempScore",
			insertTempScore: func(score repo.TempScore) (int, error) { return 0, errors.New("test") },
			body:            `{"score":1,"time":2,"results":["new zealand"],"recents": ["new zealand"]}`,
			status:          http.StatusInternalServerError,
		},
		{
			name:            "happy path",
			insertTempScore: func(score repo.TempScore) (int, error) { return 1, nil },
			body:            `{"score":1,"time":2,"results":["new zealand"],"recents": ["new zealand"]}`,
			status:          http.StatusCreated,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.InsertTempScore = tc.insertTempScore

			request, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(tc.body)))
			if err != nil {
				t.Fatalf("could not create POST request: %v", err)
			}

			writer := httptest.NewRecorder()
			CreateTempScore(writer, request)
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

				var parsed int
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}
