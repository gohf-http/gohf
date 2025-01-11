package gohf

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type Request struct {
	req         *http.Request
	body        RequestBody
	timestamp   time.Time
	ctx         context.Context
	rootContext context.Context
}

func newRequest(req *http.Request) *Request {
	ctx := req.Context()
	return &Request{
		req:         req,
		body:        newRequestBody(req.Body),
		timestamp:   time.Now(),
		ctx:         ctx,
		rootContext: ctx,
	}
}

func (req *Request) GetTimestamp() time.Time {
	return req.timestamp
}

func (req *Request) Method() string {
	return req.req.Method
}

func (req *Request) RemoteAddr() string {
	return req.req.RemoteAddr
}

func (req *Request) Host() string {
	return req.req.Host
}

func (req *Request) RequestURI() string {
	return req.req.RequestURI
}

func (req *Request) RootContext() context.Context {
	return req.rootContext
}

func (req *Request) Context() context.Context {
	return req.ctx
}

func (req *Request) SetContext(ctx context.Context) {
	req.ctx = ctx
}

func (req *Request) GetHeader(key string) string {
	return req.req.Header.Get(key)
}

func (req *Request) PathValue(name string) string {
	return req.req.PathValue(name)
}

func (req *Request) GetQuery(key string) string {
	return req.req.URL.Query().Get(key)
}

func (req *Request) GetBody() RequestBody {
	return req.body
}

func (req *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return req.req.FormFile(key)
}

func (req *Request) FormValue(key string) string {
	return req.req.FormValue(key)
}

func (req *Request) Cookies() []*http.Cookie {
	return req.req.Cookies()
}

func (req *Request) Cookie(name string) (*http.Cookie, error) {
	return req.req.Cookie(name)
}

func (req *Request) AddCookie(c *http.Cookie) {
	req.req.AddCookie(c)
}

func (req *Request) GetHttpRequest() *http.Request {
	return req.req
}

type RequestBody struct {
	body io.ReadCloser
}

func newRequestBody(body io.ReadCloser) RequestBody {
	return RequestBody{body: body}
}

func (body RequestBody) Close() error {
	return body.body.Close()
}

func (body RequestBody) Read(p []byte) (n int, err error) {
	return body.body.Read(p)
}

func (body RequestBody) JsonDecode(v any) error {
	return json.NewDecoder(body).Decode(v)
}
