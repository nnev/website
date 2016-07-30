package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/nnev/website/data"
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

func RunNext() error {
	if cmdNext.Flag.NArg() < 1 {
		return ErrUsage
	}

	n, err := strconv.Atoi(cmdNext.Flag.Arg(0))
	if err != nil {
		log.Printf("Kann %q nicht als Nummer parsen.\n", cmdNext.Flag.Arg(0))
		return ErrUsage
	}

	for _, d := range getNextThursdays(n) {
		var t data.Termin
		t.Date = data.NullTime{
			Valid: true,
			Time:  d,
		}
		t.Stammtisch = sql.NullBool{
			Valid: true,
			Bool:  d.Day() < 8,
		}
		if err := t.Insert(cmdNext.Tx); err != nil {
			return fmt.Errorf("Kann termin für %v nicht einfügen: %v", d, err)
		}
	}
	return nil
}
