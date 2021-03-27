package database

import "time"

// Score is the database object for a score entry.
type Score struct {
	ID     int       `json:"id"`
	UserID int       `json:"userId"`
	QuizID int       `json:"quizId"`
	Score  int       `json:"score"`
	Time   int       `json:"time"`
	Added  time.Time `json:"added"`
}

// ScoreDto provides us with additional fields used in the user profile.
type ScoreDto struct {
	ID           int       `json:"id"`
	UserID       int       `json:"userId"`
	QuizID       int       `json:"quizId"`
	QuizName     string    `json:"quizName"`
	QuizImageUrl string    `json:"quizImageUrl"`
	QuizMaxScore int       `json:"quizMaxScore"`
	Score        int       `json:"score"`
	Time         int       `json:"time"`
	Added        time.Time `json:"added"`
}

// GetScores returns all scores for a given user.
var GetScores = func(userID int) ([]ScoreDto, error) {
	rows, err := Connection.Query("SELECT s.id, s.userid, q.id, q.name, q.imageUrl, q.maxScore, s.score, s.time, s.added FROM scores s JOIN quizzes q ON q.id = s.quizid WHERE userId = $1;", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores = []ScoreDto{}
	for rows.Next() {
		var score ScoreDto
		if err = rows.Scan(&score.ID, &score.UserID, &score.QuizID, &score.QuizName, &score.QuizImageUrl, &score.QuizMaxScore, &score.Score, &score.Time, &score.Added); err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}
	return scores, rows.Err()
}

// GetScore return a score with the matching id.
var GetScore = func(id int) (ScoreDto, error) {
	statement := "SELECT s.id, s.userid, q.id, q.name, q.imageUrl, q.maxScore, s.score, s.time, s.added FROM scores s JOIN quizzes q ON q.id = s.quizid WHERE id = $1;"
	var score ScoreDto
	err := Connection.QueryRow(statement, id).Scan(&score.ID, &score.UserID, &score.QuizID, &score.QuizName, &score.QuizImageUrl, &score.QuizMaxScore, &score.Score, &score.Time, &score.Added)
	return score, err
}

// GetUserScore returns a score matching a user and quiz ID.
var GetUserScore = func(userID, quizID int) (Score, error) {
	statement := "SELECT * FROM scores WHERE userid = $1 AND quizid = $2;"
	var score Score
	err := Connection.QueryRow(statement, userID, quizID).Scan(&score.ID, &score.UserID, &score.QuizID, &score.Score, &score.Time, &score.Added)
	return score, err
}

// InsertScore inserts a score entry into the scores table.
var InsertScore = func(score Score) (int, error) {
	statement := "INSERT INTO scores (userId, quizId, score, time, added) VALUES ($1, $2, $3, $4, $5) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, score.UserID, score.QuizID, score.Score, score.Time, score.Added).Scan(&id)
	return id, err
}

// UpdateScore updates a score entry.
var UpdateScore = func(score Score) error {
	statement := "UPDATE scores set userid = $1, quizid = $2, score = $3, time = $4, added = $5 WHERE id = $6 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, score.UserID, score.QuizID, score.Score, score.Time, score.Added, score.ID).Scan(&id)
}

// DeleteScore deletes a score entry.
var DeleteScore = func(scoreID int) error {
	statement := "DELETE FROM scores WHERE id = $1 RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, scoreID).Scan(&id)
	return err
}
