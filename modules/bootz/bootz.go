package bootz

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/TheAldriguzDev/Bonsai_solutions/packages/database"
	"github.com/TheAldriguzDev/Bonsai_solutions/packages/log"
	proxyManager "github.com/TheAldriguzDev/Bonsai_solutions/packages/proxy"
	filereader "github.com/TheAldriguzDev/Bonsai_solutions/packages/reader"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

var location string = "Bootz"

func InitTls(proxy string) (bootzSession, error) {
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
		return bootzSession{}, err
	}

	return bootzSession{
		httpClient: client,
	}, nil
}

func InitBootz() {
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

	module.Host = "https://www.boozt.com/"
	module.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"
	module.Store_Region = "com"
	module.ProxyArr = proxyArr
	log.Success("Monitor initialised!", location)

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
			time.Sleep(7 * time.Second)
			continue
		}
		time.Sleep(6 * time.Second)
	}
}

func (tls *bootzSession) GetProducts() error {
	url := "https://www.boozt.com/eu/en/search?search_key=adidas%20campus"
	reqUrl := fmt.Sprintf("%s&ts=%s", url, GetTimeStamp())
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)

	if err != nil {
		return err
	}

	req.Header = http.Header{
		"authority":                 {"www.boozt.com"},
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"accept-language":           {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
		"cache-control":             {"max-age=0"},
		"sec-ch-ua":                 {"\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\""},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {"\"macOS\""},
		"sec-fetch-dest":            {"document"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-site":            {"none"},
		"sec-fetch-user":            {"?1"},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
	}

	res, err := tls.httpClient.Do(req)

	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		fmt.Println(res)
		return fmt.Errorf("Search bas status code! [%v]", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	err = tls.ParseResponse(doc)

	return nil
}

func (tls *bootzSession) ParseResponse(doc *goquery.Document) error {
	var productList []database.Product
	var wg sync.WaitGroup
	var mu sync.Mutex

	doc.Find("div[class*='palette-product-card--boozt']").Each(func(i int, elem *goquery.Selection) {
		wg.Add(1)

		url, _ := elem.Find("a[class*='palette-product-card-main-image__link']").Attr("href")
		price := elem.Find("span[class*='palette-product-card-price__price-tag']").Text()
		image, _ := elem.Find("img[class*='palette-image']").Attr("src")

		tempProducts := database.Product{
			Site:  "Bootz",
			Sku:   url,
			Name:  "Adidas Originals",
			Color: "None",
			Price: price,
			Image: image,
		}

		go tls.GetInfoUtils(tempProducts, &productList, &wg, &mu)
	})

	wg.Wait()

	err := database.StartDataBase("./monitorData/bootz.json", productList)
	if err != nil {
		return err
	}

	return nil
}

func (tls *bootzSession) GetInfoUtils(product database.Product, productList *[]database.Product, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()

	tempProduct := database.Product{
		Site:  "Bootz",
		Sku:   product.Sku,
		Name:  "Adidas Originals",
		Price: product.Price,
		Image: product.Image,
	}

	sizes, color, err := tls.GetSize(product.Sku)
	if err != nil {
		log.Error(err.Error(), location)
		return
	}

	tempProduct.Color = color
	tempProduct.Sizes = sizes

	mu.Lock()
	*productList = append(*productList, tempProduct)
	mu.Unlock()
}

func (tls *bootzSession) GetSize(path string) ([]database.Size, string, error) {
	var sizes []database.Size
	url := "https://www.boozt.com" + path
	reqUrl := fmt.Sprintf("%s?ts=%s", url, GetTimeStamp())
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)

	if err != nil {
		return sizes, "", err
	}

	req.Header = http.Header{
		"authority":                 {"www.boozt.com"},
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"accept-language":           {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
		"cache-control":             {"max-age=0"},
		"sec-ch-ua":                 {"\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\""},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {"\"macOS\""},
		"sec-fetch-dest":            {"document"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-site":            {"none"},
		"sec-fetch-user":            {"?1"},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
	}

	res, err := tls.httpClient.Do(req)

	if err != nil {
		return sizes, "", err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		fmt.Println(res)
		return sizes, "", fmt.Errorf("Search bas status code! [%v]", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return sizes, "", err
	}

	color := doc.Find("span[class*='pp-colour-selector__label-colour-name']").Text()
	doc.Find("button[class*='size-btn']").Each(func(i int, elem *goquery.Selection) {
		sizes = append(sizes, database.Size{
			Size:  strings.TrimSpace(elem.Text()),
			Value: "IN_STOCK",
		})
	})

	return sizes, color, nil
}
