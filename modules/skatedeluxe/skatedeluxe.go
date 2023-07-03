package skatedeluxe

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	http "github.com/bogdanfinn/fhttp"

	"github.com/TheAldriguzDev/Bonsai_solutions/packages/database"
	"github.com/TheAldriguzDev/Bonsai_solutions/packages/log"
	proxyManager "github.com/TheAldriguzDev/Bonsai_solutions/packages/proxy"
	filereader "github.com/TheAldriguzDev/Bonsai_solutions/packages/reader"

	tls_client "github.com/bogdanfinn/tls-client"
)

var location string = "SkateDeluxe"

func InitTls(proxy string) (skateSession, error) {
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
		return skateSession{}, err
	}

	return skateSession{
		httpClient: client,
	}, nil
}

func InitSkate() {
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

	module.Host = "https://www.skatedeluxe.com/"
	module.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"
	module.Store_Region = "it"
	module.ProxyArr = proxyArr
	log.Success("Monitor initialised!", location)

	var tokenStatus bool = false
	var authToken string = ""

	for !tokenStatus {
		authToken, err = module.GetToken()

		if err != nil {
			log.Error(err.Error(), location)
			time.Sleep(7 * time.Second)
			continue
		}

		if authToken != "" {
			log.Success("Got auth token", location)
			module.AuthToken = authToken
			break
		}

		time.Sleep(7 * time.Second)
	}

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

func (skate *skateSession) GetToken() (string, error) {
	log.Info("Getting auth token...", location)

	req, err := http.NewRequest(http.MethodGet, "https://www.skatedeluxe.com/it/login", nil)

	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		"Host":                      {"www.skatedeluxe.com"},
		"Connection":                {"keep-alive"},
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"Accept-Language":           {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
		"Cache-Control":             {"max-age=0"},
		"Sec-Ch-Ua":                 {"\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\""},
		"Sec-Ch-Ua-Mobile":          {"?0"},
		"Sec-Ch-Ua-Platform":        {"\"macOS\""},
		"Sec-Fetch-Dest":            {"document"},
		"Sec-Fetch-Mode":            {"navigate"},
		"Sec-Fetch-Site":            {"none"},
		"Sec-Fetch-User":            {"?1"},
		"Upgrade-Insecure-Requests": {"1"},
		"User-Agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
	}

	res, err := skate.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("Error while getting the token! [%v]", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return "", err
	}

	htmlString, _ := doc.Html()

	var authToken string

	tokenRegex := regexp.MustCompile(`data-token="([\s\S]*?)"`)
	token := tokenRegex.FindStringSubmatch(htmlString)

	if len(token) > 1 {
		authToken = token[1]
	} else {
		authToken = ""
	}

	return authToken, err
}

func (skate *skateSession) SearchProducts() error {
	url := "https://www.skatedeluxe.com/api2/search?q=campus"
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s&token=%s", url, skate.AuthToken), nil)

	if err != nil {
		return err
	}

	req.Header = http.Header{
		"Host":                      {"www.skatedeluxe.com"},
		"Connection":                {"keep-alive"},
		"Accept":                    {"*/*"},
		"Accept-Language":           {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
		"Cache-Control":             {"max-age=0"},
		"Sec-Ch-Ua":                 {"\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\""},
		"Sec-Ch-Ua-Mobile":          {"?0"},
		"Sec-Ch-Ua-Platform":        {"\"macOS\""},
		"Sec-Fetch-Dest":            {"empty"},
		"Sec-Fetch-Mode":            {"cors"},
		"Sec-Fetch-Site":            {"same-origin"},
		"Sec-Fetch-User":            {"?1"},
		"Upgrade-Insecure-Requests": {"1"},
		"User-Agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
	}

	res, err := skate.httpClient.Do(req)

	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Search products bad status code! [%v]", res.StatusCode)
	}

	var SearchRes SearchResponse

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &SearchRes)

	if err != nil {
		return err
	}

	// We need to parse the response here
	err = skate.ParseSearch(SearchRes)

	if err != nil {
		return err
	}
	return nil
}

func (skate *skateSession) ParseSearch(res SearchResponse) error {
	var productList []database.Product
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, product := range res.Products {
		wg.Add(1)
		// Get product info

		go skate.GetProduct(product, &productList, &wg, &mu)
	}

	wg.Wait()

	err := database.StartDataBase("./monitorData/skatedeluxe.json", productList)
	if err != nil {
		return err
	}

	return nil
}

func (skate skateSession) GetProduct(product Product, productList *[]database.Product, wg *sync.WaitGroup, mu *sync.Mutex) error {
	defer wg.Done()

	if strings.Contains(product.ProductsName, "00") && strings.Contains(strings.ToLower(product.ProductsName), "campus") {
		proxie, err := proxyManager.RandomProxy(skate.ProxyArr)
		proxie = strings.ReplaceAll(proxie, "\r", "")
		if err != nil {
			log.Error(err.Error(), location)
		} else {
			skate.httpClient.SetProxy(proxie)
			skate.Proxy = proxie
		}

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.skatedeluxe.com/it/_p%s", strconv.Itoa(product.ProductsID)), nil)

		if err != nil {
			return err
		}

		req.Header = http.Header{
			"Host":                      {"www.skatedeluxe.com"},
			"Connection":                {"keep-alive"},
			"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
			"Accept-Language":           {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
			"Cache-Control":             {"max-age=0"},
			"Sec-Ch-Ua":                 {"\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\""},
			"Sec-Ch-Ua-Mobile":          {"?0"},
			"Sec-Ch-Ua-Platform":        {"\"macOS\""},
			"Sec-Fetch-Dest":            {"document"},
			"Sec-Fetch-Mode":            {"navigate"},
			"Sec-Fetch-Site":            {"none"},
			"Sec-Fetch-User":            {"?1"},
			"Upgrade-Insecure-Requests": {"1"},
			"User-Agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
		}

		res, err := skate.httpClient.Do(req)

		if err != nil {
			return err
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		if res.StatusCode != 200 {
			return fmt.Errorf("Get product bad status code: [%v]", res.StatusCode)
		}

		// we can now parse the product page:
		doc, err := goquery.NewDocumentFromReader(res.Body)

		if err != nil {
			return err
		}

		var sizes []database.Size

		doc.Find("select option").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			value, _ := s.Attr("value")

			if value != "" && text != "" && containsUS(text) {
				sizes = append(sizes, database.Size{
					Size:  text,
					Value: "IN_STOCK",
				})
			}
		})

		tempProduct := database.Product{
			Site:  "SkateDeluxe",
			Sku:   strconv.Itoa(product.ProductsID),
			Name:  product.ProductsName,
			Price: strconv.FormatFloat(product.ProductsPriceCalced, 'f', -1, 64),
			Image: product.ProductImages[0],
			Sizes: sizes,
			Color: "None",
		}

		mu.Lock()
		*productList = append(*productList, tempProduct)
		mu.Unlock()
	}

	return nil
}

func containsUS(text string) bool {
	return text != "" && text[:2] == "US"
}
