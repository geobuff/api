package src

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geobuff/api/repo"
)

func TestGetAvatars(t *testing.T) {
	savedGetAvatars := repo.GetAvatars

	defer func() {
		repo.GetAvatars = savedGetAvatars
	}()

	tt := []struct {
		name       string
		getAvatars func() ([]repo.AvatarDto, error)
		status     int
	}{
		{
			name:       "error on GetAvatars",
			getAvatars: func() ([]repo.AvatarDto, error) { return nil, errors.New("test") },
			status:     http.StatusInternalServerError,
		},
		{
			name:       "happy path",
			getAvatars: func() ([]repo.AvatarDto, error) { return []repo.AvatarDto{}, nil },
			status:     http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetAvatars = tc.getAvatars

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			writer := httptest.NewRecorder()
			s := getMockServer()
			s.getAvatars(writer, request)
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

				var parsed []repo.Avatar
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}
