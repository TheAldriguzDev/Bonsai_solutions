package ftlkiller

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
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
	log.Success("We ready to fuck footlocker!", location)

	// We initialise the monitor with the cid

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

	// Monitor cycle
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
			time.Sleep(15 * time.Second)
			continue
		}
		time.Sleep(10 * time.Second)
	}
}

// Function to generate the cid, a valid dd cookie
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

func (ftl *footlockerSession) GetProducts() error {
	var searchData SearchResponse
	searchUrl := fmt.Sprintf("https://www.footlocker.it/zgw/search-core/products/v2/search?query=campus&pageSize=48&pageType=search&timestamp=%s", GetTimeStamp())

	req, err := http.NewRequest(http.MethodGet, searchUrl, nil)

	if err != nil {
		return err
	}

	newCookie := ftl.GetDataDome()

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
		"cookie":           {"datadome=" + newCookie},
	}

	cookie := &http.Cookie{
		Name:   "datadome",
		Value:  newCookie,
		MaxAge: 604800,
	}

	// Modify a cookie value
	cookies := ftl.httpClient.GetCookies(req.URL)
	for _, cookie := range cookies {
		if cookie.Name == "datadome" {
			cookie.Value = newCookie
			break
		}
	}

	ftl.httpClient.SetCookies(req.URL, cookies)
	req.AddCookie(cookie)

	res, err := ftl.httpClient.Do(req)

	if err != nil {
		return err
	}

	if res != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Search bad status code! %v", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &searchData)

	if err != nil {
		return err
	}

	// We need to parse the response...
	err = ftl.ParseSearch(searchData)

	if err != nil {
		return err
	}

	return nil
}

func (ftl *footlockerSession) ParseSearch(sData SearchResponse) error {
	var productList []database.Product
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, product := range sData.Products {
		if strings.Contains(strings.ToLower(product.Name), "campus") {
			if checkSku(product.SKU) {
				wg.Add(1)

				// Get the product and parse the response
				go ftl.GetProductInfo(product, &productList, &wg, &mu)
			}
			// Bad sku are not included
		}
	}

	wg.Wait()

	err := database.StartDataBase("./monitorData/ftlKiller.json", productList)
	if err != nil {
		return err
	}

	return nil
}

func checkSku(sku string) bool {
	badSku := []string{
		"314311015904",
		"314311294004",
		"314310886404",
		"314311003504",
	}
	for _, forbiddenSku := range badSku {
		if sku == forbiddenSku {
			return false
		}
	}
	return true
}

func (ftl footlockerSession) GetProductInfo(product Product, productList *[]database.Product, wg *sync.WaitGroup, mu *sync.Mutex) error {
	defer wg.Done()
	// Ruotiamo i proxy per evitare ip ban

	proxie, err := proxyManager.RandomProxy(ftl.ProxyArr)
	proxie = strings.ReplaceAll(proxie, "\r", "")
	proxyData := strings.Split(strings.TrimSpace(proxie), ":")
	proxyUrl := fmt.Sprintf("http://%s:%s@%s:%s", proxyData[2], proxyData[3], proxyData[0], proxyData[1])
	if err != nil {
		log.Error(err.Error(), location)
	} else {
		err = ftl.httpClient.SetProxy(proxyUrl)
		if err != nil {
			log.Error(err.Error(), location)
		}
		ftl.Proxy = proxie
	}

	refer := fmt.Sprintf("https://www.footlocker.it/it/product/~/%s.html", product.SKU)
	productUrl := fmt.Sprintf("https://www.footlocker.it/zgw/product-core/v1/pdp/sku/%s?ts=%s", product.SKU, GetTimeStamp())

	req, err := http.NewRequest(http.MethodGet, productUrl, nil)

	req.Header = http.Header{
		"authority":                   {"www.footlocker.it"},
		"accept":                      {"application/json"},
		"accept-language":             {"it-IT,it;q=0.8"},
		"cookie":                      {"----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------"},
		"referer":                     {refer},
		"sec-ch-device-memory":        {"8"},
		"sec-ch-ua":                   {"\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\""},
		"sec-ch-ua-arch":              {"\"x86\""},
		"sec-ch-ua-full-version-list": {"\"Google Chrome\";v=\"113.0.5672.126\", \"Chromium\";v=\"113.0.5672.126\", \"Not-A.Brand\";v=\"24.0.0.0\""},
		"sec-ch-ua-mobile":            {"?0"},
		"sec-ch-ua-model":             {"\"\""},
		"sec-ch-ua-platform":          {"\"macOS\""},
		"sec-fetch-dest":              {"empty"},
		"sec-fetch-mode":              {"cors"},
		"sec-fetch-site":              {"same-origin"},
		"user-agent":                  {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
		"x-api-lang":                  {"it-IT"},
		"x-fl-request-id":             {GenerateUUID()},
	}

	res, err := ftl.httpClient.Do(req)

	if err != nil {
		log.Error(err.Error(), location)
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error(err.Error(), location)
		return err
	}

	var productData ProductData

	err = json.Unmarshal(body, &productData)

	if err != nil {
		log.Error(err.Error(), location)
		return err
	}

	// Now we need to parse the product

	sizes := []database.Size{}

	for _, size := range productData.Sizes {
		if size.Active == true {
			sizes = append(sizes, database.Size{
				Size:  size.Size,
				Value: "IN_STOCK",
			})
		}
	}

	tempProduct := database.Product{
		Site:  "FTLKiller",
		Sku:   product.SKU,
		Name:  product.Name,
		Image: fmt.Sprintf("https://images.footlocker.com/is/image/FLEU/%s_01?wid=520&hei=520&fmt=png-alpha", product.SKU),
		Sizes: sizes,
		Color: productData.Style.Color,
		Price: strconv.FormatFloat(product.Price.Value, 'f', -1, 64),
	}

	mu.Lock()
	*productList = append(*productList, tempProduct)
	mu.Unlock()

	return nil
}
