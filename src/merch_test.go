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

func TestGetMerch(t *testing.T) {
	savedGetMerch := repo.GetMerch

	defer func() {
		repo.GetMerch = savedGetMerch
	}()

	tt := []struct {
		name     string
		getMerch func() ([]repo.MerchDto, error)
		status   int
	}{
		{
			name:     "error on GetMerch",
			getMerch: func() ([]repo.MerchDto, error) { return nil, errors.New("test") },
			status:   http.StatusInternalServerError,
		},
		{
			name:     "happy path",
			getMerch: func() ([]repo.MerchDto, error) { return []repo.MerchDto{}, nil },
			status:   http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo.GetMerch = tc.getMerch

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			writer := httptest.NewRecorder()
			GetMerch(writer, request)
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

				var parsed []repo.Merch
				err = json.Unmarshal(body, &parsed)
				if err != nil {
					t.Errorf("could not unmarshal response body: %v", err)
				}
			}
		})
	}
}
