package capitals

import (
	"encoding/json"
	"net/http"

	"github.com/geobuff/mapping"
)

// GetCapitals gets the map of accepted capitals.
func GetCapitals(writer http.ResponseWriter, request *http.Request) {
	continentParam := request.URL.Query().Get("continent")
	writer.Header().Set("Content-Type", "application/json")

	if continentParam != "" {
		json.NewEncoder(writer).Encode(mapping.WorldCapitals[continentParam])
	} else {
		json.NewEncoder(writer).Encode(mapping.WorldCapitals)
	}
}
