package words

import (
	"gitlab.com/nilsanderselde/funetik-ingglish/dbconnect"
	"gitlab.com/nilsanderselde/funetik-ingglish/params"
)

// GetWords calls the database io function and passes arguments to template
func GetWords(args params.TemplateParams) [][]string {
	return dbconnect.GetWords(args.Query, args.Start, args.Num)
}
