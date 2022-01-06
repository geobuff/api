package dailytrivia

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

func GenerateDailyTrivia(writer http.ResponseWriter, request *http.Request) {
	err := repo.CreateDailyTrivia()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

func GetAllDailyTrivia(writer http.ResponseWriter, request *http.Request) {
	trivia, err := repo.GetAllDailyTrivia()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(trivia)
}

func GetDailyTriviaByDate(writer http.ResponseWriter, request *http.Request) {
	trivia, err := repo.GetDailyTrivia(mux.Vars(request)["date"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(trivia)
}

func IncrementPlays(writer http.ResponseWriter, request *http.Request) {
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
