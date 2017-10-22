// Nils Elde
// https://gitlab.com/nilsanderselde

package dbconnect

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"gitlab.com/nilsanderselde/funetik-ingglish/wordtools"
)

// UpdateAllAutoValues automatically generates values for all rows
// (About 15 minutes with all tasks enabled and 50,000 words)
func UpdateAllAutoValues(fun bool, numsil bool, funsort bool, dist bool, ipa bool, ritin bool) {
	UpdateAutoValues(fun, numsil, funsort, dist, ipa, ritin, false, -1)
}

// UpdateFlaggedAutoValues automatically generates values for all flaagd rows
// (About 15 minutes with all tasks enabled and 50,000 words)
func UpdateFlaggedAutoValues(fun bool, numsil bool, funsort bool, dist bool, ipa bool, ritin bool) {
	UpdateAutoValues(fun, numsil, funsort, dist, ipa, ritin, true, -1)
}

// UpdateAutoValues automatically fills values for:
// 1. fun      (remove syllable markings from funsil)
// 2. numsil   (count syllable markings from funsil)
// 3. funsort  (substitution cipher on fun)
// 4. dist     (calc lev dist between fun and trud)
// onlyFlaagd means only flaagd will be processed
// rowID used to specify row to update, use -1 for all
func UpdateAutoValues(fun bool, numsil bool, funsort bool, dist bool, ipa bool, ritin bool, onlyFlaagd bool, rowID int) {
	if !(fun || numsil || funsort || dist || ipa || ritin) {
		return // don't bother connecting if no updates will occur
	}

	queryFrom := " FROM words"
	if rowID != -1 || onlyFlaagd {
		queryFrom += " WHERE"
		if rowID != -1 {
			queryFrom += " id = " + strconv.Itoa(rowID)
			if onlyFlaagd {
				queryFrom += " AND"
			}
		}
		if onlyFlaagd {
			queryFrom += " flaagd"
		}
	}
	queryFrom += ";"

	if fun || numsil {
		start := time.Now()
		message := "Updating"
		if fun {
			message += " fun"
			if numsil {
				message += " and"
			}
		}
		if numsil {
			message += " numsil"
		}
		fmt.Print(message + "... ")
		s := spinner.New(spinner.CharSets[13], 100*time.Millisecond)
		s.Start()

		// update fun and/or numsil with values generated using funsil
		rows, err := DB.Query("SELECT id, funsil" + queryFrom)
		if err != nil {
			log.Println(err)
		}

		for rows.Next() {
			if fun {
				updateFun(rows)
			}
			if numsil {
				updateNumsil(rows)
			}
		}
		rows.Close()
		t := time.Now()
		elapsed := t.Sub(start)
		s.Stop()
		fmt.Printf("Done. (%v)\n", elapsed)
	}

	if funsort || dist {
		start := time.Now()
		message := "Updating"
		if funsort {
			message += " funsort"
			if dist {
				message += " and"
			}
		}
		if dist {
			message += " dist"
		}
		fmt.Print(message + "... ")
		s := spinner.New(spinner.CharSets[13], 100*time.Millisecond)
		s.Start()

		// update funsort and dist with values generated using fun and trud (for dist, use written form if different)
		rows, err := DB.Query("SELECT id, fun, trud, COALESCE(COALESCE(ritin, fun), '') as funritin" + queryFrom)
		if err != nil {
			log.Println(err)
		}

		for rows.Next() {
			if funsort {
				updateFunsort(rows)
			}
			if dist {
				updateDist(rows)
			}
		}
		rows.Close()
		t := time.Now()
		elapsed := t.Sub(start)
		s.Stop()
		fmt.Printf("Done. (%v)\n", elapsed)

		// clear all flags
		_, err = DB.Exec("UPDATE words SET flaagd = false WHERE flaagd;")
		if err != nil {
			log.Println(err)
		}
		if ipa {
			updateIPA()
		}
		if ritin {
			updateRitin()
		}
	}
}

// Update IPA generates the rough IPA representation of words by
// consolidating diphthongs and digraphs into single characters,
// and by splitting y and w into semivowels and vowels
func updateIPA() {
	fmt.Println("Updating ipa...")

	_, err := DB.Exec(`update words set ipa = funsil;
	
update words set ipa = regexp_replace(funsil, 'ng', 'ŋ')
where id in (select id from words where funsil similar TO '%ng%');

update words set ipa = regexp_replace(funsil, 'dh', 'ð')
where id in (select id from words where funsil similar TO '%dh%');

update words set ipa = regexp_replace(funsil, 'th', 'θ')
where id in (select id from words where funsil similar TO '%th%');

update words set ipa = regexp_replace(funsil, 'dž', 'ʤ')
where id in (select id from words where funsil similar TO '%dž%');

update words set ipa = regexp_replace(funsil, 'tš', 'ʧ')
where id in (select id from words where funsil similar TO '%tš%');

update words set ipa = regexp_replace(funsil, 'ai', 'ī')
where id in (select id from words where funsil similar TO '%ai%');

update words set ipa = regexp_replace(funsil, 'aw', 'ã')
where id in (select id from words where funsil similar TO '%aw%');

update words set ipa = regexp_replace(funsil, 'ei', 'ā')
where id in (select id from words where funsil similar TO '%ei%');

update words set ipa = regexp_replace(funsil, 'oi', 'õ')
where id in (select id from words where funsil similar TO '%oi%');

update words set ipa = regexp_replace(funsil, 'y([aäeoøiu])', 'j\1')
where id in (select id from words where funsil similar TO '%y[aäeoøiu]%');

update words set ipa = regexp_replace(funsil, 'w([aäeoøiu])', 'ʍ\1')
where id in (select id from words where funsil similar TO '%w[aäeoøiu]%');`)

	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("Done.")
	}

}

// updateRitin updates the written form of words that that differ from their pronounced form
// where it can be easily derived (in the case of syllable boundaries between d/t and h, insert dash)
func updateRitin() {
	fmt.Println("Updating ritin...")
	_, err := DB.Exec(`update words set ritin = regexp_replace(fun, 'th', 't-h') where funsil similar to '%t[ˈˌ·]h%';
update words set ritin = regexp_replace(fun, 'dh', 'd-h') where funsil similar to '%d[ˈˌ·]h%';`)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("Done.")
	}

}

// updateFun updates passed row with generated
// value for column "fun" (funetik spelling)
func updateFun(row *sql.Rows) {
	var id int
	var funsil string
	err := row.Scan(&id, &funsil)
	// Generate value for "fun"
	fun := wordtools.RemoveSyllableMarkers(funsil)

	// update with new fun value
	var updateString string
	updateFun := DB.QueryRow("UPDATE words SET fun = '" + fun + "' WHERE id = " + strconv.Itoa(id) + ";")
	err = updateFun.Scan(&updateString)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
	}
}

// updateNumsil updates passed row with generated
// value for column  "numsil" (number of syllables)
func updateNumsil(row *sql.Rows) {
	var id int
	var funsil string
	err := row.Scan(&id, &funsil)

	// Generate value for "numsil"
	numsil := wordtools.CountSyllables(funsil)

	// update with new numsil value
	var updateString string
	updateNumsil := DB.QueryRow("UPDATE words SET numsil = '" + strconv.Itoa(numsil) + "' WHERE id = " + strconv.Itoa(id) + ";")
	err = updateNumsil.Scan(&updateString)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
	}
}

// updateFunsort updates passed row with generated
// value for column "funsort" (funetik sort order)
func updateFunsort(row *sql.Rows) {
	var id int
	var fun string
	var trud string
	var funritin string
	err := row.Scan(&id, &fun, &trud, &funritin)

	// Generate value for "funsort"
	funsort := wordtools.SubstitutionCypher(fun)

	// update with new funsort value
	var updateString string
	updateFunsort := DB.QueryRow("UPDATE words SET funsort = '" + funsort + "' WHERE id = " + strconv.Itoa(id) + ";")
	err = updateFunsort.Scan(&updateString)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
	}
}

// updateDist updates passed row with calculated
// value for column "levdist" (Levenshtein distance)
func updateDist(row *sql.Rows) {
	var id int
	var fun string
	var trud string
	var funritin string
	err := row.Scan(&id, &fun, &trud, &funritin)

	// Generate value for "dist" (true means flipping two adjacent letters is considered one move)
	dist := wordtools.FindDistance(funritin, trud, true)

	// update with new dist value
	var updateString string
	updateDist := DB.QueryRow("UPDATE words SET dist = '" + strconv.Itoa(dist) + "' WHERE id = " + strconv.Itoa(id) + ";")
	err = updateDist.Scan(&updateString)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
	}
}
