package gohf_responses

import (
	"net/http"

	"github.com/gohf-http/gohf/v4"
)

type DummyResponse struct {
}

func NewDummyResponse() DummyResponse {
	return DummyResponse{}
}

func (res DummyResponse) Send(_ http.ResponseWriter, _ *gohf.Request) {}
