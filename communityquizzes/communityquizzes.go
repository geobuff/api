package communityquizzes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/geobuff/api/auth"
	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

type CommunityQuizPageDto struct {
	Quizzes []repo.CommunityQuizDto `json:"quizzes"`
	HasMore bool                    `json:"hasMore"`
}

func GetCommunityQuizzes(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var filter repo.GetCommunityQuizzesFilter
	err = json.Unmarshal(requestBody, &filter)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	quizzes, err := repo.GetCommunityQuizzes(filter)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	switch _, err := repo.GetFirstCommunityQuizID((filter.Page + 1) * filter.Limit); err {
	case sql.ErrNoRows:
		entriesDto := CommunityQuizPageDto{quizzes, false}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entriesDto)
	case nil:
		entriesDto := CommunityQuizPageDto{quizzes, true}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entriesDto)
	default:
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
	}
}

func CreateCommunityQuiz(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var quiz repo.CreateCommunityQuizDto
	err = json.Unmarshal(requestBody, &quiz)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	if code, err := auth.ValidUser(request, quiz.UserID); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	if err := repo.InsertCommunityQuiz(quiz); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

func UpdateCommunityQuiz(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var quiz repo.UpdateCommunityQuizDto
	err = json.Unmarshal(requestBody, &quiz)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	if code, err := auth.ValidUser(request, quiz.UserID); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	if err := repo.UpdateCommunityQuiz(id, quiz); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

func DeleteCommunityQuiz(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	quiz, err := repo.GetCommunityQuiz(id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	if code, err := auth.ValidUser(request, quiz.UserID); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	err = repo.DeleteCommunityQuiz(id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}
