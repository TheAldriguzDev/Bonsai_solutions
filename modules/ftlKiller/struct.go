package ftlkiller

import tls_client "github.com/bogdanfinn/tls-client"

type footlockerSession struct {
	httpClient   tls_client.HttpClient
	Store_Region string
	ProxyArr     []string `json:"proxy"`
	ProxyFile    string
	Host         string
	UserAgent    string
	DataDome     string
	Cid          string
	Proxy        string
}

type Query struct {
	URL   string `json:"url"`
	Query struct {
		Value string `json:"value"`
	} `json:"query"`
	Sort string `json:"sort"`
}

type FacetValue struct {
	Name     string `json:"name"`
	Count    int    `json:"count"`
	Selected bool   `json:"selected"`
	Hide     bool   `json:"hide"`
	Query    Query  `json:"query"`
}

type Facet struct {
	Name        string       `json:"name"`
	MultiSelect bool         `json:"multiSelect"`
	Values      []FacetValue `json:"values"`
	Visible     bool         `json:"visible"`
	Category    bool         `json:"category"`
	Priority    int          `json:"priority"`
	Code        string       `json:"code"`
}

type Pagination struct {
	CurrentPage  int    `json:"currentPage"`
	PageSize     int    `json:"pageSize"`
	Sort         string `json:"sort"`
	TotalPages   int    `json:"totalPages"`
	TotalResults int    `json:"totalResults"`
}

type Badge struct {
	IsFLOnly     bool `json:"isFLOnly"`
	IsNewProduct bool `json:"isNewProduct"`
	IsOnlineOnly bool `json:"isOnlineOnly"`
	IsSale       bool `json:"isSale"`
}

type SelectedOption struct {
	MapEnable bool   `json:"mapEnable"`
	Style     string `json:"style"`
}

type Image struct {
	Format  string `json:"format"`
	URL     string `json:"url"`
	AltText string `json:"altText"`
}

type VariantOption struct {
	Images   []Image `json:"images"`
	SKU      string  `json:"sku"`
	ImageSKU string  `json:"imageSku"`
}

type Product struct {
	Badges        Badge            `json:"badges"`
	BaseOptions   []SelectedOption `json:"baseOptions"`
	BaseProduct   string           `json:"baseProduct"`
	Images        []Image          `json:"images"`
	Name          string           `json:"name"`
	OriginalPrice struct {
		Value          float64 `json:"value"`
		FormattedValue string  `json:"formattedValue"`
	} `json:"originalPrice"`
	Price struct {
		Value          float64 `json:"value"`
		FormattedValue string  `json:"formattedValue"`
	} `json:"price"`
	SKU            string          `json:"sku"`
	URL            string          `json:"url"`
	ImageSKU       string          `json:"imageSku"`
	VariantOptions []VariantOption `json:"variantOptions"`
	VariantsCount  int             `json:"variantsCount"`
	IsSaleProduct  bool            `json:"isSaleProduct"`
	IsNewProduct   bool            `json:"isNewProduct"`
	LaunchProduct  bool            `json:"launchProduct"`
}

type SearchResponse struct {
	CurrentQuery   Query      `json:"currentQuery"`
	Facets         []Facet    `json:"facets"`
	FreeTextSearch string     `json:"freeTextSearch"`
	Pagination     Pagination `json:"pagination"`
	Products       []Product  `json:"products"`
}

type FilterSearch struct {
	Sku   string
	Name  string
	Image string
	Price string
}

type GotShoes struct {
	Name  string   `json:"name"`
	Image string   `json:"image"`
	Price string   `json:"price"`
	Sku   string   `json:"sku"`
	Sizes []string `json:"sizes"`
}

type ProductList struct {
	Products []GotShoes `json:"products"`
}

// Product request information:

type Style struct {
	ID              string   `json:"id"`
	SKU             string   `json:"sku"`
	LanguageIsoCode string   `json:"languageIsoCode"`
	Active          bool     `json:"active"`
	Color           string   `json:"color"`
	PrimaryColor    string   `json:"primaryColor"`
	Price           Price    `json:"price"`
	ImageURL        ImageURL `json:"imageUrl"`
}

type Price struct {
	FormattedSalePrice string `json:"formattedSalePrice"`
}

type ImageURL struct {
	Base     string   `json:"base"`
	ImageSKU string   `json:"imageSku"`
	Variants []string `json:"variants"`
}

type Size struct {
	Active        bool      `json:"active"`
	ProductNumber string    `json:"productNumber"`
	Size          string    `json:"size"`
	Inventory     Inventory `json:"inventory"`
}

type Inventory struct {
	InventoryAvailable          bool        `json:"inventoryAvailable"`
	StoreInventoryAvailable     bool        `json:"storeInventoryAvailable"`
	WarehouseInventoryAvailable bool        `json:"warehouseInventoryAvailable"`
	DropshipInventoryAvailable  bool        `json:"dropshipInventoryAvailable"`
	InventoryAvailableLocations []string    `json:"inventoryAvailableLocations"`
	PreSell                     bool        `json:"preSell"`
	BackOrder                   bool        `json:"backOrder"`
	PurchaseOrderDate           interface{} `json:"purchaseOrderDate"` // Se necessario, sostituisci con un tipo di data adeguato
}

type ProductData struct {
	ID    string `json:"id"`
	Style Style  `json:"style"`
	Sizes []Size `json:"sizes"`
}
