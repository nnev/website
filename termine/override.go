package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var cmdOverride = &Command{
	UsageLine: "override datum kurzbeschreibung",
	Short:     "Überschreibt Beschreibung eines Termins",
	Long: `Überschreibt die Beschreibung eines Termins.

Kann z.B. genutzt werden, um einen ausfallenden Treff oder eine MV zu
markieren. Die Kurzbeschreibung taucht auf der Homepage unter Aktuelles auf.
Der Befehl liest ausserdem von stdin eine lange Beschreibung, die für die
E-Mail Ankündigung benutzt wird. Soll der override aufgehoben werden, sollte
beides der leere String sein.

Das Datum muss im Format 2006-02-28 angegeben werden.

Bei Erfolg wird nichts zurück gegeben`,
	Flag:         flag.NewFlagSet("override", flag.ExitOnError),
	NeedsDB:      true,
	RegenWebsite: true,
}

func init() {
	cmdOverride.Run = RunOverride
}

func RunOverride() {
	if cmdOverride.Flag.NArg() < 2 {
		log.Printf("Nicht genug Argumente. Siehe %s help override\n", os.Args[0])
		return
	}

	date, err := time.ParseInLocation("2006-01-02", cmdOverride.Flag.Arg(0), time.Local)
	if err != nil {
		log.Printf("Kann \"%s\" nicht als Datum parsen. Siehe %s help override\n", cmdNext.Flag.Arg(0), os.Args[0])
		return
	}

	override_long, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Println("Kann nicht von stdin lesen:", err)
		return
	}

	result, err := db.Exec("UPDATE termine SET override = $2, override_long = $3 WHERE date = $1", date, cmdOverride.Flag.Arg(1), string(override_long))
	if err != nil {
		log.Println("Kann Eintrag nicht ändern:", err)
		return
	}

	if n, err := result.RowsAffected(); err != nil && n == 0 {
		log.Println("Termin noch nicht vorhanden.")
		log.Println("Füge ihn erst mittels next hinzu.")
	}
}
