package countries

import (
	"encoding/json"
	"net/http"

	"github.com/geobuff/mapping"
)

// GetCountries gets the map of accepted countries.
func GetCountries(writer http.ResponseWriter, request *http.Request) {
	continentParam := request.URL.Query().Get("continent")
	writer.Header().Set("Content-Type", "application/json")

	if continentParam != "" {
		json.NewEncoder(writer).Encode(mapping.Countries[continentParam])
	} else {
		json.NewEncoder(writer).Encode(mapping.Countries)
	}
}
