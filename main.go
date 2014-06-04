package main

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"
)

var (
	addr    = flag.String("listen", "0.0.0.0:6725", "The address to listen on")
	driver  = flag.String("driver", "postgres", "The database driver to use")
	connect = flag.String("connect", "dbname=nnev host=/var/run/postgresql sslmode=disable", "The connection string to use")
	gettpl  = flag.String("template", "/var/www/www.noname-ev.de/edit_c14.html", "The template to serve for editing c¼")
	hook    = flag.String("hook", "", "A hook to run on every change")

	loc  *time.Location
	tpl  *template.Template
	idRe = regexp.MustCompile(`^\d*$`)
)

type Vortrag struct {
	Id       int
	Date     CustomTime
	HasDate  bool
	Topic    string
	Abstract string
	Speaker  string
	InfoURL  string
	Password sql.NullString
}

type CustomTime time.Time

func (t CustomTime) String() string {
	return time.Time(t).Format("2006-01-02")
}

func (t CustomTime) IsZero() bool {
	return time.Time(t).IsZero()
}

func writeError(errno int, res http.ResponseWriter, format string, args ...interface{}) {
	res.WriteHeader(errno)
	fmt.Fprintf(res, format, args...)
}

func genPassword() (string, error) {
	// 120 bits of entropy should be enough for a password
	buf := make([]byte, 15)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}
	pw := base64.URLEncoding.EncodeToString(buf)
	log.Println("Generated password:", pw)
	return pw, nil
}

func verifyPassword(a, b string) bool {
	// Since our passwords all have the same length, this does not actually
	// leak any information
	if len(a) != len(b) {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

func C14Handler(res http.ResponseWriter, req *http.Request) {
	var err error
	tpl, err = template.New("").Delims("((", "))").ParseFiles(*gettpl)
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
	pw := req.FormValue("pw")

	log.Printf("Incoming GET request: id=\"%s\" pw=\"%s\"\n", idStr, pw)

	if idStr == "" {
		dateStr := req.FormValue("date")
		v := Vortrag{Id: -1}
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
	if id <= 0 {
		writeError(400, res, "Invalid id")
		return
	}

	vortrag, err := Load(id)
	if err != nil {
		log.Printf("Could not read Vortrag %d: %v\n", id, err)
		writeError(400, res, "Could not load")
		return
	}

	if vortrag.Password.Valid && !verifyPassword(vortrag.Password.String, pw) {
		log.Println("Unauthorized edit")
		writeError(401, res, "Unauthorized")
		return
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
	pw := req.PostFormValue("pw")

	log.Printf("Incoming POST request: id=\"%s\", pw=\"%s\", date=\"%s\", topic=\"%s\", abstract=\"%s\", speaker=\"%s\", infourl=\"%s\"\n", idStr, pw, dateStr, topic, abstract, speaker, infourl)

	if topic == "" || speaker == "" {
		writeError(400, res, "You need to supply at least a speaker and a topic")
		return
	}

	date, _ := time.ParseInLocation("2006-01-02", dateStr, loc)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		id = -1
	}

	vortrag := &Vortrag{
		Id:       id,
		Date:     CustomTime(date),
		HasDate:  false,
		Topic:    topic,
		Abstract: abstract,
		Speaker:  speaker,
		InfoURL:  infourl,
	}

	if id != -1 {
		// We need to verify the password, therefore we load the talk from the
		// db. But we don't want to overwrite anything else, so we use a new variable
		v, err := Load(id)
		if err != nil {
			log.Printf("Could not read Vortrag %d: %v\n", id, err)
			writeError(400, res, "Could not load")
			return
		}

		if v.Password.Valid && !verifyPassword(v.Password.String, pw) {
			log.Println("Unauthorized edit")
			writeError(401, res, "Unauthorized")
			return
		}

		vortrag.Password = v.Password
	} else {
		newPw, err := genPassword()
		if err != nil {
			log.Println("Could not generate password:", err)
			writeError(500, res, "Could not generate password")
			return
		}

		vortrag.Password = sql.NullString{newPw, true}
	}

	err = vortrag.Put()
	if err != nil {
		log.Printf("Could not update: %v\n", err)
		writeError(400, res, "Error: %v", err)
		return
	}

	RunHook()

	url := fmt.Sprintf("/edit_c14.html?id=%d&pw=%s", vortrag.Id, vortrag.Password.String)

	fmt.Println("Redirecting:", url)

	http.Redirect(res, req, url, 303)
}

func main() {
	flag.Parse()

	// We ignore SIGPIPE, because it might be generated by the hook
	// occasionally. Since we don't care for it, we don't need a buffered
	// channel
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGPIPE)

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
