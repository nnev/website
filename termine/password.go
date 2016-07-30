package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

var cmdPassword = &Command{
	UsageLine: "password <id>",
	Short:     "Gibt Passwort zum editieren von Vorträgen",
	Long: `Gibt Passwort zum editieren von Vorträgen aus, sowie einen direkten
Link zum bearbeiten. id ist die id des betreffenden Vortrages (über den Link
auf der webseite in Erfahrung zu bringen)`,
	Flag:         flag.NewFlagSet("password", flag.ExitOnError),
	NeedsDB:      true,
	RegenWebsite: false,
}

func init() {
	cmdPassword.Run = RunPassword
}

func RunPassword() {
	if cmdPassword.Flag.NArg() < 1 {
		log.Printf("Nicht genug Argumente. Siehe %s help password\n", os.Args[0])
		return
	}

	id, err := strconv.Atoi(cmdPassword.Flag.Arg(0))
	if err != nil {
		log.Printf("Kann \"%s\" nicht als Nummer parsen. Siehe %s help password\n", cmdPassword.Flag.Arg(0), os.Args[0])
		return
	}

	var pw sql.NullString

	err = db.QueryRow("SELECT password FROM vortraege WHERE id = $1", id).Scan(&pw)
	if err == sql.ErrNoRows {
		log.Println("Vortrag existiert nicht")
		return
	}

	if !pw.Valid {
		log.Println("Kein Passwort gesetzt")
		return
	}

	fmt.Println("Passwort:", pw.String)
	fmt.Printf("Link: https://www.noname-ev.de/edit_c14.html?id=%d&pw=%s\n", id, pw.String)
}
