package badges

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geobuff/api/database"
)

func TestGetBadges(t *testing.T) {
	savedGetBadges := database.GetBadges

	defer func() {
		database.GetBadges = savedGetBadges
	}()

	tt := []struct {
		name      string
		getBadges func() ([]database.Badge, error)
		status    int
	}{
		{
			name:      "error on GetBadges",
			getBadges: func() ([]database.Badge, error) { return nil, errors.New("Test") },
			status:    http.StatusInternalServerError,
		},
		{
			name:      "happy path",
			getBadges: func() ([]database.Badge, error) { return []database.Badge{}, nil },
			status:    http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			database.GetBadges = tc.getBadges

			request, err := http.NewRequest("GET", "", nil)
			if err != nil {
				t.Fatalf("could not create GET request: %v", err)
			}

			writer := httptest.NewRecorder()
			GetBadges(writer, request)
			result := writer.Result()
			defer result.Body.Close()

			if result.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, result.StatusCode)
			}
		})
	}
}
