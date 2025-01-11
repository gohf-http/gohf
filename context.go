package gohf

import "net/http"

type Response interface {
	Send(http.ResponseWriter, *Request)
}

type HandlerFunc func(c *Context) Response

type NextFunc func() Response

type Context struct {
	w         http.ResponseWriter
	ResHeader http.Header
	Req       *Request
	Next      NextFunc
}

func newContext(w http.ResponseWriter, r *Request) *Context {
	return &Context{
		w:         w,
		ResHeader: w.Header(),
		Req:       r,
		Next:      func() Response { return nil },
	}
}
