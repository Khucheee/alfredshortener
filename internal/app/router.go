package app

import "github.com/go-chi/chi"

func (b *BaseController) Route() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", gzipMiddleware(b.WithLogging(b.solvePost)))
		r.Get("/{shorturl}", gzipMiddleware(b.WithLogging(b.solveGet)))
		r.Get("/ping", b.solvePing)

		r.Post("/api/shorten", gzipMiddleware(b.WithLogging(b.solveJSON)))
	})
	return r
}
