// Nils Elde
// https://gitlab.com/nilsanderselde

package dbconnect

import (
	"database/sql"
	"fmt"
)

// UpdateAutoValues automatically fills values for:
// 1. fun      (remove syllable markings from funsil)
// 2. numsil   (count syllable markings from funsil)
// 3. funsort  (substitution cipher on fun)
// 4. dist     (calc lev dist between fun and trud)
func RestoreBackup() {
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
	// update fun and/or numsil with values generated using funsil
	_, err = db.Exec(``)
	if err != nil {
		// log.Fatal(err)
		fmt.Println(err)
	}

}
