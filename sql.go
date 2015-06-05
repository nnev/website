package main

import (
	"database/sql"
	"errors"
	"log"
	"time"

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

func (v *Vortrag) Put() (err error) {
	log.Println("Put", v)
	tx, _ := db.Begin()

	defer tx.Rollback()

	var stmt *sql.Stmt
	if v.Id < 0 {
		if v.Date.IsZero() {
			stmt, err = tx.Prepare("INSERT INTO vortraege (topic, abstract, speaker, password) VALUES ($1, $2, $3, $4) RETURNING id")
			if err != nil {
				return err
			}

			err = stmt.QueryRow(v.Topic, v.Abstract, v.Speaker, v.Password).Scan(&v.Id)
			if err != nil {
				return err
			}

			stmt, err = tx.Prepare("INSERT INTO vortrag_links (vortrag, kind, url) VALUES ($1, $2, $3)")
			if err != nil {
				return err
			}
			for _, link := range v.Links {
				_, err = stmt.Exec(v.Id, link.Kind, link.Url)
				if err != nil {
					return err
				}
			}
		} else {
			stmt, err = tx.Prepare("INSERT INTO vortraege (topic, abstract, speaker, password, date) VALUES ($1, $2, $3, $4, $5) RETURNING id")
			if err != nil {
				return err
			}

			err = stmt.QueryRow(v.Topic, v.Abstract, v.Speaker, v.Password, time.Time(v.Date)).Scan(&v.Id)
			if err != nil {
				return err
			}

			stmt, err = tx.Prepare("INSERT INTO vortrag_links (vortrag, kind, url) VALUES ($1, $2, $3)")
			if err != nil {
				return err
			}

			for _, link := range v.Links {
				_, err = stmt.Exec(v.Id, link.Kind, link.Url)
				if err != nil {
					return err
				}
			}
		}
	} else {
		if v.Date.IsZero() {
			stmt, err = tx.Prepare("UPDATE vortraege SET topic = $1, abstract = $2, speaker = $3, date = NULL WHERE id = $4")
			if err != nil {
				return err
			}

			_, err = stmt.Exec(v.Topic, v.Abstract, v.Speaker, v.Id)
			if err != nil {
				return err
			}

			_, err = tx.Exec("DELETE FROM vortrag_links WHERE vortrag = $1;", v.Id)
			if err != nil {
				return err
			}

			stmt, err = tx.Prepare("INSERT INTO vortrag_links (vortrag, kind, url) VALUES ($1, $2, $3);")
			if err != nil {
				return err
			}

			for _, link := range v.Links {
				_, err = stmt.Exec(v.Id, link.Kind, link.Url)
				if err != nil {
					return err
				}
			}
		} else {
			stmt, err = tx.Prepare("UPDATE vortraege SET topic = $1, abstract = $2, speaker = $3, date = $4 WHERE id = $5")
			if err != nil {
				return err
			}

			_, err = stmt.Exec(v.Topic, v.Abstract, v.Speaker, time.Time(v.Date), v.Id)
			if err != nil {
				return err
			}

			_, err = tx.Exec("DELETE FROM vortrag_links WHERE vortrag = $1;", v.Id)
			if err != nil {
				return err
			}

			stmt, err = tx.Prepare("INSERT INTO vortrag_links (vortrag, kind, url) VALUES ($1, $2, $3);")
			if err != nil {
				return err
			}

			for _, link := range v.Links {
				_, err = stmt.Exec(v.Id, link.Kind, link.Url)
				if err != nil {
					return err
				}
			}
		}
	}

	return tx.Commit()
}

func Load(id int) (*Vortrag, error) {
	rows, err := db.Query("SELECT date, topic, abstract, speaker, password FROM vortraege WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("No such id")
	}

	vortrag := Vortrag{Id: id}
	var date *time.Time

	err = rows.Scan(&date, &vortrag.Topic, &vortrag.Abstract, &vortrag.Speaker, &vortrag.Password)
	if err != nil {
		return nil, err
	}
	if date != nil {
		vortrag.Date = CustomTime(*date)
	}

	if !vortrag.Date.IsZero() {
		vortrag.HasDate = true
	}

	rows, err = db.Query("SELECT kind, url FROM vortrag_links WHERE vortrag = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var kind, url string

		if err = rows.Scan(&kind, &url); err != nil {
			return nil, err
		}
		vortrag.Links = append(vortrag.Links, Link{Kind: kind, Url: url})
	}

	return &vortrag, nil
}

func Delete(id int) error {
	_, err := db.Exec("DELETE FROM vortraege WHERE id=$1", id)
	if err != nil {
		return err
	}
	return nil
}
