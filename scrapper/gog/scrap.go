package gog

import (
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"mtg-price-scrapper-sg/scrapper"
)

const StoreName = "Grey Ogre Games"
const StoreBaseURL = "https://www.greyogregames.com"
const StoreSearchURL = "/search?q="

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
	var (
		err   error
		cards []scrapper.Card
	)

	pagesMap := make(map[int]string)
	searchURL := s.BaseUrl + s.SearchUrl + url.QueryEscape(searchStr)

	c := colly.NewCollector()

	c.OnHTML("body", func(e *colly.HTMLElement) {
		// get pagination
		e.ForEach("ol.pagination li", func(_ int, el *colly.HTMLElement) {
			elStr := strings.Replace(el.Text, "«", "", -1)
			elStr = strings.Replace(elStr, "page", "", -1)
			elStr = strings.Replace(elStr, "Next", "", -1)
			elStr = strings.Replace(elStr, "Previous", "", -1)
			elStr = strings.Replace(elStr, "»", "", -1)
			elStr = strings.TrimSpace(elStr)
			if elStr != "" && elStr != "1" && el.ChildAttr("a", "href") != "" {
				elInt, strConvErr := strconv.Atoi(elStr)
				if strConvErr == nil {
					pagesMap[elInt] = el.ChildAttr("a", "href")
				}
			}
		})

		// get cards
		e.ForEach("div.productCard__card", func(_ int, el *colly.HTMLElement) {
			var (
				isInstock bool
				price     float64
			)

			// in stock
			if len(el.ChildTexts("div.productCard__button--outOfStock")) == 0 {
				isInstock = true
			}

			if isInstock {
				el.ForEach("ul.productChip__grid li", func(_ int, el2 *colly.HTMLElement) {
					if el2.Attr("data-variantavailable") == "true" && el2.Attr("data-variantqty") != "0" {
						priceStr := el2.Attr("data-variantprice")
						priceStr = strings.Replace(priceStr, "$", "", -1)
						priceStr = strings.Replace(priceStr, ",", "", -1)
						priceStr = strings.Replace(priceStr, "SGD", "", -1)
						price, _ = strconv.ParseFloat(strings.TrimSpace(priceStr), 64)
						price = price / 100

						if price > 0 {
							cards = append(cards, scrapper.Card{
								Name:    strings.TrimSpace(el.ChildText("p.productCard__title")),
								Url:     strings.TrimSpace(s.BaseUrl + el.ChildAttr("a", "href")),
								InStock: isInstock,
								Price:   price,
								Source:  s.Name,
								Img:     strings.TrimSpace("https:" + el.ChildAttr("img", "data-src")),
								Quality: el2.Attr("data-varianttitle"),
							})
						}
					}
				})
			}
		})
	})

	err = c.Visit(searchURL)
	if err != nil {
		return []scrapper.Card{}, err
	}

	if len(pagesMap) > 0 {
		log.Println("Pagination exists for "+s.Name+": ", len(pagesMap))

		c2 := colly.NewCollector(
			colly.Async(true),
		)

		for _, url := range pagesMap {
			searchURL = s.BaseUrl + url

			c2.OnHTML("div.collectionGrid", func(e *colly.HTMLElement) {
				e.ForEach("div.productCard__card", func(_ int, el *colly.HTMLElement) {
					var (
						isInstock bool
						price     float64
					)

					// in stock
					if len(el.ChildTexts("div.productCard__button--outOfStock")) == 0 {
						isInstock = true
					}

					if isInstock {
						el.ForEach("ul.productChip__grid li", func(_ int, el2 *colly.HTMLElement) {
							if el2.Attr("data-variantavailable") == "true" && el2.Attr("data-variantqty") != "0" {
								priceStr := el2.Attr("data-variantprice")
								priceStr = strings.Replace(priceStr, "$", "", -1)
								priceStr = strings.Replace(priceStr, ",", "", -1)
								priceStr = strings.Replace(priceStr, "SGD", "", -1)
								price, _ = strconv.ParseFloat(strings.TrimSpace(priceStr), 64)
								price = price / 100

								if price > 0 {
									cards = append(cards, scrapper.Card{
										Name:    strings.TrimSpace(el.ChildText("p.productCard__title")),
										Url:     strings.TrimSpace(s.BaseUrl + el.ChildAttr("a", "href")),
										InStock: isInstock,
										Price:   price,
										Source:  s.Name,
										Img:     strings.TrimSpace("https:" + el.ChildAttr("img", "data-src")),
										Quality: el2.Attr("data-varianttitle"),
									})
								}
							}
						})
					}
				})
			})

			err = c2.Visit(searchURL)
			if err != nil {
				break
			}
		}
		c2.Wait()
	}

	return cards, err
}
