package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"hash"
	"io"
	"log"
	"mime"
	"mime/quotedprintable"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	// mtx protects new builds, so that concurrent runs of the webhook don't
	// interfere with each other.
	mtx = new(sync.Mutex)

	// reqId contains a unique id of a request, used for logging of concurrent
	// requests.
	reqId id

	hook   = flag.String("hook", "/srv/git/website.git/post-update", "The hook to run when the website repo is pushed")
	listen = flag.String("listen", "localhost:12345", "The interface/port to listen on")
	repo   = flag.String("repo", "nnev/website", "Only rebuild on push to this repository")
	ref    = flag.String("ref", "refs/heads/master", "Only rebuild on push to this ref")

	secret []byte

	events chan<- event
)

type id uint64

func (i *id) String() string {
	return fmt.Sprintf("%x", uint64(*i))
}

func (i *id) Next() id {
	return id(atomic.AddUint64((*uint64)(i), 1))
}

// verifier implements the verification of the signature by github. This is
// factored into it's own type, to ease security-review.
type verifier struct {
	io.Reader
	io.Closer
	sig []byte
	h   hash.Hash
}

func newVerifier(req *http.Request, secret []byte) (v *verifier, err error) {
	v = new(verifier)
	if s := req.Header.Get("X-Hub-Signature"); s == "" {
		return nil, errors.New("no signature provided")
	} else {
		if !strings.HasPrefix(s, "sha1=") {
			// According to https://developer.github.com/webhooks/securing/ the
			// signature *always* uses sha1.
			return nil, errors.New("malformed signature: must start with sha1=")
		}
		if v.sig, err = hex.DecodeString(s[5:]); err != nil {
			return nil, fmt.Errorf("malformed signature: %v", err)
		}
	}

	v.h = hmac.New(sha1.New, secret)
	v.Reader = io.TeeReader(req.Body, v.h)
	v.Closer = req.Body

	return v, nil
}

func (v *verifier) Verify() bool {
	return hmac.Equal(v.h.Sum(nil), v.sig)
}

type event struct {
	Ref        string
	Head       string
	Before     string
	Size       int
	Repository struct {
		FullName string `json:"full_name"`
	}
}

func HandleHook(r http.ResponseWriter, req *http.Request) {
	rid := reqId.Next()
	l := log.New(os.Stderr, rid.String()+": ", log.LstdFlags)

	l.Printf("Event guid is %q", req.Header.Get("X-GitHub-Delivery"))

	if req.Method != http.MethodPost {
		l.Printf("Ignoring method %q", req.Method)
		http.Error(r, "", http.StatusMethodNotAllowed)
		return
	}

	if ct := req.Header.Get("Content-Type"); ct != "application/json" {
		l.Printf("Ignoring invalid content type %q", ct)
		http.Error(r, "", http.StatusUnsupportedMediaType)
		return
	}

	l.Printf("Signature is %q\n", req.Header.Get("X-Hub-Signature"))

	// github signs their events, we need to verify them.
	v, err := newVerifier(req, secret)
	if err != nil {
		l.Print(err)
		http.Error(r, "", http.StatusForbidden)
		return
	}

	// payload size is limited to 5 MB by github. We also enforce this limit,
	// to prevent DoS.
	body := http.MaxBytesReader(r, v, 5*(1<<20))

	var ev event

	dec := json.NewDecoder(body)
	if err := dec.Decode(&ev); err != nil {
		l.Printf("Could not unmarshal json: %v", err)
		http.Error(r, "invalid json: "+err.Error(), http.StatusBadRequest)
		return
	}

	if !v.Verify() {
		l.Printf("Invalid signature")
		http.Error(r, "", http.StatusForbidden)
		return
	}

	l.Printf("Got push event: %+v", ev)

	if ev.Ref != *ref {
		l.Printf("Ignoring ref %q", ev.Ref)
		return
	}

	if ev.Repository.FullName != *repo {
		l.Printf("Ignoring ref %q", ev.Repository.FullName)
		return
	}

	// The events channels contains a one-element buffer and we drop events
	// that don't fit into that. The buffer of one guarantees, that if we have
	// multiple updates in quick succession, there will always be another
	// hook-run scheduled. We don't need more than one buffered event, because
	// the next run of the hook will pull *all* updates, not just the one this
	// event was for, so one run is enough.
	select {
	case events <- ev:
		l.Printf("Dispatched event")
	default:
		l.Printf("Buffer is full, dropped update")
	}

	l.Printf("Done")
}

func Build(ch <-chan event) {
	for range ch {
		stdout := new(bytes.Buffer)

		cmd := exec.Command(*hook)
		cmd.Stdout = stdout
		cmd.Stderr = stdout
		err := cmd.Run()
		if err == nil {
			continue
		}
		log.Printf("Building website failed: %v", err)

		contentType := mime.FormatMediaType("text/plain", map[string]string{"charset": "utf-8"})

		mail := new(bytes.Buffer)
		fmt.Fprintf(mail, "To: root\r\n")
		fmt.Fprintf(mail, "From: webmaster@eris.noname-ev.de\r\n")
		fmt.Fprintf(mail, "Subject: Failed website build\r\n")
		fmt.Fprintf(mail, "Content-Type: %s\r\n", contentType)
		fmt.Fprintf(mail, "Content-Transfer-Encoding: quoted-printable\r\n")
		fmt.Fprintf(mail, "\r\n")

		body := quotedprintable.NewWriter(mail)
		fmt.Fprintf(body, "The website failed to build in response to a github hook:\n")
		fmt.Fprintf(body, "%v\n\n", err)
		io.Copy(body, stdout)
		body.Close()

		cmd = exec.Command("/usr/sbin/sendmail", "-t")
		cmd.Stdin = mail
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Printf("Could not send failure mail: %v", err)
		}
	}
}

func main() {
	flag.Parse()

	if s := os.Getenv("WEBHOOK_SECRET"); s == "" {
		log.Fatal("WEBHOOK_SECRET is a required environment variable")
	} else {
		secret = []byte(s)
	}

	ch := make(chan event, 1)
	events = ch
	go Build(ch)

	http.HandleFunc("/", HandleHook)
	if err := http.ListenAndServe(*listen, nil); err != nil {
		log.Fatal(err)
	}
}
