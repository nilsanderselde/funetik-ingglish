// Nils Elde
// https://gitlab.com/nilsanderselde

package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"gitlab.com/nilsanderselde/funetik-ingglish/dbconnect"
	"gitlab.com/nilsanderselde/funetik-ingglish/global"
)

// templateHandler contains all fields needed to process and execute templates
type templateHandler struct {
	filenames []string
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
	// Show 404 for unknown paths
	if t.filenames[0] == "home.html" && r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// List files that are to only show traditional English, with no transliteration option
	if t.filenames[0] != "about.html" {
		t.args.MultipleOrth = true
	}

	// special processing for words list based on query strings
	if t.filenames[0] == "words.html" {
		handleWordList(t, r)
		displayOrth(t, r, true)

		// join the template files for the wordlist
		t.templ = template.Must(template.New(t.filenames[0]).Funcs(funcMap).ParseFiles(
			filepath.Join("templates", t.filenames[0]),
			filepath.Join("templates", "words_sorted.html"),
			filepath.Join("templates", "_header.html"),
			filepath.Join("templates", "_footer.html"),
		))

	} else {
		// for keyboard page, decide which keyboard to display based on query string,
		// then decide which orthagraphy to use based on concatenative query string
		if t.filenames[0] == "kbd.html" {
			pickKeyboard(t, r)
			displayOrth(t, r, true)
		} else {
			// for pages without their own, unique query strings,
			// decide which orthagraphy to use based on a new query string (non-concat)
			displayOrth(t, r, false)

			// for tranliteration page, open channel, send input,
			// wait for output
			if t.filenames[0] == "translit.html" {
				cha := make(chan dbconnect.Output)
				go dbconnect.ProcessTrud(cha, r)
				outStruct := <-cha
				t.args.TranslitOutput = outStruct.OutputLines
				t.args.TranslitInput = outStruct.PrevInput
			}
		}
		templateName := strings.TrimSuffix(t.filenames[0], "*.html")
		for i := range t.filenames {
			t.filenames[i] = "templates/" + t.filenames[i]
		}
		t.filenames = append(t.filenames, "templates/_header.html", "templates/_footer.html")

		// not word list, just a regular page, join the template files
		t.templ = template.Must(template.New(templateName).Funcs(funcMap).ParseFiles(t.filenames...))
	}
	t.templ.Execute(w, t.args)
}

// displayOrth determines which orthography to display page in based on query string.
// concat is true if query string should be concatenated to an existing query string
// using & instead of ?
func displayOrth(t *templateHandler, r *http.Request, additive bool) {
	t.args.ChangeOrth = t.args.CurrentPage
	if r.URL.Query()["orth"] != nil && r.URL.Query()["orth"][0] == "fun" {
		if additive {
			t.args.CurrentPage += "&orth=fun"
			t.args.ChangeOrth += "&orth=trud"
		} else {
			t.args.CurrentPage = "?orth=fun"
			t.args.ChangeOrth = "?orth=trud"
		}
		t.args.DisplayTrud = false
	} else {
		if additive {
			t.args.CurrentPage += "&orth=trud"
			t.args.ChangeOrth += "&orth=fun"
		} else {
			t.args.CurrentPage = "?orth=trud"
			t.args.ChangeOrth = "?orth=fun"
		}
		t.args.DisplayTrud = true
	}
}
