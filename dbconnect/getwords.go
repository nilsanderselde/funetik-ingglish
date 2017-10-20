// Nils Elde
// https://gitlab.com/nilsanderselde

package dbconnect

import (
	"strconv"

	"gitlab.com/nilsanderselde/funetik-ingglish/global"

	// postgres drivers
	_ "github.com/lib/pq"
)

var busy bool
var wordList [][]string

// ShowWords calls the database io function and passes arguments to template
func ShowWords(args global.TemplateParams) [][]string {
	// fmt.Println("HTTP request")
	for busy {
		busy = false
		// fmt.Println("Too many requests. Here's the cached version.")
		return wordList
	}
	return GetWords(args.PQuery, args.Start, args.Num)
}

// GetWords receives SQL requests and returns requested data from the database. This package will not be shared in the repository for security reasons.
// returns true if more results can be found (to avoid next button to empty page)
func GetWords(query string, start int, num int) [][]string {
	busy = true
	// fmt.Println("Getting words")

	results := [][]string{}
	noresults := [][]string{{"-1", "-1", "-1", "-1", "-1", "-1", "-1", "-1"}}

	// get query
	rows, err := DB.Query(query)
	if err != nil {
		goto loadfailure
	}
	defer rows.Close()

	for i := 0; rows.Next() && i < start+num; i++ {
		if i >= start {
			var id int
			var fun string
			var funsil string
			var trud string
			var pus string
			var numsil int
			var dist int
			var funsort string
			var fläg bool
			err = rows.Scan(&id, &fun, &funsil, &trud, &pus, &numsil, &dist, &funsort, &fläg)
			results = append(results, []string{strconv.Itoa(id), fun, funsil, trud, strconv.Itoa(dist),
				pus, strconv.Itoa(numsil), strconv.FormatBool(fläg)})
			if err != nil {
				goto loadfailure
			}
		}
	}
	err = rows.Err()
	if err != nil {
		goto loadfailure
	}
	// Only return results if they are all good, otherwise return array with one row of -1s
	if len(results) != 0 {
		busy = false
		wordList = results
		return results
	}
loadfailure:
	return noresults
}

// CountRows counts the number of rows returned by a query that is passed to it
func CountRows() int {

	// get total number of rows
	rowcount, err := DB.Query("SELECT COUNT(*) from words;")
	if err != nil {
		// log.Fatal(err)
		return -1
	}
	numrows := 0
	for rowcount.Next() {
		err := rowcount.Scan(&numrows)
		if err != nil {
			// log.Fatal(err)
			return -1
		}
	}
	return numrows
}
