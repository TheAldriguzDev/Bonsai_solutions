package zalando

import tls_client "github.com/bogdanfinn/tls-client"

type zalandoSession struct {
	httpClient   tls_client.HttpClient
	Store_Region string
	ProxyArr     string `json:"proxy"`
	ProxyFile    string
	Host         string
	UserAgent    string
	Proxy        string
}
