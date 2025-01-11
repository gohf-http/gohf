package gohf

import (
	"net/http"
)

type ResponseWriter struct {
	w http.ResponseWriter
}

func newResponseWriter(w http.ResponseWriter) ResponseWriter {
	return ResponseWriter{
		w: w,
	}
}

func (res ResponseWriter) SetHeader(key, value string) {
	res.w.Header().Set(key, value)
}

func (res ResponseWriter) SetStatus(statusCode int) {
	res.w.WriteHeader(statusCode)
}

func (res ResponseWriter) Write(p []byte) (n int, err error) {
	return res.w.Write(p)
}

func (res ResponseWriter) GetHttpResponseWriter() http.ResponseWriter {
	return res.w
}
