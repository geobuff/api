package scores

import (
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
