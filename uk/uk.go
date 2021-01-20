package uk

import (
	"encoding/json"
	"net/http"

	"github.com/geobuff/mapping"
)

// GetCounties gets the list of accepted counties.
func GetCounties(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(mapping.UKCounties)
}
