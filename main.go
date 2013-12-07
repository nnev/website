package main

import (
	"flag"
	"log"
	"net/http"
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
	hook = flag.String("hook", "", "A hook to run on every change")

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
		dateStr := req.FormValue("date")
		v := Vortrag{}
		if dateStr != "" {
			date, err := time.ParseInLocation("2006-01-02", dateStr, loc)
			if err == nil {
				v.Date = CustomTime(date)
				v.HasDate = true
			}
		}

		err := tpl.ExecuteTemplate(res, "edit_c14.html", v)
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

	vortrag, err := Load(id)
	if err != nil {
		log.Printf("Could not read Vortrag %d: %v\n", id, err)
		writeError(400, res, "Could not load")
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

	date, _ := time.ParseInLocation("2006-01-02", dateStr, loc)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		id = -1
	}

	vortrag := Vortrag{ id, CustomTime(date), false, topic, abstract, speaker, infourl }
	err = vortrag.Put()
	if err != nil {
		log.Printf("Could not update: %v\n", err)
		writeError(400, res, "Error")
	}

	RunHook()

	http.Redirect(res, req, "/chaotische_viertelstunde.html", 303)
}

func main() {
	flag.Parse()

	err := OpenDB()
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
