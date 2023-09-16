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
		cfg := Configure{"localhost:8080", "http://localhost:8080/", "", ""}
		keepe := NewKeeper(cfg.FilePath)
		str := Storage{make(map[string]string), keepe}
		log := Logger{}
		log.CreateSuggarLogger()
		controller := NewBaseController(cfg, str, log)
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://neal.fun/deep-sea/"))
			w := httptest.NewRecorder()
			controller.solvePost(w, request)

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
		body string
		want want
	}{
		{
			name: "Check get positive",
			body: "https://neal.fun/deep-sea/",
			want: want{
				code:     307,
				location: "https://neal.fun/deep-sea/",
				path:     "/4CWoMo83vssWiq4zcx51eCiTMVVH7yFaB1ft",
			},
		}, {
			name: "Check get negative",
			body: "https://neal.fun/deep-sea/",
			want: want{
				code:     400,
				location: "",
				path:     "/4CWoMo83vssWiq4zcx51eCiTMVVH7yFaB1f",
			},
		},
	}

	for _, test := range tests {
		cfg := Configure{"localhost:8080", "http://localhost:8080/", "", ""}
		keepe := NewKeeper(cfg.FilePath)
		str := Storage{make(map[string]string), keepe}
		log := Logger{}
		controller := NewBaseController(cfg, str, log)
		t.Run(test.name, func(t *testing.T) {

			req1 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.body))
			w1 := httptest.NewRecorder()
			controller.solvePost(w1, req1)
			req2 := httptest.NewRequest(http.MethodGet, test.want.path, nil)
			w2 := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("shorturl", test.want.path[1:])
			req2 = req2.WithContext(context.WithValue(req2.Context(), chi.RouteCtxKey, rctx))

			controller.solveGet(w2, req2)
			res := w2.Result()
			res.Body.Close()
			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.location, res.Header.Get("Location"))
		})
	}
}
