package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"mtg-price-scrapper-sg/scrapper"
	"mtg-price-scrapper-sg/scrapper/agora"
	"mtg-price-scrapper-sg/scrapper/flagship"
	"mtg-price-scrapper-sg/scrapper/gog"
	"mtg-price-scrapper-sg/scrapper/hideout"
	"mtg-price-scrapper-sg/scrapper/manapro"
	"mtg-price-scrapper-sg/scrapper/onemtg"
	"mtg-price-scrapper-sg/scrapper/sanctuary"
)

func main() {
	var cards []scrapper.Card
	var searchString string

	flag.StringVar(&searchString, "s", "", "Card name to search for")
	flag.Parse()

	// searchString = "Ghalta, Primal Hunger"

	if searchString == "" {
		log.Println("No card name supplied. Supply with s=\"Ghalta, Primal Hunger\"")
		os.Exit(0)
	}

	shopScrapperMap := initAndMapScrappers()

	for shopName, shopScrapper := range shopScrapperMap {
		c, err := shopScrapper.Scrap(searchString)
		if err != nil {
			log.Fatal("err fetching from "+shopName, err)
		}

		if len(c) > 0 {
			cards = append(cards, c...)
		}
	}

	// Output Results
	if len(cards) > 0 {
		// Sort by price ASC
		sort.SliceStable(cards, func(i, j int) bool {
			return cards[i].Price < cards[j].Price
		})

		log.Println("---------------------------------------------------------------")
		i := 1
		for _, c := range cards {
			if c.InStock {
				log.Println(fmt.Sprintf("%v. ", i), c.Name, " | "+fmt.Sprintf("S$ %.2f", c.Price*100/100), " | "+c.Source)
				i++
			}
		}
		log.Println("---------------------------------------------------------------")
	}
}

func initAndMapScrappers() map[string]scrapper.Scrapper {
	return map[string]scrapper.Scrapper{
		agora.StoreName:     agora.NewScrapper(),
		flagship.StoreName:  flagship.NewScrapper(),
		onemtg.StoreName:    onemtg.NewScrapper(),
		manapro.StoreName:   manapro.NewScrapper(),
		gog.StoreName:       gog.NewScrapper(),
		hideout.StoreName:   hideout.NewScrapper(),
		sanctuary.StoreName: sanctuary.NewScrapper(),
	}
}
