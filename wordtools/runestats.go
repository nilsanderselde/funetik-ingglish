// Nils Elde
// https://gitlab.com/nilsanderselde

package wordtools

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// GetStats counts occurrences of characters in a text file given an alphabet and
// saves results to different text file.
func GetStats() [][]string {
	// Runes in alphabetical order
	runes := []rune("aäeoøiuywlrmnbpvfgkdtzsžšh")

	// Create map of all runes to count (same runes, so copy above map)
	var allRunes = make(map[rune]int)
	// Create map of word-initial runes to count
	var wordInit = make(map[rune]int)

	// Initialize maps
	for _, r := range runes {
		allRunes[r] = 0
		wordInit[r] = 0
	}

	// Open file in which to count runes
	file, err := os.Open("C:/Users/Nils/Go/io/words_for_runestats.txt")

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

	results := make([][]string, len(runes))

	for i, k := range runes {
		results[i] = make([]string, 4)
		results[i][0] = string(k)
		results[i][1] = fmt.Sprintf("%d", allRunes[k])
		results[i][2] = fmt.Sprintf("%d", wordInit[k])
		results[i][3] = fmt.Sprintf("%.1f%%", float64(wordInit[k])/float64(allRunes[k])*100)
	}
	return results
}
