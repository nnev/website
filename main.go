package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	driver  = flag.String("driver", "postgres", "The SQL-driver to use")
	connect = flag.String("connect", "dbname=nnev host=/var/run/postgresql sslmode=disable", "The connection string to use")
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
}

var Commands = []*Command{
	cmdNext,
	cmdHelp,
}

func (cmd *Command) Name() string {
	name := cmd.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Printf("Usage: %s [flags] cmd [args]\n", os.Args[0])
		os.Exit(1)
	}

	for _, cmd := range Commands {
		if cmd.Name() != flag.Arg(0) {
			continue
		}

		args := flag.Args()
		cmd.Flag.Parse(args[1:])

		if cmd.NeedsDB {
			err := OpenDB()
			if err != nil {
				fmt.Println("Fehler beim Verbinden zur Datenbank:", err)
				return
			}
		}

		cmd.Run()
		return
	}

	fmt.Printf("Unknown command \"%s\"\n", flag.Arg(0))
	os.Exit(1)
}
