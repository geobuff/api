package repo

type PlaysDto struct {
	QuizName string `json:"name"`
	Plays    int    `json:"plays"`
}

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

func GetTopFiveQuizPlays() ([]PlaysDto, error) {
	rows, err := Connection.Query("SELECT q.name, p.value FROM plays p JOIN quizzes q ON q.id = p.quizid ORDER BY value DESC LIMIT 5;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plays = []PlaysDto{}
	for rows.Next() {
		var play PlaysDto
		if err = rows.Scan(&play.QuizName, &play.Plays); err != nil {
			return nil, err
		}
		plays = append(plays, play)
	}
	return plays, rows.Err()
}

func GetLastFiveTriviaPlays() ([]PlaysDto, error) {
	rows, err := Connection.Query("SELECT name, plays FROM dailytrivia ORDER BY date DESC LIMIT 5;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plays = []PlaysDto{}
	for rows.Next() {
		var play PlaysDto
		if err = rows.Scan(&play.QuizName, &play.Plays); err != nil {
			return nil, err
		}
		plays = append(plays, play)
	}
	return plays, rows.Err()
}
