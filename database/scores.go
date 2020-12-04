package database

// Score is the database object for a score entry.
type Score struct {
	ID     int `json:"id"`
	UserID int `json:"userId"`
	QuizID int `json:"quizId"`
	Score  int `json:"score"`
}

// GetScores returns all scores for a given user.
func GetScores(userID int) ([]Score, error) {
	rows, err := Connection.Query("SELECT * FROM scores WHERE userId = $1;", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores = []Score{}
	for rows.Next() {
		var score Score
		if err = rows.Scan(&score.ID, &score.UserID, &score.QuizID, &score.Score); err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}
	return scores, rows.Err()
}

// GetScoreID returns the id for a score for a given user and quiz.
func GetScoreID(userID int, quizID int) (int, error) {
	statement := "SELECT id FROM scores WHERE userId = $1 AND quizId = $2;"
	var id int
	err := Connection.QueryRow(statement, userID, quizID).Scan(&id)
	return id, err
}

// GetScore return a score with the matching id.
func GetScore(id int) (Score, error) {
	statement := "SELECT * FROM scores WHERE id = $1;"
	var score Score
	err := Connection.QueryRow(statement, id).Scan(&score.ID, &score.UserID, &score.QuizID, &score.Score)
	return score, err
}

// InsertScore inserts a score entry into the scores table.
func InsertScore(score Score) (int, error) {
	statement := "INSERT INTO scores (userId, quizId, score) VALUES ($1, $2, $3) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, score.UserID, score.QuizID, score.Score).Scan(&id)
	return id, err
}

// UpdateScore updates a score entry.
func UpdateScore(score Score) (int, error) {
	statement := "UPDATE scores set score = $1 WHERE userId = $2 AND quizId = $3 RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, score.Score, score.UserID, score.QuizID).Scan(&id)
	return id, err
}

// DeleteScore deletes a score entry.
func DeleteScore(scoreID int) error {
	statement := "DELETE FROM scores WHERE id = $1 RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, scoreID).Scan(&id)
	return err
}
