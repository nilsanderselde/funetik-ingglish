// nils elde
// count occurrences of characters in a text file given an alphabet

package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	// get bytes from file
	b, err := ioutil.ReadFile("uniquewords.txt")
	if err != nil {
		fmt.Print(err)
	}

	// create string from bytes
	words := string(b)

	// create list of runes to count
	runes := map[rune]int{
		'a': 0,
		'ā': 0,
		'æ': 0,
		'b': 0,
		'd': 0,
		'e': 0,
		'ē': 0,
		'f': 0,
		'g': 0,
		'h': 0,
		'i': 0,
		'ī': 0,
		'k': 0,
		'l': 0,
		'm': 0,
		'n': 0,
		'o': 0,
		'p': 0,
		'r': 0,
		's': 0,
		'ʃ': 0,
		't': 0,
		'u': 0,
		'ū': 0,
		'v': 0,
		'w': 0,
		'y': 0,
		'z': 0,
		'ʒ': 0,
	}

	// for each rune in the string
	for _, element := range words {
		// if it's in the map
		if _, ok := runes[element]; ok {
			//increment its count
			runes[element]++
		}
	}

	// print counts for each rune
	for r, k := range runes {
		fmt.Println(string(r), ": ", k)
	}

}
