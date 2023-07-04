package sotf

import (
	"bytes"
	"fmt"
	"io"
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

func InitTls(proxy string) (sotfSession, error) {
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
		return sotfSession{}, err
	}

	return sotfSession{
		httpClient: client,
	}, nil
}

func InitSotf() {
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

	module.Host = "https://www.sotf.com/"
	module.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"
	module.Store_Region = "com"
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

func (sotf *sotfSession) GetProducts() error {
	url := "https://www.sotf.com/it/search.html?searchkey=campus%2000&orderBy=Arrivo&orderWay=ASC&num_art_per_pagina=30&mod_itemview=&groupBy=LevelCat1&resultpage=1&filters%5B41%5D%5B0%5D=116"

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s&ts=%s", url, GetTimeStamp()), nil)

	if err != nil {
		return err
	}

	req.Header = http.Header{
		"authority":          {"www.sotf.com"},
		"accept":             {"*/*"},
		"accept-language":    {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
		"referer":            {"https://www.sotf.com/it/search.html?searchkey=campus%2000&orderBy=Arrivo&orderWay=ASC&num_art_per_pagina=30&mod_itemview=&groupBy=LevelCat1&resultpage=1"},
		"sec-ch-ua":          {"\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\""},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"\"macOS\""},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"same-origin"},
		"user-agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"x-requested-with":   {"XMLHttpRequest"},
	}

	res, err := sotf.httpClient.Do(req)

	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Search products bad status code! %v", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	err = sotf.ParseSearch(body)

	if err != nil {
		return err
	}
	return nil
}

func (sotf *sotfSession) ParseSearch(body []byte) error {
	var productList []database.Product
	var wg sync.WaitGroup
	var mu sync.Mutex

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return err
	}

	doc.Find("div[class*='vaschetta_item col-lg-4 col-sm-6 col-6 p-2 retina_hoverimg']").Each(func(i int, s *goquery.Selection) {
		wg.Add(1)
		url := s.Find("a").AttrOr("href", "")

		go sotf.GetProduct(url, &productList, &wg, &mu)
	})

	wg.Wait()

	err = database.StartDataBase("./monitorData/sotf.json", productList)
	if err != nil {
		return err
	}

	return nil
}

func (sotf *sotfSession) GetProduct(url string, productList *[]database.Product, wg *sync.WaitGroup, mu *sync.Mutex) error {
	defer wg.Done()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?ts=%s", url, GetTimeStamp()), nil)

	if err != nil {
		log.Error(err.Error(), location)
		return err
	}

	req.Header = http.Header{
		"authority":          {"www.sotf.com"},
		"accept":             {"*/*"},
		"accept-language":    {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
		"referer":            {"https://www.sotf.com/it/search.html?searchkey=campus%2000&orderBy=Arrivo&orderWay=ASC&num_art_per_pagina=30&mod_itemview=&groupBy=LevelCat1&resultpage=1"},
		"sec-ch-ua":          {"\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\""},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"\"macOS\""},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"same-origin"},
		"user-agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		"x-requested-with":   {"XMLHttpRequest"},
	}

	res, err := sotf.httpClient.Do(req)

	if err != nil {
		log.Error(err.Error(), location)
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		log.Error("Search bad status code!", location)
		return fmt.Errorf("Get product bad status code! %v", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Error(err.Error(), location)
		return err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Error(err.Error(), location)
		return err
	}

	// Now we need to parse the pdp

	name := doc.Find("div[class*='details_info_brand marginBottom0']").Text()

	price := doc.Find("div[class*='details_price']").Text()

	image, exists := doc.Find("img[id*='ImgDetails']").Attr("src")
	if !exists {
		image = "https://www.sotf.com/common/images/logo.svg?v=3.0"
	}

	var sizes []database.Size
	doc.Find("ul[role*='radiogroup'] li").Each(func(i int, s *goquery.Selection) {
		option := strings.TrimSpace(s.Text())

		sizes = append(sizes, database.Size{
			Size:  option,
			Value: "IN_STOCK",
		})
	})

	tempProduct := database.Product{
		Site:  "Sotf",
		Sku:   url,
		Name:  name,
		Price: price,
		Image: image,
		Sizes: sizes,
		Color: "None",
	}

	mu.Lock()
	*productList = append(*productList, tempProduct)
	mu.Unlock()

	return nil
}
