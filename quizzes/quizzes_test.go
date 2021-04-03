package quizzes

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

func TestGetQuizzes(t *testing.T) {
	savedGetQuizzes := repo.GetQuizzes

	defer func() {
		repo.GetQuizzes = savedGetQuizzes
	}()

	tt := []struct {
		Name       string
		getQuizzes func(filter string) ([]repo.Quiz, error)
		status     int
	}{
		{
			Name:       "error on getQuizzes",
			getQuizzes: func(filter string) ([]repo.Quiz, error) { return nil, errors.New("test") },
			status:     http.StatusInternalServerError,
		},
		{
			Name:       "happy path",
			getQuizzes: func(filter string) ([]repo.Quiz, error) { return []repo.Quiz{}, nil },
			status:     http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			repo.GetQuizzes = tc.getQuizzes

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			writer := httptest.NewRecorder()
			GetQuizzes(writer, request)
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

				var parsed []repo.Quiz
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}

func TestGetQuiz(t *testing.T) {
	savedGetQuiz := repo.GetQuiz

	defer func() {
		repo.GetQuiz = savedGetQuiz
	}()

	tt := []struct {
		Name    string
		getQuiz func(id int) (repo.Quiz, error)
		id      string
		status  int
	}{
		{
			Name:    "invalid id",
			id:      "testing",
			getQuiz: repo.GetQuiz,
			status:  http.StatusBadRequest,
		},
		{
			Name:    "valid id, error on getQuizzes",
			id:      "1",
			getQuiz: func(id int) (repo.Quiz, error) { return repo.Quiz{}, errors.New("test") },
			status:  http.StatusInternalServerError,
		},
		{
			Name:    "happy path",
			id:      "1",
			getQuiz: func(id int) (repo.Quiz, error) { return repo.Quiz{}, nil },
			status:  http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			repo.GetQuiz = tc.getQuiz

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"id": tc.id,
			})

			writer := httptest.NewRecorder()
			GetQuiz(writer, request)
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

				var parsed repo.Quiz
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}
