package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/supachai-sukd/assessment/pkg/config"
	"log"
	"sync"
)

var onceDb sync.Once

var db *sql.DB

func GetInstance() *sql.DB {
	onceDb.Do(func() {
		databaseConfig := config.DatabaseNew().(*config.DatabaseConfig)
		//db, err := pq.Open(fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		//	databaseConfig.Psql.DbHost,
		//	databaseConfig.Psql.DbPort,
		//	databaseConfig.Psql.DbUsername,
		//	databaseConfig.Psql.DbDatabase,
		//	databaseConfig.Psql.DbPassword,
		//))
		instance, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
			databaseConfig.Psql.DbHost,
			databaseConfig.Psql.DbPort,
			databaseConfig.Psql.DbUsername,
			databaseConfig.Psql.DbDatabase,
			databaseConfig.Psql.DbPassword,
		))

		if err != nil {
			log.Fatalf("Could not connect to database :%v", err)
		}
		db = instance
	})

	return db
}
