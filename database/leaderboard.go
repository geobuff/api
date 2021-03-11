package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/geobuff/api/models"
)

// LeaderboardEntry is the database object for leaderboard entries.
type LeaderboardEntry struct {
	ID     int       `json:"id"`
	UserID int       `json:"userId"`
	Score  int       `json:"score"`
	Time   int       `json:"time"`
	Added  time.Time `json:"added"`
}

// LeaderboardEntryDto contains additional fields for display in the leaderboard table.
type LeaderboardEntryDto struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userId"`
	Username    string    `json:"username"`
	QuizID      int       `json:"quizId"`
	QuizName    string    `json:"quizName"`
	CountryCode string    `json:"countryCode"`
	Score       int       `json:"score"`
	Time        int       `json:"time"`
	Added       time.Time `json:"added"`
	Ranking     int       `json:"ranking"`
}

const (
	// CountriesTable is the name of the countries leaderboard table in the database.
	CountriesTable = "countries_leaderboard"

	// CapitalsTable is the name of the countries leaderboard table in the database.
	CapitalsTable = "capitals_leaderboard"
)

// GetLeaderboardEntries returns a page of leaderboard entries.
var GetLeaderboardEntries = func(table string, filterParams models.GetEntriesFilterParams) ([]LeaderboardEntryDto, error) {
	query := "SELECT l.id, userid, u.username, u.countrycode, score, time FROM " + table + " l JOIN users u on u.id = l.userid WHERE u.username ILIKE '%' || $1 || '%' " + getRangeFilter(filterParams.Range) + " ORDER BY score DESC, time LIMIT $2 OFFSET $3;"

	rows, err := Connection.Query(query, filterParams.User, filterParams.Limit, filterParams.Page*filterParams.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries = []LeaderboardEntryDto{}
	for rows.Next() {
		var entry LeaderboardEntryDto
		if err = rows.Scan(&entry.ID, &entry.UserID, &entry.Username, &entry.CountryCode, &entry.Score, &entry.Time); err != nil {
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

// GetLeaderboardEntryID returns the first ID for a given page.
var GetLeaderboardEntryID = func(table string, filterParams models.GetEntriesFilterParams) (int, error) {
	query := "SELECT l.id FROM " + table + " l JOIN users u on u.id = l.userid WHERE u.username ILIKE '%' || $1 || '%' " + getRangeFilter(filterParams.Range) + " ORDER BY score DESC, time LIMIT 1 OFFSET $2;"
	var id int
	err := Connection.QueryRow(query, filterParams.User, (filterParams.Page+1)*filterParams.Limit).Scan(&id)
	return id, err
}

// GetLeaderboardEntry returns the leaderboard entry with a given id.
var GetLeaderboardEntry = func(table string, userID int) (LeaderboardEntryDto, error) {
	statement := fmt.Sprintf("SELECT * from (SELECT l.id, userid, u.username, q.id, q.name, u.countrycode, score, l.time, added, RANK () OVER (ORDER BY score desc, l.time) rank FROM %s l JOIN users u on u.id = l.userId JOIN quizzes q on q.id = $1) c WHERE c.userid = $2;", table)
	var entry LeaderboardEntryDto
	name := strings.Split(table, "_")[0]
	quizID, err := GetQuizID(name)
	if err != nil {
		return entry, err
	}

	err = Connection.QueryRow(statement, quizID, userID).Scan(&entry.ID, &entry.UserID, &entry.Username, &entry.QuizID, &entry.QuizName, &entry.CountryCode, &entry.Score, &entry.Time, &entry.Added, &entry.Ranking)
	return entry, err
}

// InsertLeaderboardEntry inserts a new leaderboard entry into the countries_leaderboard table.
var InsertLeaderboardEntry = func(table string, entry LeaderboardEntry) (int, error) {
	statement := fmt.Sprintf("INSERT INTO %s (userId, score, time, added) VALUES ($1, $2, $3, $4) RETURNING id;", table)
	var id int
	err := Connection.QueryRow(statement, entry.UserID, entry.Score, entry.Time, entry.Added).Scan(&id)
	return id, err
}

// UpdateLeaderboardEntry updates an existing leaderboard entry.
var UpdateLeaderboardEntry = func(table string, entry LeaderboardEntry) error {
	statement := fmt.Sprintf("UPDATE %s set userId = $2, score = $3, time = $4, added = $5 where id = $1 RETURNING id;", table)
	var id int
	return Connection.QueryRow(statement, entry.ID, entry.UserID, entry.Score, entry.Time, entry.Added).Scan(&id)
}

// DeleteLeaderboardEntry deletes a leaderboard entry.
var DeleteLeaderboardEntry = func(table string, entryID int) error {
	statement := fmt.Sprintf("DELETE FROM %s WHERE id = $1 RETURNING id;", table)
	var id int
	return Connection.QueryRow(statement, entryID).Scan(&id)
}
