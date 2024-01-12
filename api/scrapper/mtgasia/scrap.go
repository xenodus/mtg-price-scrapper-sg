package mtgasia

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"mtg-price-scrapper-sg/scrapper"
)

const StoreName = "MTG Asia"
const StoreBaseURL = "http://www.asianmagickards.com/"
const StoreSearchURL = "store.cfm?vSS=%s&iGC=1"

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

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("cookie", "CURRENCY=SGD;")
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		e.ForEach("table#cardTable div.accordion", func(_ int, el *colly.HTMLElement) {
			// only proceed if name found
			if len(el.ChildTexts("span.missingText")) > 0 {
				el.ForEach("table.cardHeader", func(_ int, el2 *colly.HTMLElement) {
					var (
						isInstock bool
						price     float64
						quality   string
					)
					if len(el2.ChildTexts("td")) == 2 {
						priceStr := el2.ChildTexts("td")[0]
						priceStr = strings.Replace(priceStr, "SG$", "", -1)
						priceStr = strings.Replace(priceStr, ",", "", -1)
						price, _ = strconv.ParseFloat(strings.TrimSpace(priceStr), 64)

						if price > 0 {
							isInstock = true
							quality = el2.ChildTexts("td")[1]
						}
					}

					if isInstock {
						if el.ChildText("span.foilText") != "" {
							quality += " " + el.ChildText("span.foilText")
						}

						cards = append(cards, scrapper.Card{
							Name:    strings.TrimSpace(el.ChildTexts("span.missingText")[0]),
							Url:     searchURL,
							InStock: isInstock,
							Price:   price,
							Source:  s.Name,
							// MTG Asia is not on https so can't use their image
							// Img:     strings.TrimSpace(StoreBaseURL + el.ChildAttr("img", "src")),
							Quality: strings.TrimSpace(quality),
						})
					}
				})
			}
		})
	})

	return cards, c.Visit(searchURL)
}
