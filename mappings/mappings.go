package mappings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/geobuff/api/auth"
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

func EditMapping(writer http.ResponseWriter, request *http.Request) {
	if code, err := auth.IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var mapping repo.UpdateMappingDto
	err = json.Unmarshal(requestBody, &mapping)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	if err = repo.UpdateMapping(mux.Vars(request)["key"], mapping); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

func DeleteMapping(writer http.ResponseWriter, request *http.Request) {
	if code, err := auth.IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	if err := repo.DeleteMapping(mux.Vars(request)["key"]); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
	}
}
