package main

import (
	"database/sql"
	"log"

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
	_, err = db.Exec("SELECT 1")
	if err != nil {
		return err
	}

	return nil
}

func (z *Zusage) Put() (err error) {
	log.Println("Put", z)

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	old := Zusage{}
	err = tx.QueryRow("SELECT nick, kommt, kommentar FROM zusagen WHERE nick = $1", z.Nick).Scan(&old.Nick, &old.Kommt, &old.Kommentar)

	switch {
	case err == sql.ErrNoRows:
		_, err = tx.Exec("INSERT INTO zusagen (nick, kommt, kommentar, registered) VALUES ($1, $2, $3, NOW())", z.Nick, z.Kommt, z.Kommentar)
	case err != nil:
		return err
	default:
		_, err = tx.Exec("UPDATE zusagen SET kommt = $2, kommentar = $3 WHERE nick = $1 ", z.Nick, z.Kommt, z.Kommentar)
	}

	if err != nil {
		return err
	}

	return tx.Commit()
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
