package countries

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/geobuff/api/auth"
	"github.com/geobuff/api/models"
	"github.com/geobuff/api/repo"
	"github.com/gorilla/mux"
)

// EntriesDto is used to display a paged result of leaderboard entries.
type EntriesDto struct {
	Entries []repo.LeaderboardEntryDto `json:"entries"`
	HasMore bool                       `json:"hasMore"`
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

	entries, err := repo.GetLeaderboardEntries(repo.CountriesTable, filterParams)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	switch _, err := repo.GetLeaderboardEntryID(repo.CountriesTable, filterParams); err {
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

	switch entry, err := repo.GetLeaderboardEntry(repo.CountriesTable, userID); err {
	case sql.ErrNoRows:
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusNoContent)
	case nil:
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entry)
	default:
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
	}
}

// CreateEntry creates a new leaderboard entry.
func CreateEntry(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var newEntry repo.LeaderboardEntry
	err = json.Unmarshal(requestBody, &newEntry)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	if code, err := auth.ValidUser(request, newEntry.UserID); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	newEntry.Added = time.Now()
	id, err := repo.InsertLeaderboardEntry(repo.CountriesTable, newEntry)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	newEntry.ID = id
	json.NewEncoder(writer).Encode(newEntry)
}

// UpdateEntry updates an existing leaderboard entry.
func UpdateEntry(writer http.ResponseWriter, request *http.Request) {
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

	var updatedEntry repo.LeaderboardEntry
	err = json.Unmarshal(requestBody, &updatedEntry)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	if code, err := auth.ValidUser(request, updatedEntry.UserID); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	updatedEntry.ID = id
	updatedEntry.Added = time.Now()
	err = repo.UpdateLeaderboardEntry(repo.CountriesTable, updatedEntry)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(updatedEntry)
}

// DeleteEntry deletes an existing leaderboard entry.
func DeleteEntry(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	entry, err := repo.GetLeaderboardEntry(repo.CountriesTable, id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	if code, err := auth.ValidUser(request, entry.UserID); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	err = repo.DeleteLeaderboardEntry(repo.CountriesTable, entry.ID)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(entry)
}
