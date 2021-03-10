package capitals

import (
	"encoding/json"
	"net/http"

	"github.com/geobuff/mapping"
)

// GetCapitals gets the list of accepted capitals.
func GetCapitals(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(mapping.WorldCapitals)
}
