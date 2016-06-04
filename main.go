package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	driver      = flag.String("driver", "postgres", "Der benutzte sql-Treiber")
	connect     = flag.String("connect", "dbname=nnev user=anon host=/var/run/postgresql sslmode=disable", "Die Verbindusgsspezifikation")
	websitehook = flag.String("hook", "/srv/git/website.git/hooks/post-update", "Hook zum neu Bauen der Website")
	_           = flag.Bool("help", false, "Zeige Hilfe")
)

type Command struct {
	// Run runs the command.
	Run func()

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'go help' output.
	Short string

	// Long is the long message shown in the 'go help <this-command>' output.
	Long string

	Flag *flag.FlagSet

	// NeedsDB is true, if the command needs to connect to the database
	NeedsDB bool

	// RegenWebsite is true, if the website needs to be rebuild after the command
	RegenWebsite bool
}

var Commands = []*Command{
	cmdLocation,
	cmdNext,
	cmdClear,
	cmdOverride,
	cmdPassword,
	cmdYarpNarp,
	cmdHelp,
	cmdAnnounce,
}

func (cmd *Command) Name() string {
	name := cmd.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (cmd *Command) parseAndRun() {
	args := flag.Args()
	cmd.Flag.Parse(args[1:])

	if cmd.NeedsDB {
		err := OpenDB()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Fehler beim Verbinden zur Datenbank:", err)
			return
		}
	}

	cmd.Run()

	if cmd.RegenWebsite {
		cmd := exec.Command(*websitehook)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Hook fehlgeschlagen:", err)
			fmt.Fprintln(os.Stderr, "Output:")
			fmt.Fprint(os.Stderr, string(output))
		}
	}
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		cmdHelp.Run()
		os.Exit(1)
	}

	for _, cmd := range Commands {
		if cmd.Name() != flag.Arg(0) {
			continue
		}

		cmd.parseAndRun()
		return
	}

	cmdHelp.parseAndRun()
	os.Exit(1)
}
