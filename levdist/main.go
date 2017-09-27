// Nils Elde
// https://gitlab.com/nilsanderselde
//
// Calculates Levenshtein distances between words stored in a tabular text file.
// Includes commented out print statements which can help demonstrate how this
// algorithm works

package levdist

import (
	"bufio"
	"log"
	"os"
	"strconv"
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
		distance := FindDistance([]rune(words[0]), []rune(words[1]), true)
		results = append(results, []string{ /*words[0]*/ "" /*words[1]*/, "", strconv.Itoa(distance)})
	}
	return results
}

// FindDistance calculates Levenshtein distance between two words
// If flipping is true, flipping two adjacent letters counts as one move
func FindDistance(word1 []rune, word2 []rune, flipping bool) int {

	// Get length of each word
	var length1 = len(word1)
	var length2 = len(word2)
	// fmt.Printf("\nlen(%s) = %d\n", string(word1), length1)
	// fmt.Printf("len(%s) = %d\n", string(word2), length2)

	// Create a 2D array whose dimensions are the length of each word + 1
	// fmt.Println("initial array: ")
	var pathArray = make([][]int, length1+1)
	for i := 0; i < length1+1; i++ {
		pathArray[i] = make([]int, length2+1)
	}
	for i := 0; i < length1+1; i++ {
		pathArray[i][0] = i // column 0: 0,1,2,3,4,...
	}
	for j := 0; j < length2+1; j++ {
		pathArray[0][j] = j // row 0: 0,1,2,3,4,...
	}
	// for _, element := range pathArray {
	// 	fmt.Println(element)
	// }

	// Compare each letter in first word with each letter in second word
	for i := 0; i < length1; i++ {
		for j := 0; j < length2; j++ {
			rune1 := word1[i]
			rune2 := word2[j]

			// Deleting a letter
			del := pathArray[i][j+1] + 1
			// fmt.Print("del = " + strconv.Itoa(del))

			// Inserting a letter
			ins := pathArray[i+1][j] + 1
			// fmt.Print("; ins = " + strconv.Itoa(ins))

			// Replacing a letter
			rep := pathArray[i][j]
			if rune1 != rune2 {
				rep++
			}
			// fmt.Print("; rep = " + strconv.Itoa(rep))

			// Flipping letters (ab -> ba) (if enabled)
			flp := rep + 1
			if flipping && i+1 > 1 && j+1 > 1 {
				if word1[i-1] == rune2 && word2[j-1] == rune1 {
					flp = pathArray[i-1][j-1] + 1
				}
			}
			// fmt.Print("; flp = " + strconv.Itoa(flp))

			// Set current array value to whichever path is shortest
			min1 := del
			if ins < del {
				min1 = ins
			}
			min2 := rep
			if flp < rep {
				min2 = flp
			}
			if min1 < min2 {
				pathArray[i+1][j+1] = min1
				// fmt.Println("; min = " + strconv.Itoa(min1))
			} else {
				pathArray[i+1][j+1] = min2
				// fmt.Println("; min = " + strconv.Itoa(min2))
			}
		}
	}
	// fmt.Println("final array: ")
	// for _, element := range pathArray {
	// 	fmt.Println(element)
	// }

	return pathArray[length1][length2]
}
