package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
)

func createUserTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR NOT NULL,
			email VARCHAR UNIQUE NOT NULL,
			github_url VARCHAR,
			linkedin_url VARCHAR,
			bio TEXT
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}
	return  nil
}

func createCVTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS cvs (
			cv_id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR NOT NULL,
			position VARCHAR NOT NULL,
			summary TEXT,
			skills TEXT,
			experience TEXT,
			education TEXT
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}
	return  nil
}