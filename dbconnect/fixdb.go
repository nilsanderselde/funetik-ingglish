// Nils Elde
// https://gitlab.com/nilsanderselde

package dbconnect

import (
	"database/sql"
	"fmt"
)

// RestoreBackup executes a massive amount of SQL insert statements
// to rebuild database
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
