package repo

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"jobsity-code-challenge/config"
)

type RepoDB struct {
	db *sql.DB
}

func SetupDB(database config.Database) *RepoDB {
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		database.Host,
		database.Port,
		database.User,
		database.Password,
		database.Name,
	)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		panic(fmt.Sprintf("Fatal error creating DB connection: %v \n", err))
	}
	return &RepoDB{db}
}

func (d RepoDB) Close() error {
	return d.db.Close()
}
