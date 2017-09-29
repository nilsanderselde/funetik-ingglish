package main

import (
	"fmt"
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
	query     string
	queryFrom string
	args      params.TemplateParams
}

const (
	// DefaultNum is default number of words per page
	DefaultNum int = 20
)

// Redirect redirects to the passed URL
func Redirect(url string) {
	fmt.Println("test")
}

// handle http request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// pointers to TemplateParams fields for more readable code
	aNew, aOld, aDist, aID := &t.args.New, &t.args.Old, &t.args.Dist, &t.args.ID
	aRev, aQ, aStart, aNum := &t.args.Reverse, &t.args.Query, &t.args.Start, &t.args.Num
	aSort, aCurr, aNext, aPrev := &t.args.Sort, &t.args.CurrentPage, &t.args.NextPage, &t.args.PreviousPage

	funcMap := template.FuncMap{
		"GetStats":     runestats.GetStats,
		"GetDistances": levdist.GetDistances,
		"GetWords":     words.GetWords,
	}

	if t.filename == "words.html" {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		// put two pieces of query together
		*aQ = t.query + t.queryFrom

		// reset prev page to force handler to recreate it if needed
		*aPrev = ""

		// if there is a sortby value
		if r.URL.Query()["sortby"] != nil {
			// if it's set to sort by "id"
			if r.URL.Query()["sortby"][0] == "id" {
				*aNew, *aOld, *aDist, *aID = false, false, false, true
				*aSort = "?sortby=id"
				*aQ += " ORDER BY id"
				if r.URL.Query()["order"] != nil && r.URL.Query()["order"][0] == "desc" {
					*aQ += " DESC;"
					*aSort += "&order=desc"
					*aRev = true
				} else {
					*aQ += " ASC;"
					*aSort += "&order=asc"
					*aRev = false
				}
				// if it's set to sort by "new"
			} else if r.URL.Query()["sortby"][0] == "new" {
				*aNew, *aOld, *aDist, *aID = true, false, false, false
				*aSort = "?sortby=new"
				*aQ += " ORDER BY funsort"
				if r.URL.Query()["order"] != nil && r.URL.Query()["order"][0] == "desc" {
					*aQ += " DESC;"
					*aSort += "&order=desc"
					*aRev = true
				} else {
					*aQ += " ASC;"
					*aSort += "&order=asc"
					*aRev = false
				}
				// if it's set to sort by "old"
			} else if r.URL.Query()["sortby"][0] == "old" {
				*aNew, *aOld, *aDist, *aID = false, true, false, false
				*aSort = "?sortby=old"
				*aQ += " ORDER BY trud"
				if r.URL.Query()["order"] != nil && r.URL.Query()["order"][0] == "desc" {
					*aQ += " DESC;"
					*aSort += "&order=desc"
					*aRev = true
				} else {
					*aQ += " ASC;"
					*aSort += "&order=asc"
					*aRev = false
				}
				// if it's set to sort by "dist"
			} else if r.URL.Query()["sortby"][0] == "dist" {
				*aNew, *aOld, *aDist, *aID = false, false, true, false
				*aSort = "?sortby=dist"
				*aQ += " ORDER BY dist"
				if r.URL.Query()["order"] != nil && r.URL.Query()["order"][0] == "desc" {
					*aQ += " DESC;"
					*aSort += "&order=desc"
					*aRev = true
				} else {
					*aQ += " ASC;"
					*aSort += "&order=asc"
					*aRev = false
				}
			}
		} else {
			// default values: sort by new spelling, ascending, starting at zero, incrementing by DefaultNum
			*aNew, *aOld, *aDist, *aID = true, false, false, false
			*aSort = "?sortby=new&order=asc"
			*aRev = false
			*aQ += " ORDER BY funsort ASC;"
		}
		// if there is a valid "num" query string
		if r.URL.Query()["num"] != nil {

			// get the current value of it
			currNum, err := strconv.Atoi(r.URL.Query()["num"][0])
			if err != nil {
				// log.Fatal(err)
				currNum = DefaultNum
			}

			// set the template's args object's value to it
			*aNum = currNum

			// create the next page query string for the "next page" link
			*aNext = *aSort + "&num=" + strconv.Itoa(currNum)
			*aCurr = *aNext

			// if there is a valid "start" query string
			if r.URL.Query()["start"] != nil {

				// get the current value of it
				currStart, err := strconv.Atoi(r.URL.Query()["start"][0])
				if err != nil {
					// log.Fatal(err)
					currStart = 0
				}

				// set the template's args object's value to it
				*aStart = currStart

				// add to the next page query string
				*aNext += "&start=" + strconv.Itoa(currStart+currNum)
				*aCurr += "&start=" + strconv.Itoa(currStart)
				// if the current starting number is creater than or equal to the
				// number of words per page, add a previous page button
				if currStart >= currNum {
					*aPrev = *aSort + "&num=" + strconv.Itoa(currNum)
					*aPrev += "&start=" + strconv.Itoa(currStart-currNum)
					// PreviousPage will be nil on first page of results
				} else {

				}
			} else { // if there isn't a valid start query string, set the starting
				// position to 0 in both the template's args and query string
				*aStart = 0
				*aNext += "&start=" + strconv.Itoa(0+DefaultNum)
				*aCurr += "&start=" + strconv.Itoa(0)
				*aPrev = ""
			}
		} else { // if missing a valid "num" query string
			// set everything to defaults (start at 0, DefaultNum words per page)
			*aNum = DefaultNum
			*aStart = 0
			// create the default nextpage query string for the "next page" link
			*aNext = *aSort + "&num=" + strconv.Itoa(DefaultNum) + "&start=" + strconv.Itoa(0+DefaultNum)
			*aCurr = *aSort + "&num=" + strconv.Itoa(DefaultNum) + "&start=" + strconv.Itoa(0)
		}

		if r.URL.Query()["id"] != nil {
			// if it's set to sort by "new"
			id := r.URL.Query()["id"][0]
			// *aSort += "&id=" + id
			// fmt.Println(*aSort)

			dbconnect.FlagRow(id)

			// fmt.Printf("%s%s%s\n", r.Host, r.URL.Path, *aSort)

			// http.Redirect(w, r, fmt.Sprintf("%s%s%s", r.Host, r.URL.Path, *aSort), 201)
		}

		// join the template files for the wordlist
		t.templ = template.Must(template.New(t.filename).Funcs(funcMap).ParseFiles(
			filepath.Join("templates", t.filename),
			filepath.Join("templates", "words_sorted.html"),
			filepath.Join("templates", "_header.html"),
			filepath.Join("templates", "_footer.html"),
		))
		// see if next page button should be hidden before fetching rows
		numrows := dbconnect.CountRows(t.queryFrom)
		// if total number of rows is less than or equal to the number the next page would start on,
		//	 dont show a next page link.
		if numrows < *aStart+*aNum {
			*aNext = ""
		}
		if numrows < *aStart {
			*aPrev = ""
		}
	} else { // not word list, just a regular page, join the template files

		t.templ = template.Must(template.New(t.filename).Funcs(funcMap).ParseFiles(
			filepath.Join("templates", t.filename),
			filepath.Join("templates", "_header.html"),
			filepath.Join("templates", "_footer.html"),
		))
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
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.Handle("/", &templateHandler{filename: "words.html",
		query: `
SELECT
	id,
    COALESCE(fun, ''),
    COALESCE(funsil, ''),
    COALESCE(trud, ''),
    COALESCE(pus, ''),
    COALESCE(numsil, '0'),
    COALESCE(dist, '0'),
	COALESCE(funsort, ''),
	COALESCE(flaagd, 'false')`,
		queryFrom: `
FROM words
	`},
	/* SQL scratch area
	COALESCE(COALESCE(ritin, fun), '') as new,
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
