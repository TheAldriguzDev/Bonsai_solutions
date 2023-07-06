package ftlRestock

// General search
type ProductSearch struct {
	Name  string `json:"name"`
	Price struct {
		FormattedValue string `json:"formattedValue"`
	} `json:"price"`
	SKU string `json:"sku"`
}

type SearchData struct {
	Products []ProductSearch `json:"products"`
}

// Product response

type Query struct {
	URL   string `json:"url"`
	Query struct {
		Value string `json:"value"`
	} `json:"query"`
}

type Value struct {
	Name     string `json:"name"`
	Count    int    `json:"count"`
	Selected bool   `json:"selected"`
	Hide     bool   `json:"hide"`
	Query    Query  `json:"query"`
}

type Facet struct {
	Name        string  `json:"name"`
	MultiSelect bool    `json:"multiSelect"`
	Values      []Value `json:"values"`
	Visible     bool    `json:"visible"`
	Category    bool    `json:"category"`
	Priority    int     `json:"priority"`
	Code        string  `json:"code"`
}

type Price struct {
	Value          float64 `json:"value"`
	FormattedValue string  `json:"formattedValue"`
}

type Product struct {
	Name      string `json:"name"`
	LongTitle string `json:"longTitle"`
	Price     Price  `json:"price"`
	Sku       string `json:"sku"`
	ImageSku  string `json:"imageSku"`
}

type ProductData struct {
	Facets         []Facet   `json:"facets"`
	FreeTextSearch string    `json:"freeTextSearch"`
	Products       []Product `json:"products"`
}
