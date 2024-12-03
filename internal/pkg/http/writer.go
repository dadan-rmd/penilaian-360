package http

import (
	"encoding/json"
	"net/http"
)

type ResponseSuccess struct {
	ResponseCode   string      `json:"responseCode"`
	ResponseDesc   string      `json:"responseDescription"`
	Data           interface{} `json:"data"`
	AdditionalInfo interface{} `json:"additionalInfo"`
}

type emptyStruct struct{}

type ResponseError struct {
	ResponseCode string   `json:"Code"`
	ExternalID   string   `json:"externalId"`
	Message      string   `json:"Message"`
	Errors       []string `json:"Errors"`
}

type ErrorMapping struct {
	ResponseCode string
	Message      string
	HTTPStatus   int
}

const (
	successCode  = "0000000"
	successDesc  = "Success"
	defErrorCode = "9999999"
	defErrorMsg  = "Ups, terjadi kesalahan"
)

type HttpHandlerContext struct {
	ErrMap map[error]ErrorMapping
}

func NewContextHandler() HttpHandlerContext {
	return HttpHandlerContext{
		ErrMap: map[error]ErrorMapping{},
	}
}

func (hctx HttpHandlerContext) AddErr(key error, val ErrorMapping) {
	hctx.ErrMap[key] = val
}

func (hctx HttpHandlerContext) LookupErr(key error) (val ErrorMapping) {
	if e, ok := hctx.ErrMap[key]; ok {
		val = e
	} else {
		val = ErrorMapping{
			ResponseCode: defErrorCode,
			Message:      defErrorMsg,
			HTTPStatus:   http.StatusInternalServerError,
		}
	}

	return
}

type CustomWriter struct {
	C HttpHandlerContext
}

func (c *CustomWriter) Write(w http.ResponseWriter, data interface{}, additionalInfo interface{}) {
	if additionalInfo == nil {
		additionalInfo = &emptyStruct{}
	}
	resp := ResponseSuccess{
		ResponseCode:   successCode,
		ResponseDesc:   successDesc,
		Data:           data,
		AdditionalInfo: additionalInfo,
	}

	writeResponse(w, resp, http.StatusOK)
}

// WriteError sending error response based on err type
func (c *CustomWriter) WriteError(w http.ResponseWriter, extID string, err error, errCustom []string) {
	if len(errCustom) == 0 {
		errCustom = make([]string, 0)
	}

	respMap := c.C.LookupErr(err)
	respErr := ResponseError{
		ResponseCode: respMap.ResponseCode,
		ExternalID:   extID,
		Message:      respMap.Message,
		Errors:       errCustom,
	}

	writeResponse(w, respErr, respMap.HTTPStatus)
}

func writeResponse(w http.ResponseWriter, response interface{}, httpStatus int) {
	res, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error parsing body response"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(res)
}
