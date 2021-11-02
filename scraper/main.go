package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

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

//func scrapPage3(i int, wg *sync.WaitGroup, c3 *colly.Collector) {
//	fmt.Printf("Scraping Page: %d\n", i)
//	c3.Visit("https://www.rottentomatoes.com/browse/dvd-streaming-all" + strconv.Itoa(i))

//	wg.Done()
//}

func main() {
	fName := "data.csv"
	file, err := os.Create(fName)
	//var rows [][]string
	var termo_pesquisa = "dune"
	//var pesquisa = "https://www.metacritic.com/search/all/" + termo_pesquisa + "/results"
	termo_pesquisa_adaptada := strings.ReplaceAll(termo_pesquisa, " ", "+")
	termo_pesquisa_adaptada2 := strings.ReplaceAll(termo_pesquisa, " ", "%20")
	var pesquisa_meta = "https://www.metacritic.com/search/movie/" + termo_pesquisa_adaptada2 + "/results"
	var pesquisa_rotten = "https://www.rottentomatoes.com/search?search=" + termo_pesquisa_adaptada2

	pesquisa_letter := "https://letterboxd.com/search/films/" + termo_pesquisa_adaptada + "/"

	if err != nil {
		log.Fatalf("Erro ao criar arquivo, err: %q", err)
		return
	}
	defer file.Close()

	fName2 := "data2.csv"
	file2, err := os.Create(fName2)
	if err != nil {
		log.Fatalf("Erro ao criar arquivo, err: %q", err)
		return
	}
	defer file2.Close()

	fName3 := "data3.csv"
	file3, err := os.Create(fName3)
	if err != nil {
		log.Fatalf("Erro ao criar arquivo, err: %q", err)
		return
	}
	defer file3.Close()

	fName4 := "data4.csv"
	file4, err := os.Create(fName4)
	if err != nil {
		log.Fatalf("Erro ao criar arquivo, err: %q", err)
		return
	}
	defer file4.Close()

	fName5 := "data5.csv"
	file5, err := os.Create(fName5)
	if err != nil {
		log.Fatalf("Erro ao criar arquivo, err: %q", err)
		return
	}
	defer file5.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer2 := csv.NewWriter(file2)
	defer writer2.Flush()

	writer3 := csv.NewWriter(file3)
	defer writer3.Flush()

	writer4 := csv.NewWriter(file4)
	defer writer4.Flush()

	writer5 := csv.NewWriter(file5)
	defer writer5.Flush()

	c := colly.NewCollector(
		colly.AllowedDomains("www.metacritic.com"),
	)

	c2 := colly.NewCollector(
		colly.AllowedDomains("www.imdb.com"),
	)

	c3 := colly.NewCollector(
		//colly.AllowedDomains("www.rottentomatoes.com"),
		colly.AllowedDomains("www.metacritic.com"),
	)

	c4 := colly.NewCollector(
		//colly.AllowedDomains("www.rottentomatoes.com"),
		colly.AllowedDomains("www.rottentomatoes.com"),
	)

	c5 := colly.NewCollector(
		//colly.AllowedDomains("www.rottentomatoes.com"),
		colly.AllowedDomains("letterboxd.com"),
	)

	c.OnHTML(".clamp-summary-wrap", func(e *colly.HTMLElement) {
		writer.Write([]string{
			//e.ChildText("a"),
			e.ChildText("h3"),
			e.ChildText("span"),
			//e.ChildText("text"),
		})
	})

	c2.OnHTML(".lister-item-content", func(e *colly.HTMLElement) {
		writer2.Write([]string{
			//e.ChildText("a"),
			e.ChildText("h3"),
			e.ChildText("p"),
			//e.ChildText("div"),
		})
	})

	//c3.OnHTML(".mb_info", func(e *colly.HTMLElement) {
	//	writer3.Write([]string{
	//e.ChildText("a"),

	//		e.ChildText("a"),
	//		e.ChildText("span"),
	//e.ChildText("consensus"),
	//e.ChildText("div"),
	//	})
	//})
	c3.OnHTML(".main_stats", func(e *colly.HTMLElement) {
		writer3.Write([]string{
			e.ChildText(".metascore_w"),
			e.ChildText("a"),
			//fmt.Println(string(e.Text))
			//rows = append(rows, []string{string(titulo), string(nota)})
		})
	})
	// Rotten html
	c4.OnHTML("ul search-page-media-row", func(e *colly.HTMLElement) {
		writer4.Write([]string{
			e.ChildAttr("score-icon-critic", "percentage"),
			e.ChildText("a"),
		})

	})

	c5.OnHTML(".film-detail-content", func(e *colly.HTMLElement) {
		writer5.Write([]string{
			e.ChildText("a href"),
			e.ChildText(".film-title-wrapper"),
		})

	})

	var waitGroup sync.WaitGroup
	waitGroup.Add(11)

	for i := 0; i < 11; i++ {
		go scrapPage(i, &waitGroup, c)

		//go scrapPage2(i, &waitGroup, c2)
		// agora vai ser em paralelo
	}

	waitGroup.Wait()

	var waitGroup2 sync.WaitGroup
	waitGroup2.Add(11)
	for i := 0; i < 11; i++ {

		go scrapPage2(i, &waitGroup2, c2)

	}
	waitGroup2.Wait()

	//var waitGroup3 sync.WaitGroup
	//waitGroup3.Add(11)
	//for i := 0; i < 11; i++ {

	//	go scrapPage2(i, &waitGroup3, c3)

	//}
	//waitGroup3.Wait()
	log.Printf("Scraping " + termo_pesquisa)
	c3.Visit(pesquisa_meta)
	c4.Visit(pesquisa_rotten)
	c5.Visit(pesquisa_letter)

	log.Printf("Scraping Completo")
	log.Println(c)
	log.Println(c2)
	log.Println(c3)
	log.Println(c4)
}
