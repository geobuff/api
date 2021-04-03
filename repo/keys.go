package repo

// Key is the database object for a key entry.
type Key struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

// GetKey returns a key with a given name.
var GetKey = func(name string) (Key, error) {
	statement := "SELECT * FROM keys WHERE name = $1;"
	var key Key
	err := Connection.QueryRow(statement, name).Scan(&key.ID, &key.Name, &key.Key)
	return key, err
}
