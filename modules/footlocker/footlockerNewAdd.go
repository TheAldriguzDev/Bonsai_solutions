package footlocker

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"

	"github.com/TheAldriguzDev/Bonsai_solutions/packages/database"
	"github.com/TheAldriguzDev/Bonsai_solutions/packages/log"
	proxyManager "github.com/TheAldriguzDev/Bonsai_solutions/packages/proxy"
	filereader "github.com/TheAldriguzDev/Bonsai_solutions/packages/reader"
	tls_client "github.com/bogdanfinn/tls-client"
)

var location string = "Footlocker"

func InitTls(proxy string) (footlockerSession, error) {
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
		return footlockerSession{}, err
	}

	return footlockerSession{
		httpClient: client,
	}, nil
}

func InitFootlocker() {
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

	module.Host = "https://www.footlocker.it"
	module.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"
	module.Store_Region = "it"
	module.ProxyArr = proxyArr
	log.Success("Monitor initialised!", location)

	var cidStep bool = false
	var cid string = ""

	for !cidStep {
		log.General("Generating cid...", location)
		cid, err = module.GenerateCid()

		if err != nil {
			log.Error(err.Error(), location)
			proxie, err := proxyManager.RandomProxy(proxyArr)
			proxie = strings.ReplaceAll(proxie, "\r", "")
			if err != nil {
				log.Error(err.Error(), location)
			} else {
				module.httpClient.SetProxy(proxie)
				module.Proxy = proxie
			}
			time.Sleep(10 * time.Second)
			continue
		}

		if cid == "" {
			log.Error("Error while getting the cid!", location)
			proxie, err := proxyManager.RandomProxy(proxyArr)
			proxie = strings.ReplaceAll(proxie, "\r", "")
			if err != nil {
				log.Error(err.Error(), location)
			} else {
				module.httpClient.SetProxy(proxie)
				module.Proxy = proxie
			}
			time.Sleep(10 * time.Second)
			continue
		}

		cidStep = true
	}
	module.Cid = cid

	var run bool = true

	for run {
		log.General("Monitoring...", location)
		proxie, err := proxyManager.RandomProxy(proxyArr)
		proxie = strings.ReplaceAll(proxie, "\r", "")
		if err != nil {
			log.Error(err.Error(), location)
		} else {
			module.httpClient.SetProxy(proxie)
			module.Proxy = proxie
		}

		err = module.GetProducts()
		if err != nil {
			log.Error(err.Error(), location)
			time.Sleep(15 * time.Second)
			continue
		}
		time.Sleep(15 * time.Second)
	}
}

func (ftl *footlockerSession) GetProducts() error {
	var searchData SearchResponse
	searchUrl := fmt.Sprintf("https://www.footlocker.it/zgw/search-core/products/v2/search?query=campus&pageSize=48&pageType=search&timestamp=%s", GetTimeStamp())

	req, err := http.NewRequest(http.MethodGet, searchUrl, nil)

	if err != nil {
		return err
	}

	//newCookie := ftl.GetDataDome()

	req.Header = http.Header{
		"Host":             {ftl.Host},
		"user-agent":       {ftl.UserAgent},
		"Accept":           {"application/json, text/javascript, */*; q=0.01"},
		"Accept-Language":  {"de-DE,de;q=0.9"},
		"Accept-Encoding":  {"gzip, deflate"},
		"Content-Type":     {"application/x-www-form-urlencoded"},
		"X-Requested-With": {"XMLHttpRequest"},
		"Origin":           {ftl.Host},
		"Referer":          {searchUrl},
		"Sec-Fetch-Dest":   {"empty"},
		"Sec-Fetch-Mode":   {"cors"},
		"Sec-Fetch-Site":   {"same-origin"},
		"Pragma":           {"no-cache"},
		"Cache-Control":    {"no-cache"},
		"Te":               {"trailers"},
		//"cookie":           {"datadome=" + newCookie},
	}

	// cookie := &http.Cookie{
	// 	Name:   "datadome",
	// 	Value:  newCookie,
	// 	MaxAge: 604800,
	// }

	// Modify a cookie value
	cookies := ftl.httpClient.GetCookies(req.URL)
	for _, cookie := range cookies {
		if cookie.Name == "datadome" {
			fmt.Println(cookie.Value)
			break
		}
	}

	// ftl.httpClient.SetCookies(req.URL, cookies)
	// req.AddCookie(cookie)

	res, err := ftl.httpClient.Do(req)

	if err != nil {
		return err
	}

	if res != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Search bad status code!")
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &searchData)

	if err != nil {
		return err
	}

	err = ftl.SearchParser(searchData)

	if err != nil {
		return err
	}

	return nil
}

func (ftl *footlockerSession) SearchParser(searchRes SearchResponse) error {
	var availableSizes []database.Size
	var productList []database.Product
	for _, product := range searchRes.Products {
		tempProduct := database.Product{
			Site:  "Footlocker",
			Sku:   product.SKU,
			Name:  product.Name,
			Color: "None",
			Price: product.Price.FormattedValue,
			Image: fmt.Sprintf("https://images.footlocker.com/is/image/FLEU/%s_01?wid=520&hei=520&fmt=png-alpha", product.SKU),
			Sizes: availableSizes,
		}

		productList = append(productList, tempProduct)
	}

	err := database.StartDataBase("./monitorData/footlocker.json", productList)
	if err != nil {
		return err
	}

	return nil
}

// func (ftl *footlockerSession) SearchParser(searchRes SearchResponse) error {
// 	var productList []database.Product
// 	var wg sync.WaitGroup
// 	var mu sync.Mutex

// 	for _, product := range searchRes.Products {
// 		wg.Add(1)

// 		ftl.GetInfoUtils(product, &productList, &wg, &mu)
// 	}

// 	wg.Wait()

// 	err := database.StartDataBase("./monitorData/footlocker.json", productList)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (ftl *footlockerSession) GetInfoUtils(product Product, productList *[]database.Product, wg *sync.WaitGroup, mu *sync.Mutex) {
// 	defer wg.Done()

// 	tempProduct := database.Product{
// 		Site:  "Footlocker",
// 		Sku:   product.SKU,
// 		Name:  product.Name,
// 		Color: "None",
// 		Price: product.Price.FormattedValue,
// 		Image: fmt.Sprintf("https://images.footlocker.com/is/image/FLEU/%s_01?wid=520&hei=520&fmt=png-alpha", product.SKU),
// 	}

// 	sizes, color, err := ftl.GetStock(product)
// 	if err != nil {
// 		log.Error(err.Error(), location)
// 		return
// 	}

// 	tempProduct.Sizes = sizes
// 	tempProduct.Color = color
// 	mu.Lock()
// 	*productList = append(*productList, tempProduct)
// 	mu.Unlock()
// }

// func (ftl *footlockerSession) GetStock(product Product) ([]database.Size, string, error) {
// 	proxie, err := proxyManager.RandomProxy(ftl.ProxyArr)
// 	proxie = strings.ReplaceAll(proxie, "\r", "")
// 	if err != nil {
// 		log.Error(err.Error(), location)
// 	} else {
// 		ftl.httpClient.SetProxy(proxie)
// 		ftl.Proxy = proxie
// 	}

// 	var productData ProductData
// 	var availableSizes []database.Size

// 	refer := fmt.Sprintf("https://www.footlocker.it/it/product/~/%s.html", product.SKU)
// 	productUrl := fmt.Sprintf("https://www.footlocker.it/zgw/product-core/v1/pdp/sku/%s?%s", product.SKU, GetTimeStamp())

// 	req, err := http.NewRequest(http.MethodGet, productUrl, nil)

// 	//ddCookie := ftl.GetDataDome(refer)

// 	req.Header = http.Header{
// 		"authority":                   {"www.footlocker.it"},
// 		"accept":                      {"application/json"},
// 		"accept-language":             {"it-IT,it;q=0.8"},
// 		"cookie":                      {"datadome=7Fq53HCbsVW0UnrHldZ2341XX8kCDbVmSK-d~6I1Xlub-h2o8Ui6YFXnCpPphQhSQtzKqjJz7IsYcDDRkeyO4gPqIVtYsaN5b6ySQufIIhOIQ3~bhbuuP9RsxfuKa27o"},
// 		"referer":                     {refer},
// 		"sec-ch-device-memory":        {"8"},
// 		"sec-ch-ua":                   {"\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\""},
// 		"sec-ch-ua-arch":              {"\"x86\""},
// 		"sec-ch-ua-full-version-list": {"\"Google Chrome\";v=\"113.0.5672.126\", \"Chromium\";v=\"113.0.5672.126\", \"Not-A.Brand\";v=\"24.0.0.0\""},
// 		"sec-ch-ua-mobile":            {"?0"},
// 		"sec-ch-ua-model":             {"\"\""},
// 		"sec-ch-ua-platform":          {"\"macOS\""},
// 		"sec-fetch-dest":              {"empty"},
// 		"sec-fetch-mode":              {"cors"},
// 		"sec-fetch-site":              {"same-origin"},
// 		"user-agent":                  {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
// 		"x-api-lang":                  {"it-IT"},
// 		"x-fl-request-id":             {GenerateUUID()},
// 	}

// 	// Modify a cookie value
// 	cookies := ftl.httpClient.GetCookies(req.URL)
// 	for _, cookie := range cookies {
// 		if cookie.Name == "datadome" {
// 			cookie.Value = "7Fq53HCbsVW0UnrHldZ2341XX8kCDbVmSK-d~6I1Xlub-h2o8Ui6YFXnCpPphQhSQtzKqjJz7IsYcDDRkeyO4gPqIVtYsaN5b6ySQufIIhOIQ3~bhbuuP9RsxfuKa27o"
// 			break
// 		}
// 	}

// 	ftl.httpClient.SetCookies(req.URL, cookies)

// 	fmt.Println(ftl.httpClient.GetCookies(req.URL))

// 	res, err := ftl.httpClient.Do(req)

// 	if err != nil {
// 		return availableSizes, "", err
// 	}

// 	if res.Body != nil {
// 		defer res.Body.Close()
// 	}

// 	if res.StatusCode != 200 {
// 		return availableSizes, "", fmt.Errorf("Get product info bad status code! [%v]", res.StatusCode)
// 	}

// 	body, err := io.ReadAll(res.Body)

// 	if err != nil {
// 		return availableSizes, "", err
// 	}

// 	err = json.Unmarshal(body, &productData)

// 	if err != nil {
// 		return availableSizes, "", err
// 	}

// 	// Now we can parse the response:

// 	for _, size := range productData.Sizes {
// 		if size.Active != false {
// 			availableSizes = append(availableSizes, database.Size{
// 				Size:  size.Size,
// 				Value: "Available",
// 			})
// 		}
// 	}

// 	color := productData.Style.Color

// 	return availableSizes, color, nil
// }

func (ftl *footlockerSession) GenerateCid() (string, error) {
	homePageUrl := fmt.Sprintf("https://www.footlocker.%s/?%s", ftl.Store_Region, GetTimeStamp())

	req, err := http.NewRequest(http.MethodGet, homePageUrl, nil)

	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"Host":             {ftl.Host},
		"user-agent":       {ftl.UserAgent},
		"Accept":           {"application/json, text/javascript, */*; q=0.01"},
		"Accept-Language":  {"de-DE,de;q=0.9"},
		"Accept-Encoding":  {"gzip, deflate"},
		"Content-Type":     {"application/x-www-form-urlencoded"},
		"X-Requested-With": {"XMLHttpRequest"},
		"Origin":           {ftl.Host},
		"Referer":          {homePageUrl},
		"Sec-Fetch-Dest":   {"empty"},
		"Sec-Fetch-Mode":   {"cors"},
		"Sec-Fetch-Site":   {"same-origin"},
		"Pragma":           {"no-cache"},
		"Cache-Control":    {"no-cache"},
		"Te":               {"trailers"},
	}

	res, err := ftl.httpClient.Do(req)

	if err != nil {
		return "", err
	}

	if res != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		fmt.Println(res.StatusCode)
		return "", fmt.Errorf("Get home page bad status code! [%v]", res.StatusCode)
	}

	if !strings.Contains(res.Header.Get("Set-Cookie"), "datadome") {
		return "", fmt.Errorf("No datadome cid has been found!")
	}

	ddCookie := strings.Split(res.Header.Get("Set-Cookie"), ";")[0]
	ddCookie = strings.Split(ddCookie, "=")[1]

	return ddCookie, nil
}
