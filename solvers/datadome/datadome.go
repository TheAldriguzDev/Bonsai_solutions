package datadome

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
)

func InitClient(proxy string) (DataDomeSession, error) {
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
		return DataDomeSession{}, err
	}

	return DataDomeSession{
		httpClient: client,
	}, nil
}

func GetCookie(genData GenerateData) (DDResponse, error) {
	var ddRes DDResponse

	if genData.DataDomeKey == "" {
		return ddRes, fmt.Errorf("Datadome key expected but received nil!")
	}

	if genData.DataDomeVersion == "" {
		return ddRes, fmt.Errorf("DataDome version expected, but received nil!")
	}

	if genData.Proxy == "" {
		return ddRes, fmt.Errorf("Proxy expected, but received nil!")
	}

	if genData.Refer == "" {
		return ddRes, fmt.Errorf("Refer expected, but received nil!")
	}

	if genData.UserAgent == "" {
		return ddRes, fmt.Errorf("UserAgent expected but received nil!")
	}

	module, err := InitClient(genData.Proxy)

	if err != nil {
		return ddRes, err
	}

	if genData.Cid == "" || genData.Cid == "null" {
		cid, err := module.GetCid(genData)

		if err != nil {
			fmt.Println("Error while generating the cid!")
			genData.Cid = "null"
		} else {
			genData.Cid = cid
		}
	}

	if strings.Contains(genData.Refer, "hermes") {
		module.DataDomeDomain = "https://d.digital.hermes/js/"
	} else {
		module.DataDomeDomain = "https://api-js.datadome.co/js/"
	}

	parsedUrl, err := url.Parse(genData.Refer)

	if err != nil {
		return ddRes, err
	}

	genData.Path = parsedUrl.Path
	genData.Host = parsedUrl.Host

	payload, err := module.GenFingerPrint(genData)

	if err != nil {
		return ddRes, err
	}

	module.Payload = payload

	ddRes, err = module.RequestCookie()

	if err != nil {
		return ddRes, err
	}

	return ddRes, nil
}

func (dd *DataDomeSession) GetCid(genData GenerateData) (string, error) {
	fmt.Println("Generating cid for: " + genData.Refer)

	req, err := http.NewRequest(http.MethodGet, genData.Refer, nil)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	req.Header = http.Header{
		"sec-ch-ua":          {"\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\""},
		"sec-ch-ua-mobile":   {"?0"},
		"user-agent":         {genData.UserAgent},
		"sec-ch-ua-platform": {"\"macOS\""},
		"accept":             {"*/*"},
		"sec-fetch-site":     {"cross-site"},
		"sec-fetch-mode":     {"no-cors"},
		"sec-fetch-dest":     {"script"},
		"referer":            {genData.Refer},
		"accept-encoding":    {"gzip, deflate, br"},
		"accept-language":    {"de-DE,de;q=0.9"},
	}

	res, err := dd.httpClient.Do(req)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	if res.StatusCode != 200 && res.StatusCode != 302 {
		return "", fmt.Errorf("Bad status code!")
	}

	cookies := res.Cookies()
	var datadomeCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "datadome" {
			datadomeCookie = cookie
			break
		}
	}

	fmt.Println(res.StatusCode)

	return datadomeCookie.Value, err
}

func (ddModule *DataDomeSession) GenFingerPrint(data GenerateData) ([]Param, error) {
	var param []Param
	t := time.Now()
	outerHeight, outerWidth := GenerateWindow()
	screenHeight, screenWidth := GenerateScreen()
	_, offset := t.Zone()

	uaLower := strings.ToLower(data.UserAgent)

	var fingerEval int
	if strings.Contains(uaLower, "Chrome") || strings.Contains(uaLower, "Opera") || strings.Contains(uaLower, "Other") {
		fingerEval = 33
	} else if strings.Contains(uaLower, "Internet Explorer") || strings.Contains(uaLower, "Other") {
		fingerEval = 39
	} else if strings.Contains(uaLower, "Safari") || strings.Contains(uaLower, "Firefox") || strings.Contains(uaLower, "Other") {
		fingerEval = 37
	}

	jsData := JSData{
		TimeSinceLastEvent: GenerateNumber(15.0, 230.0, 14),
		Ifov:               false,
		TimeToDoWebGL:      GenerateNumber(3.0, 13, 15),
		Glvd:               "Google Inc. (Intel Inc.)",
		Glrd:               "ANGLE (Intel Inc., Intel(R) Iris(TM) Plus Graphics 645, OpenGL 4.1)",
		Hc:                 generateCpu(),
		BrOh:               outerHeight,
		BrOw:               outerWidth,
		Ua:                 data.UserAgent,
		Wbd:                false,
		Wdif:               false,
		Wdifrm:             false,
		Npmtm:              false,
		BrH:                969,
		BrW:                1341,
		CookieLength:       1,
		RsH:                screenHeight,
		RsW:                screenWidth,
		RsCd:               24,
		Phe:                false,
		Nm:                 false,
		Jsf:                false,
		Lg:                 "it-IT",
		Pr:                 1,
		ArsH:               screenHeight,
		ArsW:               screenWidth,
		TimezoneOffset:     offset,
		StrSs:              true,
		StrLs:              true,
		StrIdb:             true,
		StrOdb:             true,
		Plgod:              false,
		NavPluginLength:    5,
		Plgne:              true,
		Plgre:              true,
		Plgof:              false,
		Plggt:              false,
		Pltod:              false,
		Hcovdr:             false,
		Hcovdr2:            false,
		Plovdr:             false,
		Plovdr2:            false,
		Ftsovdr:            false,
		Ftsovdr2:           false,
		Lb:                 false,
		EvalLength:         fingerEval,
		Lo:                 false,
		TsMtp:              0,
		TsTec:              false,
		TsTsa:              false,
		Vnd:                "Google Inc.",
		Bid:                "NA",
		Mmt:                "application/pdf,text/pdf",
		Plu:                "PDF Viewer,Chrome PDF Viewer,Chromium PDF Viewer,Microsoft Edge PDF Viewer,WebKit built-in PDF",
		Hdn:                false,
		Awe:                false,
		Geb:                false,
		Dat:                false,
		Med:                "defined",
		Aco:                "probably",
		Acots:              false,
		Acmp:               "probably",
		Acmpts:             true,
		Acw:                "probably",
		Acwts:              false,
		Acma:               "maybe",
		Acmats:             false,
		Acaa:               "probably",
		Acaats:             true,
		Ac3:                "",
		Ac3ts:              false,
		Acf:                "probably",
		Acfts:              false,
		Acmp4:              "maybe",
		Acmp4ts:            false,
		Acmp3:              "probably",
		Acmp3ts:            false,
		Acwm:               "maybe",
		Acwmts:             false,
		Ocpt:               false,
		Vco:                "probably",
		Vcots:              false,
		Vch:                "probably",
		Vchts:              true,
		Vcw:                "probably",
		Vcwts:              true,
		Vc3:                "maybe",
		Vc3ts:              false,
		Vcmp:               "",
		Vcmpts:             false,
		Vcq:                "",
		Vcqts:              false,
		Vc1:                "probably",
		Vc1ts:              true,
		Dvm:                8,
		Sqt:                false,
		So:                 "landscape-primary",
		Wdw:                true,
		Cokys:              "bG9hZFRpbWVzY3NpYXBwL=",
		Ecpc:               false,
		Lgs:                true,
		Lgsod:              false,
		Psn:                true,
		Edp:                true,
		Addt:               true,
		Wsdc:               true,
		Ccsr:               true,
		Nuad:               true,
		Bcda:               true,
		Idn:                true,
		Capi:               false,
		Svde:               false,
		Vpbq:               true,
		Ucdv:               false,
		Spwn:               false,
		Emt:                false,
		Bfr:                false,
		Dbov:               false,
		Cfpfe:              "RXJyb3I6IENhbm5vdCByZWFkIHByb3BlcnRpZXMgb2YgbnVsbCAocmVhZGluZyAndG9TdHJpbmcnKQ==", // ricorda che gli hai modificati per testare liverpool
		Stcfp:              "IChodHRwczovL2Nkbi5jb29raWVsYXcub3JnL3NjcmlwdHRlbXBsYXRlcy82LjI5LjAvb3RCYW5uZXJTZGsuanM6NzoxNzM1KQogICAgYXQgaHR0cHM6Ly9jZG4uY29va2llbGF3Lm9yZy9zY3JpcHR0ZW1wbGF0ZXMvNi4yOS4wL290QmFubmVyU2RrLmpzOjc6Njcz",
		Prm:                true,
		Tzp:                "Europe/Rome",
		Cvs:                true,
		Usb:                "defined",
		TimeSeconds:        (t.UnixMilli() / 1000) - 1000,
	}

	jsonData, err := json.Marshal(jsData)

	if err != nil {
		return param, fmt.Errorf("Error while creating the json data!")
	}

	params := []Param{
		{key: "jsData", value: string(jsonData)},
		{key: "eventCounters", value: "[]"},
		{key: "jsType", value: "ch"},
		{key: "cid", value: data.Cid},
		{key: "ddk", value: data.DataDomeKey},
		{key: "Referer", value: data.Refer},
		{key: "request", value: data.Path},
		{key: "responsePage", value: "origin"},
		{key: "ddv", value: data.DataDomeVersion},
	}

	return params, nil
}

func (ddModule *DataDomeSession) RequestCookie() (DDResponse, error) {
	var ddRes DDResponse

	req, err := http.NewRequest(http.MethodPost, ddModule.DataDomeDomain, nil)

	bodyPayload := &strings.Builder{}
	for i, p := range ddModule.Payload {
		if i > 0 {
			bodyPayload.WriteByte('&')
		}
		bodyPayload.WriteString(p.key)
		bodyPayload.WriteByte('=')
		bodyPayload.WriteString(url.QueryEscape(p.value))
	}

	req.Body = io.NopCloser(strings.NewReader(bodyPayload.String()))
	req.ContentLength = int64(len(bodyPayload.String()))

	if err != nil {
		return ddRes, err
	}

	req.Header = http.Header{
		"authority":          {"api-js.datadome.co"},
		"accept":             {"*/*"},
		"accept-language":    {"it-IT,it;q=0.9,en-IT;q=0.8,en;q=0.7,si-LK;q=0.6,si;q=0.5,en-US;q=0.4"},
		"content-type":       {"application/x-www-form-urlencoded"},
		"origin":             {ddModule.GenData.Refer},
		"referer":            {ddModule.GenData.Refer},
		"sec-ch-ua":          {"\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\""},
		"sec-ch-ua-mobile":   {"?0"},
		"sec-ch-ua-platform": {"\"macOS\""},
		"sec-fetch-dest":     {"empty"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-site":     {"cross-site"},
		"user-agent":         {ddModule.GenData.UserAgent},
	}

	res, err := ddModule.httpClient.Do(req)

	if err != nil {
		return ddRes, err
	}

	if res != nil {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		return ddRes, fmt.Errorf("Generate datadome bad status code!")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ddRes, err
	}

	json.Unmarshal(body, &ddRes)
	return ddRes, nil
}
