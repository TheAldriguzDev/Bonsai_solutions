package rinascente

type AlgoliaSearchRequest struct {
	OperationName string    `json:"operationName"`
	Variables     Variables `json:"variables"`
	Query         string    `json:"query"`
}

type Variables struct {
	Query                   string   `json:"query"`
	Page                    int      `json:"page"`
	Size                    int      `json:"size"`
	SortOrder               string   `json:"sortOrder"`
	FilterGroups            []string `json:"filterGroups"`
	NumericFilterGroups     []string `json:"numericFilterGroups"`
	CategoryID              string   `json:"categoryId"`
	UseTargetedProductIndex bool     `json:"useTargetedProductIndex"`
}
