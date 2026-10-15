package ftlRestock

import (
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
	"github.com/TheAldriguzDev/Bonsai_solutions/solvers/datadome"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

var location string = "Footlocker Restock"

func generateRandomNumber(min, max int) int {
	// Genera un seme univoco utilizzando l'ora corrente
	rand.Seed(time.Now().UnixNano())

	// Genera un numero casuale compreso tra min e max
	number := rand.Intn(max-min+1) + min

	return number
}

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

func InitRestock() {
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

	// Idea
	// Mettere un contatore, che se arriva a 3 403 --> riparte tutto da qua, dal get cid ricreando anche un nuovo client tls cosi da avere una nuova sessione

	var cidStep bool = false
	var cid string = ""

	for !cidStep {
		log.General("Generating cid...", location)
		cid, err = module.GenerateCid()

		if err != nil {
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
			time.Sleep(10 * time.Second)
			continue
		}

		if cid == "" {
			log.Error("Error while getting the cid!", location)
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
			time.Sleep(10 * time.Second)
			continue
		}

		cidStep = true
	}
	module.Cid = cid

	// Start the process here

	// var run bool = true

	var run bool = true

	for run {
		log.General("Monitoring...", location)

		// we have to rotate the proxy for every request
		err = module.RotateProxy()
		if err != nil {
			log.Error(err.Error(), location)
		}

		// Now we can monitor the new add
		err = module.ReadFile()
		if err != nil {
			log.Error(err.Error(), location)
			time.Sleep(time.Duration(generateRandomNumber(15, 30)) * time.Second)
			continue
		}

		time.Sleep(time.Duration(generateRandomNumber(50, 80)) * time.Second)
	}

}

func (ftl *footlockerSession) RotateProxy() error {
	proxie, err := proxyManager.RandomProxy(ftl.ProxyArr)
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
		err = ftl.httpClient.SetProxy(proxyUrl)
		if err != nil {
			return err
		}
		ftl.Proxy = proxie
	}
	return nil
}

type ProductFile struct {
	SKU string `json:"sku"`
}

type File struct {
	Products []ProductFile `json:"products"`
}

func (ftl *footlockerSession) ReadFile() error {
	var productList []database.Product
	var wg sync.WaitGroup
	var mu sync.Mutex

	data, err := ioutil.ReadFile("./monitorData/ftlNewAdd.json")
	if err != nil {
		return err
	}

	// Decodifica il contenuto JSON
	var jsonData File
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return err
	}

	if len(jsonData.Products) == 0 {
		return fmt.Errorf("No product found to be monitored!")
	}

	for _, product := range jsonData.Products {
		wg.Add(1)

		go ftl.GetProducts(product.SKU, &productList, &wg, &mu)
	}

	wg.Wait()

	err = database.StartDataBase("./monitorData/ftlRestock.json", productList)
	if err != nil {
		return err
	}

	return nil
}

func (ftl footlockerSession) GetProducts(sku string, productList *[]database.Product, wg *sync.WaitGroup, mu *sync.Mutex) error {
	defer wg.Done()

	// We rotate the proxy for every product we need to scrape
	proxie, err := proxyManager.RandomProxy(ftl.ProxyArr)
	if err != nil {
		log.Error(err.Error(), location)
		return err
	}
	proxie = strings.ReplaceAll(proxie, "\r", "")
	proxyData := strings.Split(strings.TrimSpace(proxie), ":")
	if len(proxyData) < 4 {
		log.Error("The proxy format is not correct!", location)
		return fmt.Errorf("The proxy format is not correct! %s", proxie)
	}
	proxyUrl := fmt.Sprintf("http://%s:%s@%s:%s", proxyData[2], proxyData[3], proxyData[0], proxyData[1])
	if err != nil {
		log.Error(err.Error(), location)
		return err
	} else {
		err = ftl.httpClient.SetProxy(proxyUrl)
		if err != nil {
			log.Error(err.Error(), location)
			return err
		}
		ftl.Proxy = proxie
	}
	// Here we get the search ep with keyword to get the products loaded

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.footlocker.it/apigate/products/search?currentPage=0&pageSize=24&query=%s&ts=%s", sku, GetTimeStamp()), nil)

	if err != nil {
		log.Error(err.Error(), location)
		return err
	}

	newCookie := ftl.GetDataDome()
	ftl.Cid = newCookie

	req.Header = http.Header{
		"Host":                   {"www.footlocker.it"},
		"x-newrelic-id":          {"VgAPVVdRDRAIV1lXAAEEXlQ="},
		"tracestate":             {"@nr=0-2-2684125-838570380-409878c121a2c448--0--" + GenTracestate()},
		"x-fl-request-id":        {GenerateUUID()},
		"accept":                 {"application/json"},
		"x-flapi-session-id":     {"avmpbbf4mkxxw55fin8hxeov.fzcxwefapipdb428880"},
		"x-fl-app-version":       {"5.9.1"},
		"x-flapi-api-identifier": {"921B2b33cAfba5WWcb0bc32d5ix89c6b0f614"},
		"newrelic":               {"ewoiZCI6IHsKImFjIjogIjI2ODQxMjUiLAoiYXAiOiAiODM4NTcwMzgwIiwKImlkIjogIjQwOTg3OGMxMjFhMmM0NDgiLAoidGkiOiAxNjcyNjY3MDc5MDYxLAoidHIiOiAiYWQ0OGMyMWJkYzhlYjY2M2ZmYzc1MjRmNzNlZTZjM2YiLAoidHkiOiAiTW9iaWxlIgp9LAoidiI6IFsKMCwKMgpdCn0="},
		"x-fl-device-id":         {"A53993E0-4E99-4765-A5DC-08D5ED91A58C"},
		"accept-language":        {"it-IT,it;q=0.8"},
		"x-api-key":              {"R9wAA5tmfidnRM7DG7qGZ6myVHnZb3QU"},
		"user-agent":             {"FLEU/CFNetwork/Darwin"},
		"x-api-lang":             {"it-IT"},
		"traceparent":            {"00-ad48c21bdc8eb663ffc7524f73ee6c3f-409878c121a2c448-00"},
		"cookie":                 {newCookie},
	}

	if hasDatadomeCookie(req) {
		// We need to modify the existing datadome cookie
		cookies := ftl.httpClient.GetCookies(req.URL)
		for _, cookie := range cookies {
			if cookie.Name == "datadome" {
				cookie.Value = newCookie
				ftl.httpClient.SetCookies(req.URL, cookies)
				break
			}
		}
	} else {
		// No datadome cookie found, we need to set it!
		cookie := &http.Cookie{
			Name:   "datadome",
			Value:  newCookie,
			MaxAge: 604800,
		}

		ftl.httpClient.SetCookies(req.URL, []*http.Cookie{cookie}) // Aggiunta del nuovo cookie
		req.Header.Set("Cookie", req.Header.Get("Cookie"))
	}

	res, err := ftl.httpClient.Do(req)

	if err != nil {
		log.Error(err.Error(), location)
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode == 403 {
		log.Error("Datadome block!", location)
		return fmt.Errorf("Datadome block! [%v]", res.StatusCode)
	}

	if res.StatusCode == 431 {
		log.Error("Product not loaded!", location)
		return fmt.Errorf("Product not loaded! [%v]", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Error(err.Error(), location)
		return err
	}

	var productResponse ProductData

	err = json.Unmarshal(body, &productResponse)

	if err != nil {
		log.Error(err.Error(), location)
		return err
	}

	// Here we need to parse the product search response

	tempProduct, err := ftl.ParseSearch(productResponse, sku)
	if err != nil {
		log.Error(err.Error(), location)
		return err
	}

	mu.Lock()
	*productList = append(*productList, tempProduct)
	mu.Unlock()

	return nil
}

func (ftl *footlockerSession) ParseSearch(searchData ProductData, sku string) (database.Product, error) {

	if len(searchData.Products) != 1 {
		log.Error("More than one product found!", location)
	}

	product := searchData.Products[0]
	var color string = "None"

	sizes := []database.Size{}

	for _, facet := range searchData.Facets {
		if facet.Name == "Taglia" {
			for _, tempSize := range facet.Values {
				sizes = append(sizes, database.Size{
					Size:  tempSize.Name,
					Value: "IN_STOCK",
				})
			}
		}

		if facet.Name == "Colore" {
			color = facet.Values[0].Name
		}
	}

	tempProduct := database.Product{
		Site:  "Footlocker Restock",
		Sku:   sku,
		Name:  product.Name,
		Price: product.Price.FormattedValue,
		Image: fmt.Sprintf("https://images.footlocker.com/is/image/FLEU/%s_01?wid=520&hei=520&fmt=png-alpha", sku),
		Sizes: sizes,
		Color: color,
	}

	return tempProduct, nil
}

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
		return "", fmt.Errorf("Get home page bad status code! [%v]", res.StatusCode)
	}

	if !strings.Contains(res.Header.Get("Set-Cookie"), "datadome") {
		return "", fmt.Errorf("No datadome cid has been found!")
	}

	ddCookie := strings.Split(res.Header.Get("Set-Cookie"), ";")[0]
	ddCookie = strings.Split(ddCookie, "=")[1]

	return ddCookie, nil
}

func (ftl *footlockerSession) GetDataDome() string {
	var ddStep bool = false
	var ddRes datadome.DDResponse

	ddData := datadome.GenerateData{
		DataDomeKey:     "A55FBF4311ED6F1BF9911EB71931D5",
		DataDomeVersion: "4.12.1",
		Cid:             ftl.Cid,
		Refer:           "https://www.footlocker.it/",
		UserAgent:       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
		Proxy:           ftl.Proxy,
	}

	for !ddStep {
		ddResTemp, err := datadome.GetCookie(ddData)

		if err != nil {
			log.Error(err.Error(), location)
			time.Sleep(5 * time.Second)
			continue
		}

		ddRes = ddResTemp
		ddStep = true
	}

	newCookie := strings.Split(strings.Split(ddRes.Cookie, ";")[0], "=")[1]
	return newCookie
}

func GenTracestate() string {
	return fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))
}
