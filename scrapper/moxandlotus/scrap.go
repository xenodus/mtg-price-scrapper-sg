package moxandlotus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"mtg-price-scrapper-sg/scrapper"
)

const StoreName = "Mox & Lotus"
const StoreBaseURL = "https://moxandlotus.sg"
const StoreSearchURL = "/products?title="
const StoreApiURL = "/api/products?&limit=48&full_search=true&showStatus=false&is_paginated=true&in_stock=true&sortVariation=true&&category_id=1&variation_code=regular&order_by=Price%20Low%20to%20High&search="
const CardImageURL = "https://d3nmvyqkci0c2u.cloudfront.net/%s/%s.png"

type apiResponse struct {
	CurrentPage int `json:"current_page"`
	Data        []struct {
		ID                    int       `json:"id"`
		ImagePath             any       `json:"image_path"`
		Title                 string    `json:"title"`
		CardNumber            string    `json:"card_number"`
		SaleTag               any       `json:"sale_tag"`
		CategoryID            int       `json:"category_id"`
		Artist                string    `json:"artist"`
		ExpansionCode         string    `json:"expansion_code"`
		RarityCode            string    `json:"rarity_code"`
		ColorCode             string    `json:"color_code"`
		ColorIdentity         string    `json:"color_identity"`
		ManaCost              string    `json:"mana_cost"`
		ManaValue             int       `json:"mana_value"`
		Weight                any       `json:"weight"`
		VariationCode         string    `json:"variation_code"`
		TypeCode              string    `json:"type_code"`
		CustomText            any       `json:"custom_text"`
		Tag                   any       `json:"tag"`
		IsPopular             int       `json:"is_popular"`
		IsRecommended         int       `json:"is_recommended"`
		Ability               string    `json:"ability"`
		Ruling                any       `json:"ruling"`
		CreatedAt             string    `json:"created_at"`
		UpdatedAt             time.Time `json:"updated_at"`
		Limit                 any       `json:"limit"`
		LockStock             int       `json:"lock_stock"`
		LockPrice             int       `json:"lock_price"`
		NotForSale            any       `json:"not_for_sale"`
		CrawlCode             any       `json:"crawl_code"`
		CrawlSite             any       `json:"crawl_site"`
		CrawlDate             any       `json:"crawl_date"`
		CrawlSource           any       `json:"crawl_source"`
		PreOrder              any       `json:"pre_order"`
		Status                string    `json:"status"`
		Price                 string    `json:"price"`
		TotalSold             string    `json:"total_sold"`
		TotalStocks           string    `json:"totalStocks"`
		Stocks                int       `json:"stocks"`
		DefaultConditionCode  string    `json:"default_condition_code"`
		DefaultConditionIndex int       `json:"default_condition_index"`
		Rarity                string    `json:"rarity"`
		Expansion             string    `json:"expansion"`
		ReviewCount           int       `json:"review_count"`
		Rating                int       `json:"rating"`
		Conditions            []struct {
			ID        int       `json:"id"`
			ProductID int       `json:"product_id"`
			Code      string    `json:"code"`
			Stocks    int       `json:"stocks"`
			Price     string    `json:"price"`
			Sold      int       `json:"sold"`
			LockStock int       `json:"lock_stock"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Limit     int       `json:"limit"`
		} `json:"conditions"`
		Inventories []any `json:"inventories"`
	} `json:"data"`
	FirstPageURL string `json:"first_page_url"`
	From         int    `json:"from"`
	LastPage     int    `json:"last_page"`
	LastPageURL  string `json:"last_page_url"`
	Links        []struct {
		URL    any    `json:"url"`
		Label  string `json:"label"`
		Active bool   `json:"active"`
	} `json:"links"`
	NextPageURL any    `json:"next_page_url"`
	Path        string `json:"path"`
	PerPage     string `json:"per_page"`
	PrevPageURL any    `json:"prev_page_url"`
	To          int    `json:"to"`
	Total       int    `json:"total"`
}

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
		apiResponse apiResponse
		cards       []scrapper.Card
	)

	searchURL := s.BaseUrl + s.SearchUrl + url.QueryEscape(searchStr)
	apiURL := s.BaseUrl + StoreApiURL + url.QueryEscape(searchStr)

	resp, err := http.Get(apiURL)
	if err != nil {
		return cards, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return cards, err
	}

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return cards, err
	}

	if len(apiResponse.Data) > 0 {
		for _, card := range apiResponse.Data {
			price, _ := strconv.ParseFloat(strings.TrimSpace(card.Price), 64)
			cardNo, err := strconv.Atoi(card.CardNumber)
			if err != nil {
				continue
			}

			cards = append(cards, scrapper.Card{
				Name:    strings.TrimSpace(card.Title),
				Url:     searchURL,
				InStock: true,
				Price:   price,
				Source:  s.Name,
				Img:     fmt.Sprintf(CardImageURL, card.ExpansionCode, fmt.Sprintf("%03d", cardNo)),
			})
		}
	}

	return cards, nil
}
