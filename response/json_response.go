package response

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gohf-http/gohf/v6"
)

type JsonResponse struct {
	Status int
	Data   interface{}
}

func JSON(statusCode int, data interface{}) JsonResponse {
	return JsonResponse{
		Status: statusCode,
		Data:   data,
	}
}

func (res JsonResponse) Send(w http.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Status)
	//nolint:errcheck
	json.NewEncoder(w).Encode(res.Data)
}
