package main

import (
	"bytes"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/quotedprintable"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/nnev/website/data"
)

var cmdAnnounce = &Command{
	UsageLine: "announce",
	Short:     "Kündigt nächsten Stammtisch oder nächste c¼h an",
	Long: `Kündigt den nächsten Stammtisch oder die nächste c¼h per E-Mail an,
je nachdem, was am nächsten Donnerstag ist.`,
	Flag:         flag.NewFlagSet("announce", flag.ExitOnError),
	NeedsDB:      true,
	RegenWebsite: false,
}

var (
	targetmailaddr  string
	confirmAnnounce bool
)

func init() {
	cmdAnnounce.Run = RunAnnounce
	cmdAnnounce.Flag.StringVar(&targetmailaddr, "address", "ccchd@ccchd.de", "Mailadresse, an die Ankündigungen gehen sollen.")
	cmdAnnounce.Flag.BoolVar(&confirmAnnounce, "confirm", false, "Frage nach Bestätigung bevor die Mail gesendet wird")
}

func announceStammtisch(t *data.Termin) error {
	maildraft := `Liebe Treffler,

am kommenden Donnerstag ist wieder Stammtisch. Diesmal sind wir bei {{.Location}}.

Damit wir passend reservieren können, tragt bitte bis Dienstag Abend,
18:00 Uhr unter [0] ein, ob ihr kommt oder nicht.


[0] https://www.noname-ev.de/yarpnarp.html
	`

	mailtmpl := template.Must(template.New("maildraft").Parse(maildraft))
	mailbuf := new(bytes.Buffer)
	if err := mailtmpl.Execute(mailbuf, t); err != nil {
		return fmt.Errorf("Fehler beim Füllen des Templates: %v", err)
	}
	mail := mailbuf.Bytes()

	return sendAnnouncement("Bitte für Stammtisch eintragen", mail)
}

func announceC14(t *data.Termin) error {
	vortrag, err := t.GetVortrag(cmdAnnounce.Tx)
	if err == sql.ErrNoRows {
		fmt.Println("Es gibt nächsten Donnerstag noch keine c¼h. :(")
		return nil
	}
	if err != nil {
		log.Fatal("Kann vortrag nicht lesen:", err)
	}

	maildraft := `Liebe Treffler,

am kommenden Donnerstag wird {{.Speaker}} eine c¼h zum Thema

    {{.Topic}}

halten.

Kommet zahlreich!


Wer mehr Informationen möchte:

{{.Abstract}}
	`

	mailtmpl := template.Must(template.New("maildraft").Parse(maildraft))
	mailbuf := new(bytes.Buffer)
	if err := mailtmpl.Execute(mailbuf, vortrag); err != nil {
		return fmt.Errorf("Fehler beim Füllen des Templates: %v", err)
	}
	mail := mailbuf.Bytes()
	return sendAnnouncement(vortrag.Topic, mail)
}

func sendAnnouncement(subject string, msg []byte) error {
	mail := new(bytes.Buffer)
	fmt.Fprintf(mail, "From: frank@noname-ev.de\r\n")
	fmt.Fprintf(mail, "To: %s\r\n", mime.QEncoding.Encode("utf-8", targetmailaddr))
	fmt.Fprintf(mail, "Subject: %s\r\n", mime.QEncoding.Encode("utf-8", subject))
	fmt.Fprintf(mail, "Content-Type: %s\r\n", mime.FormatMediaType("text/plain", map[string]string{"charset": "utf-8"}))
	fmt.Fprintf(mail, "Content-Transfer-Encoding: quoted-printable\r\n")
	fmt.Fprintf(mail, "\r\n")

	body := quotedprintable.NewWriter(mail)
	body.Write(msg)
	body.Close()

	ok, err := getConfirmation(mail.Bytes())
	if !ok || err != nil {
		log.Println("Abgebrochen")
		return err
	}

	cmd := exec.Command("/bin/true", "-t")

	cmd.Stdin = mail

	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stdout

	if err := cmd.Run(); err != nil {
		io.Copy(os.Stderr, stdout)
		return fmt.Errorf("Fehler beim Senden der Mail: %v", err)
	}
	return nil
}

func getConfirmation(mail []byte) (ok bool, err error) {
	if !confirmAnnounce {
		return true, nil
	}
	fmt.Println("Folgende E-Mail wird versendet:")
	os.Stdout.Write(mail)
	fmt.Println()
	fmt.Print("Senden? [j/N] ")
	var answer string
	if _, err = fmt.Scan(&answer); err != nil {
		return false, err
	}
	answer = strings.ToLower(answer)
	return (strings.HasPrefix(answer, "j") || strings.HasPrefix(answer, "y")), nil
}

func RunAnnounce() error {
	t, err := data.FutureTermine(cmdAnnounce.Tx).First()
	if err == sql.ErrNoRows {
		return errors.New("Keine termine gefunden")
	}
	if err != nil {
		return fmt.Errorf("Kann nächsten Termin nicht auslesen: %v", err)
	}

	if t.Stammtisch.Bool {
		return announceStammtisch(t)
	}
	return announceC14(t)
}
