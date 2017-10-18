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

	// put two pieces of postgres query together
	t.args.PQuery = t.query + t.queryFrom

	// reset prev page to force templateHandler to recreate it if needed
	t.args.PreviousPage = ""

	urlQ := r.URL.Query()

	// // if set to update automatically generated values
	{ // if urlQ["updeit"] != nil {
		// 	if urlQ["updeit"][0] == "al" {
		// 		dbconnect.UpdateAllAutoValues(fun, numsil, funsort, dist, false)
		// 	} else if urlQ["updeit"][0] == "flaagd" {
		// 		dbconnect.UpdateAllAutoValues(fun, numsil, funsort, dist, true)
		// 	}
		// }
	}

	// if set to sort by value
	if urlQ["sortby"] != nil {
		// sort by id, truditional spelling, distance, or default (funsort)
		columns := []string{"id", "trud", "dist"}
		sortby := "funsort"
		for _, s := range columns {
			if urlQ["sortby"][0] == s {
				sortby = s
			}
		}
		t.args.SortBy = sortby
		t.args.SortQ = "?sortby=" + sortby
		t.args.PQuery += " ORDER BY " + sortby
	}
	// ascending or descending
	if urlQ["order"] != nil {
		reverse := "desc"
		order := "asc"
		if urlQ["order"][0] == reverse {
			order = reverse
			t.args.Reverse = true
		} else {
			t.args.Reverse = false
		}
		t.args.PQuery += " " + order + ", id " + order + ";"
		t.args.SortQ += "&order=" + order
	}

	// if there is a valid "num" query string...
	if urlQ["num"] != nil {
		// ...get the current value of it, preventing massive queries
		currNum, err := strconv.Atoi(urlQ["num"][0])
		if err != nil || currNum > 1000 {
			currNum = DefaultNum
		}
		// ...set the template's args object's value to it
		t.args.Num = currNum

		// create the next page query string for the "next page" link
		t.args.NextPage = t.args.SortQ + "&num=" + strconv.Itoa(currNum)
		t.args.CurrentPage = t.args.NextPage

		// if there is a valid "start" query string...
		if urlQ["start"] != nil {

			// ...get the current value of it
			currStart, err := strconv.Atoi(urlQ["start"][0])
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
				t.args.PreviousPage = t.args.SortQ + "&num=" + strconv.Itoa(currNum)
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
		t.args.NextPage = t.args.SortQ + "&num=" + strconv.Itoa(DefaultNum) + "&start=" + strconv.Itoa(0+DefaultNum)
		t.args.CurrentPage = t.args.SortQ + "&num=" + strconv.Itoa(DefaultNum) + "&start=" + strconv.Itoa(0)
	}

	if urlQ["id"] != nil {
		// if it's set to sort by "new"
		id := urlQ["id"][0]
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
