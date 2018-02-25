package test

import (
	"net/http"
	"net/http/httptest"
)

type ResponseWriter struct {
	HandlerInfo HandlerInfo
	URLVars     map[string]string
	W           *httptest.ResponseRecorder
	r           *http.Request
}

func NewResponseWriter(w *httptest.ResponseRecorder, r *http.Request) *ResponseWriter {
	return &ResponseWriter{
		W: w,
		r: r,
	}
}

func (rw *ResponseWriter) Header() http.Header {
	return rw.W.Header()
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.setHandlerInfo()
	return rw.W.Write(b)
}

func (rw *ResponseWriter) WriteHeader(c int) {
	rw.W.WriteHeader(c)
}

func (rw *ResponseWriter) setHandlerInfo() {
	handlerInfo := infoFunc(rw.r)

	if handlerInfo == nil {
		return
	}

	rw.HandlerInfo = *handlerInfo
}
