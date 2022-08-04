package maps

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

func GetMaps(writer http.ResponseWriter, request *http.Request) {
	maps, err := repo.GetMaps()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(maps)
}

func GetMap(writer http.ResponseWriter, request *http.Request) {
	svgMap, err := repo.GetMap(mux.Vars(request)["className"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(svgMap)
}

func GetMapHighlightedRegions(writer http.ResponseWriter, request *http.Request) {
	regions, err := repo.GetMapHighlightedRegions(mux.Vars(request)["className"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(regions)
}
