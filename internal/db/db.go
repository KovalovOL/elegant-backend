package db


import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"	
	"os"
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

func ConnectDB() (*sql.DB, error) {
	dsn := fmt.Sprintf(
    	"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
    	os.Getenv("DB_USER"),
    	os.Getenv("DB_PASSWORD"),
    	os.Getenv("DB_NAME"),
    	os.Getenv("DB_HOST"),
    	os.Getenv("DB_PORT"),
	)

    db, err := sql.Open("postgres", dsn)
    if err != nil {
       log.Fatal(err)
	   return nil, err
    }	

	if err := db.Ping(); err != nil {
		return  nil, err
	}

	err = createUserTable(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}