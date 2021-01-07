package world

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/geobuff/mapping"
)

func TestGetCountries(t *testing.T) {
	request, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatalf("could not create GET request: %v", err)
	}

	writer := httptest.NewRecorder()
	GetCountries(writer, request)
	result := writer.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", result.StatusCode)
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}

	var parsed map[string][]mapping.Country
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		t.Fatalf("could not unmarshal response body: %v", err)
	}

	if !reflect.DeepEqual(parsed, mapping.Countries) {
		t.Fatalf("response body does not match expected map")
	}
}

func TestGetCountriesQueryParam(t *testing.T) {
	continent := "asia"
	request, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/api/world/countries?continent=%v", continent), nil)
	if err != nil {
		t.Fatalf("could not create GET request: %v", err)
	}

	writer := httptest.NewRecorder()
	GetCountries(writer, request)
	result := writer.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", result.StatusCode)
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}

	var parsed []mapping.Country
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		t.Fatalf("could not unmarshal response body: %v", err)
	}

	if !reflect.DeepEqual(parsed, mapping.Countries[continent]) {
		t.Fatalf("response body does not match expected list of countries")
	}
}
