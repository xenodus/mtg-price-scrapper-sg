package onemtg

import (
	"fmt"
	"net/url"
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
		e.ForEach("div.Norm", func(_ int, el *colly.HTMLElement) {
			var isInstock bool

			if len(el.ChildTexts("div.addNow")) > 0 {
				for i := 0; i < len(el.ChildTexts("div.addNow")); i++ {
					isInstock = el.ChildTexts("div.addNow")[i] != ""

					if isInstock {
						priceStr := strings.TrimSpace(el.ChildTexts("div.addNow")[i])

						price, quality, err := parsePriceAndQuality(priceStr)
						if err != nil {
							continue
						}

						if price > 0 {
							cards = append(cards, scrapper.Card{
								Name:    strings.TrimSpace(el.ChildText("p.productTitle")),
								Url:     strings.TrimSpace(s.BaseUrl + el.ChildAttr("a", "href")),
								InStock: isInstock,
								Price:   price,
								Source:  s.Name,
								Img:     strings.TrimSpace("https:" + el.ChildAttr("img", "src")),
								Quality: quality,
							})
						}
					}
				}
			}
		})
	})

	return cards, c.Visit(searchURL)
}

func parsePriceAndQuality(priceQualityStr string) (float64, string, error) {
	priceQualityStrSlice := strings.Split(priceQualityStr, " - ")
	if len(priceQualityStrSlice) == 2 {
		quality := strings.TrimSpace(priceQualityStrSlice[0])

		priceStr := strings.TrimSpace(priceQualityStrSlice[1])
		priceStr = strings.Replace(priceStr, "$", "", -1)
		priceStr = strings.Replace(priceStr, ",", "", -1)
		price, err := strconv.ParseFloat(priceStr, 64)

		return price, quality, err
	}
	return 0, "", nil
}
