package sanctuary

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"mtg-price-scrapper-sg/scrapper"
)

const StoreName = "Sanctuary Gaming"
const StoreBaseURL = "https://sanctuary-gaming.com"

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

	c.OnHTML("div.container", func(e *colly.HTMLElement) {
		e.ForEach("div.Mob", func(_ int, el *colly.HTMLElement) {
			var (
				isInstock bool
				price     float64
			)

			if len(el.ChildTexts("div.addNow")) > 0 {
				priceStr := strings.TrimSpace(el.ChildTexts("div.addNow")[0])
				price, _ = parsePriceRegex(priceStr)
				if price > 0 {
					isInstock = true
				}
			}

			if price > 0 {
				cards = append(cards, scrapper.Card{
					Name:    el.ChildText("p.productTitle"),
					Url:     s.BaseUrl + el.ChildAttr("a", "href"),
					InStock: isInstock,
					Price:   price,
					Source:  s.Name,
				})
			}
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
