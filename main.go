package main

import (
	"flag"
	"log"
	"net/http"
	_ "github.com/lib/pq"
	"database/sql"
	"fmt"
	"regexp"
	"time"
	"html/template"
	"strconv"
)

var (
	addr = flag.String("listen", "0.0.0.0:6725", "The address to listen on")
	driver = flag.String("driver", "postgres", "The database driver to use")
	connect = flag.String("connect", "dbname=nnev host=/var/run/postgresql sslmode=disable", "The connection string to use")
	gettpl = flag.String("template", "/var/www/www.noname-ev.de/edit_c14.html", "The template to serve for editing c¼")

	db *sql.DB
	loc *time.Location
	tpl *template.Template
	idRe = regexp.MustCompile(`^\d*$`)
)

type Vortrag struct {
	Id int
	Date CustomTime
	HasDate bool
	Topic string
	Abstract string
	Speaker string
	InfoURL string
}

type CustomTime time.Time

func (t CustomTime) String() string {
	return time.Time(t).Format("2006-01-02")
}

func (t CustomTime) IsZero() bool {
	return time.Time(t).IsZero()
}

func writeError(errno int, res http.ResponseWriter, format string, args... interface{}) {
	res.WriteHeader(errno)
	fmt.Fprintf(res, format, args...)
}

func C14Handler(res http.ResponseWriter, req *http.Request) {
	var err error
	tpl, err = template.New("").Delims("<<", ">>").ParseFiles(*gettpl)
	if err != nil {
		log.Fatal("Could not parse template:", err)
	}

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

	if idStr == "" {
		err := tpl.ExecuteTemplate(res, "edit_c14.html", Vortrag{})
		if err != nil {
			log.Println(err)
			return
		}
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(400, res, "Could not parse \"%d\" as int", idStr)
		return
	}

	rows, err := db.Query("SELECT date, topic, abstract, speaker, infourl FROM vortraege WHERE id = $1", id)
	if err != nil {
		writeError(500, res, "Could not query db: %v", err)
		return
	}

	if !rows.Next() {
		writeError(400, res, "No such id", err)
		return
	}

	vortrag := Vortrag{Id: id}

	var Date time.Time
	err = rows.Scan(&Date, &vortrag.Topic, &vortrag.Abstract, &vortrag.Speaker, &vortrag.InfoURL)
	if err != nil {
		writeError(500, res, "Could not scan row: %v", err)
		return
	}
	vortrag.Date = CustomTime(Date)
	if !vortrag.Date.IsZero() {
		vortrag.HasDate = true
	}

	err = tpl.ExecuteTemplate(res, "edit_c14.html", vortrag)
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func handlePost(res http.ResponseWriter, req *http.Request) {
	idStr := req.PostFormValue("id")
	dateStr := req.PostFormValue("date")
	topic := req.PostFormValue("topic")
	abstract := req.PostFormValue("abstract")
	speaker := req.PostFormValue("speaker")
	infourl := req.PostFormValue("infourl")

	log.Printf("Incoming POST request: id=\"%s\", date=\"%s\", topic=\"%s\", abstract=\"%s\", speaker=\"%s\", infourl=\"%s\"\n", idStr, dateStr, topic, abstract, speaker, infourl)

	if topic == "" || speaker == "" {
		writeError(400, res, "You need to supply at least a speaker and a topic")
		return
	}

	date, err := time.ParseInLocation("2006-01-02", dateStr, loc)
	if err != nil {
		writeError(400, res, "Could not parse \"%s\" as date (expected Format: YYYY-MM-DD): %v", dateStr, err)
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

	tpl, err = template.New("").Delims("<<", ">>").ParseFiles(*gettpl)
	if err != nil {
		log.Fatal("Could not parse template:", err)
	}

	http.HandleFunc("/", C14Handler)

	log.Println("Listening on", *addr)
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Println("Could not listen:", err)
	}
}
