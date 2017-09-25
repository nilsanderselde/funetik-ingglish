package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"gitlab.com/nilsanderselde/funetik-ingglish/customsort"
	"gitlab.com/nilsanderselde/funetik-ingglish/levdist"
	"gitlab.com/nilsanderselde/funetik-ingglish/runestats"
)

type templateHandler struct {
	// once     sync.Once
	filename string
	templ    *template.Template
	order    []rune
	path     string
}

// handle http request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	funcMap := template.FuncMap{
		"GetStats":       runestats.GetStats,
		"GetDistances":   levdist.GetDistances,
		"SortWords":      customsort.SortWords,
		"SortByDistance": customsort.SortByDistance,
	}

	// var layout string

	if t.path == "/customsort" {
		fmt.Println(t.path)
		if r.URL.Query()["fun"] != nil {
			t.order = []rune("aäeiywuøolrmnbpvfgkdtzsžšh")
			t.filename = "words.html"
			fmt.Println("words")
		}

		if r.URL.Query()["trud"] != nil {
			t.order = []rune("aäbdefghiklmnoøprsštuvwyzž")
			t.filename = "customsort.html"
		}

		if r.URL.Query()["dist"] != nil {
			t.filename = "distsort.html"
		}

		// wordsLayout, err := ioutil.ReadFile(filepath.Join("templates", "words.html"))
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// layout = string(wordsLayout)
	}

	// t.once.Do(func() { // do once prevents dynamic rendering of templates based on query string

	// data, err := ioutil.ReadFile(filepath.Join("templates", t.filename))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// filepath.Join("templates", "words.html"),

	t.templ = template.Must(template.New(t.filename).Funcs(funcMap).ParseFiles(filepath.Join("templates", t.filename)))
	// })

	t.templ.Execute(w, t.order)
}

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/", &templateHandler{filename: "menu.html"})
	http.Handle("/runestats", &templateHandler{filename: "runestats.html"})
	http.Handle("/levdist", &templateHandler{filename: "levdist.html"})
	http.Handle("/customsort", &templateHandler{path: "/customsort"})

	// start server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
