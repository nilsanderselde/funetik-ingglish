// Nils Elde
// https://github.com/nilsanderselde

package dbconnect

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
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
// 5. ipa      (substitute digraphs for single characters)
// 6. ritin    (generates written form for certain words)
// onlyFlaagd means only flaagd will be processed
// rowID used to specify row to update, use -1 for all
func UpdateAutoValues(fun bool, numsil bool, funsort bool, dist bool, ipa bool, ritin bool, onlyFlaagd bool, rowID int) {
	if !(fun || numsil || funsort || dist || ipa || ritin) {
		return // don't bother connecting if no updates will occur
	}

	var queryEnd string
	if rowID != -1 || onlyFlaagd {
		if rowID != -1 {
			queryEnd += " id = " + strconv.Itoa(rowID)
			if onlyFlaagd {
				queryEnd += " AND"
			}
		}
		if onlyFlaagd {
			queryEnd += " flaagd"
		}
	}
	queryEnd += ";"

	if fun || numsil {
		var processes string
		if fun {
			processes += " fun"
			if numsil {
				processes += " and"
			}
		}
		if numsil {
			processes += " numsil"
		}
		p := updateBegin(processes)

		qWhere := queryEnd
		if queryEnd != ";" {
			qWhere = " WHERE" + queryEnd
		}
		// update fun and/or numsil with values generated using funsil
		rows, err := DB.Query("SELECT id, funsil FROM words" + qWhere)
		if err != nil {
			log.Fatal(err)
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
		updateEnd(p)
	}

	if funsort || dist {
		var processes string
		if funsort {
			processes += " funsort"
			if dist {
				processes += " and"
			}
		}
		if dist {
			processes += " dist"
		}
		p := updateBegin(processes)

		qWhere := queryEnd
		if queryEnd != ";" {
			qWhere = " WHERE" + queryEnd
		}
		// update funsort and dist with values generated using fun and trud (for dist, use written form if different)
		rows, err := DB.Query("SELECT id, fun, trud, COALESCE(COALESCE(ritin, fun), '') as funritin FROM words" + qWhere)
		if err != nil {
			log.Fatal(err)
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
		updateEnd(p)
	}
	if ipa {
		p := updateBegin(" ipa")
		updateIPA(queryEnd)
		updateEnd(p)
	}
	if ritin {
		p := updateBegin(" ritin")
		updateRitin(queryEnd)
		updateEnd(p)
	}

	if onlyFlaagd {
		// clear all flags
		_, err := DB.Exec("UPDATE words SET flaagd = false WHERE flaagd;")
		if err != nil {
			log.Fatal(err)
		}
	}

}

type progress struct {
	start time.Time
	s     *spinner.Spinner
}

func updateBegin(processes string) (p progress) {
	start := time.Now()
	fmt.Print("Updating" + processes)
	anim := []string{"       ", ".      ", "..     ", "...    ", "....   ", ".....  ", "...... ", ".......", "...... ", ".....  ", "....   ", "...    ", "..     ", ".      "}
	s := spinner.New(anim, 100*time.Millisecond)
	s.Start()
	return progress{start, s}
}

func updateEnd(p progress) {
	t := time.Now()
	elapsed := t.Sub(p.start)
	p.s.Stop()
	fmt.Printf("... Done. (%v)\n", elapsed)
}

// Update IPA generates the rough IPA representation of words by
// consolidating diphthongs and digraphs into single characters,
// and by splitting y and w into semivowels and vowels
func updateIPA(queryEnd string) {
	queryWhere := queryEnd
	if queryEnd != ";" {
		queryEnd = " AND" + queryEnd
		queryWhere = " WHERE" + queryWhere
	}

	_, err := DB.Exec(`update words set ipa = funsil` + queryWhere + `
	
update words set ipa = regexp_replace(funsil, 'ng', 'ŋ', 'g')
where id in (select id from words where funsil similar TO '%ng%')` + queryEnd + `

update words set ipa = regexp_replace(funsil, 'dh', 'ð', 'g')
where id in (select id from words where funsil similar TO '%dh%')` + queryEnd + `

update words set ipa = regexp_replace(funsil, 'th', 'θ', 'g')
where id in (select id from words where funsil similar TO '%th%')` + queryEnd + `

update words set ipa = regexp_replace(funsil, 'dž', 'ʤ', 'g')
where id in (select id from words where funsil similar TO '%dž%')` + queryEnd + `

update words set ipa = regexp_replace(funsil, 'tš', 'ʧ', 'g')
where id in (select id from words where funsil similar TO '%tš%')` + queryEnd + `

update words set ipa = regexp_replace(funsil, 'ai', 'ā', 'g')
where id in (select id from words where funsil similar TO '%ai%')` + queryEnd + `

update words set ipa = regexp_replace(funsil, 'aw', 'å', 'g')
where id in (select id from words where funsil similar TO '%aw%')` + queryEnd + `

update words set ipa = regexp_replace(funsil, 'ei', 'ē', 'g')
where id in (select id from words where funsil similar TO '%ei%')` + queryEnd + `

update words set ipa = regexp_replace(funsil, 'oi', 'ō', 'g')
where id in (select id from words where funsil similar TO '%oi%')` + queryEnd + `

update words set ipa = regexp_replace(funsil, 'y([aäeoøiu])', 'j\1', 'g')
where id in (select id from words where funsil similar TO '%y[aäeoøiu]%')` + queryEnd + `

update words set ipa = regexp_replace(funsil, 'w([aäeoøiu])', 'ʍ\1', 'g')
where id in (select id from words where funsil similar TO '%w[aäeoøiu]%')` + queryEnd)

	if err != nil {
		log.Fatal(err)
	}
}

// updateRitin updates the written form of words that that differ from their pronounced form
// where it can be easily derived (in the case of syllable boundaries between d/t and h, insert dash)
func updateRitin(queryEnd string) {

	if queryEnd != ";" {
		queryEnd = " AND" + queryEnd
	}
	_, err := DB.Exec("update words set ritin = regexp_replace(fun, '([tdsz])h', '\\1-h', 'g') where funsil similar to '%[tdsz][ˈˌ·]h%'" + queryEnd)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec("update words set ritin = regexp_replace(fun, 'ng', 'n-g', 'g') where funsil similar to '%n[ˈˌ·]g%'" + queryEnd)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
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
		log.Fatal(err)
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
		log.Fatal(err)
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
	dist := wordtools.FindDistance(strings.ToLower(funritin), strings.ToLower(trud), true)

	// update with new dist value
	var updateString string
	updateDist := DB.QueryRow("UPDATE words SET dist = '" + strconv.Itoa(dist) + "' WHERE id = " + strconv.Itoa(id) + ";")
	err = updateDist.Scan(&updateString)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
}
