// Nils Elde
// https://gitlab.com/nilsanderselde
// Count occurrences of characters in a text file given an alphabet.
// Used to determine location of new keys in keyboard layout.

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	// Runes in alphabetical order
	runes := []rune{
		'a', 'á', 'b', 'c', 'd', 'e', 'é', 'f', 'g', 'i', 'í', 'j', 'k', 'l', 'm',
		'n', 'o', 'ó', 'ø', 'p', 'r', 's', 'š', 't', 'u', 'ú', 'v', 'z', 'ž', 'h',
	}

	// Create map of all runes to count (same runes, so copy above map)
	var allRunes = make(map[rune]int)
	// Create map of word-initial runes to count
	var wordInit = make(map[rune]int)

	// Initialize maps
	for _, k := range runes {
		allRunes[k] = 0
		wordInit[k] = 0
	}

	// Get first letter of each word for counts
	file, err := os.Open("dictionary.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		current := scanner.Text()

		// Add all letters to map of total letter counts
		index := 0
		for _, element := range current {
			// For each rune in the string, if it's in the map...
			if _, ok := allRunes[element]; ok {
				// ...increment its count
				allRunes[element]++
				// If first, add first letter to map of first letter counts
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

	// Print counts for each rune
	fmt.Println("rune,total,initial")
	for _, k := range runes {
		fmt.Print(string(k), ",", allRunes[k], ",", wordInit[k], "\n")
	}
}
