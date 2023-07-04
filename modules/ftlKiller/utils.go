package ftlkiller

import (
	"strconv"
	"time"

	http "github.com/bogdanfinn/fhttp"
	"github.com/google/uuid"
)

type BypassHeader struct {
	Host      string
	UserAgent string
	Refer     string
	DataDome  string
}

func GetBypassHeader(info BypassHeader) http.Header {
	headers := http.Header{
		"Host":             {info.Host},
		"user-agent":       {info.UserAgent},
		"Accept":           {"application/json, text/javascript, */*; q=0.01"},
		"Accept-Language":  {"de-DE,de;q=0.9"},
		"Accept-Encoding":  {"gzip, deflate"},
		"Content-Type":     {"application/x-www-form-urlencoded"},
		"X-Requested-With": {"XMLHttpRequest"},
		"Origin":           {info.Host},
		"Referer":          {info.Refer},
		"Sec-Fetch-Dest":   {"empty"},
		"Sec-Fetch-Mode":   {"cors"},
		"Sec-Fetch-Site":   {"same-origin"},
		"Pragma":           {"no-cache"},
		"Cache-Control":    {"no-cache"},
		"Te":               {"trailers"},
		"cookie":           {"datadome=" + info.DataDome},
	}

	return headers
}

func GenerateUUID() string {
	uuidV4 := uuid.New()
	return uuidV4.String()
}

func GetTimeStamp() string {
	currentTime := time.Now()
	timestamp := currentTime.UnixNano() / int64(time.Millisecond)
	cacheBuster := strconv.Itoa(int(timestamp))
	return cacheBuster
}
