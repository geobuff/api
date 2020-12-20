package world

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetCountries(t *testing.T) {
	request, err := http.NewRequest("GET", "localhost:8080/api/world/countries", nil)
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

	var parsed []string
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		t.Fatalf("could not unmarshal response body: %v", err)
	}

	if !reflect.DeepEqual(parsed, countries) {
		t.Fatalf("response body does not match expected countries")
	}
}

func TestGetAlternativeNamings(t *testing.T) {
	request, err := http.NewRequest("GET", "localhost:8080/api/world/countries/alternatives", nil)
	if err != nil {
		t.Fatalf("could not create GET request: %v", err)
	}

	writer := httptest.NewRecorder()
	GetAlternativeNamings(writer, request)
	result := writer.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", result.StatusCode)
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}

	var parsed []string
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		t.Fatalf("could not unmarshal response body: %v", err)
	}

	if !reflect.DeepEqual(parsed, alternativeNamings) {
		t.Fatalf("response body does not match expected alternative namings")
	}
}

func TestGetPrefixes(t *testing.T) {
	request, err := http.NewRequest("GET", "localhost:8080/api/world/countries/prefixes", nil)
	if err != nil {
		t.Fatalf("could not create GET request: %v", err)
	}

	writer := httptest.NewRecorder()
	GetPrefixes(writer, request)
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

	if !reflect.DeepEqual(parsed, prefixes) {
		t.Fatalf("response body does not match expected prefixes")
	}
}

func TestGetCountriesMap(t *testing.T) {
	request, err := http.NewRequest("GET", "localhost:8080/api/world/countries/map", nil)
	if err != nil {
		t.Fatalf("could not create GET request: %v", err)
	}

	writer := httptest.NewRecorder()
	GetCountriesMap(writer, request)
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

	if !reflect.DeepEqual(parsed, countriesMap) {
		t.Fatalf("response body does not match expected countries map")
	}
}
