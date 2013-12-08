package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"
)

var (
	addr    = flag.String("listen", "0.0.0.0:5417", "The address to listen on")
	driver  = flag.String("driver", "postgres", "The database driver to use")
	connect = flag.String("connect", "dbname=nnev host=/var/run/postgresql sslmode=disable", "The connection string to use")
	gettpl  = flag.String("template", "/var/www/www.noname-ev.de/yarpnarp.html", "The template to serve for editing zusagen")
	hook    = flag.String("hook", "", "A hook to run on every change")

	loc  *time.Location
	tpl  *template.Template
	idRe = regexp.MustCompile(`^\d*$`)
)

type Zusage struct {
	Nick      string
	Kommt     bool
	Kommentar string
	HasKommt  bool
}

func writeError(errno int, res http.ResponseWriter, format string, args ...interface{}) {
	res.WriteHeader(errno)
	fmt.Fprintf(res, format, args...)
}

func YarpNarpHandler(res http.ResponseWriter, req *http.Request) {
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
	z := Zusage{}
	if cookie, _ := req.Cookie("nick"); cookie != nil {
		z.Nick = cookie.Value
	}
	z = GetZusage(z.Nick)
	if cookie, _ := req.Cookie("kommentar"); cookie != nil {
		z.Kommentar = cookie.Value
	}

	err := tpl.ExecuteTemplate(res, "yarpnarp.html", z)
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func handlePost(res http.ResponseWriter, req *http.Request) {
	nick := req.FormValue("nick")
	kommt := (req.FormValue("kommt") == "Yarp")
	kommentar := req.FormValue("kommentar")

	if nick == "" {
		writeError(400, res, "Nick darf nicht leer sein")
		return
	}

	log.Printf("Incoming POST request: nick=\"%s\", kommt=%v, kommentar=\"%s\"\n", nick, kommt, kommentar)

	zusage := Zusage{nick, kommt, kommentar, true}

	err := zusage.Put()
	if err != nil {
		log.Printf("Could not update: %v\n", err)
		writeError(400, res, "Error: %v", err)
	}

	RunHook()

	http.SetCookie(res, &http.Cookie{ Name: "nick", Value: nick, Expires: time.Date(2030, 0, 0, 0, 0, 0, 0, time.Local) })
	http.SetCookie(res, &http.Cookie{ Name: "kommentar", Value: kommentar, Expires: time.Date(2030, 0, 0, 0, 0, 0, 0, time.Local) })
	http.Redirect(res, req, "/yarpnarp.html", 303)
}

func main() {
	flag.Parse()

	err := OpenDB()
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	tpl, err = template.New("").Delims("<<", ">>").ParseFiles(*gettpl)
	if err != nil {
		log.Fatal("Could not parse template:", err)
	}

	http.HandleFunc("/", YarpNarpHandler)

	log.Println("Listening on", *addr)
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Println("Could not listen:", err)
	}
}
