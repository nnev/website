package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var (
	db *sql.DB
)

func OpenDB() (err error) {
	db, err = sql.Open(*driver, *connect)
	if err != nil {
		return err
	}
	_, err = db.Exec("SELECT 1")
	if err != nil {
		return err
	}

	return nil
}

func (z *Zusage) Put() (err error) {
	log.Println("Put", z)
	_, _ = db.Exec("DELETE FROM zusagen WHERE nick = $1", z.Nick)
	_, err = db.Exec("INSERT INTO zusagen (nick, kommt, kommentar) VALUES ($1, $2, $3)", z.Nick, z.Kommt, z.Kommentar)
	if err != nil {
		return err
	}
	return nil
}

func GetZusage(nick string) (z Zusage) {
	err := db.QueryRow("SELECT kommt, kommentar FROM zusagen WHERE nick = $1", nick).Scan(&z.Kommt, &z.Kommentar)
	if err != nil {
		z.Nick = nick
		z.HasKommt = false
		return z
	}

	z.Nick = nick
	z.HasKommt = true
	return z
}
