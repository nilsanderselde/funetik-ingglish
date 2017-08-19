// nils elde
// count occurrences of characters in a text file given an alphabet

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	// runes in alphabetical order
	runes := []rune{
		'a', 'ā', 'æ', 'b', 'ч', 'd', 'ð',
		'e', 'ē', 'f', 'g', 'h', 'i', 'ī', 'j',
		'k', 'l', 'm', 'n', 'o', 'p', 'r', 's',
		'ʃ', 't', 'θ', 'u', 'ū', 'v', 'w', 'y', 'z', 'ʒ',
	}

	// create map of all runes to count (same runes, so copy above map)
	var allRunes = make(map[rune]int)
	// create map of word-initial runes to count
	var wordInit = make(map[rune]int)

	for _, k := range runes {
		allRunes[k] = 0
		wordInit[k] = 0
	}

	// get first letter of each word for counts
	file, err := os.Open("uniquewords.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		current := scanner.Text()

		// add all letters to map of total letter counts
		index := 0
		for _, element := range current {
			// for each rune in the string, if it's in the map
			if _, ok := allRunes[element]; ok {
				//increment its count
				allRunes[element]++
				// if first, add first letter to map of first letter counts
				if index == 0 {
					wordInit[element]++
					index++
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// print counts for each rune
	fmt.Println("rune,total,initial")
	for _, k := range runes {
		fmt.Print(string(k), ",", allRunes[k], ",", wordInit[k], "\n")
	}
}
