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

	// column to sort by
	sortby := "funsort"
	if urlQ["sortby"] != nil {
		// sort by id, truditional spelling, distance, or default (funsort)
		columns := []string{"id", "trud", "dist"}
		for _, s := range columns {
			if urlQ["sortby"][0] == s {
				sortby = s
			}
		}
	}
	t.args.SortBy = sortby
	t.args.SortQ = "?sortby=" + sortby
	t.args.PQuery += " ORDER BY " + sortby

	// ascending or descending
	order := "asc"
	if urlQ["order"] != nil {
		reverse := "desc"
		if urlQ["order"][0] == reverse {
			order = reverse
			t.args.Reverse = true
		} else {
			t.args.Reverse = false
		}
	}
	t.args.PQuery += " " + order + ", id " + order + ";"
	t.args.SortQ += "&order=" + order

	// number of words per page
	num := DefaultNum
	if urlQ["num"] != nil {
		curr, err := strconv.Atoi(urlQ["num"][0])
		if err == nil || curr <= 100 {
			num = curr
		}
	}
	t.args.Num = num
	t.args.NextPage = t.args.SortQ + "&num=" + strconv.Itoa(num)
	t.args.CurrentPage = t.args.NextPage

	// offset from beginning of results
	var start int
	if urlQ["start"] != nil {
		curr, err := strconv.Atoi(urlQ["start"][0])
		if err == nil {
			start = curr
		}
	}
	t.args.Start = start

	// set next page and current page query strings
	t.args.NextPage += "&start=" + strconv.Itoa(start+num)
	t.args.CurrentPage += "&start=" + strconv.Itoa(start)

	// if there is a previous page, create the query string for the link to it
	if start >= num {
		t.args.PreviousPage = t.args.SortQ + "&num=" + strconv.Itoa(num) + "&start=" + strconv.Itoa(start-num)
	} else {
		t.args.PreviousPage = ""
	}

	// // flag row if URL contains query string for id
	// if urlQ["id"] != nil {
	// 	dbconnect.FlagRow(urlQ["id"][0])
	// }

	// see if next page button should be hidden because there's no more results
	numrows := dbconnect.CountRows(t.queryFrom)
	if numrows < start+num {
		t.args.NextPage = ""
	}
	// if the number of rows returned is less than the starting number, the starting number is too high
	// and backwards results navigation should also be disabled (this would only achievable by manually
	// entering a huge starting number, but due to the growing nature of the database, setting a hard limit
	// to start doesn't make sense)
	if numrows < start {
		t.args.PreviousPage = ""
	}
}
