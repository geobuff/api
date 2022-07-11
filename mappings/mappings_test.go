package mappings

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

func TestGetMapping(t *testing.T) {
	key := "countries"
	request, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatalf("could not create GET request: %v", err)
	}

	request = mux.SetURLVars(request, map[string]string{
		"key": key,
	})

	writer := httptest.NewRecorder()
	GetMapping(writer, request)
	result := writer.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", result.StatusCode)
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}

	var parsed []repo.Mapping
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		t.Errorf("could not unmarshal response body: %v", err)
	}

	if !reflect.DeepEqual(parsed, repo.Mappings[key]) {
		t.Fatalf("response body does not match expected list of mappings")
	}
}
