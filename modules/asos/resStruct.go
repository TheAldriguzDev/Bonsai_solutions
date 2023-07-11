package asos

// Search endpoint

type Price struct {
	Current struct {
		Value float64 `json:"value"`
		Text  string  `json:"text"`
	} `json:"current"`
}

type ProductSearch struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       Price  `json:"price"`
	Colour      string `json:"colour"`
	BrandName   string `json:"brandName"`
	ProductCode int    `json:"productCode"`
	URL         string `json:"url"`
	ImageURL    string `json:"imageUrl"`
}

type Data struct {
	SearchTerm   string          `json:"searchTerm"`
	CategoryName *string         `json:"categoryName"`
	ItemCount    int             `json:"itemCount"`
	RedirectURL  string          `json:"redirectUrl"`
	Products     []ProductSearch `json:"products"`
}

// Product endpoint

type LocalisedData struct {
	Locale string `json:"locale"`
	Title  string `json:"title"`
	PdpUrl string `json:"pdpUrl"`
}

type PriceData struct {
	Value        float64 `json:"value"`
	Text         string  `json:"text"`
	VersionId    string  `json:"versionId"`
	ConversionId string  `json:"conversionId"`
}

type Variant struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	SizeID          int    `json:"sizeId"`
	DisplaySizeText string `json:"displaySizeText"`
	SizeOrder       int    `json:"sizeOrder"`
	SKU             string `json:"sku"`
	IsInStock       bool   `json:"isInStock"`
	IsAvailable     bool   `json:"isAvailable"`
	Colour          string `json:"colour"`
	Price           struct {
		Current  PriceData `json:"current"`
		Currency string    `json:"currency"`
	} `json:"price"`
}

type Product struct {
	ID            int             `json:"id"`
	Name          string          `json:"name"`
	LocalisedData []LocalisedData `json:"localisedData"`
	ProductCode   string          `json:"productCode"`
	Variants      []Variant       `json:"variants"`
	Price         struct {
		Current PriceData `json:"current"`
	} `json:"price"`
}
