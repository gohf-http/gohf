package response

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gohf-http/gohf/v6"
)

type ServeContentResponse struct {
	Name    string
	Modtime time.Time
	Content io.ReadSeeker
}

func ServeContent(name string, modtime time.Time, content io.ReadSeeker) ServeContentResponse {
	return ServeContentResponse{
		Name:    name,
		Modtime: modtime,
		Content: content,
	}
}

func (res ServeContentResponse) Send(w http.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	http.ServeContent(
		w,
		req.GetHttpRequest(),
		res.Name,
		res.Modtime,
		res.Content,
	)
}
