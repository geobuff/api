package mappings

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

func GetMappingGroups(writer http.ResponseWriter, request *http.Request) {
	mapping, err := repo.GetMappingGroups()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(mapping)
}

func GetMappingEntries(writer http.ResponseWriter, request *http.Request) {
	mapping, err := repo.GetMappingEntries(mux.Vars(request)["key"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(mapping)
}

func GetMappingsWithoutFlags(writer http.ResponseWriter, request *http.Request) {
	keys, err := repo.GetMappingsWithoutFlags()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(keys)
}
