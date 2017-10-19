// Nils Elde
// https://gitlab.com/nilsanderselde
// This file contains the word list template logic

package main

import (
	"net/http"
	"strconv"

	"gitlab.com/nilsanderselde/funetik-ingglish/global"

	"gitlab.com/nilsanderselde/funetik-ingglish/dbconnect"
)

const (
	// DefaultNum is default number of words per page
	DefaultNum int = 32
)

// Options for dbconnect.Update...AutoValues
const (
	fun     bool = true
	numsil  bool = true
	funsort bool = true
	dist    bool = true
)

func handleWordList(t *templateHandler, r *http.Request) {

	// SQL query split into two parts because the FROM section is used twice
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
	wordsQueryFrom := "FROM words"
	t.args.PQuery = wordsQuery + wordsQueryFrom

	// get URL query string
	urlQ := r.URL.Query()

	// if set to update automatically generated values
	if global.IsDev && urlQ["updeit"] != nil {
		if urlQ["updeit"][0] == "al" {
			dbconnect.UpdateAllAutoValues(fun, numsil, funsort, dist)
		} else if urlQ["updeit"][0] == "flaagd" {
			dbconnect.UpdateFlaggedAutoValues(fun, numsil, funsort, dist)
		}
	}
	// flag row if URL contains query string for id
	if global.IsDev && urlQ["id"] != nil {
		dbconnect.FlagRow(urlQ["id"][0])
	}

	// column to sort by
	sortby := "funsort"
	if urlQ["sortby"] != nil {
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
	t.args.NextPage += "&start=" + strconv.Itoa(start+num)
	t.args.CurrentPage += "&start=" + strconv.Itoa(start)

	// if there is a previous page, create the query string for the link to it
	if start >= num {
		t.args.PreviousPage = t.args.SortQ + "&num=" + strconv.Itoa(num) + "&start=" + strconv.Itoa(start-num)
	} else {
		t.args.PreviousPage = ""
	}

	// see if next page link should be hidden because there's no more results
	numrows := dbconnect.CountRows(wordsQueryFrom)
	if numrows < start+num {
		t.args.NextPage = ""
	}
	// if num rows returned less than start num, disable previous page link
	if numrows < start {
		t.args.PreviousPage = ""
	}
}
