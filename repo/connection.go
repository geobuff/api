package repo

import (
	"database/sql"
	"os"
)

// Connection is the connection handle for the database.
var Connection *sql.DB

// OpenConnection initializes the Connection variable.
var OpenConnection = func() error {
	connection, err := sql.Open("postgres", os.Getenv("CONNECTION_STRING"))
	if err != nil {
		return err
	}

	Connection = connection
	return Connection.Ping()
}
