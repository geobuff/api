package repo

import "database/sql"

func GetLastWeekTriviaPlays() ([]PlaysDto, error) {
	rows, err := Connection.Query("SELECT q.name, p.plays FROM triviaplays p JOIN trivia q ON q.id = p.triviaid ORDER BY q.date DESC LIMIT 7;")
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

	// Reverse slice for correct order to display in line graph.
	for i, j := 0, len(plays)-1; i < j; i, j = i+1, j-1 {
		plays[i], plays[j] = plays[j], plays[i]
	}

	return plays, rows.Err()
}

func IncrementTriviaPlays(triviaId int) error {
	var id int
	statement := "SELECT id FROM triviaplays WHERE triviaId = $1;"
	err := Connection.QueryRow(statement, triviaId).Scan(&id)

	if err == sql.ErrNoRows {
		statement = "INSERT INTO triviaplays (triviaId, plays) VALUES ($1, $2) RETURNING id;"
		return Connection.QueryRow(statement, triviaId, 1).Scan(&id)
	} else if err != nil {
		return err
	}

	statement = "UPDATE triviaplays set plays = plays + 1 WHERE id = $1 RETURNING id;"
	return Connection.QueryRow(statement, id).Scan(&id)
}

func DeleteTriviaPlays(triviaId int) error {
	statement := "DELETE FROM triviaplays WHERE triviaid = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, triviaId).Scan(&id)
}
