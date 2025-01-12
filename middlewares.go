package gohf

import "net/http"

func MaxBytesMiddleware(n int64) HandlerFunc {
	return func(c *Context) Response {
		w := GetResponseWriter(c)
		maxBytesBody := http.MaxBytesReader(w, c.Req.req.Body, n)
		c.Req.req.Body = maxBytesBody
		c.Req.body = newRequestBody(maxBytesBody)
		return c.Next()
	}
}
