package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gohf-http/gohf/v6"
)

func TestJsonResponse(t *testing.T) {
	createServerMux := func(data interface{}) *http.ServeMux {
		router := gohf.New()
		router.Use(func(c *gohf.Context) gohf.Response {
			return JSON(http.StatusOK, data)
		})
		return router.CreateServeMux()
	}

	getJsonResponse := func(mux *http.ServeMux, dest interface{}) error {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		mux.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()
		return json.NewDecoder(res.Body).Decode(dest)
	}

	t.Run("test json response: map", func(t *testing.T) {
		type bodyNested struct {
			Foo string `json:"foo"`
		}

		type body struct {
			Hello  string     `json:"hello"`
			Number int        `json:"number"`
			Nested bodyNested `json:"nested"`
		}

		mux := createServerMux(map[string]interface{}{
			"hello":  "gohf",
			"number": 123,
			"nested": map[string]interface{}{
				"foo": "bar",
			},
		})

		var dest body
		err := getJsonResponse(mux, &dest)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if got, want := dest.Hello, "gohf"; got != want {
			t.Errorf("got:%v want:%v", got, want)
		}
		if got, want := dest.Number, 123; got != want {
			t.Errorf("got:%v want:%v", got, want)
		}
		if got, want := dest.Nested.Foo, "bar"; got != want {
			t.Errorf("got:%v want:%v", got, want)
		}
	})

	t.Run("test json response: struct", func(t *testing.T) {
		type bodyNested struct {
			Foo string `json:"foo"`
		}

		type body struct {
			Hello  string     `json:"hello"`
			Number int        `json:"number"`
			Nested bodyNested `json:"nested"`
		}

		mux := createServerMux(body{
			Hello:  "GoHF!",
			Number: 1234,
			Nested: bodyNested{
				Foo: "Doe",
			},
		})

		var dest body
		err := getJsonResponse(mux, &dest)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if got, want := dest.Hello, "GoHF!"; got != want {
			t.Errorf("got:%v want:%v", got, want)
		}
		if got, want := dest.Number, 1234; got != want {
			t.Errorf("got:%v want:%v", got, want)
		}
		if got, want := dest.Nested.Foo, "Doe"; got != want {
			t.Errorf("got:%v want:%v", got, want)
		}
	})
}
