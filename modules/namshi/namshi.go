package namshi

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/TheAldriguzDev/Bonsai_solutions/packages/database"
	"github.com/TheAldriguzDev/Bonsai_solutions/packages/log"
	proxyManager "github.com/TheAldriguzDev/Bonsai_solutions/packages/proxy"
	filereader "github.com/TheAldriguzDev/Bonsai_solutions/packages/reader"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

var location string = "NAMSHI"

func InitTls(proxy string) (namshiSession, error) {
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
		return namshiSession{}, err
	}

	return namshiSession{
		httpClient: client,
	}, nil
}

func InitNamshi() {
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

	module.Host = "namshi.it"
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
		err = module.GetCampus()
		if err != nil {
			log.Error(err.Error(), location)
			time.Sleep(15 * time.Second)
			module, err = InitTls(proxie)
			module.Proxy = proxie
			module.Host = "namshi.it"
			module.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"
			module.Store_Region = "it"
			module.ProxyArr = proxyArr
			continue
		}
		time.Sleep(15 * time.Second)
	}
}

func (namshi *namshiSession) RotateProxy() error {
	proxie, err := proxyManager.RandomProxy(namshi.ProxyArr)
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
		err = namshi.httpClient.SetProxy(proxyUrl)
		if err != nil {
			return err
		}
		namshi.Proxy = proxie
	}
	return nil
}

func (namshi *namshiSession) GetCampus() error {
	url := "https://en-global.namshi.com/adidas_originals/?q=adidas%20campus%2000s"
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s&ts=%s", url, GetTimeStamp()), nil)

	if err != nil {
		return err
	}

	req.Header = http.Header{
		"authority":                 {"en-global.namshi.com"},
		"accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
		"accept-language":           {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
		"referer":                   {"https://en-global.namshi.com/catalog/?q=adidas%20campus%2000s"},
		"sec-ch-ua":                 {"\"\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"114\", \"Google Chrome\";v=\"114\"\""},
		"sec-ch-ua-mobile":          {"?0"},
		"sec-ch-ua-platform":        {"\"macOS\""},
		"sec-fetch-dest":            {"document"},
		"sec-fetch-mode":            {"navigate"},
		"sec-fetch-site":            {"same-origin"},
		"upgrade-insecure-requests": {"1"},
		"user-agent":                {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"},
	}

	res, err := namshi.httpClient.Do(req)

	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Get products bad status code! [%v]", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	// Now we need to parse the products

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))

	if err != nil {
		return err
	}

	var productList []database.Product

	doc.Find("li[class*='listing closed']").Each(func(i int, elem *goquery.Selection) {
		url, _ := elem.Find("a").Attr("href")

		image, _ := elem.Find("a").Find("div[class*='image_container']").Find("img[class*='list-product-image']").Attr("id")

		parts := strings.Split(image, "/")

		var value string

		value = parts[1]

		var imageUrl string = ""

		if parts[0] == "A0" || parts[0] == "A4" {
			imageUrl = fmt.Sprintf("https://pa.namshicdn.com/product/%s/%s/1-web-desktop-list.webp", parts[0], value)
		} else {
			imageUrl = fmt.Sprintf("https://pa.namshicdn.com/product/A1/%s/1-web-desktop-list.jpg", value)
		}

		brand := elem.Find("p[class*='brand']").Text()
		description := elem.Find("p[class*='description']").Text()
		name := brand + " " + description

		var price string
		if specialPrice := elem.Find("span[class*='specialPrice']"); specialPrice.Length() > 0 {
			price = specialPrice.Text()
		} else {
			price = elem.Find("span[class*='original_price']").Text()
		}

		var sizes = []database.Size{}
		elem.Find("li[class*='size_listing simple']").Each(func(i int, sizeElem *goquery.Selection) {
			sizes = append(sizes, database.Size{
				Size:  strings.TrimSpace(sizeElem.Text()),
				Value: "IN_STOCK",
			})
		})

		tempProduct := database.Product{
			Site:  "Namshi",
			Sku:   url,
			Name:  name,
			Price: price,
			Image: imageUrl,
			Sizes: sizes,
			Color: "nil",
		}

		productList = append(productList, tempProduct)
	})

	err = database.StartDataBase("./monitorData/namshi.json", productList)
	if err != nil {
		return err
	}

	return nil
}
