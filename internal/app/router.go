package app

import "github.com/go-chi/chi"

func (b *BaseController) Route() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", b.WithLogging(b.solvePost))
		r.Get("/{shorturl}", b.WithLogging(b.solveGet))
		r.Post("/api/shorten", b.WithLogging(b.solveJSON))
	})
	return r
}
