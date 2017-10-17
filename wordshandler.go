// Nils Elde
// https://gitlab.com/nilsanderselde
// This file contains the word list template logic

package main

import (
	"net/http"
	"strconv"

	"gitlab.com/nilsanderselde/funetik-ingglish/dbconnect"
)

const (
	// DefaultNum is default number of words per page
	DefaultNum int = 32
)

// Options for ?update=all (func UpdateAllAutoValues)
const (
	fun        bool = true
	numsil     bool = true
	funsort    bool = true
	dist       bool = true
	onlyFlaagd bool = true
)

func handleWordList(t *templateHandler, r *http.Request) {

	// put two pieces of query together
	t.args.Query = t.query + t.queryFrom

	// reset prev page to force templateHandler to recreate it if needed
	t.args.PreviousPage = ""

	// // if set to update automatically generated values
	// if r.URL.Query()["updeit"] != nil {
	// 	if r.URL.Query()["updeit"][0] == "al" {
	// 		dbconnect.UpdateAllAutoValues(fun, numsil, funsort, dist, false)
	// 	} else if r.URL.Query()["updeit"][0] == "flaagd" {
	// 		dbconnect.UpdateAllAutoValues(fun, numsil, funsort, dist, true)
	// 	}
	// }

	// if set to sort by value
	if r.URL.Query()["sortby"] != nil {

		desc := "&order=desc"
		asc := "&order=asc"

		// sort by id, old spelling, distance, or default (new spelling)
		if r.URL.Query()["sortby"][0] == "id" {
			t.args.New, t.args.Old, t.args.Dist, t.args.ID = false, false, false, true
			t.args.Sort = "?sortby=id"
			t.args.Query += " ORDER BY id"
		} else if r.URL.Query()["sortby"][0] == "old" {
			t.args.New, t.args.Old, t.args.Dist, t.args.ID = false, true, false, false
			t.args.Sort = "?sortby=old"
			t.args.Query += " ORDER BY trud"
		} else if r.URL.Query()["sortby"][0] == "dist" {
			t.args.New, t.args.Old, t.args.Dist, t.args.ID = false, false, true, false
			t.args.Sort = "?sortby=dist"
			t.args.Query += " ORDER BY dist"
		} else {
			t.args.New, t.args.Old, t.args.Dist, t.args.ID = true, false, false, false
			t.args.Sort = "?sortby=new"
			t.args.Query += " ORDER BY funsort"
		}
		// ascending or descending
		if r.URL.Query()["order"] != nil && r.URL.Query()["order"][0] == "desc" {
			t.args.Query += " DESC, id DESC;"
			t.args.Sort += desc
			t.args.Reverse = true
		} else {
			t.args.Query += " ASC, id ASC;"
			t.args.Sort += asc
			t.args.Reverse = false
		}
	} else {
		// default values: sort by new spelling, ascending, starting at zero, incrementing by DefaultNum
		t.args.New, t.args.Old, t.args.Dist, t.args.ID = true, false, false, false
		t.args.Sort = "?sortby=new&order=asc"
		t.args.Reverse = false
		t.args.Query += " ORDER BY funsort ASC, id asc;"
	}
	// if there is a valid "num" query string...
	if r.URL.Query()["num"] != nil {
		// ...get the current value of it, preventing massive queries
		currNum, err := strconv.Atoi(r.URL.Query()["num"][0])
		if err != nil || currNum > 1000 {
			currNum = DefaultNum
		}
		// ...set the template's args object's value to it
		t.args.Num = currNum

		// create the next page query string for the "next page" link
		t.args.NextPage = t.args.Sort + "&num=" + strconv.Itoa(currNum)
		t.args.CurrentPage = t.args.NextPage

		// if there is a valid "start" query string...
		if r.URL.Query()["start"] != nil {

			// ...get the current value of it
			currStart, err := strconv.Atoi(r.URL.Query()["start"][0])
			if err != nil {
				currStart = 0
			}
			// ...set the template's args object's value to it
			t.args.Start = currStart

			// add to the next page query string
			t.args.NextPage += "&start=" + strconv.Itoa(currStart+currNum)
			t.args.CurrentPage += "&start=" + strconv.Itoa(currStart)
			// if curr starting number >= num of words per page, add prev page button
			if currStart >= currNum {
				t.args.PreviousPage = t.args.Sort + "&num=" + strconv.Itoa(currNum)
				t.args.PreviousPage += "&start=" + strconv.Itoa(currStart-currNum)
				// PreviousPage will be nil on first page of results
			} else {

			}
		} else { // if there isn't a valid start query string, set starting position to 0 in both template's args and query string
			t.args.Start = 0
			t.args.NextPage += "&start=" + strconv.Itoa(0+DefaultNum)
			t.args.CurrentPage += "&start=" + strconv.Itoa(0)
			t.args.PreviousPage = ""
		}
	} else { // if missing a valid "num" query string, set everything to defaults (start at 0, DefaultNum words per page)
		t.args.Num = DefaultNum
		t.args.Start = 0
		// create the default nextpage query string for the "next page" link
		t.args.NextPage = t.args.Sort + "&num=" + strconv.Itoa(DefaultNum) + "&start=" + strconv.Itoa(0+DefaultNum)
		t.args.CurrentPage = t.args.Sort + "&num=" + strconv.Itoa(DefaultNum) + "&start=" + strconv.Itoa(0)
	}

	if r.URL.Query()["id"] != nil {
		// if it's set to sort by "new"
		id := r.URL.Query()["id"][0]
		dbconnect.FlagRow(id)
	}

	// see if next page button should be hidden before fetching rows
	numrows := dbconnect.CountRows(t.queryFrom)
	if numrows < t.args.Start+t.args.Num {
		t.args.NextPage = ""
	}
	if numrows < t.args.Start {
		t.args.PreviousPage = ""
	}
}
