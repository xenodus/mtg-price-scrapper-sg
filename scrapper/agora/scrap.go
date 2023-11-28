package agora

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"mtg-price-scrapper-sg/scrapper"
)

const StoreName = "Agora Hobby"
const StoreBaseURL = "https://agorahobby.com"
const StoreSearchURL = "/store/search?category=mtg&searchfield="

type Store struct {
	Name      string
	BaseUrl   string
	SearchUrl string
}

func NewScrapper() scrapper.Scrapper {
	return Store{
		Name:      StoreName,
		BaseUrl:   StoreBaseURL,
		SearchUrl: StoreSearchURL,
	}
}

func (s Store) Scrap(searchStr string) ([]scrapper.Card, error) {
	searchURL := s.BaseUrl + s.SearchUrl + url.QueryEscape(searchStr)
	var cards []scrapper.Card

	c := colly.NewCollector()

	c.OnHTML("div#store_listingcontainer", func(e *colly.HTMLElement) {
		e.ForEach("div.store-item", func(_ int, el *colly.HTMLElement) {
			var (
				isInstock bool
				price     float64
			)

			// in stock
			if el.ChildText("div.store-item-stock") != "Stock: 0" {
				isInstock = true
			}

			// price
			priceStr := strings.TrimSpace(el.ChildText("div.store-item-price"))
			priceStr = strings.Replace(priceStr, "$", "", -1)
			priceStr = strings.Replace(priceStr, ",", "", -1)
			price, _ = strconv.ParseFloat(strings.TrimSpace(priceStr), 64)

			// name
			name := el.ChildText("div.store-item-title")

			// Exclude Japanese cards
			if price > 0 && !strings.Contains(name, "Japanese") {
				cards = append(cards, scrapper.Card{
					Name:    strings.TrimSpace(el.ChildText("div.store-item-title")),
					Url:     strings.TrimSpace(s.BaseUrl + "/store/search?category=mtg&searchfield=" + url.QueryEscape(searchStr)),
					InStock: isInstock,
					Price:   price,
					Source:  s.Name,
					Img:     strings.TrimSpace(el.ChildAttr("div.store-item-img", "data-img")),
				})
			}
		})
	})

	return cards, c.Visit(searchURL)
}
