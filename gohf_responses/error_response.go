package gohf_responses

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gohf-http/gohf/v4"
)

type ErrorResponse[T interface{}] struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Err     T      `json:"err"`
}

func NewErrorResponse[TError error](statusCode int, err TError) ErrorResponse[TError] {
	return ErrorResponse[TError]{
		Status:  statusCode,
		Message: err.Error(),
		Err:     err,
	}
}

func (res ErrorResponse[T]) Error() string {
	return fmt.Sprintf("http error %d: %s", res.Status, res.Message)
}

func (res ErrorResponse[T]) Send(w http.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Status)
	//nolint:errcheck
	json.NewEncoder(w).Encode(res)
}
