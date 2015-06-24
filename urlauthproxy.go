package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var urlmap map[string]string

// Commandline
var listenAddr = flag.String("http", "localhost:9999", "http host and port to listen to")
var certfile = flag.String("cert", "", "tls certificate file")
var keyfile = flag.String("key", "", "tls key file")
var urlmapfile = flag.String("urlmap", "urlmap", "urlmap file, key = value set of values")
var logfile = flag.String("logfile", "", "log file (leave empty for stdout)")

func LoadUrlMap() {
	urlmap = make(map[string]string)

	f, err := os.Open(*urlmapfile)
	if err != nil {
		log.Fatalf("Could not open configuration file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		pieces := strings.SplitN(line, "=", 2)

		urlmap[pieces[0]] = pieces[1]
	}
}

func httpBaseHandler(w http.ResponseWriter, r *http.Request) {
	url, ok := urlmap[r.URL.Path]
	if !ok {
		log.Printf("[%s] %s: Not found", r.RemoteAddr, r.URL.Path)
		http.NotFound(w, r)
		return
	}

	/* Ok, now make an outgoing call to get this resource */
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[%s] %s ->%s: ERR:%s", r.RemoteAddr, r.URL, url, err)
		http.Error(w, "Error fetching resource", 500)
		return
	}

	/* Copy all headers */
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}

	/* Copy data */
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[%s] %s ->%s: ERR COPY: %s", r.RemoteAddr, r.URL, url, err)
		http.Error(w, "Error copying response", 500)
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(result)

	log.Printf("[%s] %s ->%s: %d", r.RemoteAddr, r.URL, url, len(result))
}


func main() {
	flag.Parse()

	if *logfile != "" {
		lf, err := os.OpenFile(*logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("Could not open logfile: %s", err)
		}
		defer lf.Close()

		log.SetOutput(lf)
	}

	LoadUrlMap()

	http.HandleFunc("/", httpBaseHandler)

	log.Printf("Started with urlmap of %d entries.", len(urlmap))

	err := http.ListenAndServeTLS(*listenAddr, *certfile, *keyfile, nil)
	if err != nil {
		log.Fatalf("Could not listen: %s", err)
	}
}
