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
	log.Println("Nutzung:")
	log.Println()
	log.Println("    ", cmd.UsageLine)
	log.Println()
	log.Println(cmd.Long)
}

func showGlobalHelp() {
	log.Println("Tool zum Bearbeiten der nnev-Termin Datenbank.")
	log.Println()
	log.Println("Nutzung:")
	log.Println()
	log.Printf("    %s [flags] befehl [argumente]\n\n", os.Args[0])
	log.Println("Die vorhandenen Befehle sind:")
	log.Println()

	w := tabwriter.NewWriter(os.Stderr, 8, 4, 2, ' ', 0)

	for _, cmd := range Commands {
		if cmd.Name() == "help" {
			continue
		}

		fmt.Fprintf(w, "    %s\t%s\n", cmd.Name(), cmd.Short)
	}

	w.Flush()

	log.Printf("\nDie Benutzung eines Befehls zeigt dir \"%s help [befehl]\" an.\n", os.Args[0])

	log.Printlnr()
	log.Println("Flags:")
	log.Println()

	flag.PrintDefaults()
}

func RunHelp() error {
	if cmdHelp.Flag.NArg() < 1 {
		showGlobalHelp()
		return nil
	}

	for _, cmd := range Commands {
		if cmd.Name() == cmdHelp.Flag.Arg(0) {
			showCmdHelp(cmd)
			return nil
		}
	}

	return fmt.Errorf("Unbekannter Befehl \"%s\"\n", cmdHelp.Flag.Arg(0))
}
