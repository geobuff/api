package keys

import (
	"github.com/geobuff/geobuff-api/database"
)

// Key is the database object for a key entry.
type Key struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

// ValidAuth0Key checks whether a key matches the Auth0 key stored in the database.
func ValidAuth0Key(key string) bool {
	statement := "SELECT * FROM keys WHERE name = $1;"
	row := database.DBConnection.QueryRow(statement, "auth0")

	var entry Key
	if err := row.Scan(&entry.ID, &entry.Name, &entry.Key); err != nil {
		panic(err)
	}

	return entry.Key == key
}
