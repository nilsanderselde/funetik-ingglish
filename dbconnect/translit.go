// Nils Elde
// https://gitlab.com/nilsanderselde

package dbconnect

import (
	"bufio"
	"fmt"
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
					var isBlank bool
					for _, trud := range newWords {
						if trud != "" { // don't store empty strings as words
							words = append(words, trud) // add truditional spelling to list of words
						} else {
							isBlank = true
						}
					}
					// Add a newline character as a "word" to be reinserted before displaying results
					if !isBlank {
						words = append(words, "\n")
					}
				}

				var output string
				// transliterate words
				for _, trud := range words {
					// if it's a single symbol by itself, just add it to output
					if r := []rune(trud)[0]; len(trud) == 1 && (unicode.IsPunct(r) || unicode.IsSymbol(r) || r == '\n') {
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

	trudR := []rune(trud)
	var leading string
	var trailing string

	// as long as first character is punctuation, add it to leading symbol string
	// and remove it from word string (trud)
	for len(trudR) > 1 && (unicode.IsPunct(trudR[0]) || unicode.IsSymbol(trudR[0])) {
		leading += string(trudR[0])
		trudR = trudR[1:]
		// fmt.Printf("{%s},{%s}\n", leading, string(trudR))
	}
	// as long as last character is punctuation, add it to trailing symbol string
	// and remove it from word string (trud)
	for len(trudR) > 1 && (unicode.IsPunct(trudR[len(trudR)-1]) || unicode.IsSymbol(trudR[len(trudR)-1])) {
		trailing = string(trudR[len(trudR)-1]) + trailing
		trudR = trudR[0 : len(trudR)-1]
		// fmt.Printf("{%s},{%s}\n", string(trudR), trailing)
	}
	trud = replaceRunes(trudR) //

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
	row := DB.QueryRow("SELECT COALESCE(ritin, fun) FROM words WHERE trud = $1;", trud)
	err := row.Scan(&fun)
	if err != nil {
		// fmt.Println("not found, checking lowercase:", trud)
		row = DB.QueryRow("SELECT COALESCE(ritin, fun) FROM words WHERE trud = $1;", strings.ToLower(trud))
		err = row.Scan(&fun)
		if err != nil {
			// renamed variable as it may be either an unknown word returned in original form or a hyphenated/slashed word
			// return as transliterated sum of its parts
			word := trud

			// if word contains hyphen or slash in middle, dont add to list of unknown words (as it is compound)
			// if word contains hyphen in middle
			if strings.ContainsRune(word, '-') {
				// split word by hyphen and recursively process each part
				wordSplit := strings.Split(word, "-")
				newWord := getFun(wordSplit[0])
				for i := 1; i < len(wordSplit); i++ {
					newWord += "-" + getFun(wordSplit[i])
				}
				word = newWord
			}

			// if word contains slash in middle
			if strings.ContainsRune(word, '/') {
				// split word by slash and recursively process each part
				wordSplit := strings.Split(word, "/")
				newWord := getFun(wordSplit[0])
				for i := 1; i < len(wordSplit); i++ {
					newWord += "/" + getFun(wordSplit[i])
				}
				word = newWord
			}

			// cases for productive contractions
			if strings.HasSuffix(word, "'s") {
				newWordR := []rune(getFun(strings.TrimSuffix(word, "'s")))
				if strings.ContainsAny(string(newWordR[len(newWordR)-1]), "pfkt") || strings.HasSuffix(string(newWordR), "th") {
					word = string(newWordR) + "'s"
				} else {
					word = string(newWordR) + "'z"
				}
			} else if strings.HasSuffix(word, "'re") {
				newWordR := []rune(getFun(strings.TrimSuffix(word, "'re")))
				word = string(newWordR) + "'r"
			} else if strings.HasSuffix(word, "'d") {
				newWordR := []rune(getFun(strings.TrimSuffix(word, "'d")))
				word = string(newWordR) + "'d"
			} else if strings.HasSuffix(word, "'ve") {
				newWordR := []rune(getFun(strings.TrimSuffix(word, "'ve")))
				word = string(newWordR) + "'v"
			} else if strings.HasSuffix(word, "'d've") {
				newWordR := []rune(getFun(strings.TrimSuffix(word, "'d've")))
				word = string(newWordR) + "'d'v"
			} else if strings.HasSuffix(word, "'ll") {
				newWordR := []rune(getFun(strings.TrimSuffix(word, "'ll")))
				word = string(newWordR) + "'l"
			}

			addToList := true
			wordR := []rune(word)
			for _, r := range wordR {
				if unicode.IsPunct(r) ||
					unicode.IsSymbol(r) ||
					unicode.IsDigit(r) {
					addToList = false
					break
				}
			}

			// add to unknown word table
			if addToList {
				_, err = DB.Exec(`INSERT INTO unknown (trud)
SELECT '` + word + `'
WHERE NOT EXISTS (SELECT trud FROM unknown WHERE LOWER(trud) = '` + strings.ToLower(word) + `')
AND NOT EXISTS (SELECT trud FROM words WHERE LOWER(trud) = '` + strings.ToLower(word) + `')`)
				if err != nil {
					fmt.Println("word not saved")
				}
			}

			return leading + word + trailing
		}
	}

	// format funetik spelling with saved case pattern
	if isUpperCase && len(trud) > 1 {
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
