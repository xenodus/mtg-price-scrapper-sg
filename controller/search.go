package controller

import (
	"fmt"
	"log"
	"slices"
	"sort"
	"strings"
	"time"

	"mtg-price-scrapper-sg/scrapper"
	"mtg-price-scrapper-sg/scrapper/agora"
	"mtg-price-scrapper-sg/scrapper/cardscitadel"
	"mtg-price-scrapper-sg/scrapper/duellerpoint"
	"mtg-price-scrapper-sg/scrapper/flagship"
	"mtg-price-scrapper-sg/scrapper/gameshaven"
	"mtg-price-scrapper-sg/scrapper/gog"
	"mtg-price-scrapper-sg/scrapper/hideout"
	"mtg-price-scrapper-sg/scrapper/manapro"
	"mtg-price-scrapper-sg/scrapper/moxandlotus"
	"mtg-price-scrapper-sg/scrapper/mtgasia"
	"mtg-price-scrapper-sg/scrapper/onemtg"
	"mtg-price-scrapper-sg/scrapper/sanctuary"
)

type SearchInput struct {
	SearchString string
	Lgs          []string
}

func Search(input SearchInput) ([]scrapper.Card, error) {
	var cards, inStockCards, inStockExactMatchCards, inStockPartialMatchCards, inStockPrefixMatchCards []scrapper.Card

	shopScrapperMap := initAndMapScrappers(input.Lgs)

	if len(shopScrapperMap) > 0 {
		// Create a channel with a buffer size of shopScrapperMap
		done := make(chan bool, len(shopScrapperMap))

		log.Println("Start checking shops...")
		for shopName, shopScrapper := range shopScrapperMap {
			shopName := shopName
			shopScrapper := shopScrapper
			go func() {
				start := time.Now()
				c, _ := shopScrapper.Scrap(input.SearchString)
				log.Println(fmt.Sprintf("Done: %s. Took: %s", shopName, time.Since(start)))

				if len(c) > 0 {
					cards = append(cards, c...)
				}

				// Signal that the goroutine is done
				done <- true
			}()
		}

		// Wait for all goroutines to finish
		for i := 0; i < len(shopScrapperMap); i++ {
			<-done
		}
		log.Println("End checking shops...")

		if len(cards) > 0 {
			// Sort by price ASC
			sort.SliceStable(cards, func(i, j int) bool {
				return cards[i].Price < cards[j].Price
			})

			// Only showing in stock, contains searched string and not art card
			for _, c := range cards {
				if c.InStock && !strings.Contains(strings.ToLower(c.Name), "art card") {
					cleanCardName := c.Name

					// if string has [, get index of it to strip [*] away
					squareBracketIndex := strings.Index(cleanCardName, "[")
					if squareBracketIndex > 1 {
						cleanCardName = strings.TrimSpace(cleanCardName[:squareBracketIndex-1])
					}

					// if string has (, get index of it to strip (*) away
					roundBracketIndex := strings.Index(cleanCardName, "(")
					if roundBracketIndex > 1 {
						cleanCardName = strings.TrimSpace(cleanCardName[:roundBracketIndex-1])
					}

					// todo: make this configurable
					// increase accuracy by only including cards which contains searched string in names
					/*
						if !strings.Contains(strings.ToLower(cleanCardName), strings.ToLower(input.SearchString)) {
							continue
						}
					*/

					// exact match
					if strings.ToLower(cleanCardName) == strings.ToLower(input.SearchString) {
						inStockExactMatchCards = append(inStockExactMatchCards, c)
						continue
					}
					// fall back check for exact card name
					cardNameSlice := strings.Split(cleanCardName, " ")
					if len(cardNameSlice) > 1 {
						if strings.ToLower(cardNameSlice[0]) == strings.ToLower(input.SearchString) {
							inStockExactMatchCards = append(inStockExactMatchCards, c)
							continue
						}
					}

					if strings.HasPrefix(strings.ToLower(cleanCardName), strings.ToLower(input.SearchString)) {
						inStockPrefixMatchCards = append(inStockPrefixMatchCards, c)
						continue
					}

					inStockPartialMatchCards = append(inStockPartialMatchCards, c)
				}
			}

			// order of results: exact > prefix > partial match
			inStockCards = append(inStockExactMatchCards, inStockPrefixMatchCards...)
			inStockCards = append(inStockCards, inStockPartialMatchCards...)

		}
	}
	return inStockCards, nil
}

func initAndMapScrappers(lgs []string) map[string]scrapper.Scrapper {
	storeScrappers := map[string]scrapper.Scrapper{
		agora.StoreName:        agora.NewScrapper(),
		cardscitadel.StoreName: cardscitadel.NewScrapper(),
		duellerpoint.StoreName: duellerpoint.NewScrapper(),
		flagship.StoreName:     flagship.NewScrapper(),
		gameshaven.StoreName:   gameshaven.NewScrapper(),
		gog.StoreName:          gog.NewScrapper(),
		hideout.StoreName:      hideout.NewScrapper(),
		manapro.StoreName:      manapro.NewScrapper(),
		moxandlotus.StoreName:  moxandlotus.NewScrapper(),
		mtgasia.StoreName:      mtgasia.NewScrapper(),
		onemtg.StoreName:       onemtg.NewScrapper(),
		sanctuary.StoreName:    sanctuary.NewScrapper(),
	}

	if len(lgs) > 0 {
		for storeName := range storeScrappers {
			if !slices.Contains(lgs, storeName) {
				delete(storeScrappers, storeName)
			}
		}
	}
	return storeScrappers
}
