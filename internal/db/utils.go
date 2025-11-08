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
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
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
			title VARCHAR(255) NOT NULL,
			position VARCHAR(255) NOT NULL,
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

func createTagTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tags (
			tag_id SERIAL PRIMARY KEY,
			name VARCHAR UNIQUE NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}
	return  nil
}

// func createTestTable(db *sql.DB) error {
// 	_, err := db.Exec(`
// 		CREATE TABLE IF NOT EXISTS tests (
// 			test_id SERIAL PRIMARY KEY,	
// 			title VARCHAR(255) UNIQUE NOT NULL,
// 			time_limit INTEGER,
// 			type VARCHAR NOT NULL
// 		)
// 	`)
// 	if err != nil {
// 		return fmt.Errorf("error creating users table: %w", err)
// 	}
// 	return  nil
// }

// func createTestTagTable(db *sql.DB) error {
// 	_, err := db.Exec(`
// 		CREATE TABLE IF NOT EXISTS test_tags (
// 			id SERIAL PRIMARY KEY,
// 			tag_id INTEGER REFERENCES tags(tag_id) ON DELETE CASCADE,
// 			test_id INTEGER REFERENCES tests(test_id) ON DELETE CASCADE,
// 			UNIQUE (test_id, tag_id)
// 		)
// 	`)
// 	if err != nil {
// 		return fmt.Errorf("error creating users table: %w", err)
// 	}
// 	return  nil
// }