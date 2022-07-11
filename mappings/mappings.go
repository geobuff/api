package mappings

import (
	"encoding/json"
	"net/http"

	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

func GetMapping(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(repo.Mappings[mux.Vars(request)["key"]])
}
