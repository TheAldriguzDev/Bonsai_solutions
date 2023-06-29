package adidas

type ThumbnailImage struct {
	Href string `json:"href,omitempty"`
}

type Product struct {
	ProductID       string        `json:"product_id,omitempty"`
	ModelNumber     string        `json:"model_number,omitempty"`
	ProductName     string        `json:"product_name,omitempty"`
	ProductFullName string        `json:"product_full_name,omitempty"`
	OriginalPrice   int           `json:"original_price,omitempty"`
	OriginalColor   string        `json:"original_color,omitempty"`
	Gender          string        `json:"gender,omitempty"`
	PurchaseLimit   int           `json:"purchase_limit,omitempty"`
	Orderable       bool          `json:"orderable,omitempty"`
	Embedded        *EmbeddedData `json:"_embedded,omitempty"`
}

type EmbeddedData struct {
	Products       []Product      `json:"products,omitempty"`
	ThumbnailImage ThumbnailImage `json:"thumbnail_image,omitempty"`
}

type SearchData struct {
	Embedded EmbeddedData `json:"_embedded,omitempty"`
}

// Stock response:
type Variation struct {
	SKU                string `json:"sku"`
	Size               string `json:"size"`
	Availability       int    `json:"availability"`
	AvailabilityStatus string `json:"availability_status"`
}

type ProductStock struct {
	ID            string      `json:"id"`
	Availability  string      `json:"availability_status"`
	VariationList []Variation `json:"variation_list"`
}
