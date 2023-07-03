package adidas

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	database "github.com/TheAldriguzDev/Bonsai_solutions/packages/database"
	"github.com/TheAldriguzDev/Bonsai_solutions/packages/log"
	proxyManager "github.com/TheAldriguzDev/Bonsai_solutions/packages/proxy"
	filereader "github.com/TheAldriguzDev/Bonsai_solutions/packages/reader"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

var location = "Adidas"

func InitTls(proxy string) (adidasSession, error) {
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
		return adidasSession{}, err
	}

	return adidasSession{
		httpClient: client,
	}, nil
}

func InitAdidas() {
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

	module.Host = "https://www.adidas.it/"
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

		err = module.SearchProducts()
		if err != nil {
			log.Error(err.Error(), location)
			time.Sleep(7 * time.Second)
			continue
		}
		time.Sleep(6 * time.Second)
	}
}

func (adidas *adidasSession) SearchProducts() error {
	var searchData SearchData

	url := "https://api.3stripes.net/gw-api/v2/products/search?experiment_product_data=A&limit=24&page=1&query=Campus%2000"
	searchUrl := fmt.Sprintf("%s?%s", url, GetTimeStamp())

	req, err := http.NewRequest(http.MethodGet, searchUrl, nil)

	if err != nil {
		return err
	}

	req.Header = http.Header{
		//"cookie": {"akacd_api_3stripes=3864663956~rv%3D7~id%3D342369957e78c748be1595b33bf5a68b"},
		"Host":                {"api.3stripes.net"},
		"Cookie":              {"akacd_api_3stripes=3864632107~rv=18~id=a123da15ff5ed528e454280629888100"},
		"accept":              {"application/hal+json"},
		"x-device-info":       {"app/adidas; os/iOS; os-version/16.5; app-version/5.25.1; buildnumber/2023.6.13.12.5; type/iPhone14,2; fingerprint/D440FE93-EB37-49BB-A337-CA990667CD66"},
		"x-market":            {"IT"},
		"accept-version":      {"3.0.0"},
		"x-signature":         {"D492AF5CBDF6A5C73921D613CDAAC15301D25D625933B6B9023A1EFB5137A2D4"},
		"x-api-key":           {"m79qyapn2kbucuv96ednvh22"},
		"accept-language":     {"it-IT"},
		"user-agent":          {"adidas/2023.6.13.12.5 CFNetwork/1408.0.4 Darwin/22.5.0"},
		"x-app-info":          {"platform/iOS version/5.25.1"},
		"x-forter-mobile-uid": {"D440FE93-EB37-49BB-A337-CA990667CD66"},
	}

	res, err := adidas.httpClient.Do(req)

	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Error, bad statud code!")
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &searchData)
	if err != nil {
		return err
	}

	// Here we have the product search json
	err = adidas.parseSearch(searchData)
	if err != nil {
		return err
	}

	return nil
}

func (adidas *adidasSession) parseSearch(search SearchData) error {
	var productList []database.Product
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, product := range search.Embedded.Products {
		if strings.Contains(product.ProductFullName, "00") {
			wg.Add(1)

			// In this way, we can easily get all the product in the same time
			go adidas.getInfoUtils(product, &productList, &wg, &mu)
		}
	}

	wg.Wait()

	err := database.StartDataBase("./monitorData/adidas.json", productList)
	if err != nil {
		return err
	}

	return nil
}

func (adidas *adidasSession) getInfoUtils(product Product, productList *[]database.Product, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()

	tempProduct := database.Product{
		Site:  "Adidas",
		Sku:   product.ProductID,
		Name:  product.ProductFullName,
		Color: product.OriginalColor,
		Price: strconv.Itoa(product.OriginalPrice),
		Image: product.Embedded.ThumbnailImage.Href,
	}

	sizes, err := adidas.getProductInfo(tempProduct.Sku)
	if err != nil {
		log.Error(err.Error(), location)
		return
	}

	tempProduct.Sizes = sizes

	mu.Lock()
	*productList = append(*productList, tempProduct)
	mu.Unlock()
}

func (adidas *adidasSession) getProductInfo(sku string) ([]database.Size, error) {
	// we can change the ip for every request:
	proxie, err := proxyManager.RandomProxy(adidas.ProxyArr)
	proxie = strings.ReplaceAll(proxie, "\r", "")
	if err != nil {
		log.Error(err.Error(), location)
	} else {
		adidas.httpClient.SetProxy(proxie)
		adidas.Proxy = proxie
	}

	var productResp ProductStock
	var availableSizes []database.Size

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.adidas.it/api/products/%s/availability", sku), nil)
	if err != nil {
		return availableSizes, err
	}

	req.Header = http.Header{
		"authority":          {"www.adidas.it"},
		"accept":             {"*/*"},
		"accept-language":    {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
		"content-type":       {"application/json"},
		"glassversion":       {"c4e5ca0"},
		"referer":            {fmt.Sprintf("https://www.adidas.it/scarpe/%s.html", sku)},
		"sec-ch-ua":          {`"Not.A/Brand";v="8", "Chromium";v="114", "Google Chrome";v="114"`},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {`"macOS"`},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"same-origin"},
		"user-agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"x-instana-l":        {"1,correlationType=web;correlationId=b3946c491273721a"},
		"x-instana-s":        {"b3946c491273721a"},
		"x-instana-t":        {"b3946c491273721a"},
	}

	res, err := adidas.httpClient.Do(req)

	if err != nil {
		return availableSizes, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		return availableSizes, fmt.Errorf("Bad status code, stock!")
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return availableSizes, err
	}

	err = json.Unmarshal(body, &productResp)

	if err != nil {
		return availableSizes, err
	}

	// We can now parse the response

	for _, tempSize := range productResp.VariationList {
		if tempSize.Availability != 0 {
			availableSizes = append(availableSizes, database.Size{
				Size:  tempSize.Size,
				Value: strconv.Itoa(tempSize.Availability),
			})
		}
	}

	return availableSizes, nil
}
