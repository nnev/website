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
	connect = flag.String("connect", "", "The connection string to use")
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
	Links    []Link
	Password sql.NullString
}

type Link struct {
	Kind string
	Url  string
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

func verifyCaptcha(req *http.Request) bool {
	log.Printf("Entered as captcha: %q", req.FormValue("captcha"))
	return req.FormValue("captcha") == "NoName e.V."
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

	writeError(http.StatusMethodNotAllowed, res, "")
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
		writeError(http.StatusBadRequest, res, "Could not parse \"%s\" as int", idStr)
		return
	}
	if id <= 0 {
		writeError(http.StatusBadRequest, res, "Invalid id")
		return
	}

	vortrag, err := Load(id)
	if err != nil {
		log.Printf("Could not read Vortrag %d: %v\n", id, err)
		writeError(http.StatusBadRequest, res, "Could not load")
		return
	}

	if vortrag.Password.Valid && !verifyPassword(vortrag.Password.String, pw) {
		log.Println("Unauthorized edit")
		writeError(http.StatusUnauthorized, res, "Unauthorized")
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
	err := req.ParseForm()
	if err != nil {
		// Not sure what to do, so we just return 400
		log.Printf("ParseMultipartForm returned %v", err)
		writeError(http.StatusBadRequest, res, "Bad request.")
		return
	}

	idStr := req.PostFormValue("id")
	dateStr := req.PostFormValue("date")
	topic := req.PostFormValue("topic")
	abstract := req.PostFormValue("abstract")
	speaker := req.PostFormValue("speaker")
	pw := req.PostFormValue("pw")
	del := req.PostFormValue("delete")
	kinds := req.Form["kind"]
	urls := req.Form["url"]

	if len(kinds) != len(urls) {
		log.Printf("Got different numbers of kinds (%d) and urls (%d)", len(kinds), len(urls))
		if len(kinds) < len(urls) {
			urls = urls[:len(kinds)]
		} else {
			kinds = kinds[:len(urls)]
		}
	}

	var links []Link

	for i := range kinds {
		if urls[i] == "" {
			continue
		}
		if kinds[i] == "" {
			kinds[i] = "Sonstiges"
		}
		links = append(links, Link{
			Kind: kinds[i],
			Url:  urls[i],
		})
	}

	log.Printf("Incoming POST request: id=\"%s\", pw=\"%s\", date=\"%s\", topic=\"%s\", abstract=\"%s\", speaker=\"%s\", links=\"%+v\", delete=\"%s\"\n", idStr, pw, dateStr, topic, abstract, speaker, links, del)

	if !verifyCaptcha(req) {
		writeError(http.StatusBadRequest, res, "Bitte fülle das CAPTCHA korrekt aus.")
		return
	}

	if topic == "" || speaker == "" {
		writeError(http.StatusBadRequest, res, "You need to supply at least a speaker and a topic")
		return
	}

	date, _ := time.ParseInLocation("2006-01-02", dateStr, loc)

	if !date.IsZero() && date.Weekday() != time.Thursday {
		writeError(400, res, "The date is not a thursday, we currently (and for the forseeable future) don't have talks on non-thursdays.")
		return
	}

	if !date.IsZero() && date.Day() < 8 {
		writeError(400, res, "This is the first thursday of the month. Since we currently have our Stammtisch there, you can't give a talk.")
		return
	}

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
		Links:    links,
	}

	if id != -1 {
		// We need to verify the password, therefore we load the talk from the
		// db. But we don't want to overwrite anything else, so we use a new variable
		v, err := Load(id)
		if err != nil {
			log.Printf("Could not read Vortrag %d: %v\n", id, err)
			writeError(http.StatusBadRequest, res, "Could not load")
			return
		}

		if v.Password.Valid && !verifyPassword(v.Password.String, pw) {
			log.Println("Unauthorized edit")
			writeError(http.StatusUnauthorized, res, "Unauthorized")
			return
		}

		vortrag.Password = v.Password

		if del != "" {
			if err = Delete(id); err != nil {
				log.Printf("Could not delete Vortrag %d: %v\n", id, err)
				writeError(http.StatusInternalServerError, res, "Could not delete Vortrag")
				return
			}
			log.Printf("Deleted Vortrag %d", id)

			RunHook()
			url := fmt.Sprintf("/chaotische_viertelstunde.html?ts=%d", time.Now().UnixNano())
			http.Redirect(res, req, url, http.StatusSeeOther)
			return
		}

	} else {
		newPw, err := genPassword()
		if err != nil {
			log.Println("Could not generate password:", err)
			writeError(http.StatusInternalServerError, res, "Could not generate password")
			return
		}

		vortrag.Password = sql.NullString{String: newPw, Valid: true}
	}

	err = vortrag.Put()
	if err != nil {
		log.Printf("Could not update: %v\n", err)
		writeError(http.StatusBadRequest, res, "Error: %v", err)
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
