package ftlkiller

import (
	"strconv"
	"strings"
	"time"

	"github.com/TheAldriguzDev/Bonsai_solutions/packages/log"
	"github.com/TheAldriguzDev/Bonsai_solutions/solvers/datadome"
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

func (ftl *footlockerSession) GetDataDome() string {
	var ddStep bool = false
	var ddRes datadome.DDResponse

	ddData := datadome.GenerateData{
		DataDomeKey:     "A55FBF4311ED6F1BF9911EB71931D5",
		DataDomeVersion: "4.10.2",
		Cid:             ftl.Cid,
		Refer:           "https://www.footlocker.it/",
		UserAgent:       "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36",
		Proxy:           ftl.Proxy,
	}

	for !ddStep {
		ddResTemp, err := datadome.GetCookie(ddData)

		if err != nil {
			log.Error(err.Error(), location)
			time.Sleep(5 * time.Second)
			continue
		}

		ddRes = ddResTemp
		ddStep = true
	}

	newCookie := strings.Split(strings.Split(ddRes.Cookie, ";")[0], "=")[1]
	return newCookie
}
