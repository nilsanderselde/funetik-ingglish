package wordtools

import "strings"

// SubstitutionCypher substitutes letters from the first row below to the letter
// directly below it to allow sorting based on a custom alphabet in SQL.
// Results sorted in column "funsort"
func SubstitutionCypher(fun string) (funsort string) {
	funRunes := []rune(strings.ToLower(fun))

	cypher := [][]rune{
		[]rune("aäeiywuøolrmnbpvfgkdtzsžšh"),
		[]rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")}

	for _, r1 := range funRunes {
		for i, r2 := range cypher[0] {
			if r1 == r2 {
				funsort += string(cypher[1][i])
			}
		}
	}
	return funsort
}
