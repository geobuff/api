package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/geobuff/api/config"
	"github.com/geobuff/api/database"
	"github.com/geobuff/api/permissions"
	"github.com/geobuff/auth"
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
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	up := auth.UserPermission{
		Request:    request,
		Permission: permissions.ReadUsers,
	}

	if hasPermission, err := auth.HasPermission(up); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	} else if !hasPermission {
		http.Error(writer, "invalid permissions to make request", http.StatusUnauthorized)
		return
	}

	users, err := database.GetUsers(10, page*10)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
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
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
	}
})

// GetUser gets a user entry by id.
var GetUser = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	user, err := database.GetUser(id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	uv := auth.UserValidation{
		Request:    request,
		UserID:     id,
		Permission: permissions.ReadUsers,
		Identifier: config.Values.Auth0.Identifier,
		Key:        user.Username,
	}

	if code, err := auth.ValidUser(uv); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(user)
})

// CreateUser creates a new user entry.
func CreateUser(writer http.ResponseWriter, request *http.Request) {
	if entry, err := database.GetKey("auth0"); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	} else if entry.Key != request.Header.Get("x-api-key") {
		http.Error(writer, fmt.Sprintf("%v\n", "invalid api key"), http.StatusUnauthorized)
		return
	}

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var newUser database.User
	err = json.Unmarshal(requestBody, &newUser)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	id, err := database.InsertUser(newUser)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
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
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	user, err := database.GetUser(id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	uv := auth.UserValidation{
		Request:    request,
		UserID:     id,
		Permission: permissions.ReadUsers,
		Identifier: config.Values.Auth0.Identifier,
		Key:        user.Username,
	}

	if code, err := auth.ValidUser(uv); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	err = database.DeleteUser(id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(user)
})
