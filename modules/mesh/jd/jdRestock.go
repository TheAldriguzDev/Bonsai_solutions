package jd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/TheAldriguzDev/Bonsai_solutions/packages/database"
	"github.com/TheAldriguzDev/Bonsai_solutions/packages/log"
	proxyManager "github.com/TheAldriguzDev/Bonsai_solutions/packages/proxy"
	filereader "github.com/TheAldriguzDev/Bonsai_solutions/packages/reader"
	http "github.com/bogdanfinn/fhttp"
)

//./monitorData/jdNewAdd.json

var locale string = "JD Restock"

func InitJdRestock() {
	proxyFile, err := filereader.ReadFolder("proxies", "proxies")

	if err != nil {
		log.Error(err.Error(), locale)
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}

	proxyArr, err := proxyManager.ReadProxyFile(proxyFile)

	if err != nil {
		log.Error(err.Error(), locale)
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}

	proxie, err := proxyManager.RandomProxy(proxyArr)

	if err != nil {
		log.Error(err.Error(), locale)
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}

	module, err := InitTls(proxie)
	module.Proxy = proxie
	if err != nil {
		log.Error(err.Error(), locale)
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}

	module.Host = "jdsports.it"
	module.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"
	module.Store_Region = "it"
	module.ProxyArr = proxyArr
	log.Success("Monitor initialised!", locale)

	var run bool = true
	for run {
		log.General("Monitoring...", locale)

		// we have to rotate the proxy for every request
		err = module.RotateProxy()
		if err != nil {
			log.Error(err.Error(), locale)
		}

		// Now we can monitor the new add
		err = module.ReadFile()
		if err != nil {
			log.Error(err.Error(), locale)
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

type ProductFile struct {
	SKU string `json:"sku"`
}

type File struct {
	Products []ProductFile `json:"products"`
}

func (mesh *meshSession) ReadFile() error {
	var productList []database.Product
	var wg sync.WaitGroup
	var mu sync.Mutex

	data, err := ioutil.ReadFile("./monitorData/jdNewAdd.json")
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
		go mesh.GetSizes(product.SKU, &productList, &wg, &mu)
	}

	wg.Wait()

	err = database.StartDataBase("./monitorData/jdRestock.json", productList)
	if err != nil {
		return err
	}

	return nil
}

func (mesh *meshSession) GetSizes(sku string, productList *[]database.Product, wg *sync.WaitGroup, mu *sync.Mutex) error {
	defer wg.Done()

	// Here i need to add proxy rotation to prevent proxy bans

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://prod.jdgroupmesh.cloud/stores/jdsportsit/products/%s?api_key=%s&channel=iphone-app&debug=1&locale=it", sku, getRandomString(jdKeys)), nil)

	if err != nil {
		log.Error(err.Error(), locale)
		return err
	}

	req.Header = http.Header{
		"authority":                 {"prod.jdgroupmesh.cloud"},
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"accept-language":           {"fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7"},
		"cache-control":             {"no-cache"},
		"pragma":                    {"no-cache"},
		"sec-ch-ua":                 {"\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\""},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {"\"Windows\""},
		"sec-fetch-dest":            {"document"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-site":            {"none"},
		"sec-fetch-user":            {"?1"},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
		//"cookie":                    {"----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------"},
	}

	res, err := mesh.httpClient.Do(req)

	if err != nil {
		log.Error(err.Error(), locale)
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		log.Error(fmt.Sprintf("Get product bad status code! %v", res.StatusCode), locale)
		return fmt.Errorf("Get product bad status code! [%v]", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Error(err.Error(), locale)
		return err
	}

	var productRes ItemData

	err = json.Unmarshal(body, &productRes)

	if err != nil {
		log.Error(err.Error(), locale)
		return err
	}

	// Here we need to parse the products

	product, err := mesh.ParseProduct(productRes)

	if err != nil {
		log.Error(err.Error(), locale)
		return err
	}

	mu.Lock()
	*productList = append(*productList, product)
	mu.Unlock()

	return nil
}

func (mesh *meshSession) ParseProduct(product ItemData) (database.Product, error) {

	sizes := []database.Size{}

	for _, size := range product.SortedOptions {
		if size.Product.StockStatus == "IN STOCK" {
			sizes = append(sizes, database.Size{
				Size:  size.Name,
				Value: size.Product.SKU,
			})
		}
	}

	if len(sizes) == 0 {
		tempProduct := database.Product{}
		return tempProduct, fmt.Errorf("No sizes has been found!")
	}

	tempProduct := database.Product{
		Site:  "JD Restock",
		Sku:   product.SKU,
		Name:  product.Name,
		Price: product.Price.Amount,
		Image: product.MainImage,
		Sizes: sizes,
		Color: product.Colour,
	}

	return tempProduct, nil
}
