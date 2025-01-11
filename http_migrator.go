package gohf

import "net/http"

func FromHttpHandleFunc(httpHandleFunc func(http.ResponseWriter, *http.Request)) HandlerFunc {
	return func(c *Context) Response {
		httpHandleFunc(c.Res.GetHttpResponseWriter(), c.Req.GetHttpRequest())
		return dummyResponse{}
	}
}

func FromHttpHandler(httpHandler http.Handler) HandlerFunc {
	return func(c *Context) Response {
		httpHandler.ServeHTTP(c.Res.GetHttpResponseWriter(), c.Req.GetHttpRequest())
		return dummyResponse{}
	}
}

type dummyResponse struct{}

func (response dummyResponse) Send(_ ResponseWriter, _ *Request) {}
