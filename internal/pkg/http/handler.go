package http

import (
	"net/http"

	"central-auth/internal/pkg/log"
)

type HandlerOption func(*HttpHandler)

type HttpHandler struct {
	// H is handler, with return interface{} as data object, *string for token next page, error for error type
	H func(w http.ResponseWriter, r *http.Request) (interface{}, interface{}, error, []string)
	CustomWriter
	ServiceName string
}

func NewHttpHandler(c HttpHandlerContext, opts ...HandlerOption) func(handler func(w http.ResponseWriter, r *http.Request) (interface{}, interface{}, error, []string)) HttpHandler {
	return func(handler func(w http.ResponseWriter, r *http.Request) (interface{}, interface{}, error, []string)) HttpHandler {
		h := HttpHandler{H: handler, CustomWriter: CustomWriter{C: c}}

		// Option paremeters values:
		for _, opt := range opts {
			opt(&h)
		}

		return h
	}
}

func (h HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, pageToken, err, errCustom := h.H(w, r)
	extID := r.Header.Get("X-EXTERNAL-ID")

	if err != nil {
		log.Zlogger(r.Context()).Err(err).Msgf("Response: %+v", data)
		h.WriteError(w, extID, err, errCustom)
		return
	}

	h.Write(w, data, pageToken)
}

const paramSign = "PARAM"
