// +build ignore

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

var (
	db      *sql.DB
	driver  = flag.String("driver", "postgres", "The database driver to use")
	connect = flag.String("connect", "", "The connection string to use")
)

func openDB() (err error) {
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

func guessKind(infourl string) string {
	if strings.HasSuffix(infourl, ".pdf") {
		fmt.Println("Guessing kind 'Folien' for", infourl)
		return "Folien"
	}

	if strings.Contains(infourl, "youtu.be") {
		fmt.Println("Guessing kind 'Aufzeichnung' for", infourl)
		return "Aufzeichung"
	}

	if strings.Contains(infourl, "noname-ev.de/w/") {
		fmt.Println("Guessing kind 'Wiki-Link' for", infourl)
		return "Wikilink"
	}

	return "Sonstiges"
}

func main() {
	flag.Parse()

	err := openDB()
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	rows, err := tx.Query("SELECT id, infourl FROM vortraege;")
	if err != nil {
		log.Fatal(err)
	}

	infos := make(map[int][]string)

	for rows.Next() {
		var id int
		var infourl sql.NullString

		if err := rows.Scan(&id, &infourl); err != nil {
			log.Fatal(err)
		}

		if infourl.String == "" {
			continue
		}

		kind := guessKind(infourl.String)
		infos[id] = []string{kind, infourl.String}
	}

	for id, info := range infos {
		_, err = tx.Exec("INSERT INTO vortrag_links (vortrag, kind, url) values($1, $2, $3);", id, info[0], info[1])
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
