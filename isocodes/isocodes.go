package isocodes

import (
	"encoding/json"
	"net/http"

	mapping "github.com/geobuff/location-mapping/isocodes"
)

// GetCodes gets the map of country name to ISO-3166 code.
func GetCodes(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(mapping.Codes)
}
