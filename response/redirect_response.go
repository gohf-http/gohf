package response

import (
	"context"
	"errors"
	"net/http"

	"github.com/gohf-http/gohf/v6"
)

type RedirectResponse struct {
	Status int
	Url    string
}

func Redirect(statusCode int, url string) RedirectResponse {
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
