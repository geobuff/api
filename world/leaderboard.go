package world

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

// EntriesDto is used to display a paged result of leaderboard entries.
type EntriesDto struct {
	Entries []database.Entry `json:"entries"`
	HasMore bool             `json:"hasMore"`
}

// GetEntries gets the leaderboard entries for a given page.
func GetEntries(writer http.ResponseWriter, request *http.Request) {
	pageParam := request.URL.Query().Get("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	entries, err := database.GetWorldLeaderboardEntries(10, page*10)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	switch _, err := database.GetWorldLeaderboardEntryID(1, (page+1)*10); err {
	case sql.ErrNoRows:
		entriesDto := EntriesDto{entries, false}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entriesDto)
	case nil:
		entriesDto := EntriesDto{entries, true}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entriesDto)
	default:
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
	}
}

// GetEntry gets a leaderboard entry by id.
func GetEntry(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	switch entry, err := database.GetWorldLeaderboardEntry(id); err {
	case sql.ErrNoRows:
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "%v\n", err)
	case nil:
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entry)
	default:
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
	}
}

// CreateEntry creates a new leaderboard entry.
var CreateEntry = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	var newEntry database.Entry
	err = json.Unmarshal(requestBody, &newEntry)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	if code, err := auth.ValidUser(request, newEntry.UserID, auth.WriteLeaderboard); err != nil {
		writer.WriteHeader(code)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	id, err := database.InsertWorldLeaderboardEntry(newEntry)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
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
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	var updatedEntry database.Entry
	err = json.Unmarshal(requestBody, &updatedEntry)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	if code, err := auth.ValidUser(request, updatedEntry.UserID, auth.WriteLeaderboard); err != nil {
		writer.WriteHeader(code)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	switch id, err = database.UpdateWorldLeaderboardEntry(id, updatedEntry); err {
	case sql.ErrNoRows:
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "%v\n", err)
	case nil:
		writer.Header().Set("Content-Type", "application/json")
		updatedEntry.ID = id
		json.NewEncoder(writer).Encode(updatedEntry)
	default:
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
	}
})

// DeleteEntry deletes an existing leaderboard entry.
var DeleteEntry = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	switch entry, err := database.GetWorldLeaderboardEntry(id); err {
	case sql.ErrNoRows:
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "%v\n", err)
	case nil:
		if code, err := auth.ValidUser(request, entry.UserID, auth.WriteLeaderboard); err != nil {
			writer.WriteHeader(code)
			fmt.Fprintf(writer, "%v\n", err)
			return
		}

		_, err := database.DeleteWorldLeaderboardEntry(entry.ID)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(writer, "%v\n", err)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entry)
	default:
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
	}
})
