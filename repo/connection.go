package repo

import (
	"database/sql"
	"os"
)

var Connection *sql.DB

var OpenConnection = func() error {
	connection, err := sql.Open("postgres", os.Getenv("CONNECTION_STRING"))
	if err != nil {
		return err
	}

	Connection = connection
	return Connection.Ping()
}
