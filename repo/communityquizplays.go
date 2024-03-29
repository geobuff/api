package repo

import "database/sql"

type CommunityQuizPlay struct {
	ID              int           `json:"id"`
	CommunityQuizID sql.NullInt64 `json:"communityQuizId"`
	Plays           int           `json:"plays"`
}

func IncrementCommunityQuizPlays(communityQuizID int) error {
	var id int
	statement := "SELECT id FROM communityquizplays WHERE communityQuizId = $1;"
	err := Connection.QueryRow(statement, communityQuizID).Scan(&id)

	if err == sql.ErrNoRows {
		statement = "INSERT INTO communityquizplays (communityQuizId, plays) VALUES ($1, $2) RETURNING id;"
		return Connection.QueryRow(statement, communityQuizID, 1).Scan(&id)
	} else if err != nil {
		return err
	}

	statement = "UPDATE communityquizplays set plays = plays + 1 WHERE id = $1 RETURNING id;"
	return Connection.QueryRow(statement, id).Scan(&id)
}

func DeleteCommunityQuizPlay(communityQuizID int) error {
	var id int
	statement := "DELETE FROM communityquizplays WHERE communityQuizId = $1 RETURNING id;"
	return Connection.QueryRow(statement, communityQuizID).Scan(&id)
}

func ClearCommunityQuizPlayCommunityQuizId(communityQuizID int) error {
	var id int
	statement := "UPDATE communityquizplays set communityquizid = null WHERE communityquizid = $1 RETURNING id;"
	return Connection.QueryRow(statement, communityQuizID).Scan(&id)
}
