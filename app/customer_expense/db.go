package customer_expense

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func InitDB() {
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		"localhost",
		"1112",
		"root",
		"golang",
		"root",
		"disable"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	createTb := `CREATE TABLE IF NOT EXISTS expenses (id SERIAL PRIMARY KEY,title TEXT,amount FLOAT,note TEXT,tags TEXT[]);`
	_, err = db.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}

}
