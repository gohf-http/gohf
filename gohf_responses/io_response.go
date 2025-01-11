package gohf_responses

import (
	"context"
	"errors"
	"io"

	"github.com/gohf-http/gohf/v3"
)

type IoResponse struct {
	Status int
	Reader io.Reader
}

func NewIoResponse(statusCode int, reader io.Reader) IoResponse {
	return IoResponse{
		Status: statusCode,
		Reader: reader,
	}
}

func (response IoResponse) Send(res gohf.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	res.SetStatus(response.Status)
	//nolint:errcheck
	io.Copy(res, response.Reader)
}
