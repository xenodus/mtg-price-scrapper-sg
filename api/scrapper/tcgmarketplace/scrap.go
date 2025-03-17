package tcgmarketplace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"mtg-price-scrapper-sg/scrapper"
)

const StoreName = "The TCG Marketplace"
const StoreBaseURL = "https://www.thetcgmarketplace.com/"

const cardLinkAPI = "https://thetcgmarketplace.com:3501/encoder/advancedsearch"
const cardInfoAPI = "https://thetcgmarketplace.com:3501/product/advancedfilter"
const mtgCategoryNo = 3
const accessTokenKey = "TCG_MARKETPLACE_ACCESS_TOKEN"

type cardLinkResponse struct {
	Status int `json:"status"`
	Data   struct {
		Message string `json:"message"`
		Data    []struct {
			Name                  string      `json:"name"`
			Setcode               string      `json:"setcode"`
			Setname               string      `json:"setname"`
			Language              string      `json:"language"`
			CrdFoilType           interface{} `json:"crd_foil_type"`
			Rarity                string      `json:"rarity"`
			Available             interface{} `json:"available"`
			From                  interface{} `json:"from"`
			NonFoilReferencePrice interface{} `json:"non_foil_reference_price"`
			FoilReferencePrice    interface{} `json:"foil_reference_price"`
			URL                   string      `json:"url"`
		} `json:"data"`
	} `json:"data"`
	Meta struct {
		Total int `json:"total"`
	} `json:"meta"`
}

type cardInfoResponse struct {
	Status int `json:"status"`
	Data   struct {
		Message string `json:"message"`
		Data    []struct {
			ID          int         `json:"id"`
			Name        string      `json:"name"`
			Setcode     string      `json:"setcode"`
			Setname     string      `json:"setname"`
			Language    string      `json:"language"`
			Image       string      `json:"image"`
			CrdFoilType interface{} `json:"crd_foil_type"`
			SymbolImage string      `json:"symbol_image"`
			Rarity      string      `json:"rarity"`
			Available   interface{} `json:"available"`
			From        interface{} `json:"from"`
		} `json:"data"`
	} `json:"data"`
	Meta struct {
		Total int `json:"total"`
	} `json:"meta"`
}

type Store struct {
	Name      string
	BaseUrl   string
	SearchUrl string
}

type payload struct {
	AccessToken string `json:"access_token"`
	Name        string `json:"name"`
	Category    int32  `json:"category"`
	Order       string `json:"order"`
}

func NewScrapper() scrapper.Scrapper {
	return Store{
		Name:    StoreName,
		BaseUrl: StoreBaseURL,
	}
}

func (s Store) Scrap(searchStr string) ([]scrapper.Card, error) {
	var (
		cardLinkRes cardLinkResponse
		cardInfoRes cardInfoResponse
		cards       []scrapper.Card
		accessToken string
	)

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	accessToken = os.Getenv(accessTokenKey)

	reqPayload, err := json.Marshal(payload{
		AccessToken: accessToken,
		Name:        searchStr,
		Category:    mtgCategoryNo,
		Order:       "name_asc",
	})
	if err != nil {
		return cards, err
	}

	// todo: check if can get both in 1 api request
	// 1st request to get card link and price
	resp, err := http.Post(cardLinkAPI, "application/json", bytes.NewBuffer(reqPayload))
	if err != nil {
		return cards, err
	}
	defer resp.Body.Close()

	cardLinkBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return cards, err
	}

	err = json.Unmarshal(cardLinkBody, &cardLinkRes)
	if err != nil {
		return cards, err
	}

	// 2nd request to get card img
	resp, err = http.Post(cardInfoAPI, "application/json", bytes.NewBuffer(reqPayload))
	if err != nil {
		return cards, err
	}
	defer resp.Body.Close()

	cardInfoBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return cards, err
	}

	err = json.Unmarshal(cardInfoBody, &cardInfoRes)
	if err != nil {
		return cards, err
	}

	if len(cardLinkRes.Data.Data) > 0 {
		for _, card := range cardLinkRes.Data.Data {
			if card.Available != nil && card.From != nil {
				price, err := strconv.ParseFloat(fmt.Sprint(card.From), 64)
				if err != nil {
					continue
				}

				// Strip [XXX] prefix from card name
				// e.g. [CMM] Deflecting Swat (V2)(Etched foil)
				name := strings.TrimSpace(card.Name)
				squareBracketIndex := strings.Index(name, "]")
				if squareBracketIndex > 1 {
					name = strings.TrimSpace(name[squareBracketIndex+1:])
				}

				cards = append(cards, scrapper.Card{
					Name:    strings.TrimSpace(name),
					Url:     card.URL,
					InStock: true,
					Price:   price,
					Source:  s.Name,
					// attempt to get image from 2nd request
					Img: cardInfoRes.getImageURLByName(card.Name),
				})
			}
		}
	}

	log.Println(cards)

	return cards, nil
}

func (c cardInfoResponse) getImageURLByName(name string) string {
	if len(c.Data.Data) > 0 {
		for _, card := range c.Data.Data {
			if strings.TrimSpace(card.Name) == strings.TrimSpace(name) {
				return card.Image
			}
		}
	}
	return ""
}
