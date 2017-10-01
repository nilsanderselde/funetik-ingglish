// Nils Elde
// https://gitlab.com/nilsanderselde

package wordtools

// CountSyllables counts occurrences of syllable markers ˈˌ· in a word
func CountSyllables(funsil string) (numsil int) {
	funsilRunes := []rune(funsil)
	count := 0

	for _, r := range funsilRunes {
		if r == 'ˈ' || r == 'ˌ' || r == '·' {
			count++
		}
	}
	return count
}
