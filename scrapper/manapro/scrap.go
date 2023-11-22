package manapro

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"mtg-price-scrapper-sg/scrapper"
)

const StoreName = "ManaPro"
const StoreBaseURL = "https://sg-manapro.com"

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
	searchURL := s.BaseUrl + "/search?type=product&q=" + url.QueryEscape(searchStr)
	var cards []scrapper.Card

	c := colly.NewCollector()

	c.OnHTML("div.products-display", func(e *colly.HTMLElement) {
		e.ForEach("div.product-card-list2", func(_ int, el *colly.HTMLElement) {
			var (
				isInstock bool
				price     float64
			)

			// in stock
			if len(el.ChildTexts("a.addToCart span.value")) > 0 {
				isInstock = el.ChildTexts("a.addToCart span.value")[len(el.ChildTexts("a.addToCart span.value"))-1] != "SOLD OUT"
			}

			// price
			var priceStr string

			if strings.TrimSpace(el.ChildText("span.qv-discountprice")) != "" {
				priceStr = el.ChildText("span.qv-discountprice")
			} else {
				priceStr = el.ChildText("span.qv-regularprice")
			}

			priceStr = strings.Replace(priceStr, "$", "", -1)
			price, _ = strconv.ParseFloat(priceStr, 64)

			if price > 0 {
				cards = append(cards, scrapper.Card{
					Name:    el.ChildText("div.grid-view-item__title"),
					Url:     s.BaseUrl + el.ChildAttr("a", "href"),
					InStock: isInstock,
					Price:   price,
					Source:  s.Name,
					Img:     "https:" + el.ChildAttr("img", "src"),
				})
			}
		})
	})

	return cards, c.Visit(searchURL)
}
