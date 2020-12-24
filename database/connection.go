package database

import (
	"database/sql"
	"fmt"

	"github.com/geobuff/geobuff-api/config"
)

// Connection is the connection handle for the database.
var Connection *sql.DB

// OpenConnection initializes the Connection variable.
var OpenConnection = func() error {
	Connection, err := sql.Open("postgres", getConnectionString())
	if err != nil {
		return err
	}
	return Connection.Ping()
}

func getConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Values.Database.Host,
		config.Values.Database.Port, config.Values.Database.User, config.Values.Database.Password,
		config.Values.Database.Name)
}
