package repo

import (
	"database/sql"
	"fmt"
	"os"
)

// Connection is the connection handle for the database.
var Connection *sql.DB

// OpenConnection initializes the Connection variable.
var OpenConnection = func() error {
	connection, err := sql.Open("postgres", getConnectionString())
	if err != nil {
		return err
	}

	Connection = connection
	return Connection.Ping()
}

func getConnectionString() string {
	return fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_NAME"))
}
