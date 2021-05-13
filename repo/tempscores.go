package repo

import "time"

// TempScore is the database object for a tempscore entry.
type TempScore struct {
	ID    int       `json:"id"`
	Score int       `json:"score"`
	Time  int       `json:"time"`
	Added time.Time `json:"added"`
}

// GetTempScore returns a tempscore with the matching id.
var GetTempScore = func(id int) (TempScore, error) {
	statement := "SELECT * from tempscores WHERE id = $1;"
	var score TempScore
	err := Connection.QueryRow(statement, id).Scan(&score.ID, &score.Score, &score.Time, &score.Added)
	return score, err
}

// InsertTempScore inserts a tempscore entry into the tempscores table.
var InsertTempScore = func(score TempScore) (int, error) {
	statement := "INSERT INTO tempscores (score, time, added) VALUES ($1, $2, $3) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, score.Score, score.Time, score.Added).Scan(&id)
	return id, err
}

// DeleteTempScore deletes a tempscore entry.
var DeleteTempScore = func(id int) error {
	statement := "DELETE FROM tempscores WHERE id = $1 RETURNING id;"
	var tempScoreID int
	err := Connection.QueryRow(statement, id).Scan(&tempScoreID)
	return err
}
