// Nils Elde
// https://gitlab.com/nilsanderselde
// This package sorts a list of words according to a
// custom alphabetical order supplied by the user

package customsort

import (
	"bufio"
	"log"
	"os"
	"sort"
	"strings"

	"gitlab.com/nilsanderselde/funetik-ingglish/params"
)

var (
	letters  []rune
	alphabet []rune
)

// CustomAlphabeticalOrder is the alias for array of strings to be sorted
type CustomAlphabeticalOrder []string

func (s CustomAlphabeticalOrder) Len() int {
	return len(s)
}

func (s CustomAlphabeticalOrder) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func findIndex(r rune) int {
	for i := 0; i < len(alphabet); i++ {
		if r == alphabet[i] {
			return i
		}
	}
	return -1
}

func (s CustomAlphabeticalOrder) Less(i, j int) bool {

	word1 := []rune(strings.Split(s[i], "\t")[0])
	word2 := []rune(strings.Split(s[j], "\t")[0])
	length1 := len(word1)
	length2 := len(word2)

	var minlength int
	if length1 < length2 {
		minlength = length1
	} else if length1 > length2 {
		minlength = length2
	} else {
		minlength = length1
	}
	for k := 0; k < minlength; k++ {
		letter1 := []rune(strings.ToLower(string(word1)))[k]
		letter1order := findIndex(letter1)
		letter2 := []rune(strings.ToLower(string(word2)))[k]
		letter2order := findIndex(letter2)

		// if on last letter and word is the same so far, return true if first
		// word is shorter
		if k == minlength-1 {
			if letter1order == letter2order {
				return length1 < length2
			}
		}

		if letter1order > letter2order {
			return false
		} else if letter1order < letter2order {
			return true
		}
	}
	return false
}

// IgnoreCaseTrudOrder is the type for an array of strings to be case insensitively sorted
type IgnoreCaseTrudOrder []string

func (s IgnoreCaseTrudOrder) Len() int {
	return len(s)
}

func (s IgnoreCaseTrudOrder) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s IgnoreCaseTrudOrder) Less(i, j int) bool {

	word1 := []rune(strings.Split(s[i], "\t")[0])
	word2 := []rune(strings.Split(s[j], "\t")[0])
	length1 := len(word1)
	length2 := len(word2)

	var minlength int
	if length1 < length2 {
		minlength = length1
	} else if length1 > length2 {
		minlength = length2
	} else {
		minlength = length1
	}

	for k := 0; k < minlength; k++ {
		letter1 := []rune(strings.ToLower(string(word1)))[k]
		letter2 := []rune(strings.ToLower(string(word2)))[k]

		// if on last letter and word is the same so far, return true if first
		// word is shorter
		if k == minlength-1 {
			if letter1 == letter2 {
				return length1 < length2
			}
		}
		if letter1 > letter2 {
			return false
		} else if letter1 < letter2 {
			return true
		}
	}
	return false
}

// SortWords sorts a list of words
//
// It splits a tab-delimited text file into lines, and sorts by the
// first word in each line (useful for dictionary/glossary/encyclopedia sorting)
//
func SortWords(args params.Params) [][]string {

	alphabet = args.Order

	// Open file containing rows of words
	file, err := os.Open("C:/Users/Nils/Go/io/words_to_sort.txt")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if args.Reverse {
		sort.Sort(sort.Reverse(CustomAlphabeticalOrder(lines)))
	} else {
		sort.Sort(CustomAlphabeticalOrder(lines))
	}

	// Return results as 2D array of strings, sorted by the first string in each subarray
	var splitLines [][]string
	for i := 0; i < len(lines); i++ {
		splitLines = append(splitLines, strings.Split(lines[i], "\t"))
	}
	return splitLines
}

// SortByTrud sorts a list of words by the funetik spelling in pseudo-traditional order
func SortByTrud(args params.Params) [][]string {

	// Open file containing rows of words
	file, err := os.Open("C:/Users/Nils/Go/io/words_to_sort.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// adds extra traditional words added for quick sorting
		lines = append(lines, strings.Split(scanner.Text(), "\t")[1]+"\t"+scanner.Text())
	}
	if args.Reverse {
		sort.Sort(sort.Reverse(IgnoreCaseTrudOrder(lines)))
	} else {
		sort.Sort(IgnoreCaseTrudOrder(lines))
	}

	// Return results as 2D array of strings, sorted by Levenshtein distance
	var splitLines [][]string
	for i := 0; i < len(lines); i++ {
		// removes extra distance added for quick sorting
		trimmedLine := strings.Split(lines[i], "\t")
		trimmedLine = []string{trimmedLine[1], trimmedLine[2], trimmedLine[3], trimmedLine[4]}
		splitLines = append(splitLines, trimmedLine)
	}
	return splitLines
}

// SortByDistance sorts a list of words by Levenshtein distance
func SortByDistance(args params.Params) [][]string {

	// Open file containing rows of words
	file, err := os.Open("C:/Users/Nils/Go/io/words_to_sort.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// adds extra distance added for quick sorting
		lines = append(lines, strings.Split(scanner.Text(), "\t")[2]+"\t"+scanner.Text())
	}
	if args.Reverse {
		sort.Sort(sort.Reverse(sort.StringSlice(lines)))
	} else {
		sort.Strings(lines)
	}

	// Return results as 2D array of strings, sorted by levenshtein distance
	var splitLines [][]string
	for i := 0; i < len(lines); i++ {
		// removes extra distance added for quick sorting
		trimmedLine := strings.Split(lines[i], "\t")
		trimmedLine = []string{trimmedLine[1], trimmedLine[2], trimmedLine[3], trimmedLine[4]}
		splitLines = append(splitLines, trimmedLine)
	}
	return splitLines
}
