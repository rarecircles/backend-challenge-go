package sql

import (
	"os"

	"github.com/go-pg/pg/v10"
)

// TODO: Make middleware

func Connect() *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     "postgres:5432",
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB"),
	})

	return db
}
