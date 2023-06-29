package database

type ProductList struct {
	Products []Product `json:"products"`
}

type Product struct {
	Site  string `json:"site"`
	Sku   string `json:"sku"`
	Name  string `json:"name"`
	Price string `json:"price"`
	Image string `json:"image"`
	Sizes []Size `json:"sizes"`
	Color string `json:"color"`
}

type Size struct {
	Size  string `json:"size"`
	Value string `json:"value"`
}
