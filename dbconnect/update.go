package dbconnect

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"gitlab.com/nilsanderselde/funetik-ingglish/wordtools"
)

// UpdateAllAutoValues automatically fills values for:
// 1. fun      (remove syllable markings from funsil)
// 2. numsil   (count syllable markings from funsil)
// 3. funsort  (substitution cipher on fun)
// 4. dist     (calc lev dist between fun and trud)
func UpdateAllAutoValues() {
	db, err := sql.Open("postgres", DBInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// update fun and numsil with values generated using funsil
	rows, err := db.Query("SELECT id, funsil FROM words;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		UpdateFun(rows, db)
		UpdateNumsil(rows, db)
	}

	// update funsort and dist with values generated using fun and trud
	rows, err = db.Query("SELECT id, fun, trud FROM words;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		UpdateFunsort(rows, db)
		UpdateDist(rows, db)
	}
	fmt.Println("All automatically generated columns updated.")
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
	err := row.Scan(&id, &fun)

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
	err := row.Scan(&id, &fun, &trud)

	// Generate value for "dist" (true means flipping two adjacent letters is considered one move)
	dist := wordtools.FindDistance(fun, trud, true)

	// update with new dist value
	var updateString string
	updateDist := db.QueryRow("UPDATE words SET dist = '" + strconv.Itoa(dist) + "' WHERE id = " + strconv.Itoa(id) + ";")
	err = updateDist.Scan(&updateString)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
}
