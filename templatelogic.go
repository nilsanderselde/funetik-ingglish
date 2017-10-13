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

		displayOrth(t, r, true)

		// join the template files for the wordlist
		t.templ = template.Must(template.New(t.filename).Funcs(funcMap).ParseFiles(
			filepath.Join("templates", t.filename),
			filepath.Join("templates", "words_sorted.html"),
			filepath.Join("templates", "_header.html"),
			filepath.Join("templates", "_footer.html"),
		))

	} else if t.filename == "keyboard.html" {
		pickKeyboard(t, r)
		displayOrth(t, r, true)

		// join the template files for the wordlist
		t.templ = template.Must(template.New(t.filename).Funcs(funcMap).ParseFiles(
			filepath.Join("templates", t.filename),
			filepath.Join("templates", "keyboards", "1.html"),
			filepath.Join("templates", "keyboards", "2.html"),
			filepath.Join("templates", "keyboards", "3.html"),
			filepath.Join("templates", "keyboards", "4.html"),
			filepath.Join("templates", "keyboards", "5.html"),
			filepath.Join("templates", "keyboards", "6.html"),
			filepath.Join("templates", "keyboards", "7.html"),
			filepath.Join("templates", "keyboards", "8.html"),
			filepath.Join("templates", "keyboards", "9.html"),
			filepath.Join("templates", "keyboards", "10.html"),
			filepath.Join("templates", "keyboards", "11.html"),
			filepath.Join("templates", "keyboards", "12.html"),
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

		displayOrth(t, r, false)

		// not word list, just a regular page, join the template files
		t.templ = template.Must(template.New(t.filename).Funcs(funcMap).ParseFiles(
			filepath.Join("templates", t.filename),
			filepath.Join("templates", "_header.html"),
			filepath.Join("templates", "_footer.html"),
		))
	}
	t.templ.Execute(w, t.args)
}

// displayOrth determines which orthography to display page in based on query string.
// concat is true if query string should be concatenated to an existing query string
// using & instead of ?
func displayOrth(t *templateHandler, r *http.Request, additive bool) {

	if r.URL.Query()["orth"] != nil && r.URL.Query()["orth"][0] == "trud" {
		if additive {
			t.args.CurrentPage += "&orth=fun"
		} else {
			t.args.CurrentPage = "?orth=fun"
		}
		t.args.DisplayTrud = true
	} else {
		if additive {
			t.args.CurrentPage += "&orth=trud"
		} else {
			t.args.CurrentPage = "?orth=trud"
		}
		t.args.DisplayTrud = false
	}
}

// pickKeyboard determines which keyboard layout to show on keyboard page
func pickKeyboard(t *templateHandler, r *http.Request) {

	if r.URL.Query()["v"] != nil {
		v := &r.URL.Query()["v"][0]

		if *v == "1" || *v == "2" || *v == "3" || *v == "4" ||
			*v == "5" || *v == "6" || *v == "7" || *v == "8" ||
			*v == "9" || *v == "10" || *v == "11" {
			t.args.CurrentPage = "?v=" + *v
		} else {
			t.args.CurrentPage = "?v=12"
		}
	} else {
		t.args.CurrentPage = "?v=12"
	}
}
