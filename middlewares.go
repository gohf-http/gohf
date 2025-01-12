package gohf

import "net/http"

func MaxBytesMiddleware(n int64) HandlerFunc {
	return func(c *Context) Response {
		c.Req.body = newRequestBody(http.MaxBytesReader(c.w, c.Req.body, n))
		return c.Next()
	}
}
