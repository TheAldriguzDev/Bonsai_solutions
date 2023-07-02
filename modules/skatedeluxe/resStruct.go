package skatedeluxe

// Search response

type Product struct {
	ProductsID           int      `json:"products_id"`
	ProductsImage        string   `json:"products_image"`
	ProductsNameFull     string   `json:"products_name_full"`
	ProductsName         string   `json:"products_name"`
	ProductMaskID        string   `json:"product_mask_id"`
	ProductGroupID       string   `json:"product_group_id"`
	ProductsPriceOld     float64  `json:"products_price_old"`
	ProductsPriceCalced  float64  `json:"products_price_calced"`
	ProductsPreviewImage string   `json:"products_previewimage"`
	ProductsLink         string   `json:"products_link"`
	ProductImages        []string `json:"product_images"`
}

type SearchResponse struct {
	Status     string    `json:"status"`
	StatusCode int       `json:"status_code"`
	Products   []Product `json:"products"`
}
