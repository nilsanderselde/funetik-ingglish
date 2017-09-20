// Nils Elde
// https://gitlab.com/nilsanderselde
//
// editDistance and editDistStep converted to Go from Python, based on:
//
// Natural Language Toolkit: Distance Metrics
//
// Copyright (C) 2001-2017 NLTK Project
// Author: Edward Loper <edloper@gmail.com>
//         Steven Bird <stevenbird1@gmail.com>
//         Tom Lippincott <tom@cs.columbia.edu>
// URL: <http://nltk.org/>
//
package levdist

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// GetDistances calculates Levenshtein distances between words stored in a tabular text file
func GetDistances() [][]string {
	file, err := os.Open("C:/Users/Nils/Go/io/words_for_distance.txt")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var results [][]string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		current := scanner.Text()
		words := strings.Split(current, "\t")
		distance := EditDistance(words[0], words[1], 1, false)
		results = append(results, []string{words[0], words[1], fmt.Sprintf("%d", distance)})

	}

	return results
}

// EditDistance calculates Levenshtein distance between two words
func EditDistance(s1 string, s2 string, substitutionCost int, transpositions bool) int {
	// convert to rune array to access each unicode letter
	r1 := []rune(s1)
	r2 := []rune(s2)

	// set up a 2-D array
	var len1 = len(r1)
	var len2 = len(r2)
	var lev = make([][]int, len1+1)
	for i := 0; i < len1+1; i++ {
		lev[i] = make([]int, len2+1)
	}

	for i := 0; i < len1; i++ {
		lev[i][0] = i // column 0: 0,1,2,3,4,...
	}

	for j := 0; j < len2; j++ {
		lev[0][j] = j // row 0: 0,1,2,3,4,...
	}

	// iterate over the array
	for i := 0; i < len1; i++ {
		for j := 0; j < len2; j++ {
			editDistStep(lev, i+1, j+1, r1, r2, substitutionCost, transpositions)

			// print out progress in console to illustrate how it works
			// fmt.Print(lev[i][j])
			// if i < len1-1 {
			// 	fmt.Print("\n")
			// }
		}
	}
	return lev[len1][len2]
}

func editDistStep(lev [][]int, i int, j int, r1 []rune, r2 []rune, substitutionCost int, transpositions bool) {
	c1 := r1[i-1]
	c2 := r2[j-1]

	// skipping a character in s1
	a := lev[i-1][j] + 1
	// skipping a character in s2
	b := lev[i][j-1] + 1

	// substitution
	c := lev[i-1][j-1]

	if c1 != c2 {
		c += substitutionCost
	}

	// transposition
	d := c + 1 // never picked by default
	if transpositions && i > 1 && j > 1 {
		if r1[i-2] == c2 && r2[j-2] == c1 {
			d = lev[i-2][j-2] + 1
		}
	}

	// pick the cheapest

	min1 := a
	if b < a {
		min1 = b
	}

	min2 := c
	if d < c {
		min2 = d
	}

	if min1 < min2 {
		lev[i][j] = min1
	} else {
		lev[i][j] = min2
	}

	// lev[i][j] = int(math.Min(math.Min(float64(a), float64(b)), math.Min(float64(c), float64(d))))

}
