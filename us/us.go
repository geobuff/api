package us

import (
	"encoding/json"
	"net/http"

	"github.com/geobuff/mapping"
)

// GetStates gets the list of accepted states.
func GetStates(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(mapping.USStates)
}
