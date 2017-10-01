package dbconnect

import (
	"database/sql"
	"log"

	"gitlab.com/nilsanderselde/funetik-ingglish/wordtools"
)

// GetStats connects to the database, gets all the
// funetik spellings of words, sends the list to
// CalculateStats, which counts the letters,
// and then returns the stats
func GetStats() [][]string {
	db, err := sql.Open("postgres", DBInfo)
	if err != nil {
		log.Fatal(err)
		// fmt.Println(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		// fmt.Println(err)
	}

	rows, err := db.Query("SELECT fun FROM words;")
	if err != nil {
		log.Fatal(err)
		// fmt.Println(err)
	}
	defer rows.Close()

	words := []string{}
	for rows.Next() {
		var word string
		err := rows.Scan(&word)
		if err != nil {
			log.Fatal(err)
		}
		words = append(words, word)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return wordtools.CalculateStats(words)
}
