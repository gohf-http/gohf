package gohf_responses

import (
	"context"
	"errors"
	"net/http"

	"github.com/gohf-http/gohf/v5"
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

func (res RedirectResponse) Send(w http.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	http.Redirect(
		w,
		req.GetHttpRequest(),
		res.Url,
		res.Status,
	)
}
