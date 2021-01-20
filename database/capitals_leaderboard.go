package database

// CapitalsLeaderboardEntry is the database object for a capitals leaderboard entry.
type CapitalsLeaderboardEntry struct {
	ID       int    `json:"id"`
	UserID   int    `json:"userId"`
	Country  string `json:"country"`
	Capitals int    `json:"capitals"`
	Time     int    `json:"time"`
}

// GetCapitalsLeaderboardEntries returns a page of leaderboard entries.
var GetCapitalsLeaderboardEntries = func(limit int, offset int) ([]CapitalsLeaderboardEntry, error) {
	rows, err := Connection.Query("SELECT * FROM capitals_leaderboard ORDER BY capitals DESC, time LIMIT $1 OFFSET $2;", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries = []CapitalsLeaderboardEntry{}
	for rows.Next() {
		var entry CapitalsLeaderboardEntry
		if err = rows.Scan(&entry.ID, &entry.UserID, &entry.Country, &entry.Capitals, &entry.Time); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

// GetCapitalsLeaderboardEntryID returns the first ID for a given page.
var GetCapitalsLeaderboardEntryID = func(limit int, offset int) (int, error) {
	statement := "SELECT id FROM capitals_leaderboard ORDER BY capitals DESC, time LIMIT $1 OFFSET $2;"
	var id int
	err := Connection.QueryRow(statement, limit, offset).Scan(&id)
	return id, err
}

// GetCapitalsLeaderboardEntry returns the leaderboard entry with a given id.
var GetCapitalsLeaderboardEntry = func(userID int) (CapitalsLeaderboardEntry, error) {
	statement := "SELECT * FROM capitals_leaderboard WHERE userId = $1;"
	var entry CapitalsLeaderboardEntry
	err := Connection.QueryRow(statement, userID).Scan(&entry.ID, &entry.UserID, &entry.Country, &entry.Capitals, &entry.Time)
	return entry, err
}

// InsertCapitalsLeaderboardEntry inserts a new leaderboard entry into the capitals_leaderboard table.
var InsertCapitalsLeaderboardEntry = func(entry CapitalsLeaderboardEntry) (int, error) {
	statement := "INSERT INTO capitals_leaderboard (userId, country, capitals, time) VALUES ($1, $2, $3, $4) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, entry.UserID, entry.Country, entry.Capitals, entry.Time).Scan(&id)
	return id, err
}

// UpdateCapitalsLeaderboardEntry updates an existing leaderboard entry.
var UpdateCapitalsLeaderboardEntry = func(entry CapitalsLeaderboardEntry) error {
	statement := "UPDATE capitals_leaderboard set userId = $2, country = $3, capitals = $4, time = $5 where id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, entry.ID, entry.UserID, entry.Country, entry.Capitals, entry.Time).Scan(&id)
}

// DeleteCapitalsLeaderboardEntry deletes a leaderboard entry.
var DeleteCapitalsLeaderboardEntry = func(entryID int) error {
	statement := "DELETE FROM capitals_leaderboard WHERE id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, entryID).Scan(&id)
}
