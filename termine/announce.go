package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"text/template"
	"time"
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

func isStammtisch(date time.Time) (stammt bool, err error) {
	err = db.QueryRow("SELECT stammtisch FROM termine WHERE date = $1", date).Scan(&stammt)
	return
}

func announceStammtisch(date time.Time) {
	loc, err := getLocation(date)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Kann Location nicht auslesen:", err)
		return
	}

	maildraft := `Liebe Treffler,

am kommenden Donnerstag ist wieder Stammtisch. Diesmal sind wir bei {{.Location}}.

Damit wir passend reservieren können, tragt bitte bis Dienstag Abend,
18:00 Uhr unter [0] ein, ob ihr kommt oder nicht.


[0] https://www.noname-ev.de/yarpnarp.html
	`

	mailtmpl := template.Must(template.New("maildraft").Parse(maildraft))
	mailbuf := new(bytes.Buffer)
	type data struct {
		Location string
	}
	if err = mailtmpl.Execute(mailbuf, data{loc}); err != nil {
		fmt.Fprintln(os.Stderr, "Fehler beim Füllen des Templates:", err)
		return
	}
	mail := mailbuf.Bytes()

	sendAnnouncement("Bitte für Stammtisch eintragen", mail)
}

func announceC14(date time.Time) {
	var data struct {
		Topic,
		Abstract,
		Speaker string
	}

	if err := db.QueryRow("SELECT topic FROM vortraege WHERE date = $1", date).Scan(&data.Topic); err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Es gibt nächsten Donnerstag noch keine c¼h. :(")
			return
		}

		fmt.Fprintln(os.Stderr, "Kann topic nicht auslesen:", err)
		return
	}

	if err := db.QueryRow("SELECT abstract FROM vortraege WHERE date = $1", date).Scan(&data.Abstract); err != nil {
		fmt.Fprintln(os.Stderr, "Kann abstract nicht auslesen:", err)
		return
	}

	if err := db.QueryRow("SELECT speaker FROM vortraege WHERE date = $1", date).Scan(&data.Speaker); err != nil {
		fmt.Fprintln(os.Stderr, "Kann speaker nicht auslesen:", err)
		return
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
	if err := mailtmpl.Execute(mailbuf, data); err != nil {
		fmt.Fprintln(os.Stderr, "Fehler beim Füllen des Templates:", err)
		return
	}
	mail := mailbuf.Bytes()
	sendAnnouncement(data.Topic, mail)
}

func sendAnnouncement(subject string, msg []byte) {
	fullmail := new(bytes.Buffer)
	fmt.Fprintf(fullmail, `From: frank@noname-ev.de
To: %s
Subject: %s

%s`, *targetmailaddr, subject, msg)

	cmd := exec.Command("/usr/sbin/sendmail", "-t")

	cmd.Stdin = fullmail

	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stdout

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Fehler beim Senden der Mail: ", err)
		fmt.Fprintln(os.Stderr, "Output von Sendmail:")
		io.Copy(os.Stderr, stdout)
	}
}

func RunAnnounce() {
	var nextRelevantDate time.Time

	if err := db.QueryRow("SELECT date FROM termine WHERE date > NOW() AND override = '' ORDER BY date ASC LIMIT 1").Scan(&nextRelevantDate); err != nil {
		fmt.Fprintln(os.Stderr, "Kann nächsten Termin nicht auslesen:", err)
		return
	}

	isStm, err := isStammtisch(nextRelevantDate)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Kann stammtischiness nicht auslesen:", err)
		return
	}

	if isStm {
		announceStammtisch(nextRelevantDate)
	} else {
		announceC14(nextRelevantDate)
	}
}
