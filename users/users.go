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

type UserPageDto struct {
	Users   []repo.UserDto `json:"users"`
	HasMore bool           `json:"hasMore"`
}

func GetUsers(writer http.ResponseWriter, request *http.Request) {
	if code, err := auth.IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	requestBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	var filterParams repo.GetUsersFilterParams
	err = json.Unmarshal(requestBody, &filterParams)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	users, err := repo.GetUsers(filterParams)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	switch _, err := repo.GetFirstUserID(filterParams); err {
	case sql.ErrNoRows:
		entriesDto := UserPageDto{users, false}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entriesDto)
	case nil:
		entriesDto := UserPageDto{users, true}
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(entriesDto)
	default:
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
	}
}

func GetUser(writer http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	switch user, err := repo.GetUser(id); err {
	case sql.ErrNoRows:
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusNoContent)
	case nil:
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(user)
	default:
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
	}
}

func GetUserByEmail(writer http.ResponseWriter, request *http.Request) {
	switch user, err := repo.GetUserByEmail(mux.Vars(request)["email"]); err {
	case sql.ErrNoRows:
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusNoContent)
	case nil:
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(user)
	default:
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
	}
}

func GetLastWeekTotalUsers(writer http.ResponseWriter, request *http.Request) {
	if code, err := auth.IsAdmin(request); err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), code)
		return
	}

	data, err := repo.GetLastWeekTotalUsers()
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(data)
}

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
	updatedUser.AvatarDescription = avatar.Description
	updatedUser.AvatarPrimaryImageUrl = avatar.PrimaryImageUrl
	updatedUser.AvatarSecondaryImageUrl = avatar.SecondaryImageUrl

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(updatedUser)
}

func UpdateUserXP(writer http.ResponseWriter, request *http.Request) {
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

	var dto repo.UpdateUserXPDto
	err = json.Unmarshal(requestBody, &dto)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusBadRequest)
		return
	}

	err = validation.Validator.Struct(dto)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	increase, err := repo.UpdateUserXP(id, dto.Score, dto.MaxScore)
	if err != nil {
		http.Error(writer, fmt.Sprintf("%v\n", err), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(increase)
}

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
