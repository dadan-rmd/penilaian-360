package loggers

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// RecordThridParty ...
func RecordThridParty(record *Data, req *http.Request, start time.Time, service string, status int, body io.Reader, response []byte) {
	var (
		payload string
	)
	t := time.Since(start)
	if body != nil {

		bd, _ := req.GetBody()
		reqBody, err := ioutil.ReadAll(bd)
		if err != nil {
			payload = ""
		} else {

			payload = string(reqBody)
			req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
		}
	}
	third := ThirdParty{}

	third.Service = service
	third.URL = req.Host + req.URL.Path
	third.Response = string(response)
	third.StatusCode = status
	third.RequestHeader = DumpRequest(req)
	third.RequestBody = payload
	third.ExecTime = t.Seconds()
	third.RequestID = record.RequestID
	record.ThirdParty = append(record.ThirdParty, third)

}

// RecordThridPartyFailed ...
func RecordThridPartyFailed(record *Data, req *http.Request, start time.Time, service string, status int, body io.Reader, messages string) {
	var (
		url     = req.Host + req.URL.Path
		payload string
	)

	t := time.Since(start)
	if req == nil {
		url = ""
	}

	if body != nil {
		bd, _ := req.GetBody()
		reqBody, err := ioutil.ReadAll(bd)
		if err != nil {
			payload = ""
		} else {
			payload = string(reqBody)
			req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
		}
	}
	third := ThirdParty{}

	third.Service = service
	third.URL = url
	third.Response = messages
	third.StatusCode = status
	third.RequestHeader = DumpRequest(req)
	third.RequestBody = payload
	third.ExecTime = t.Seconds()

	third.RequestID = record.RequestID
	record.ThirdParty = append(record.ThirdParty, third)
}
