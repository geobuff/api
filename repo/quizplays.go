package repo

import (
	"database/sql"
	"strings"
)

type PlaysDto struct {
	QuizName string `json:"name"`
	Plays    int    `json:"plays"`
}

var GetAllQuizPlays = func() (int, error) {
	var quizPlays int
	err := Connection.QueryRow("SELECT SUM(plays) from quizplays;").Scan(&quizPlays)
	if err != nil && strings.Contains(err.Error(), "sql: Scan error on column index 0") {
		quizPlays = 0
	} else if err != nil {
		return 0, err
	}

	var triviaPlays int
	err = Connection.QueryRow("SELECT SUM(plays) from triviaplays;").Scan(&triviaPlays)
	if err != nil && strings.Contains(err.Error(), "sql: Scan error on column index 0") {
		triviaPlays = 0
	} else if err != nil {
		return 0, err
	}

	var communityQuizPlays int
	err = Connection.QueryRow("SELECT SUM(plays) from communityquizplays;").Scan(&communityQuizPlays)
	if err != nil && strings.Contains(err.Error(), "sql: Scan error on column index 0") {
		communityQuizPlays = 0
	} else if err != nil {
		return 0, err
	}

	return quizPlays + triviaPlays + communityQuizPlays, nil
}

var GetQuizPlayCount = func(quizID int) (int, error) {
	statement := "SELECT plays from quizplays WHERE quizId = $1;"
	var plays int
	err := Connection.QueryRow(statement, quizID).Scan(&plays)
	return plays, err
}

var IncrementQuizPlayCount = func(quizID int) error {
	var id int
	statement := "SELECT id FROM quizplays WHERE quizId = $1;"
	err := Connection.QueryRow(statement, quizID).Scan(&id)

	if err == sql.ErrNoRows {
		statement = "INSERT INTO quizplays (quizId, plays) VALUES ($1, $2) RETURNING id;"
		return Connection.QueryRow(statement, quizID, 1).Scan(&id)
	} else if err != nil {
		return err
	}

	statement = "UPDATE quizplays set plays = plays + 1 WHERE id = $1 RETURNING id;"
	return Connection.QueryRow(statement, id).Scan(&id)
}

func GetTopFiveQuizPlays() ([]PlaysDto, error) {
	rows, err := Connection.Query("SELECT q.name, p.plays FROM quizplays p JOIN quizzes q ON q.id = p.quizid ORDER BY plays DESC LIMIT 5;")
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
