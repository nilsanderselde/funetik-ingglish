// Nils Elde
// https://gitlab.com/nilsanderselde

package dbconnect

import (
	"bufio"
	"database/sql"
	"log"
	"net/http"
	"strings"
	"unicode"
)

var (
	// UnknownWords contains words that were unsuccessfully transliterated
	UnknownWords []string
)

// Output encapsulates the two return values of ProcessTrud:
// a slice of strings, each a line of the output text;
// and the original input string.
type Output struct {
	OutputLines []string
	PrevInput   string
}

// ProcessTrud tries to process form input by transliterating
// traditional to funetik spellings.
func ProcessTrud(ch chan Output, r *http.Request) {
	outStruct := Output{}
	// attempt to parse form
	err := r.ParseForm()
	// if form was parsed successfully
	if err == nil {
		// if the input text field was submitted
		if r.Form["inputtext"] != nil {
			// if the input text field is not blank
			if r.Form["inputtext"][0] != "" {
				// save first 2500 runes of input text
				input := r.Form["inputtext"][0]
				inputR := []rune(input)
				if len(inputR) > 2500 {
					input = string(inputR[0:2500])
				}

				// make a copy to place back in translit text area
				outStruct.PrevInput = input

				// for collecting words after transliteration attempt
				var words []string

				// load text into scanner to split it in lines
				scanner := bufio.NewScanner(strings.NewReader(input))
				for scanner.Scan() { // for each line, split it into words by spaces
					newWords := strings.Fields(scanner.Text())
					for _, trud := range newWords {
						if trud != "" { // don't store empty strings as words
							words = append(words, trud) // add truditional spelling to list of words
						}
					}
					// Add a newline character as a "word" to be reinserted before displaying results
					words = append(words, "\n")
				}

				var output string
				// transliterate words
				for _, trud := range words {
					// if it's a single symbol by itself, just add it to output
					if len(trud) == 1 && (unicode.IsPunct([]rune(trud)[0]) || trud == "\n") {
						output += trud
						continue

					}
					output += getFun(trud) + " "
				}
				// Split one-line results string into separate lines
				scanner = bufio.NewScanner(strings.NewReader(output))
				for scanner.Scan() {
					outStruct.OutputLines = append(outStruct.OutputLines, scanner.Text())
				}
				// fmt.Println(UnknownWords)
			}
		}
	}
	ch <- outStruct
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
	for unicode.IsPunct(stringR[0]) {
		leading += string(stringR[0])
		stringR = stringR[1:]
		// fmt.Printf("{%s},{%s}\n", leading, string(stringR))

	}
	// as long as last character is punctuation, add it to trailing symbol string
	// and remove it from word string (trud)
	for unicode.IsPunct(stringR[len(stringR)-1]) {
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

	// determine if case is upper or title, optimisticly
	isUpperCase := true
	isTitleCase := true
	for i, r := range []rune(trud) {
		if !unicode.IsUpper(r) {
			if i == 0 {
				isTitleCase = false
			}
			isUpperCase = false
		}
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

	// format funetik spelling with saved case pattern
	if isUpperCase {
		fun = strings.ToUpper(fun)
	} else if isTitleCase {
		fun = capitalizeContraction(fun)
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

// prevents the first letter after the apostrophe in a contraction
// from being capitalized by strings.Title (as in "Don'T")
func capitalizeContraction(word string) (capitalized string) {
	wordR := []rune(word)
	capitalized += strings.Title(string(wordR[0]))
	capitalized += string(wordR[1:])
	return capitalized
}
