package plays

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
	savedGetAllPlays := repo.GetAllPlays

	defer func() {
		repo.GetAllPlays = savedGetAllPlays
	}()

	tt := []struct {
		name        string
		getAllPlays func() (int, error)
		status      int
	}{
		{
			name:        "error on GetMerch",
			getAllPlays: func() (int, error) { return 0, errors.New("test") },
			status:      http.StatusInternalServerError,
		},
		{
			name:        "happy path",
			getAllPlays: func() (int, error) { return 1, nil },
			status:      http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetAllPlays = tc.getAllPlays

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			writer := httptest.NewRecorder()
			GetAllPlays(writer, request)
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

func TestGetPlays(t *testing.T) {
	savedGetPlayCount := repo.GetPlayCount

	defer func() {
		repo.GetPlayCount = savedGetPlayCount
	}()

	tt := []struct {
		name         string
		getPlayCount func(quizID int) (int, error)
		quizID       string
		status       int
	}{
		{
			name:         "invalid quizID",
			getPlayCount: repo.GetPlayCount,
			quizID:       "test",
			status:       http.StatusBadRequest,
		},
		{
			name:         "valid quizId, error on GetPlayCount",
			getPlayCount: func(quizID int) (int, error) { return 0, errors.New("test") },
			quizID:       "1",
			status:       http.StatusInternalServerError,
		},
		{
			name:         "happy path",
			getPlayCount: func(quizID int) (int, error) { return 1, nil },
			quizID:       "1",
			status:       http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetPlayCount = tc.getPlayCount

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"quizId": tc.quizID,
			})

			writer := httptest.NewRecorder()
			GetPlays(writer, request)
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
	savedIncrementPlayCount := repo.IncrementPlayCount

	defer func() {
		repo.IncrementPlayCount = savedIncrementPlayCount
	}()

	tt := []struct {
		name               string
		incrementPlayCount func(quizID int) error
		quizID             string
		status             int
	}{
		{
			name:               "invalid quizID",
			incrementPlayCount: repo.IncrementPlayCount,
			quizID:             "test",
			status:             http.StatusBadRequest,
		},
		{
			name:               "valid quizId, error on IncrementPlayCount",
			incrementPlayCount: func(quizID int) error { return errors.New("test") },
			quizID:             "1",
			status:             http.StatusInternalServerError,
		},
		{
			name:               "happy path",
			incrementPlayCount: func(quizID int) error { return nil },
			quizID:             "1",
			status:             http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.IncrementPlayCount = tc.incrementPlayCount

			request, err := http.NewRequest("PUT", "", nil)
			if err != nil {
				t.Fatalf("could not create PUT request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"quizId": tc.quizID,
			})

			writer := httptest.NewRecorder()
			IncrementPlays(writer, request)
			result := writer.Result()
			defer result.Body.Close()

			if result.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, result.StatusCode)
			}
		})
	}
}
