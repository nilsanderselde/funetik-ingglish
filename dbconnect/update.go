// Nils Elde
// https://gitlab.com/nilsanderselde

package dbconnect

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/briandowns/spinner"
	"gitlab.com/nilsanderselde/funetik-ingglish/wordtools"
)

// UpdateAutoValues automatically fills values for:
// 1. fun      (remove syllable markings from funsil)
// 2. numsil   (count syllable markings from funsil)
// 3. funsort  (substitution cipher on fun)
// 4. dist     (calc lev dist between fun and trud)
// onlyFlaagd means only flaagd will be processed
// rowID used to specify row to update, use -1 for all
func UpdateAutoValues(fun bool, numsil bool, funsort bool, dist bool, onlyFlaagd bool, rowID int) {
	if !(fun || numsil || funsort || dist) {
		return // don't bother connecting if no updates will occur
	}
	db, err := sql.Open("postgres", DBInfo)
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err)
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
		rows, err := db.Query("SELECT id, funsil" + queryFrom)
		if err != nil {
			// log.Fatal(err)
			fmt.Println(err)
		}

		for rows.Next() {
			if fun {
				UpdateFun(rows, db)
			}
			if numsil {
				UpdateNumsil(rows, db)
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
		rows, err := db.Query("SELECT id, fun, trud, COALESCE(COALESCE(ritin, fun), '') as funritin" + queryFrom)
		if err != nil {
			// log.Fatal(err)
			fmt.Println(err)
		}

		for rows.Next() {
			if funsort {
				UpdateFunsort(rows, db)
			}
			if dist {
				UpdateDist(rows, db)
			}
		}
		rows.Close()
		t := time.Now()
		elapsed := t.Sub(start)
		s.Stop()
		fmt.Printf("Done. (%v)\n", elapsed)

		// clear all flags
		_, err = db.Exec("UPDATE words SET flaagd = false WHERE flaagd;")
		if err != nil {
			// log.Fatal(err)
			fmt.Println(err)
		}
	}
}

// UpdateFun updates passed row with generated
// value for column "fun" (funetik spelling)
func UpdateFun(row *sql.Rows, db *sql.DB) {
	var id int
	var funsil string
	err := row.Scan(&id, &funsil)
	// Generate value for "fun"
	fun := wordtools.RemoveSyllableMarkers(funsil)

	// update with new fun value
	var updateString string
	updateFun := db.QueryRow("UPDATE words SET fun = '" + fun + "' WHERE id = " + strconv.Itoa(id) + ";")
	err = updateFun.Scan(&updateString)
	if err != nil && err != sql.ErrNoRows {
		// log.Fatal(err)
		fmt.Println(err)
	}
}

// UpdateNumsil updates passed row with generated
// value for column  "numsil" (number of syllables)
func UpdateNumsil(row *sql.Rows, db *sql.DB) {
	var id int
	var funsil string
	err := row.Scan(&id, &funsil)

	// Generate value for "numsil"
	numsil := wordtools.CountSyllables(funsil)

	// update with new numsil value
	var updateString string
	updateNumsil := db.QueryRow("UPDATE words SET numsil = '" + strconv.Itoa(numsil) + "' WHERE id = " + strconv.Itoa(id) + ";")
	err = updateNumsil.Scan(&updateString)
	if err != nil && err != sql.ErrNoRows {
		// log.Fatal(err)
		fmt.Println(err)
	}
}

// UpdateFunsort updates passed row with generated
// value for column "funsort" (funetik sort order)
func UpdateFunsort(row *sql.Rows, db *sql.DB) {
	var id int
	var fun string
	var trud string
	var funritin string
	err := row.Scan(&id, &fun, &trud, &funritin)

	// Generate value for "funsort"
	funsort := wordtools.SubstitutionCypher(fun)

	// update with new funsort value
	var updateString string
	updateFunsort := db.QueryRow("UPDATE words SET funsort = '" + funsort + "' WHERE id = " + strconv.Itoa(id) + ";")
	err = updateFunsort.Scan(&updateString)
	if err != nil && err != sql.ErrNoRows {
		// log.Fatal(err)
		fmt.Println(err)
	}
}

// UpdateDist updates passed row with calculated
// value for column "levdist" (Levenshtein distance)
func UpdateDist(row *sql.Rows, db *sql.DB) {
	var id int
	var fun string
	var trud string
	var funritin string
	err := row.Scan(&id, &fun, &trud, &funritin)

	// Generate value for "dist" (true means flipping two adjacent letters is considered one move)
	dist := wordtools.FindDistance(funritin, trud, true)

	// update with new dist value
	var updateString string
	updateDist := db.QueryRow("UPDATE words SET dist = '" + strconv.Itoa(dist) + "' WHERE id = " + strconv.Itoa(id) + ";")
	err = updateDist.Scan(&updateString)
	if err != nil && err != sql.ErrNoRows {
		// log.Fatal(err)
		fmt.Println(err)
	}
}

// UpdateAllAutoValues automatically generates values for all rows
// (About 15 minutes with all tasks enabled and 50,000 words)
func UpdateAllAutoValues(fun bool, numsil bool, funsort bool, dist bool) {
	UpdateAutoValues(fun, numsil, funsort, dist, false, -1)
}

// UpdateFlaggedAutoValues automatically generates values for all flaagd rows
// (About 15 minutes with all tasks enabled and 50,000 words)
func UpdateFlaggedAutoValues(fun bool, numsil bool, funsort bool, dist bool) {
	UpdateAutoValues(fun, numsil, funsort, dist, true, -1)
}
