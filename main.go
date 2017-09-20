package main

// import (
// 	"fmt"

// 	"gitlab.com/nilsanderselde/funetik-ingglish/levdist"
// )

// func main() {
// 	fmt.Println(levdist.EditDistance("at", "ta", 1, true))
// }

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"gitlab.com/nilsanderselde/funetik-ingglish/levdist"
	"gitlab.com/nilsanderselde/funetik-ingglish/runestats"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// handle http request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{
		"GetStats":     runestats.GetStats,
		"GetDistances": levdist.GetDistances,
	}

	t.once.Do(func() {
		data, err := ioutil.ReadFile(filepath.Join("templates", t.filename))
		if err != nil {
			log.Fatal(err)
		}
		t.templ = template.Must(template.New(t.filename).Funcs(funcMap).Parse(string(data)))
	})
	t.templ.Execute(w, nil)
}

func main() {

	http.Handle("/", &templateHandler{filename: "menu.html"})
	http.Handle("/runestats", &templateHandler{filename: "runestats.html"})
	http.Handle("/levdist", &templateHandler{filename: "levdist.html"})

	// start server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
