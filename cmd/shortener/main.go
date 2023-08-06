package main

import (
	"github.com/Khucheee/alfredshortener.git/internal/app"
	"github.com/go-chi/chi"
	"net/http"
)

func main() {
	app.SetFlags()
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", app.SolvePost)
		r.Get("/{shorturl}", app.SolveGet)
	})
	err := http.ListenAndServe(*app.Host, r)
	if err != nil {
		panic(err)
	}
}
