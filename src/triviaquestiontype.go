package src

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/geobuff/api/repo"
)

func GetTriviaQuestionTypes(writer http.ResponseWriter, request *http.Request) {
	types, err := repo.GetTriviaQuestionTypes()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(types)
}
