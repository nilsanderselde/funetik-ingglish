package dbconnect

import (
	"database/sql"
	"log"
	"math/rand"
	"strconv"
	"strings"
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
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		// log.Fatal(err)
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
	var updateString string

	updateFläg := db.QueryRow(`
	BEGIN;
	UPDATE words SET flaagd = ` + strconv.FormatBool(!fläg) + ` WHERE id = ` + id + `;
	SELECT flaagd FROM words WHERE id = ` + id + `;
	COMMIT;
	`)
	err = updateFläg.Scan(&updateString)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
	// fmt.Println("changing from " + flägString + " to " + updateString)
}

// UpdateFunSort updates all rows with new data, performing
// automatic calculation of funsort
func UpdateFunSort() {

	// automatically fill values for:
	// Fun     string	(remove syllable markings from funsil)
	// Funsort string  (substitution cipher on fun)
	// Numsil  int		(count syllable markings from funsil)
	// Dist    int		(calc lev dist between fun and trud)

	db, err := sql.Open("postgres", DBInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// get query
	rows, err := db.Query("SELECT id, fun FROM words;")
	if err != nil {
		log.Fatal(err)

	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var fun string
		err = rows.Scan(&id, &fun)

		funsort := SubstitutionCypher(fun)

		// update with new funsort value
		var updateString string
		updateSort := db.QueryRow(`
BEGIN;
UPDATE words SET funsort = '` + funsort + `' WHERE id = ` + strconv.Itoa(id) + `;
SELECT funsort FROM words WHERE id = ` + strconv.Itoa(id) + `;
COMMIT;
		`)
		err = updateSort.Scan(&updateString)
		// fmt.Println(updateString)
		if err != nil && err != sql.ErrNoRows {
			log.Fatal(err)
		}
	}

}

// SubstitutionCypher substitutes letters from the first row
// below to the letter directly below it to faciliate sorting
// based on a custom alphabet in SQL.
// 	aäeiywuøolrmnbpvfgkdtzsžšh
// 	ABCDEFGHIJKLMNOPQRSTUVWXYZ
func SubstitutionCypher(fun string) (funSort string) {
	funRunes := []rune(strings.ToLower(fun))
	for _, r := range funRunes {
		switch r {
		case 'a':
			funSort += string('A')
		case 'ä':
			funSort += string('B')
		case 'e':
			funSort += string('C')
		case 'i':
			funSort += string('D')
		case 'y':
			funSort += string('E')
		case 'w':
			funSort += string('F')
		case 'u':
			funSort += string('G')
		case 'ø':
			funSort += string('H')
		case 'o':
			funSort += string('I')
		case 'l':
			funSort += string('J')
		case 'r':
			funSort += string('K')
		case 'm':
			funSort += string('L')
		case 'n':
			funSort += string('M')
		case 'b':
			funSort += string('N')
		case 'p':
			funSort += string('O')
		case 'v':
			funSort += string('P')
		case 'f':
			funSort += string('Q')
		case 'g':
			funSort += string('R')
		case 'k':
			funSort += string('S')
		case 'd':
			funSort += string('T')
		case 't':
			funSort += string('U')
		case 'z':
			funSort += string('V')
		case 's':
			funSort += string('W')
		case 'ž':
			funSort += string('X')
		case 'š':
			funSort += string('Y')
		case 'h':
			funSort += string('Z')
		}
	}
	return funSort
}

// numsil

// dist

// flaagd
