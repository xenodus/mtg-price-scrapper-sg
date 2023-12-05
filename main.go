package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"mtg-price-scrapper-sg/scrapper"
	"mtg-price-scrapper-sg/scrapper/agora"
	"mtg-price-scrapper-sg/scrapper/flagship"
	"mtg-price-scrapper-sg/scrapper/gameshaven"
	"mtg-price-scrapper-sg/scrapper/gog"
	"mtg-price-scrapper-sg/scrapper/hideout"
	"mtg-price-scrapper-sg/scrapper/manapro"
	"mtg-price-scrapper-sg/scrapper/moxandlotus"
	"mtg-price-scrapper-sg/scrapper/onemtg"
	"mtg-price-scrapper-sg/scrapper/sanctuary"
)

const isTestEnv = false

type webResponse struct {
	Data []scrapper.Card `json:"data"`
}

func main() {
	if isTestEnv {
		start := time.Now()
		log.Println(handler(context.Background(), events.APIGatewayProxyRequest{}))
		log.Println(fmt.Sprintf("Took: %s", time.Since(start)))
	} else {
		lambda.Start(handler)
	}
}

func handler(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var cards, inStockCards []scrapper.Card
	var apiRes events.APIGatewayProxyResponse
	var webRes webResponse
	var lgs []string

	searchString, err := url.QueryUnescape(strings.TrimSpace(request.QueryStringParameters["s"]))
	if err != nil {
		searchString = ""
	}
	lgsString, err := url.QueryUnescape(strings.TrimSpace(request.QueryStringParameters["lgs"]))
	if err != nil {
		lgsString = ""
	}

	if isTestEnv {
		searchString = "Sol Ring"
		lgsString, _ = url.QueryUnescape("Flagship%20Games%2CGames%20Haven%2CGrey%20Ogre%20Games%2CHideout%2CManaPro%2CMox%20%26%20Lotus%2COneMtg%2CSanctuary%20Gaming")
	}

	if searchString == "" {
		apiRes.StatusCode = http.StatusBadRequest
		return lambdaApiResponse(apiRes, webRes)
	}

	if lgsString != "" {
		lgs = strings.Split(lgsString, ",")
	}

	shopScrapperMap := initAndMapScrappers(lgs)

	if len(shopScrapperMap) > 0 {
		// Create a channel with a buffer size of numGoroutines
		done := make(chan bool, len(shopScrapperMap))

		log.Println("Start checking shops...")
		for shopName, shopScrapper := range shopScrapperMap {
			shopName := shopName
			shopScrapper := shopScrapper
			go func() {
				start := time.Now()
				c, _ := shopScrapper.Scrap(searchString)
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

		apiRes.StatusCode = http.StatusOK

		if len(cards) > 0 {
			// Sort by price ASC
			sort.SliceStable(cards, func(i, j int) bool {
				return cards[i].Price < cards[j].Price
			})

			// Only showing in stock and not art card
			for _, c := range cards {
				if c.InStock && !strings.Contains(strings.ToLower(c.Name), "art card") {
					inStockCards = append(inStockCards, c)
				}
			}

			if len(inStockCards) > 0 {
				webRes.Data = inStockCards
			}
		}
	}

	return lambdaApiResponse(apiRes, webRes)
}

func lambdaApiResponse(apiResponse events.APIGatewayProxyResponse, webResponse webResponse) (events.APIGatewayProxyResponse, error) {
	bodyBytes, err := json.MarshalIndent(webResponse, "", "    ")
	if err != nil {
		apiResponse.StatusCode = http.StatusInternalServerError
		apiResponse.Body = "err marshalling to json result"
		return apiResponse, nil
	}

	apiResponse.Body = strings.Replace(string(bodyBytes), "\\u0026", "&", -1)

	return apiResponse, nil
}

func initAndMapScrappers(lgs []string) map[string]scrapper.Scrapper {
	storeScrappers := map[string]scrapper.Scrapper{
		agora.StoreName:       agora.NewScrapper(),
		flagship.StoreName:    flagship.NewScrapper(),
		onemtg.StoreName:      onemtg.NewScrapper(),
		manapro.StoreName:     manapro.NewScrapper(),
		gog.StoreName:         gog.NewScrapper(),
		hideout.StoreName:     hideout.NewScrapper(),
		sanctuary.StoreName:   sanctuary.NewScrapper(),
		gameshaven.StoreName:  gameshaven.NewScrapper(),
		moxandlotus.StoreName: moxandlotus.NewScrapper(),
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
