package lvr

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/TheAldriguzDev/Bonsai_solutions/packages/log"
	proxyManager "github.com/TheAldriguzDev/Bonsai_solutions/packages/proxy"
	filereader "github.com/TheAldriguzDev/Bonsai_solutions/packages/reader"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

var location string = "ASOS"

func InitTls(proxy string) (lvrSession, error) {
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
		return lvrSession{}, err
	}

	return lvrSession{
		httpClient: client,
	}, nil
}

func InitLvr() {
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

	// var run bool = true
	// for run {
	// 	log.General("Monitoring...", location)

	// 	// we have to rotate the proxy for every request
	// 	err = module.RotateProxy()
	// 	if err != nil {
	// 		log.Error(err.Error(), location)
	// 	}

	// 	// Now we can monitor the new add
	// 	err = module.SearchCampus()
	// 	if err != nil {
	// 		log.Error(err.Error(), location)
	// 		time.Sleep(15 * time.Second)
	// 		module, err = InitTls(proxie)
	// 		module.Proxy = proxie
	// 		module.Host = "jdsports.it"
	// 		module.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"
	// 		module.Store_Region = "it"
	// 		module.ProxyArr = proxyArr
	// 		continue
	// 	}
	// 	time.Sleep(15 * time.Second)
	// }
}

func (lvr *lvrSession) RotateProxy() error {
	proxie, err := proxyManager.RandomProxy(lvr.ProxyArr)
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
		err = lvr.httpClient.SetProxy(proxyUrl)
		if err != nil {
			return err
		}
		lvr.Proxy = proxie
	}
	return nil
}

func (lvr *lvrSession) SearchProducts() error {
	// For every request we generate a new auth token --> less cache

	url := "https://api.luisaviaroma.com/lvrapprk/public/v1/catalog/bysearch?Gender=women&Search=Campus%2080&Season=actual&IncludeTreeNavigation=true&Language=IT&Country=IT&CurrencyView=EUR&CurrencyFatt=EUR&App=true&format=json"

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(url+"&ts=", GetTimeStamp()), nil)

	if err != nil {
		return err
	}

	req.Header = http.Header{
		"Host":                     {"api.luisaviaroma.com"},
		"__lvr_mobile_api_token__": {TokenGen()},
		"accept":                   {"application/json, text/plain, */*"},
		"accept-language":          {"it-IT,it;q=0.9"},
		"env":                      {"PROD"},
		"user-agent":               {"LVR/2.1.90 (iPhone; iOS 16.5)"},
	}

	res, err := lvr.httpClient.Do(req)

	if err != nil {
		return err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Get products bad status code! %v", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	var searchData Data

	err = json.Unmarshal(body, &searchData)

	if err != nil {
		return err
	}

	// Here we need to parse the response

	

	return nil

}
