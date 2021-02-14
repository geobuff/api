package database

import (
	"fmt"
	"time"

	"github.com/geobuff/api/models"
)

// LeaderboardEntry is the database object for leaderboard entries.
type LeaderboardEntry struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userId"`
	CountryCode string    `json:"countryCode"`
	Score       int       `json:"score"`
	Time        int       `json:"time"`
	Added       time.Time `json:"added"`
}

const (
	// CountriesTable is the name of the countries leaderboard table in the database.
	CountriesTable = "countries_leaderboard"

	// CapitalsTable is the name of the countries leaderboard table in the database.
	CapitalsTable = "capitals_leaderboard"
)

// GetLeaderboardEntries returns a page of leaderboard entries.
var GetLeaderboardEntries = func(table string, filterParams models.GetEntriesFilterParams) ([]LeaderboardEntry, error) {
	query := fmt.Sprintf("SELECT l.id, userid, countrycode, score, time, added FROM %s l %s ORDER BY score DESC, time LIMIT $1 OFFSET $2;", table, getFilterSection(filterParams))

	rows, err := Connection.Query(query, filterParams.Limit, filterParams.Page*filterParams.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries = []LeaderboardEntry{}
	for rows.Next() {
		var entry LeaderboardEntry
		if err = rows.Scan(&entry.ID, &entry.UserID, &entry.CountryCode, &entry.Score, &entry.Time, &entry.Added); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func getFilterSection(filterParams models.GetEntriesFilterParams) string {
	if filterParams.User == "" {
		switch filterParams.Range {
		case "day":
			date := time.Now().AddDate(0, 0, -1)
			return fmt.Sprintf("WHERE added > '%v-%v-%v' ", date.Year(), date.Month(), date.Day())
		case "week":
			date := time.Now().AddDate(0, 0, -8)
			return fmt.Sprintf("WHERE added > '%v-%v-%v' ", date.Year(), date.Month(), date.Day())
		}
	}

	result := fmt.Sprintf("JOIN users u on u.id = l.userid WHERE u.username = '%s' ", filterParams.User)
	switch filterParams.Range {
	case "day":
		date := time.Now().AddDate(0, 0, -1)
		return result + fmt.Sprintf("AND added > '%v-%v-%v' ", date.Year(), date.Month(), date.Day())
	case "week":
		date := time.Now().AddDate(0, 0, -8)
		return result + fmt.Sprintf("AND added > '%v-%v-%v' ", date.Year(), date.Month(), date.Day())
	}
	return result
}

// GetLeaderboardEntryID returns the first ID for a given page.
var GetLeaderboardEntryID = func(table string, filterParams models.GetEntriesFilterParams) (int, error) {
	statement := fmt.Sprintf("SELECT l.id FROM %s l %s ORDER BY score DESC, time LIMIT 1 OFFSET $1;", table, getFilterSection(filterParams))
	var id int
	err := Connection.QueryRow(statement, (filterParams.Page+1)*filterParams.Limit).Scan(&id)
	return id, err
}

// GetLeaderboardEntry returns the leaderboard entry with a given id.
var GetLeaderboardEntry = func(table string, userID int) (LeaderboardEntry, error) {
	statement := fmt.Sprintf("SELECT * FROM %s WHERE userId = $1;", table)
	var entry LeaderboardEntry
	err := Connection.QueryRow(statement, userID).Scan(&entry.ID, &entry.UserID, &entry.CountryCode, &entry.Score, &entry.Time)
	return entry, err
}

// InsertLeaderboardEntry inserts a new leaderboard entry into the countries_leaderboard table.
var InsertLeaderboardEntry = func(table string, entry LeaderboardEntry) (int, error) {
	statement := fmt.Sprintf("INSERT INTO %s (userId, countryCode, score, time) VALUES ($1, $2, $3, $4) RETURNING id;", table)
	var id int
	err := Connection.QueryRow(statement, entry.UserID, entry.CountryCode, entry.Score, entry.Time).Scan(&id)
	return id, err
}

// UpdateLeaderboardEntry updates an existing leaderboard entry.
var UpdateLeaderboardEntry = func(table string, entry LeaderboardEntry) error {
	statement := fmt.Sprintf("UPDATE %s set userId = $2, countryCode = $3, score = $4, time = $5 where id = $1 RETURNING id;", table)
	var id int
	return Connection.QueryRow(statement, entry.ID, entry.UserID, entry.CountryCode, entry.Score, entry.Time).Scan(&id)
}

// DeleteLeaderboardEntry deletes a leaderboard entry.
var DeleteLeaderboardEntry = func(table string, entryID int) error {
	statement := fmt.Sprintf("DELETE FROM %s WHERE id = $1 RETURNING id;", table)
	var id int
	return Connection.QueryRow(statement, entryID).Scan(&id)
}
