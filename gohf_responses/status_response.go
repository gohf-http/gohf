package gohf_responses

import (
	"context"
	"errors"
	"net/http"

	"github.com/gohf-http/gohf/v3"
)

type StatusResponse struct {
	Status int
}

func NewStatusResponse(statusCode int) StatusResponse {
	return StatusResponse{
		Status: statusCode,
	}
}

func (response StatusResponse) Send(res gohf.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	res.SetHeader("Content-Type", "text/plain")
	res.SetStatus(response.Status)
	//nolint:errcheck
	res.Write([]byte(http.StatusText(response.Status)))
}
