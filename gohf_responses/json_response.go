package gohf_responses

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gohf-http/gohf/v4"
)

type JsonResponse[T interface{}] struct {
	Status int
	Data   T
}

func NewJsonResponse[T interface{}](statusCode int, data T) JsonResponse[T] {
	return JsonResponse[T]{
		Status: statusCode,
		Data:   data,
	}
}

func (res JsonResponse[T]) Send(w http.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Status)
	//nolint:errcheck
	json.NewEncoder(w).Encode(res.Data)
}
