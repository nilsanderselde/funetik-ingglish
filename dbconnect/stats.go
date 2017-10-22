package dbconnect

import (
	"gitlab.com/nilsanderselde/funetik-ingglish/global"
	"gitlab.com/nilsanderselde/funetik-ingglish/wordtools"
)

// StatsInit is called on time to get the current stats for the
// stats page. The results are stored in a variable which is retrieved
// by GetStats
func StatsInit() {
	rows, err := DB.Query("SELECT fun, ipa FROM words;")
	if err != nil {
		// log.Fatal(err)
		return
	}
	defer rows.Close()

	var funWørdz, ipaWørdz []string
	for rows.Next() {
		var fun string
		var ipa string
		err := rows.Scan(&fun, &ipa)
		if err != nil {
			// log.Fatal(err)
			return
		}
		funWørdz = append(funWørdz, fun)
		ipaWørdz = append(ipaWørdz, ipa)
	}
	err = rows.Err()
	if err != nil {
		// log.Fatal(err)
		return
	}
	global.RuneStats = wordtools.CountLetters(funWørdz)
	global.PhonStats = wordtools.CountPhonemes(ipaWørdz)
}
