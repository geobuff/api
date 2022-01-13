package trivia

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

func GenerateTrivia(writer http.ResponseWriter, request *http.Request) {
	err := repo.CreateTrivia()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

func GetAllTrivia(writer http.ResponseWriter, request *http.Request) {
	trivia, err := repo.GetAllTrivia()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(trivia)
}

func GetTriviaByDate(writer http.ResponseWriter, request *http.Request) {
	trivia, err := repo.GetTrivia(mux.Vars(request)["date"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(trivia)
}
