package gohf_responses

import (
	"context"
	"errors"
	"net/http"

	"github.com/gohf-http/gohf/v3"
)

type RedirectResponse struct {
	Status int
	Url    string
}

func NewRedirectResponse(statusCode int, url string) RedirectResponse {
	return RedirectResponse{
		Status: statusCode,
		Url:    url,
	}
}

func (response RedirectResponse) Send(res gohf.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	http.Redirect(
		res.GetHttpResponseWriter(),
		req.GetHttpRequest(),
		response.Url,
		response.Status,
	)
}
