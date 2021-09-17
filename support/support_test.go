package support

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geobuff/api/email"
	"github.com/sendgrid/rest"
)

func TestSendSupportRequest(t *testing.T) {
	savedSendSupportRequest := email.SendSupportRequest

	defer func() {
		email.SendSupportRequest = savedSendSupportRequest
	}()

	tt := []struct {
		name               string
		sendSupportRequest func(request email.SupportRequest) (*rest.Response, error)
		body               string
		status             int
	}{
		{
			name:               "invalid body",
			sendSupportRequest: email.SendSupportRequest,
			body:               "testing",
			status:             http.StatusBadRequest,
		},
		{
			name:               "valid body, error on SendSupportRequest",
			sendSupportRequest: func(request email.SupportRequest) (*rest.Response, error) { return nil, errors.New("test") },
			body:               `{"from": "testing@gmail.com","subject": "Testing", "message": "testing testing"}`,
			status:             http.StatusInternalServerError,
		},
		{
			name:               "happy path",
			sendSupportRequest: func(request email.SupportRequest) (*rest.Response, error) { return nil, nil },
			body:               `{"from": "testing@gmail.com","subject": "Testing", "message": "testing testing"}`,
			status:             http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			email.SendSupportRequest = tc.sendSupportRequest

			request, err := http.NewRequest("POST", "", bytes.NewBuffer([]byte(tc.body)))
			if err != nil {
				t.Fatalf("could not create POST request: %v", err)
			}

			writer := httptest.NewRecorder()
			SendSupportRequest(writer, request)
			result := writer.Result()
			defer result.Body.Close()

			if result.StatusCode != tc.status {
				t.Errorf("expected status %v; got %v", tc.status, result.StatusCode)
			}
		})
	}
}
