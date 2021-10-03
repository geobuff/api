package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/geobuff/api/auth"
	"github.com/geobuff/api/repo"
	"github.com/geobuff/api/validation"
	"github.com/gorilla/mux"
)

// PageDto is used to display a paged result of user entries.
type PageDto struct {
	Users   []repo.UserDto `json:"users"`
	HasMore bool           `json:"hasMore"`
}

// GetUsers gets the user entries for a given page.
func GetUsers(writer http.ResponseWriter, request *http.Request) {
	pageParam := request.URL.Query().Get("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	if code, err := auth.IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	users, err := repo.GetUsers(10, page*10)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	switch _, err := repo.GetFirstID(1, (page+1)*10); err {
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
}

// GetUser gets a user entry by id.
func GetUser(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	if code, err := auth.ValidUser(request, id); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	user, err := repo.GetUser(id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(user)
}

// UpdateUser creates a new user entry.
func UpdateUser(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	if code, err := auth.ValidUser(request, id); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var updatedUser repo.UpdateUserDto
	err = json.Unmarshal(requestBody, &updatedUser)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	err = validation.Validator.Struct(updatedUser)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	usernameExists, err := repo.AnotherUserWithUsername(id, updatedUser.Username)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	if usernameExists {
		http.Error(writer, "Username already in use. Please choose another and try again.", http.StatusBadRequest)
		return
	}

	emailExists, err := repo.AnotherUserWithEmail(id, updatedUser.Email)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	if emailExists {
		http.Error(writer, "Email already in use. Please choose another and try again.", http.StatusBadRequest)
		return
	}

	if err = repo.UpdateUser(id, updatedUser); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	avatar, err := repo.GetAvatar(updatedUser.AvatarId)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	updatedUser.AvatarName = avatar.Name
	updatedUser.AvatarImageUrl = avatar.ImageUrl
	updatedUser.AvatarBackground = avatar.Background
	updatedUser.AvatarBorder = avatar.Border
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(updatedUser)
}

// DeleteUser deletes an existing user entry.
func DeleteUser(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	if code, err := auth.ValidUser(request, id); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	user, err := repo.GetUser(id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	err = repo.DeleteUser(id)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(user)
}
