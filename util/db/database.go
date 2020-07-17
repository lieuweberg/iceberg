package db

import (
	"database/sql"
	"log"

	// Import the DB driver
	_ "github.com/mattn/go-sqlite3"
)

// DB is the *sql.DB
var db *sql.DB

func init() {
	database, err := sql.Open("sqlite3", "./iceberg.db")
	if err != nil {
		log.Fatal(err)
	}
	
	err = database.Ping()
	if err != nil {
		log.Fatalf("Unable to connect to database: %s", err)
	}

	_, err = database.Exec(`
		CREATE TABLE IF NOT EXISTS users(
			id	STRING PRIMARY KEY	NOT NULL,
			mcName 			STRING	NOT NULL,
			birthday		STRING
		)
	`)
	if err != nil {
		log.Fatalf("Could not create table users: %s", err)
	}

	db = database;
}

// GetDB returns the *sql.DB
func GetDB() *sql.DB {
	return db
}

// User is the database entry of a user
type User struct {
	ID string
	McName string
	Birthday string
}

// GetUser fetches the database entry of id
func GetUser(id string) (user User, err error) {
	user = User{}
	err = db.QueryRow("SELECT * FROM users WHERE id=?", id).Scan(&user.ID, &user.McName, &user.Birthday)
	return
}