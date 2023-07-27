package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

func GetFlagGroups(writer http.ResponseWriter, request *http.Request) {
	groups, err := repo.GetFlagGroups()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(groups)
}

func GetFlagEntries(writer http.ResponseWriter, request *http.Request) {
	entries, err := repo.GetFlagEntries(mux.Vars(request)["key"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(entries)
}

func GetFlagUrl(writer http.ResponseWriter, request *http.Request) {
	url, err := repo.GetFlagUrl(mux.Vars(request)["code"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(url)
}

func CreateFlags(writer http.ResponseWriter, request *http.Request) {
	if code, err := IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var newFlags repo.CreateFlagsDto
	err = json.Unmarshal(requestBody, &newFlags)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	if err := repo.CreateFlags(newFlags); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}
