// Código baseado em exemplo: https://golang.org/doc/articles/wiki/

package main

import (
	"html/template"
	_ "io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var validPath = regexp.MustCompile("^/([a-zA-Z0-9]+)$")

func makeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.Redirect(w, r, "/index", 302)
			return
		}
		var pages = template.Must(template.ParseFiles("pages/index.html"))
		var pageName = m[1]

		if pageName == "index" {
			r.ParseForm()
			var filmes = r.Form["filme"]
			if filmes != nil {
				var search = filmes[0]
				executarScraping(search)
			}
		}

		err := pages.ExecuteTemplate(w, pageName+".html", "data")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func main() {
	http.HandleFunc("/", makeHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
