package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/nnev/website/data"
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

func RunOverride() error {
	if cmdOverride.Flag.NArg() != 2 {
		log.Printf("Falsche Anzahl an Argumenten.")
		return ErrUsage
	}

	date, err := time.ParseInLocation("2006-01-02", cmdOverride.Flag.Arg(0), time.Local)
	if err != nil {
		log.Printf("Kann \"%s\" nicht als Datum parsen.", cmdNext.Flag.Arg(0))
		return ErrUsage
	}

	override_long, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("Kann nicht von stdin lesen: %v", err)
	}

	t, err := data.GetTermin(cmdOverride.Tx, date)
	if err == sql.ErrNoRows {
		return errors.New("Termin nicht vorhanden. Füge ihn erst mittels next hinzu.")
	}
	if err != nil {
		return fmt.Errorf("Kann Termin nicht lesen:", err)
	}

	t.Override = cmdOverride.Flag.Arg(1)
	t.OverrideLong = string(override_long)
	return t.Update(cmdOverride.Tx)
}
