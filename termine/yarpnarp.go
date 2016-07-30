package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/nnev/website/data"
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

type Zusagen []*data.Zusage

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

	// If one or both registered fields are NULL, just treat them as the zero
	// time (so they sort before all other values)
	if z[i].Registered.Time.Before(z[j].Registered.Time) {
		return true
	} else if z[j].Registered.Time.Before(z[i].Registered.Time) {
		return false
	}

	// If one or both nick fields are NULL, just tream them as the empty string
	if z[i].Nick.String < z[j].Nick.String {
		return true
	}

	return false
}

func (z Zusagen) minWidth() int {
	var nick int
	for _, zusage := range z {
		if len(zusage.Nick.String) > nick {
			nick = len(zusage.Nick.String)
		}
	}
	// nick + padding + Kommt + Date
	return nick + 1 + 6 + 20
}

func formatBool(b bool) string {
	if b {
		return "Ja"
	}
	return "Nein"
}

func maybeTruncate(s string, width int, truncate bool) string {
	if !truncate {
		return s
	}

	if width <= 0 {
		return "…"
	}

	if len(s) > width {
		return s[:width-1] + "…"
	}
	return s
}

func RunYarpNarp() error {
	var zusagen Zusagen

	it := data.Zusagen(cmdYarpNarp.Tx)
	for it.Next() {
		zusagen = append(zusagen, it.Val())
	}
	if err := it.Err(); err != nil {
		return err
	}

	sort.Sort(zusagen)

	width, trunc := TermWidth()

	width -= zusagen.minWidth()

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 1, ' ', tabwriter.DiscardEmptyColumns)
	fmt.Fprintf(w, "Nick\tKommt\tLetzte Änderung\t%s\n", maybeTruncate("Kommentar", width, trunc))
	for _, z := range zusagen {
		fmt.Fprintf(w, "%s\t%v\t%s\t%s\n", z.Nick, formatBool(z.Kommt), z.Registered.Time.In(time.Local).Format("2006-01-02 15:04:05"), maybeTruncate(z.Kommentar, width, trunc))
	}
	w.Flush()
	return nil
}
