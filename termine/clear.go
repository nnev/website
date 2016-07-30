package main

import (
	"flag"
	"fmt"
)

var cmdClear = &Command{
	UsageLine: "clear",
	Short:     "Löscht Zu- und Absagen zum Stammtisch",
	Long: `Löscht Zu- und Absagen zum Stammtisch. Zu Benutzen nach dem
Stammtisch um für den nächsten Stammtisch die Zusagen zu managen.

Bei Erfolg gibt der Befehl nichts aus`,
	Flag:         flag.NewFlagSet("clear", flag.ExitOnError),
	NeedsDB:      true,
	RegenWebsite: false,
}

func init() {
	cmdClear.Run = RunClear
}

func RunClear() error {
	_, err := cmdClear.Tx.Exec("DELETE FROM zusagen")
	if err != nil {
		return fmt.Errorf("Kann Tabelle nicht leeren: %v", err)
	}
	return nil
}
