package gohf

import "net/http"

func MaxBytesMiddleware(n int64) HandlerFunc {
	return func(c *Context) Response {
		w := GetResponseWriter(c)
		httpReq := c.Req.GetHttpRequest()

		maxBytesBody := http.MaxBytesReader(w, httpReq.Body, n)

		httpReq.Body = maxBytesBody
		c.Req.body = newRequestBody(maxBytesBody)
		return c.Next()
	}
}
