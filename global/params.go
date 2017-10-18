// Nils Elde
// https://gitlab.com/nilsanderselde

package global

// TemplateParams encapsulates data to be
// passed to mapped functions in templates
type TemplateParams struct {
	// words page
	SortBy  string
	Reverse bool
	Start   int
	Num     int

	Words  [][]string
	SortQ  string
	PQuery string

	NextPage     string
	PreviousPage string
	// end words page

	//translit page
	TranslitOutput []string
	TranslitInput  string

	// kbd page
	KbdVer string
	Kbd    [][][]string

	// all pages
	CurrentPage string
	SingleOrth  bool
	ChangeOrth  string
	DisplayTrud bool
}

var (
	// CurrRand stores the current random rune
	CurrRand = 'a'
	// LastRand stores the previous random rune so generator doesn't repeat itself
	LastRand = 'a'
)
