package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"gitlab.com/nilsanderselde/funetik-ingglish/levdist"
	"gitlab.com/nilsanderselde/funetik-ingglish/params"
	"gitlab.com/nilsanderselde/funetik-ingglish/runestats"
	"gitlab.com/nilsanderselde/funetik-ingglish/words"
)

type templateHandler struct {
	filename  string
	templ     *template.Template
	path      string
	query     string
	queryFrom string
	args      params.Params
	sortBy    string // template filename fragment
}

const (
	// DefaultNum is default number of words per page
	DefaultNum int = 50
)

// handle http request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	t.args.Query = t.query + t.queryFrom

	funcMap := template.FuncMap{
		"GetStats":     runestats.GetStats,
		"GetDistances": levdist.GetDistances,
		"ShowWords":    words.GetWords,
	}

	if t.path == "/" {
		t.sortBy = "words_sort_new.html"
		t.args.New = true
		t.args.Old = false
		t.args.Dist = false
		t.args.Reverse = false
		t.args.CurrentPage = "?sortby=new"

		if r.URL.Query()["sortby"] != nil {
			if r.URL.Query()["sortby"][0] == "old" {
				t.args.New = false
				t.args.Old = true
				t.sortBy = "words_sort_old.html"
				t.args.CurrentPage = "?sortby=old"

				if r.URL.Query()["order"] != nil {
					if r.URL.Query()["order"][0] == "desc" {
						t.args.Query += " ORDER BY old DESC;"
						t.args.Reverse = true
					} else {
						t.args.Query += " ORDER BY old ASC;"
					}
				} else {
					t.args.Query += " ORDER BY old ASC;"
				}

			} else if r.URL.Query()["sortby"][0] == "dist" {
				t.args.New = false
				t.args.Dist = true
				t.sortBy = "words_sort_dist.html"
				t.args.CurrentPage = "?sortby=dist"

				if r.URL.Query()["order"] != nil {
					if r.URL.Query()["order"][0] == "desc" {
						t.args.Query += " ORDER BY dist DESC;"
						t.args.Reverse = true
					} else {
						t.args.Query += " ORDER BY dist ASC;"
					}
				} else {
					t.args.Query += " ORDER BY dist ASC;"
				}
			} else if r.URL.Query()["order"] != nil {
				if r.URL.Query()["order"][0] == "desc" {
					t.args.Query += " ORDER BY words.funsort DESC;"
					t.args.Reverse = true
				} else {
					t.args.Query += " ORDER BY words.funsort ASC;"
				}
			}
			if r.URL.Query()["order"][0] == "desc" {
				t.args.CurrentPage += "&order=desc"
			} else {
				t.args.CurrentPage += "&order=asc"
			}

		} else {
			t.args.Query += ";"
		}
		if r.URL.Query()["num"] != nil {
			currNum, err := strconv.Atoi(r.URL.Query()["num"][0])
			if err != nil {
				log.Fatal(err)
			}
			t.args.Num = currNum
			t.args.NextPage = t.args.CurrentPage + "&num=" + strconv.Itoa(currNum)

			if r.URL.Query()["start"] != nil {
				currStart, err := strconv.Atoi(r.URL.Query()["start"][0])
				if err != nil {
					log.Fatal(err)
				}
				t.args.Start = currStart
				t.args.NextPage += "&start=" + strconv.Itoa(currStart+currNum)
				if currStart >= currNum {
					// PreviousPage will be nil on first page of results
					t.args.PreviousPage = t.args.CurrentPage + "&num=" + strconv.Itoa(currNum)
					t.args.PreviousPage += "&start=" + strconv.Itoa(currStart-currNum)
				} else {
					t.args.PreviousPage = ""
				}

				if err != nil {
					log.Fatal(err)
				}
			} else {
				t.args.Start = 0
				t.args.NextPage += "&start=" + strconv.Itoa(0+DefaultNum)
			}
		} else { // if missing either num or start query string
			t.args.Num = DefaultNum
			t.args.Start = 0
			t.args.NextPage = t.args.CurrentPage + "&num=" + strconv.Itoa(DefaultNum)
			t.args.NextPage += "&start=" + strconv.Itoa(0+DefaultNum)
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

	// if numrows := dbconnect.CountRows(t.queryFrom); numrows <= t.args.Start {
	// 	t.args.NextPage = ""
	// }

	t.templ.Execute(w, t.args)
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
	http.ServeFile(w, r, "relative/path/to/favicon.ico")
}

func main() {
	http.Handle("/static/", setHeaders(http.StripPrefix("/static/", http.FileServer(http.Dir("static")))))
	http.Handle("/", &templateHandler{filename: "words.html", path: "/",
		query: `SELECT words.fun,
				words.funsil,
				words.trud,
				words.pus,
				words.numsil,
				words.dist,
				words.funsort`,
		queryFrom: ` FROM words
					WHERE words.kamin = true`},
	)
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.Handle("/runestats", &templateHandler{filename: "runestats.html"})
	http.Handle("/levdist", &templateHandler{filename: "levdist.html"})

	// start server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
