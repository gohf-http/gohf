package gohf

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type textResponse struct {
	status int
	text   string
}

func (res textResponse) Send(w http.ResponseWriter, req *Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(res.status)
	//nolint:errcheck
	w.Write([]byte(res.text))
}

func TestHelloGoHF(t *testing.T) {
	type result struct {
		body   string
		status int
	}
	testcases := map[string]result{
		"/": {
			body:   "Page not found",
			status: http.StatusNotFound,
		},
		"/asd": {
			body:   "Page not found",
			status: http.StatusNotFound,
		},
		"/greeting": {
			body:   "Name is required",
			status: http.StatusBadRequest,
		},
		"/greeting?name": {
			body:   "Name is required",
			status: http.StatusBadRequest,
		},
		"/greeting?name=": {
			body:   "Name is required",
			status: http.StatusBadRequest,
		},
		"/greeting?name=GoHF": {
			body:   "Hello, GoHF!",
			status: http.StatusOK,
		},
	}

	router := New()

	router.Handle("GET /greeting", func(c *Context) Response {
		name := c.Req.GetQuery("name")
		if name == "" {
			return textResponse{
				http.StatusBadRequest,
				"Name is required",
			}
		}

		greeting := fmt.Sprintf("Hello, %s!", name)
		return textResponse{http.StatusOK, greeting}
	})

	router.Use(func(c *Context) Response {
		return textResponse{
			http.StatusNotFound,
			"Page not found",
		}
	})

	mux := router.CreateServeMux()

	for url, result := range testcases {
		name := fmt.Sprintf("Test Hello GoHF: %s", url)
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			data, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("unexpected error %v", err)
			}
			if got, want := string(data), result.body; got != want {
				t.Errorf("body error. got:\"%v\" want:\"%v\"", got, want)
			}
			if got, want := res.StatusCode, result.status; got != want {
				t.Errorf("status error. got:\"%v\" want:\"%v\"", got, want)
			}
		})
	}
}
