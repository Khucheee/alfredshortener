package app

import (
	"bytes"
	"encoding/json"
	"github.com/btcsuite/btcutil/base58"
	"github.com/go-chi/chi"
	"io"
	"net/http"
	"strings"
	"time"
)

type BaseController struct {
	config Configure
	urls   map[string]string //мапа содержит сокращенный урл и полный
	logger Logger
}
type Jsonquery struct {
	URL string
}
type Jsonresponse struct {
	Response string `json:"result"`
}

func (b *BaseController) addURL(shorturl, url string) { //добавляем значение в мапу
	b.urls[shorturl] = url
}

func (b *BaseController) searchURL(shorturl string) string { //ищем значение в мапе, если ""то не нашли
	url := b.urls[shorturl]
	return url
}

func NewBaseController(c Configure, l Logger) *BaseController {
	return &BaseController{
		config: c,
		urls:   make(map[string]string),
		logger: l,
	}
}

func (b *BaseController) Route() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", gzipMiddleware(b.WithLogging(b.solvePost)))
		r.Get("/{shorturl}", gzipMiddleware(b.WithLogging(b.solveGet)))
		r.Post("/api/shorten", gzipMiddleware(b.WithLogging(b.solveJSON)))
	})
	return r
}

func (b *BaseController) WithLogging(h http.HandlerFunc) http.HandlerFunc {
	logfn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		responseData := &responseData{status: 0, size: 0}
		lw := loggingResponseWriter{w, responseData}
		h.ServeHTTP(&lw, r)
		duration := time.Since(start)
		b.logger.sugar.Infoln(
			"URI", r.RequestURI,
			"duration", duration,
			"method", r.Method,
			"status", responseData.status,
			"size", responseData.size,
			"storage", b.urls)
	}
	return logfn
}

func (b *BaseController) solvePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	reqBody, _ := io.ReadAll(r.Body)
	reqBodyEncoded := base58.Encode(reqBody)
	defer r.Body.Close()
	respBody := b.config.Address + reqBodyEncoded
	w.Write([]byte(respBody))
	b.addURL(reqBodyEncoded, string(reqBody))
}

func (b *BaseController) solveGet(w http.ResponseWriter, r *http.Request) {
	if b.searchURL(chi.URLParam(r, "shorturl")) == "" { //если ключ в мапе пустой, то 400
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", b.urls[chi.URLParam(r, "shorturl")]) //если дошли до сюда, то в location суем значение из мапы по ключу
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (b *BaseController) solveJSON(w http.ResponseWriter, r *http.Request) {
	var jsonquery Jsonquery
	var jsonresponse Jsonresponse
	var buf bytes.Buffer
	var shorturl string
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &jsonquery); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shorturl = base58.Encode([]byte(jsonquery.URL))
	jsonresponse.Response = b.config.Address + shorturl
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp, _ := json.Marshal(jsonresponse)
	w.Write(resp)
	b.addURL(shorturl, jsonquery.URL)
}

func gzipMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ow := w
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer cw.Close()
		}
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}
		h.ServeHTTP(ow, r)
	}
}
