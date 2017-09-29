package main

import (
	"log"
	"net/http"

	"gitlab.com/nilsanderselde/funetik-ingglish/dbconnect"
)

const (
	// DefaultNum is default number of words per page
	DefaultNum int = 20
)

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

func main() {
	dbconnect.DBInfo = dbconnect.GetDBInfo()

	http.Handle("/static/", setHeaders(http.StripPrefix("/static/", http.FileServer(http.Dir("static")))))
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.Handle("/", &templateHandler{filename: "words.html",
		query: `
SELECT
	id,
    COALESCE(COALESCE(ritin, fun), '') as fun,
    COALESCE(funsil, ''),
    COALESCE(trud, ''),
    COALESCE(pus, ''),
    COALESCE(numsil, '0'),
    COALESCE(dist, '0'),
	COALESCE(funsort, ''),
	COALESCE(flaagd, 'false')`,
		queryFrom: `
FROM words WHERE fun LIKE '%ar%' AND trud LIKE '%or%'
	`},
	/* SQL scratch area
	COALESCE(fun, '') as new,
	WHERE kamin = true
	*/
	)

	http.Handle("/runestats", &templateHandler{filename: "runestats.html"})
	http.Handle("/levdist", &templateHandler{filename: "levdist.html"})

	// start server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
