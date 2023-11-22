package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"mtg-price-scrapper-sg/scrapper"
	"mtg-price-scrapper-sg/scrapper/agora"
	"mtg-price-scrapper-sg/scrapper/flagship"
	"mtg-price-scrapper-sg/scrapper/gog"
	"mtg-price-scrapper-sg/scrapper/hideout"
	"mtg-price-scrapper-sg/scrapper/manapro"
	"mtg-price-scrapper-sg/scrapper/onemtg"
	"mtg-price-scrapper-sg/scrapper/sanctuary"
)

type webResponse struct {
	Data []scrapper.Card `json:"data"`
}

func handler(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var cards, inStockCards []scrapper.Card
	var apiRes events.APIGatewayProxyResponse
	var webRes webResponse

	searchString := strings.TrimSpace(request.QueryStringParameters["s"])

	if searchString == "" {
		apiRes.StatusCode = http.StatusBadRequest
		return lambdaApiResponse(apiRes, webRes)
	}

	log.Println("Searching for card: " + searchString)
	shopScrapperMap := initAndMapScrappers()

	log.Println("Start checking shops...")
	for shopName, shopScrapper := range shopScrapperMap {
		start := time.Now()
		c, _ := shopScrapper.Scrap(searchString)
		log.Println(fmt.Sprintf("Done: %s. Took: %s", shopName, time.Since(start)))

		if len(c) > 0 {
			cards = append(cards, c...)
		}
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
			return lambdaApiResponse(apiRes, webRes)
		}
	}

	return apiRes, nil
}

func main() {
	lambda.Start(handler)
}

func lambdaApiResponse(apiResponse events.APIGatewayProxyResponse, webResponse webResponse) (events.APIGatewayProxyResponse, error) {
	bodyBytes, err := json.MarshalIndent(webResponse, "", "    ")
	if err != nil {
		apiResponse.StatusCode = http.StatusInternalServerError
		apiResponse.Body = "err marshalling to json result"
		return apiResponse, nil
	}

	apiResponse.Body = string(bodyBytes)

	return apiResponse, nil
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
