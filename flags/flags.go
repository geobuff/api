package flags

import (
	"encoding/json"
	"net/http"

	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

func GetFlagGroups(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(repo.FlagGroups)
}

func GetFlagGroup(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	key := mux.Vars(request)["key"]
	for _, group := range repo.FlagGroups {
		if group.Key == key {
			json.NewEncoder(writer).Encode(group)
			return
		}
	}

	http.Error(writer, "invalid key; does not match existing flag group", http.StatusBadRequest)
}
