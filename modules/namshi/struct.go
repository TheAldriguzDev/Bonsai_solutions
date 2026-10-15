package namshi

import (
	"strconv"
	"time"

	tls_client "github.com/bogdanfinn/tls-client"
)

type namshiSession struct {
	httpClient   tls_client.HttpClient
	Store_Region string
	ProxyArr     []string `json:"proxy"`
	ProxyFile    string
	Host         string
	UserAgent    string
	Proxy        string
}

func GetTimeStamp() string {
	currentTime := time.Now()
	timestamp := currentTime.UnixNano() / int64(time.Millisecond)
	cacheBuster := strconv.Itoa(int(timestamp))
	return cacheBuster
}
