package gog

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"mtg-price-scrapper-sg/scrapper"
)

const StoreName = "Grey Ogre Games"
const StoreBaseURL = "https://www.greyogregames.com"

type Store struct {
	Name    string
	BaseUrl string
}

func NewScrapper() scrapper.Scrapper {
	return Store{
		Name:    StoreName,
		BaseUrl: StoreBaseURL,
	}
}

func (s Store) Scrap(searchStr string) ([]scrapper.Card, error) {
	searchURL := s.BaseUrl + "/search?q=" + url.QueryEscape(searchStr)
	var cards []scrapper.Card

	c := colly.NewCollector()

	c.OnHTML("div.collectionGrid", func(e *colly.HTMLElement) {
		e.ForEach("div.productCard__card", func(_ int, el *colly.HTMLElement) {
			var (
				isInstock bool
				price     float64
			)

			// in stock
			if len(el.ChildTexts("div.productCard__button--outOfStock")) == 0 {
				isInstock = true
			}

			// price
			var priceStr string

			if strings.TrimSpace(el.ChildText("p.productCard__price")) != "" {
				priceStr = el.ChildText("p.productCard__price")
			} else {
				priceStr = el.ChildText("p.productCard__price")
			}

			priceStr = strings.Replace(priceStr, "$", "", -1)
			priceStr = strings.Replace(priceStr, "SGD", "", -1)
			price, _ = strconv.ParseFloat(strings.TrimSpace(priceStr), 64)

			if price > 0 {
				cards = append(cards, scrapper.Card{
					Name:    strings.TrimSpace(el.ChildText("p.productCard__title")),
					Url:     strings.TrimSpace(s.BaseUrl + el.ChildAttr("a", "href")),
					InStock: isInstock,
					Price:   price,
					Source:  s.Name,
					Img:     strings.TrimSpace("https:" + el.ChildAttr("img", "data-src")),
				})
			}
		})
	})

	return cards, c.Visit(searchURL)
}
