package gohf_responses

import "github.com/gohf-http/gohf/v2"

type DummyResponse struct {
}

func NewDummyResponse() DummyResponse {
	return DummyResponse{}
}

func (response DummyResponse) Send(_ gohf.ResponseWriter, _ *gohf.Request) {}
