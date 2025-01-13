package response

import (
	"net/http"

	"github.com/gohf-http/gohf/v6"
)

type DummyResponse struct {
}

func Dummy() DummyResponse {
	return DummyResponse{}
}

func (res DummyResponse) Send(_ http.ResponseWriter, _ *gohf.Request) {}
