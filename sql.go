package main

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log"
	"time"
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

func (v *Vortrag) Put() (err error) {
	log.Println("Put", v)
	tx, _ := db.Begin()

	var stmt *sql.Stmt
	if v.Id < 0 {
		if v.Date.IsZero() {
			stmt, err = tx.Prepare("INSERT INTO vortraege (topic, abstract, speaker, infourl, password) VALUES ($1, $2, $3, $4, $5) RETURNING id")
			if err != nil {
				return err
			}

			err = stmt.QueryRow(v.Topic, v.Abstract, v.Speaker, v.InfoURL, v.Password).Scan(&v.Id)
		} else {
			stmt, err = tx.Prepare("INSERT INTO vortraege (topic, abstract, speaker, infourl, password, date) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id")
			if err != nil {
				return err
			}

			err = stmt.QueryRow(v.Topic, v.Abstract, v.Speaker, v.InfoURL, v.Password, time.Time(v.Date)).Scan(&v.Id)
		}
	} else {
		if v.Date.IsZero() {
			stmt, err = tx.Prepare("UPDATE vortraege SET topic = $1, abstract = $2, speaker = $3, infourl = $4 WHERE id = $5")
			if err != nil {
				return err
			}

			_, err = stmt.Exec(v.Topic, v.Abstract, v.Speaker, v.InfoURL, v.Id)
		} else {
			stmt, err = tx.Prepare("UPDATE vortraege SET topic = $1, abstract = $2, speaker = $3, infourl = $4, date = $5 WHERE id = $6")
			if err != nil {
				return err
			}

			_, err = stmt.Exec(v.Topic, v.Abstract, v.Speaker, v.InfoURL, time.Time(v.Date), v.Id)
		}
	}

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}

func Load(id int) (*Vortrag, error) {
	rows, err := db.Query("SELECT date, topic, abstract, speaker, infourl, password FROM vortraege WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("No such id")
	}

	vortrag := Vortrag{Id: id}
	var date *time.Time

	err = rows.Scan(&date, &vortrag.Topic, &vortrag.Abstract, &vortrag.Speaker, &vortrag.InfoURL, &vortrag.Password)
	if err != nil {
		return nil, err
	}
	if date != nil {
		vortrag.Date = CustomTime(*date)
	}

	if !vortrag.Date.IsZero() {
		vortrag.HasDate = true
	}

	return &vortrag, nil
}
