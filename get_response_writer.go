package gohf

import "net/http"

func GetResponseWriter(c *Context) http.ResponseWriter {
	return c.w
}
