// Nils Elde
// https://gitlab.com/nilsanderselde

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"sync"

	"gitlab.com/nilsanderselde/funetik-ingglish/dbconnect"
	"gitlab.com/nilsanderselde/funetik-ingglish/global"
)

// templateHandler contains all fields needed to process and execute templates
type templateHandler struct {
	once      sync.Once
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
	// Show 404 for unknown paths
	if strings.TrimPrefix(t.filenames[0], "templates/") == "home.html" && r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// List files that are to only show traditional English, with no transliteration option
	if strings.TrimPrefix(t.filenames[0], "templates/") != "about.html" {
		t.args.MultipleOrth = true
	}

	fmt.Println(t.filenames[0])
	// special processing for words list based on query strings
	if strings.TrimPrefix(t.filenames[0], "templates/") == "words.html" {
		handleWordList(t, r)
		displayOrth(t, r, true)

	} else {
		// for keyboard page, decide which keyboard to display based on query string,
		// then decide which orthagraphy to use based on concatenative query string
		if strings.TrimPrefix(t.filenames[0], "templates/") == "kbd.html" {
			pickKeyboard(t, r)
			displayOrth(t, r, true)
		} else {
			// for pages without their own, unique query strings,
			// decide which orthagraphy to use based on a new query string (non-concat)
			displayOrth(t, r, false)

			// for tranliteration page, open channel, send input,
			// wait for output
			if strings.TrimPrefix(t.filenames[0], "templates/") == "translit.html" {
				cha := make(chan dbconnect.Output)
				go dbconnect.ProcessTrud(cha, r)
				outStruct := <-cha
				t.args.TranslitOutput = outStruct.OutputLines
				t.args.TranslitInput = outStruct.PrevInput
			}
		}

	}

	t.once.Do(func() {

		funcMap := template.FuncMap{
			"GetStats":  dbconnect.GetStats,
			"ShowWords": dbconnect.ShowWords,
			"Random":    randomRune,
		}

		templateName := strings.TrimSuffix(t.filenames[0], "*.html")
		for i := range t.filenames {
			t.filenames[i] = "templates/" + t.filenames[i]
		}
		t.filenames = append(t.filenames, "templates/_header.html", "templates/_footer.html")

		t.templ = template.Must(template.New(templateName).Funcs(funcMap).ParseFiles(t.filenames...))

	})
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
