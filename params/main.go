package params

// TemplateParams encapsulates data to be passed to mapped functions
// in templates
type TemplateParams struct {
	New          bool
	Old          bool
	Dist         bool
	ID           bool
	Reverse      bool
	Query        string
	Start        int
	Num          int
	Sort         string
	CurrentPage  string
	NextPage     string
	PreviousPage string
	Words        [][]string
}

// Row contains fields for all columns in a row to
// all rows to be updated
type Row struct {
	ID     int
	Funsil string
	Trud   string
	Pus    string
	Ritin  string
	Kamin  bool
	Tshekt bool
	Flaagd bool
}
