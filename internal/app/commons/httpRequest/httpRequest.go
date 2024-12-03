package httpRequest

import (
	"bytes"
	"strconv"
	"time"

	"central-auth/internal/app/commons/loggers"
	"central-auth/internal/app/commons/utils"
)

func Curl(record *loggers.Data, serviceName, method string, url string, body []byte, arrheader map[string]string) ([]byte, int, error) {
	var rs utils.Request
	urls, timer := url, "30"

	rs.Service = serviceName
	rs.Method = method
	rs.URL = urls
	rs.Header = arrheader
	t, _ := strconv.Atoi(timer)
	rs.Payload = bytes.NewBuffer(body)
	rs.Timeout = time.Duration(t) * time.Second
	return rs.DoRequest(record)
}
