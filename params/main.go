package params

// Params encapsulates data to be passed to mapped functions
// in templates
type Params struct {
	New          bool
	Old          bool
	Dist         bool
	Reverse      bool
	Query        string
	Start        int
	Num          int
	CurrentPage  string
	NextPage     string
	PreviousPage string
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
