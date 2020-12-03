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

// Score is the database object for a score entry.
type Score struct {
	ID     int `json:"id"`
	UserID int `json:"userId"`
	QuizID int `json:"quizId"`
	Score  int `json:"score"`
}

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

	rows, err := database.DBConnection.Query("SELECT * FROM scores WHERE userId = $1;", id)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
		return
	}
	defer rows.Close()

	var scores = []Score{}
	for rows.Next() {
		var score Score
		err = rows.Scan(&score.ID, &score.UserID, &score.QuizID, &score.Score)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(writer, "%v\n", err)
			return
		}
		scores = append(scores, score)
	}

	err = rows.Err()
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

	var score Score
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

	statement := "SELECT id FROM scores WHERE userId = $1 AND quizId = $2;"
	row := database.DBConnection.QueryRow(statement, score.UserID, score.QuizID)

	var id int
	switch err = row.Scan(&id); err {
	case sql.ErrNoRows:
		statement := "INSERT INTO scores (userId, quizId, score) VALUES ($1, $2, $3) RETURNING id;"
		err = database.DBConnection.QueryRow(statement, score.UserID, score.QuizID, score.Score).Scan(&id)
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
		statement := "UPDATE scores set score = $1 WHERE userId = $2 AND quizId = $3 RETURNING id;"
		err = database.DBConnection.QueryRow(statement, score.Score, score.UserID, score.QuizID).Scan(&id)
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

	statement := "SELECT * FROM scores WHERE id = $1;"
	row := database.DBConnection.QueryRow(statement, id)

	var score Score
	switch err = row.Scan(&score.ID, &score.UserID, &score.QuizID, &score.Score); err {
	case sql.ErrNoRows:
		writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(writer, "%v\n", err)
		return
	case nil:
		if code, err := auth.ValidUser(request, score.UserID, auth.WriteScores); err != nil {
			writer.WriteHeader(code)
			fmt.Fprintf(writer, "%v\n", err)
			return
		}

		statement = "DELETE FROM scores WHERE id = $1;"
		database.DBConnection.QueryRow(statement, id)

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(score)
	default:
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "%v\n", err)
	}
})
