package scores

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/geobuff/geobuff-api/auth"
	"github.com/geobuff/geobuff-api/database"
	"github.com/gorilla/mux"
)

// GetScores returns all scores for a given user.
var GetScores = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	if code, err := auth.ValidUser(request, id, auth.ReadScores); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	scores, err := database.GetScores(id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(scores)
})

// CreateScore creates a new score.
var CreateScore = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var score database.Score
	err = json.Unmarshal(requestBody, &score)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	if code, err := auth.ValidUser(request, score.UserID, auth.WriteScores); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	id, err := database.InsertScore(score)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	score.ID = id
	json.NewEncoder(writer).Encode(score)
})

// UpdateScore updates an existing score.
var UpdateScore = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var score database.Score
	err = json.Unmarshal(requestBody, &score)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	if code, err := auth.ValidUser(request, score.UserID, auth.WriteScores); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	id, err := database.GetScoreID(score.UserID, score.QuizID)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	id, err = database.UpdateScore(score)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	score.ID = id
	json.NewEncoder(writer).Encode(score)
})

// DeleteScore deletes an existing score entry.
var DeleteScore = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	switch score, err := database.GetScore(id); err {
	case sql.ErrNoRows:
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusNotFound)
	case nil:
		if code, err := auth.ValidUser(request, score.UserID, auth.WriteScores); err != nil {
			http.Error(writer, fmt.Sprintf("%v\n", err), code)
			return
		}

		err = database.DeleteScore(id)
		if err != nil {
			http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(score)
	default:
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
	}
})
