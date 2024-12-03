package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"

	"central-auth/internal/app/commons/loggers"

	"github.com/sony/gobreaker"
	"github.com/spf13/cast"
)

const (
	hidemsisdn = `((628|08)(31|32|33|38|591|598)\d{6,10}|\d{5,13})`
)

// Request struct
type Request struct {
	Method  string
	URL     string
	Payload io.Reader
	Header  map[string]string
	Service string
	Timeout time.Duration
	Breaker *gobreaker.CircuitBreaker
}
type respon struct {
	start       time.Time
	result      []byte
	status      int
	dumprequest string
	message     []string
}

func (r *respon) setMsg(value string) {
	r.message = append(r.message, value)
}

// ComposeRequest is func to create request
func ComposeRequest(r Request) (*http.Request, error) {
	req, err := http.NewRequest(r.Method, r.URL, r.Payload)
	if err != nil {
		// logger.Errorf("could not make request %v", err)
		return nil, err
	}

	if r.Header != nil {
		for k, v := range r.Header {
			req.Header.Set(k, v)
		}
	}
	return req, nil
}

/*Curl is func
 *
 */
func Curl(record *loggers.Data, r Request, temp interface{}) (interface{}, error) {
	ctx := context.Background()
	loc, _ := time.LoadLocation("Asia/Jakarta")
	d := new(respon)
	d.start = time.Now().In(loc)

	req, err := ComposeRequest(r)
	if err != nil {

		d.result = nil
		d.status = http.StatusInternalServerError
		loggers.RecordThridPartyFailed(record, nil, d.start, r.Service, http.StatusInternalServerError, r.Payload, err.Error())
		return d, err
	}

	//dump request
	d.dumprequest = dumpRequest(req)
	cons, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	res, err := defclient.Do(req.WithContext(cons))
	if err != nil {

		d.result = nil
		d.status = http.StatusInternalServerError

		loggers.RecordThridPartyFailed(record, req, d.start, r.Service, d.status, r.Payload, err.Error())

		return d, err
	}
	defer res.Body.Close()

	p, err := selectReader(res.Body, temp)
	if err != nil {

		d.result = nil
		d.status = http.StatusInternalServerError

		loggers.RecordThridPartyFailed(record, req, d.start, r.Service, d.status, r.Payload, err.Error())

		return d, err
	}
	d.result = p
	d.status = res.StatusCode

	loggers.RecordThridParty(record, req, d.start, r.Service, d.status, r.Payload, d.result)

	return d, nil
}

/*chooseBreaker dengan circuit breaker
 * @author : mff
 */
func chooseBreaker(record *loggers.Data, r Request, temp interface{}) ([]byte, int, error) {
	var (
		body interface{}
		err  error
	)
	rs := new(respon)
	rs.start = time.Now()
	body, err = Curl(record, r, temp)
	rs = body.(*respon)
	(cast.ToString(rs.result))

	if err != nil {

		rs.status = http.StatusInternalServerError
		rs.result = []byte(err.Error())

		return nil, rs.status, err
	}
	return rs.result, rs.status, nil
}

func dumpRequest(r *http.Request) string {
	dump, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		return ""
	}
	trim := bytes.ReplaceAll(dump, []byte("\r\n"), []byte("   "))

	data := string(trim)

	return string(data)
}

func recordToPrometheus(d time.Duration, method, url, status string) {
	loggers.Getprometheus().MetricRecord(status, method, url, loggers.GetName(), d)
}

func setBreaker(b *gobreaker.CircuitBreaker) *gobreaker.CircuitBreaker {
	flags := "true"
	switch flags {
	case "true":
		return nil
	}
	return b
}

// ReadBody function to read body and marshal
func readBody(result io.ReadCloser, temp interface{}) ([]byte, error) {

	err := json.NewDecoder(result).Decode(&temp)
	if err != nil {
		return nil, err
	}

	response, err := json.Marshal(&temp)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// DoRequest is wrap function choose breaker
func (r Request) DoRequest(record *loggers.Data) ([]byte, int, error) {
	return chooseBreaker(record, r, nil)
}

func selectReader(result io.ReadCloser, temp interface{}) ([]byte, error) {
	if temp != nil {
		return readBody(result, temp)
	}
	return ioutil.ReadAll(result)
}
