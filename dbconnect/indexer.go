// Nils Elde
// https://github.com/nilsanderselde

package dbconnect

import (
	"fmt"
	"strconv"

	"gitlab.com/nilsanderselde/funetik-ingglish/global"
)

// Indexer gets the row number of the first word of each group
// that start each phonetic letter, traditional letter, and also
// indexes by distance
func Indexer() {
	queryBegin := `SELECT rnum FROM (
	SELECT rnum FROM (
		SELECT fun, row_number() OVER (order by funsort, trud, id) as rnum
		FROM words order by funsort, trud, id
	) as inoor where lower(fun) like '`
	queryEnd := `%'
) as aawtoor limit 1;`

	global.FunetikIndex = make([]global.OrderedIndexMap, 26, 26)

	for i := 0; i < 26; i++ {
		var start int
		err := DB.QueryRow(queryBegin + string(global.Älfubit[i]) + queryEnd).Scan(&start)
		if err == nil {
			global.FunetikIndex[i].Value = global.Älfubit[i]
			global.FunetikIndex[i].Offset = strconv.Itoa(start - 1)
		}

	}
	queryBegin = `SELECT rnum FROM (
	SELECT rnum FROM (
		SELECT trud, row_number() OVER (order by trud, funsort, id) as rnum
		FROM words order by trud, funsort, id
	) as inoor where lower(trud) like '`
	queryEnd = `%'
) as aawtoor limit 1;`

	global.TrudIndex = make([]global.OrderedIndexMap, 26, 26)

	for i := 0; i < 26; i++ {
		var start int
		err := DB.QueryRow(queryBegin + string(global.Alphabet[i]) + queryEnd).Scan(&start)
		if err == nil {
			global.TrudIndex[i].Value = global.Alphabet[i]
			global.TrudIndex[i].Offset = strconv.Itoa(start - 1)
		}
	}

	var maxDist int

	err := DB.QueryRow("SELECT max(dist) FROM words").Scan(&maxDist)
	if err != nil {
		fmt.Println(err)
	}

	queryBegin = `SELECT rnum FROM (
	SELECT rnum FROM (
		SELECT dist, row_number() OVER (order by dist, funsort, trud, id) as rnum
		FROM words order by dist, funsort, trud, id
	) as inoor where dist =`
	queryEnd = `
) as aawtoor limit 1;`

	global.DistIndex = make([]global.OrderedIndexMap, maxDist+1, maxDist+1)

	for i := 0; i < maxDist+1; i++ {
		var start int
		err := DB.QueryRow(queryBegin + strconv.Itoa(i) + queryEnd).Scan(&start)
		if err == nil {
			global.DistIndex[i].Value = strconv.Itoa(i)
			global.DistIndex[i].Offset = strconv.Itoa(start - 1)
		}
	}
}
