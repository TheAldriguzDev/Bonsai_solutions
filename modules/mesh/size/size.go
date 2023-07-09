package size

import (
	"encoding/json"
	"fmt"
	"io"
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

var location string = "SIZE IT"

func InitTls(proxy string) (meshSession, error) {
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
		return meshSession{}, err
	}

	return meshSession{
		httpClient: client,
	}, nil
}

func InitSize() {
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

	module.Host = "size.it"
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
			time.Sleep(15 * time.Second)
			continue
		}
		time.Sleep(15 * time.Second)
	}
}

func (mesh *meshSession) SearchProducts() error {
	req, err := http.NewRequest(http.MethodGet, "https://prod.jdgroupmesh.cloud/stores/sizeit/products/search?max=24&customerLoyalty=1&from=0&channel=android-app-phone&query=campus", nil)

	if err != nil {
		return err
	}

	req.Header = http.Header{
		"accept":                {"application/json"},
		"accept-encoding":       {"gzip"},
		"cache-control":         {"no-cache"},
		"connection":            {"Keep-Alive"},
		"host":                  {"prod.jdgroupmesh.cloud"},
		"mesh-commerce-channel": {"android-app-phone"},
		"newrelic":              {"eyJ2IjpbMCwyXSwiZCI6eyJkLnR5IjoiTW9iaWxlIiwiZC5hYyI6IjEwNDEzNTUiLCJkLmFwIjoiMTI0MTg4NDUxIiwiZC50ciI6IjFhYWZkZWZjYTFiOTQ1NzA5YmU0MzI5NDJjOTk4MTI5IiwiZC5pZCI6ImZlNjAyYzlmNmJlNjRjNzYiLCJkLnRpIjoxNjg4MjQ0MzE1Njk1fX0="},
		"traceparent":           {"00-1aafdefca1b945709be432942c998129-fe602c9f6be64c76-00"},
		"tracestate":            {fmt.Sprintf("@nr=0-2-1041355-124188451-fe602c9f6be64c76----%s", GenTracestate())},
		"user-agent":            {"sizeit/4.2.2.10964 (android-app-phone; Android 7.1.2; Build/dream2ltexx-user 7.1.2 NRD90M 700210421 release-keys)"},
		"x-api-key":             {"78F4021E82C64DBA858F49E210192EE0"},
		"x-newrelic-id":         {"VQYDUFVWDRABVlVTDwkDU1U="},
		"x-request-auth":        {fmt.Sprintf("Hawk id=\"c766afaefd\", mac=\"+N7AZ/a5uDNmivSHP3bF2vX5L2c+2bLMR72LgVtUQmc=\", ts=\"%s\", nonce=\"9d93c5", GenHawk())},
		"x-trace":               {"d052b189-bf3f-4ed2-be01-eca818a3cb03"},
		"cookie":                {"----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------"},
	}

	res, err := mesh.httpClient.Do(req)

	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Search product bad status code! [%v]", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	var searchRes SearchRes

	err = json.Unmarshal(body, &searchRes)

	if err != nil {
		return err
	}

	err = mesh.ParseSearch(searchRes)

	if err != nil {
		return err
	}

	return nil
}

func (mesh *meshSession) ParseSearch(searchData SearchRes) error {
	var productList []database.Product
	var wg sync.WaitGroup
	//var mu sync.Mutex

	for _, product := range searchData.Products {

		if strings.Contains(product.Name, "00") || strings.Contains(strings.ToLower(product.Name), "junior") {
			wg.Add(1)
			//go mesh.GetProduct(product, &productList, &wg, &mu)
		}
	}

	wg.Wait()

	err := database.StartDataBase("./monitorData/size.json", productList)
	if err != nil {
		return err
	}

	return nil
}

// func (mesh *meshSession) GetProduct(product Product, productList *[]database.Product, wg *sync.WaitGroup, mu *sync.Mutex) error {
// 	defer wg.Done()

// 	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://prod.jdgroupmesh.cloud/stores/sizeit/products/%s?api_key=78F4021E82C64DBA858F49E210192EE0", product.SKU), nil)

// 	if err != nil {
// 		return err
// 	}

// 	req.Header = http.Header{
// 		"accept":                {"application/json"},
// 		"accept-encoding":       {"gzip"},
// 		"cache-control":         {"no-cache"},
// 		"connection":            {"Keep-Alive"},
// 		"host":                  {"prod.jdgroupmesh.cloud"},
// 		"mesh-commerce-channel": {"android-app-phone"},
// 		"newrelic":              {"eyJ2IjpbMCwyXSwiZCI6eyJkLnR5IjoiTW9iaWxlIiwiZC5hYyI6IjEwNDEzNTUiLCJkLmFwIjoiMTI0MTg4NDUxIiwiZC50ciI6IjFhYWZkZWZjYTFiOTQ1NzA5YmU0MzI5NDJjOTk4MTI5IiwiZC5pZCI6ImZlNjAyYzlmNmJlNjRjNzYiLCJkLnRpIjoxNjg4MjQ0MzE1Njk1fX0="},
// 		"traceparent":           {"00-1aafdefca1b945709be432942c998129-fe602c9f6be64c76-00"},
// 		"tracestate":            {fmt.Sprintf("@nr=0-2-1041355-124188451-fe602c9f6be64c76----%s", GenTracestate())},
// 		"user-agent":            {"sizeit/4.2.2.10964 (android-app-phone; Android 7.1.2; Build/dream2ltexx-user 7.1.2 NRD90M 700210421 release-keys)"},
// 		"x-api-key":             {"78F4021E82C64DBA858F49E210192EE0"},
// 		"x-newrelic-id":         {"VQYDUFVWDRABVlVTDwkDU1U="},
// 		"x-request-auth":        {fmt.Sprintf("Hawk id=\"c766afaefd\", mac=\"+N7AZ/a5uDNmivSHP3bF2vX5L2c+2bLMR72LgVtUQmc=\", ts=\"%s\", nonce=\"9d93c5", GenHawk())},
// 		"x-trace":               {"d052b189-bf3f-4ed2-be01-eca818a3cb03"},
// 		"cookie":                {"----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------"},
// 	}

// 	res, err := mesh.httpClient.Do(req)

// 	if err != nil {
// 		log.Error(err.Error(), location)
// 		return err
// 	}

// 	if res.Body != nil {
// 		defer res.Body.Close()
// 	}

// 	if res.StatusCode != 200 {
// 		log.Error("Bad status code!", location)
// 		return fmt.Errorf("Get product bad statuc code! [%v]", res.StatusCode)
// 	}

// 	body, err := io.ReadAll(res.Body)

// 	if err != nil {
// 		log.Error(err.Error(), location)
// 		return err
// 	}

// 	var productRes ProductRes

// 	err = json.Unmarshal(body, &productRes)

// 	if err != nil {
// 		log.Error(err.Error(), location)
// 		return err
// 	}

// 	// Now we can parse and create the temp product

// 	var sizes []database.Size

// 	for _, size := range productRes.SortedOptions {
// 		if size.Product.StockStatus == "IN STOCK" {
// 			sizes = append(sizes, database.Size{
// 				Size:  size.Name,
// 				Value: strings.Split(size.Product.SKU, ".")[0],
// 			})
// 		}
// 	}

// 	tempProduct := database.Product{
// 		Site:  "SIZE",
// 		Sku:   product.SKU,
// 		Name:  product.ProductName,
// 		Price: product.Price.Amount,
// 		Image: product.MainImage,
// 		Color: product.Colour,
// 		Sizes: sizes,
// 	}

// 	mu.Lock()
// 	*productList = append(*productList, tempProduct)
// 	mu.Unlock()
// 	return nil
// }

func GenHawk() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

func GenTracestate() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
}
