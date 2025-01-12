package gohf

import "net/http"

type Response interface {
	Send(http.ResponseWriter, *Request)
}

type HandlerFunc func(c *Context) Response

type NextFunc func() Response

type Context struct {
	w    http.ResponseWriter
	Req  *Request
	Next NextFunc
}

func newContext(w http.ResponseWriter, r *Request) *Context {
	return &Context{
		w:    w,
		Req:  r,
		Next: func() Response { return nil },
	}
}

func (c *Context) ResHeader() http.Header {
	return c.w.Header()
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.w, cookie)
}
