package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var options struct {
	port     int
	hostname string
}

func bail(status int, template string, args ...interface{}) {
	if !strings.HasSuffix(template, "\n") {
		template += "\n"
	}
	if status == 0 {
		fmt.Fprintf(os.Stdout, template, args...)
	} else {
		fmt.Fprintf(os.Stderr, template, args...)
	}
	os.Exit(status)
}

func main() {
	flag.Parse()
	cwd, err := os.Getwd()
	if err != nil {
		bail(1, "unable to get working directory: %v", err)
	}
	dir := http.Dir(cwd)
	http.Handle("/", &logWrapper{http.FileServer(dir)})
	addr := fmt.Sprintf("%s:%d", options.hostname, options.port)
	go func() {
		started := false
		for i := 0; i < 30; i++ {
			start := time.Now()
			err := http.ListenAndServe(addr, nil)
			if time.Since(start) > time.Second {
				started = true
				break
			} else {
				fmt.Fprintf(os.Stderr, "failed to start listener: %v\n", err)
				time.Sleep(time.Second)
			}
		}
		if !started {
			bail(1, "never started successfully")
		}
	}()

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT)
	<-s
	fmt.Println("caught SIGINIT, shutting down")
}

func init() {
	flag.IntVar(&options.port, "port", 8000, "port to serve on")
	flag.StringVar(&options.hostname, "host", "", "hostname")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Lmicroseconds)
}
