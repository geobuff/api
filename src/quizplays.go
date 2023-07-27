package src

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

func GetAllQuizPlays(writer http.ResponseWriter, request *http.Request) {
	plays, err := repo.GetAllQuizPlays()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(plays)
}

func GetQuizPlays(writer http.ResponseWriter, request *http.Request) {
	quizID, err := strconv.Atoi(mux.Vars(request)["quizId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	count, err := repo.GetQuizPlayCount(quizID)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(count)
}

func GetTopFiveQuizPlays(writer http.ResponseWriter, request *http.Request) {
	if code, err := IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	plays, err := repo.GetTopFiveQuizPlays()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(plays)
}

func IncrementQuizPlays(writer http.ResponseWriter, request *http.Request) {
	quizID, err := strconv.Atoi(mux.Vars(request)["quizId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	err = repo.IncrementQuizPlayCount(quizID)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}
