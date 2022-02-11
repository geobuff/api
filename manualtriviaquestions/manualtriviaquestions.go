package manualtriviaquestions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/geobuff/api/auth"
	"github.com/geobuff/api/repo"
)

type ManualQuestionsDto struct {
	Questions []repo.ManualTriviaQuestion `json:"questions"`
	HasMore   bool                        `json:"hasMore"`
}

func GetManualTriviaQuestions(writer http.ResponseWriter, request *http.Request) {
	pageParam := request.URL.Query().Get("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	if code, err := auth.IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	questions, err := repo.GetAllManualTriviaQuestions(10, page*10)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	switch _, err := repo.GetFirstManualTriviaQuestionID((page + 1) * 10); err {
	case sql.ErrNoRows:
		entriesDto := ManualQuestionsDto{questions, false}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entriesDto)
	case nil:
		entriesDto := ManualQuestionsDto{questions, true}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entriesDto)
	default:
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
	}
}
