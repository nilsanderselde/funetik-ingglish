package wordtools

// RemoveSyllableMarkers generates "fun" field value by stripping "funsil" of syllable markers ˈˌ·
func RemoveSyllableMarkers(funsil string) (fun string) {
	funsilRunes := []rune(funsil)

	for _, r := range funsilRunes {
		if r != 'ˈ' && r != 'ˌ' && r != '·' {
			fun += string(r)
		}
	}
	return fun
}
