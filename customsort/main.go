// Nils Elde
// https://gitlab.com/nilsanderselde
// This package sorts a list of words according to a
// custom alphabetical order supplied by the user

package customsort

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

var letters []rune
var alphabet map[rune]int

// CustomAlphabeticalOrder is the alias for array of strings to be sorted
type CustomAlphabeticalOrder []string

func (s CustomAlphabeticalOrder) Len() int {
	return len(s)
}

func (s CustomAlphabeticalOrder) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
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
	// fmt.Print(string(word1) + ": " + strconv.Itoa(length1))
	// fmt.Print("; " + string(word2) + ": " + strconv.Itoa(length2))
	// fmt.Println("; minlength: " + strconv.Itoa(minlength))

	for k := 0; k < minlength; k++ {
		letter1 := []rune(strings.ToLower(string(word1)))[k]
		letter1order := alphabet[letter1]
		letter2 := []rune(strings.ToLower(string(word2)))[k]
		letter2order := alphabet[letter2]

		// fmt.Print(string(letter1) + ": " + strconv.Itoa(letter1order))
		// fmt.Println("\t" + string(letter2) + ": " + strconv.Itoa(letter2order))

		// if on last letter and word is the same so far, return true if first
		// word is shorter
		if k == minlength-1 {
			if letter1order == letter2order {
				return length1 < length2
			}
		}

		if letter1order > letter2order {
			// fmt.Println(">> " + string(word1) + " is after " + string(word2))
			// fmt.Println()
			return false
		} else if letter1order < letter2order {
			// fmt.Println(">> " + string(word1) + " is before " + string(word2))
			// fmt.Println()
			return true
		}

	}
	// fmt.Println(">> " + string(word1) + " is equal to " + string(word2))
	// fmt.Println()
	return false
}

// SortWords sorts a list of words
//
// It splits a tab-delimited text file into lines, and sorts by the
// first word in each line (useful for dictionary/glossary/encyclopedia sorting)
//
// If param trud is false, custom alphabetical order is used
// If param trud is true, traditional alphabetical order is used.
func SortWords(letters []rune) [][]string {

	// Create map of all letters in order
	alphabet = make(map[rune]int)
	for index, r := range letters {
		alphabet[r] = index
	}

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
	sort.Sort(CustomAlphabeticalOrder(lines))

	// Return results as 2D array of strings, sorted by the first string in each subarray
	var splitLines [][]string
	for i := 0; i < len(lines); i++ {
		splitLines = append(splitLines, strings.Split(lines[i], "\t"))
	}

	return splitLines
}

// SortByDistance sorts a list of words by Levenshtein distance
func SortByDistance() [][]string {
	// Open file containing rows of words
	file, err := os.Open("C:/Users/Nils/Go/io/words_to_sort.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fmt.Println("file loaded")
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// adds extra distance added for quick sorting (this will be eliminated when database is implemented)
		lines = append(lines, strings.Split(scanner.Text(), "\t")[2]+"\t"+scanner.Text())
	}
	sort.Strings(lines)
	fmt.Println("lines split; test: \n" + lines[0] + "\n" + lines[1] + "\n" + lines[2] + "\n")
	// Return results as 2D array of strings, sorted by levenshtein distance
	var splitLines [][]string
	for i := 0; i < len(lines); i++ {
		// removes extra distance added for quick sorting (this will be eliminated when database is implemented)
		trimmedLine := strings.Split(lines[i], "\t")
		// fmt.Printf("%s %s %s %s\n", trimmedLine[1], trimmedLine[2], trimmedLine[3], trimmedLine[4])
		trimmedLine = []string{trimmedLine[1], trimmedLine[2], trimmedLine[3], trimmedLine[4]}
		// fmt.Printf("%s %s %s %s\n", trimmedLine[0], trimmedLine[1], trimmedLine[2], trimmedLine[3])
		splitLines = append(splitLines, trimmedLine)
	}

	return splitLines

}
