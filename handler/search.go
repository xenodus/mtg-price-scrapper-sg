package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"mtg-price-scrapper-sg/controller"
	"mtg-price-scrapper-sg/pkg/config"
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

type WebResponse struct {
	Data []scrapper.Card `json:"data"`
}

func Search(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var apiRes events.APIGatewayProxyResponse
	var webRes WebResponse
	var lgs []string

	searchString, err := url.QueryUnescape(strings.TrimSpace(request.QueryStringParameters["s"]))
	if err != nil {
		searchString = ""
	}
	lgsString, err := url.QueryUnescape(strings.TrimSpace(request.QueryStringParameters["lgs"]))
	if err != nil {
		lgsString = ""
	}

	if config.IsTestEnv {
		searchString = "Sol Ring"
		lgsString, _ = url.QueryUnescape("Flagship%20Games%2CGames%20Haven%2CGrey%20Ogre%20Games%2CHideout%2CMana%20Pro%2CMox%20%26%20Lotus%2COneMtg%2CSanctuary%20Gaming")
	}

	if searchString == "" || len(searchString) < 3 {
		apiRes.StatusCode = http.StatusBadRequest
		return lambdaApiResponse(apiRes, webRes)
	}

	if lgsString != "" {
		lgs = strings.Split(lgsString, ",")
	}

	inStockCards, _ := controller.Search(searchString, lgs)

	if len(inStockCards) > 0 {
		apiRes.StatusCode = http.StatusOK
		webRes.Data = inStockCards
	}

	return lambdaApiResponse(apiRes, webRes)
}

func lambdaApiResponse(apiResponse events.APIGatewayProxyResponse, webResponse WebResponse) (events.APIGatewayProxyResponse, error) {
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
