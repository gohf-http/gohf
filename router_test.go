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

func TestNestedRouter(t *testing.T) {
	type result struct {
		status         int
		test1          bool
		test2          bool
		test3          bool
		test4          bool
		test5          bool
		test1notfound  bool
		routernotfound bool
	}
	testcases := map[string]result{
		"/": {
			status:         http.StatusNotFound,
			test1:          false,
			test2:          false,
			test3:          false,
			test4:          false,
			test5:          false,
			test1notfound:  false,
			routernotfound: true,
		},
		"/test": {
			status:         http.StatusNotFound,
			test1:          false,
			test2:          false,
			test3:          false,
			test4:          false,
			test5:          false,
			test1notfound:  false,
			routernotfound: true,
		},
		"/test1": {
			status:         http.StatusMovedPermanently,
			test1:          false,
			test2:          false,
			test3:          false,
			test4:          false,
			test5:          false,
			test1notfound:  false,
			routernotfound: false,
		},
		"/test1/": {
			status:         http.StatusNotFound,
			test1:          true,
			test2:          false,
			test3:          false,
			test4:          false,
			test5:          false,
			test1notfound:  true,
			routernotfound: false,
		},
		"/test1/test": {
			status:         http.StatusNotFound,
			test1:          true,
			test2:          false,
			test3:          false,
			test4:          false,
			test5:          false,
			test1notfound:  true,
			routernotfound: false,
		},
		"/test1/test2": {
			status:         http.StatusMovedPermanently,
			test1:          false,
			test2:          false,
			test3:          false,
			test4:          false,
			test5:          false,
			test1notfound:  false,
			routernotfound: false,
		},
		"/test1/test2/": {
			status:         http.StatusNotFound,
			test1:          true,
			test2:          true,
			test3:          false,
			test4:          false,
			test5:          false,
			test1notfound:  true,
			routernotfound: false,
		},
		"/test1/test2/test": {
			status:         http.StatusNotFound,
			test1:          true,
			test2:          true,
			test3:          false,
			test4:          false,
			test5:          false,
			test1notfound:  true,
			routernotfound: false,
		},
		"/test1/test2/test3": {
			status:         http.StatusOK,
			test1:          true,
			test2:          true,
			test3:          true,
			test4:          false,
			test5:          false,
			test1notfound:  false,
			routernotfound: false,
		},
		"/test1/test2/test3/": {
			status:         http.StatusNotFound,
			test1:          true,
			test2:          true,
			test3:          false,
			test4:          false,
			test5:          false,
			test1notfound:  true,
			routernotfound: false,
		},
		"/test1/test2/test3/aaa": {
			status:         http.StatusNotFound,
			test1:          true,
			test2:          true,
			test3:          false,
			test4:          false,
			test5:          false,
			test1notfound:  true,
			routernotfound: false,
		},
		"/test1/test4": {
			status:         http.StatusNotFound,
			test1:          true,
			test2:          false,
			test3:          false,
			test4:          false,
			test5:          false,
			test1notfound:  true,
			routernotfound: false,
		},
		"/test1/test4/asd": {
			status:         http.StatusMovedPermanently,
			test1:          false,
			test2:          false,
			test3:          false,
			test4:          false,
			test5:          false,
			test1notfound:  false,
			routernotfound: false,
		},
		"/test1/test4/asd/": {
			status:         http.StatusNotFound,
			test1:          true,
			test2:          false,
			test3:          false,
			test4:          true,
			test5:          false,
			test1notfound:  true,
			routernotfound: false,
		},
		"/test1/test4/asd/test": {
			status:         http.StatusNotFound,
			test1:          true,
			test2:          false,
			test3:          false,
			test4:          true,
			test5:          false,
			test1notfound:  true,
			routernotfound: false,
		},
		"/test1/test4/asd/test5": {
			status:         http.StatusOK,
			test1:          true,
			test2:          false,
			test3:          false,
			test4:          true,
			test5:          true,
			test1notfound:  false,
			routernotfound: false,
		},
		"/test1/test4/asd/test5/": {
			status:         http.StatusNotFound,
			test1:          true,
			test2:          false,
			test3:          false,
			test4:          true,
			test5:          false,
			test1notfound:  true,
			routernotfound: false,
		},
		"/test1/test4/asd/test5/asd": {
			status:         http.StatusNotFound,
			test1:          true,
			test2:          false,
			test3:          false,
			test4:          true,
			test5:          false,
			test1notfound:  true,
			routernotfound: false,
		},
	}

	router := New()

	test1Router := router.SubRouter("/test1")
	test1Router.Use(func(c *Context) Response {
		c.ResHeader().Set("test1", "1")
		return c.Next()
	})

	test2Router := test1Router.SubRouter("/test2")
	test2Router.Use(func(c *Context) Response {
		c.ResHeader().Set("test2", "1")
		return c.Next()
	})

	test2Router.Handle("/test3", func(c *Context) Response {
		c.ResHeader().Set("test3", "1")
		return textResponse{http.StatusOK, ""}
	})

	test4Router := test1Router.SubRouter("/test4/{id}")
	test4Router.Use(func(c *Context) Response {
		c.ResHeader().Set("test4", "1")
		return c.Next()
	})

	test4Router.Handle("/test5", func(c *Context) Response {
		c.ResHeader().Set("test5", "1")
		return textResponse{http.StatusOK, ""}
	})

	test1Router.Use(func(c *Context) Response {
		c.ResHeader().Set("test1-not-found", "1")
		return textResponse{http.StatusNotFound, ""}
	})

	router.Use(func(c *Context) Response {
		c.ResHeader().Set("router-not-found", "1")
		return textResponse{http.StatusNotFound, ""}
	})

	mux := router.CreateServeMux()

	assertHeader := func(res *http.Response, key string, exist bool) {
		if exist {
			if got, want := res.Header.Get(key), "1"; got != want {
				t.Errorf("header %s error. got:\"%v\" want:\"%v\"", key, got, want)
			}
		} else {
			if got, want := res.Header.Get(key), ""; got != want {
				t.Errorf("header %s error. got:\"%v\" want:\"%v\"", key, got, want)
			}
		}
	}

	for url, result := range testcases {
		name := fmt.Sprintf("Test nested router: %s", url)
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			res := w.Result()
			defer res.Body.Close()

			if got, want := res.StatusCode, result.status; got != want {
				t.Errorf("status error. got:\"%v\" want:\"%v\"", got, want)
			}
			assertHeader(res, "test1", result.test1)
			assertHeader(res, "test2", result.test2)
			assertHeader(res, "test3", result.test3)
			assertHeader(res, "test4", result.test4)
			assertHeader(res, "test5", result.test5)
			assertHeader(res, "test1-not-found", result.test1notfound)
			assertHeader(res, "router-not-found", result.routernotfound)
		})
	}
}
