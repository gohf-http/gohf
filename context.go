package gohf

type Response interface {
	Send(ResponseWriter, *Request)
}

type HandlerFunc func(c *Context) Response

type NextFunc func() Response

type Context struct {
	Res  ResponseWriter
	Req  *Request
	Next NextFunc
}

func newContext(res ResponseWriter, req *Request) *Context {
	return &Context{
		Res:  res,
		Req:  req,
		Next: func() Response { return nil },
	}
}
