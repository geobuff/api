package users

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

// PageDto is used to display a paged result of user entries.
type PageDto struct {
	Users   []database.User `json:"users"`
	HasMore bool            `json:"hasMore"`
}

// GetUsers gets the user entries for a given page.
var GetUsers = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	pageParam := request.URL.Query().Get("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	if hasPermission, err := auth.HasPermission(request, auth.ReadUsers); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
		return
	} else if !hasPermission {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	users, err := database.GetUsers(10, page*10)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	switch _, err := database.GetUserID(1, (page+1)*10); err {
	case sql.ErrNoRows:
		entriesDto := PageDto{users, false}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entriesDto)
	case nil:
		entriesDto := PageDto{users, true}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entriesDto)
	default:
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
	}
})

// GetUser gets a user entry by id.
var GetUser = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	if code, err := auth.ValidUser(request, id, auth.ReadUsers); err != nil {
		writer.WriteHeader(code)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	switch user, err := database.GetUser(id); err {
	case sql.ErrNoRows:
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "%v\n", err)
	case nil:
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(user)
	default:
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
	}
})

// CreateUser creates a new user entry.
func CreateUser(writer http.ResponseWriter, request *http.Request) {
	if entry, err := database.GetKey("auth0"); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
		return
	} else if entry.Key != request.Header.Get("x-api-key") {
		writer.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(writer, "%v\n", "Invalid api key.")
		return
	}

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	var newUser database.User
	err = json.Unmarshal(requestBody, &newUser)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	id, err := database.InsertUser(newUser)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	newUser.ID = id
	json.NewEncoder(writer).Encode(newUser)
}

// DeleteUser deletes an existing user entry.
var DeleteUser = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	if code, err := auth.ValidUser(request, id, auth.WriteUsers); err != nil {
		writer.WriteHeader(code)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	switch user, err := database.DeleteUser(id); err {
	case sql.ErrNoRows:
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "%v\n", err)
	case nil:
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(user)
	default:
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
	}
})
