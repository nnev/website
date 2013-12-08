package main

import (
	"flag"
	"fmt"
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
	fmt.Println("Nutzung:\n")
	fmt.Println("    ", cmd.UsageLine, "\n")
	fmt.Println(cmd.Long)
}

func showGlobalHelp() {
	fmt.Println("Tool zum Bearbeiten der nnev-Termin Datenbank\n")
	fmt.Println("Nutzung:\n")
	fmt.Printf("    %s [flags] befehl [argumente]\n\n", os.Args[0])
	fmt.Println("Die vorhandenen Befehle sind:\n")

	w := tabwriter.NewWriter(os.Stdout, 8, 4, 2, ' ', 0)

	for _, cmd := range Commands {
		if cmd.Name() == "help" {
			continue
		}

		fmt.Fprintf(w, "    %s\t%s\n", cmd.Name(), cmd.Short)
	}

	w.Flush()

	fmt.Printf("\nDie Benutzung eines Befehls zeigt dir \"%s help [befehl]\" an.\n", os.Args[0])
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

	fmt.Printf("Unbekannter Befehl \"%s\"\n", cmdHelp.Flag.Arg(0))
}
