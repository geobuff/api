package countries

import (
	"encoding/json"
	"net/http"

	"github.com/geobuff/mapping"
)

// GetCountries gets the list of accepted countries.
func GetCountries(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(mapping.WorldCountries)
}
