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
	runes := []rune("aäeoøiuywlrmnbpvfgkdtzsžšh")

	// Create map of all runes to count (same runes, so copy above map)
	var allRunes = make(map[rune]int)
	// Create map of word-initial runes to count
	var wordInit = make(map[rune]int)

	// Initialize maps
	for _, k := range runes {
		allRunes[k] = 0
		wordInit[k] = 0
	}

	// Open file in which to count runes
	file, err := os.Open("words_for_runestats.txt")
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
	for _, k := range runes {
		fmt.Print(string(k), "\t", allRunes[k], "\t", wordInit[k], "\r\n")
	}

	// Save counts for each rune to file

	f, err := os.Create("out.txt")
	check(err)

	defer f.Close()

	for i, k := range runes {
		f.Write([]byte(fmt.Sprintf("%s\t%d\t%d", string(k), allRunes[k], wordInit[k])))
		if i < len(runes)-1 {
			f.Write([]byte("\r\n"))
		}
	}
	fmt.Println("Results saved in out.txt")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
