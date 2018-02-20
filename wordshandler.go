// Nils Elde
// https://github.com/nilsanderselde
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
	DefaultNum int = 20
)

// Options for dbconnect.Update...AutoValues
const (
	fun     bool = true
	numsil  bool = true
	funsort bool = true
	dist    bool = true
	ipa     bool = true
	ritin   bool = true
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
	fullQuery := wordsQuery + wordsQueryFrom
	// get URL query string
	urlQ := r.URL.Query()

	// if set to update automatically generated values
	if global.IsDev && urlQ["updeit"] != nil {
		if urlQ["updeit"][0] == "al" {
			dbconnect.UpdateAllAutoValues(fun, numsil, funsort, dist, ipa, ritin)
		} else if urlQ["updeit"][0] == "flaagd" {
			dbconnect.UpdateFlaggedAutoValues(fun, numsil, funsort, dist, ipa, ritin)
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
	sortQ := "?sortby=" + sortby
	fullQuery += " ORDER BY " + sortby

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
	fullQuery += " " + order
	switch sortby {
	case "id":
		fullQuery += ", funsort " + order + ", trud " + order
	case "trud":
		fullQuery += ", funsort " + order + ", id " + order
	case "dist":
		fullQuery += ", funsort " + order + ", trud " + order + ", id " + order
	default:
		fullQuery += ", trud " + order + ", id " + order
	}
	fullQuery += ";"
	sortQ += "&order=" + order

	// number of words per page
	num := DefaultNum
	if urlQ["num"] != nil {
		curr, err := strconv.Atoi(urlQ["num"][0])
		if err == nil && curr <= 100 {
			num = curr
		}
	}
	t.args.Num = num

	// offset from beginning of results
	var start int
	if urlQ["start"] != nil {
		curr, err := strconv.Atoi(urlQ["start"][0])
		if err == nil {
			start = curr
		}
	}
	t.args.Start = start

	pagePrefix := sortQ + "&num=" + strconv.Itoa(num) + "&start="
	t.args.CurrentPage = pagePrefix + strconv.Itoa(start)

	// if there is a previous page, create the query string for the link to it
	if start >= num && global.RowCount >= start {
		t.args.PreviousPage = pagePrefix + strconv.Itoa(start-num)
	} else {
		t.args.PreviousPage = ""
		// t.args.PreviousPage = pagePrefix + "0"
	}

	// see if next page link should be hidden because there's no more results
	if global.RowCount > start+num {
		t.args.NextPage = pagePrefix + strconv.Itoa(start+num)
	} else {
		t.args.NextPage = ""
	}

	t.args.PQuery = fullQuery
	t.args.SortQ = sortQ
}
