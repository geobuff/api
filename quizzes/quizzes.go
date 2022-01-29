package quizzes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/geobuff/api/auth"
	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

type ToggleQuizDto struct {
	quizID int `json:"quizId"`
}

// GetQuizzes returns all quizzes.
func GetQuizzes(writer http.ResponseWriter, request *http.Request) {
	filterParam := request.URL.Query().Get("filter")
	quizzes, err := repo.GetQuizzes(filterParam)
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

	quiz, err := repo.GetQuiz(id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(quiz)
}

func ToggleQuizEnabled(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	if code, err := auth.IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	_, err = repo.ToggleQuizEnabled(id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}
