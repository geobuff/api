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

func TestGetQuizzes(t *testing.T) {
	savedGetQuizzes := repo.GetQuizzes
	savedGetFirstQuizID := repo.GetFirstQuizID

	defer func() {
		repo.GetQuizzes = savedGetQuizzes
		repo.GetFirstQuizID = savedGetFirstQuizID
	}()

	tt := []struct {
		Name           string
		getQuizzes     func(filter repo.QuizzesFilterDto) ([]repo.Quiz, error)
		getFirstQuizID func(offset int) (int, error)
		body           string
		status         int
	}{
		{
			Name:           "invalid body",
			getQuizzes:     repo.GetQuizzes,
			getFirstQuizID: repo.GetFirstQuizID,
			body:           "testing",
			status:         http.StatusBadRequest,
		},
		{
			Name:           "error on getQuizzes",
			getQuizzes:     func(filter repo.QuizzesFilterDto) ([]repo.Quiz, error) { return nil, errors.New("test") },
			getFirstQuizID: repo.GetFirstQuizID,
			body:           `{"filter":"","page":0,"limit":10}`,
			status:         http.StatusInternalServerError,
		},
		{
			Name:           "error on getFirstQuizID",
			getQuizzes:     func(filter repo.QuizzesFilterDto) ([]repo.Quiz, error) { return []repo.Quiz{}, nil },
			getFirstQuizID: func(offset int) (int, error) { return 0, errors.New("test") },
			body:           `{"filter":"","page":0,"limit":10}`,
			status:         http.StatusInternalServerError,
		},
		{
			Name:           "happy path",
			getQuizzes:     func(filter repo.QuizzesFilterDto) ([]repo.Quiz, error) { return []repo.Quiz{}, nil },
			getFirstQuizID: func(offset int) (int, error) { return 1, nil },
			body:           `{"filter":"","page":0,"limit":10}`,
			status:         http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			repo.GetQuizzes = tc.getQuizzes
			repo.GetFirstQuizID = tc.getFirstQuizID

			request, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(tc.body)))
			if err != nil {
				t.Fatalf("could not create POST request: %v", err)
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

				var parsed QuizPageDto
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
