package database

import "fmt"

// LeaderboardEntry is the database object for leaderboard entries.
type LeaderboardEntry struct {
	ID          int    `json:"id"`
	UserID      int    `json:"userId"`
	CountryCode string `json:"countryCode"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
}

const (
	// CountriesTable is the name of the countries leaderboard table in the database.
	CountriesTable = "countries_leaderboard"

	// CapitalsTable is the name of the countries leaderboard table in the database.
	CapitalsTable = "capitals_leaderboard"
)

// GetLeaderboardEntries returns a page of leaderboard entries.
var GetLeaderboardEntries = func(table string, limit, offset int) ([]LeaderboardEntry, error) {
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY score DESC, time LIMIT $1 OFFSET $2;", table)
	rows, err := Connection.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries = []LeaderboardEntry{}
	for rows.Next() {
		var entry LeaderboardEntry
		if err = rows.Scan(&entry.ID, &entry.UserID, &entry.CountryCode, &entry.Score, &entry.Time); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

// GetLeaderboardEntryID returns the first ID for a given page.
var GetLeaderboardEntryID = func(table string, limit, offset int) (int, error) {
	statement := fmt.Sprintf("SELECT id FROM %s ORDER BY score DESC, time LIMIT $1 OFFSET $2;", table)
	var id int
	err := Connection.QueryRow(statement, limit, offset).Scan(&id)
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
