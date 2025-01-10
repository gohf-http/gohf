package gohf

import (
	"io"
	"net/http"
	"time"
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

func (res ResponseWriter) ServeFile(req *Request, filepath string) {
	http.ServeFile(res.w, req.req, filepath)
}

func (res ResponseWriter) ServeContent(req *Request, name string, modtime time.Time, content io.ReadSeeker) {
	http.ServeContent(res.w, req.req, name, modtime, content)
}
