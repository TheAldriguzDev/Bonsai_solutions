package asos

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
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

var location string = "ASOS"

func InitTls(proxy string) (asosSession, error) {
	proxy = strings.ReplaceAll(proxy, "\r", "")
	proxyData := strings.Split(strings.TrimSpace(proxy), ":")
	proxyUrl := fmt.Sprintf("http://%s:%s@%s:%s", proxyData[2], proxyData[3], proxyData[0], proxyData[1])

	jarOptions := []tls_client.CookieJarOption{}

	jar := tls_client.NewCookieJar(jarOptions...)

	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(30),
		tls_client.WithClientProfile(tls_client.MeshIos2),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
		tls_client.WithRandomTLSExtensionOrder(),
		tls_client.WithProxyUrl(proxyUrl),
		tls_client.WithInsecureSkipVerify(),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)

	if err != nil {
		fmt.Println("Error while creating the tls client!")
		return asosSession{}, err
	}

	return asosSession{
		httpClient: client,
	}, nil
}

func InitAsos() {
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

	module.Host = "jdsports.it"
	module.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"
	module.Store_Region = "it"
	module.ProxyArr = proxyArr
	log.Success("Monitor initialised!", location)

	var run bool = true
	for run {
		log.General("Monitoring...", location)

		// we have to rotate the proxy for every request
		err = module.RotateProxy()
		if err != nil {
			log.Error(err.Error(), location)
		}

		// Now we can monitor the new add
		err = module.SearchCampus()
		if err != nil {
			log.Error(err.Error(), location)
			time.Sleep(15 * time.Second)
			module, err = InitTls(proxie)
			module.Proxy = proxie
			module.Host = "jdsports.it"
			module.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"
			module.Store_Region = "it"
			module.ProxyArr = proxyArr
			continue
		}
		time.Sleep(15 * time.Second)
	}
}

func (asos *asosSession) RotateProxy() error {
	proxie, err := proxyManager.RandomProxy(asos.ProxyArr)
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
		err = asos.httpClient.SetProxy(proxyUrl)
		if err != nil {
			return err
		}
		asos.Proxy = proxie
	}
	return nil
}

func (asos *asosSession) SearchCampus() error {

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.asos.com/api/product/search/v2/?offset=0&q=campus&store=IT&lang=it-IT&currency=EUR&rowlength=2&channel=desktop-web&country=IT&keyStoreDataversion=ornjx7v-36&limit=%v", genNum(70, 1000000000)), nil)

	if err != nil {
		return err
	}

	req.Header = http.Header{
		"authority":          {"www.asos.com"},
		"accept":             {"application/json, text/plain, */*"},
		"accept-language":    {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
		"asos-c-name":        {"@asosteam/asos-web-product-listing-page"},
		"asos-c-plat":        {"web"},
		"asos-c-ver":         {"1.2.0-e16f1be4ec2f-9399"},
		"asos-cid":           {"e7bdd551-5dc9-4ab0-a15e-3e72202af608"},
		"referer":            {"https://www.asos.com/it/search/?q=campus+"},
		"sec-ch-ua":          {"\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\""},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"\"macOS\""},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"same-origin"},
		"user-agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
	}

	res, err := asos.httpClient.Do(req)

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

	var searchRes Data

	err = json.Unmarshal(body, &searchRes)

	if err != nil {
		return err
	}

	// Here we need to parse the response!

	err = asos.ParseSearch(searchRes)

	if err != nil {
		return err
	}

	return nil
}

func (asos *asosSession) ParseSearch(searchData Data) error {
	var productList []database.Product
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, product := range searchData.Products {
		wg.Add(1)

		// get the product info
		go asos.GetProduct(product, &productList, &wg, &mu)
	}

	wg.Wait()

	err := database.StartDataBase("./monitorData/asos.json", productList)
	if err != nil {
		return err
	}

	return nil
}

func (asos asosSession) GetProduct(product ProductSearch, productList *[]database.Product, wg *sync.WaitGroup, mu *sync.Mutex) error {
	defer wg.Done()

	if strings.Contains(product.Name, "00") && product.ID != 203414519 {
		// Here we rotate proxies to reduce proxy bans

		proxie, err := proxyManager.RandomProxy(asos.ProxyArr)
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
			err = asos.httpClient.SetProxy(proxyUrl)
			if err != nil {
				log.Error(err.Error(), location)
				return err
			}
			asos.Proxy = proxie
		}

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://api.asos.com/product/catalogue/v3/products/%s?currency=EUR&keyStoreDataversion=ornjx7v-36&sizeSchema=IT&lang=it-IT&store=IT", strconv.Itoa(product.ID)), nil)

		if err != nil {
			log.Error(err.Error(), location)
			return err
		}

		req.Header = http.Header{
			"authority":          {"www.asos.com"},
			"accept":             {"application/json, text/plain, */*"},
			"accept-language":    {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
			"asos-c-name":        {"@asosteam/asos-web-product-listing-page"},
			"asos-c-plat":        {"web"},
			"asos-c-ver":         {"1.2.0-e16f1be4ec2f-9399"},
			"asos-cid":           {"e7bdd551-5dc9-4ab0-a15e-3e72202af608"},
			"referer":            {"https://www.asos.com/it/search/?q=campus+"},
			"sec-ch-ua":          {"\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\""},
			"sec-ch-ua-mobile":   {"?0"},
			"sec-ch-ua-platform": {"\"macOS\""},
			"sec-fetch-dest":     {"empty"},
			"sec-fetch-mode":     {"cors"},
			"sec-fetch-site":     {"same-origin"},
			"user-agent":         {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
		}

		res, err := asos.httpClient.Do(req)

		if err != nil {
			log.Error(err.Error(), location)
			return err
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		if res.StatusCode != 200 {
			log.Error("Get product bad status code!", location)
			return fmt.Errorf("Get product bad status code! %v", res.StatusCode)
		}

		body, err := io.ReadAll(res.Body)

		if err != nil {
			log.Error(err.Error(), location)
			return err
		}

		var productRes Product

		err = json.Unmarshal(body, &productRes)

		if err != nil {
			log.Error(err.Error(), location)
			return err
		}

		// Here we need to parse the product
		sizes := []database.Size{}

		for _, size := range productRes.Variants {
			if size.IsInStock {
				sizes = append(sizes, database.Size{
					Size:  size.DisplaySizeText,
					Value: strconv.Itoa(size.ID),
				})
			}
		}

		tempProduct := database.Product{
			Site:  "Asos",
			Sku:   strconv.Itoa(product.ID),
			Name:  product.Name,
			Price: product.Price.Current.Text,
			Image: product.ImageURL,
			Color: product.Colour,
			Sizes: sizes,
		}

		if len(sizes) != 0 {
			mu.Lock()
			*productList = append(*productList, tempProduct)
			mu.Unlock()
		}
	}

	return nil

}

func genNum(min, max int) int {
	// Imposta un seme diverso per ogni esecuzione
	rand.Seed(time.Now().UnixNano())

	// Genera un numero casuale compreso tra min e max
	return rand.Intn(max-min+1) + min
}
