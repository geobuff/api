package badges

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

func TestGetBadges(t *testing.T) {
	savedGetUserBadges := repo.GetUserBadges

	defer func() {
		repo.GetUserBadges = savedGetUserBadges
	}()

	tt := []struct {
		name          string
		getUserBadges func(userID int) ([]repo.BadgeDto, error)
		userId        string
		status        int
	}{
		{
			name:          "invalid userId",
			getUserBadges: repo.GetUserBadges,
			userId:        "test",
			status:        http.StatusBadRequest,
		},
		{
			name:          "error on GetBadges",
			getUserBadges: func(userID int) ([]repo.BadgeDto, error) { return nil, errors.New("test") },
			userId:        "1",
			status:        http.StatusInternalServerError,
		},
		{
			name:          "happy path",
			getUserBadges: func(userID int) ([]repo.BadgeDto, error) { return []repo.BadgeDto{}, nil },
			userId:        "1",
			status:        http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetUserBadges = tc.getUserBadges

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			request = mux.SetURLVars(request, map[string]string{
				"userId": tc.userId,
			})

			writer := httptest.NewRecorder()
			GetUserBadges(writer, request)
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

				var parsed []repo.Badge
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}
