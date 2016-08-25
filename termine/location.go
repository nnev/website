package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/nnev/website/data"
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

func RunLocation() error {
	if cmdLocation.Flag.NArg() > 1 {
		showCmdHelp(cmdLocation)
		os.Exit(1)
	}

	t, err := data.QueryTermine(cmdLocation.Tx, "WHERE date >= $1 AND stammtisch = true", time.Now()).First()
	if err == sql.ErrNoRows {
		return errors.New("Termin muss erst mittels next hinzugefügt werden.")
	}
	if err != nil {
		return fmt.Errorf("Kann Termin nicht lesen: %v", err)
	}
	if cmdLocation.Flag.NArg() == 0 {
		fmt.Println(t.Location)
		return nil
	}
	t.Location = cmdLocation.Flag.Arg(0)
	if err = t.Update(cmdLocation.Tx); err != nil {
		return fmt.Errorf("Kann Location nicht setzen: %v", err)
	}
	return nil
}
