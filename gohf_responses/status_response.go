package gohf_responses

import (
	"context"
	"errors"
	"net/http"

	"github.com/gohf-http/gohf/v5"
)

type StatusResponse struct {
	Status int
}

func NewStatusResponse(statusCode int) StatusResponse {
	return StatusResponse{
		Status: statusCode,
	}
}

func (res StatusResponse) Send(w http.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(res.Status)
	//nolint:errcheck
	w.Write([]byte(http.StatusText(res.Status)))
}
