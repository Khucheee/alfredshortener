package main

import (
	"github.com/Khucheee/alfredshortener.git/internal/app"
	"github.com/go-chi/chi"
	"net/http"
)

func main() {

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/", app.SolvePost)
		r.Get("/{shorturl}", app.SolveGet)
	})

	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		panic(err)
	}
}
