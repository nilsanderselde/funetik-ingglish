// Nils Elde
// https://gitlab.com/nilsanderselde

package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"gitlab.com/nilsanderselde/funetik-ingglish/dbconnect"
	"gitlab.com/nilsanderselde/funetik-ingglish/global"
)

// Options for ?update=all (func UpdateAllAutoValues)
const (
	fun     bool = false
	numsil  bool = false
	funsort bool = true
	dist    bool = false
)

// templateHandler contains all fields needed to process and execute templates
type templateHandler struct {
	filename  string
	templ     *template.Template
	query     string
	queryFrom string
	args      global.TemplateParams
}

// creates random number to prompt page reloads after db changes
func randomRune() string {
	return string(global.CurrRand)
}

// handle http request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// pointers to TemplateParams fields for more readable code
	aNew, aOld, aDist, aID := &t.args.New, &t.args.Old, &t.args.Dist, &t.args.ID
	aRev, aQ, aStart, aNum := &t.args.Reverse, &t.args.Query, &t.args.Start, &t.args.Num
	aSort, aCurr, aNext, aPrev := &t.args.Sort, &t.args.CurrentPage, &t.args.NextPage, &t.args.PreviousPage

	funcMap := template.FuncMap{
		"GetStats":  dbconnect.GetStats,
		"ShowWords": dbconnect.ShowWords,
		"Random":    randomRune,
	}
	// special processing for words list based on query strings
	if t.filename == "words.html" {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		// put two pieces of query together
		*aQ = t.query + t.queryFrom

		// reset prev page to force templateHandler to recreate it if needed
		*aPrev = ""

		// if set to update automatically generated values
		if r.URL.Query()["update"] != nil {
			if r.URL.Query()["update"][0] == "all" {
				dbconnect.UpdateAllAutoValues(fun, numsil, funsort, dist)
			}
		}

		// if set to sort by value
		if r.URL.Query()["sortby"] != nil {
			sortBy := &r.URL.Query()["sortby"][0]
			desc := "&order=desc"
			asc := "&order=asc"

			// sort by id, old spelling, distance, or default (new spelling)
			if *sortBy == "id" {
				*aNew, *aOld, *aDist, *aID = false, false, false, true
				*aSort = "?sortby=id"
				*aQ += " ORDER BY id"
			} else if *sortBy == "old" {
				*aNew, *aOld, *aDist, *aID = false, true, false, false
				*aSort = "?sortby=old"
				*aQ += " ORDER BY trud"
			} else if *sortBy == "dist" {
				*aNew, *aOld, *aDist, *aID = false, false, true, false
				*aSort = "?sortby=dist"
				*aQ += " ORDER BY dist, funsort"
			} else {
				*aNew, *aOld, *aDist, *aID = true, false, false, false
				*aSort = "?sortby=new"
				*aQ += " ORDER BY funsort"
			}
			// ascending or descending
			if r.URL.Query()["order"] != nil && r.URL.Query()["order"][0] == "desc" {
				*aQ += " DESC;"
				*aSort += desc
				*aRev = true
			} else {
				*aQ += " ASC;"
				*aSort += asc
				*aRev = false
			}
		} else {
			// default values: sort by new spelling, ascending, starting at zero, incrementing by DefaultNum
			*aNew, *aOld, *aDist, *aID = true, false, false, false
			*aSort = "?sortby=new&order=asc"
			*aRev = false
			*aQ += " ORDER BY funsort ASC;"
		}
		// if there is a valid "num" query string...
		if r.URL.Query()["num"] != nil {
			// ...get the current value of it, preventing massive queries
			currNum, err := strconv.Atoi(r.URL.Query()["num"][0])
			if err != nil || currNum > 1000 {
				currNum = DefaultNum
			}
			// ...set the template's args object's value to it
			*aNum = currNum

			// create the next page query string for the "next page" link
			*aNext = *aSort + "&num=" + strconv.Itoa(currNum)
			*aCurr = *aNext

			// if there is a valid "start" query string...
			if r.URL.Query()["start"] != nil {

				// ...get the current value of it
				currStart, err := strconv.Atoi(r.URL.Query()["start"][0])
				if err != nil {
					currStart = 0
				}
				// ...set the template's args object's value to it
				*aStart = currStart

				// add to the next page query string
				*aNext += "&start=" + strconv.Itoa(currStart+currNum)
				*aCurr += "&start=" + strconv.Itoa(currStart)
				// if curr starting number >= num of words per page, add prev page button
				if currStart >= currNum {
					*aPrev = *aSort + "&num=" + strconv.Itoa(currNum)
					*aPrev += "&start=" + strconv.Itoa(currStart-currNum)
					// PreviousPage will be nil on first page of results
				} else {

				}
			} else { // if there isn't a valid start query string, set starting position to 0 in both template's args and query string
				*aStart = 0
				*aNext += "&start=" + strconv.Itoa(0+DefaultNum)
				*aCurr += "&start=" + strconv.Itoa(0)
				*aPrev = ""
			}
		} else { // if missing a valid "num" query string, set everything to defaults (start at 0, DefaultNum words per page)
			*aNum = DefaultNum
			*aStart = 0
			// create the default nextpage query string for the "next page" link
			*aNext = *aSort + "&num=" + strconv.Itoa(DefaultNum) + "&start=" + strconv.Itoa(0+DefaultNum)
			*aCurr = *aSort + "&num=" + strconv.Itoa(DefaultNum) + "&start=" + strconv.Itoa(0)
		}

		if r.URL.Query()["id"] != nil {
			// if it's set to sort by "new"
			id := r.URL.Query()["id"][0]
			dbconnect.FlagRow(id)
		}

		// see if next page button should be hidden before fetching rows
		numrows := dbconnect.CountRows(t.queryFrom)
		if numrows < *aStart+*aNum {
			*aNext = ""
		}
		if numrows < *aStart {
			*aPrev = ""
		}

		// join the template files for the wordlist
		t.templ = template.Must(template.New(t.filename).Funcs(funcMap).ParseFiles(
			filepath.Join("templates", t.filename),
			filepath.Join("templates", "words_sorted.html"),
			filepath.Join("templates", "_header.html"),
			filepath.Join("templates", "_footer.html"),
		))

	} else {
		// not word list, just a regular page, join the template files
		t.templ = template.Must(template.New(t.filename).Funcs(funcMap).ParseFiles(
			filepath.Join("templates", t.filename),
			filepath.Join("templates", "_header.html"),
			filepath.Join("templates", "_footer.html"),
		))
	}
	t.templ.Execute(w, t.args)
}
