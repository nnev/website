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
	"text/template"

	"github.com/nnev/website/data"
)

var cmdAnnounce = &Command{
	UsageLine: "announce",
	Short:     "Kündigt nächsten Stammtisch oder nächste c¼h an",
	Long: `Kündigt den nächsten Stammtisch oder die nächste c¼h an,
je nachdem, was am nächsten Donnerstag ist.`,
	Flag:         flag.NewFlagSet("announce", flag.ExitOnError),
	NeedsDB:      true,
	RegenWebsite: false,
}

var targetmailaddr = flag.String("announceAddress", "ccchd@ccchd.de", "Mailadresse, an die Ankündigungen gehen sollen.")

func init() {
	cmdAnnounce.Run = RunAnnounce
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
	v, err := t.GetVortrag(cmdAnnounce.Tx)
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
	if err := mailtmpl.Execute(mailbuf, v); err != nil {
		return fmt.Errorf("Fehler beim Füllen des Templates: %v", err)
	}
	mail := mailbuf.Bytes()
	return sendAnnouncement(v.Topic, mail)
}

func sendAnnouncement(subject string, msg []byte) error {
	mail := new(bytes.Buffer)
	fmt.Fprintf(mail, "From: frank@noname-ev.de\r\n")
	fmt.Fprintf(mail, "To: %s\r\n", mime.QEncoding.Encode("utf-8", *targetmailaddr))
	fmt.Fprintf(mail, "Subject: %s\r\n", mime.QEncoding.Encode("utf-8", subject))
	fmt.Fprintf(mail, "Content-Type: %s\r\n", mime.FormatMediaType("text/plain", map[string]string{"charset": "utf-8"}))
	fmt.Fprintf(mail, "Content-Transfer-Encoding: quoted-printable\r\n")
	fmt.Fprintf(mail, "\r\n")

	body := quotedprintable.NewWriter(mail)
	body.Write(msg)
	body.Close()

	cmd := exec.Command("/usr/sbin/sendmail", "-t")

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
