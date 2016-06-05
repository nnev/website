package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var cmdLocation = &Command{
	UsageLine: "location [locname]",
	Short:     "Zeigt oder ändert Location des nächsten Stammtisches",
	Long: `Zeigt oder ändert die Location des nächsten Stammtisches.
Schreibweise muss mit Feld 'locname' im stammtisch_*.md der entsprechenden
Stammtischseite im website-git übereinstimmen!

Bei Erfolg gibt der Befehl nichts aus.`,
	Flag:         flag.NewFlagSet("location", flag.ExitOnError),
	NeedsDB:      true,
	RegenWebsite: true,
}

func init() {
	cmdLocation.Run = RunLocation
}

func getLocation(date time.Time) (location string, err error) {
	err = db.QueryRow("SELECT location FROM termine WHERE date = $1", date).Scan(&location)
	return
}

func setLocation(date time.Time, location string) (updated bool, err error) {
	result, err := db.Exec("UPDATE termine SET location = $2 WHERE date = $1", date, location)
	if err != nil {
		return false, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return n > 0, nil
}

func RunLocation() {
	// Wir holen uns die nächsten 5 Donnerstage -- darunter muss ein Stammtisch
	// sein
	var stammtisch time.Time
	for _, d := range getNextThursdays(5) {
		if d.Day() < 8 {
			stammtisch = d
			break
		}
	}

	if cmdLocation.Flag.NArg() == 0 {
		loc, err := getLocation(stammtisch)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Kann Location nicht auslesen:", err)
			return
		}
		fmt.Println(loc)
	} else {
		updated, err := setLocation(stammtisch, cmdLocation.Flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Kann Location nicht setzen:", err)
			return
		}
		if !updated {
			fmt.Fprintln(os.Stderr, "Termin noch nicht vorhanden.")
			fmt.Fprintln(os.Stderr, "Füge ihn erst mittels next hinzu.")
		}
	}
}
