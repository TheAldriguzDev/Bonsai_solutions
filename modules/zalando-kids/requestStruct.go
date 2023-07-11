package zalando_kids

import (
	"strconv"
	"time"
)

type PersistedQuery struct {
	Sha256Hash string `json:"sha256Hash"`
	Version    string `json:"version"`
}

type Variables struct {
	BeautyColorImageWidth                       int    `json:"beautyColorImageWidth"`
	BenefitsLogoWidth                           int    `json:"benefitsLogoWidth"`
	BrandedBarLogoWidth                         int    `json:"brandedBarLogoWidth"`
	BrandLogoWidth                              int    `json:"brandLogoWidth"`
	ColorImageWidth                             int    `json:"colorImageWidth"`
	ConfigSku                                   string `json:"configSku"`
	ExperienceLogoWidth                         int    `json:"experienceLogoWidth"`
	FullScreenGalleryWidth                      int    `json:"fullScreenGalleryWidth"`
	FullScreenHdGalleryWidth                    int    `json:"fullScreenHdGalleryWidth"`
	MaxFlagCount                                int    `json:"maxFlagCount"`
	PortraitGalleryWidth                        int    `json:"portraitGalleryWidth"`
	SegmentedBannerHeaderLogoWidth              int    `json:"segmentedBannerHeaderLogoWidth"`
	ShouldIncludeFlagInfo                       bool   `json:"shouldIncludeFlagInfo"`
	ShouldIncludeHistogramValues                bool   `json:"shouldIncludeHistogramValues"`
	ShouldIncludeOfferSelectionValues           bool   `json:"shouldIncludeOfferSelectionValues"`
	ShouldIncludeOmnibusConfigModeChanges       bool   `json:"shouldIncludeOmnibusConfigModeChanges"`
	ShouldIncludeOmnibusPrice                   bool   `json:"shouldIncludeOmnibusPrice"`
	ShouldIncludePrestigeBeauty                 bool   `json:"shouldIncludePrestigeBeauty"`
	ShouldIncludePurchaseRestriction            bool   `json:"shouldIncludePurchaseRestriction"`
	ShouldIncludeSamples                        bool   `json:"shouldIncludeSamples"`
	ShouldIncludeSegmentation                   bool   `json:"shouldIncludeSegmentation"`
	ShouldIncludeSegmentGatedComingSoonReminder bool   `json:"shouldIncludeSegmentGatedComingSoonReminder"`
	ShouldIncludeSizeAdvicewBM                  bool   `json:"shouldIncludeSizeAdvicewBM"`
	ShouldIncludeSubscriptionValues             bool   `json:"shouldIncludeSubscriptionValues"`
	ShouldIncludeSustainabilityClaims           bool   `json:"shouldIncludeSustainabilityClaims"`
	ShouldIncludeTargetedCoupons                bool   `json:"shouldIncludeTargetedCoupons"`
	ShouldIncludeUSPReturnPolicy                bool   `json:"shouldIncludeUSPReturnPolicy"`
}

type Data struct {
	Extensions struct {
		PersistedQuery PersistedQuery `json:"persistedQuery"`
	} `json:"extensions"`
	ID            string    `json:"id"`
	OperationName string    `json:"operationName"`
	Variables     Variables `json:"variables"`
}

// Search struct

type RangeFilter struct {
	Key     string `json:"key"`
	Minimum int    `json:"minimum"`
	Maximum int    `json:"maximum"`
}

type Filters struct {
	DiscreteFilters []interface{} `json:"discreteFilters"`
	RangeFilters    []RangeFilter `json:"rangeFilters"`
	FullTextQuery   string        `json:"fullTextQuery"`
}

type SearchVariables struct {
	ShouldIncludeWatermarkExperience      bool    `json:"shouldIncludeWatermarkExperience"`
	AepNumberOfEntities                   int     `json:"aepNumberOfEntities"`
	SegmentGatedLogoWidth                 int     `json:"segmentGatedLogoWidth"`
	CategoryImageWidth                    int     `json:"categoryImageWidth"`
	First                                 int     `json:"first"`
	ShouldIncludeOmnibusPrice             bool    `json:"shouldIncludeOmnibusPrice"`
	GardenExperienceLogoWidth             int     `json:"gardenExperienceLogoWidth"`
	IsFirstPage                           bool    `json:"isFirstPage"`
	ProductCardImageWidth                 int     `json:"productCardImageWidth"`
	ShouldIncludeSustainabilityClaims     bool    `json:"shouldIncludeSustainabilityClaims"`
	ShouldIncludeSegmentGatedEnabled      bool    `json:"shouldIncludeSegmentGatedEnabled"`
	ShouldIncludeCampaignReleventEntities bool    `json:"shouldIncludeCampaignReleventEntities"`
	ShouldIncludeExperience               bool    `json:"shouldIncludeExperience"`
	BrandedBarLogoWidth                   int     `json:"brandedBarLogoWidth"`
	ShouldIncludeHightlightCollectionData bool    `json:"shouldIncludeHightlightCollectionData"`
	ID                                    string  `json:"id"`
	IsFDBERequested                       bool    `json:"isFDBERequested"`
	ShouldIncludeThemeExperience          bool    `json:"shouldIncludeThemeExperience"`
	HighlightNumberOfEntities             int     `json:"highlightNumberOfEntities"`
	ShouldSkipCatalogInGridCampaign       bool    `json:"shouldSkipCatalogInGridCampaign"`
	IsRecoCatalog                         bool    `json:"isRecoCatalog"`
	ShouldIncludeFullExperience           bool    `json:"shouldIncludeFullExperience"`
	Order                                 string  `json:"order"`
	ExperienceWatermarkLogoWidth          int     `json:"experienceWatermarkLogoWidth"`
	Filters                               Filters `json:"filters"`
}

type SearchData struct {
	ID         string `json:"id"`
	Extensions struct {
		PersistedQuery PersistedQuery `json:"persistedQuery"`
	} `json:"extensions"`
	OperationName string          `json:"operationName"`
	Variables     SearchVariables `json:"variables"`
}

func GetTimeStamp() string {
	currentTime := time.Now()
	timestamp := currentTime.UnixNano() / int64(time.Millisecond)
	cacheBuster := strconv.Itoa(int(timestamp))
	return cacheBuster
}
