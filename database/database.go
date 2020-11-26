package database

import (
	"database/sql"
	"fmt"

	"github.com/geobuff/geobuff-api/config"
)

// DBConnection is the connection handle for the database.
var DBConnection *sql.DB

// GetConnectionString returns the database connection string.
func GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Values.Database.Host,
		config.Values.Database.Port, config.Values.Database.User, config.Values.Database.Password,
		config.Values.Database.Name)
}
