package scrapper

type Card struct {
	Name    string
	Url     string
	Img     string
	Price   float64
	InStock bool
	Source  string
}

type Scrapper interface {
	Scrap(searchStr string) ([]Card, error)
}
