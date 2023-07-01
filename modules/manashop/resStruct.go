package manashop

type SearchRes struct {
	Result Result `json:"result"`
}

type Result struct {
	Product Product `json:"product"`
}

type Product struct {
	Data []Data `json:"data"`
}

type Data struct {
	Name             string           `json:"name"`
	Image            string           `json:"image"`
	Price            string           `json:"price"`
	URL              string           `json:"url"`
	ShortDescription ShortDescription `json:"0"`
	ReviewsRating    ReviewsRating    `json:"1"`
}

type ShortDescription struct {
	Description string `json:"short_description"`
}

type ReviewsRating struct {
	Rating string `json:"reviews_rating"`
}

// // Product res

type Option struct {
	ID       string   `json:"id"`
	Label    string   `json:"label"`
	Products []string `json:"products"`
}

type Attribute struct {
	ID       string   `json:"id"`
	Code     string   `json:"code"`
	Label    string   `json:"label"`
	Options  []Option `json:"options"`
	Position string   `json:"position"`
}

type OptionPrice struct {
	OldPrice   Price   `json:"oldPrice"`
	BasePrice  Price   `json:"basePrice"`
	FinalPrice Price   `json:"finalPrice"`
	TierPrices []Price `json:"tierPrices"`
	MsrpPrice  Price   `json:"msrpPrice"`
}

type Price struct {
	Amount float64 `json:"amount"`
}

type JSONConfig struct {
	Attributes   map[string]Attribute   `json:"attributes"`
	OptionPrices map[string]OptionPrice `json:"optionPrices"`
	Prices       struct {
		OldPrice   Price `json:"oldPrice"`
		BasePrice  Price `json:"basePrice"`
		FinalPrice Price `json:"finalPrice"`
	} `json:"prices"`
	ProductID string `json:"productId"`
}

type SwatchOptions struct {
	SwatchRenderer struct {
		JSONConfig JSONConfig `json:"jsonConfig"`
	} `json:"Magento_Swatches/js/swatch-renderer"`
}

type MyData struct {
	SwatchOptions SwatchOptions `json:"[data-role=swatch-options]"`
}
