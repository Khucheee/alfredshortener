package app

import "github.com/go-chi/chi"

func (b *BaseController) Route() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Use(CookieMiddleware)
		r.Use(b.WithLogging)
		r.Use(gzipMiddleware)

		r.Post("/", b.solvePost)
		r.Get("/{shorturl}", b.solveGet)
		r.Get("/ping", b.solvePing)
		r.Post("/api/shorten", b.solveJSON)
		r.Post("/api/shorten/batch", b.solveBatch)
		r.Get("/api/user/urls", b.solveUserLinks)
		r.Delete("/api/user/urls", b.DeleteUserLinks)
	})
	return r
}
