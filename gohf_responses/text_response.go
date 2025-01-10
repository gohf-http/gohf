package gohf_responses

import (
	"context"
	"errors"

	"github.com/gohf-http/gohf"
)

type TextResponse struct {
	Status int
	Text   string
}

func NewTextResponse(statusCode int, text string) TextResponse {
	return TextResponse{
		Status: statusCode,
		Text:   text,
	}
}

func (response TextResponse) Send(res gohf.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	res.SetHeader("Content-Type", "text/plain")
	res.SetStatus(response.Status)
	//nolint:errcheck
	res.Write([]byte(response.Text))
}
