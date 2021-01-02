package isocodes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mapping "github.com/geobuff/location-mapping/isocodes"
)

func TestGetCodes(t *testing.T) {
	request, err := http.NewRequest("GET", "localhost:8080/api/isocodes", nil)
	if err != nil {
		t.Fatalf("could not create GET request: %v", err)
	}

	writer := httptest.NewRecorder()
	GetCodes(writer, request)
	result := writer.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", result.StatusCode)
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}

	var parsed map[string]string
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		t.Fatalf("could not unmarshal response body: %v", err)
	}

	if !reflect.DeepEqual(parsed, mapping.Codes) {
		t.Fatalf("response body does not match expected codes")
	}
}
