package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	_ "github.com/lib/pq"

	"github.com/nnev/website/data"
)

var (
	ErrUsage = errors.New("wrong usage")
)

type Command struct {
	// Run runs the command.
	Run func() error

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

	// Tx will be set to a transaction that should be used for all database
	// accesses, if NeedsDB is true.
	Tx *sql.Tx

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

func (cmd *Command) parseAndRun(hook string) (err error) {
	args := flag.Args()
	helpShort := cmd.Flag.Bool("h", false, "")
	help := cmd.Flag.Bool("help", false, "Zeige diese Hilfe")
	cmd.Flag.Parse(args[1:])

	if *help || *helpShort {
		showCmdHelp(cmd)
		return nil
	}

	if cmd.NeedsDB {
		db, err := data.OpenDB()
		if err != nil {
			return fmt.Errorf("Fehler beim Verbinden zur Datenbank: %v", err)
		}
		defer db.Close()

		cmd.Tx, err = db.Begin()
		if err != nil {
			return fmt.Errorf("Kann keine Transaktion starten: %v", err)
		}
		defer func() {
			if err != nil {
				cmd.Tx.Rollback()
			}
			if err = cmd.Tx.Commit(); err != nil {
				err = fmt.Errorf("Kann Transaktion nicht committen: %v", err)
			}
		}()
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	if cmd.RegenWebsite {
		cmd := exec.Command(hook)
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Println("Hook fehlgeschlagen:", err)
			log.Println("Output:")
			log.Fatal(string(output))
		}
	}
	return nil
}

func ExpectNArg(fs *flag.FlagSet, n int) error {
	if fs.NArg() < n {
		log.Printf("Nicht genug Argumente.")
		return ErrUsage
	}
	if fs.NArg() > n {
		log.Printf("Zu viele Argumente.")
		return ErrUsage
	}
	return nil
}

func main() {
	websitehook := flag.String("hook", "/usr/bin/update-website", "Hook zum neu Bauen der Website")
	helpShort := flag.Bool("h", false, "")
	help := flag.Bool("help", false, "Zeige Hilfe")

	log.SetFlags(0)

	flag.Parse()
	if *help || *helpShort || flag.NArg() < 1 {
		cmdHelp.Run()
		os.Exit(1)
	}

	for _, cmd := range Commands {
		if cmd.Name() != flag.Arg(0) {
			continue
		}

		if err := cmd.parseAndRun(*websitehook); err != nil {
			if err == ErrUsage {
				if cmd != cmdHelp {
					showCmdHelp(cmd)
					return
				}
				os.Exit(2)
			}
			log.Fatal(err)
		}
		return
	}
	log.Printf("Unbekannter Befehl %q", flag.Arg(0))

	showGlobalHelp()
	os.Exit(2)
}
