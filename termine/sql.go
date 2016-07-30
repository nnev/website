package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var (
	db *sql.DB
)

func OpenDB() (err error) {
	db, err = sql.Open(*driver, *connect)
	if err != nil {
		return err
	}
	return nil
}
