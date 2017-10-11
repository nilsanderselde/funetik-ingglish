package dbconnect

import (
	"database/sql"

	"gitlab.com/nilsanderselde/funetik-ingglish/wordtools"
)

// GetStats connects to the database, gets all the
// funetik spellings of words, sends the list to
// CalculateStats, which counts the letters,
// and then returns the stats
func GetStats() [][]string {

	notfound := [][]string{{"-1", "-1", "-1", "-1", "-1", "-1", "-1", "-1"}}

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

	rows, err := db.Query("SELECT fun FROM words;")
	if err != nil {
		// log.Fatal(err)
		return notfound
	}
	defer rows.Close()

	words := []string{}
	for rows.Next() {
		var word string
		err := rows.Scan(&word)
		if err != nil {
			// log.Fatal(err)
			return notfound
		}
		words = append(words, word)
	}
	err = rows.Err()
	if err != nil {
		// log.Fatal(err)
		return notfound
	}

	return wordtools.CalculateStats(words)
}
