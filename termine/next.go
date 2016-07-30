package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"
)

var cmdNext = &Command{
	UsageLine: "next n",
	Short:     "Fügt die nächsten Termine automatisch ein",
	Long: `Fügt die nächsten n Termine automatisch ein.

Stammtische werden jeweils am ersten Donnerstag eines Monats ohne Ort angelegt,
alle anderen Termine werden ohne Stammtisch oder chaotische Viertelstunde
angelegt. Existente Termine werden nicht geändert.

Bei Erfolg gibt der Befehl nichts aus`,
	NeedsDB:      true,
	RegenWebsite: true,
}

func init() {
	cmdNext.Flag = flag.NewFlagSet("next", flag.ExitOnError)
	cmdNext.Run = RunNext
}

func getNextThursdays(n int) (next []time.Time) {
	cur := time.Now()
	// 0 == Sonntag
	wd := cur.Weekday()

	// Wir bestimmen zuerst letzten Donnerstag, indem wir vom aktuellen Datum
	// die richtige Anzahl an Tagen abziehen
	date := cur.AddDate(0, 0, -((3 + int(wd)) % 7))

	// Jetzt addieren wir immer 7 Tage auf dieses Datum drauf und erhalten so
	// die nächsten n Donnerstage
	for i := 0; i < n; i++ {
		date = date.AddDate(0, 0, 7)
		next = append(next, date)
	}

	return next
}

func RunNext() {
	if cmdNext.Flag.NArg() < 1 {
		log.Printf("Nicht genug Argumente. Siehe %s help next\n", os.Args[0])
		return
	}

	n, err := strconv.Atoi(cmdNext.Flag.Arg(0))
	if err != nil {
		log.Printf("Kann \"%s\" nicht als Nummer parsen. Siehe %s help next\n", cmdNext.Flag.Arg(0), os.Args[0])
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("SQL-Fehler:", err)
		return
	}
	for _, d := range getNextThursdays(n) {
		_, err := tx.Exec("INSERT INTO termine (stammtisch, date, override, location, override_long) SELECT $2, $1, '', '', '' WHERE NOT EXISTS (SELECT 1 FROM termine WHERE date = $1)", d, d.Day() < 8)
		if err != nil {
			log.Println("SQL-Fehler:", err)
			tx.Rollback()
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Println("SQL-Fehler:", err)
	}
}
