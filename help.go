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
	fmt.Fprintln(os.Stderr, "Nutzung:\n")
	fmt.Fprintln(os.Stderr, "    ", cmd.UsageLine, "\n")
	fmt.Fprintln(os.Stderr, cmd.Long)
}

func showGlobalHelp() {
	fmt.Fprintln(os.Stderr, "Tool zum Bearbeiten der nnev-Termin Datenbank\n")
	fmt.Fprintln(os.Stderr, "Nutzung:\n")
	fmt.Fprintf(os.Stderr, "    %s [flags] befehl [argumente]\n\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "Die vorhandenen Befehle sind:\n")

	w := tabwriter.NewWriter(os.Stderr, 8, 4, 2, ' ', 0)

	for _, cmd := range Commands {
		if cmd.Name() == "help" {
			continue
		}

		fmt.Fprintf(w, "    %s\t%s\n", cmd.Name(), cmd.Short)
	}

	w.Flush()

	fmt.Fprintf(os.Stderr, "\nDie Benutzung eines Befehls zeigt dir \"%s help [befehl]\" an.\n", os.Args[0])

	fmt.Fprintln(os.Stderr, "\nFlags:\n")

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

	fmt.Fprintf(os.Stderr, "Unbekannter Befehl \"%s\"\n", cmdHelp.Flag.Arg(0))
}
