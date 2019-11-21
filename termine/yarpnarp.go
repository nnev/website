package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/nnev/website/data"
)

type OrderPriority = uint8

const (
	Asc OrderPriority = iota
	Desc
)

type ZusageOrderBy = uint8

const (
	ZusageByNick ZusageOrderBy = iota
	ZusageByKommt
	ZusageByKommentar
	ZusageByRegistered
)

type ZusageOrder struct {
	by       ZusageOrderBy
	priority OrderPriority
}

var cmdYarpNarp = &Command{
	UsageLine:    "yarpnarp",
	Short:        "Zeige Zu- und Absagen für den Stammtisch",
	Long:         `Zeigt die Zu- und Absagen für den Stammtisch an`,
	Flag:         flag.NewFlagSet("yarpnarp", flag.ExitOnError),
	NeedsDB:      true,
	RegenWebsite: false,
}

var sortOrder string

func init() {
	cmdYarpNarp.Flag.StringVar(&sortOrder, "sort", "-kommt,+registered", "Bestimmt die Sortierreihenfolge")
	cmdYarpNarp.Run = RunYarpNarp
}

type Zusagen []*data.Zusage

func (z Zusagen) Len() int {
	return len(z)
}

func (z Zusagen) Swap(i, j int) {
	z[i], z[j] = z[j], z[i]
}

func (z Zusagen) Less(i, j int, key []ZusageOrder) bool {
	for _, order := range key {
		i_, j_ := i, j
		if order.priority == Desc {
			i_, j_ = j, i
		}
		switch order.by {
		case ZusageByKommentar:
			if z[i_].Kommentar < z[j_].Kommentar {
				return true
			}
		case ZusageByKommt:
			if !z[i_].Kommt && z[j_].Kommt {
				return true
			}
		case ZusageByNick:
			if z[i_].Nick.String < z[j_].Nick.String {
				return true
			}
		case ZusageByRegistered:
			if z[i_].Registered.Time.Before(z[j_].Registered.Time) {
				return true
			}
		default:
			panic("Invalid ZusageOrderBy field")
		}
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

func parseSortOrder(s string) []ZusageOrder {
	elements := strings.Split(s, ",")
	result := make([]ZusageOrder, len(elements))
	for i, element := range elements {
		result[i].priority = Asc

		if element[0] == '-' {
			result[i].priority = Desc
		}
		if element[0] == '+' || element[0] == '-' {
			element = element[1:]
		}

		switch element {
		case "nick":
			result[i].by = ZusageByNick
		case "kommentar":
			result[i].by = ZusageByKommentar
		case "kommt":
			result[i].by = ZusageByKommt
		case "registered":
			result[i].by = ZusageByRegistered
		default:
			// FIXME: Probably not panic, because it is an expected error.
			panic("Invalid sort order element")
		}
	}
	return result
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
	if err := it.Close(); err != nil {
		return err
	}

	order := parseSortOrder(sortOrder)
	sort.Slice(zusagen, func(i, j int) bool {
		return zusagen.Less(i, j, order)
	})

	width, trunc := TermWidth()

	width -= zusagen.minWidth()

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 1, ' ', tabwriter.DiscardEmptyColumns)
	fmt.Fprintf(w, "Nick\tKommt\tLetzte Änderung\t%s\n", maybeTruncate("Kommentar", width, trunc))
	for _, z := range zusagen {
		fmt.Fprintf(w, "%s\t%v\t%s\t%s\n", z.Nick.String, formatBool(z.Kommt), z.Registered.Time.In(time.Local).Format("2006-01-02 15:04:05"), maybeTruncate(z.Kommentar, width, trunc))
	}
	w.Flush()
	return nil
}
