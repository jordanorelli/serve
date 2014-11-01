package main

import (
	"github.com/jordanorelli/serve/termcolors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type logWrapper struct {
	http.Handler
}

func (l *logWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := newRequestId()
	if IsTerminal(os.Stdout.Fd()) {
		log.Println(termcolors.WrapString(termcolors.Magenta, id.String()), r.Method, r.URL)
	} else {
		log.Println(id, r.Method, r.URL)
	}
	statsThing := &writerWatcher{ResponseWriter: w}
	start := time.Now()
	l.Handler.ServeHTTP(statsThing, r)
	log.Println(id, statsThing.PrettyStatus(), statsThing.written, time.Since(start))
}

type writerWatcher struct {
	http.ResponseWriter
	written       int
	headerWritten int
}

func (w *writerWatcher) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.written += n
	return n, err
}

func (w *writerWatcher) WriteHeader(status int) {
	if w.headerWritten != 0 {
		log.Println("somehow we wrote two headers.")
	}
	w.headerWritten = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *writerWatcher) PrettyStatus() string {
	s := strconv.Itoa(w.headerWritten)
	if !IsTerminal(os.Stdout.Fd()) {
		return s
	}
	switch {
	case w.headerWritten < 100:
		panic("wtf")
	case w.headerWritten < 200:
		return termcolors.WrapString(termcolors.White, s)
	case w.headerWritten < 300:
		return termcolors.WrapString(termcolors.Green, s)
	case w.headerWritten < 400:
		return termcolors.WrapString(termcolors.White, s)
	case w.headerWritten < 500:
		return termcolors.WrapString(termcolors.Yellow, s)
	default:
		return termcolors.WrapString(termcolors.Red, s)
	}
}
