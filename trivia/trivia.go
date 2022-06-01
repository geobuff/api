package trivia

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

type GetTriviaDto struct {
	Trivia  []repo.Trivia `json:"trivia"`
	HasMore bool          `json:"hasMore"`
}

func GenerateTrivia(writer http.ResponseWriter, request *http.Request) {
	err := repo.CreateTrivia()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

func RegenerateTrivia(writer http.ResponseWriter, request *http.Request) {
	if code, err := auth.IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	err := repo.RegenerateTrivia(mux.Vars(request)["date"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

func DeleteTrivia(writer http.ResponseWriter, request *http.Request) {
	if code, err := auth.IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	err := repo.DeleteTriviaByDate(mux.Vars(request)["date"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

func DeleteOldTrivia(writer http.ResponseWriter, request *http.Request) {
	if code, err := auth.IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	newTriviaCount, err := strconv.Atoi(mux.Vars(request)["newTriviaCount"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	err = repo.DeleteOldTrivia(newTriviaCount)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}
}

func GetAllTrivia(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var filter repo.GetTriviaFilter
	err = json.Unmarshal(requestBody, &filter)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	trivia, err := repo.GetAllTrivia(filter)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	switch _, err := repo.GetFirstTriviaID((filter.Page + 1) * filter.Limit); err {
	case sql.ErrNoRows:
		entriesDto := GetTriviaDto{trivia, false}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entriesDto)
	case nil:
		entriesDto := GetTriviaDto{trivia, true}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entriesDto)
	default:
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
	}
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
