package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"gitlab.com/nilsanderselde/funetik-ingglish/dbconnect"
	"gitlab.com/nilsanderselde/funetik-ingglish/levdist"
	"gitlab.com/nilsanderselde/funetik-ingglish/params"
	"gitlab.com/nilsanderselde/funetik-ingglish/runestats"
)

type templateHandler struct {
	filename string
	templ    *template.Template
	path     string
	query    string
	start    int
	num      int
	args     params.Params
	sortBy   string // template filename fragment
}

// handle http request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var err error
	funcMap := template.FuncMap{
		"GetStats":     runestats.GetStats,
		"GetDistances": levdist.GetDistances,
		"ShowWords":    dbconnect.PostgresIO,
	}

	if t.path == "/" {
		t.sortBy = "words_sort_fun.html"
		t.args.Fun = true
		t.args.Trud = false
		t.args.Dist = false

		if r.URL.Query()["sortby"] != nil {
			if r.URL.Query()["sortby"][0] == "trud" {
				t.args.Fun = false
				t.args.Trud = true
				t.sortBy = "words_sort_trud.html"
			} else if r.URL.Query()["sortby"][0] == "dist" {
				t.args.Fun = false
				t.args.Dist = true
				t.sortBy = "words_sort_dist.html"
			}
		}
		if r.URL.Query()["order"] != nil {
			if r.URL.Query()["order"][0] == "desc" {
				t.query += "ORDER BY subcipher DESC;"
			}
		} else {
			t.query += ";"
		}
		if r.URL.Query()["start"] != nil {
			t.start, err = strconv.Atoi(r.URL.Query()["start"][0])
			if err != nil {
				log.Fatal(err)
			}
		} else {
			t.start = 0
		}
		if r.URL.Query()["num"] != nil {
			t.num, err = strconv.Atoi(r.URL.Query()["num"][0])
			if err != nil {
				log.Fatal(err)
			}
		} else {
			t.num = 10
		}

		t.templ = template.Must(template.New(t.filename).Funcs(funcMap).ParseFiles(
			filepath.Join("templates", t.filename),
			filepath.Join("templates", t.sortBy),
			filepath.Join("templates", "_header.html"),
			filepath.Join("templates", "_footer.html"),
		))
	} else {
		t.templ = template.Must(template.New(t.filename).Funcs(funcMap).ParseFiles(
			filepath.Join("templates", t.filename),
			filepath.Join("templates", "_header.html"),
			filepath.Join("templates", "_footer.html"),
		))
	}

	t.templ.Execute(w, dbconnect.PostgresIO(t.query, t.start, t.num))
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

func main() {
	http.Handle("/static/", setHeaders(http.StripPrefix("/static/", http.FileServer(http.Dir("static")))))
	http.Handle("/", setHeaders(&templateHandler{filename: "words.html", path: "/", query: "SELECT * FROM formatted"}))
	http.Handle("/runestats", setHeaders(&templateHandler{filename: "runestats.html"}))
	http.Handle("/levdist", setHeaders(&templateHandler{filename: "levdist.html"}))

	// start server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
