// Nils Elde
// https://github.com/nilsanderselde/funetik-ingglish

package global

var (
	// IsDev toggles configuration of program in several places depending on if
	// it is in development or production mode
	IsDev bool

	// CurrRand stores the current random rune
	CurrRand = 'a'

	// LastRand stores the previous random rune so generator doesn't repeat itself
	LastRand = 'a'

	// FunetikIndex stores the offset amounts for the first word starting with each letter
	// of the alphabet, to enable browsing by letter on the words page
	FunetikIndex []OrderedIndexMap

	// TrudIndex stores the offset amounts for the first word starting with each letter
	// of the alphabet, to enable browsing by letter on the words page
	TrudIndex []OrderedIndexMap

	// DistIndex stores the offset amounts for the first word starting with each letter
	// of the alphabet, to enable browsing by letter on the words page
	DistIndex []OrderedIndexMap

	// Älfubit lists Funetik Inggliš letters in order. Used to allow jumping to letter on word page.
	Älfubit = []string{"a", "ä", "e", "i", "y", "w", "u", "ø", "o", "r", "l", "n", "m", "b", "p", "v", "f", "g", "k", "d", "t", "z", "s", "ž", "š", "h"}

	// Alphabet lists traditional English letters in order..
	Alphabet = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

	// RowCount stores number of rows in DB
	RowCount int

	// RowCountF stores formatted number of rows in DB
	RowCountF string

	// PhonStats stores the phoneme stats about the words in the database
	PhonStats [][]string

	// RuneStats stores the rune stats about the words in the database
	RuneStats [][]string
)

// OrderedIndexMap stores the start number for the first word starting with the stored value
type OrderedIndexMap struct {
	Value  string
	Offset string
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
	FunetikIndex []OrderedIndexMap
	TrudIndex    []OrderedIndexMap
	DistIndex    []OrderedIndexMap

	NextPage     string
	PreviousPage string

	// stats page
	PhonStats [][]string
	RuneStats [][]string
	RowCount  int
	RowCountF string

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
