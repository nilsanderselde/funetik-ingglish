// Nils Elde
// https://gitlab.com/nilsanderselde

package funetiksort

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var alphabet = "aäeoøiuywlrmnbpvfgkdtzsžšh"

// Alias for array of strings to be sorted
type BaiFunetikOrdør []string

func (s BaiFunetikOrdør) Len() int {
	return len(s)
}

func (s BaiFunetikOrdør) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s BaiFunetikOrdør) Less(i, j int) bool {
	word1 := []rune(strings.Split(s[i], "\t")[0])
	word2 := []rune(strings.Split(s[j], "\t")[0])

	var minlength int
	var longest string
	if len(word1) < len(word2) {
		minlength = len(word1)
		longest = "word2"
	} else if len(word1) > len(word2) {
		minlength = len(word2)
		longest = "word1"
	} else {
		minlength = len(word1)
		longest = "equal length"
	}
	fmt.Println("word1: " + string(word1) + "\tword2: " + string(word2) + "\t minlength: " + strconv.Itoa(minlength) + "\t longest: " + longest)

	for k := 0; k < minlength; k++ {
		letter1 := []rune(strings.ToLower(string(word1)))[k]
		letter1order := strings.IndexRune(alphabet, letter1)
		letter2 := []rune(strings.ToLower(string(word2)))[k]
		letter2order := strings.IndexRune(alphabet, letter2)

		// fmt.Print("k: " + strconv.Itoa(k) + "\trune1: " + string(letter1) + "\torder1: " + strconv.Itoa(letter1order))

		// fmt.Println("\trune2: " + string(letter2) + "\torder2: " + strconv.Itoa(letter2order))

		if letter1order < letter2order {
			return true
		}
		return false
	}
	return false
}

// SortWords sorts a list of words in custom alphabetical order
func SortWords() [][]string {
	// Open file containing rows of words
	file, err := os.Open("C:/Users/Nils/Go/io/words_to_sort.txt")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// fmt.Println(scanner.Text())
		lines = append(lines, scanner.Text())
	}
	sort.Sort(BaiFunetikOrdør(lines))

	var splitLines [][]string
	for i := 0; i < len(lines); i++ {
		splitLines = append(splitLines, strings.Split(lines[i], "\t"))
	}

	return splitLines
}
