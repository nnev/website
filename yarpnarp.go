package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/jmoiron/sqlx"
)

var cmdYarpNarp = &Command{
	UsageLine:    "yarpnarp",
	Short:        "Zeige Zu- und Absagen für den Stammtisch",
	Long:         `Zeigt die Zu- und Absagen für den Stammtisch an`,
	Flag:         flag.NewFlagSet("clear", flag.ExitOnError),
	NeedsDB:      true,
	RegenWebsite: false,
}

func init() {
	cmdYarpNarp.Run = RunYarpNarp
}

type Zusage struct {
	Nick       string
	Kommt      bool
	Kommentar  string
	Registered time.Time
}

type Zusagen []Zusage

func (z Zusagen) Len() int {
	return len(z)
}

func (z Zusagen) Swap(i, j int) {
	z[i], z[j] = z[j], z[i]
}

func (z Zusagen) Less(i, j int) bool {
	if z[i].Kommt && !z[j].Kommt {
		return true
	} else if z[j].Kommt && !z[i].Kommt {
		return false
	}

	if z[i].Registered.Before(z[j].Registered) {
		return true
	} else if z[j].Registered.Before(z[i].Registered) {
		return false
	}

	if z[i].Nick < z[j].Nick {
		return true
	}

	return false
}

func formatBool(b bool) string {
	if b {
		return "Ja"
	} else {
		return "Nein"
	}
}

func RunYarpNarp() {
	var zusagen Zusagen
	dbx := sqlx.NewDb(db, *driver)

	if err := sqlx.Select(dbx, &zusagen, "SELECT nick, kommt, kommentar, registered FROM zusagen"); err != nil {
		fmt.Fprintln(os.Stderr, "Datenbankfehler:", err)
		return
	}

	sort.Sort(zusagen)

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 1, ' ', tabwriter.DiscardEmptyColumns)
	fmt.Fprintf(w, "Nick\tKommt\tLetzte Änderung\n")
	for _, z := range zusagen {
		fmt.Fprintf(w, "%s\t%v\t%s\n", z.Nick, formatBool(z.Kommt), z.Registered.In(time.Local).Format("2006-01-02 15:04:05"))
	}
	w.Flush()
}
