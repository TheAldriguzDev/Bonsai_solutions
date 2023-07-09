package size

// Search endpoint

type SearchRes struct {
	Products []Product `json:"products"`
}

type Product struct {
	ID                string      `json:"id,omitempty"`
	OriginalID        string      `json:"ID,omitempty"`
	SKU               string      `json:"SKU,omitempty"`
	Href              string      `json:"href,omitempty"`
	TrackingSKU       string      `json:"trackingSKU,omitempty"`
	Name              string      `json:"name,omitempty"`
	BrandName         string      `json:"brandName,omitempty"`
	MainImage         string      `json:"mainImage,omitempty"`
	Price             *Price      `json:"price,omitempty"`
	StockStatus       string      `json:"stockStatus,omitempty"`
	Size              interface{} `json:"size,omitempty"`
	Colour            string      `json:"colour,omitempty"`
	Weight            int         `json:"weight,omitempty"`
	ColourDescription string      `json:"colourDescription,omitempty"`
	SortOrder         int         `json:"sortOrder,omitempty"`
	URL               string      `json:"URL,omitempty"`
	ProductReference  string      `json:"productReference,omitempty"`
	ProductName       string      `json:"productName,omitempty"`
}

// Price represents the price details
type Price struct {
	Amount   string `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

// Product ep

type ProductRes struct {
	SKU         string `json:"SKU"`
	StockStatus string `json:"stockStatus"`
}

type SortedOption struct {
	Name    string     `json:"name"`
	Product ProductRes `json:"product"`
}

type ItemData struct {
	SKU           string         `json:"SKU"`
	Href          string         `json:"href"`
	Name          string         `json:"name"`
	BrandName     string         `json:"brandName"`
	MainImage     string         `json:"mainImage"`
	Price         Price          `json:"price"`
	Colour        string         `json:"colour"`
	SortedOptions []SortedOption `json:"sortedOptions"`
	URL           string         `json:"URL"`
}
