package size

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TheAldriguzDev/Bonsai_solutions/packages/database"
	"github.com/TheAldriguzDev/Bonsai_solutions/packages/log"
	proxyManager "github.com/TheAldriguzDev/Bonsai_solutions/packages/proxy"
	filereader "github.com/TheAldriguzDev/Bonsai_solutions/packages/reader"
	http "github.com/bogdanfinn/fhttp"
)

var loc string = "JD NEW ADD"

func InitMeshNewAdd() {
	proxyFile, err := filereader.ReadFolder("proxies", "proxies")

	if err != nil {
		log.Error(err.Error(), loc)
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}

	proxyArr, err := proxyManager.ReadProxyFile(proxyFile)

	if err != nil {
		log.Error(err.Error(), loc)
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}

	proxie, err := proxyManager.RandomProxy(proxyArr)

	if err != nil {
		log.Error(err.Error(), loc)
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}

	module, err := InitTls(proxie)
	module.Proxy = proxie
	if err != nil {
		log.Error(err.Error(), loc)
		time.Sleep(2 * time.Second)
		os.Exit(0)
	}

	module.Host = "size.it"
	module.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"
	module.Store_Region = "it"
	module.ProxyArr = proxyArr
	log.Success("Monitor initialised!", loc)

	var run bool = true
	for run {
		log.General("Monitoring...", loc)

		// we have to rotate the proxy for every request
		err = module.RotateProxy()
		if err != nil {
			log.Error(err.Error(), loc)
		}

		// Now we can monitor the new add
		err = module.GetNewAdd()
		if err != nil {
			log.Error(err.Error(), loc)
			time.Sleep(15 * time.Second)
			module, err = InitTls(proxie)
			module.Proxy = proxie
			module.Host = "size.it"
			module.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"
			module.Store_Region = "it"
			module.ProxyArr = proxyArr
			continue
		}
		time.Sleep(15 * time.Second)
	}
}

func (mesh *meshSession) RotateProxy() error {
	proxie, err := proxyManager.RandomProxy(mesh.ProxyArr)
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
		err = mesh.httpClient.SetProxy(proxyUrl)
		if err != nil {
			return err
		}
		mesh.Proxy = proxie
	}
	return nil
}

func (mesh *meshSession) GetNewAdd() error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://prod.jdgroupmesh.cloud/stores/sizeit/products/search?api_key=%s&max=%s&customerLoyalty=1&from=0&channel=iphone-app&debug=1&query=campus&locale=it", getRandomString(sizeKeys), genNumber(24, 128)), nil)

	if err != nil {
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
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		fmt.Println(mesh.httpClient.GetProxy())
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

	err = mesh.ParseNewAdd(searchRes)

	if err != nil {
		return err
	}

	return nil
}

func (mesh *meshSession) ParseNewAdd(searchData SearchRes) error {
	var productList []database.Product
	//var mu sync.Mutex

	for _, product := range searchData.Products {

		if strings.Contains(product.Name, "00") || strings.Contains(strings.ToLower(product.Name), "junior") {
			// We parse the response and save in the jdNewAdd file, than another script will monitor the saved pids
			sizes := []database.Size{}
			tempProduct := database.Product{
				Site:  "SIZE New",
				Sku:   product.SKU,
				Name:  product.Name,
				Price: product.Price.Amount,
				Image: product.MainImage,
				Color: product.Colour,
				Sizes: sizes,
			}

			productList = append(productList, tempProduct)
		}
	}

	err := database.StartDataBase("./monitorData/sizeNewAdd.json", productList)
	if err != nil {
		return err
	}

	return nil
}

func genNumber(min, max int) string {
	rand.Seed(time.Now().UnixNano())

	randomNumber := rand.Intn(max-min+1) + min

	randomNumberString := strconv.Itoa(randomNumber)

	return randomNumberString
}
