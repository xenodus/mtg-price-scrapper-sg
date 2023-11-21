package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

const (
	FlagshipBaseURL        = "https://www.flagshipgames.sg"
	OneMtgBaseURL          = "https://onemtg.com.sg"
	ManaProBaseURL         = "https://sg-manapro.com"
	GreyOgreGamesBaseURL   = "https://www.greyogregames.com"
	HideoutBaseURL         = "https://www.hideout-online.com/"
	AgoraHobbyBaseURL      = "https://agorahobby.com/"
	SanctuaryGamingBaseURL = "https://sanctuary-gaming.com/"
)

type card struct {
	name    string
	url     string
	price   float64
	inStock bool
	source  string
}

func main() {
	var cards []card
	var searchString string

	flag.StringVar(&searchString, "s", "", "Card name to search for")
	flag.Parse()

	// searchString = "Ghalta, Primal Hunger"

	if searchString == "" {
		log.Println("No card name supplied. Supply with s=\"Ghalta, Primal Hunger\"")
		os.Exit(0)
	}

	// map shop to their web scrap func
	shopMap := map[string]func(string) ([]card, error){
		AgoraHobbyBaseURL:      agora,
		FlagshipBaseURL:        flagship,
		OneMtgBaseURL:          oneMtg,
		ManaProBaseURL:         manapro,
		GreyOgreGamesBaseURL:   gog,
		HideoutBaseURL:         hideout,
		SanctuaryGamingBaseURL: sanctuaryGaming,
	}

	for shopUrl, fn := range shopMap {
		c, err := fn(searchString)
		if err != nil {
			log.Fatal("err fetching from "+shopUrl, err)
		}

		if len(c) > 0 {
			cards = append(cards, c...)
		}
	}

	// Output Results
	if len(cards) > 0 {
		// Sort by price ASC
		sort.SliceStable(cards, func(i, j int) bool {
			return cards[i].price < cards[j].price
		})

		log.Println("---------------------------------------------------------------")
		for _, c := range cards {
			if c.inStock {
				log.Println(c.name, " | ", fmt.Sprintf("S$ %.2f", c.price*100/100), " | ", c.source+c.url)
			}
		}
		log.Println("---------------------------------------------------------------")
	}
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

func oneMtg(searchStr string) ([]card, error) {
	searchURL := OneMtgBaseURL + "/search?q=*" + url.QueryEscape(searchStr) + "*"
	var cards []card

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

			cards = append(cards, card{
				name:    el.ChildText("p.productTitle"),
				url:     el.ChildAttr("a", "href"),
				inStock: isInstock,
				price:   price,
				source:  OneMtgBaseURL,
			})
		})
	})

	return cards, c.Visit(searchURL)
}

func flagship(searchStr string) ([]card, error) {
	searchURL := FlagshipBaseURL + "/search?type=product&q=" + url.QueryEscape(searchStr)
	var cards []card

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
				cards = append(cards, card{
					name:    el.ChildText("div.grid-view-item__title"),
					url:     el.ChildAttr("a", "href"),
					inStock: isInstock,
					price:   price,
					source:  FlagshipBaseURL,
				})
			}
		})
	})

	return cards, c.Visit(searchURL)
}

func manapro(searchStr string) ([]card, error) {
	searchURL := ManaProBaseURL + "/search?type=product&q=" + url.QueryEscape(searchStr)
	var cards []card

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
				cards = append(cards, card{
					name:    el.ChildText("div.grid-view-item__title"),
					url:     el.ChildAttr("a", "href"),
					inStock: isInstock,
					price:   price,
					source:  GreyOgreGamesBaseURL,
				})
			}
		})
	})

	return cards, c.Visit(searchURL)
}

func sanctuaryGaming(searchStr string) ([]card, error) {
	searchURL := SanctuaryGamingBaseURL + "/search?type=product&q=" + url.QueryEscape(searchStr)
	var cards []card

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
				cards = append(cards, card{
					name:    el.ChildText("p.productTitle"),
					url:     el.ChildAttr("a", "href"),
					inStock: isInstock,
					price:   price,
					source:  SanctuaryGamingBaseURL,
				})
			}
		})
	})

	return cards, c.Visit(searchURL)
}

func gog(searchStr string) ([]card, error) {
	searchURL := GreyOgreGamesBaseURL + "/search?q=" + url.QueryEscape(searchStr)
	var cards []card

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
				cards = append(cards, card{
					name:    el.ChildText("p.productCard__title"),
					url:     el.ChildAttr("a", "href"),
					inStock: isInstock,
					price:   price,
					source:  GreyOgreGamesBaseURL,
				})
			}
		})
	})

	return cards, c.Visit(searchURL)
}

func hideout(searchStr string) ([]card, error) {
	searchURL := HideoutBaseURL + "/search?type=product&q=" + url.QueryEscape(searchStr)
	var cards []card

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
				cards = append(cards, card{
					name:    el.ChildText("div.grid-view-item__title"),
					url:     el.ChildAttr("a", "href"),
					inStock: isInstock,
					price:   price,
					source:  HideoutBaseURL,
				})
			}
		})
	})

	return cards, c.Visit(searchURL)
}

func agora(searchStr string) ([]card, error) {
	searchURL := AgoraHobbyBaseURL + "/store/search?category=mtg&searchfield=" + url.QueryEscape(searchStr)
	var cards []card

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
			price, _ = strconv.ParseFloat(priceStr, 64)

			// name
			name := el.ChildText("div.store-item-title")

			// Exclude Japanese cards
			if price > 0 && !strings.Contains(name, "Japanese") {
				cards = append(cards, card{
					name:    el.ChildText("div.store-item-title"),
					url:     "/store/search?category=mtg&searchfield=" + url.QueryEscape(searchStr),
					inStock: isInstock,
					price:   price,
					source:  AgoraHobbyBaseURL,
				})
			}
		})
	})

	return cards, c.Visit(searchURL)
}
