package gohf_responses

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gohf-http/gohf"
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

func (response ErrorResponse[T]) Error() string {
	return fmt.Sprintf("default http error [%d]: %s", response.Status, response.Message)
}

func (response ErrorResponse[T]) Send(res gohf.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	res.SetHeader("Content-Type", "application/json")
	res.SetStatus(response.Status)
	//nolint:errcheck
	json.NewEncoder(res).Encode(response)
}
