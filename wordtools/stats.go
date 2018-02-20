// Nils Elde
// https://github.com/nilsanderselde/funetik-ingglish

package wordtools

import (
	"fmt"
	"strings"
)

// CountLetters counts occurrences of specified characters in a slice of strings and
// returns the results in a slice of slices of strings.
func CountLetters(words []string) [][]string {
	// Runes in alphabetical order
	runes := []rune("aäeiywuøorlnmbpvfgkdtzsžšh")
	totalRunes := 0

	// Create map of all runes to count (same runes, so copy above map)
	var allRunes = make(map[rune]int)
	// Create map of word-initial runes to count
	var wordInit = make(map[rune]int)

	// Initialize maps
	for _, r := range runes {
		allRunes[r] = 0
		wordInit[r] = 0
	}

	for _, word := range words {
		index := 0
		for _, element := range strings.ToLower(word) {
			// For each rune in the string, if it's in the map...
			if _, ok := allRunes[element]; ok {
				// ...increment total count
				totalRunes++
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

	results := make([][]string, len(runes))

	for i, k := range runes {
		results[i] = make([]string, 6)
		results[i][0] = string(k)
		results[i][1] = fmt.Sprintf("%d", allRunes[k])
		results[i][2] = fmt.Sprintf("%.1f%%", float64(allRunes[k])/float64(totalRunes)*100)
		results[i][3] = fmt.Sprintf("%d", wordInit[k])
		results[i][4] = fmt.Sprintf("%.1f%%", float64(wordInit[k])/float64(allRunes[k])*100)
		results[i][5] = fmt.Sprintf("%d", i+1)

	}
	return results
}

// CountPhonemes counts occurrences of specified phonemes in a slice of strings and
// returns the results in a slice of slices of strings.
func CountPhonemes(words []string) [][]string {
	// Runes in alphabetical order
	runes := []rune("aāåäeēiyjwʍuøoōrlnŋmbpvfgkdʤðtʧθzsžšh")

	totalRunes := 0

	// Create map of all runes to count (same runes, so copy above map)
	var allRunes = make(map[rune]int)
	// Create map of word-initial runes to count
	var wordInit = make(map[rune]int)

	// Initialize maps
	for _, r := range runes {
		allRunes[r] = 0
		wordInit[r] = 0
	}

	for _, word := range words {
		index := 0
		for _, element := range strings.ToLower(word) {
			// For each rune in the string, if it's in the map...
			if _, ok := allRunes[element]; ok {
				// ...increment total count
				totalRunes++
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

	results := make([][]string, len(runes))

	for i, r := range runes {
		var rwn string
		switch r {
		case 'ŋ':
			rwn = "ng"
		case 'ð':
			rwn = "dh"
		case 'θ':
			rwn = "th"
		case 'ʤ':
			rwn = "dž"
		case 'ʧ':
			rwn = "tš"
		case 'ā':
			rwn = "ai"
		case 'å':
			rwn = "aw"
		case 'ē':
			rwn = "ei"
		case 'ō':
			rwn = "oi"
		case 'j':
			rwn = "y (K)"
		case 'ʍ':
			rwn = "w (K)"
		case 'y':
			rwn = "y (V)"
		case 'w':
			rwn = "w (V)"
		default:
			rwn = string(r)
		}

		results[i] = make([]string, 6)
		results[i][0] = rwn
		results[i][1] = fmt.Sprintf("%d", allRunes[r])
		results[i][2] = fmt.Sprintf("%.1f%%", float64(allRunes[r])/float64(totalRunes)*100)
		results[i][3] = fmt.Sprintf("%d", wordInit[r])
		results[i][4] = fmt.Sprintf("%.1f%%", float64(wordInit[r])/float64(allRunes[r])*100)
		results[i][5] = fmt.Sprintf("%d", i+1)
	}
	return results
}
