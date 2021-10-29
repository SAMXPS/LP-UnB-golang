package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gocolly/colly"
)

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
	for i := 0; i < 11; i++ {
		fmt.Printf("Scraping Page: %d\n", i)

		c.Visit("https://www.metacritic.com/browse/movies/score/metascore/all/filtered/netflix?page=" + strconv.Itoa(i))
	}
	log.Printf("Scraping Completo")
	log.Println(c)
}
