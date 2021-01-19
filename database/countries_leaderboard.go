package database

// CountriesLeaderboardEntry is the database object for a countries leaderboard entry.
type CountriesLeaderboardEntry struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userId"`
	Country   string `json:"country"`
	Countries int    `json:"countries"`
	Time      int    `json:"time"`
}

// GetCountriesLeaderboardEntries returns a page of leaderboard entries.
var GetCountriesLeaderboardEntries = func(limit int, offset int) ([]CountriesLeaderboardEntry, error) {
	rows, err := Connection.Query("SELECT * FROM countries_leaderboard ORDER BY countries DESC, time LIMIT $1 OFFSET $2;", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries = []CountriesLeaderboardEntry{}
	for rows.Next() {
		var entry CountriesLeaderboardEntry
		if err = rows.Scan(&entry.ID, &entry.UserID, &entry.Country, &entry.Countries, &entry.Time); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

// GetCountriesLeaderboardEntryID returns the first ID for a given page.
var GetCountriesLeaderboardEntryID = func(limit int, offset int) (int, error) {
	statement := "SELECT id FROM countries_leaderboard ORDER BY countries DESC, time LIMIT $1 OFFSET $2;"
	var id int
	err := Connection.QueryRow(statement, limit, offset).Scan(&id)
	return id, err
}

// GetCountriesLeaderboardEntry returns the leaderboard entry with a given id.
var GetCountriesLeaderboardEntry = func(userID int) (CountriesLeaderboardEntry, error) {
	statement := "SELECT * FROM countries_leaderboard WHERE userId = $1;"
	var entry CountriesLeaderboardEntry
	err := Connection.QueryRow(statement, userID).Scan(&entry.ID, &entry.UserID, &entry.Country, &entry.Countries, &entry.Time)
	return entry, err
}

// InsertCountriesLeaderboardEntry inserts a new leaderboard entry into the countries_leaderboard table.
var InsertCountriesLeaderboardEntry = func(entry CountriesLeaderboardEntry) (int, error) {
	statement := "INSERT INTO countries_leaderboard (userId, country, countries, time) VALUES ($1, $2, $3, $4) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, entry.UserID, entry.Country, entry.Countries, entry.Time).Scan(&id)
	return id, err
}

// UpdateCountriesLeaderboardEntry updates an existing leaderboard entry.
var UpdateCountriesLeaderboardEntry = func(entry CountriesLeaderboardEntry) error {
	statement := "UPDATE countries_leaderboard set userId = $2, country = $3, countries = $4, time = $5 where id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, entry.ID, entry.UserID, entry.Country, entry.Countries, entry.Time).Scan(&id)
}

// DeleteCountriesLeaderboardEntry deletes a leaderboard entry.
var DeleteCountriesLeaderboardEntry = func(entryID int) error {
	statement := "DELETE FROM countries_leaderboard WHERE id = $1 RETURNING id;"
	var id int
	return Connection.QueryRow(statement, entryID).Scan(&id)
}
