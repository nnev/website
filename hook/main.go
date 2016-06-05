package main

import (
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

	l.Println("Signature is %q", req.Header.Get("X-Hub-Signature"))

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

	var ev struct {
		Ref        string
		Head       string
		Before     string
		Size       int
		Repository struct {
			FullName string `json:"full_name"`
		}
	}

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

	if err := RunHook(); err != nil {
		l.Printf("Could not run hook: %v", err)
		http.Error(r, "internal server error", http.StatusInternalServerError)
		return
	}
	l.Printf("Done")
}

func RunHook() error {
	mtx.Lock()
	defer mtx.Unlock()

	cmd := exec.Command(*hook)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	flag.Parse()

	if s := os.Getenv("WEBHOOK_SECRET"); s == "" {
		log.Fatal("WEBHOOK_SECRET is a required environment variable")
	} else {
		secret = []byte(s)
	}

	http.HandleFunc("/", HandleHook)
	if err := http.ListenAndServe(*listen, nil); err != nil {
		log.Fatal(err)
	}
}
