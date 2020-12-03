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
	"github.com/geobuff/geobuff-api/keys"
	"github.com/gorilla/mux"
)

// User is the database object for a user entry.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// PageDto is used to display a paged result of user entries.
type PageDto struct {
	Users   []User `json:"users"`
	HasMore bool   `json:"hasMore"`
}

// GetUsers gets the user entries for a given page.
var GetUsers = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	if hasPermission, err := auth.HasPermission(request, auth.ReadUsers); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
		return
	} else if !hasPermission {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	pageParam := request.URL.Query().Get("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	rows, err := database.DBConnection.Query("SELECT * FROM users LIMIT $1 OFFSET $2;", 10, strconv.Itoa(page*10))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}
	defer rows.Close()

	var users = []User{}
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Username)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(writer, "%v\n", err)
			return
		}

		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	statement := "SELECT id FROM users LIMIT $1 OFFSET $2;"
	row := database.DBConnection.QueryRow(statement, 1, strconv.Itoa((page+1)*10))

	var id int
	switch err = row.Scan(&id); err {
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

	statement := "SELECT * FROM users WHERE id = $1;"
	row := database.DBConnection.QueryRow(statement, id)

	var user User
	switch err = row.Scan(&user.ID, &user.Username); err {
	case sql.ErrNoRows:
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "%v\n", err)
	case nil:
		if code, err := auth.ValidUser(request, user.ID, auth.ReadUsers); err != nil {
			writer.WriteHeader(code)
			fmt.Fprintf(writer, "%v\n", err)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(user)
	default:
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
	}
})

// CreateUser creates a new user entry.
func CreateUser(writer http.ResponseWriter, request *http.Request) {
	key := request.Header.Get("x-api-key")
	if !keys.ValidAuth0Key(key) {
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

	var newUser User
	err = json.Unmarshal(requestBody, &newUser)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}

	statement := "INSERT INTO users (username) VALUES ($1) RETURNING id;"

	var id int
	err = database.DBConnection.QueryRow(statement, newUser.Username).Scan(&id)
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

	scoresStatement := "DELETE FROM scores WHERE userId = $1;"
	database.DBConnection.QueryRow(scoresStatement, id)
	leaderboardStatement := "DELETE FROM world_leaderboard WHERE userId = $1;"
	database.DBConnection.QueryRow(leaderboardStatement, id)
	usersStatement := "DELETE FROM users WHERE id = $1 RETURNING *;"
	row := database.DBConnection.QueryRow(usersStatement, id)

	var user User
	switch err = row.Scan(&user.ID, &user.Username); err {
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
