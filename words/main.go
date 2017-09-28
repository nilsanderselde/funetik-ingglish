package words

import (
	"fmt"

	"gitlab.com/nilsanderselde/funetik-ingglish/dbconnect"
	"gitlab.com/nilsanderselde/funetik-ingglish/params"
)

// GetWords calls the database io function and pass arguments to template
func GetWords(args params.Params) [][]string {
	fmt.Println(args.Query)
	fmt.Println(args.Start)
	fmt.Println(args.Num)
	words := dbconnect.PostgresIO(args.Query, args.Start, args.Num)
	fmt.Println("test2")

	return words
}
