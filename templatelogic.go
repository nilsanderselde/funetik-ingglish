// Nils Elde
// https://gitlab.com/nilsanderselde

package main

import (
	"html/template"
	"net/http"
	"path/filepath"

	"gitlab.com/nilsanderselde/funetik-ingglish/dbconnect"
	"gitlab.com/nilsanderselde/funetik-ingglish/global"
)



// templateHandler contains all fields needed to process and execute templates
type templateHandler struct {
	filename  string
	templ     *template.Template
	query     string
	queryFrom string
	args      global.TemplateParams
}

// creates random number to prompt page reloads after db changes
func randomRune() string {
	return string(global.CurrRand)
}

// handle http request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	funcMap := template.FuncMap{
		"GetStats":  dbconnect.GetStats,
		"ShowWords": dbconnect.ShowWords,
		"Random":    randomRune,
	}
	if t.filename == "home.html" {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
	}

	// special processing for words list based on query strings
	if t.filename == "words.html" {
		handleWordList(t, r)

		// join the template files for the wordlist
		t.templ = template.Must(template.New(t.filename).Funcs(funcMap).ParseFiles(
			filepath.Join("templates", t.filename),
			filepath.Join("templates", "words_sorted.html"),
			filepath.Join("templates", "_header.html"),
			filepath.Join("templates", "_footer.html"),
		))

	} else {
		if t.filename == "translit.html" {

			cha := make(chan dbconnect.Output)

			go dbconnect.ProcessTrud(cha, r)

			outStruct := <-cha // waits till getA() returns

			t.args.TranslitOutput = outStruct.OutputLines
			t.args.TranslitInput = outStruct.PrevInput
		}

		// not word list, just a regular page, join the template files
		t.templ = template.Must(template.New(t.filename).Funcs(funcMap).ParseFiles(
			filepath.Join("templates", t.filename),
			filepath.Join("templates", "_header.html"),
			filepath.Join("templates", "_footer.html"),
		))
	}
	t.templ.Execute(w, t.args)
}
