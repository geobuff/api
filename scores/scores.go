package scores

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/geobuff/api/config"
	"github.com/geobuff/api/database"
	"github.com/geobuff/api/permissions"
	"github.com/geobuff/auth"
	"github.com/gorilla/mux"
)

// GetScores returns all scores for a given user.
var GetScores = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	userID, err := strconv.Atoi(mux.Vars(request)["userId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	user, err := database.GetUser(userID)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	uv := auth.UserValidation{
		Request:    request,
		UserID:     userID,
		Permission: permissions.ReadScores,
		Identifier: config.Values.Auth0.Identifier,
		Key:        user.Username,
	}

	if code, err := auth.ValidUser(uv); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	scores, err := database.GetScores(userID)
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

	user, err := database.GetUser(score.UserID)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	uv := auth.UserValidation{
		Request:    request,
		UserID:     score.UserID,
		Permission: permissions.WriteScores,
		Identifier: config.Values.Auth0.Identifier,
		Key:        user.Username,
	}

	if code, err := auth.ValidUser(uv); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	score.Added = time.Now()
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

	var score database.Score
	err = json.Unmarshal(requestBody, &score)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	user, err := database.GetUser(score.UserID)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	uv := auth.UserValidation{
		Request:    request,
		UserID:     score.UserID,
		Permission: permissions.WriteScores,
		Identifier: config.Values.Auth0.Identifier,
		Key:        user.Username,
	}

	if code, err := auth.ValidUser(uv); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	score.ID = id
	score.Added = time.Now()
	err = database.UpdateScore(score)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(score)
})

// DeleteScore deletes an existing score entry.
var DeleteScore = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	score, err := database.GetScore(id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	user, err := database.GetUser(score.UserID)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	uv := auth.UserValidation{
		Request:    request,
		UserID:     score.UserID,
		Permission: permissions.WriteScores,
		Identifier: config.Values.Auth0.Identifier,
		Key:        user.Username,
	}

	if code, err := auth.ValidUser(uv); err != nil {
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
})
