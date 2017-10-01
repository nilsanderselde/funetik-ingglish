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

// UpdateAllAutoValues automatically fills values for:
// 1. fun      (remove syllable markings from funsil)
// 2. numsil   (count syllable markings from funsil)
// 3. funsort  (substitution cipher on fun)
// 4. dist     (calc lev dist between fun and trud)
func UpdateAllAutoValues(fun bool, numsil bool, funsort bool, dist bool) {
	if !(fun || numsil || funsort || dist) {
		return // don't bother connecting if no updates will occur
	}
	db, err := sql.Open("postgres", DBInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

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
		rows, err := db.Query("SELECT id, funsil FROM words;")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			if fun {
				UpdateFun(rows, db)
			}
			if numsil {
				UpdateNumsil(rows, db)
			}
		}
		t := time.Now()
		elapsed := t.Sub(start)
		s.Stop()
		fmt.Println("Done. (", elapsed, ")")
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

		// update funsort and dist with values generated using fun and trud (for dist, use written form if different)
		rows, err := db.Query("SELECT id, fun, trud, COALESCE(COALESCE(ritin, fun), '') as funritin FROM words where id > 49152;")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			if funsort {
				UpdateFunsort(rows, db)
			}
			if dist {
				UpdateDist(rows, db)
			}
		}
		t := time.Now()
		elapsed := t.Sub(start)
		fmt.Println("Done. (", elapsed, ")")
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
		log.Fatal(err)
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
		log.Fatal(err)
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
		log.Fatal(err)
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
		log.Fatal(err)
	}
}
