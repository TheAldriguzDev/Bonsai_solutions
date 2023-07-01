package rinascente

type AlgoliaSearchResponse struct {
	Data struct {
		AlgoliaSearch struct {
			Error    interface{} `json:"error,omitempty"`
			Products []Product   `json:"products,omitempty"`
		} `json:"algoliaSearch,omitempty"`
	} `json:"data,omitempty"`
}

type Product struct {
	Name                string              `json:"name,omitempty"`
	VisibilitySearch    int                 `json:"visibility_search,omitempty"`
	VisibilityCatalog   int                 `json:"visibility_catalog,omitempty"`
	TypeID              string              `json:"type_id,omitempty"`
	SKU                 string              `json:"sku,omitempty"`
	Brand               string              `json:"brand,omitempty"`
	SizeInt             []string            `json:"size_int,omitempty"`
	MainColour          string              `json:"main_colour,omitempty"`
	URLKey              string              `json:"url_key,omitempty"`
	Price               int                 `json:"price,omitempty"`
	MediaGalleryEntries []MediaGalleryEntry `json:"media_gallery_entries,omitempty"`
}

type MediaGalleryEntry struct {
	ID                  string      `json:"id,omitempty"`
	Disabled            string      `json:"disabled,omitempty"`
	Position            string      `json:"position,omitempty"`
	MediaType           string      `json:"media_type,omitempty"`
	Label               string      `json:"label,omitempty"`
	File                string      `json:"file,omitempty"`
	Types               []string    `json:"types,omitempty"`
	ExtensionAttributes interface{} `json:"extension_attributes,omitempty"`
	Typename            string      `json:"__typename,omitempty"`
}
