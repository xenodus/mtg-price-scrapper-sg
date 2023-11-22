package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"mtg-price-scrapper-sg/scrapper"
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

	// Create an HTTP server with Gorilla Mux router.
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Start the server in a separate Goroutine.
	go func() {
		log.Println("Starting the server on :8080")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Implement graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down the server...")

	// Set a timeout for shutdown (for example, 5 seconds).
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server gracefully stopped")
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

			// Only showing in stock and not art card
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
		// agora.StoreName:     agora.NewScrapper(),
		flagship.StoreName:  flagship.NewScrapper(),
		onemtg.StoreName:    onemtg.NewScrapper(),
		manapro.StoreName:   manapro.NewScrapper(),
		gog.StoreName:       gog.NewScrapper(),
		hideout.StoreName:   hideout.NewScrapper(),
		sanctuary.StoreName: sanctuary.NewScrapper(),
	}
}
