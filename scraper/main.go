package main

import (
    "encoding/csv"                  // Utilizado para construir banco CSV
    "fmt"                           // Utilizado para formatação de strings  
    "log"                           // Utilizado para log de mensagens 
    "os"                            // Utilizado para abrir e escrever em arquivos
    "strconv"                       // Utilizado para converter tipos (ex: int) em strings
    "strings"                       // Utilizado para textos em strings
    "sync"                          // Utilizado para paralelismo  
    "github.com/gocolly/colly"      // Utilizado para fazer scraping das páginas  
)

// Struct de contexto para scraping
type scrap_context struct {
	colly_meta 	    *colly.Collector    // Scraper do website metacritc 
    colly_imdb 		*colly.Collector    // Scraper do website IMBD
    colly_rotten    *colly.Collector    // Scraper do website rottentomatoes
    colly_letter    *colly.Collector    // Scraper do website letterboxd
    files           [5]*os.File         // Arquivos temporários CSV
    writers         [5]*csv.Writer      // Writes do CSV
}

// Função que cria contexto de execução do scraping
func criarContexto() (*scrap_context, error){
    var writers [5]*csv.Writer
    var files [5]*os.File

    // Abrimos 5 arquivos e 5 writers
    for i := 0; i < 5; i++ {
        fName := "data" + strconv.Itoa(i) + ".csv"
        file, err := os.Create(fName)
        if err != nil {
            log.Fatalf("Erro ao criar arquivo, err: %q", err)
            return nil, err
        }
        writers[i] = csv.NewWriter(file)
        files[i] = file
	}

    // Inicializamos os scrapers
    context := scrap_context{
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
        writers: writers,
        files: files,
    }
    return &context, nil;
} 

// Função que configura o contexto de scraping
func configurarContexto(context *scrap_context) {
    (*context).colly_meta.OnHTML(".clamp-summary-wrap", func(e *colly.HTMLElement) {
        (*context).writers[0].Write([]string{
            e.ChildText("h3"),
            e.ChildText("span"),
        })
    })

    (*context).colly_imdb.OnHTML(".lister-item-content", func(e *colly.HTMLElement) {
        //fmt.Fprint(context.response, "Resposta IMDB:")
        (*context).writers[1].Write([]string{
            e.ChildText("h3"),
            e.ChildText("p"),
        })
    })

    (*context).colly_meta.OnHTML(".main_stats", func(e *colly.HTMLElement) {
        //fmt.Fprint(context.response, "Resposta META2:")
        (*context).writers[2].Write([]string{
            e.ChildText(".metascore_w"),
            e.ChildText("a"),
        })
    })

    (*context).colly_rotten.OnHTML("ul search-page-media-row", func(e *colly.HTMLElement) {
        //fmt.Fprint(context.response, "Resposta rotten:")
        (*context).writers[3].Write([]string{
            e.ChildAttr("score-icon-critic", "percentage"),
            e.ChildText("a"),
        })
    })

    (*context).colly_letter.OnHTML(".film-detail-content", func(e *colly.HTMLElement) {
        //fmt.Fprint(context.response, "Resposta letter:")
        (*context).writers[4].Write([]string{
            e.ChildText("a href"),
            e.ChildText(".film-title-wrapper"),
        })
    })
}

// Função que faz scraping de página do Metacritic
func realizarScrapMetacritic(i int, wg *sync.WaitGroup, c *colly.Collector) {
	fmt.Printf("Scraping Page: %d\n", i)
	c.Visit("https://www.metacritic.com/browse/movies/score/metascore/all/filtered/netflix?page=" + strconv.Itoa(i))
	wg.Done()
}

// Funçao que faz scraping de página do IMDB
func realizarScrapImdb(i int, wg *sync.WaitGroup, c2 *colly.Collector) {
	fmt.Printf("Scraping Page: %d\n", i)
	c2.Visit("https://www.imdb.com/search/title/?count=100&groups=top_1000&sort=user_rating" + strconv.Itoa(i))
    wg.Done()
}

// Função que realiza pesquisa na url e colly selecionados
func realizarPesquisaIndividual(wait_group *sync.WaitGroup, colly_atual* colly.Collector, url_pesquisa string) {
    colly_atual.Visit(url_pesquisa)
    wait_group.Done()
}

// Função que pesquisa um filme em 3 sites diferentes, de forma paralela
func realizarPesquisaParalelo(termo_pesquisa string, context *scrap_context, wait_group *sync.WaitGroup) {
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
}

// Função que faz scraping e monta banco de dados de melhores filmes
func montarDatabase(context *scrap_context, wait_group *sync.WaitGroup) {
	for i := 0; i < 11; i++ {
        wait_group.Add(2)
		go realizarScrapMetacritic(i, &wait_group, context.colly_meta)
		go realizarScrapImdb(i, &waitGroup, context.colly_imdb)
	}
}

// Função para fechar writers e arquivos
func encerrarContext(context *scrap_context) {
    for i := 0; i < 5; i++ {
        context.writers[i].Flush()
        context.files[i].Close()
    }
}

func main() {
    context, err := criarContexto()
    if err != nil {
        return 
    } else {
        var wait_group sync.WaitGroup

        // Criamos o contexto de scraping
        configurarContexto(context)

        // Rodamos as tarefas pendentes
        go realizarPesquisaParalelo("Titanic", context, &wait_group)
        go montarDatabase(context, &wait_group)

        // Esperamos todas as tarefas finalizarem
        wait_group.Wait()

        // Encerramos o contexto de scraping
        encerrarContext(context)
    }
}
