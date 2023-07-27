package src

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

func GetLastWeekTriviaPlays(writer http.ResponseWriter, request *http.Request) {
	if code, err := IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	plays, err := repo.GetLastWeekTriviaPlays()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(plays)
}

func IncrementTriviaPlays(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	err = repo.IncrementTriviaPlays(id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}
