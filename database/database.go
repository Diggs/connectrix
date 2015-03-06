package database

import (
	"database/sql"
	"github.com/diggs/connectrix/config"
	_ "github.com/lib/pq"
)

// database contains a reference to the database, populated via Connect.
var database *sql.DB

// GetDatabase returns a database reference that can be used to query postgres.
func GetDatabase() *sql.DB {
	return database
}

// Connect connects to the configured postgres database and verifies the connection.
func Connect() error {

	db, err := sql.Open("postgres", config.Get().DatabaseConnection)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	database = db
	return nil
}
