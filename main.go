// Nils Elde
// https://gitlab.com/nilsanderselde

package main

import (
	"log"
	"net/http"

	"gitlab.com/nilsanderselde/funetik-ingglish/dbconnect"
)

const (
	// ROOT is the subdirectory to be appended to all URLs on this server
	ROOT = "/funing"
	// SOCK is the UNIX socket for this server
	SOCK = "/web/nec/funing/go.sock"
)

func main() {
	// Load database connection info and calculate data for stats page
	dbconnect.DBInfo = dbconnect.GetDBInfo()
	dbconnect.StatsInit()

	http.Handle(ROOT+"/static/", setHeaders(http.StripPrefix(ROOT+"/static/", http.FileServer(http.Dir("./static")))))
	http.HandleFunc(ROOT+"/favicon.ico", faviconHandler)
	http.Handle(ROOT+"/", &templateHandler{filenames: []string{"home.html"}})
	http.Handle(ROOT+"/kybord", &templateHandler{filenames: []string{"kbd.html"}})
	http.Handle(ROOT+"/woordz", &templateHandler{filenames: []string{"words.html", "words_sorted.html"}})
	http.Handle(ROOT+"/staats", &templateHandler{filenames: []string{"stats.html"}})
	http.Handle(ROOT+"/traanzlit", &templateHandler{filenames: []string{"translit.html"}})
	http.Handle(ROOT+"/ubaawt", &templateHandler{filenames: []string{"about.html"}})

	// start server (development)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

	// start server (production)

	// unix, err := net.Listen("unix", SOCK)
	// if err != nil {
	// 	log.Fatal("Listen error: ", err)
	// }
	// sigchan := make(chan os.Signal, 1)
	// signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	// go func(unix net.Listener, c chan os.Signal) {
	// 	sig := <-c
	// 	log.Printf("Caught signal %s: shutting down.", sig)
	// 	unix.Close()
	// 	os.Exit(0)
	// }(unix, sigchan)

	// http.Serve(unix, nil)

	// if err := os.Remove(SOCK); err != nil {
	// 	log.Fatal(err)
	// }
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
