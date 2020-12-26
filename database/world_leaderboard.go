package database

// Entry is the database object for a leaderboard entry.
type Entry struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userId"`
	Country   string `json:"country"`
	Countries int    `json:"countries"`
	Time      int    `json:"time"`
}

// GetWorldLeaderboardEntries returns a page of leaderboard entries.
var GetWorldLeaderboardEntries = func(limit int, offset int) ([]Entry, error) {
	rows, err := Connection.Query("SELECT * FROM world_leaderboard ORDER BY countries DESC, time LIMIT $1 OFFSET $2;", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries = []Entry{}
	for rows.Next() {
		var entry Entry
		if err = rows.Scan(&entry.ID, &entry.UserID, &entry.Country, &entry.Countries, &entry.Time); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

// GetWorldLeaderboardEntryID returns the first ID for a given page.
var GetWorldLeaderboardEntryID = func(limit int, offset int) (int, error) {
	statement := "SELECT id FROM world_leaderboard ORDER BY countries DESC, time LIMIT $1 OFFSET $2;"
	var id int
	err := Connection.QueryRow(statement, limit, offset).Scan(&id)
	return id, err
}

// GetWorldLeaderboardEntry returns the leaderboard entry with a given id.
func GetWorldLeaderboardEntry(id int) (Entry, error) {
	statement := "SELECT * FROM world_leaderboard WHERE id = $1;"
	var entry Entry
	err := Connection.QueryRow(statement, id).Scan(&entry.ID, &entry.UserID, &entry.Country, &entry.Countries, &entry.Time)
	return entry, err
}

// InsertWorldLeaderboardEntry inserts a new leaderboard entry into the world_leaderboard table.
func InsertWorldLeaderboardEntry(entry Entry) (int, error) {
	statement := "INSERT INTO world_leaderboard (userId, country, countries, time) VALUES ($1, $2, $3, $4) RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, entry.UserID, entry.Country, entry.Countries, entry.Time).Scan(&id)
	return id, err
}

// UpdateWorldLeaderboardEntry updates an existing leaderboard entry.
func UpdateWorldLeaderboardEntry(entryID int, entry Entry) (int, error) {
	statement := "UPDATE world_leaderboard set userId = $2, country = $3, countries = $4, time = $5 where id = $1 RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, entryID, entry.UserID, entry.Country, entry.Countries, entry.Time).Scan(&id)
	return id, err
}

// DeleteWorldLeaderboardEntry deletes a leaderboard entry.
func DeleteWorldLeaderboardEntry(entryID int) (int, error) {
	statement := "DELETE FROM world_leaderboard WHERE id = $1 RETURNING id;"
	var id int
	err := Connection.QueryRow(statement, entryID).Scan(&id)
	return id, err
}
