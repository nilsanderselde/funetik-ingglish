package dbconnect

import (
	"gitlab.com/nilsanderselde/funetik-ingglish/wordtools"
)

var stats [][]string

// GetStats connects to the database, gets all the
// funetik spellings of words, sends the list to
// CalculateStats, which counts the letters,
// and then returns the stats
func GetStats() [][]string {
	return stats
}

// StatsInit is called on time to get the current stats for the
// stats page. The results are stored in a variable which is retrieved
// by GetStats
func StatsInit() {
	stats = [][]string{{"-1", "-1", "-1", "-1", "-1", "-1"}}

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
	stats = wordtools.CalculateStats(words)
}
