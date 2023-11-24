package onemtg

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"mtg-price-scrapper-sg/scrapper"
)

const StoreName = "OneMtg"
const StoreBaseURL = "https://onemtg.com.sg"
const StoreSearchURL = "/search?q=*%s*"

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
	searchURL := s.BaseUrl + fmt.Sprintf(s.SearchUrl, url.QueryEscape(searchStr))
	var cards []scrapper.Card

	c := colly.NewCollector()

	c.OnHTML("div.container", func(e *colly.HTMLElement) {
		e.ForEach("div.Mob", func(_ int, el *colly.HTMLElement) {
			var (
				isInstock bool
				price     float64
			)

			// in stock + price
			// only have in stock data
			if len(el.ChildTexts("div.addNow")) > 0 {
				isInstock = el.ChildTexts("div.addNow")[len(el.ChildTexts("div.addNow"))-1] != ""
				if isInstock {
					priceStr := strings.TrimSpace(el.ChildTexts("div.addNow")[len(el.ChildTexts("div.addNow"))-1])
					price, _ = parsePriceRegex(priceStr)
				}
			}

			cards = append(cards, scrapper.Card{
				Name:    strings.TrimSpace(el.ChildText("p.productTitle")),
				Url:     strings.TrimSpace(s.BaseUrl + el.ChildAttr("a", "href")),
				InStock: isInstock,
				Price:   price,
				Source:  s.Name,
				Img:     strings.TrimSpace("https:" + el.ChildAttr("img", "src")),
			})
		})
	})

	return cards, c.Visit(searchURL)
}

func parsePriceRegex(price string) (float64, error) {
	re := regexp.MustCompile(`(?s)\((.*)\)`)
	m := re.FindAllStringSubmatch(price, -1)
	if len(m) > 0 && len(m[0]) > 1 && len(m[0][1]) > 0 {
		m[0][1] = strings.Replace(m[0][1], "$", "", -1)
		return strconv.ParseFloat(m[0][1], 64)
	}
	return 0, nil
}
