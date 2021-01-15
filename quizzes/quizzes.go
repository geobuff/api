package quizzes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/geobuff/api/database"
	"github.com/gorilla/mux"
)

// GetQuizzes returns all quizzes.
func GetQuizzes(writer http.ResponseWriter, request *http.Request) {
	filterParam := request.URL.Query().Get("filter")
	quizzes, err := database.GetQuizzes(filterParam)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(quizzes)
}

// GetQuiz gets a quiz entry by id.
func GetQuiz(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	quiz, err := database.GetQuiz(id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(quiz)
}
