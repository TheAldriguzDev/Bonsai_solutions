package ftlRestock

import (
	"strconv"
	"time"

	http "github.com/bogdanfinn/fhttp"

	"github.com/google/uuid"
)

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

func hasDatadomeCookie(req *http.Request) bool {
	cookies := req.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "datadome" {
			return true
		}
	}
	return false
}
