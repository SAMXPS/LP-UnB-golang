// Código baseado em exemplo: https://golang.org/doc/articles/wiki/

package main

import (
	"html/template"
	"io/ioutil"
	_ "io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sync"
)

type ResultadoPesquisa struct {
	Resultados [3]string
}

var validPath = regexp.MustCompile("^/([a-zA-Z0-9]+)$")

func makeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.Redirect(w, r, "/index", http.StatusFound)
			return
		}
		var pages = template.Must(template.ParseFiles("pages/index.html", "pages/resultado.html"))
		var pageName = m[1]
		var executado = false
		var err error
		err = nil

		if pageName == "index" {
			r.ParseForm()
			var filmes = r.Form["filme"]
			if filmes != nil {
				var search = filmes[0]
				var wait_group sync.WaitGroup
				executado = true

				context, err2 := executarScrapingPersonalizado(search, &wait_group)
				if err2 != nil {
					err = err2
				} else {
					encerrarContext(context)
					var resultados [3]string
					for i := 2; i < 5; i++ {
						content, err3 := ioutil.ReadFile(context.files[i].Name())
						if err3 != nil {
							err = err3
							break
						} else {
							resultados[i-2] = (string(content))
						}
					}
					if err == nil {
						dados := ResultadoPesquisa{
							Resultados: resultados,
						}
						err = pages.ExecuteTemplate(w, "resultado.html", dados)
					}
				}

			}
		}

		if !executado {
			err = pages.ExecuteTemplate(w, pageName+".html", "data")
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func executarScrapingPersonalizado(termo_pesquisa string, wait_group *sync.WaitGroup) (*scrap_context, error) {
	context, err := criarContexto()
	if err != nil {
		return nil, err
	} else {
		// Criamos o contexto de scraping
		configurarContexto(context)

		wait_group.Add(1)

		// Rodamos as tarefas pendentes
		go realizarPesquisaParalelo(termo_pesquisa, context, wait_group)

		// Esperamos todas as tarefas finalizarem
		wait_group.Wait()

		return context, nil
	}
}

func main() {
	http.HandleFunc("/", makeHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
