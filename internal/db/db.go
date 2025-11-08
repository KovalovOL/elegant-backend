package db


import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"	
	"os"
	"fmt"
)


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

	if err = createUserTable(db); err != nil {
		return nil, err
	}
	if err = createCVTable(db); err != nil {
		return nil, err
	}
	if err = createTagTable(db); err != nil {
		return nil, err
	}
	if err = createTestTable(db); err != nil {
		return nil, err
	}
	if err = createTestTagTable(db); err != nil {
		return nil, err
	}

	return db, nil
}