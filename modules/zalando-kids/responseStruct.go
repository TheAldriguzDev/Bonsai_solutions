package zalando_kids

import "github.com/TheAldriguzDev/Bonsai_solutions/packages/database"

//////// Search response

type SearchResponse struct {
	Data struct {
		Collection CollectionData `json:"collection"`
	} `json:"data"`
}

type CollectionData struct {
	Typename                string            `json:"__typename"`
	ID                      string            `json:"id"`
	Title                   string            `json:"title"`
	Subtitle                string            `json:"subtitle"`
	CTA                     string            `json:"cta"`
	Entities                EntitiesData      `json:"entities"`
	Tags                    []TagData         `json:"tags"`
	CollectionURL           CollectionURLData `json:"collectionURL"`
	TrackingContextEntityID string            `json:"trackingContext_entity_id"`
	TrackingContextConfig   TrackingConfig    `json:"trackingContext_config"`
	PrimaryImage            ImageData         `json:"primaryImage"`
	SuggestedQuery          interface{}       `json:"suggestedQuery"`
	Brand                   BrandData         `json:"brand"`
	DisplayFlags            []FlagData        `json:"displayFlags"`
	Tracking                TrackingData      `json:"tracking"`
	Breadcrumbs             BreadcrumbsData   `json:"breadcrumbs"`
}

type EntitiesData struct {
	Typename          string       `json:"__typename"`
	OrderedBy         string       `json:"orderedBy"`
	CatalogResultKind string       `json:"catalogResultKind"`
	Experience        interface{}  `json:"experience"`
	HiddenFilters     interface{}  `json:"hiddenFilters"`
	TotalCount        int          `json:"totalCount"`
	Nodes             []NodeData   `json:"nodes"`
	PageInfo          PageInfoData `json:"pageInfo"`
}

type NodeData struct {
	Typename                 string               `json:"__typename"`
	ID                       string               `json:"id"`
	SKU                      string               `json:"sku"`
	Name                     string               `json:"name"`
	SupplierName             string               `json:"supplierName"`
	ShortDescription         interface{}          `json:"shortDescription"`
	BrandRestriction         BrandRestrictionData `json:"brandRestriction"`
	ShowPlusFlag             bool                 `json:"showPlusFlag"`
	Benefits                 []interface{}        `json:"benefits"`
	Simples                  []SimpleData         `json:"simples"`
	Fulfillment              interface{}          `json:"fulfillment"`
	Brand                    BrandData            `json:"brand"`
	DisplayFlags             []FlagData           `json:"displayFlags"`
	DisplayPrice             DisplayPriceData     `json:"displayPrice"`
	Condition                interface{}          `json:"condition"`
	BasePrice                interface{}          `json:"basePrice"`
	CatalogDefaultImage      ImageData            `json:"catalogDefaultImage"`
	PrimaryImage             ImageData            `json:"primaryImage"`
	NavigationTargetGroup    string               `json:"navigationTargetGroup"`
	Reason                   interface{}          `json:"reason"`
	Sustainability           interface{}          `json:"sustainability"`
	ExperienceWatermarkImage interface{}          `json:"experienceWatermarkImage"`
	Experience               interface{}          `json:"experience"`
	TrackingContextEntityID  string               `json:"trackingContext_entity_id"`
	TrackingContextConfig    TrackingConfig       `json:"trackingContext_config"`
	PlusBenefitBadge         PlusBadgeData        `json:"plusBenefitBadge"`
	ProductDetailsTracking   ProductTrackingData  `json:"productDetailsTracking"`
	CurrentPrice             PriceData            `json:"currentPrice"`
	PriceValue               PriceData            `json:"priceValue"`
}

type BrandRestrictionData struct {
	Typename      string      `json:"__typename"`
	ThemeOverride interface{} `json:"themeOverride"`
}

type SimpleData struct {
	Typename string `json:"__typename"`
	Size     string `json:"size"`
}

type BrandData struct {
	Typename string `json:"__typename"`
	ID       string `json:"id"`
	Name     string `json:"name"`
}

type FlagData struct {
	Typename string `json:"__typename"`
	ID       string `json:"id"`
	Label    string `json:"label"`
}

type DisplayPriceData struct {
	Typename     string `json:"__typename"`
	Formatted    string `json:"formatted"`
	DisplayValue string `json:"displayValue"`
}

type ImageData struct {
	Typename  string `json:"__typename"`
	URL       string `json:"uri"`
	Thumbnail string `json:"thumbnail"`
}

type TrackingConfig struct {
	Typename string `json:"__typename"`
	ID       string `json:"id"`
}

type PageInfoData struct {
	Typename      string `json:"__typename"`
	HasNextPage   bool   `json:"hasNextPage"`
	EndCursor     string `json:"endCursor"`
	RecordPerPage int    `json:"recordPerPage"`
	TotalRecords  int    `json:"totalRecords"`
}

type TagData struct {
	Typename string `json:"__typename"`
	Label    string `json:"label"`
	Value    string `json:"value"`
}

type CollectionURLData struct {
	Typename string `json:"__typename"`
	URL      string `json:"url"`
}

type BreadcrumbsData struct {
	Typename string           `json:"__typename"`
	Nodes    []BreadcrumbData `json:"nodes"`
}

type BreadcrumbData struct {
	Typename string `json:"__typename"`
	Label    string `json:"label"`
	URL      string `json:"url"`
}

type TrackingData struct {
	Typename string `json:"__typename"`
	PageType string `json:"pageType"`
	PageID   string `json:"pageId"`
}

type PlusBadgeData struct {
	Typename string `json:"__typename"`
	ID       string `json:"id"`
	Label    string `json:"label"`
}

type ProductTrackingData struct {
	Typename string `json:"__typename"`
	ID       string `json:"id"`
}

type PriceData struct {
	Typename     string  `json:"__typename"`
	Value        float64 `json:"value"`
	CurrencyCode string  `json:"currencyCode"`
	Formatted    string  `json:"formatted"`
}

//////// Product

type ProductStock struct {
	Quantity string `json:"quantity"`
}

type P_Price struct {
	Typename         string            `json:"__typename"`
	DiscountLabel    interface{}       `json:"discountLabel"`
	DisplayMode      interface{}       `json:"displayMode"`
	Promotional      interface{}       `json:"promotional"`
	Original         ProductPriceValue `json:"original"`
	Previous         interface{}       `json:"previous"`
	HistoryExplainer interface{}       `json:"historyExplainer"`
}

type ProductPriceValue struct {
	Typename      string      `json:"__typename"`
	Formatted     string      `json:"formatted"`
	Amount        int         `json:"amount"`
	Label         interface{} `json:"label"`
	DiscountLabel interface{} `json:"discountLabel"`
}

type ProductOffer struct {
	Stock ProductStock `json:"stock"`
	Price P_Price      `json:"price"`
}

type ProductSimple struct {
	Size                string       `json:"size"`
	BeautySize          string       `json:"beautySize"`
	SKU                 string       `json:"sku"`
	SupplierSize        string       `json:"supplierSize"`
	HasMeaningfulOffers bool         `json:"hasMeaningfulOffers"`
	Offer               ProductOffer `json:"offer"`
}

type ProductColor struct {
	Name string `json:"name"`
}

type Product struct {
	TypeName              string          `json:"__typename"`
	SKU                   string          `json:"sku"`
	ComingSoon            bool            `json:"comingSoon"`
	NavigationTargetGroup string          `json:"navigationTargetGroup"`
	ShowPlusFlag          bool            `json:"showPlusFlag"`
	ShowPrixMiniFlag      bool            `json:"showPrixMiniFlag"`
	URI                   string          `json:"uri"`
	Name                  string          `json:"name"`
	Group                 string          `json:"group"`
	AvailabilityStatus    string          `json:"availabilityStatus"`
	Experience            interface{}     `json:"experience"`
	IsActive              bool            `json:"isActive"`
	Color                 ProductColor    `json:"color"`
	Simples               []ProductSimple `json:"simples"`
}

type ProductData struct {
	Product Product `json:"product"`
}

type ProductResponse struct {
	Data ProductData `json:"data"`
}

type CreateProduct struct {
	Sku   string
	Name  string
	Price string
	Image string
	Data  UpdateProduct
}

type AvailableSizes struct {
	Size  string
	Value string
}

type UpdateProduct struct {
	ComingSoon bool
	Active     bool
	Color      string
	Sizes      []database.Size
}

type ProductList struct {
	Products []CreateProduct `json:"products"`
}
