package app

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSolvePost(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
		location    string
	}
	tests := []struct {
		name string
		want want
	}{{
		name: "Check post",
		want: want{
			code:        201,
			response:    "http://localhost:8080/4CWoMo83vssWiq4zcx51eCiTMVVH7yFaB1ft",
			contentType: "text/plain",
			location:    "",
		},
	},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://neal.fun/deep-sea/"))
			w := httptest.NewRecorder()
			SolvePost(w, request)

			res := w.Result()
			resBody, _ := io.ReadAll(res.Body)
			res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, test.want.location, res.Header.Get("Location"))
		})
	}

}
func TestSolveGet(t *testing.T) {
	type want struct {
		code     int
		location string
		path     string
	}
	tests := []struct {
		name string
		want want
	}{{name: "Check get positive",
		want: want{code: 307, location: "https://neal.fun/deep-sea/", path: "/4CWoMo83vssWiq4zcx51eCiTMVVH7yFaB1ft"}}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, test.want.path, nil)
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("shorturl", "4CWoMo83vssWiq4zcx51eCiTMVVH7yFaB1ft")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			SolveGet(w, req)
			res := w.Result()
			res.Body.Close()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.location, res.Header.Get("Location"))
		})
	}
}