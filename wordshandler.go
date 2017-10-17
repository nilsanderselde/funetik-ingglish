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
	// pointers to TemplateParams fields for more readable code
	aNew, aOld, aDist, aID := &t.args.New, &t.args.Old, &t.args.Dist, &t.args.ID
	aRev, aQ, aStart, aNum := &t.args.Reverse, &t.args.Query, &t.args.Start, &t.args.Num
	aSort, aCurr, aNext, aPrev := &t.args.Sort, &t.args.CurrentPage, &t.args.NextPage, &t.args.PreviousPage

	// put two pieces of query together
	*aQ = t.query + t.queryFrom

	// reset prev page to force templateHandler to recreate it if needed
	*aPrev = ""

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
			*aQ += " ORDER BY dist"
		} else {
			*aNew, *aOld, *aDist, *aID = true, false, false, false
			*aSort = "?sortby=new"
			*aQ += " ORDER BY funsort"
		}
		// ascending or descending
		if r.URL.Query()["order"] != nil && r.URL.Query()["order"][0] == "desc" {
			*aQ += " DESC, id DESC;"
			*aSort += desc
			*aRev = true
		} else {
			*aQ += " ASC, id ASC;"
			*aSort += asc
			*aRev = false
		}
	} else {
		// default values: sort by new spelling, ascending, starting at zero, incrementing by DefaultNum
		*aNew, *aOld, *aDist, *aID = true, false, false, false
		*aSort = "?sortby=new&order=asc"
		*aRev = false
		*aQ += " ORDER BY funsort ASC, id asc;"
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
}
