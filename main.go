// Nils Elde
// https://gitlab.com/nilsanderselde

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"gitlab.com/nilsanderselde/funetik-ingglish/dbconnect"
	"gitlab.com/nilsanderselde/funetik-ingglish/global"
)

const (
	// ROOT is the subdirectory to be appended to all URLs on this server
	ROOT = "/funing"
	// SOCK is the UNIX socket for this server
	SOCK = "/web/nec/tmp/funing.sock"
)

func main() {
	// First, determine if program should run in dev or prod mode
	file, err := os.Open("env/isdev")
	if err == nil {
		fmt.Println("Development mode enabled.")
		global.IsDev = true
	}
	defer file.Close()

	// Load database connection info and calculate data for stats page
	dbconnect.DBInfo = dbconnect.GetDBInfo()
	dbconnect.StatsInit()
	dbconnect.IndexByInitial()

	http.Handle(ROOT+"/static/", setHeaders(http.StripPrefix(ROOT+"/static/", http.FileServer(http.Dir("./static")))))
	http.HandleFunc(ROOT+"/favicon.ico", faviconHandler)
	http.Handle(ROOT+"/", &templateHandler{filenames: []string{"home.html"}})
	http.Handle(ROOT+"/kybord", &templateHandler{filenames: []string{"kbd.html"}})
	http.Handle(ROOT+"/woordz", &templateHandler{filenames: []string{"words.html", "words_sorted.html"}})
	http.Handle(ROOT+"/staats", &templateHandler{filenames: []string{"stats.html"}})
	http.Handle(ROOT+"/traanzlit", &templateHandler{filenames: []string{"translit.html"}})
	http.Handle(ROOT+"/ubaawt", &templateHandler{filenames: []string{"about.html"}})

	if global.IsDev {
		// start server (development)
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal("ListenAndServe:", err)
		}
	} else {
		// start server (production)
		os.Remove(SOCK)
		unixListener, err := net.Listen("unix", SOCK)
		if err != nil {
			log.Fatal(err)
		}
		defer unixListener.Close()
		http.Serve(unixListener, nil)
	}
}

// Sets HTTP headers for handler passed to function
func setHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// cache-control headers
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		h.ServeHTTP(w, r)
	})
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/favicon.ico")
}
