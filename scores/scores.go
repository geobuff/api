package scores

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
