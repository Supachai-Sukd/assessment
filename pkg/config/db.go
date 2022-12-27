package config

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {

	var err error

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://hgqvnwpr:F50_9sky10ii2OVedWnRhdJWvm66iSW7@tiny.db.elephantsql.com/hgqvnwpr?sslmode=disable"
	}

	DB, err = sql.Open("postgres", databaseURL)

	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	createTb := `CREATE TABLE IF NOT EXISTS expenses (id SERIAL PRIMARY KEY,title TEXT,amount FLOAT,note TEXT,tags TEXT[]);`
	_, err = DB.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}

	return nil, nil
}
