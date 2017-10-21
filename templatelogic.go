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
	args      global.TemplateParams
}

// creates random number to prompt page reloads after db changes
func randomRune() string {
	return string(global.CurrRand)
}

// handle http request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	filename := strings.TrimPrefix(t.filenames[0], "templates/")
	var additive bool
	t.args.Root = ROOT
	t.args.IsDev = global.IsDev

	switch filename {
	case "about.html":
		t.args.TitleTrud = "About"
		t.args.TitleFun = "Ubäwt"
	case "home.html":
		if r.URL.Path != t.args.Root+"/" {
			http.NotFound(w, r)
			return
		}
		t.args.TitleTrud = "Home"
		t.args.TitleFun = "Hom"
	case "kbd.html":
		pickKeyboard(t, r)
		additive = true
		t.args.TitleTrud = "Keyboard"
		t.args.TitleFun = "Kybord"
	case "stats.html":
		t.args.Stats = global.Stats
		t.args.TitleTrud = "Stats"
		t.args.TitleFun = "Stäts"
	case "translit.html":
		ch := make(chan dbconnect.Output)
		go dbconnect.ProcessTrud(ch, r)
		outStruct := <-ch
		t.args.TranslitOutput = outStruct.OutputLines
		t.args.TranslitInput = outStruct.PrevInput
		t.args.TitleTrud = "Transliterator"
		t.args.TitleFun = "Tränzlitøreitør"
	case "words.html":
		handleWordList(t, r)
		additive = true
		t.args.InitialIndex = global.InitialIndex
		t.args.TitleTrud = "Words"
		t.args.TitleFun = "Wørdz"
	}
	displayOrth(t, r, additive)

	t.once.Do(func() {
		funcMap := template.FuncMap{
			"ShowWords": dbconnect.ShowWords,
			"Random":    randomRune,
		}
		templateName := strings.TrimSuffix(t.filenames[0], "*.html")
		t.filenames = append(t.filenames, "_header.html", "_footer.html")
		for i := range t.filenames {
			t.filenames[i] = "templates/" + t.filenames[i]
		}
		t.templ = template.Must(template.New(templateName).Funcs(funcMap).ParseFiles(t.filenames...))
		fmt.Println("Parsing", t.filenames)
	})
	t.templ.Execute(w, t.args)
}

// displayOrth determines which orthography to display page in based on query string.
// concat is true if query string should be concatenated to an existing query string
// using & instead of ?
func displayOrth(t *templateHandler, r *http.Request, additive bool) {

	t.args.ChangeOrth = t.args.CurrentPage
	if additive {
		t.args.CurrentPage += "&orth="
		t.args.ChangeOrth += "&orth="
	} else {
		t.args.CurrentPage = "?orth="
		t.args.ChangeOrth = "?orth="
	}
	if r.URL.Query()["orth"] != nil && r.URL.Query()["orth"][0] == "fun" {
		t.args.CurrentPage += "fun"
		t.args.ChangeOrth += "trud"
		t.args.DisplayTrud = false
	} else {
		t.args.CurrentPage += "trud"
		t.args.ChangeOrth += "fun"
		t.args.DisplayTrud = true
	}
}
