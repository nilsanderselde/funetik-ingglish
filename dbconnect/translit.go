// Nils Elde
// https://gitlab.com/nilsanderselde

package dbconnect

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
)

// ProcessTrud tries to process form input by transliterating
// traditional to funetik spellings.
func ProcessTrud(r *http.Request) string {
	err := r.ParseForm()
	if err == nil {
		if r.Form["inputtext"] != nil {
			if r.Form["inputtext"][0] != "" {
				text := r.Form["inputtext"][0]
				words := strings.Split(text, " ")
				last := ' '
				result := ""
				for i, trud := range words {
					if i == 0 {
						result += strings.Title(GetFun(trud)) + " "
					} else {
						if last == '.' {
							result += strings.Title(GetFun(trud)) + " "
						} else {
							result += GetFun(trud) + " "
						}
					}
					stringR := []rune(trud)
					last = stringR[len(stringR)-1]
				}
				return result
			}
			return ""
		}
		return ""
	}
	return ""
}

func hasPunc(last rune) bool {
	if last == '.' ||
		last == ',' ||
		last == '!' ||
		last == '?' ||
		last == ':' ||
		last == ';' ||
		last == '\'' ||
		last == '"' ||
		last == '-' ||
		last == '_' ||
		last == '+' ||
		last == '=' ||
		last == '/' {
		return true
	}
	return false
}

// GetFun tries to return the corresponding funetik spelling of an English word
func GetFun(trud string) (fun string) {
	stringR := []rune(trud)
	first := stringR[0]
	last := stringR[len(stringR)-1]
	firstIsPunc := hasPunc(first)
	lastIsPunc := hasPunc(last)

	db, err := sql.Open("postgres", DBInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	if firstIsPunc {
		stringR = []rune(trud)
		stringR = stringR[1:]
		trud = string(stringR)
	}
	if lastIsPunc {
		stringR = []rune(trud)
		stringR = stringR[0 : len(stringR)-1]
		trud = string(stringR)
	}

	// update fun and/or numsil with values generated using funsil
	rows := db.QueryRow("SELECT COALESCE(ritin, fun) FROM words WHERE trud = $1;", trud)
	err = rows.Scan(&fun)
	if err != nil {
		rows = db.QueryRow("SELECT COALESCE(ritin, fun) FROM words WHERE trud = $1;", strings.ToLower(trud))
		err = rows.Scan(&fun)
		if err != nil {
			return trud
		}
	}
	if firstIsPunc {
		fun = string(first) + fun
	}
	if lastIsPunc {
		fun += string(last)
	}
	return fun
}
