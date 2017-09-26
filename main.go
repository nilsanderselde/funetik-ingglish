package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"gitlab.com/nilsanderselde/funetik-ingglish/customsort"
	"gitlab.com/nilsanderselde/funetik-ingglish/levdist"
	"gitlab.com/nilsanderselde/funetik-ingglish/params"
	"gitlab.com/nilsanderselde/funetik-ingglish/runestats"
)

type templateHandler struct {
	// once      sync.Once
	filename  string
	templ     *template.Template
	path      string
	args      params.Params
	sortBy    string // template filename fragment
	sortOrder string // template filename fragment
}

// handle http request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// t.once.Do(func() {
	funcMap := template.FuncMap{
		"GetStats":       runestats.GetStats,
		"GetDistances":   levdist.GetDistances,
		"SortWords":      customsort.SortWords,
		"SortByTrud":     customsort.SortByTrud,
		"SortByDistance": customsort.SortByDistance,
	}

	if t.path == "/" {
		t.sortBy = "words_sort_fun.html"
		t.args.Fun = true
		t.args.Altfun = false
		t.args.Trud = false
		t.args.Dist = false
		t.args.Order = []rune("aäeiywuøolrmnbpvfgkdtzsžšh")
		t.args.Reverse = false

		if r.URL.Query()["sortby"] != nil {
			if r.URL.Query()["sortby"][0] == "altfun" {
				t.args.Fun = false
				t.args.Altfun = true
				// t.args.Trud = false
				// t.args.Dist = false
				t.args.Order = []rune("aäbdefghiklmnoøprsštuvwyzž") // sudo-trudišinul ordør
			} else if r.URL.Query()["sortby"][0] == "trud" {
				t.args.Fun = false
				// t.args.Altfun = false
				t.args.Trud = true
				// t.args.Dist = false
				t.sortBy = "words_sort_trud.html"
			} else if r.URL.Query()["sortby"][0] == "dist" {
				t.args.Fun = false
				// t.args.Altfun = false
				// t.args.Trud = false
				t.args.Dist = true
				t.sortBy = "words_sort_dist.html"
			}
		}
		if r.URL.Query()["order"] != nil {
			if r.URL.Query()["order"][0] == "desc" {
				t.args.Reverse = true
			} else {
				// t.args.Reverse = false
				// fmt.Println("test")
			}
		}
		t.templ = template.Must(template.New(t.filename).Funcs(funcMap).ParseFiles(
			filepath.Join("templates", t.filename),
			filepath.Join("templates", t.sortBy),
			filepath.Join("templates", "_header.html"),
			filepath.Join("templates", "_footer.html"),
		))
	} else {
		t.templ = template.Must(template.New(t.filename).Funcs(funcMap).ParseFiles(
			filepath.Join("templates", t.filename),
			filepath.Join("templates", "_header.html"),
			filepath.Join("templates", "_footer.html"),
		))
	}
	// })
	// do once prevents dynamic rendering of templates based on query string
	t.templ.Execute(w, t.args)

}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/", &templateHandler{filename: "words.html", path: "/"})
	http.Handle("/runestats", &templateHandler{filename: "runestats.html"})
	http.Handle("/levdist", &templateHandler{filename: "levdist.html"})

	// start server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
