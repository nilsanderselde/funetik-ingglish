// Nils Elde
// https://gitlab.com/nilsanderselde

package global

var (
	// IsDev toggles configuration of program in several places depending on if
	// it is in development or production mode
	IsDev bool
	// CurrRand stores the current random rune
	CurrRand = 'a'
	// LastRand stores the previous random rune so generator doesn't repeat itself
	LastRand = 'a'
	// InitialIndex stores the offset amounts for the first word starting with each letter
	// of the alphabet, to enable browsing by letter on the words page
	InitialIndex []InitialIndexValue
	// Alphabet lists Funetik Inggliš letters in order. Used to allow jumping to letter on word page.
	Alphabet = []string{"a", "ä", "e", "i", "y", "w", "u", "ø", "o", "r", "l", "n", "m", "b", "p", "v", "f", "g", "k", "d", "t", "z", "s", "ž", "š", "h"}
	// RowCount counts number of rows in DB
	RowCount int
	// Stats stores the stats about the words in the database
	Stats [][]string
)

// InitialIndexValue stores the start number for the first word starting with the stored letter
type InitialIndexValue struct {
	Letter string
	Index  string
}

// TemplateParams encapsulates data to be
// passed to mapped functions in templates
type TemplateParams struct {
	// words page
	SortBy       string
	Reverse      bool
	Start        int
	Num          int
	Alphabet     []string
	SortQ        string
	PQuery       string
	InitialIndex []InitialIndexValue
	NextPage     string
	PreviousPage string

	// stats page
	Stats [][]string

	// translit page
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
	TitleTrud   string
	TitleFun    string
	Root        string

	// Development mode on or off
	IsDev bool
}
