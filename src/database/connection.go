package database

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
)

func Connect() (*sql.DB, error) {
	cfg := mysql.Config{
		User:   "test",
		Passwd: "test",
		Addr:   "localhost",
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		log.Fatal(err)
		return db, err
	}

	if err := db.Ping(); err != nil {
		log.Print("Connected")
	}

	return db, nil
}
