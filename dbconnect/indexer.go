package dbconnect

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"gitlab.com/nilsanderselde/funetik-ingglish/global"
)

// IndexByInitial gets the row number of the first word of each group
// that start with the same letter
func IndexByInitial() {
	queryBegin := `SELECT rnum FROM (
	SELECT rnum FROM (
		SELECT fun, row_number() OVER (order by funsort, id) as rnum
		FROM words order by FUNSORT, ID
	) as inoor where lower(fun) like '`
	queryEnd := `%'
) as aawtoor limit 1;`

	global.InitialIndex = make([]global.InitialIndexValue, 26, 26)

	for i := 0; i < 26; i++ {
		var start int
		err := DB.QueryRow(queryBegin + string(global.Alphabet[i]) + queryEnd).Scan(&start)
		switch {
		case err == sql.ErrNoRows:
			fmt.Println("No words start with that letter.")
		case err != nil:
			log.Fatal("IndexByInitial: ", err)
		default:
			global.InitialIndex[i].Letter = global.Alphabet[i]
			global.InitialIndex[i].Index = strconv.Itoa(start - 1)
		}

	}
}
