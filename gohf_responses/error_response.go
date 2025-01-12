package gohf_responses

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gohf-http/gohf/v5"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Err     error  `json:"err"`
}

func NewErrorResponse(statusCode int, err error) ErrorResponse {
	return ErrorResponse{
		Status:  statusCode,
		Message: err.Error(),
		Err:     err,
	}
}

func (res ErrorResponse) Error() string {
	return fmt.Sprintf("http error %d: %s", res.Status, res.Message)
}

func (res ErrorResponse) Send(w http.ResponseWriter, req *gohf.Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Status)
	//nolint:errcheck
	json.NewEncoder(w).Encode(res)
}
