package main

import (
	"flag"
	"fmt"
	"os"
)

var cmdClear = &Command{
	UsageLine: "clear",
	Short:     "Löscht Zu- und Absagen zum Stammtisch",
	Long: `Löscht Zu- und Absagen zum Stammtisch. Zu Benutzen nach dem
Stammtisch um für den nächsten Stammtisch die Zusagen zu managen.

Bei Erfolg gibt der Befehl nichts aus`,
	Flag:    flag.NewFlagSet("clear", flag.ExitOnError),
	NeedsDB: true,
}

func init() {
	cmdClear.Run = RunClear
}

func RunClear() {
	_, err := db.Exec("DELETE FROM zusagen")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Kann Tabelle nicht leeren:", err)
	}
}
