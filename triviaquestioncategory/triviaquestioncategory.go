package triviaquestioncategory

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/geobuff/api/repo"
)

func GetTriviaQuestionCategories(writer http.ResponseWriter, request *http.Request) {
	categories, err := repo.GetTriviaQuestionCategories()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(categories)
}
