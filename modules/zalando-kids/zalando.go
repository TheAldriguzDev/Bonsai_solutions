package zalando_kids

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/TheAldriguzDev/Bonsai_solutions/packages/database"
	"github.com/TheAldriguzDev/Bonsai_solutions/packages/log"
	proxyManager "github.com/TheAldriguzDev/Bonsai_solutions/packages/proxy"
	filereader "github.com/TheAldriguzDev/Bonsai_solutions/packages/reader"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

var location string = "Zalando Kids"

func InitTls(proxy string) (zalandoSession, error) {
	proxy = strings.ReplaceAll(proxy, "\r", "")
	proxyData := strings.Split(strings.TrimSpace(proxy), ":")
	proxyUrl := fmt.Sprintf("http://%s:%s@%s:%s", proxyData[2], proxyData[3], proxyData[0], proxyData[1])

	jarOptions := []tls_client.CookieJarOption{}

	jar := tls_client.NewCookieJar(jarOptions...)

	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(tls_client.Chrome_112),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
		tls_client.WithRandomTLSExtensionOrder(),
		tls_client.WithProxyUrl(proxyUrl),
		tls_client.WithInsecureSkipVerify(),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)

	if err != nil {
		fmt.Println("Error while creating the tls client!")
		return zalandoSession{}, err
	}

	return zalandoSession{
		httpClient: client,
	}, nil
}

func InitZalando() {
	proxyFile, err := filereader.ReadFolder("proxies", "proxies")

	if err != nil {
		log.Error(err.Error(), location)
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}

	proxyArr, err := proxyManager.ReadProxyFile(proxyFile)

	if err != nil {
		log.Error(err.Error(), location)
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}

	proxie, err := proxyManager.RandomProxy(proxyArr)

	if err != nil {
		log.Error(err.Error(), location)
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}

	module, err := InitTls(proxie)
	module.Proxy = proxie

	if err != nil {
		log.Error(err.Error(), location)
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}

	module.Host = "https://www.zalando.it/"
	module.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"
	module.Store_Region = "it"
	module.ProxyArr = proxyArr
	log.Success("Monitor initialised!", location)

	var run bool = true

	for run {
		log.General("Monitoring...", location)

		err = module.RotateProxy()
		if err != nil {
			log.Error(err.Error(), location)
		}

		err = module.SearchProducts("ern:collection:cat:categ:bambini")
		if err != nil {
			log.Error(err.Error(), location)
			time.Sleep(3 * time.Second)
			continue
		}

		time.Sleep(3 * time.Second)
	}
}

func (zala *zalandoSession) RotateProxy() error {
	proxie, err := proxyManager.RandomProxy(zala.ProxyArr)
	if err != nil {
		return err
	}
	proxie = strings.ReplaceAll(proxie, "\r", "")
	proxyData := strings.Split(strings.TrimSpace(proxie), ":")
	if len(proxyData) < 4 {
		return fmt.Errorf("The proxy format is not correct! %s", proxie)
	}
	proxyUrl := fmt.Sprintf("http://%s:%s@%s:%s", proxyData[2], proxyData[3], proxyData[0], proxyData[1])
	if err != nil {
		return err
	} else {
		err = zala.httpClient.SetProxy(proxyUrl)
		if err != nil {
			return err
		}
		zala.Proxy = proxie
	}
	return nil
}

func (zalando *zalandoSession) SearchProducts(id string) error {

	var searchResponse SearchResponse
	searchData := SearchData{
		ID: "9a92659a9266bba8fe0ab012cb343fae9df622df9fc2a5cf368aa22e503a56db",
		Extensions: struct {
			PersistedQuery PersistedQuery `json:"persistedQuery"`
		}{
			PersistedQuery: PersistedQuery{
				Version:    "1",
				Sha256Hash: "9a92659a9266bba8fe0ab012cb343fae9df622df9fc2a5cf368aa22e503a56db",
			},
		},
		OperationName: "appCatalog",
		Variables: SearchVariables{
			ShouldIncludeWatermarkExperience:      true,
			AepNumberOfEntities:                   7,
			SegmentGatedLogoWidth:                 100,
			CategoryImageWidth:                    300,
			First:                                 genNum(400, 1000),
			ShouldIncludeOmnibusPrice:             true,
			GardenExperienceLogoWidth:             300,
			IsFirstPage:                           true,
			ProductCardImageWidth:                 500,
			ShouldIncludeSustainabilityClaims:     true,
			ShouldIncludeSegmentGatedEnabled:      true,
			ShouldIncludeCampaignReleventEntities: true,
			ShouldIncludeExperience:               true,
			BrandedBarLogoWidth:                   300,
			ShouldIncludeHightlightCollectionData: true,
			ID:                                    id,
			IsFDBERequested:                       false,
			ShouldIncludeThemeExperience:          true,
			HighlightNumberOfEntities:             genNum(20, 300),
			ShouldSkipCatalogInGridCampaign:       true,
			IsRecoCatalog:                         false,
			ShouldIncludeFullExperience:           true,
			Order:                                 "NEW_IN",
			ExperienceWatermarkLogoWidth:          300,
			Filters: Filters{
				DiscreteFilters: []interface{}{},
				RangeFilters: []RangeFilter{
					{
						Key:     "prices",
						Minimum: 20,
						Maximum: genNum(130, 170),
					},
				},
				FullTextQuery: "campus 00",
			},
		},
	}

	jsonPayload, err := json.Marshal(searchData)

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://www.zalando.%s/api/graphql/mobile?ts=%s", zalando.Store_Region, GetTimeStamp()), bytes.NewBuffer(jsonPayload))

	if err != nil {
		return err
	}

	req.Header = http.Header{
		"Host":                         {"www.zalando.it"},
		"accept":                       {"application/json"},
		"accept-language":              {"it-IT"},
		"apollographql-client-name":    {"de.zalando.iphone-apollo-ios"},
		"apollographql-client-version": {"23.8.0-60621"},
		"content-type":                 {"application/json"},
		"user-agent":                   {"zalando/23.8.0 (iPhone; iOS 16.5; Scale/3.00)"},
		"x-apollo-operation-id":        {"f80a3f6d34066528edc1f9c97abf4be5cc02ceec02455c145d7e2e3db9b6e65e"},
		"x-apollo-operation-name":      {"Pdp"},
		"x-apollo-operation-type":      {"query"},
		"x-app-domain":                 {"15"},
		"x-app-version":                {"23.8.0"},
		"x-device-os":                  {"ios"},
		"x-device-platform":            {"ios"},
		"x-device-type":                {"smartphone"},
		"x-frontend-type":              {"mobile-app"},
		"x-logged-in":                  {"true"},
		"x-os-version":                 {"16.5"},
		//"x-sales-channel": "ebf57ebf-e26d-4ebd-8009-6ad519073d2a",
		"x-uuid": {"5A2C34BE-CDC4-4B5F-9FA8-2924FAD51AD9"},
		//"x-zalando-client-id":      {"f39a28b9-56c5-4481-817f-6684dcd97a98"},
		"x-zalando-debug":   {"app-pdp-recos=true"},
		"x-zalando-feature": {"pdp"},
		//"x-zalando-intent-context": {"navigationTargetGroup=MEN;experience=RESOLVE"},
		"x-zalando-mobile-app": {"3580f92a4bafb890i"},
		//"Cookie":                   {"Zalando-Client-Id=f39a28b9-56c5-4481-817f-6684dcd97a98; ak_bmsc=D8D00258A828D590AC2E10A9E1AA9151~00000000000000000"},
	}

	res, err := zalando.httpClient.Do(req)
	if err != nil {
		return err
	}

	if res != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Bad status code [SEARCH]")
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &searchResponse)

	if err != nil {
		return err
	}

	err = zalando.ParseSearch(searchResponse)
	if err != nil {
		return err
	}
	return nil
}

func (zalando *zalandoSession) ParseSearch(respData SearchResponse) error {
	var productList []database.Product

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, product := range respData.Data.Collection.Entities.Nodes {
		if product.SKU != "AD115O1JQ-C12" && product.SKU != "AD115O1J6-B11" && product.SKU != "AD115O1J6-Q11" && product.SKU != "AD115O1KC-K11" && product.SKU != "" {
			wg.Add(1)

			go zalando.GetInfoUtils(product, &productList, &wg, &mu)
		}
	}

	// forceProductList := []ForceProduct{
	// 	{
	// 		Sku:   "AD111A28H-A14",
	// 		Name:  "CAMPUS 00S",
	// 		Price: "109.99 €",
	// 		Image: "https://images-ext-1.discordapp.net/external/9AVhc8FZJtW9Vcs349upbOnKxlFmbJOo3r6Li0GGTsg/%3Fimwidth%3D500%26filter%3Dpackshot/https/img01.ztat.net/article/spp-media-p1/b1e9facf3a01443c8b0633716bfca5ea/07008daa57da4e22b78e8b4ba88acd50.jpg?width=718&height=1036",
	// 	},
	// }

	// for _, product := range forceProductList {
	// 	wg.Add(1)
	// 	go zalando.GetForceInfo(product, &productList, &wg, &mu)
	// }

	wg.Wait()

	err := database.StartDataBase("./monitorData/zalando_kids.json", productList)
	if err != nil {
		return err
	}

	return nil
}

func (zalando *zalandoSession) GetInfoUtils(product NodeData, productList *[]database.Product, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()

	tempProduct := database.Product{
		Site:  "Zalando",
		Sku:   product.SKU,
		Name:  product.SupplierName,
		Price: product.DisplayPrice.Formatted,
		Image: product.CatalogDefaultImage.URL,
	}

	info, err := zalando.GetStock(tempProduct.Sku)
	if err != nil {
		log.Error(err.Error(), location)
		return
	}

	tempProduct.Sizes = info.Sizes
	tempProduct.Color = info.Color

	if len(tempProduct.Sizes) != 0 {
		mu.Lock()
		*productList = append(*productList, tempProduct)
		mu.Unlock()
	}
}

type ForceProduct struct {
	Name  string
	Image string
	Price string
	Sku   string
}

func (zalando *zalandoSession) GetForceInfo(product ForceProduct, productList *[]database.Product, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()

	tempProduct := database.Product{
		Site:  "Zalando",
		Sku:   product.Sku,
		Name:  product.Name,
		Price: product.Price,
		Image: product.Image,
	}

	info, err := zalando.GetStock(tempProduct.Sku)
	if err != nil {
		log.Error(err.Error(), location)
		return
	}

	tempProduct.Sizes = info.Sizes
	tempProduct.Color = info.Color

	if len(tempProduct.Sizes) != 0 {
		mu.Lock()
		*productList = append(*productList, tempProduct)
		mu.Unlock()
	}
}

func (zalando zalandoSession) GetStock(sku string) (UpdateProduct, error) {
	var productResponse ProductResponse
	var productData UpdateProduct

	// Rotate proxies here
	// We rotate the proxy for every product we need to scrape
	proxie, err := proxyManager.RandomProxy(zalando.ProxyArr)
	if err != nil {
		log.Error(err.Error(), location)
		return productData, err
	}
	proxie = strings.ReplaceAll(proxie, "\r", "")
	proxyData := strings.Split(strings.TrimSpace(proxie), ":")
	if len(proxyData) < 4 {
		log.Error("The proxy format is not correct!", location)
		return productData, fmt.Errorf("The proxy format is not correct! %s", proxie)
	}
	proxyUrl := fmt.Sprintf("http://%s:%s@%s:%s", proxyData[2], proxyData[3], proxyData[0], proxyData[1])
	if err != nil {
		log.Error(err.Error(), location)
		return productData, err
	} else {
		err = zalando.httpClient.SetProxy(proxyUrl)
		if err != nil {
			log.Error(err.Error(), location)
			return productData, err
		}
		zalando.Proxy = proxie
	}

	payload := Data{
		Extensions: struct {
			PersistedQuery PersistedQuery `json:"persistedQuery"`
		}{
			PersistedQuery: PersistedQuery{
				Sha256Hash: "f80a3f6d34066528edc1f9c97abf4be5cc02ceec02455c145d7e2e3db9b6e65e",
				Version:    "1",
			},
		},
		ID:            "f80a3f6d34066528edc1f9c97abf4be5cc02ceec02455c145d7e2e3db9b6e65e",
		OperationName: "Pdp",
		Variables: Variables{
			BeautyColorImageWidth:                       1,
			BenefitsLogoWidth:                           genNum(1, 10000),
			BrandedBarLogoWidth:                         genNum(200, 300),
			BrandLogoWidth:                              genNum(600, 800),
			ColorImageWidth:                             76,
			ConfigSku:                                   sku,
			ExperienceLogoWidth:                         genNum(200, 400),
			FullScreenGalleryWidth:                      1200,
			FullScreenHdGalleryWidth:                    2600,
			MaxFlagCount:                                genNum(3, 5),
			PortraitGalleryWidth:                        1170,
			SegmentedBannerHeaderLogoWidth:              108,
			ShouldIncludeFlagInfo:                       false,
			ShouldIncludeHistogramValues:                true,
			ShouldIncludeOfferSelectionValues:           true,
			ShouldIncludeOmnibusConfigModeChanges:       true,
			ShouldIncludeOmnibusPrice:                   true,
			ShouldIncludePrestigeBeauty:                 false,
			ShouldIncludePurchaseRestriction:            true,
			ShouldIncludeSamples:                        true,
			ShouldIncludeSegmentation:                   true,
			ShouldIncludeSegmentGatedComingSoonReminder: true,
			ShouldIncludeSizeAdvicewBM:                  false,
			ShouldIncludeSubscriptionValues:             false,
			ShouldIncludeSustainabilityClaims:           true,
			ShouldIncludeTargetedCoupons:                true,
			ShouldIncludeUSPReturnPolicy:                true,
		},
	}

	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return productData, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://www.zalando.%s/api/graphql/mobile?ts=%s", zalando.Store_Region, GetTimeStamp()), bytes.NewBuffer(jsonPayload))

	if err != nil {
		return productData, err
	}

	req.Header = http.Header{
		"Host":                         {"www.zalando.it"},
		"accept":                       {"application/json"},
		"accept-language":              {"it-IT"},
		"apollographql-client-name":    {"de.zalando.iphone-apollo-ios"},
		"apollographql-client-version": {"23.8.0-60621"},
		"content-type":                 {"application/json"},
		"user-agent":                   {"zalando/23.8.0 (iPhone; iOS 16.5; Scale/3.00)"},
		"x-apollo-operation-id":        {"f80a3f6d34066528edc1f9c97abf4be5cc02ceec02455c145d7e2e3db9b6e65e"},
		"x-apollo-operation-name":      {"Pdp"},
		"x-apollo-operation-type":      {"query"},
		"x-app-domain":                 {"15"},
		"x-app-version":                {"23.8.0"},
		"x-device-os":                  {"ios"},
		"x-device-platform":            {"ios"},
		"x-device-type":                {"smartphone"},
		"x-frontend-type":              {"mobile-app"},
		"x-logged-in":                  {"true"},
		"x-os-version":                 {"16.5"},
		"x-sales-channel":              {"ebf57ebf-e26d-4ebd-8009-6ad519073d2a"},
		"x-uuid":                       {"5A2C34BE-CDC4-4B5F-9FA8-2924FAD51AD9"},
		"x-zalando-client-id":          {"f39a28b9-56c5-4481-817f-6684dcd97a98"},
		"x-zalando-debug":              {"app-pdp-recos=true"},
		"x-zalando-feature":            {"pdp"},
		//"x-zalando-intent-context": {"navigationTargetGroup=MEN;experience=RESOLVE"},
		"x-zalando-mobile-app": {"3580f92a4bafb890i"},
		//"Cookie":                   {"Zalando-Client-Id=f39a28b9-56c5-4481-817f-6684dcd97a98; ak_bmsc=D8D00258A828D590AC2E10A9E1AA9151~00000000000000000"},
	}

	res, err := zalando.httpClient.Do(req)

	if err != nil {
		return productData, err
	}

	if res != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		return productData, fmt.Errorf("Bad status code")
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return productData, err
	}

	err = json.Unmarshal(body, &productResponse)

	if err != nil {
		return productData, err
	}

	productData = productParser(productResponse)

	if len(productData.Sizes) == 0 {
		return productData, fmt.Errorf("No sizes has been found!")
	}

	return productData, nil
}

func productParser(product ProductResponse) UpdateProduct {
	updateData := UpdateProduct{
		ComingSoon: product.Data.Product.ComingSoon,
		Active:     product.Data.Product.IsActive,
		Color:      product.Data.Product.Color.Name,
	}

	var availableSizes []database.Size

	for _, size := range product.Data.Product.Simples {
		if size.Offer.Stock.Quantity != "OUT_OF_STOCK" {
			availableSizes = append(availableSizes, database.Size{
				Size:  size.Size,
				Value: size.Offer.Price.Original.Formatted,
			})
		}
	}

	updateData.Sizes = availableSizes

	return updateData
}

func genNum(min, max int) int {
	// Imposta un seme diverso per ogni esecuzione
	rand.Seed(time.Now().UnixNano())

	// Genera un numero casuale compreso tra min e max
	return rand.Intn(max-min+1) + min
}
