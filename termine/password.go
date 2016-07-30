package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/nnev/website/data"
)

var cmdPassword = &Command{
	UsageLine: "password <id>",
	Short:     "Gibt Passwort zum editieren von Vorträgen",
	Long: `Gibt Passwort zum editieren von Vorträgen aus, sowie einen direkten
Link zum bearbeiten. id ist die id des betreffenden Vortrages (über den Link
auf der webseite in Erfahrung zu bringen)`,
	Flag:         flag.NewFlagSet("password", flag.ExitOnError),
	NeedsDB:      true,
	RegenWebsite: false,
}

func init() {
	cmdPassword.Run = RunPassword
}

func RunPassword() error {
	if err := ExpectNArg(cmdPassword.Flag, 1); err != nil {
		return err
	}

	id, err := strconv.Atoi(cmdPassword.Flag.Arg(0))
	if err != nil {
		log.Printf("Kann %q nicht als Nummer parsen.", cmdPassword.Flag.Arg(0))
		return ErrUsage
	}

	v, err := data.GetVortrag(cmdPassword.Tx, id)
	if err == sql.ErrNoRows {
		return errors.New("Vortrag existiert nicht")
	}
	if err != nil {
		return fmt.Errorf("Kann Vortrag nicht lesen: %v", err)
	}

	if v.Password == "" {
		fmt.Println("Kein Password gesetzt")
		return nil
	}
	fmt.Println("Passwort:", v.Password)
	fmt.Printf("Link: https://www.noname-ev.de/edit_c14.html?id=%d&pw=%s\n", id, v.Password)
	return nil
}
