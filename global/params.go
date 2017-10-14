// Nils Elde
// https://gitlab.com/nilsanderselde

package global

// TemplateParams encapsulates data to be
// passed to mapped functions in templates
type TemplateParams struct {
	New            bool
	Old            bool
	Dist           bool
	ID             bool
	Reverse        bool
	Query          string
	Start          int
	Num            int
	Sort           string
	CurrentPage    string
	NextPage       string
	PreviousPage   string
	Words          [][]string
	TranslitOutput []string
	TranslitInput  string
	DisplayTrud    bool
	KbdVer         string
	Kbd            [][][]string
}

var (
	// CurrRand stores the current random rune
	CurrRand = 'a'
	// LastRand stores the previous random rune so generator doesn't repeat itself
	LastRand = 'a'
)
