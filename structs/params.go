package structs

// TemplateParams encapsulates data to be
// passed to mapped functions in templates
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
