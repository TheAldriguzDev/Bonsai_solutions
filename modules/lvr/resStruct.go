package lvr

type Item struct {
	ItemCode                   string         `json:"ItemCode"`
	SeasonId                   string         `json:"SeasonId"`
	CollectionId               string         `json:"CollectionId"`
	ItemId                     string         `json:"ItemId"`
	VendorColorId              string         `json:"VendorColorId"`
	Description                string         `json:"Description"`
	Designer                   string         `json:"Designer"`
	Image                      string         `json:"Image"`
	Sizes                      string         `json:"Sizes"`
	MultiPrice                 bool           `json:"MultiPrice"`
	ListPrice                  string         `json:"ListPrice"`
	ListPriceDiscounted        string         `json:"ListPriceDiscounted"`
	FinalPrice                 interface{}    `json:"FinalPrice"`
	FinalPriceEur              int            `json:"FinalPriceEur"`
	Discount                   int            `json:"Discount"`
	FurtherReduction           int            `json:"FurtherReduction"`
	FurtherReductionPriceLabel string         `json:"FurtherReductionPriceLabel"`
	PromoReduction             int            `json:"PromoReduction"`
	PromoReductionPriceLabel   interface{}    `json:"PromoReductionPriceLabel"`
	PromoReductionType         string         `json:"PromoReductionType"`
	Tags                       []string       `json:"Tags"`
	Badges                     []string       `json:"Badges"`
	SortedTags                 []string       `json:"SortedTags"`
	DYTag                      string         `json:"DYTag"`
	Note                       interface{}    `json:"Note"`
	Variants                   []string       `json:"Variants"`
	UniqueId                   string         `json:"UniqueId"`
	ItemParameters             ItemParameters `json:"ItemParameters"`
	ExtraInfo                  ExtraInfo      `json:"ExtraInfo"`
	IsMadeToMeasure            bool           `json:"IsMadeToMeasure"`
	OfferMetaInfo              OfferMetaInfo  `json:"OfferMetaInfo"`
	PrivateSaleInfoList        interface{}    `json:"PrivateSaleInfoList"`
	QOS                        interface{}    `json:"QOS"`
}

type ItemParameters struct {
	ItemCode      string `json:"ItemCode"`
	SeasonId      string `json:"SeasonId"`
	CollectionId  string `json:"CollectionId"`
	ItemId        string `json:"ItemId"`
	VendorColorId string `json:"VendorColorId"`
	SeasonMemo    string `json:"SeasonMemo"`
	GenderMemo    string `json:"GenderMemo"`
}

type ExtraInfo struct {
	Available    string `json:"Available"`
	Size         string `json:"Size"`
	Colors       string `json:"Colors"`
	PriceTooltip string `json:"PriceTooltip"`
}

type OfferMetaInfo struct {
	Availability string `json:"Availability"`
	Currency     string `json:"Currency"`
	FinalPrice   string `json:"FinalPrice"`
}

type Data struct {
	Items []Item `json:"Items"`
}
