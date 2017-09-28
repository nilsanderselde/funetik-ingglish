package params

// Params encapsulates data to be passed to mapped functions
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
