package repo

import "database/sql"

func GetLastFiveTriviaPlays() ([]PlaysDto, error) {
	rows, err := Connection.Query("SELECT q.name, p.plays FROM dailytriviaplays p JOIN on dailytrivia q ORDER BY q.date DESC LIMIT 5;")
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

func IncrementTriviaPlays(triviaId int) error {
	var id int
	statement := "SELECT id FROM dailytriviaplays WHERE triviaId = $1;"
	err := Connection.QueryRow(statement, triviaId).Scan(&id)

	if err == sql.ErrNoRows {
		statement = "INSERT INTO dailytriviaplays (triviaId, plays) VALUES ($1, $2) RETURNING id;"
		return Connection.QueryRow(statement, triviaId, 1).Scan(&id)
	} else if err != nil {
		return err
	}

	statement = "UPDATE dailytriviaplays set plays = plays + 1 WHERE id = $1 RETURNING id;"
	return Connection.QueryRow(statement, id).Scan(&id)
}
