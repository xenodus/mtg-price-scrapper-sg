package duellerpoint

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"mtg-price-scrapper-sg/scrapper"
)

const StoreName = "Dueller's Point"
const StoreBaseURL = "https://www.duellerspoint.com"
const StoreSearchURL = "/products/search?search_text=%s"

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
	var cards []scrapper.Card
	apiURL := s.BaseUrl + fmt.Sprintf(s.SearchUrl, url.QueryEscape(searchStr))

	resp, err := http.Get(apiURL)
	if err != nil {
		return cards, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return cards, err
	}

	doc.Find("div.dp_item").Each(func(i int, se *goquery.Selection) {
		if se.Find("label.label-info").Text() != "" {
			cardTitle := se.Find("a.dp_title").Text()
			cardUrl := se.Find("a.dp_title").AttrOr("href", "")
			cardImg := se.Find("a.dp_title").AttrOr("data-image-url", "")

			cardPriceStr := se.Find("div.dp_price").Text()
			cardPriceStr = strings.Replace(cardPriceStr, "S$", "", -1)
			cardPriceStr = strings.Replace(cardPriceStr, ",", "", -1)
			cardPriceStr = strings.Replace(cardPriceStr, "SGD", "", -1)
			cardPrice, _ := strconv.ParseFloat(strings.TrimSpace(cardPriceStr), 64)

			if cardTitle != "" && cardUrl != "" && cardImg != "" && cardPrice != 0 {
				cards = append(cards, scrapper.Card{
					Name:    strings.TrimSpace(cardTitle),
					Url:     StoreBaseURL + cardUrl,
					InStock: true,
					Price:   cardPrice,
					Source:  s.Name,
					Img:     StoreBaseURL + cardImg,
				})
			}
		}
	})

	return cards, nil
}
