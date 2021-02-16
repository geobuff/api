package countries

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/geobuff/api/config"
	"github.com/geobuff/api/database"
	"github.com/geobuff/api/models"
	"github.com/geobuff/api/permissions"
	"github.com/geobuff/auth"
	"github.com/gorilla/mux"
)

// EntriesDto is used to display a paged result of leaderboard entries.
type EntriesDto struct {
	Entries []database.LeaderboardEntry `json:"entries"`
	HasMore bool                        `json:"hasMore"`
}

// GetEntries gets the leaderboard entries for a given page.
func GetEntries(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var filterParams models.GetEntriesFilterParams
	err = json.Unmarshal(requestBody, &filterParams)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	entries, err := database.GetLeaderboardEntries(database.CountriesTable, filterParams)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	switch _, err := database.GetLeaderboardEntryID(database.CountriesTable, filterParams); err {
	case sql.ErrNoRows:
		entriesDto := EntriesDto{entries, false}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entriesDto)
	case nil:
		entriesDto := EntriesDto{entries, true}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entriesDto)
	default:
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
	}
}

// GetEntry gets a leaderboard entry by user id.
func GetEntry(writer http.ResponseWriter, request *http.Request) {
	userID, err := strconv.Atoi(mux.Vars(request)["userId"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	entry, err := database.GetLeaderboardEntry(database.CountriesTable, userID)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(entry)
}

// CreateEntry creates a new leaderboard entry.
var CreateEntry = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var newEntry database.LeaderboardEntry
	err = json.Unmarshal(requestBody, &newEntry)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	user, err := database.GetUser(newEntry.UserID)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	uv := auth.UserValidation{
		Request:    request,
		Permission: permissions.WriteLeaderboard,
		Identifier: config.Values.Auth0.Identifier,
		Key:        user.Username,
	}

	if code, err := auth.ValidUser(uv); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	newEntry.Added = time.Now()
	id, err := database.InsertLeaderboardEntry(database.CountriesTable, newEntry)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	newEntry.ID = id
	json.NewEncoder(writer).Encode(newEntry)
})

// UpdateEntry updates an existing leaderboard entry.
var UpdateEntry = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
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

	var updatedEntry database.LeaderboardEntry
	err = json.Unmarshal(requestBody, &updatedEntry)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	user, err := database.GetUser(updatedEntry.UserID)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	uv := auth.UserValidation{
		Request:    request,
		Permission: permissions.WriteLeaderboard,
		Identifier: config.Values.Auth0.Identifier,
		Key:        user.Username,
	}

	if code, err := auth.ValidUser(uv); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	updatedEntry.ID = id
	updatedEntry.Added = time.Now()
	err = database.UpdateLeaderboardEntry(database.CountriesTable, updatedEntry)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(updatedEntry)
})

// DeleteEntry deletes an existing leaderboard entry.
var DeleteEntry = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	entry, err := database.GetLeaderboardEntry(database.CountriesTable, id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	user, err := database.GetUser(entry.UserID)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	uv := auth.UserValidation{
		Request:    request,
		Permission: permissions.WriteLeaderboard,
		Identifier: config.Values.Auth0.Identifier,
		Key:        user.Username,
	}

	if code, err := auth.ValidUser(uv); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	err = database.DeleteLeaderboardEntry(database.CountriesTable, entry.ID)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(entry)
})
