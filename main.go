package main

import (
	"flag"
	"log"
	"net/http"
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	addr = flag.String("listen", "0.0.0.0:6725", "The address to listen on")
	driver = flag.String("driver", "postgres", "The database driver to use")
	connect = flag.String("connect", "dbname=nnev host=/var/run/postgresql sslmode=disable", "The connection string to use")
	gettpl = flag.String("template", "/var/www/www.noname-ev.de/edit_c14.html", "The template to serve for editing c¼")

	db *sql.DB
	loc *time.Location
	tpl string
	idRe = regexp.MustCompile(`^\d*$`)
)

func writeError(errno int, res http.ResponseWriter, format string, args... interface{}) {
	res.WriteHeader(errno)
	fmt.Fprintf(res, format, args...)
}

func C14Handler(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		handlePost(res, req)
		return
	}

	if req.Method == "GET" {
		handleGet(res, req)
		return
	}

	writeError(405, res, "")
	return
}

func handleGet(res http.ResponseWriter, req *http.Request) {
	idStr := req.FormValue("id")

	log.Printf("Incoming GET request: id=\"%s\"\n", idStr)

	if !idRe.MatchString(idStr) {
		writeError(400, res, "id must be numerical")
		return
	}

	io.Copy(res, strings.NewReader(strings.Replace(tpl, "__C14_ID__", idStr, -1)))
}

func handlePost(res http.ResponseWriter, req *http.Request) {
	idStr := req.PostFormValue("id")
	dateStr := req.PostFormValue("date")
	topic := req.PostFormValue("topic")
	abstract := req.PostFormValue("abstract")
	speaker := req.PostFormValue("speaker")
	infourl := req.PostFormValue("infourl")

	log.Printf("Incoming POST request: id=\"%s\", date=\"%s\", topic=\"%s\", abstract=\"%s\", speaker=\"%s\", infourl=\"%s\"\n", idStr, dateStr, topic, abstract, speaker, infourl)

	date, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		writeError(400, res, "Could not parse \"%s\" as date: %v", dateStr, err)
		return
	}

	if idStr == "" {
		stmt, err := db.Prepare("INSERT INTO vortraege (date, topic, abstract, speaker, infourl) VALUES ($1, $2, $3, $4, $5)")
		if err != nil {
			writeError(500, res, "Could not prepare insert statement: %v", err)
			return
		}

		_, err = stmt.Exec(date, topic, abstract, speaker, infourl)
		if err != nil {
			writeError(500, res, "Could not insert: %v", err)
			return
		}
	} else {
		stmt, err := db.Prepare("UPDATE vortraege SET date = $1, topic = $2, abstract = $3, speaker = $4, infourl = $5 WHERE id = $6")
		if err != nil {
			writeError(500, res, "Could not prepare update statement: %v", err)
			return
		}

		_, err = stmt.Exec(date, topic, abstract, speaker, infourl, idStr)
		if err != nil {
			writeError(500, res, "Could not update: %v", err)
			return
		}
	}
}

func main() {
	flag.Parse()

	var err error

	db, err = sql.Open(*driver, *connect)
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	loc, err = time.LoadLocation("CET")
	if err != nil {
		log.Fatal("Could not load timezone-data:", err)
	}

	file, err := os.Open(*gettpl)
	if err != nil {
		log.Fatal("Could not open template:", err)
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal("Could not read template:", err)
	}
	file.Close()

	tpl = string(bytes)

	http.HandleFunc("/", C14Handler)
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Println("Could not listen:", err)
	}
}
