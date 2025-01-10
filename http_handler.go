package gohf

import (
	"net/http"
)

type httpHandler struct {
	handlerFuncs []HandlerFunc
}

func newHttpHandler() *httpHandler {
	return &httpHandler{}
}

func (httpHandler *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := newResponseWriter(w)
	req := newRequest(r)
	c := newContext(res, req)

	var handle func(idx int) Response
	handle = func(idx int) Response {
		if idx == len(httpHandler.handlerFuncs) {
			return nil
		}

		c.Next = func() Response { return handle(idx + 1) }

		return httpHandler.handlerFuncs[idx](c)
	}

	if response := handle(0); response != nil {
		response.Send(c.Res, c.Req)
	}
}

func (httpHandler *httpHandler) addHandlerFunc(handlerFunc HandlerFunc) {
	httpHandler.handlerFuncs = append(
		httpHandler.handlerFuncs,
		handlerFunc,
	)
}

func (httpHandler *httpHandler) len() int {
	return len(httpHandler.handlerFuncs)
}
