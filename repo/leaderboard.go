package repo

import (
	"database/sql"
	"fmt"
	"time"
)

type LeaderboardEntry struct {
	ID     int       `json:"id"`
	QuizID int       `json:"quizId"`
	UserID int       `json:"userId"`
	Score  int       `json:"score"`
	Time   int       `json:"time"`
	Added  time.Time `json:"added"`
}

type LeaderboardEntryDto struct {
	ID          int       `json:"id"`
	QuizID      int       `json:"quizId"`
	UserID      int       `json:"userId"`
	Username    string    `json:"username"`
	QuizName    string    `json:"quizName"`
	CountryCode string    `json:"countryCode"`
	Score       int       `json:"score"`
	Time        int       `json:"time"`
	Added       time.Time `json:"added"`
	Rank        int       `json:"rank"`
}

type GetLeaderboardEntriesFilterParams struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Range string `json:"range"`
	User  string `json:"user"`
	Rank  int    `json:"rank"`
}

type UserLeaderboardEntryDto struct {
	ID           int       `json:"id"`
	UserID       int       `json:"userId"`
	QuizID       int       `json:"quizId"`
	BadgeID      int       `json:"badgeId"`
	QuizName     string    `json:"quizName"`
	QuizImageUrl string    `json:"quizImageUrl"`
	Score        int       `json:"score"`
	Time         int       `json:"time"`
	Added        time.Time `json:"added"`
	Rank         int       `json:"rank"`
}

var GetLeaderboardEntries = func(quizID int, filterParams GetLeaderboardEntriesFilterParams) ([]LeaderboardEntryDto, error) {
	var rows *sql.Rows
	var err error
	if filterParams.Rank != 0 {
		query := "SELECT * FROM (SELECT l.id, l.quizid, l.userid, u.username, u.countrycode, l.score, l.time, l.added, RANK () OVER (PARTITION BY l.quizid ORDER BY score desc, l.time) rank FROM leaderboard l JOIN users u on u.id = l.userid) a WHERE quizid = $1 AND username ILIKE '%' || $2 || '%' " + getRangeFilter(filterParams.Range) + " AND rank BETWEEN $3 AND $4 ORDER BY score DESC, time"
		lower := filterParams.Rank - (filterParams.Rank % 10)
		upper := lower + filterParams.Limit
		rows, err = Connection.Query(query, quizID, filterParams.User, lower+1, upper)
	} else {
		query := "SELECT * FROM (SELECT l.id, l.quizid, l.userid, u.username, u.countrycode, l.score, l.time, l.added, RANK () OVER (PARTITION BY l.quizid ORDER BY score desc, l.time) rank FROM leaderboard l JOIN users u on u.id = l.userid) a WHERE quizid = $1 AND username ILIKE '%' || $2 || '%' " + getRangeFilter(filterParams.Range) + " ORDER BY score DESC, time LIMIT $3 OFFSET $4;"
		rows, err = Connection.Query(query, quizID, filterParams.User, filterParams.Limit, filterParams.Page*filterParams.Limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries = []LeaderboardEntryDto{}
	for rows.Next() {
		var entry LeaderboardEntryDto
		if err = rows.Scan(&entry.ID, &entry.QuizID, &entry.UserID, &entry.Username, &entry.CountryCode, &entry.Score, &entry.Time, &entry.Added, &entry.Rank); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func getRangeFilter(rangeFilter string) string {
	switch rangeFilter {
	case "day":
		date := time.Now().AddDate(0, 0, -1)
		return fmt.Sprintf("AND added > '%v-%v-%v' ", date.Year(), date.Month(), date.Day())
	case "week":
		date := time.Now().AddDate(0, 0, -8)
		return fmt.Sprintf("AND added > '%v-%v-%v' ", date.Year(), date.Month(), date.Day())
	default:
		return ""
	}
}

var GetLeaderboardEntryID = func(quizID int, filterParams GetLeaderboardEntriesFilterParams) (int, error) {
	query := "SELECT l.id FROM leaderboard l JOIN users u on u.id = l.userid WHERE l.quizid = $1 AND u.username ILIKE '%' || $2 || '%' " + getRangeFilter(filterParams.Range) + " ORDER BY score DESC, time LIMIT 1 OFFSET $3;"
	var id int
	err := Connection.QueryRow(query, quizID, filterParams.User, (filterParams.Page+1)*filterParams.Limit).Scan(&id)
	return id, err
}

var GetUserLeaderboardEntries = func(userID int) ([]UserLeaderboardEntryDto, error) {
	query := "SELECT * from (SELECT l.id, l.userid, l.quizid, q.badgeId, q.name, q.imageUrl, l.score, l.time, l.added, RANK () OVER (PARTITION BY l.quizid ORDER BY score desc, l.time) rank FROM leaderboard l JOIN quizzes q on q.id = l.quizid) c WHERE c.userid = $1;"

	rows, err := Connection.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries = []UserLeaderboardEntryDto{}
	for rows.Next() {
		var entry UserLeaderboardEntryDto
		if err = rows.Scan(&entry.ID, &entry.UserID, &entry.QuizID, &entry.BadgeID, &entry.QuizName, &entry.QuizImageUrl, &entry.Score, &entry.Time, &entry.Added, &entry.Rank); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

var GetLeaderboardEntry = func(quizID, userID int) (LeaderboardEntryDto, error) {
	statement := "SELECT * from (SELECT l.id, l.quizid, l.userid, u.username, q.name, u.countrycode, l.score, l.time, l.added, RANK () OVER (PARTITION BY l.quizid ORDER BY score desc, l.time) rank FROM leaderboard l JOIN users u on u.id = l.userId JOIN quizzes q on q.id = l.quizid WHERE l.quizid = $1) c WHERE c.userid = $2;"
	var entry LeaderboardEntryDto
	err := Connection.QueryRow(statement, quizID, userID).Scan(&entry.ID, &entry.QuizID, &entry.UserID, &entry.Username, &entry.QuizName, &entry.CountryCode, &entry.Score, &entry.Time, &entry.Added, &entry.Rank)
	return entry, err
}

var GetLeaderboardEntryById = func(id int) (LeaderboardEntry, error) {
	statement := "SELECT * from leaderboard WHERE id = $1;"
	var entry LeaderboardEntry
	err := Connection.QueryRow(statement, id).Scan(&entry.ID, &entry.QuizID, &entry.UserID, &entry.Score, &entry.Time, &entry.Added)
	return entry, err
}

var InsertLeaderboardEntry = func(entry LeaderboardEntry) (int, error) {
	statement := "INSERT INTO leaderboard (quizId, userId, score, time, added) VALUES ($1, $2, $3, $4, $5) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, entry.QuizID, entry.UserID, entry.Score, entry.Time, entry.Added).Scan(&id)
	return id, err
}

var UpdateLeaderboardEntry = func(entry LeaderboardEntry) error {
	statement := "UPDATE leaderboard set quizId = $2, userId = $3, score = $4, time = $5, added = $6 where id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, entry.ID, entry.QuizID, entry.UserID, entry.Score, entry.Time, entry.Added).Scan(&id)
}

var DeleteLeaderboardEntry = func(entryID int) error {
	statement := "DELETE FROM leaderboard WHERE id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, entryID).Scan(&id)
}
