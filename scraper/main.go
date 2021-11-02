package main

import (
    _"encoding/csv"
    "fmt"
    _"log"
    _"os"
    _"strconv"
    "strings"
    "sync"
    "bytes"
    "github.com/gocolly/colly"
)

type scrap_context struct {
	colly_meta 	    *colly.Collector
    colly_imdb 		*colly.Collector
    colly_rotten    *colly.Collector
    colly_letter    *colly.Collector
    response        *bytes.Buffer
}

// Lista de variáveis globais


func criarCollys() scrap_context{
    return scrap_context{
        colly_meta: colly.NewCollector(
            colly.AllowedDomains("www.metacritic.com"),
        ),
        colly_imdb: colly.NewCollector(
            colly.AllowedDomains("www.imdb.com"),
        ),
        colly_rotten: colly.NewCollector(
            colly.AllowedDomains("www.rottentomatoes.com"),
        ),
        colly_letter: colly.NewCollector(
            colly.AllowedDomains("letterboxd.com"),
        ),
        response: bytes.NewBufferString(""),
    }
} 

func configurarCollys(context *scrap_context) {
    (*context).colly_meta.OnHTML(".clamp-summary-wrap", func(e *colly.HTMLElement) {
        fmt.Fprint(context.response, "Resposta meta:")
        /*writer.Write([]string{
            e.ChildText("h3"),
            e.ChildText("span"),
        })*/
    })

    (*context).colly_imdb.OnHTML(".lister-item-content", func(e *colly.HTMLElement) {
        fmt.Fprint(context.response, "Resposta IMDB:")
        /*writer2.Write([]string{
            e.ChildText("h3"),
            e.ChildText("p"),
        })*/
    })

    (*context).colly_meta.OnHTML(".main_stats", func(e *colly.HTMLElement) {
        fmt.Fprint(context.response, "Resposta META2:")
        /*writer3.Write([]string{
            e.ChildText(".metascore_w"),
            e.ChildText("a"),
        })*/
    })

    (*context).colly_rotten.OnHTML("ul search-page-media-row", func(e *colly.HTMLElement) {
        fmt.Fprint(context.response, "Resposta rotten:")
        /*writer4.Write([]string{
            e.ChildAttr("score-icon-critic", "percentage"),
            e.ChildText("a"),
        })*/
    })

    (*context).colly_letter.OnHTML(".film-detail-content", func(e *colly.HTMLElement) {
        fmt.Fprint(context.response, "Resposta letter:")
        /*writer5.Write([]string{
            e.ChildText("a href"),
            e.ChildText(".film-title-wrapper"),
        })*/
    })

}

func realizarPesquisaIndividual(wait_group *sync.WaitGroup, colly_atual* colly.Collector, url_pesquisa string) {
    colly_atual.Visit(url_pesquisa)
    wait_group.Done()
}

func scrapPage(i int, wg *sync.WaitGroup, c *colly.Collector) {
	fmt.Printf("Scraping Page: %d\n", i)
	c.Visit("https://www.metacritic.com/browse/movies/score/metascore/all/filtered/netflix?page=" + strconv.Itoa(i))
	wg.Done()
}

func scrapPage2(i int, wg *sync.WaitGroup, c2 *colly.Collector) {
	fmt.Printf("Scraping Page: %d\n", i)
	c2.Visit("https://www.imdb.com/search/title/?count=100&groups=top_1000&sort=user_rating" + strconv.Itoa(i))
    wg.Done()
}

func realizarPesquisaParalelo(termo_pesquisa string, context *scrap_context) {
    var wait_group sync.WaitGroup
    
    var termo_pesquisa_adaptada = strings.ReplaceAll(termo_pesquisa, " ", "+")
    var termo_pesquisa_adaptada2 = strings.ReplaceAll(termo_pesquisa, " ", "%20")
    
    var pesquisa_meta = "https://www.metacritic.com/search/movie/" + termo_pesquisa_adaptada2 + "/results"
    var pesquisa_rotten = "https://www.rottentomatoes.com/search?search=" + termo_pesquisa_adaptada2
    var pesquisa_letter = "https://letterboxd.com/search/films/" + termo_pesquisa_adaptada + "/"
    
    // Colocamos 3 tasks para aguardar.
    wait_group.Add(3)
    go realizarPesquisaIndividual(&wait_group, context.colly_meta, pesquisa_meta)
    go realizarPesquisaIndividual(&wait_group, context.colly_rotten, pesquisa_rotten)
    go realizarPesquisaIndividual(&wait_group, context.colly_letter, pesquisa_letter)
    
    // Aguardamos as 3 tasks serem completas.
    wait_group.Wait()
}

func montarDatabase() {
    var wait_group sync.WaitGroup
    wait_group.Add(11)

	for i := 0; i < 11; i++ {
		go scrapPage(i, &waitGroup, c)
	}

	wait_group.Add(11)
	for i := 0; i < 11; i++ {
		go scrapPage2(i, &waitGroup2, c2)
	}

    // Aguardamos as tasks serem completas.
    wait_group.Wait()
}

func main() {
    var context = criarCollys()
    configurarCollys(&context)
    realizarPesquisaParalelo("Titanic", &context)
    println(context.response.String())
}
