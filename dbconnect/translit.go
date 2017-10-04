// Nils Elde
// https://gitlab.com/nilsanderselde

package dbconnect

import (
	"bufio"
	"database/sql"
	"log"
	"net/http"
	"strings"
)

var (
	// UnknownWords contains words that were unsuccessfully transliterated
	UnknownWords []string
)

// ProcessTrud tries to process form input by transliterating
// traditional to funetik spellings.
func ProcessTrud(r *http.Request) (outputLines []string, input string) {

	// attempt to parse form
	err := r.ParseForm()
	// if form was parsed successfully
	if err == nil {
		// if the input text field was submitted
		if r.Form["inputtext"] != nil {
			// if the input text field is not blank
			if r.Form["inputtext"][0] != "" {
				// save input text as "text"
				text := r.Form["inputtext"][0]

				// make a copy to place back in translit text area
				input = text

				// for collecting words after transliteration attempt
				var words []string

				// load text into scanner to split it in lines
				scanner := bufio.NewScanner(strings.NewReader(text))
				for scanner.Scan() { // for each line, split it into words by spaces
					newWords := strings.Split(scanner.Text(), " ")
					for _, trud := range newWords {
						if trud != "" { // don't store empty strings as words
							words = append(words, trud) // add truditional spelling to list of words
						}
					}
					// Add a newline character as a "word" to be reinserted before displaying results
					words = append(words, "\n")
				}

				var last rune
				var output string
				// transliterate words
				for i, trud := range words {
					// if it's a single symbol by itself, just add it to output
					if len(trud) == 1 {
						if hasPunc([]rune(trud)[0]) || trud == "\n" {
							output += trud
							continue
						}
					} else { // don't process strings containing saved newlines
						if i == 0 { // for first word in input, store the funetik spelling in title case
							output += capitalizeContraction(getFun(trud)) + " "
						} else { // for words following sentence-terminating punctuation, store the funetik spellings in title case
							if last == '.' || last == '!' || last == '?' {
								output += capitalizeContraction(getFun(trud)) + " "
							} else { // otherwise just store the funetik spelling
								output += getFun(trud) + " "
							}
						}
					}

					stringR := []rune(trud)        // convert string to rune array to access last character
					last = stringR[len(stringR)-1] // save last character for next iteration
				}
				// Split one-line results string into separate lines
				scanner = bufio.NewScanner(strings.NewReader(output))
				for scanner.Scan() {
					outputLines = append(outputLines, scanner.Text())
				}
				// fmt.Println(UnknownWords)
				return outputLines, input
			}
		}
	}
	return outputLines, input
}

// prevents the first letter after the apostrophe in a contraction
// from being capitalized by strings.Title (as in "Don'T")
func capitalizeContraction(word string) (capitalized string) {
	wordR := []rune(word)
	capitalized += strings.Title(string(wordR[0]))
	capitalized += string(wordR[1:])
	return capitalized
}

func hasPunc(r rune) bool {
	if r == '`' || r == '~' || r == '!' || r == '@' || r == '#' || r == '$' ||
		r == '%' || r == '^' || r == '&' || r == '*' || r == '(' || r == ')' ||
		r == '_' || r == '+' || r == '-' || r == '=' || r == '[' || r == ']' ||
		r == '\\' || r == '{' || r == '}' || r == '|' || r == ';' || r == '\'' ||
		r == ':' || r == '"' || r == ',' || r == '.' || r == '/' || r == '<' ||
		r == '>' || r == '?' {
		return true
	}
	return false
}

// getFun tries to return the corresponding funetik spelling of an English word
func getFun(trud string) (fun string) {
	db, err := sql.Open("postgres", DBInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	stringR := []rune(trud)
	var leading string
	var trailing string

	// as long as first character is punctuation, add it to leading symbol string
	// and remove it from word string (trud)
	for hasPunc(stringR[0]) {
		leading += string(stringR[0])
		stringR = stringR[1:]
		// fmt.Printf("{%s},{%s}\n", leading, string(stringR))

	}
	// as long as last character is punctuation, add it to trailing symbol string
	// and remove it from word string (trud)
	for hasPunc(stringR[len(stringR)-1]) {
		trailing += string(stringR[len(stringR)-1])
		stringR = stringR[0 : len(stringR)-1]
		// fmt.Printf("{%s},{%s}\n", string(stringR), trailing)
	}

	trud = replaceRunes(stringR) //

	// attempt to split trud by hyphen and recursively process each part
	trudSplit := strings.Split(trud, "-")
	if len(trudSplit) > 1 {
		newTrud := getFun(trudSplit[0])
		for i := 1; i < len(trudSplit); i++ {
			newTrud += "-" + getFun(trudSplit[i])
		}
		trud = newTrud
	}

	// update fun and/or numsil with values generated using funsil
	// fmt.Println(">>", trud)
	rows := db.QueryRow("SELECT COALESCE(ritin, fun) FROM words WHERE trud = $1;", trud)
	err = rows.Scan(&fun)
	if err != nil {
		// fmt.Println("not found, checking lowercase:", trud)
		rows = db.QueryRow("SELECT COALESCE(ritin, fun) FROM words WHERE trud = $1;", strings.ToLower(trud))
		err = rows.Scan(&fun)
		if err != nil {
			// fmt.Println("not found, returning trud:", trud)
			UnknownWords = append(UnknownWords, trud)
			return leading + trud + trailing
		}
	}
	// fmt.Println(fun)
	return leading + fun + trailing
}

// replace ‘’“” with basic versions
func replaceRunes(trud []rune) (newTrud string) {
	for _, r := range trud {
		if r == '‘' || r == '’' {
			newTrud += "'"
		} else if r == '“' || r == '”' {
			newTrud += `"`
		} else {
			newTrud += string(r)
		}
	}
	return newTrud
}
