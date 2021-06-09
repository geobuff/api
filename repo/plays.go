package repo

// GetAllPlays returns a the total play count.
var GetAllPlays = func() (int, error) {
	statement := "SELECT SUM(value) from plays;"
	var count int
	err := Connection.QueryRow(statement).Scan(&count)
	return count, err
}

// GetPlayCount returns a play count with the matching quiz.
var GetPlayCount = func(quizID int) (int, error) {
	statement := "SELECT value from plays WHERE quizId = $1;"
	var count int
	err := Connection.QueryRow(statement, quizID).Scan(&count)
	return count, err
}

// IncrementPlayCount increments a count for a given quiz's plays.
var IncrementPlayCount = func(quizID int) error {
	statement := "UPDATE plays set value = value + 1 WHERE quizId = $1 RETURNING value;"
	var count int
	return Connection.QueryRow(statement, quizID).Scan(&count)
}
