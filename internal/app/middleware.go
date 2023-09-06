package app

import (
	"net/http"
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
