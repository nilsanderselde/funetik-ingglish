// Nils Elde
// https://gitlab.com/nilsanderselde
//
// Calculates Levenshtein distances between two strings.
//
// Based on the Python implementation in the Natural Language Toolkit,
// part of the nltk.metrics.distance module. License for the NLTK
// Levenshtein implementation below:
//
// Copyright (C) 2001-2017 NLTK Project
// Author: Edward Loper <edloper@gmail.com>
// Steven Bird <stevenbird1@gmail.com>
// Tom Lippincott <tom@cs.columbia.edu>
// URL: <http://nltk.org/>
// For license information, see LICENSE.TXT
//
// LICENSE.TXT
// Copyright (C) 2001-2014 NLTK Project
//
// Licensed under the Apache License, Version 2.0 (the 'License');
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an 'AS IS' BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package wordtools

// FindDistance calculates Levenshtein distance between two words.
// If flipping is true, flipping two adjacent letters counts as one move.
func FindDistance(funS string, trudS string, flipping bool) int {

	fun := []rune(funS)
	trud := []rune(trudS)

	// Get length of each word
	var length1 = len(fun)
	var length2 = len(trud)

	// Create a 2D array whose dimensions are the length of each word + 1
	var pathArray = make([][]int, length1+1)
	for i := 0; i < length1+1; i++ {
		pathArray[i] = make([]int, length2+1)
	}
	for i := 0; i < length1+1; i++ {
		pathArray[i][0] = i
	}
	for j := 0; j < length2+1; j++ {
		pathArray[0][j] = j
	}

	// Compare each letter in first word with each letter in second word
	for i := 0; i < length1; i++ {
		for j := 0; j < length2; j++ {
			rune1 := fun[i]
			rune2 := trud[j]

			// Deleting a letter
			del := pathArray[i][j+1] + 1

			// Inserting a letter
			ins := pathArray[i+1][j] + 1

			// Replacing a letter
			rep := pathArray[i][j]
			if rune1 != rune2 {
				rep++
			}

			// Flipping letters (ab -> ba) (if enabled)
			flp := rep + 1
			if flipping && i+1 > 1 && j+1 > 1 {
				if fun[i-1] == rune2 && trud[j-1] == rune1 {
					flp = pathArray[i-1][j-1] + 1
				}
			}

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
			} else {
				pathArray[i+1][j+1] = min2
			}
		}
	}
	return pathArray[length1][length2]
}
