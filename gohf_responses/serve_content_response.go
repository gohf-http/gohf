package gohf_responses

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gohf-http/gohf/v3"
)

type ServeContentResponse struct {
	Name    string
	Modtime time.Time
	Content io.ReadSeeker
}

func NewServeContentResponse(name string, modtime time.Time, content io.ReadSeeker) ServeContentResponse {
	return ServeContentResponse{
		Name:    name,
		Modtime: modtime,
		Content: content,
	}
}

func (response ServeContentResponse) Send(res gohf.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	http.ServeContent(
		res.GetHttpResponseWriter(),
		req.GetHttpRequest(),
		response.Name,
		response.Modtime,
		response.Content,
	)
}
