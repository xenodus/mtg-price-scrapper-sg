package scrapper

type Card struct {
	Name    string  `json:"name"`
	Url     string  `json:"url"`
	Img     string  `json:"img"`
	Price   float64 `json:"price"`
	InStock bool    `json:"inStock"`
	Source  string  `json:"src"`
}

type Scrapper interface {
	Scrap(searchStr string) ([]Card, error)
}
