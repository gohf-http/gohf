package gohf_responses

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gohf-http/gohf/v5"
)

func TestErrorResponse(t *testing.T) {
	type testcase struct {
		status int
		err    error
	}
	type result struct {
		status  int
		message string
	}

	testcases := map[testcase]result{
		{http.StatusBadRequest, errors.New("test error")}: {
			status:  http.StatusBadRequest,
			message: "test error",
		},
		{http.StatusNotFound, errors.New("not found error")}: {
			status:  http.StatusNotFound,
			message: "not found error",
		},
		{http.StatusConflict, errors.New("conflict error")}: {
			status:  http.StatusConflict,
			message: "conflict error",
		},
	}

	createServerMux := func(status int, err error) *http.ServeMux {
		router := gohf.New()
		router.Use(func(c *gohf.Context) gohf.Response {
			return NewErrorResponse(status, err)
		})
		return router.CreateServeMux()
	}

	getJsonResponseAndStatus := func(mux *http.ServeMux, dest interface{}) (int, error) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()
		err := json.NewDecoder(res.Body).Decode(dest)
		status := res.StatusCode
		return status, err
	}

	for testcase, result := range testcases {
		name := fmt.Sprintf("test error response: %v %s", testcase.status, testcase.err.Error())
		t.Run(name, func(t *testing.T) {
			mux := createServerMux(testcase.status, testcase.err)

			type body struct {
				Status  int    `json:"status"`
				Message string `json:"message"`
			}

			var dest body
			status, err := getJsonResponseAndStatus(mux, &dest)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if got, want := status, result.status; got != want {
				t.Errorf("status - got:%v want:%v", got, want)
			}
			if got, want := dest.Status, result.status; got != want {
				t.Errorf("dest.status - got:%v want:%v", got, want)
			}
			if got, want := dest.Message, result.message; got != want {
				t.Errorf("dest.message - got:%v want:%v", got, want)
			}
		})
	}
}
