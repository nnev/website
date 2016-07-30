package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
)

var cmdHelp = &Command{
	UsageLine: "help [cmd]",
	Short:     "",
	Long:      "",
}

func init() {
	cmdHelp.Flag = flag.NewFlagSet("help", flag.ExitOnError)
	cmdHelp.Run = RunHelp
}

func showCmdHelp(cmd *Command) {
	log.Println("Nutzung:\n")
	log.Println("    ", cmd.UsageLine, "\n")
	log.Println(cmd.Long)
}

func showGlobalHelp() {
	log.Println("Tool zum Bearbeiten der nnev-Termin Datenbank\n")
	log.Println("Nutzung:\n")
	log.Printf("    %s [flags] befehl [argumente]\n\n", os.Args[0])
	log.Println("Die vorhandenen Befehle sind:\n")

	w := tabwriter.NewWriter(os.Stderr, 8, 4, 2, ' ', 0)

	for _, cmd := range Commands {
		if cmd.Name() == "help" {
			continue
		}

		fmt.Fprintf(w, "    %s\t%s\n", cmd.Name(), cmd.Short)
	}

	w.Flush()

	log.Printf("\nDie Benutzung eines Befehls zeigt dir \"%s help [befehl]\" an.\n", os.Args[0])

	log.Println("\nFlags:\n")

	flag.PrintDefaults()
}

func RunHelp() {
	if cmdHelp.Flag.NArg() < 1 {
		showGlobalHelp()
		return
	}

	for _, cmd := range Commands {
		if cmd.Name() == "help" {
			continue
		}

		if cmd.Name() == cmdHelp.Flag.Arg(0) {
			showCmdHelp(cmd)
			return
		}
	}

	log.Printf("Unbekannter Befehl \"%s\"\n", cmdHelp.Flag.Arg(0))
}
