package words

import (
	"gitlab.com/nilsanderselde/funetik-ingglish/dbconnect"
	"gitlab.com/nilsanderselde/funetik-ingglish/params"
)

// GetWords calls the database io function and pass arguments to template
func GetWords(args params.TemplateParams) params.TemplateParams {
	args.Words = dbconnect.PostgresIO(args.Query, args.Start, args.Num)
	return args
}
