package quizplays

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/geobuff/api/auth"
	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

// GetAllPlays gets the total play count.
func GetAllQuizPlays(writer http.ResponseWriter, request *http.Request) {
	plays, err := repo.GetAllQuizPlays()
	if err != nil && strings.Contains(err.Error(), "sql: Scan error on column index 0") {
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(0)
		return
	} else if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(plays)
}

// GetPlays gets a play count for a given quiz.
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
	if code, err := auth.IsAdmin(request); err != nil {
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

// IncrementPlays increments the play count for a quiz.
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
