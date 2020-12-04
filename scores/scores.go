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
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	if code, err := auth.ValidUser(request, id, auth.ReadScores); err != nil {
		writer.WriteHeader(code)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	scores, err := database.GetScores(id)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(scores)
})

// SetScore upserts a score.
var SetScore = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	var score database.Score
	err = json.Unmarshal(requestBody, &score)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	if code, err := auth.ValidUser(request, score.UserID, auth.WriteScores); err != nil {
		writer.WriteHeader(code)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	switch id, err := database.GetScoreID(score.UserID, score.QuizID); err {
	case sql.ErrNoRows:
		id, err = database.InsertScore(score)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(writer, "%v\n", err)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusCreated)
		score.ID = id
		json.NewEncoder(writer).Encode(score)
	case nil:
		id, err = database.UpdateScore(score)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(writer, "%v\n", err)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		score.ID = id
		json.NewEncoder(writer).Encode(score)
	default:
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
	}
})

// DeleteScore deletes an existing score entry.
var DeleteScore = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	switch score, err := database.GetScore(id); err {
	case sql.ErrNoRows:
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "%v\n", err)
	case nil:
		if code, err := auth.ValidUser(request, score.UserID, auth.WriteScores); err != nil {
			writer.WriteHeader(code)
			fmt.Fprintf(writer, "%v\n", err)
			return
		}

		err = database.DeleteScore(id)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(writer, "%v\n", err)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(score)
	default:
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
	}
})
