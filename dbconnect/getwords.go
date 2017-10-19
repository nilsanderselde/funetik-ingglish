// Nils Elde
// https://gitlab.com/nilsanderselde

package dbconnect

import (
	"database/sql"
	"strconv"
	"time"

	// postgres drivers
	_ "github.com/lib/pq"
	"gitlab.com/nilsanderselde/funetik-ingglish/global"
)

var ready = true
var wordList = [][]string{{"-1", "-1", "-1", "-1", "-1", "-1", "-1", "-1"}}

// ShowWords calls the database io function and passes arguments to template
func ShowWords(args global.TemplateParams) [][]string {
	// fmt.Println("HTTP request")
	for !ready {
		time.Sleep(1000 * time.Millisecond)
		ready = true
		// fmt.Println("Too many requests. Here's the cached version.")
		return wordList
	}
	return GetWords(args.PQuery, args.Start, args.Num)
}

// GetWords receives SQL requests and returns requested data from the database. This package will not be shared in the repository for security reasons.
// returns true if more results can be found (to avoid next button to empty page)
func GetWords(query string, start int, num int) [][]string {
	ready = false
	// fmt.Println("Getting words")

	results := [][]string{}
	notfound := [][]string{{"-1", "-1", "-1", "-1", "-1", "-1", "-1", "-1"}}

	// get query
	rows, err := DB.Query(query)
	if err != nil {
		// log.Fatal(err)
		return notfound
	}
	defer rows.Close()

	i := 0
	for rows.Next() && i < start+num {
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
				// log.Fatal(err)
				return notfound
			}
		}
		i++
	}
	err = rows.Err()
	if err != nil {
		// log.Fatal(err)
		return notfound
	}
	// Only return results if they are all good, otherwise return array with one row of -1s
	if len(results) != 0 {
		ready = true
		wordList = results
		return results
	}
	return notfound
}

// CountRows counts the number of rows returned by a query that is passed to it
func CountRows(queryFrom string) int {

	db, err := sql.Open("postgres", DBInfo)
	if err != nil {
		// log.Fatal(err)
		return -1
	}
	defer db.Close()

	err = DB.Ping()
	if err != nil {
		// log.Fatal(err)
		return -1
	}

	// get total number of rows
	rowcount, err := DB.Query("SELECT COUNT(*)" + queryFrom + ";")
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
