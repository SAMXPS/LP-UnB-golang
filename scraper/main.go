package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"github.com/gocolly/colly"
)

func scrapPage(i int, wg *sync.WaitGroup, c *colly.Collector) {
	fmt.Printf("Scraping Page: %d\n", i)
	c.Visit("https://www.metacritic.com/browse/movies/score/metascore/all/filtered/netflix?page=" + strconv.Itoa(i))
	wg.Done()
}

func main() {
	fName := "data.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Erro ao criar arquivo, err: %q", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	c := colly.NewCollector(
		colly.AllowedDomains("www.metacritic.com"),
	)

	c.OnHTML(".clamp-summary-wrap", func(e *colly.HTMLElement) {
		writer.Write([]string{
			//e.ChildText("a"),
			e.ChildText("h3"),
			e.ChildText("span"),
			e.ChildText("summary"),
		})
	})

	var waitGroup sync.WaitGroup
	waitGroup.Add(11)

	for i := 0; i < 11; i++ {
		go scrapPage(i, &waitGroup, c)
		// agora vai ser em paralelo
	}

	waitGroup.Wait()

	log.Printf("Scraping Completo")
	log.Println(c)
}
