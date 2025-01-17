package response

import (
	"context"
	"errors"
	"net/http"

	"github.com/gohf-http/gohf/v6"
)

type TextResponse struct {
	Status int
	Text   string
}

func Text(statusCode int, text string) TextResponse {
	return TextResponse{
		Status: statusCode,
		Text:   text,
	}
}

func (res TextResponse) Send(w http.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(res.Status)
	//nolint:errcheck
	w.Write([]byte(res.Text))
}
