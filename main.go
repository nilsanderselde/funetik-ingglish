// Nils Elde
// https://gitlab.com/nilsanderselde

package main

import (
	"log"
	"net/http"

	"gitlab.com/nilsanderselde/funetik-ingglish/dbconnect"
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
	// dbconnect.RestoreBackup()

	wordsQuery := `SELECT id,
COALESCE(COALESCE(ritin, fun), '') as fun,
COALESCE(funsil, ''),
COALESCE(trud, ''),
COALESCE(pus, ''),
COALESCE(numsil, '0'),
COALESCE(dist, '0'),
COALESCE(funsort, ''),
COALESCE(flaagd, 'false')
`
	wordsQueryFrom := `FROM words` //tshekt != true` // split up because two queries must use this part

	http.Handle("/static/", setHeaders(http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))))
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.Handle("/", &templateHandler{filenames: []string{"home.html"}})
	http.Handle("/kybord", &templateHandler{filenames: []string{"kbd.html"}})
	http.Handle("/woordz", &templateHandler{filenames: []string{"words.html", "words_sorted.html"}, query: wordsQuery, queryFrom: wordsQueryFrom})
	http.Handle("/staats", &templateHandler{filenames: []string{"stats.html"}})
	http.Handle("/traanzlit", &templateHandler{filenames: []string{"translit.html"}})
	http.Handle("/ubaawt", &templateHandler{filenames: []string{"about.html"}})

	// Start Server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
