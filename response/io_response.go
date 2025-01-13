package response

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/gohf-http/gohf/v6"
)

type IoResponse struct {
	Status int
	Reader io.Reader
}

func IO(statusCode int, reader io.Reader) IoResponse {
	return IoResponse{
		Status: statusCode,
		Reader: reader,
	}
}

func (res IoResponse) Send(w http.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	w.WriteHeader(res.Status)
	//nolint:errcheck
	io.Copy(w, res.Reader)
}
