package us

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geobuff/mapping"
)

func TestGetStates(t *testing.T) {
	request, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatalf("could not create GET request: %v", err)
	}

	writer := httptest.NewRecorder()
	GetStates(writer, request)
	result := writer.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", result.StatusCode)
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}

	var parsed []mapping.Mapping
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		t.Errorf("could not unmarshal response body: %v", err)
	}
}
