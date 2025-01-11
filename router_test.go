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

func (response textResponse) Send(res ResponseWriter, req *Request) {
	if errors.Is(req.RootContext().Err(), context.Canceled) {
		return
	}

	res.SetHeader("Content-Type", "text/plain")
	res.SetStatus(response.status)
	//nolint:errcheck
	res.Write([]byte(response.text))
}

func TestRouter(t *testing.T) {
	sendRequest := func(handler http.Handler, method string, target string) ([]byte, error) {
		req := httptest.NewRequest(method, target, nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()
		return io.ReadAll(res.Body)
	}

	t.Run("should response \"Hello, GoHF!\"", func(t *testing.T) {
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

		data, err := sendRequest(mux, http.MethodGet, "/greeting?name=GoHF")

		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if got, want := string(data), "Hello, GoHF!"; got != want {
			t.Errorf("got:\"%s\" want:\"%s\"", got, want)
		}
	})

	t.Run("should response \"Name is required\"", func(t *testing.T) {
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

		data, err := sendRequest(mux, http.MethodGet, "/greeting")

		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if got, want := string(data), "Name is required"; got != want {
			t.Errorf("got:\"%s\" want:\"%s\"", got, want)
		}
	})

	t.Run("should response \"Page not found\"", func(t *testing.T) {
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

		data, err := sendRequest(mux, http.MethodGet, "/greeting2?name=GoHF")

		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if got, want := string(data), "Page not found"; got != want {
			t.Errorf("got:\"%s\" want:\"%s\"", got, want)
		}
	})
}
