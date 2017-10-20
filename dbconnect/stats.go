package dbconnect

import (
	"gitlab.com/nilsanderselde/funetik-ingglish/global"
	"gitlab.com/nilsanderselde/funetik-ingglish/wordtools"
)

// StatsInit is called on time to get the current stats for the
// stats page. The results are stored in a variable which is retrieved
// by GetStats
func StatsInit() {
	rows, err := DB.Query("SELECT fun FROM words;")
	if err != nil {
		// log.Fatal(err)
		return
	}
	defer rows.Close()

	words := []string{}
	for rows.Next() {
		var word string
		err := rows.Scan(&word)
		if err != nil {
			// log.Fatal(err)
			return
		}
		words = append(words, word)
	}
	err = rows.Err()
	if err != nil {
		// log.Fatal(err)
		return
	}
	global.Stats = wordtools.CalculateStats(words)
}
