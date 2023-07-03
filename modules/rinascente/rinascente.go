package rinascente

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TheAldriguzDev/Bonsai_solutions/packages/database"
	"github.com/TheAldriguzDev/Bonsai_solutions/packages/log"
	proxyManager "github.com/TheAldriguzDev/Bonsai_solutions/packages/proxy"
	filereader "github.com/TheAldriguzDev/Bonsai_solutions/packages/reader"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

var location string = "Rinascente"

func InitTls(proxy string) (rinascenteSession, error) {
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
		return rinascenteSession{}, err
	}

	return rinascenteSession{
		httpClient: client,
	}, nil
}

func InitRina() {
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

	module.Host = "https://www.rinascente.it/"
	module.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"
	module.Store_Region = "it"
	module.ProxyArr = proxyArr
	log.Success("Monitor initialised!", location)

	var run bool = true
	for run {
		log.General("Monitoring...", location)
		proxie, err := proxyManager.RandomProxy(proxyArr)
		proxie = strings.ReplaceAll(proxie, "\r", "")
		proxyData := strings.Split(strings.TrimSpace(proxie), ":")
		proxyUrl := fmt.Sprintf("http://%s:%s@%s:%s", proxyData[2], proxyData[3], proxyData[0], proxyData[1])
		if err != nil {
			log.Error(err.Error(), location)
		} else {
			err = module.httpClient.SetProxy(proxyUrl)
			if err != nil {
				log.Error(err.Error(), location)
			}
			module.Proxy = proxie
		}

		err = module.GetProducts()
		if err != nil {
			log.Error(err.Error(), location)
			time.Sleep(7 * time.Second)
			continue
		}
		time.Sleep(7 * time.Second)
	}
}

func (rina *rinascenteSession) GetProducts() error {
	reqData := AlgoliaSearchRequest{
		OperationName: "algoliaSearch",
		Variables: Variables{
			Query:                   "campus",
			Page:                    0,
			Size:                    20,
			SortOrder:               "",
			FilterGroups:            []string{"categoryIds:508", "brand:Adidas Originals"},
			NumericFilterGroups:     []string{},
			CategoryID:              "508",
			UseTargetedProductIndex: false,
		},
		Query: "query algoliaSearch($query: String, $page: Int, $size: Int, $sortOrder: String, $filterGroups: JSON!, $numericFilterGroups: [String], $categoryId: String!, $useTargetedProductIndex: Boolean) {\n  algoliaSearch(\n    query: $query\n    page: $page\n    size: $size\n    sortOrder: $sortOrder\n    filterGroups: $filterGroups\n    numericFilterGroups: $numericFilterGroups\n    categoryId: $categoryId\n    useTargetedProductIndex: $useTargetedProductIndex\n  ) {\n    error\n    filters {\n      attribute_code\n      name\n      position\n      items\n      __typename\n    }\n    priceFilters {\n      lowerBound\n      upperBound\n      __typename\n    }\n    sorting {\n      code\n      __typename\n    }\n    products {\n      name\n      visibility_search\n      visibility_catalog\n      type_id\n      categories {\n        level0\n        level1\n        level2\n        level3\n        __typename\n      }\n      categoryIds\n      in_stock\n      sku\n      brand\n      closure\n      size_int\n      seasonality_filter\n      material\n      main_colour\n      composition_supplier\n      commoditygroup_level1\n      commoditygroup_level2\n      commoditygroup_level3\n      url_key\n      price\n      created_at\n      color_variants {\n        hex_code\n        url_key\n        __typename\n      }\n      additional_price_data\n      price_discount\n      price_saving\n      extension_attributes {\n        brand {\n          brand_id\n          url_key\n          name\n          __typename\n        }\n        __typename\n      }\n      media_gallery_entries {\n        id\n        disabled\n        position\n        media_type\n        label\n        file\n        types\n        extension_attributes {\n          video_content {\n            video_url\n            __typename\n          }\n          __typename\n        }\n        __typename\n      }\n      custom_attributes_option\n      segmentation_price\n      breadcrumbs {\n        category_id\n        name\n        url\n        level\n        __typename\n      }\n      child_stock\n      configurable_product_options\n      meta_title\n      special_price\n      price_min\n      price_max\n      multiple_price\n      product_id\n      sold_qty\n      objectID\n      __typename\n    }\n    categoryQuickFilter {\n      parent_category_id\n      filters {\n        label\n        value\n        custom_attributes {\n          path\n          is_active\n          level\n          parent_id\n          children_count\n          product_count\n          name\n          position\n          entity_id\n          is_parent\n          url_path\n          __typename\n        }\n        __typename\n      }\n      __typename\n    }\n    total_count\n    pagination {\n      page\n      maxPage\n      size\n      __typename\n    }\n    userData {\n      products\n      __typename\n    }\n    tracking {\n      query\n      queryID\n      processingTimeMS\n      indexName\n      serverUsed\n      appliedFilters {\n        facetFilters\n        numericFilters\n        __typename\n      }\n      appliedRules {\n        objectID\n        __typename\n      }\n      serverTimeMS\n      __typename\n    }\n    __typename\n  }\n}\n",
	}

	jsonData, err := json.Marshal(reqData)

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://coreapi.rinascente.it/production/graphql?operationName=algoliaSearch&categoryId=508&query=campus&useTargetedProductIndex=false", bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	req.Header = http.Header{
		"authority":       {"coreapi.rinascente.it"},
		"accept":          {"*/*"},
		"accept-language": {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
		"authorization":   {""},
		"content-type":    {"application/json"},
		"user-agent":      {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"store":           {"it"},
	}

	res, err := rina.httpClient.Do(req)

	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Search bad status code! [%v]", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	var productResp AlgoliaSearchResponse

	err = json.Unmarshal(body, &productResp)

	if err != nil {
		return err
	}

	prodList := rina.ParseResp(productResp)

	err = database.StartDataBase("./monitorData/rinascente.json", prodList)
	if err != nil {
		return err
	}

	return nil
}

func (rina *rinascenteSession) ParseResp(searchResp AlgoliaSearchResponse) []database.Product {
	var productList []database.Product
	for _, product := range searchResp.Data.AlgoliaSearch.Products {
		var sizes []database.Size

		for _, size := range product.SizeInt {
			sizes = append(sizes, database.Size{
				Size:  size,
				Value: "IN_STOCK",
			})
		}

		tempProduct := database.Product{
			Site:  "Rinascente",
			Sku:   product.URLKey,
			Name:  product.Name,
			Price: strconv.Itoa(product.Price),
			Image: product.MediaGalleryEntries[0].File,
			Color: product.MainColour,
			Sizes: sizes,
		}

		productList = append(productList, tempProduct)
	}

	return productList
}
