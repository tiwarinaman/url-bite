package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		short_id TEXT UNIQUE NOT NULL,
		original_url TEXT NOT NULL,
		created_at TEXT NOT NULL
	);`
	if _, err := db.Exec(createTableQuery); err != nil {
		return nil, fmt.Errorf("error creating table: %v", err)
	}

	return db, nil
}
