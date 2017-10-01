// Nils Elde
// https://gitlab.com/nilsanderselde

package dbconnect

import (
	"database/sql"
	"strconv"

	// postgres drivers
	_ "github.com/lib/pq"
	"gitlab.com/nilsanderselde/funetik-ingglish/global"
)

// ShowWords calls the database io function and passes arguments to template
func ShowWords(args global.TemplateParams) [][]string {
	return GetWords(args.Query, args.Start, args.Num)
}

// GetWords receives SQL requests and returns requested data from the database. This package will not be shared in the repository for security reasons.
// returns true if more results can be found (to avoid next button to empty page)
func GetWords(query string, start int, num int) [][]string {

	results := [][]string{}
	notfound := [][]string{{"-1", "-1", "-1", "-1", "-1", "-1", "-1", "-1"}}

	// fmt.Println(query + "\nStarting at " + strconv.Itoa(start) + " with " + strconv.Itoa(num) + " results per page.")

	db, err := sql.Open("postgres", DBInfo)
	if err != nil {
		// log.Fatal(err)
		return notfound
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		// log.Fatal(err)
		return notfound
	}

	// get query
	rows, err := db.Query(query)
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
			// fmt.Printf("%s %s %s %s %s %s %s %s\n", strconv.Itoa(id), fun, funsil, trud, strconv.Itoa(dist), pus, strconv.Itoa(numsil), strconv.FormatBool(fläg))
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

	err = db.Ping()
	if err != nil {
		// log.Fatal(err)
		return -1
	}

	// get total number of rows
	rowcount, err := db.Query("SELECT COUNT(*)" + queryFrom + ";")
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
