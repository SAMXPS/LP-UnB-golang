// Código baseado em exemplo: https://golang.org/doc/articles/wiki/

package main

import (
	"html/template"
	_"io/ioutil"
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
		err := pages.ExecuteTemplate(w, m[1]+".html", "data")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func main() {
	http.HandleFunc("/", makeHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}