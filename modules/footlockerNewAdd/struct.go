package ftlNewAdd

import tls_client "github.com/bogdanfinn/tls-client"

// footlocker struct

type footlockerSession struct {
	httpClient   tls_client.HttpClient
	Store_Region string
	ProxyArr     []string `json:"proxy"`
	ProxyFile    string
	Host         string
	UserAgent    string
	DataDome     string
	Cid          string
	Proxy        string
}
