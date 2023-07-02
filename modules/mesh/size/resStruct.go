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
	ID                string         `json:"id,omitempty"`
	SKU               string         `json:"SKU,omitempty"`
	Href              string         `json:"href,omitempty"`
	TrackingSKU       string         `json:"trackingSKU,omitempty"`
	Name              string         `json:"name,omitempty"`
	BrandName         string         `json:"brandName,omitempty"`
	MainImage         string         `json:"mainImage,omitempty"`
	Price             Price          `json:"price,omitempty"`
	StockStatus       string         `json:"stockStatus,omitempty"`
	Colour            string         `json:"colour,omitempty"`
	ColourDescription string         `json:"colourDescription,omitempty"`
	Images            []string       `json:"images,omitempty"`
	SortedOptions     []SortedOption `json:"sortedOptions,omitempty"`
	URL               string         `json:"URL,omitempty"`
	ProductName       string         `json:"productName,omitempty"`
}

type SortedOption struct {
	Name    string         `json:"name,omitempty"`
	Product ProductDetails `json:"product,omitempty"`
}

type ProductDetails struct {
	SKU               string `json:"SKU,omitempty"`
	StockStatus       string `json:"stockStatus,omitempty"`
	Colour            string `json:"colour,omitempty"`
	ColourDescription string `json:"colourDescription,omitempty"`
	UPC               string `json:"UPC,omitempty"`
	Size              string `json:"size,omitempty"`
}
