package flags

import (
	"encoding/json"
	"fmt"
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
