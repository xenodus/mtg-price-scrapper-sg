package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/mux"
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
	Status string          `json:"status"`
	Code   int             `json:"code"`
	Data   []scrapper.Card `json:"data"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", searchCards).Methods("GET")
	http.ListenAndServe(":8080", r)
}

func searchCards(w http.ResponseWriter, r *http.Request) {
	var cards, inStockCards []scrapper.Card

	res := webResponse{
		Status: "Err",
		Code:   http.StatusBadRequest,
	}

	searchString := strings.TrimSpace(r.URL.Query().Get("s"))

	if searchString == "" {
		res.Status = "Error: no search string provided"
		returnWebResponse(w, res)
		return
	}

	if searchString != "" {
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

		res.Status = "OK"
		res.Code = http.StatusOK

		if len(cards) > 0 {
			// Sort by price ASC
			sort.SliceStable(cards, func(i, j int) bool {
				return cards[i].Price < cards[j].Price
			})

			for _, c := range cards {
				if c.InStock && !strings.Contains(strings.ToLower(c.Name), "art card") {
					inStockCards = append(inStockCards, c)
				}
			}

			if len(inStockCards) > 0 {
				res.Data = inStockCards
				returnWebResponse(w, res)
				return
			}
		}
	}

	returnWebResponse(w, res)
	return
}

func returnWebResponse(w http.ResponseWriter, res webResponse) {
	responseBytes, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		log.Println("err marshalling to json result", err)

		res.Status = "Error: encoding json"
		res.Code = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	w.Write(responseBytes)
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
