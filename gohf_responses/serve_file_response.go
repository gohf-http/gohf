package gohf_responses

import (
	"context"
	"errors"
	"net/http"

	"github.com/gohf-http/gohf/v3"
)

type ServeFileResponse struct {
	Filepath string
}

func NewServeFileResponse(filepath string) ServeFileResponse {
	return ServeFileResponse{
		Filepath: filepath,
	}
}

func (response ServeFileResponse) Send(res gohf.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	http.ServeFile(res.GetHttpResponseWriter(), req.GetHttpRequest(), response.Filepath)
}
