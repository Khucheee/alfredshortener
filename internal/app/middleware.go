package app

import (
	"net/http"
	"strings"
	"time"
)

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
			"storage", b.storage.Urls)
	}
	return logfn
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
