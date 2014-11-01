package main

import (
	"log"
	"net/http"
	"time"
)

type logWrapper struct {
	http.Handler
}

func (l *logWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	id := newRequestId()
	log.Println(id, r.Method, r.URL)
	statsThing := &writerWatcher{ResponseWriter: w}
	start := time.Now()
	l.Handler.ServeHTTP(statsThing, r)
	log.Println(id, statsThing.headerWritten, statsThing.written, time.Since(start))
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
