package response

import (
	"context"
	"errors"
	"net/http"

	"github.com/gohf-http/gohf/v6"
)

type ServeFileResponse struct {
	Filepath string
}

func ServeFile(filepath string) ServeFileResponse {
	return ServeFileResponse{
		Filepath: filepath,
	}
}

func (res ServeFileResponse) Send(w http.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	http.ServeFile(w, req.GetHttpRequest(), res.Filepath)
}
