package world

import (
	"encoding/json"
	"net/http"

	mapping "github.com/geobuff/location-mapping/world"
)

// GetCountries gets the list of accepted countries.
func GetCountries(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(mapping.Countries)
}

// GetAlternativeNamings gets the list of alternative names for countries.
func GetAlternativeNamings(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(mapping.AlternativeNamings)
}

// GetPrefixes gets the map of the prefix submission to alternative country name.
func GetPrefixes(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(mapping.Prefixes)
}

// GetCountriesMap gets the map of country name to correctly formatted name used in the SVG map.
func GetCountriesMap(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(mapping.CountriesMap)
}
