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
	"gitlab.com/nilsanderselde/funetik-ingglish/words"
)

type templateHandler struct {
	filename  string
	templ     *template.Template
	path      string
	query     string
	queryFrom string
	args      params.Params
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
		t.args.PreviousPage = ""
		// t.args.NextPage = ""
		t.args.New = true
		t.args.Old = false
		t.args.Dist = false
		t.args.Reverse = false
		t.args.CurrentPage = "?sortby=new"

		if r.URL.Query()["sortby"] != nil {
			if r.URL.Query()["sortby"][0] == "old" {
				t.args.New = false
				t.args.Old = true
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
		// if there is a valid "num" query string
		if r.URL.Query()["num"] != nil {

			// get the current value of it
			currNum, err := strconv.Atoi(r.URL.Query()["num"][0])
			if err != nil {
				log.Fatal(err)
			}

			// set the template's args object's value to it
			t.args.Num = currNum

			// create the next page query string for the "next page" link
			t.args.NextPage = t.args.CurrentPage + "&num=" + strconv.Itoa(currNum)

			// if there is a valid "start" query string
			if r.URL.Query()["start"] != nil {

				// get the current value of it
				currStart, err := strconv.Atoi(r.URL.Query()["start"][0])
				if err != nil {
					log.Fatal(err)
				}

				// set the template's args object's value to it
				t.args.Start = currStart

				// add to the next page query string
				t.args.NextPage += "&start=" + strconv.Itoa(currStart+currNum)

				// if the current starting number is creater than or equal to the
				// number of words per page, add a previous page button
				if currStart >= currNum {
					t.args.PreviousPage = t.args.CurrentPage + "&num=" + strconv.Itoa(currNum)
					t.args.PreviousPage += "&start=" + strconv.Itoa(currStart-currNum)
					// PreviousPage will be nil on first page of results
				} else {

				}
			} else { // if there isn't a valid start query string, set the starting
				// position to 0 in both the template's args and query string
				t.args.Start = 0
				t.args.NextPage += "&start=" + strconv.Itoa(0+DefaultNum)
				t.args.PreviousPage = ""
			}
		} else { // if missing a valid "num" query string
			// set everything to defaults (start at 0, DefaultNum words per page)
			t.args.Num = DefaultNum
			t.args.Start = 0
			// create the default nextpage query string for the "next page" link
			t.args.NextPage = t.args.CurrentPage + "&num=" + strconv.Itoa(DefaultNum)
			t.args.NextPage += "&start=" + strconv.Itoa(0+DefaultNum)
		}
		// join the template files for the wordlist
		t.templ = template.Must(template.New(t.filename).Funcs(funcMap).ParseFiles(
			filepath.Join("templates", t.filename),
			filepath.Join("templates", "words_sorted.html"),
			filepath.Join("templates", "_header.html"),
			filepath.Join("templates", "_footer.html"),
		))
	} else { // not word list, just a regular page, join the template files
		t.templ = template.Must(template.New(t.filename).Funcs(funcMap).ParseFiles(
			filepath.Join("templates", t.filename),
			filepath.Join("templates", "_header.html"),
			filepath.Join("templates", "_footer.html"),
		))
	}
	// see if next page button should be hidden before fetching rows
	numrows := dbconnect.CountRows(t.queryFrom)
	// if total number of rows is less than or equal to the number the next page would start on,
	//	 dont show a next page link.
	if numrows < t.args.Start+t.args.Num {
		t.args.NextPage = ""
	}
	if numrows < t.args.Start {
		t.args.PreviousPage = ""
	}

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
		query: `SELECT COALESCE(words.ritin, words.fun) as new,
				words.funsil,
				words.trud as old,
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
