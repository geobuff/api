package quizplays

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

func TestGetAllPlays(t *testing.T) {
	savedGetAllQuizPlays := repo.GetAllQuizPlays

	defer func() {
		repo.GetAllQuizPlays = savedGetAllQuizPlays
	}()

	tt := []struct {
		name            string
		getAllQuizPlays func() (int, error)
		status          int
		expected        int
	}{
		{
			name:            "scan error on GetAllQuizPlays",
			getAllQuizPlays: func() (int, error) { return 0, errors.New("sql: Scan error on column index 0") },
			status:          http.StatusOK,
			expected:        0,
		},
		{
			name:            "other error on GetAllQuizPlays",
			getAllQuizPlays: func() (int, error) { return 0, errors.New("test") },
			status:          http.StatusInternalServerError,
			expected:        0,
		},
		{
			name:            "happy path",
			getAllQuizPlays: func() (int, error) { return 1, nil },
			status:          http.StatusOK,
			expected:        1,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetAllQuizPlays = tc.getAllQuizPlays

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			writer := httptest.NewRecorder()
			GetAllQuizPlays(writer, request)
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

				if parsed != tc.expected {
					t.Errorf("expected %d; got %d", tc.expected, parsed)
				}
			}
		})
	}
}

func TestGetPlays(t *testing.T) {
	savedGetPlayCount := repo.GetQuizPlayCount

	defer func() {
		repo.GetQuizPlayCount = savedGetPlayCount
	}()

	tt := []struct {
		name             string
		getQuizPlayCount func(quizID int) (int, error)
		quizID           string
		status           int
	}{
		{
			name:             "invalid quizID",
			getQuizPlayCount: repo.GetQuizPlayCount,
			quizID:           "test",
			status:           http.StatusBadRequest,
		},
		{
			name:             "valid quizId, error on GetPlayCount",
			getQuizPlayCount: func(quizID int) (int, error) { return 0, errors.New("test") },
			quizID:           "1",
			status:           http.StatusInternalServerError,
		},
		{
			name:             "happy path",
			getQuizPlayCount: func(quizID int) (int, error) { return 1, nil },
			quizID:           "1",
			status:           http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetQuizPlayCount = tc.getQuizPlayCount

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"quizId": tc.quizID,
			})

			writer := httptest.NewRecorder()
			GetQuizPlays(writer, request)
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

func TestIncrementPlays(t *testing.T) {
	savedIncrementQuizPlayCount := repo.IncrementQuizPlayCount

	defer func() {
		repo.IncrementQuizPlayCount = savedIncrementQuizPlayCount
	}()

	tt := []struct {
		name                   string
		incrementQuizPlayCount func(quizID int) error
		quizID                 string
		status                 int
	}{
		{
			name:                   "invalid quizID",
			incrementQuizPlayCount: repo.IncrementQuizPlayCount,
			quizID:                 "test",
			status:                 http.StatusBadRequest,
		},
		{
			name:                   "valid quizId, error on IncrementPlayCount",
			incrementQuizPlayCount: func(quizID int) error { return errors.New("test") },
			quizID:                 "1",
			status:                 http.StatusInternalServerError,
		},
		{
			name:                   "happy path",
			incrementQuizPlayCount: func(quizID int) error { return nil },
			quizID:                 "1",
			status:                 http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.IncrementQuizPlayCount = tc.incrementQuizPlayCount

			request, err := http.NewRequest("PUT", "", nil)
			if err != nil {
				t.Fatalf("could not create PUT request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"quizId": tc.quizID,
			})

			writer := httptest.NewRecorder()
			IncrementQuizPlays(writer, request)
			result := writer.Result()
			defer result.Body.Close()

			if result.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, result.StatusCode)
			}
		})
	}
}
