package manashop

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
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

var location string = "Manashop"

func InitTls(proxy string) (manaSession, error) {
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
		return manaSession{}, err
	}

	return manaSession{
		httpClient: client,
	}, nil
}

func InitMana() {
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

	module.Host = "https://manashop.club/"
	module.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"
	module.Store_Region = "club"
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

func (mana *manaSession) GetProducts() error {

	var searchRes SearchRes

	url := "https://manashop.club/instantsearch/ajax/result?q=campus%2000"

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s&_=%s", url, GetTimeStamp()), nil)

	if err != nil {
		return err
	}

	req.Header = http.Header{
		"Accept":             {"application/json, text/javascript, */*; q=0.01"},
		"Accept-Language":    {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
		"Connection":         {"keep-alive"},
		"Referer":            {"https://manashop.club/shop-by-brand/originals"},
		"Sec-Fetch-Dest":     {"empty"},
		"Sec-Fetch-Mode":     {"cors"},
		"Sec-Fetch-Site":     {"same-origin"},
		"User-Agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"X-Requested-With":   {"XMLHttpRequest"},
		"sec-ch-ua":          {"\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\""},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"\"macOS\""},
	}

	res, err := mana.httpClient.Do(req)

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

	err = json.Unmarshal(body, &searchRes)

	if err != nil {
		return err
	}

	mana.ParseSearch(searchRes)

	return nil
}

func (mana *manaSession) ParseSearch(data SearchRes) error {
	var productList []database.Product
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, product := range data.Result.Product.Data {
		wg.Add(1)

		go mana.GetProductInfo(product, &productList, &wg, &mu)

	}

	wg.Wait()

	err := database.StartDataBase("./monitorData/manashop.json", productList)
	if err != nil {
		return err
	}

	return nil
}

func (mana *manaSession) GetProductInfo(product Data, productList *[]database.Product, wg *sync.WaitGroup, mu *sync.Mutex) error {
	defer wg.Done()
	var sizes []database.Size

	req, err := http.NewRequest(http.MethodGet, product.URL, nil)

	if err != nil {
		log.Error(err.Error(), location)
		return err
	}

	req.Header = http.Header{
		"Accept":             {"application/json, text/javascript, */*; q=0.01"},
		"Accept-Language":    {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
		"Connection":         {"keep-alive"},
		"Referer":            {"https://manashop.club/shop-by-brand/originals"},
		"Sec-Fetch-Dest":     {"empty"},
		"Sec-Fetch-Mode":     {"cors"},
		"Sec-Fetch-Site":     {"same-origin"},
		"User-Agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"X-Requested-With":   {"XMLHttpRequest"},
		"sec-ch-ua":          {"\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\""},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"\"macOS\""},
	}

	res, err := mana.httpClient.Do(req)

	if err != nil {
		log.Error(err.Error(), location)
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		log.Error(fmt.Sprintf("Get product bad status code! [%v]", res.StatusCode), location)
		return fmt.Errorf("Get product bad status code! [%v]", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Error(err.Error(), location)
		return err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		log.Error(err.Error(), location)
		return err
	}

	var jsonScript string

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		content := s.Text()
		if strings.Contains(content, "Magento_Swatches/js/swatch-renderer") {
			jsonScript = content
			return
		}
	})

	jsonStartIndex := strings.Index(jsonScript, "{")
	jsonEndIndex := strings.LastIndex(jsonScript, "}")

	if jsonStartIndex == -1 || jsonEndIndex == -1 {
		log.Error("Json not found!", location)
		return fmt.Errorf("Json not found!")
	}

	jsonContent := jsonScript[jsonStartIndex : jsonEndIndex+1]

	var productData MyData

	err = json.Unmarshal([]byte(jsonContent), &productData)
	if err != nil {
		return err
	}

	// Now we can parse the product

	color := productData.SwatchOptions.SwatchRenderer.JSONConfig.Attributes["93"].Label

	for _, size := range productData.SwatchOptions.SwatchRenderer.JSONConfig.Attributes["144"].Options {
		sizes = append(sizes, database.Size{
			Size:  size.Label,
			Value: "IN_STOCK",
		})
	}

	price := productData.SwatchOptions.SwatchRenderer.JSONConfig.Prices.FinalPrice.Amount
	//sku := productData.DataRoleSwatchOptions["[data-role=swatch-options]"].JsonConfig.ProductId

	tempProduct := database.Product{
		Site:  "Mana",
		Sku:   product.URL,
		Name:  product.Name,
		Price: strconv.FormatFloat(price, 'f', -1, 64),
		Image: product.Image,
		Color: color,
		Sizes: sizes,
	}

	mu.Lock()
	*productList = append(*productList, tempProduct)
	mu.Unlock()

	return nil
}
