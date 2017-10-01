package dbconnect

import (
	"database/sql"
	"log"
	"math/rand"
	"strconv"
	"time"

	"gitlab.com/nilsanderselde/funetik-ingglish/global"
)

// NewRandomRune creates and stores a new rune for reloading the word list after toggling fields
func NewRandomRune() {
	rand.Seed(time.Now().Unix())

	for global.CurrRand == global.LastRand {
		global.CurrRand = rune(rand.Intn(2) + 97)
	}
	global.LastRand = global.CurrRand
}

// FlagRow inverts a row's flaagd field value
func FlagRow(id string) {
	NewRandomRune()

	db, err := sql.Open("postgres", DBInfo)
	if err != nil {
		// log.Fatal(err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		// log.Fatal(err)
		return
	}
	// set flaagd of row with passed id to opposite of its current state
	flägResult := db.QueryRow("SELECT flaagd FROM words WHERE id = " + id + ";")
	var flägString string
	err = flägResult.Scan(&flägString)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
	var fläg bool
	fläg, err = strconv.ParseBool(flägString)

	// if true, make false, and vice versa
	updateFläg := db.QueryRow("	UPDATE words SET flaagd = " + strconv.FormatBool(!fläg) + " WHERE id = " + id + ";	")
	var updateString string
	err = updateFläg.Scan(&updateString)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
}
